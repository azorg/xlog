// File: "checksum.go"

package xlog

import (
	"encoding/json"
	"fmt"
	"log/slog" // go>=1.21
	"reflect"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	"github.com/sigurn/crc16"
	// FIXME: "golang.org/x/exp/slog" // экспериментальный пакет для go=1.20 только
)

const (
	// Метка времени в JSON/text журнале
	TimeKey = "time"

	// Метка уровня в JSON/text журнале
	LevelKey = "level"

	// Метка сообщения JSON/logfmt журнале
	MsgKey = "msg"

	// Метка исходных текстов в JSON журнале
	SourceKey = "source"
)

// ChecksumRes - это результат проверки контрольной суммы JSON записи.
// Пользователь может сверить поля LogSum и Sum.
// Заполняется по результатам выполнения функции ChecksumVerify().
type ChecksumRes struct {
	Time      time.Time      // метка времени записи
	Level     slog.Level     // уровень сообщения
	Source    map[string]any // ссылка на исходные тексты (если есть)
	Message   string         // сообщение журнала
	Goroutine int            // идентификатор горутины (если есть)
	LogId     uuid.UUID      // идентификатор записи в журнале
	Err       string         // ошибка в сообщении с ключом "err"
	LogSum    uint16         // контрольная сумма извлеченная их журнала
	Sum       uint16         // контрольная сумма вычисленная
}

// Используемая таблица для вычисления CRC16
// [https://pkg.go.dev/github.com/sigurn/crc16#pkg-constants]
// [https://reveng.sourceforge.io/crc-catalogue/16.htm]
var crcTable = crc16.MakeTable(crc16.CRC16_XMODEM)

// Checksum вычисляет контрольную сумму (на основе CRC16) записи в журнале.
// Используется один из двух алгоритмов (full=true/false).
//
// На входе:
//
//	full - признак для вычисления контрольной суммы по всем атрибутам рекурсивно
//	sum - значение контрольной суммы предыдущей записи или 0
//	timeOn - включить в расчет CRC метку времени
//	r - подготовленная для выдачи в slog-журнал запись
//	logId - UUID записи
func Checksum(
	sum uint16, full, timeOn bool, r slog.Record, logId uuid.UUID) uint16 {

	if !full {
		return ChecksumSimple(sum, timeOn, r, logId)
	} else {
		return ChecksumFull(sum, timeOn, r, logId)
	}
}

// ChecksumSimple вычисляет контрольную сумму (CRC16) записи в журнале.
// Контрльная сумма включает в себя:
//
//   - time (RFC3339Milli или RFC3339Micro как в журнале)
//   - level (как строка типа ERROR, WARN и т.п.)
//   - msg (как строка)
//   - logId (старшие 14 байт)
//   - err, если есть (как строку)
//
// На входе:
//
//	sum - значение контрольной суммы предыдущей записи или 0
//	timeOn - включить в расчет CRC метку времени
//	r - подготовленная для выдачи в slog-журнал запись
//	logId - UUID записи
func ChecksumSimple(
	sum uint16, timeOn bool, r slog.Record, logId uuid.UUID) uint16 {

	// Сформировать буфер "sync pool" для сбора данны для CRC
	buf := newBuffer()
	defer buf.Free()

	if timeOn {
		// Учесть в контрольной сумме метку времени как строку в формате RFC3339Milli
		buf.WriteString(r.Time.UTC().Format(RFC3339Milli))
	}

	// Учесть в CRC уровень журналирования как строку
	buf.WriteString(LevelToLabel(r.Level))

	// Учесть в CRC текст сообщения
	buf.WriteString(r.Message)

	if !logId.IsNil() {
		// Учесть в CRC первые (старшие) 14 байт UUID идентификатора
		*buf = append(*buf, logId[0:14]...)
	}

	// Учесть к CRC ошибки "err", если они есть в атрибутах
	//r.Attrs(func(attr slog.Attr) bool {
	//	if attr.Key == ErrKey {
	//		if err, ok := attr.Value.Any().(error); err != nil && ok {
	//			buf.WriteString(err.Error())
	//		}
	//	}
	//	return true
	//})

	// Вычислить CRC16
	// Учесть предыдущее значение CRC
	return sum ^ crc16.Checksum(*buf, crcTable)
}

