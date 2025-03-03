package xlog // import "github.com/azorg/xlog"

Пакет xlog реализует простую надстройку на стандартными логерами log и slog.

# Краткий экскурс

Логер slog включен в стандартную поставку Go начиная с версии 1.21 ("log/slog").
До этого логер представлен экспериментальным пакетом "golang.org/x/exp/slog".

slog - это "конструктор", "концепция", которая позволяет построить нужный именно
Вам логгер в Go приложениях.

Пакет `xlog` - это пример использования slog с возможными его расширениями
(имеется дополнительный tinted handler, для "пестрой" раскраски с помощью
ANSI/Escape последовательностей, добавлены дополнительные уровни журналирования
типа Trace, Notice, Critical, Fatal и т.п.).

Заложены "мостики" единообразного поведения стандартного логгера "log" при
работе через настроенный логгер slog/xlog.

В ряде случаев из всего пакета может быть полезна только одна функция Setup(),
которая на основе структуры конфигурации может настроить нужное поведение
глобального стандартного Go логера "log" и структурного логера "slog".

Можно предположить вариант создания структуры конфигурации "по умолчанию" с
помощью функции NewConf(), внесения необходимых изменений (управление уровнем
журналирования например) и последующим вызовом Setup(). Например так:

    conf = xlog.NewConf() // создать структура конфигурации по умолчанию
    conf.Level = "debug" // задать уровень отладки
    conf.Tint = true // включить структурный цветной логер
    conf.Src = true // включить в журнал имена файлов и номера строк исходных текстов
    xlog.Setup(conf) // настроить глобальный логер (xlog + slog + log)

Структура Conf имеет JSON тега для возможности сериализации в/из JSON.

# Возможные конфигурации логера

    | Slog  | JSON  | Tint  | Описание                        |
    |:-----:|:-----:|:-----:|:--------------------------------|
    | false | false | false | логер по умолчанию              |
    | true  | false | false | структурный TextHandler         |
    | x     | true  | false | структурный JSONHandler         |
    | x     | x     | true  | структурный цветной TintHandler |

Если планируется настраивать только стандартный логер "log", то достаточно
применения функций NewLog()/SetupLog() заместо Setup().

Для настройки только структурного логера slog стоит использовать функцию
NewSlog(). Функция получает структуру конфигурации и возвращает стандартный
`*slog.Logger`.

Для дополнительных "сахарных" функций реализована обёртка xlog.Logger,
которая хранит указатель на структурный логер `*slog.Logger`, xlog.Leveler и
xlog.Writer, предоставляет дополнительные методы типа `Infof()`, `Trace()`,
`Fatalf()` и др. Имеются функции - аналоги соответствующих методом для
использования глобального логера. Через функции и методы SetLevel()/SetLvl()
можно изменять уровень логирования "на лету".

# Структуры данных:

    Conf - Обобщенная структура конфигурации логера, имеет JSON тэги
    Logger - Структура/обёртка над slog для добавления методов типа Debugf/Noticef/Errorf/Trace

# Интерфейсы:

    Xlogger - интерфейс к структуре Logger (приведен для наглядности API)
    Leveler - интерфейс управления уровнем журналирования
    Writer - обобщенный интерфейс для записи логов, влкючая функци ротации

# Функции настройки конфигурации:

    NewConf() - заполнить обобщенную структуру конфигурации логгера значениями по умолчанию
    SetupLog() - настроить стандартный логгер в соответствии с заданной структурой конфигурации
    SetupLogEx() - настроить стандартный логгер с ротатором
    SetupLogCustom() - настроить стандартный логгер с заданным io.Writer
    NewLog() - создать стандартный логгер log.Logger в соответствии со структурой конфигурации
    New() - создать структурированный логгер на основе заданной конфигурации
    NewEx() - создать структурированный логгер с заданным ротаторм
    NewCustom() - создать структурированный логгер с заданным io.Writer
    NewSlog() - создать структурированный логгер slog.Logger в соответствии со структурой конфигурации
    Setup() - настроить стандартный и структурированный логгеры по умолчанию в соответствии с структурой конфигурации
    GetLevel()/SetLvl() - вернуть/установить текущий уровень журналирования
    GetLevel()/SetLevel() - вернуть/установить текущий уровень журналирования (как slog.Level)

# Функции для работы с надстройкой Logger:

    Default() - Создать логер на основе исходного slog.Deafult()
    Current() - Вернуть текущий глобальный логер
    Slog() - Вернуть текущий глобальный логер slog.Logger
    X() - Создать логер на основе логера slog (для доступа к "сахарным" методам xlog)

