// File: "ratelimit.go"

package xlog

import (
	"context"
	"log/slog" // go>=1.21
	"sync"
	"time"
	// FIXME: "golang.org/x/exp/slog" // экспериментальный пакет для go=1.20 только
)

// Параметры ограничителя (Rate Limiter'а) по умолчанию
const (
	DEFAULT_RATE_LIMIT_MAX_NUM     = 100  // максимальное число однотипных сообщений [шт]
	DEFAULT_RATE_LIMIT_INTERVAL_MS = 1000 // за заданный период времени [мс]

	// Период проверки наличия готовых сгруппированных сообщений [мс]
	DEFAULT_RATE_LIMIT_FLUSH_PERIOD_MS = 3000
)

// Дополнительные атрибуты для записей в журнале,
// сгруппрованных с помощью Rate Limiter'а
const (
	// Ключ для числа сообщений объединенных в одно в случае срабатывания
	// ограничений RateLimiter'а.
	// Передается значение больше 1.
	RepeatedKey = "repeated"

	// Ключ для булевого атрибута, передается, если сообщение было задержано
	// но не было объединения атрибутов с другими однотипными сообщениями.
	DelayedKey = "delayed"

	// Ключ для булевого атрибута, который свидетельствует о том,
	// что выброс задержанных сообщений в журнал был в
	// соответствии с заданным значением FlushPeriodMs
	FlushedKey = "flushed"
)

// rateLimitTopic - идентификатор (топик) для группировки однотипных сообщений
type rateLimitTopic struct {
	level slog.Level // уровень логитопикия
	msg   string     // тело сообщения
}

// rateLimitItem - структура накопления статистики для каждого класса сообщений
type rateLimitItem struct {
	// Кольцевой буфер с временными метки последних
	// сообщений данного класса (Unix time, мс) за период IntervalMs [мс]
	// из структуры RateLimitConf.
	// Размер слайса - MaxNum из структуры RateLimitConf.
	ring []int64

	next int // индекс для следующей метки в кольцевом буфере
	old  int // индекс самой "старой" метки в кольцевом буфере
	cnt  int // число временных меток в кольцевом буфере

	// Время последнего пропущенного сообщения
	time time.Time

	// Объединенные атрибуты для пропущенных сообщений
	attrs map[string]slog.Value

	// Program counter крайнего пропущенного сообщения
	pc uintptr

	// Счётчик пропущенных (сгруппированных) сообщений
	repeated int
}

// Структура с переменными состояния RateLimiter'а для Middleware
type rateLimit struct {
	// Максимальное число сообщение с однотипным level/msg допустимое
	// за заданное время
	maxNum int

	// Временной интервал (размер временного скользящего окна) в течении
	// которого однотипные сообщение до maxNum штук не группируются,
	// задается в миллисекундах
	intervalMs int

	// Период проверки и выброса в журнал сгруппированных сообщений,
	// необходимый для того, чтобы сгруппированные сообщения своевременно
	// попадали в журнал в условиях, когда других сообщений нет,
	// задается в миллисекундах.
	flushPeriodMs int

	// Карта для каждого класса (типа) сообщений
	items map[rateLimitTopic]*rateLimitItem

	// Последняя функция выдачи сообщений в журнал и контекст для неё
	handle HandleFunc
	ctx    context.Context

	// Признак того, что в данный момент есть отложенные сообщения
	// в множестве items
	limited bool

	// Признак того, что горутина goWait запущена и
	// повторный запуск не требуется
	wait bool

	// Мьютекс для защиты items/handle/ctx/wait
	mx sync.Mutex
}

// newRateLimit инициализирует структуру состояния RateLimiter'а
func newRateLimit(conf RateLimitConf) *rateLimit {
	rl := &rateLimit{
		maxNum:        conf.MaxNum,
		intervalMs:    conf.IntervalMs,
		flushPeriodMs: conf.FlushPeriodMs,
		items:         map[rateLimitTopic]*rateLimitItem{},
	}

	// Если не заданы некоторые параметры использовать умолчания
	if rl.maxNum <= 0 {
		rl.maxNum = DEFAULT_RATE_LIMIT_MAX_NUM
	}
	if rl.intervalMs <= 0 {
		rl.intervalMs = DEFAULT_RATE_LIMIT_INTERVAL_MS
	}
	if rl.flushPeriodMs == 0 {
		rl.flushPeriodMs = DEFAULT_RATE_LIMIT_FLUSH_PERIOD_MS
	}

	return rl
}

// goWait - горутина отложнной проверки наличия готовых к
// выдаче сгруппированных сообщений
func (rl *rateLimit) goWait() {
	for {
		time.Sleep(time.Duration(rl.flushPeriodMs) * time.Millisecond) // задержка
		rl.mx.Lock()
		now := time.Now().UnixMilli()
		rl.check(now, true) // flushed=true
		if !rl.limited {    // нет сгруппированных сообщений, завершить горутину
			rl.wait = false
			rl.mx.Unlock()
			return
		}
		rl.mx.Unlock()
	} // for
}

// startGoWait запускает горутину goWait, но не более одной
func (rl *rateLimit) startGoWait() {
	if !rl.wait && rl.limited { // не запущена и требуется
		rl.wait = true
		go rl.goWait()
	}
}