// ChecksumFull вычисляет контрольную сумму записи в журнале.
// Контрольная сумма вычисляется специальным образом по всем атрибутам JSON.
//
// Контрольные суммы для последовательности пар ключ/значение
// складываются по модулю 2 (XOR), таким образом обеспечивается
// инвариантность контрольную суммы в случае перестановки атрибутов
// при обработке журналов разными фильтрами.
// Контрольная сумма для каждой пары ключ/значение формируется путем простого
// сложения в дополнительном коде CRC16 сумм.
// Перед вычислением CRC16 каждое значением приводится к строке.
//
// Контрольная сумма включает в себя:
//
//   - временную метку в RFC3339Milli, если timeOn=true (time=<string>);
//   - уровень сообщения журнала (level=<int>)
//   - текст сообщения r.Message (msg=<string>)
//   - "рекурсивную сумму" всех атрибутов не зависящую от любых перестановок
//   - первые (старшие) 14 байт UUID идентификатора (logId=<bytes>)
//
// На входе:
//
//	sum - начальное значение контрольной суммы или 0
//	timeOn - включить в расчет суммы метку времени
//	r - подготовленная для выдачи в slog-журнал запись
//	logId - UUID записи
func ChecksumFull(
	sum uint16, timeOn bool, r slog.Record, logId uuid.UUID) uint16 {

	if timeOn {
		// Учесть в контрольной сумме метку времени как строку в формате RFC3339Milli
		sum ^= ChecksumAttr(TimeKey, r.Time.UTC().Format(RFC3339Milli))
	}

	// Учесть в КС уровень сообщения
	sum ^= ChecksumAttr(LevelKey, r.Level.Level())

	// Учесть в КС текст сообщения
	sum ^= ChecksumAttr(MsgKey, r.Message)

	// Добавить в контрольную сумму суммы по всем атрибутам
	r.Attrs(func(attr slog.Attr) bool {
		sum ^= ChecksumAttrSlog(attr.Key, attr.Value)
		return true
	})

	if !logId.IsNil() {
		// Учесть в КС первые (старшие) 15 байт UUID идентификатора
		sum ^= ChecksumAttr(IdKey, logId[:14])
	}

	return sum
}

// ChecksumAttrSlog - вычисляет контрольную сумму записи для одного
// атрибута slog key/value.
// Анализируется тип slog значения и если тип, не стандартный
// (см. slogg.Kind) применяется рефлексия (таким образом рассчитываем
// немного поднять производительность).
// Контрольная сумма вычисляется рекурсивно для всех вложенных структур.
// Контрольные суммы смежных атрибутов складываются по модулю 2 (XOR).
// Контрольные суммы key и value складываются по правилу сложения
// в дополнительном коде.
// Функция принимает key и slog.Value.
// Функция корректно обрабатывает slog группы.
func ChecksumAttrSlog(key string, value slog.Value) uint16 {
	var str string
	switch value.Kind() {
	case slog.KindString:
		str = value.String()

	case slog.KindBool:
		str = strconv.FormatBool(value.Bool())

	case slog.KindTime:
		str = value.Time().Format(time.RFC3339Nano)

	case slog.KindDuration:
		str = strconv.FormatInt(int64(value.Duration()), 10)

	case slog.KindInt64:
		str = strconv.FormatInt(value.Int64(), 10)

	case slog.KindUint64:
		str = strconv.FormatUint(value.Uint64(), 10)

	case slog.KindFloat64:
		//str = strconv.FormatFloat(value.Float64(), 'g', -1, 64) FIXME: так плохо!
		bytes, _ := json.Marshal(value.Float64())
		str = string(bytes)

	case slog.KindGroup:
		sum := uint16(0)
		for _, attr := range value.Group() {
			//fmt.Fprintf(os.Stderr, "group=%s key=%s value=%s sum=%04X\n",
			//  key, attr.Key, Sprint(attr.Value.Any()),
			//  ChecksumAttrSlog(attr.Key, attr.Value))
			sum ^= ChecksumAttrSlog(attr.Key, attr.Value)
		}
		return sum + crc16.Checksum([]byte(key), crcTable)

	case slog.KindLogValuer:
		// FIXME: Тут может быть проблема (!)
		// В журнал может пойти более позднее значение, чем то,
		// которое использовалось при вычислении контрольной суммы.
		return ChecksumAttrSlog(key, value.LogValuer().LogValue())

	default: // slog.KindAny
		return ChecksumAttr(key, value.Any())
	}

	sum := crc16.Checksum([]byte(key), crcTable)
	sum += crc16.Checksum([]byte(str), crcTable)
	//fmt.Fprintf(os.Stderr, ">>> key=%s str=%s sum=%04X\n", key, str, sum) //!!!
	return sum
}

