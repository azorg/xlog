// File: "idhandler.go"

package xlog

import (
	"context"
	"fmt"
	"log/slog" // go>=1.21
	//"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/gofrs/uuid"
	// FIXME: "golang.org/x/exp/slog" // экспериментальный пакет для go=1.20 только
)

// Дополнительные атрибуты для каждой записи в журнале
const (
	// Ключ идентификации горутины (если GoId=true)
	GoKey = "goroutine"

	// Ключ UUID идентификатора записи в журнале (если LogId=true)
	IdKey = "logId"

	// Ключ контрольной суммы в журнале (если SumAlone=true)
	SumKey = "logSum"
)

// Структура конфигурации для IdHandler'а
type IdOptions struct {
	// Добавить в журнал идентификатор горутины ("goroutine")
	GoId bool `json:"goId"`

	// Добавлять UUID идентификатор к каждой записи в журнале ("logId")
	LogId bool `json:"logId"`

	// Добавить подсчёт контрольной суммы (КС)
	AddSum bool `json:"addSum"`

	// Вычислять контрольную сумму по всем атрибутам рекурсивно
	SumFull bool `json:"sumFull"`

	// Включить в расчет контрольной суммы метку времени
	SumTime bool `json:"sumTime"`

	// Подсчет контрольной суммы с учётом предыдущей записи
	SumChain bool `json:"sumChain"`

	// Не упаковать КС в последний байт UUID, а добавить ключ "LogSum"
	SumAlone bool `json:"sumAlone"`
}

// Структура безопасного хранения контрольной суммы
type idSum struct {
	val uint16     // значение CRC16
	mx  sync.Mutex // мьютекс для безопасного совместного доступа к val
}

// IdHandler - это обертка заданного slog.Handler'а для возможности
// обогащения журнала дополнительными атрибутами (goroutine, logId, logSum).
// Кроме того, IdHandler поддерживает Middleware для метода Handle
// интерфейса slog.Handler.
type IdHandler struct {
	handler slog.Handler  // исходный (оборачиваемый) хендлер
	opts    *IdOptions    // заданные опции для всей цепочки
	sum     *idSum        // контрольная сумма предыдущей записи
	withSum uint16        // контрольная сумма "With" атрибутов
	valuers []slog.Attr   // корневые атрибуты содержащие slog.LogValuer'ы
	groups  []string      // цепочка открытых групп
	attrs   [][]slog.Attr // атрибуты открытых групп
	mx      sync.Mutex    // мьютекс для безопасного совместного доступа к groups/attrs
	mws     []Middleware  // обёртки для метода Handle
}

// Убедиться, что *IdHandler реализует интерфейс slog.Handler
var _ slog.Handler = (*IdHandler)(nil)

// Создать новый Id-хендлер обертку для slog.Handler'а,
// который позволяет добавлять в журнал дополнительные атрибуты
// "goroutine", "logId" (на основе UUID) и "logSum" (при необходимости),
// а также позволяет оборачивать метод Handle() заданого handler'а
// произвольной цепочной Middleware.
//
//	handler - оборачиваемый slog handler
//	opts - опции Id-хендлера или nil
//	sum - начальное значение контрольной суммы (обычно 0)
//	mws - цепочка Middleware для оборачивания метода Hanlde()
func NewIdHandler(
	handler slog.Handler, opts *IdOptions, sum uint16,
	mws ...Middleware,
) *IdHandler {
	// Optimization: avoid chains of idHandlers
	if ih, ok := handler.(*IdHandler); ok {
		handler = ih.handler
	}
	h := &IdHandler{
		handler: handler,
		opts:    &IdOptions{},
		sum:     &idSum{val: sum},
		withSum: uint16(0),
		valuers: make([]slog.Attr, 0),
		groups:  make([]string, 0),
		attrs:   make([][]slog.Attr, 0),
		mws:     mws,
	}
	if opts != nil {
		h.opts = opts
	}
	return h
}