// check производит анализ всех накопленных ранее временных меток и
// выдает в журнал сгруппированные сообщения при необходимости.
//
//	now     - текущее Unix время [мс]
//	flushed - признак задерженного запуска из горутины goWait (для отладки)
//
// Функция возвращает ошибку, если при выводе сгруппированного сообщения
// rt.handle вернул ошибку.
//
// При необходимоси с помощью функции startGoWait запускается
// горутина (1 шт) goWait.
func (rl *rateLimit) check(now int64, flushed bool) error {
	rl.limited = false
	for topic, it := range rl.items {
		for now-it.ring[it.old] > int64(rl.intervalMs) {
			if it.cnt == rl.maxNum && it.repeated != 0 {
				// Вывести в журнал сгруппированное сообщение
				rec := slog.NewRecord(it.time, topic.level, topic.msg, it.pc)
				for key, value := range it.attrs {
					rec.AddAttrs(slog.Attr{
						Key:   key,
						Value: value,
					})
				}
				if it.repeated > 1 {
					rec.AddAttrs(slog.Attr{
						Key:   RepeatedKey,
						Value: slog.IntValue(it.repeated),
					})
				} else {
					rec.AddAttrs(slog.Attr{
						Key:   DelayedKey,
						Value: slog.BoolValue(true),
					})
				}
				if flushed {
					rec.AddAttrs(slog.Attr{
						Key:   FlushedKey,
						Value: slog.BoolValue(true),
					})
				}
				if rl.handle == nil { // параноидальная проверка
					return ErrNilHandler
				}
				if err := rl.handle(rl.ctx, rec); err != nil {
					return err
				}
				it.attrs = nil
				it.repeated = 0
			}

			// Удалить метку из кольцевого буфера
			it.old = (it.old + 1) % rl.maxNum
			it.cnt--
			if it.cnt == 0 { // не хранить в ОЗУ записи с пустым кольцевым буфером
				delete(rl.items, topic)
				break
			}
		} // for

		// Если хотя бы в одном топике остались сгруппированные сообщений,
		// то необходимо поддерживать горутину goWait в запущенном состоянии
		rl.limited = rl.limited || (it.cnt == rl.maxNum)
	} // for

	// Запустить горутину отложенного вывода при необходимости
	rl.startGoWait()
	return nil
}

// mwFunc - метод реализации middleware для RateLimiter'а
//
//	ctx    - контекст для передачи в handle
//	r      - запись для журнала
//	handle - функция для выдачи записи в журнал
func (rl *rateLimit) mwFunc(ctx context.Context, r slog.Record, handle HandleFunc) error {
	now := r.Time.UnixMilli() // текущее Unix время сообщения [мс]

	rl.mx.Lock()
	defer rl.mx.Unlock()

	// Сохранить handle и контекст для него (используется в методе check)
	rl.handle = handle
	rl.ctx = ctx

	// Провести анализ всех накопленных ранее временных меток,
	// выдать в журнал сгруппированные сообщения при необходимости
	err := rl.check(now, false) // flushed=false
	if err != nil {
		return err
	}

	// Топик данного сообщения (ключ для rl.items)
	topic := rateLimitTopic{
		level: r.Level,   // уровень сообщения
		msg:   r.Message, // текст сообщения
	}

	// Проверить были ли ранее в недалеком прошлом сообщения данного класса
	it, ok := rl.items[topic]
	if !ok { // ранее подобных сообщений не отмечено
		it := &rateLimitItem{
			ring: make([]int64, rl.maxNum), // выделить память под кольцевой буфер
		}

		// Сохранить первую отметку в кольцевом буфере
		it.ring[0] = now
		it.next = 1
		it.cnt = 1

		// Сохранить новый топик в карте
		rl.items[topic] = it

		return handle(ctx, r) // вывести сообщение в журнал без прореживания
	}

	// Подобные сообщения данного класса были в недалеком прошлом
	if it.cnt < rl.maxNum { // не достигнут предел в MaxNum сообщений
		// Сохранить новую отметку в кольцевом буфере
		it.ring[it.next] = now
		it.next = (it.next + 1) % rl.maxNum
		it.cnt++
		return handle(ctx, r) // вывести сообщение в журнал без прореживания
	}

	// Достигнут предел в MaxNum сообщений за интервал IntervalMs.
	// Необходимо накапливать атрибуты с целью прореживания сообщений.
	it.time = r.Time
	if it.repeated == 0 {
		it.attrs = map[string]slog.Value{}
	}
	r.Attrs(func(attr slog.Attr) bool {
		it.attrs[attr.Key] = attr.Value
		return true
	})
	it.pc = r.PC
	it.repeated++

	// Запустить горутину отложенного вывода при необходимости
	rl.startGoWait()

	return nil // пропустить вывод сообщения в журнал сейчас
}

// NewMiddlewareRateLimit возвращает Middleware для реализации
// ограничителя большого числа повторяющихся в журнале сообщений
// т.н. RateLimiter'а.
//
// Определение повторов сообщений производится по уровню логированию (level)
// и тексту сообщения (msg).
//
//	conf - конфигурация RateLimeter'а
func NewMiddlewareRateLimit(conf RateLimitConf) Middleware {
	if conf.Disable { // Rate Limiter отключен конфигурацией
		return func(hf HandleFunc) HandleFunc {
			return hf // вернуть функцию хендлера без изменений
		}
	}

	// Инициализировать структуру состояния RateLimiter'а
	rl := newRateLimit(conf)

	// Вернуть middleware
	return NewMiddleware(rl.mwFunc)
}

// EOF: "ratelimit.go"