// ChecksumAttr вычисляет контрольную сумму записи для одного
// произвольного атрибута key/value.
// Контрольная сумма вычисляется рекурсивно для всех вложенных структур
// с использованием рефлексии.
// Контрольные суммы смежных атрибутов складываются по модулю 2 (XOR).
// Контрольные суммы key и value складываются по правилу сложения
// в дополнительном коде.
func ChecksumAttr(key string, value any) uint16 {
	if value == nil { // защита от nil
		return ChecksumAttr(key, "<nil>")
	}

	switch v := value.(type) {
	case []byte: // для последовательности байт вычислить КС как CRC8
		sum := crc16.Checksum([]byte(key), crcTable)
		sum += crc16.Checksum(v, crcTable)
		return sum

	case error: // go-ошибку представить в виде строки
		return ChecksumAttr(key, v.Error())

	case time.Time: // время представить как строку в RFC3339Nano
		//tmp := v.Format(time.RFC3339Nano)
		//fmt.Fprintf(os.Stderr, "key=%s time=%s sum=%04X\n",
		//  key, tmp, ChecksumAttr(key, tmp))
		return ChecksumAttr(key, v.Format(time.RFC3339Nano))
	}

	v := reflect.ValueOf(value)
	str := ""
	switch v.Kind() {
	case reflect.String:
		str = v.String()

	case reflect.Bool:
		val := v.Bool()
		str = strconv.FormatBool(val)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val := v.Int() // int64
		str = strconv.FormatInt(val, 10)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uintptr:
		val := v.Uint() // uint64
		str = strconv.FormatUint(val, 10)

	case reflect.UnsafePointer:
		val := v.Pointer() // uintptr
		str = strconv.FormatUint(uint64(val), 10)

	case reflect.Float32, reflect.Float64:
		val := v.Float() // float64
		//str = fmt.Sprintf("%v", val) FIXME: так плохо
		bytes, _ := json.Marshal(val)
		str = string(bytes)

	case reflect.Complex64, reflect.Complex128:
		val := v.Complex() // complex128
		str = fmt.Sprintf("%v", val)

	case reflect.Chan, reflect.Func:
		str = v.Type().Name()

	case reflect.Pointer, reflect.Interface:
		return ChecksumAttr("*"+key, v.Elem().Interface())

	case reflect.Map:
		sum := uint16(0)
		iter := v.MapRange()
		for iter.Next() {
			mkey := iter.Key().Interface()
			value := iter.Value().Interface()
			field, ok := mkey.(string)
			if !ok { // ключ - не строка => преобразовать к строке
				field = Sprint(mkey)
			}
			//fmt.Fprintf(os.Stderr, "key=%s field=%s value=%s sum=%04X\n",
			//  key, field, Sprint(value), ChecksumAttr(field, value))
			sum ^= ChecksumAttr(field, value)
		}
		sum += crc16.Checksum([]byte(key), crcTable)
		return sum

	case reflect.Array, reflect.Slice:
		sum := uint16(0)
		for i := 0; i < v.Len(); i++ {
			index := strconv.Itoa(i)
			value := v.Index(i).Interface()
			//fmt.Fprintf(os.Stderr, "key=%s index=%s value=%s sum=%04X\n",
			//  key, index, Sprint(value), ChecksumAttr(index, value))
			sum ^= ChecksumAttr(index, value)
		}
		return sum + crc16.Checksum([]byte(key), crcTable)

	case reflect.Struct:
		sum := uint16(0)
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			if !field.IsExported() {
				//fmt.Fprintf(os.Stderr, "skip unexported key=%s field=%s\n", key, field.Name)
				continue // пропустить не экспортируемые поля (избежать паники)
			}
			name := field.Name
			if tag := field.Tag.Get("json"); tag != "" {
				name = tag
			}
			value := v.Field(i).Interface()
			//fmt.Fprintf(os.Stderr, "key=%s field=%s value=%s sum=%04X\n",
			//  key, name, Sprint(value), ChecksumAttr(name, value))
			sum ^= ChecksumAttr(name, value)
		}
		return sum + crc16.Checksum([]byte(key), crcTable)
	} // switch

	sum := crc16.Checksum([]byte(key), crcTable)
	sum += crc16.Checksum([]byte(str), crcTable)
	//fmt.Fprintf(os.Stderr, "--> key=%s str=%s sum=%04X\n", key, str, sum) //!!!
	return sum
}

