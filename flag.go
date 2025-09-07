// File: "flag.go"

package xlog

import "flag"

// Префикс для флагов по умолчанию
const DefaultFlagPrefix = "log-"

// Структура управления журналированием на основе опций командной строки.
// Типовое использование:
//
//	opt := xlog.NewOpt()  // создать набор опций (*xlog.Opt)
//	conf := xlog.Conf{}   // подготовить структуру конфигурации логгера
//	xlog.Env(&conf)       // обогатить conf переменными окружения
//	flag.Parse()          // обогатить opt из опций командной строки
//	opt.UpdateConf(&conf) // обогатить conf опциями командной строки
//
//	log := xlog.New(conf) // создать логгер (*xlog.Logger)
//	logger := log.Logger  // получить указатель на *slog.Logger
//
//	log.Notice("Привет, X-logger", "version", "1.0.0")
//	mylog := logger.With("app", "helloworld")
//	mylog.Info("application started")
type Opt struct {
	Level            string // -log-level
	Pipe             string // -log-pipe
	File             string // -log-file
	FileMode         string // -log-file-mode
	Format           string // -log-format
	GoId             string // -log-goid
	Id               string // -log-id
	Sum              string // -log-sum
	SumFull          string // -log-sum-full
	SumChain         string // -log-sum-chain
	SumAlone         string // -log-sum-alone
	Time             string // -log-time
	TimeLocal        string // -log-time-local
	TimeMicro        string // -log-time-micro
	TimeFormat       string // -log-time-format
	Src              string // -log-src
	SrcPkg           string // -log-src-pkg
	SrcFunc          string // -log-src-func
	SrcExt           string // -log-src-ext
	Color            string // -log-color
	LevelOff         string // -log-level-off
	Rotate           string // -log-rotate
	RotateMaxSize    string // -log-rotate-max-size
	RotateMaxAge     string // -log-rotate-max-age
	RotateMaxBackups string // -log-rotate-max-backups
	RotateLocalTime  string // -log-rotate-local-time
	RotateCompress   string // -log-rotate-compress
}

// NewOpt создаёт набор опций командной строки с параметрами для X-logger'а.
// После создания опций Opt можно использовать стандартный вызов flag.Parse()
// для заполнения полей структуры. Булевы переменные обрабатываются так же
// как и переменные окружения.
//
//	prefixOpt - опциональный префикс (по умолчанию "log-")
//
// Приложения могут включить в свой usage-вывод следующий текст:
//
//	-log-level <level>              - log level (flood/trace/debug/info/notice/warm/error/crit)
//	-log-pipe <pipe>                - log pipe (stdout/stderr/null)
//	-log-file <file>                - log file path
//	-log-file-mode <perm>           - log file mode (0640, 0600, 0644)
//	-log-format <format>            - log format (json|prod/text|logfmt/tint|tinted|human/default|std)
//	-log-goid <on/off>              - force on/off goroutine id for each record (goroutine)
//	-log-id <on/off>                - force on/off id (UUID) for each record (logId)
//	-log-sum <on/off>               - force on/off check sum for each record
//	-log-sum-full <on/off>          - force on/off calculate full sum for earch record
//	-log-sum-chain <on/off>         - force on/off check sum chain
//	-log-sum-alone <on/off>         - force on/off add check sum as alone atribute (logSum)
//	-log-time <on/off>              - force on/off timestamp
//	-log-time-local <on/off>        - use local time (UTC by default)
//	-log-time-micro <on/off>        - force on/off microseconds in timestamp
//	-log-time-format <fmt>          - override tinted log time format (e.g. 15:04:05.999 or timeOnly)
//	-log-src <on/off>               - force on/off log source file name and line number
//	-log-src-pkg <on/off>           - force on/off log source directory/file name and line number
//	-log-src-func <on/off>          - force on/off log function name
//	-log-src-ext <on/off>           - force enable/disable show ".go" extension of source file name
//	-log-color <on/off>             - force enable/disable tinted colors (ANSI/Escape)
//	-log-level-off <true/false>     - force disable/enable level output
//	-log-rotate <on/off>            - force on/off log rotate
//	-log-rotate-max-size <mb>       - rotate max size (begabytes)
//	-log-rotate-max-age <days>      - rotate max age (days)
//	-log-rotate-max-backups <num>   - rotate max backup files
//	-log-rotate-local-time <yes/no> - use localtime (default UTC)
//	-log-rotate-compress <on/off>   - on/off compress (gzip)
func NewOpt(prefixOpt ...string) *Opt {
	prefix := DefaultFlagPrefix
	if len(prefixOpt) != 0 {
		prefix = prefixOpt[0]
	}
	opt := &Opt{}

	flag.StringVar(&opt.Level, prefix+"level", "", "override log level (flood/trace/debug/info/notice/warm/error/crit)")
	flag.StringVar(&opt.Pipe, prefix+"pipe", "", "log pipe (stdout/stderr/null)")
	flag.StringVar(&opt.File, prefix+"file", "", "log file path")
	flag.StringVar(&opt.FileMode, prefix+"file-mode", "", "log file mode (0640, 0600, 0644)")
	flag.StringVar(&opt.Format, prefix+"format", "", "log format (json|prod/text|logfmt/tint|tinted|human/std|default)")
	flag.StringVar(&opt.GoId, prefix+"goid", "", "force on/off goroutine id for each record (goroutine)")
	flag.StringVar(&opt.Id, prefix+"id", "", "force on/off id (UUID) for each record (logId)")
	flag.StringVar(&opt.Sum, prefix+"sum", "", "force on/off check sum for each record")
	flag.StringVar(&opt.SumFull, prefix+"sum-full", "", "force on/off calculate full check sum for each record")
	flag.StringVar(&opt.SumChain, prefix+"sum-chain", "", "force on/off check sum chain")
	flag.StringVar(&opt.SumAlone, prefix+"sum-alone", "", "force on/off add check sum as alone atribute (logSum)")
	flag.StringVar(&opt.Time, prefix+"time", "", "force on/off timestamp")
	flag.StringVar(&opt.TimeLocal, prefix+"time-local", "", "use local time (UTC by default)")
	flag.StringVar(&opt.TimeMicro, prefix+"time-micro", "", "force on/off microseconds in timestamp")
	flag.StringVar(&opt.TimeFormat, prefix+"time-format", "", "override tinted log time format (e.g. 15:04:05.999 or TimeOnly)")
	flag.StringVar(&opt.Src, prefix+"src", "", "force on/off log source file name and line number")
	flag.StringVar(&opt.SrcPkg, prefix+"src-pkg", "", "force on/off log source directory/file name and line number")
	flag.StringVar(&opt.SrcFunc, prefix+"src-func", "", "force enable/disable functions name")
	flag.StringVar(&opt.SrcExt, prefix+"src-ext", "", "force enable/disable show '.go' extension of source file name")
	flag.StringVar(&opt.Color, prefix+"color", "", "force enable/disable tinted colors")
	flag.StringVar(&opt.LevelOff, prefix+"level-off", "", "force disable/enable level output")
	flag.StringVar(&opt.Rotate, prefix+"rotate", "", "force enable/disable log rotate")
	flag.StringVar(&opt.RotateMaxSize, prefix+"rotate-max-size", "", "rotate max size (begabytes)")
	flag.StringVar(&opt.RotateMaxAge, prefix+"rotate-max-age", "", "rotate max age (days)")
	flag.StringVar(&opt.RotateMaxBackups, prefix+"rotate-max-backups", "", "rotate max backup files")
	flag.StringVar(&opt.RotateLocalTime, prefix+"rotate-local-time", "", "use localtime (default UTC)")
	flag.StringVar(&opt.RotateCompress, prefix+"rotate-compress", "", "compress (gzip)")

	return opt
}