# Методы для работы с Logger (методы интерфейса Xlogger):

    With() - Создать дочерний логгер с дополнительными атрибутами
    WithAttrs() - Создать дочерний логгер с дополнительными атрибутами
    WithGroup() - Создать дочерний логгер с группировкой ключей
    Slog() - Обратное преобразование *xlog.Logger -> *slog.Logger
    SetDefault() - Установить логер как xlog по умолчанию
    SetDefaultLogs() - Установить логер как log/slog/xlog по умолчанию
    GetLevel() - получить текуший уровень журналирования (slog.Level)
    SetLevel(l) - обновить текущий уровень журналирования (slog.Level)
    GetLvl() - получить текущий уровень журналирования в виде строки
    SetLvl() - обновить текущий уровень журналирования в виде строки
    NewLog(prefix string) *log.Logger - вернуть стандартный логгер с префиксом
    NewWriter(slog.Level) - создать io.Writer перенаправляемый в журнал

# Методы жля ротаци логов:

    Rotable() - проверить возможна ли ротация
    Rotate() - совершить ротацию логов (обычно по сигналу SIGHUP)

# Методы для использования xlog.Logger с дополнительными уровнями:

    Log(level slog.Level, msg string, args ...any)
    Flood(msg string, args ...any)
    Trace(msg string, args ...any)
    Debug(msg string, args ...any)
    Info(msg string, args ...any)
    Notice(msg string, args ...any)
    Warn(msg string, args ...any)
    Error(msg string, args ...any)
    Crit(msg string, args ...any)
    Fatal(msg string, args ...any)
    Panic(msg string)

    Logf(level slog.Level, format string, args ...any)
    Floodf(format string, args ...any)
    Tracef(format string, args ...any)
    Debugf(format string, args ...any)
    Infof(format string, args ...any)
    Noticef(format string, args ...any)
    Warnf(format string, args ...any)
    Errorf(format string, args ...any)
    Critf(format string, args ...any)
    Fatalf(format string, args ...any)

Примечание: имеются аналогичные глобальные функции в пакете для использования
глобального логера.

# Вспомогательные функции работы с уровнями журналирования:

    ParseLvl(lvl string) slog.Level - получить уровень из строки типа "debug"
    ParseLevel(level slog.Level) string - преобразовать уровень к строке

# Методы интерфейса Leveler:

    Level() slog.Level - получить уровень журналирования (реализация интерфейса slog.Leveler)
    Update(slog.Level) - обновить уровень журналирования
    String() string - сформировать метку для журнала
    ColorString() string - сформировать метку для журнала с ANSI/Escape подкраской

# Ротация логов

Используется популярный пакет https://github.com/natefinch/lumberjack

# Пути вывода журнала

Определено несколько вариантов вывода журнала.

 1. В никуда (если в Pipe="null", File="")
 2. Только в pipe (stdout/stderr)
 3. В кастомный io.Writer
 4. В кастомный io.Writer и дополнительно в pipe
 5. В заданный файл
 6. В заданный файл и дополнительно в pipe
 7. В заданный файл с ротацией
 8. В заданный файл с ротацией и дополнительно в pipe

Интерфейс для направления журнала доступен в поле Writer структуры Logger.

Выбор пути определяется полями Pipe, File, Rotation.Enable в структуре
конфигурации.

Для определения катомного io.Writer'а используются отдельные функции (типа
NewCustom()).

# Полезные ссылки на документацию

  - https://go.dev/blog/slog - Structured Logging with slog (August 2023)
  - https://pkg.go.dev/log - стандартный логер в Go
  - https://pkg.go.dev/log/slog - `slog` из стандартной поставки начиная с
    go-1.21.0
  - https://pkg.go.dev/golang.org/x/exp/slog - `slog` для go-1.20 и младше
  - https://github.com/golang/example/blob/master/slog-handler-guide/README.md

# Как перейти с go1.20.x на go1.21.x и старше?

 1. Убрать из секцию import "golang.org/x/exp/slog"
 2. Добавить в секцию import "log/slog"
 3. В файле "const.go" задать `OLD_SLOG_FIX = false`

Переход на "golang.org/x/exp/slog" (для go-1.20) делается в обратной
последовательности.

CONSTANTS