// Метод Enabled() реализует интерфейс slog.Handler
func (h *IdHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// addIdAndSum обогащает запись журнала дополнительными атрибутами
// (goroutine, logId, logSum)
func (h *IdHandler) addIdAndSum(r *slog.Record) {
	if h.opts.GoId { // добавить в журнал goroutine
		var buf [64]byte // FIXME: magic size
		n := runtime.Stack(buf[:], false)
		idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
		goroutine, err := strconv.Atoi(idField)
		if err == nil {
			r.AddAttrs(slog.Int(GoKey, goroutine))
		}
	}

	var logSum uint16
	if h.opts.LogId { // добавить в журнал logId
		logId, _ := uuid.NewV7()
		if h.opts.AddSum {
			logSum = h.sum.val ^ Checksum(h.withSum,
				h.opts.SumFull, h.opts.SumTime, *r, logId)
			if h.opts.SumChain {
				h.sum.val = logSum
			}

			if h.opts.SumAlone { // добавить в журнал logId и logSum
				r.AddAttrs(
					slog.String(IdKey, logId.String()),
					slog.String(SumKey, fmt.Sprintf("%04x", logSum)))
			} else { // добавить в журнал только logId с logSum внутри
				logId[14] = byte((logSum >> 8) & 0xFF)
				logId[15] = byte(logSum & 0xFF)
				r.AddAttrs(slog.String(IdKey, logId.String()))
			}
		} else { // добаить в журнал только logId без logSum
			r.AddAttrs(slog.String(IdKey, logId.String()))
		}
	} else if h.opts.AddSum {
		// Добавить в журнал только logSum
		logSum = h.sum.val ^ Checksum(h.withSum,
			h.opts.SumFull, h.opts.SumTime, *r, uuid.UUID{})
		if h.opts.SumChain {
			h.sum.val = logSum
		}
		r.AddAttrs(slog.String(SumKey, fmt.Sprintf("%04x", logSum)))
	}
}

// Добавить атрибуты в последнюю открытую группу.
// Перед этим копировать списки открытых групп, а для
// последней открытой группы создать копию списка атрибутов.
// Весьма магическая функция, требующая осмысления.
func attrsAdd(as [][]slog.Attr, attrs []slog.Attr) [][]slog.Attr {
	bs := make([][]slog.Attr, 0, len(as))
	ix := len(as) - 1
	for i := 0; i < ix; i++ {
		bs = append(bs, as[i])
	}
	b := make([]slog.Attr, 0, len(as[ix])+len(attrs))
	b = append(b, as[ix]...)
	b = append(b, attrs...)
	bs = append(bs, b)
	return bs
}

// Создать атрибут с групповым значением из всех открытых групп
func groupAttr(groups []string, attrs [][]slog.Attr) slog.Attr {
	for ix := len(attrs) - 1; ix > 0; ix-- {
		i := ix - 1
		attrs[i] = append(attrs[i], slog.Attr{
			Key:   groups[ix],
			Value: slog.GroupValue(attrs[ix]...),
		})
	}

	return slog.Attr{
		Key:   groups[0],                    // ключ
		Value: slog.GroupValue(attrs[0]...), // групповое значение
	}
}

// Применить цепочку middleware, если len(h.mws) != 0.
// Последним в цепочке может быть middleware для обогащения
// журанал дополнительными атрибутами (goroutine, logId, logSum).
func (h *IdHandler) middleware(ctx context.Context, r slog.Record) error {
	if len(h.mws) == 0 { // нет middleware
		return h.handler.Handle(ctx, r)
	}

	// Использовать цепочку middleware
	handle := h.mws[len(h.mws)-1](h.handler.Handle)
	for i := len(h.mws) - 2; i >= 0; i-- {
		handle = h.mws[i](handle)
	}
	return handle(ctx, r)
}

// Метод Handle() реализует интерфейс slog.Handler
func (h *IdHandler) Handle(ctx context.Context, r slog.Record) error {
	if h.opts.SumChain {
		// Захватить мьютек доступа к общей контрольной суммы до завершения вывода
		h.sum.mx.Lock()
		defer h.sum.mx.Unlock()
	}

	h.mx.Lock()
	defer h.mx.Unlock()

	if len(h.groups) == 0 && len(h.valuers) == 0 {
		// Нет открытых групп, нет slog.LogValuer'ов.
		// Обогатить существующую запись требуемыми полями
		// (goroutine, logId, logSum).
		h.addIdAndSum(&r)

		// Обработать цепочку middleware
		return h.middleware(ctx, r)
	}

	// Создать новую запись
	rNew := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)

	// Добавить в начало записи записи с slog.LogValuer'ами
	if len(h.valuers) != 0 {
		rNew.AddAttrs(h.valuers...)
	}

	if len(h.groups) == 0 /*&& len(h.valuers) != 0*/ {
		// Добавить атрибуты из старой записи (r -> rNew)
		r.Attrs(func(attr slog.Attr) bool {
			rNew.AddAttrs(attr)
			return true
		})

		// Обогатить новую запись требуемыми полями (goroutine, logId, logSum)
		h.addIdAndSum(&rNew)

		// Обработать цепочку middleware
		return h.middleware(ctx, rNew)
	}

	// Получить список атрибутов из старой записи (r -> attrs)
	attrs := make([]slog.Attr, 0, r.NumAttrs())
	r.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, attr)
		return true
	})

	// Добавить атрибуты из старой записи в атрибуты открытых групп
	as := attrsAdd(h.attrs, attrs)

	// Создать групповое значение из открытых групп
	grpAttr := groupAttr(h.groups, as) // slog.Attr

	// Заполнить новую запись обёрнутыми в группы атрибутами
	rNew.AddAttrs(grpAttr)

	// Обогатить новую запись требуемыми полями (goroutine, logId, logSum)
	h.addIdAndSum(&rNew)

	// Обработать цепочку middleware
	return h.middleware(ctx, rNew)
}