// ChecksumVerify производит вычисление контрольной суммы JSON записи
// и заполняет структуру ChecksumRes.
// Ошибка возвращается, если не удалось распарить данные.
// Сравнение Sum и LogSum должен делать внешний код (при несовпадении
// ошибка не возвращается).
// Если в записи журнала нет одновременно и logId, и logSum, то
// возвращается ошибка.
//
//	full - признак для вычисления контрольной суммы по всем атрибутам рекурсивно
//	rec - запись извлекаемая из журнала с помощью JSON декодера
func ChecksumVerify(full bool, rec map[string]any) (ChecksumRes, error) {
	if !full {
		return ChecksumVerifySimple(rec)
	} else {
		return ChecksumVerifyFull(rec)
	}
}

// ChecksumVerify производит вычисление контрольной суммы JSON записи
// и заполняет структуру ChecksumRes.
// Контрольная сумма вычисляется упрощенно не по всем атрибутам JSON
// (time, level, msg, logId, err).
//
// Ошибка возвращается, если не удалось распарить данные.
// Сравнение Sum и LogSum должен делать внешний код (при несовпадении
// ошибка не возвращается).
// Если в записи журнала нет одновременно и logId, и logSum, то
// возвращается ошибка.
//
//	rec - запись извлекаемая из журнала с помощью JSON декодера
func ChecksumVerifySimple(rec map[string]any) (ChecksumRes, error) {
	res := ChecksumRes{
		Source: make(map[string]any),
		LogId:  uuid.UUID{},
		Sum:    uint16(0),
	}
	logSum := ""

	// Сформировать буфер "sync pool" для сбора данны для CRC
	buf := newBuffer()
	defer buf.Free()

	// Учесть в контрольной сумме метку времени как строку в формате RFC3339Milli
	for k, v := range rec {
		if k == TimeKey { // "time"
			val, ok := v.(string)
			if ok && len(val) > 0 {
				t, err := time.Parse(RFC3339Milli, val)
				if err == nil {
					res.Time = t
					buf.WriteString(t.UTC().Format(RFC3339Milli))
					continue
				}
				return res, fmt.Errorf("can't parse timestamp: %w (time=%s)", err, val)
			}
		}

		//	if k == SourceKey { // "source"
		//		src, ok := v.(map[string]any)
		//		if ok {
		//			res.Source = src
		//		}
		//		continue
		//	}

		if k == SumKey { // "logSum"
			val, ok := v.(string)
			if ok && len(val) > 0 {
				logSum = val
				continue
			}
		}
	} // for k, v

	// Учесть в CRC уровень журналирования как строку
	for k, v := range rec {
		if k == LevelKey { // "level"
			val, ok := v.(string)
			if ok && len(val) > 0 {
				res.Level = LevelFromLabel(val)
				buf.WriteString(LevelToLabel(res.Level))
				continue
			}
		}
	} // for k, v

	// Учесть в CRC текст сообщения
	for k, v := range rec {
		if k == MsgKey { // "msg"
			val, ok := v.(string)
			if ok {
				res.Message = val
				buf.WriteString(val)
				continue
			}
			res.Message = Sprint(val)
			return res, fmt.Errorf("msg is not string: type=%t", v)
		}
	} // for k, v

	// Учесть в CRC первые (старшие) 14 байт UUID идентификатора
	for k, v := range rec {
		if k == IdKey { // "logId"
			val, ok := v.(string)
			if ok && len(val) > 0 {
				id, err := uuid.FromString(val)
				if err == nil {
					res.LogId = id
					*buf = append(*buf, id[0:14]...)
					continue
				}
				return res, fmt.Errorf("can't parse logId (UUID): %w", err)
			}
		}
	} // for k, v

	// Учесть к CRC ошибки "err", если они есть в атрибутах
	//for k, v := range rec {
	//	if k == ErrKey { // "err"
	//		val, ok := v.(string)
	//		if ok {
	//			res.Err = val
	//			buf.WriteString(val)
	//			continue
	//		}
	//		res.Err = Sprint(val)
	//		return res, fmt.Errorf("err is not string: type=%t", v)
	//	}
	//} // for k, v

	// Вычислить CRC16
	res.Sum = crc16.Checksum(*buf, crcTable)

	if res.LogId.IsNil() && logSum == "" {
		return res, fmt.Errorf("logId and logSum are nil both")
	}

	if logSum != "" { // найден отдельный атрибут "logSum"
		sum, err := strconv.ParseInt(logSum, 16, 0)
		if err != nil {
			return res, fmt.Errorf("can't parse logSum: %w", err)
		}
		res.LogSum = uint16(sum)
	} else { // извлечь контрольную сумму из "logId"
		res.LogSum = uint16(res.LogId[14]) << 8
		res.LogSum |= uint16(res.LogId[15])
	}

	return res, nil
}