const (
	AnsiReset            = "\033[0m"        // All attributes off
	AnsiFaint            = "\033[2m"        // Decreased intensity
	AnsiResetFaint       = "\033[22m"       // Normal color (reset faint)
	AnsiRed              = "\033[31m"       // Red
	AnsiGreen            = "\033[32m"       // Green
	AnsiYellow           = "\033[33m"       // Yellow
	AnsiBlue             = "\033[34m"       // blue
	AnsiMagenta          = "\033[35m"       // Magenta
	AnsiCyan             = "\033[36m"       // Cyan
	AnsiWhile            = "\033[37m"       // While
	AnsiBrightRed        = "\033[91m"       // Bright Red
	AnsiBrightRedNoFaint = "\033[91;22m"    // Bright Red and normal intensity
	AnsiBrightGreen      = "\033[92m"       // Bright Green
	AnsiBrightYellow     = "\033[93m"       // Bright Yellow
	AnsiBrightBlue       = "\033[94m"       // Bright Blue
	AnsiBrightMagenta    = "\033[95m"       // Bright Magenta
	AnsiBrightCyan       = "\033[96m"       // Bright Cyan
	AnsiBrightWight      = "\033[97m"       // Bright White
	AnsiBrightRedFaint   = "\033[91;2m"     // Bright Red and decreased intensity
	AnsiBlackOnWhite     = "\033[30;107;1m" // Black on Bright White background
	AnsiBlueOnWhite      = "\033[34;47;1m"  // Blue on Bright White background
	AnsiWhiteOnMagenta   = "\033[37;45;1m"  // Bright White on Magenta background
	AnsiWhiteOnRed       = "\033[37;41;1m"  // White on Red background
)
    ANSI modes

const (
	AnsiFlood    = AnsiGreen
	AnsiTrace    = AnsiBrightBlue
	AnsiDebug    = AnsiBrightCyan
	AnsiInfo     = AnsiBrightGreen
	AnsiNotice   = AnsiBrightMagenta
	AnsiWarn     = AnsiBrightYellow
	AnsiError    = AnsiBrightRed
	AnsiCritical = AnsiWhiteOnRed
	AnsiFatal    = AnsiWhiteOnMagenta
	AnsiPanic    = AnsiBlackOnWhite
)
    Level keys ANSI colors

const (
	AnsiTime   = AnsiYellow
	AnsiSource = AnsiMagenta
	AnsiKey    = AnsiCyan
	AnsiErrKey = AnsiRed
	AnsiErrVal = AnsiBrightRed
)
    Log part colors

const (
	PIPE      = ""      // log pipe ("stdout", "stderr", "null" or "")
	FILE      = ""      // log file path or ""
	FILE_MODE = "0640"  // log file mode
	LEVEL     = LvlInfo // log level (flood/trace/debug/info/warn/error/critical/fatal/silent)
	SLOG      = false   // use slog instead standard log (slog.TextHandler)
	JSON      = false   // use JSON log (slog.JSONHandelr)
	TINT      = false   // use tinted (colorized) log (xlog.TintHandler)
	TIME      = false   // add time stamp
	TIME_US   = false   // us time stamp (only if SLOG=false)
	TIME_TINT = ""      // tinted log time format (~time.Kitchen, "15:04:05.999")
	SRC       = false   // log file name and line number
	SRC_LONG  = false   // log long file path (directory + file name)
	SRC_FUNC  = false   // add function name to log
	NO_EXT    = false   // remove ".go" extension from file name
	NO_LEVEL  = false   // don't print log level tag to log (~level="INFO")
	NO_COLOR  = false   // don't use tinted colors (only if Tint=true)
	PREFIX    = ""      // add prefix to standard log (SLOG=false)
	ADD_KEY   = ""      // add key to structured log (SLOG=true)
	ADD_VALUE = ""      // add value to structured log (SLOG=true

	// Log rotate
	ROTATE_ENABLE      = true
	ROTATE_MAX_SIZE    = 10 // megabytes
	ROTATE_MAX_AGE     = 10 // days
	ROTATE_MAX_BACKUPS = 100
	ROTATE_LOCAL_TIME  = true
	ROTATE_COMPRESS    = true
)
    Default logger configure