// UpdateConf обогащает структуру конфигурации логгера опциями
// командной строки. Если соответствующие опции командной строки не
// заданы, то поля структуры конфигурации conf не модифицируются.
func (opt *Opt) UpdateConf(conf *Conf) {
	if opt.Level != "" {
		conf.Level = opt.Level
	}
	if opt.Pipe != "" {
		conf.Pipe = opt.Pipe
	}
	if opt.File != "" {
		conf.File = opt.File
	}
	if opt.FileMode != "" {
		conf.FileMode = opt.FileMode
	}
	if opt.Format != "" {
		conf.Format = opt.Format
	}
	if opt.GoId != "" {
		conf.GoId = StringToBool(opt.GoId)
	}
	if opt.Id != "" {
		conf.IdOn = StringToBool(opt.Id)
	}
	if opt.Sum != "" {
		conf.SumOn = StringToBool(opt.Sum)
	}
	if opt.SumFull != "" {
		conf.SumFull = StringToBool(opt.Sum)
	}
	if opt.SumChain != "" {
		conf.SumChain = StringToBool(opt.SumChain)
	}
	if opt.SumAlone != "" {
		conf.SumAlone = StringToBool(opt.SumAlone)
	}
	if opt.Src != "" {
		conf.Src = StringToBool(opt.Src)
	}
	if opt.SrcPkg != "" {
		conf.SrcPkg = StringToBool(opt.SrcPkg)
		conf.Src = conf.Src || conf.SrcPkg
	}
	if opt.SrcFunc != "" {
		conf.SrcFunc = StringToBool(opt.SrcFunc)
		conf.Src = conf.Src || conf.SrcFunc
	}
	if opt.SrcExt != "" {
		conf.SrcExt = StringToBool(opt.SrcExt)
	}
	if opt.Time != "" {
		conf.TimeOff = !StringToBool(opt.Time)
		if conf.TimeOff {
			conf.TimeFormat = timeOff
		}
	}
	if opt.TimeLocal != "" {
		conf.TimeLocal = StringToBool(opt.TimeLocal)
	}
	if opt.TimeMicro != "" {
		conf.TimeMicro = StringToBool(opt.TimeMicro)
	}
	if opt.TimeFormat != "" {
		conf.TimeOff = false
		conf.TimeFormat = opt.TimeFormat
	}
	if opt.Color != "" {
		conf.ColorOff = !StringToBool(opt.Color)
	}
	if opt.LevelOff != "" {
		conf.LevelOff = StringToBool(opt.LevelOff)
	}
	if opt.Rotate != "" {
		conf.Rotate.Enable = StringToBool(opt.Rotate)
	}
	if opt.RotateMaxSize != "" {
		conf.Rotate.MaxSize = StringToInt(opt.RotateMaxSize)
	}
	if opt.RotateMaxAge != "" {
		conf.Rotate.MaxAge = StringToInt(opt.RotateMaxAge)
	}
	if opt.RotateMaxBackups != "" {
		conf.Rotate.MaxBackups = StringToInt(opt.RotateMaxBackups)
	}
	if opt.RotateLocalTime != "" {
		conf.Rotate.LocalTime = StringToBool(opt.RotateLocalTime)
	}
	if opt.RotateCompress != "" {
		conf.Rotate.Compress = StringToBool(opt.RotateCompress)
	}
}

// EOF: "flag.go"