// ChecksumVerifyFull производит вычисление контрольной суммы JSON записи
// и заполняет структуру ChecksumRes.
// Контрольная сумма вычисляется специальным образом по всем атрибутам JSON.
//
// Ошибка возвращается, если не удалось распарить данные.
// Сравнение Sum и LogSum должен делать внешний код (при несовпадении
// ошибка не возвращается).
// Если в записи журнала нет одновременно и logId, и logSum, то
// возвращается ошибка.
//
//	rec - запись извлекаемая из журнала с помощью JSON декодера
func ChecksumVerifyFull(rec map[string]any) (ChecksumRes, error) {
	res := ChecksumRes{
		Source: make(map[string]any),
		LogId:  uuid.UUID{},
		Sum:    uint16(0),
	}
	logSum := ""

	for k, v := range rec {
		if k == TimeKey { // "time"
			val, ok := v.(string)
			if ok && len(val) > 0 {
				t, err := time.Parse(RFC3339Milli, val)
				if err == nil {
					res.Time = t
					res.Sum ^= ChecksumAttr(TimeKey, t.UTC().Format(RFC3339Milli))
					continue
				}
				return res, fmt.Errorf("can't parse timestamp: %w", err)
			}
		}

		if k == LevelKey { // "level"
			val, ok := v.(string)
			if ok && len(val) > 0 {
				res.Level = LevelFromLabel(val)
				res.Sum ^= ChecksumAttr(LevelKey, res.Level)
				continue
			}
		}

		if k == SourceKey { // "source"
			src, ok := v.(map[string]any)
			if ok {
				res.Source = src
			}
			continue // ссылки на исходные тексты в КС не входят
		}

		if k == MsgKey { // "msg"
			val, ok := v.(string)
			if ok {
				res.Message = val
				res.Sum ^= ChecksumAttr(MsgKey, val)
				continue
			}
			res.Message = Sprint(val)
			return res, fmt.Errorf("msg is not string: type=%t", v)
		}

		if k == GoKey { // "goroutine"
			switch id := v.(type) {
			case int:
				res.Goroutine = id
			case float64:
				res.Goroutine = int(id)
			case json.Number:
				id64, _ := id.Int64()
				res.Goroutine = int(id64)
			}
		}

		if k == IdKey { // "logId"
			val, ok := v.(string)
			if ok && len(val) > 0 {
				id, err := uuid.FromString(val)
				if err == nil {
					res.LogId = id
					res.Sum ^= ChecksumAttr(IdKey, id[:14])
					continue
				}
				return res, fmt.Errorf("can't parse logId (UUID): %w", err)
			}
		}

		if k == SumKey { // "logSum"
			val, ok := v.(string)
			if ok && len(val) > 0 {
				logSum = val
				continue
			}
		}

		res.Sum ^= ChecksumAttr(k, v)
	} // for k, v

	if res.LogId.IsNil() && logSum == "" {
		return res, fmt.Errorf("logId and logSum are nil both")
	}

	if logSum != "" { // найден отдельный атрибут "logSum"
		sum, err := strconv.ParseInt(logSum, 16, 0)
		if err != nil {
			return res, fmt.Errorf("can't parse logSum: %w", err)
		}
		res.LogSum = uint16(sum)
	} else { // извлечь контрольную сумму из "logId"
		res.LogSum = uint16(res.LogId[14]) << 8
		res.LogSum |= uint16(res.LogId[15])
	}

	return res, nil
}

// SourceToString - преобразуем map/JSON представление ссылки на исходные тексты
// (file/function/line) в строку вида "file:function():line".
// Функция может быть полезна для визуализации поля Source структуры
// ChecksumRes.
func (res ChecksumRes) SourceToString() string {
	str := ""
	val, ok := res.Source["file"]
	if ok {
		file, ok := val.(string)
		if ok {
			str = file + ":"
		}
	}

	val, ok = res.Source["function"]
	if ok {
		function, ok := val.(string)
		if ok {
			str += function + "():"
		}
	}

	val, ok = res.Source["line"]
	if ok {
		switch v := val.(type) {
		case json.Number:
			v64, _ := v.Int64()
			str += strconv.FormatInt(v64, 10)
		case int:
			str += strconv.Itoa(v)
		case int64:
			str += strconv.FormatInt(v, 10)
		case float64:
			str += fmt.Sprintf("%.0f", v)
		case string:
			str += v
		}
	}

	return str
}

// EOF: "checksum.go"