const (
	// Add addition log level marks (TRACE/NOTICE/FATAL/PANIC)
	ADD_LEVELS = true

	// Log file mode in error configuration
	DEFAULT_FILE_MODE = 0600 // read/write only for owner for more secure

	// Set false for go > 1.21 with log/slog
	OLD_SLOG_FIX = false // runtime.Version() < go1.21.0

	// Pretty alignment time format in tinted handler (add zeros to end)
	TINT_ALIGN_TIME = true
)
const (
	LevelFlood    = slog.Level(-12) // FLOOD    (-12)
	LevelTrace    = slog.Level(-8)  // TRACE    (-8)
	LevelDebug    = slog.LevelDebug // DEBUG    (-4)
	LevelInfo     = slog.LevelInfo  // INFO     (0)
	LevelNotice   = slog.Level(2)   // NOTICE   (2)
	LevelWarn     = slog.LevelWarn  // WARN     (4)
	LevelError    = slog.LevelError // ERROR    (8)
	LevelCritical = slog.Level(12)  // CRITICAL (12)
	LevelFatal    = slog.Level(16)  // FATAL    (16)
	LevelPanic    = slog.Level(18)  // PANIC    (18)
	LevelSilent   = slog.Level(20)  // SILENT   (20)
)
    Log levels delivered from slog.Level

const (
	LvlFlood    = "flood"
	LvlTrace    = "trace"
	LvlDebug    = "debug"
	LvlInfo     = "info"
	LvlNotice   = "notice"
	LvlWarn     = "warn"
	LvlError    = "error"
	LvlCritical = "critical"
	LvlFatal    = "fatal"
	LvlPanic    = "panic"
	LvlSilent   = "silent"
)
    Log level as string for setup

const (
	LabelFlood    = "FLOOD"
	LabelTrace    = "TRACE"
	LabelDebug    = "DEBUG"
	LabelInfo     = "INFO"
	LabelNotice   = "NOTICE"
	LabelWarn     = "WARN"
	LabelError    = "ERROR"
	LabelCritical = "CRITICAL"
	LabelFatal    = "FATAL"
	LabelPanic    = "PANIC"
	LabelSilent   = "SILENT"
)
    Log level tags

const (
	// Time OFF
	TIME_OFF = ""

	// Default time format of standard logger
	STD_TIME  = "2006/01/02 15:04:05"
	DATE_TIME = time.DateTime // "2006-01-02 15:04:05"

	// Default time format of standard logger + milliseconds
	STD_TIME_MS  = "2006/01/02 15:04:05.999"
	DATE_TIME_MS = "2006-01-02 15:04:05.999"

	// Default time format of standard logger + microseconds
	STD_TIME_US  = "2006/01/02 15:04:05.999999"
	DATE_TIME_US = "2006-01-02 15:04:05.999999"

	// RFC3339 time format + nanoseconds (slog.TextHandler by default)
	RFC3339Nano = time.RFC3339Nano // "2006-01-02T15:04:05.999999999Z07:00"

	// RFC3339 time format + microseconds
	RFC3339Micro = "2006-01-02T15:04:05.999999Z07:00"

	// RFC3339 time format + milliseconds
	RFC3339Milli = "2006-01-02T15:04:05.999Z07:00"

	// Time only format + microseconds
	TimeOnlyMicro = "15:04:05.999999"

	// Time only format + milliseconds
	TimeOnlyMilli = "15:04:05.999"

	// Time format for file names (no spaces, no ":", sorted by date/time)
	FILE_TIME_FORMAT = "2006-01-02_15.04.05"

	// Compromise time format (no spaces, no ":")
	COMPROMISE_TIME_FORMAT_DS = "2006-01-02_15.04.05.9"
	COMPROMISE_TIME_FORMAT    = "2006-01-02_15.04.05.999"
	COMPROMISE_TIME_FORMAT_US = "2006-01-02_15.04.05.999999"
	COMPROMISE_TIME_FORMAT_NS = "2006-01-02_15.04.05.999999999"

	// Digital clock
	CLOCK_TIME_FORMAT = "15:04"

	// Default (recommended) time format wuth milliseconds
	DEFAULT_TIME_FORMAT = STD_TIME_MS

	// Default (recommended) time format with microseconds
	DEFAULT_TIME_FORMAT_US = STD_TIME_US
)
    Time formats

const BUFFER_DEFAULT_CAP = 1 << 10 // 1K
const BUFFER_MAX_SIZE = 16 << 10 // 16K
const DEFAULT_LEVEL = LevelInfo
const DEFAULT_PREFIX = "LOG_"
    Default prefix

const ERR_KEY = "err"
const NEW_LINE = '\n'
    New line sequence


FUNCTIONS

func AddOpt(opt *Opt, conf *Conf)
    Add parsed command line options to logger config

func Close() error
    Close current log file

func Crit(msg string, args ...any)
    Crit logs at LevelCritical with default logger