// Метод WithAttrs() реализует интерфейс slog.Handler
func (h *IdHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.mx.Lock()
	defer h.mx.Unlock()

	// valuers - признак того, что в списке атрибутов найдено
	// значение с "отложенным" вычислением (slog.LogValuer или FieldsProvider)
	valuers := false
	for i := range attrs {
		value := attrs[i].Value.Any()
		switch value.(type) {
		case slog.LogValuer, FieldsProvider:
			valuers = true
			break
		} // switch
	} // for

	if len(h.groups) == 0 && len(h.valuers) == 0 && !valuers {
		// Нет открытых групп, нет slog.Valuer'ов
		withSum := h.withSum
		if h.opts.SumFull {
			for _, attr := range attrs {
				withSum ^= ChecksumAttrSlog(attr.Key, attr.Value)
			}
		}

		return &IdHandler{
			handler: h.handler.WithAttrs(attrs),
			opts:    h.opts,
			sum:     h.sum,
			withSum: withSum,
			valuers: h.valuers,
			groups:  h.groups,
			attrs:   h.attrs,
			mws:     h.mws,
		}
	}

	vs := h.valuers
	as := h.attrs

	if valuers || len(h.valuers) != 0 {
		vs = append(vs, attrs...)
	} else { // len(h.groups) != 0
		// Добавить атрибуты в последнюю открытую группу
		as = attrsAdd(as, attrs)
	}

	return &IdHandler{
		handler: h.handler,
		opts:    h.opts,
		sum:     h.sum,
		withSum: h.withSum,
		valuers: vs,
		groups:  h.groups,
		attrs:   as,
		mws:     h.mws,
	}
}

// Метод WithGroup() реализует интерфейс slog.Handler
func (h *IdHandler) WithGroup(name string) slog.Handler {
	if name == "" { // если группа пустая - ничего не делать
		return h
	}

	h.mx.Lock()
	defer h.mx.Unlock()

	// Открыть новую группу (добавить пустой слайс атрибутов)
	return &IdHandler{
		handler: h.handler,
		opts:    h.opts,
		sum:     h.sum,
		withSum: h.withSum,
		groups:  append(h.groups, name),
		attrs:   append(h.attrs, []slog.Attr{}),
		mws:     h.mws,
	}
}

// EOF: "idhandler.go"