func Critf(format string, args ...any)
    Critf logs at LevelCritical as standard logger with default logger

func CropFuncName(function string) string
    Crop function name

func Debug(msg string, args ...any)
    Debug logs at LevelDebug with default logger

func Debugf(format string, args ...any)
    Debugf logs at LevelDebug as standard logger with default logger

func Env(conf *Conf, prefixOpt ...string)
    Add settings from environment variables

func Err(err error) slog.Attr
    Err() returns slog.Attr with "err" key if err != nil

func Error(msg string, args ...any)
    Error logs at LevelError with default logger

func Errorf(format string, args ...any)
    Errorf logs at LevelError as standard logger with default logger

func Fatal(msg string, args ...any)
    Fatal logs at LevelFatal with default logger and os.Exit(1)

func Fatalf(format string, args ...any)
    Fatalf logs at LevelFatal as standard logger with default logger and
    os.Exit(1)

func FileMode(mode string) fs.FileMode
    Convert file mode string (oct like "0640") to fs.FileMode

func Flood(msg string, args ...any)
    Flood logs at LevelFlood with default logger

func Floodf(format string, args ...any)
    Floodf logs at LevelFlood as standard logger with default logger

func GetFuncName(skip int) string
    Return function name

func GetLevel() slog.Level
    Return current log level as int (slog.Level)

func GetLvl() string
    Return current log level as string

func Info(msg string, args ...any)
    Info logs at LevelInfo with default logger

func Infof(format string, args ...any)
    Infof logs at LevelInfo as standard logger with default logger

func Int(key string, value int) slog.Attr
    Integer return slog.Attr if key != "" and value != 0

func Log(level slog.Level, msg string, args ...any)
    Log logs at given level with default logger

func Logf(level slog.Level, format string, args ...any)
    Logf logs at given level as standard logger with default logger

func NewLog(conf Conf) *log.Logger
    Create new configured standard logger

func NewSlog(conf Conf) *slog.Logger
    Create new configured structured logger (default/text/JSON/Tinted handler)

func NewWriter(level slog.Level) io.Writer
    Create log io.Writer based on current logger

func Notice(msg string, args ...any)
    Notice logs at LevelNotice with default logger

func Noticef(format string, args ...any)
    Noticef logs at LevelNotice as standard logger with default logger

func Panic(msg string)
    Panic logs at LevelPanic with default logger and panic

func ParseLevel(level slog.Level) string
    Parse Level (num to string: Level -> Lvl)

func ParseLvl(lvl string) slog.Level
    Parse Lvl (string to num: Lvl -> Level)

func RemoveGoExt(file string) string
    Remove ".go" extension from source file name

func Rotable() bool
    Check current log rotation possible

func Rotate() error
    Rotate current log file

func SetLevel(level slog.Level)
    Set current log level as int (slog.Level)

func SetLvl(level string)
    Set current log level as string

func Setup(conf Conf)
    Setup standart and structured default global loggers

func SetupLog(logger *log.Logger, conf Conf)
    Setup standard simple logger with output to pipe/file with rotation

func SetupLogCustom(logger *log.Logger, conf Conf, writer io.Writer)
    Setup standard simple logger with output to custom io.Writer

func SetupLogEx(logger *log.Logger, conf Conf, writer io.Writer)
    Setup standard simple logger with custom io.Writer

func Slog() *slog.Logger
    Return current *slog.Logger

func String(key, value string) slog.Attr
    String return slog.Attr if key != "" and value != ""

func StringToBool(s string) bool
    String to bool converter

        true:  true, True, yes, YES, on, 1, 2
        false: false, FALSE, no, Off, 0, "Abra-Cadabra"

func StringToInt(s string) int
    String to int converter

func TimeFormat(alias string) (format string, ok bool)
    Return time format by alias

func Trace(msg string, args ...any)
    Trace logs at LevelTrace with default logger

func Tracef(format string, args ...any)
    Tracef logs at LevelTrace as standard logger with default logger

func Warn(msg string, args ...any)
    Warn logs at LevelWarn with default logger

func Warnf(format string, args ...any)
    Warnf logs at LevelWarn as standard logger with default logger


TYPES

type Buffer []byte

func NewBuffer() *Buffer

func (b *Buffer) Free()

func (b *Buffer) String() string

func (b *Buffer) Write(bytes []byte) int

func (b *Buffer) WriteByte(char byte) error

func (b *Buffer) WriteString(str string) int

func (b *Buffer) WriteStringIf(ok bool, str string) int

type Conf struct {
	Pipe     string `json:"pipe"`      // log pipe ("stdout", "stderr" or "null" / "")
	File     string `json:"file"`      // log file path or ""
	FileMode string `json:"file-mode"` // log file mode
	Level    string `json:"level"`     // log level (trace/debug/info/warn/error/fatal/silent)
	Slog     bool   `json:"slog"`      // use slog instead standard log (slog.TextHandler)
	JSON     bool   `json:"json"`      // use JSON log (slog.JSONHandler)
	Tint     bool   `json:"tint"`      // use tinted (colorized) log (xlog.TintHandler)
	Time     bool   `json:"time"`      // add timestamp
	TimeUS   bool   `json:"time-us"`   // use timestamp in microseconds
	TimeTint string `json:"time-tint"` // tinted log time format (like time.Kitchen, time.DateTime)
	Src      bool   `json:"src"   `    // log file name and line number
	SrcLong  bool   `json:"src-long"`  // log long file path (directory + file name)
	SrcFunc  bool   `json:"src-func"`  // add function name to log
	NoExt    bool   `json:"no-ext"`    // remove ".go" extension from file name
	NoLevel  bool   `json:"no-level"`  // don't print log level tag to log (~level="INFO")
	NoColor  bool   `json:"no-color"`  // disable tinted colors (only if Tint=true)
	Prefix   string `json:"preifix"`   // add prefix to standard log (log=false)
	AddKey   string `json:"add-key"`   // add key to structured log (Slog=true)
	AddValue string `json:"add-value"` // add value to structured log (Slog=true

	// Log rotate options
	Rotate RotateOpt `json:"rotate"`
}
    Logger configure structure

func NewConf() Conf
    Create default logger structure

type Level slog.Level
    xlog level delivered from slog.Level, implements slog.Leveler interface

func (lp *Level) ColorString() string
    ColorString() returns a color label for the level

func (lp *Level) Level() slog.Level
    Level() returns log level (Level() - implements slog.Leveler interface)

func (lp *Level) String() string
    String() returns a label for the level

func (lp *Level) Update(level slog.Level)
    Update log level

type Leveler interface {
	Level() slog.Level   // get log level as slog.Level (implements a slog.Leveler interface)
	Update(slog.Level)   // update log level
	String() string      // get log level as label
	ColorString() string // get log level as color label
}
    xlog leveler interface (slog.Leveler + Update() + String()/ColorString()
    methods)

type Logger struct {
	*slog.Logger        // standard slog logger
	Leveler      *Level // current log level with Leveler interface
	Writer              // log writer (rotator) interface
}
    Logger wrapper structure

func Current() *Logger
    Return current Logger

func Default() *Logger
    Create logger based on default slog.Logger

func New(conf Conf) *Logger
    Create new configured structured logger with output to pipe/file with
    rotation

func NewCustom(conf Conf, writer io.Writer) *Logger
    Create new configured structured logger (default/text/JSON/Tinted handler)
    with output to custom io.Writer

func NewEx(conf Conf, writer Writer) *Logger
    Create new configured structured logger with custom Writer

func X(logger *slog.Logger) *Logger
    Create Logger based on *slog.Logger (*slog.Logger -> Logger)

func (x *Logger) Crit(msg string, args ...any)
    Crit logs at LevelCritical

func (x *Logger) Critf(format string, args ...any)
    Critf logs at LevelCritical as standard logger

func (x *Logger) Debug(msg string, args ...any)
    Debug logs at LevelDebug

func (x *Logger) Debugf(format string, args ...any)
    Debugf logs at LevelDebug as standard logger

func (x *Logger) Error(msg string, args ...any)
    Error logs at LevelError

func (x *Logger) Errorf(format string, args ...any)
    Errorf logs at LevelError as standard logger

func (x *Logger) Fatal(msg string, args ...any)
    Fatal logs at LevelFatal and os.Exit(1)

func (x *Logger) Fatalf(format string, args ...any)
    Fatalf logs at LevelFatal as standard logger and os.Exit(1)

func (x *Logger) Flood(msg string, args ...any)
    Flood logs at LevelFlood

func (x *Logger) Floodf(format string, args ...any)
    Floodf logs at LevelFlood as standard logger

func (x *Logger) GetLevel() slog.Level
    Return log level as int (slog.Level)

func (x *Logger) GetLvl() string
    Return log level as string

func (x *Logger) Info(msg string, args ...any)
    Info logs at LevelInfo

func (x *Logger) Infof(format string, args ...any)
    Infof logs at LevelInfo as standard logger

func (x *Logger) Log(level slog.Level, msg string, args ...any)
    Log logs at given level

func (x *Logger) Logf(level slog.Level, format string, args ...any)
    Logf logs at given level as standard logger

func (x *Logger) NewLog(prefix string) *log.Logger
    Return standard logger with prefix

func (x *Logger) NewWriter(level slog.Level) io.Writer
    Create log io.Writer

func (x *Logger) Notice(msg string, args ...any)
    Notice logs at LevelNotice

func (x *Logger) Noticef(format string, args ...any)
    Noticef logs at LevelNotice as standard logger

func (x *Logger) Panic(msg string)
    Panic logs at LevelPanic and panic

func (x *Logger) SetDefault()
    Set logger as default logger

func (x *Logger) SetDefaultLogs()
    Set logger as default xlog/log/slog loggers

func (x *Logger) SetLevel(level slog.Level)
    Set log level as int (slog.Level)

func (x *Logger) SetLvl(level string)
    Set log level as string

func (x *Logger) Slog() *slog.Logger
    Extract *slog.Logger (*xlog.Logger -> *slog.Logger)

func (x *Logger) Trace(msg string, args ...any)
    Trace logs at LevelTrace

func (x *Logger) Tracef(format string, args ...any)
    Tracef logs at LevelTrace as standard logger

func (x *Logger) Warn(msg string, args ...any)
    Warn logs at LevelWarn

func (x *Logger) Warnf(format string, args ...any)
    Warnf logs at LevelWarn as standard logger

func (x *Logger) With(args ...any) *Logger
    Create logger that includes the given attributes in each output

func (x *Logger) WithAttrs(attrs []slog.Attr) *Logger
    Create logger that includes the given attributes in each output

func (x *Logger) WithGroup(name string) *Logger
    Create logger that starts a group

type Opt struct {
	Level    string // -log <level>
	Pipe     string // -lpipe <pipe>
	File     string // -lfile <file>
	SLog     bool   // -slog
	JLog     bool   // -jlog
	TLog     bool   // -tlog
	Src      bool   // -lsrc
	NoSrc    bool   // -lnosrc
	Pkg      bool   // -lpkg
	NoPkg    bool   // -lnopkg
	Func     bool   // -lfunc
	NoFunc   bool   // -lnofunc
	Ext      bool   // -lext
	NoExt    bool   // -lnoext
	Time     bool   // -ltime
	NoTime   bool   // -lnotime
	TimeFmt  string // -ltimefmt <fmt>
	OnLevel  bool   // -lonlevel
	NoLevel  bool   // -lnolevel
	Color    bool   // -lcolor
	NoColor  bool   // -lnocolor
	Rotate   bool   // -lrotate
	NoRotate bool   // -lnorotate
}
    Command line logger option structure

func NewOpt() *Opt
    Setup command line logger options Usage:

        -log <level>        - log level (flood/trace/debug/info/notice/warm/error/critical)
        -lpipe              - log pipe (stdout/stderr/null)
        -lfile <file>       - log file path
        -slog               - use structured text logger (slog)
        -jlog               - use structured JSON logger (slog)
        -tlog               - use tinted (colorized) logger (tint)
        -lsrc|-lnosrc       - force on/off log source file name and line number
        -lpkg|-lnopkg       - force on/off log source directory/file name and line number
        -lfunc|-lnofunc     - force on/off log function name
        -lext|-lnoext       - force enable/disable remove ".go" extension from source file name
        -ltime|-lnotime     - force on/off timestamp
        -ltimefmt <format>  - override log time format (e.g. 15:04:05.999 or TimeOnly)
        -lnolevel|lonlevel  - disable/enable log level tag (~level="INFO")
        -lcolor|-lnocolor   - force enable/disable tinted colors
        -lrotate|-lnorotate - force on/off log rotate

type RotateOpt struct {
	// Enable log rotation
	Enable bool `json:"enable"`

	// Maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int `json:"max-size"`

	// Maximum number of days to retain old log files based on the
	// timestamp encoded in their filename. Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int `json:"max-age"`

	// Maximum number of old log files to retain. The default/ is to retain
	// all old log files (though MaxAge may still cause them to get deleted)
	MaxBackups int `json:"max-backups"`

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.  The default is to use UTC
	// time.
	LocalTime bool `json:"local-time"`

	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool `json:"compress"`
}
    Log rotate options (delivered from lumberjack)

type TintHandler struct {
	// Has unexported fields.
}
    Tinted (colorized) handler implements a slog.Handler

func NewTintHandler(w io.Writer, opts *TintOptions) *TintHandler
    Create new tinted (colorized) handler

func (h *TintHandler) Enabled(_ context.Context, level slog.Level) bool
    Enabled() implements slog.Handler interface

func (h *TintHandler) Format(r slog.Record) string
    Format record to byte array

func (h *TintHandler) Handle(ctx context.Context, r slog.Record) error
    Handle() implements slog.Handler interface

func (h *TintHandler) WithAttrs(attrs []slog.Attr) slog.Handler
    WithAttrs() implements slog.Handler interface

func (h *TintHandler) WithGroup(name string) slog.Handler
    WithGroup() implements slog.Handler interface

type TintOptions struct {
	// Minimum level to log (Default: slog.LevelInfo)
	Level slog.Leveler

	// Enable source code location
	AddSource bool

	// Log long file path (directory + file name)
	SourceLong bool

	// Log functions name
	SourceFunc bool

	// Remove ".go" extension from file name
	NoExt bool

	// Off level keys
	NoLevel bool

	// Time format
	TimeFormat string

	// Disable color
	NoColor bool

	// ReplaceAttr is called to rewrite each non-group attribute before it is logged.
	// See https://pkg.go.dev/log/slog#HandlerOptions for details.
	ReplaceAttr func(groups []string, attr slog.Attr) slog.Attr
}

type Writer interface {
	Rotable() bool // check rotation possible
	Rotate() error // rotate log
	Close() error  // close file
	io.Writer
}
    Interface of log writer

type Xlogger interface {
	// Create Logger that includes the given attributes in each output
	With(args ...any) *Logger

	// Create logger that includes the given attributes in each output
	WithAttrs(attrs []slog.Attr) *Logger

	// Create logger that starts a group
	WithGroup(name string) *Logger

	// Extract *slog.Logger (*xlog.Logger -> *slog.Logger)
	Slog() *slog.Logger

	// Set logger as default xlog logger
	SetDefault()

	// Set logger as default xlog/log/slog loggers
	SetDefaultLogs()

	// Return log level as int (slog.Level)
	GetLevel() slog.Level

	// Set log level as int (slog.Level)
	SetLevel(level slog.Level)

	// Return log level as string
	GetLvl() string

	// Set log level as string
	SetLvl(level string)

	// Return standard logger with prefix
	NewLog(prefix string) *log.Logger

	// Log logs at given level
	Log(level slog.Level, msg string, args ...any)

	// Flood logs at LevelFlood
	Flood(msg string, args ...any)

	// Trace logs at LevelTrace
	Trace(msg string, args ...any)

	// Debug logs at LevelDebug
	Debug(msg string, args ...any)

	// Info logs at LevelInfo
	Info(msg string, args ...any)

	// Notice logs at LevelNotice
	Notice(msg string, args ...any)

	// Warn logs at LevelWarn
	Warn(msg string, args ...any)

	// Error logs at LevelError
	Error(msg string, args ...any)

	// Crit logs at LevelCritical
	Crit(msg string, args ...any)

	// Fatal logs at LevelFatal and os.Exit(1)
	Fatal(msg string, args ...any)

	// Panic logs at LevelPanic and panic
	Panic(msg string)

	// Logf logs at given level as standard logger
	Logf(level slog.Level, format string, args ...any)

	// Floodf logs at LevelFlood as standard logger
	Floodf(format string, args ...any)

	// Tracef logs at LevelTrace as standard logger
	Tracef(format string, args ...any)

	// Debugf logs at LevelDebug as standard logger
	Debugf(format string, args ...any)

	// Infof logs at LevelInfo as standard logger
	Infof(format string, args ...any)

	// Noticef logs at LevelNotice as standard logger
	Noticef(format string, args ...any)

	// Warnf logs at LevelWarn as standard logger
	Warnf(format string, args ...any)

	// Errorf logs at LevelError as standard logger
	Errorf(format string, args ...any)

	// Critf logs at LevelCritical as standard logger
	Critf(format string, args ...any)

	// Fatalf logs at LevelFatal as standard logger and os.Exit(1)
	Fatalf(format string, args ...any)

	// Check log rotation possible
	Rotable() bool

	// Close the existing log file and immediately create a new one
	Rotate() error

	// Close current log file
	Close() error

	// Create log io.Writer
	NewWriter(slog.Level) io.Writer
}
    Logger interface

