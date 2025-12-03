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
//	log.Notice("Привет, Логгер", "version", "1.0.0")
//	mylog := logger.With("app", "helloworld")
//	mylog.Info("application started")
type Opt struct {
	Level                  string // -log-level
	Pipe                   string // -log-pipe
	File                   string // -log-file
	FileMode               string // -log-file-mode
	Format                 string // -log-format
	GoId                   string // -log-goid
	Id                     string // -log-id
	Sum                    string // -log-sum
	SumFull                string // -log-sum-full
	SumChain               string // -log-sum-chain
	SumAlone               string // -log-sum-alone
	Time                   string // -log-time
	TimeLocal              string // -log-time-local
	TimeMicro              string // -log-time-micro
	TimeFormat             string // -log-time-format
	Src                    string // -log-src
	SrcPkg                 string // -log-src-pkg
	SrcFunc                string // -log-src-func
	Source                 bool   // -log-source
	SrcExt                 string // -log-src-ext
	Color                  string // -log-color
	LevelOff               string // -log-level-off
	RateLimit              string // -log-rate-limit
	RateLimitMaxNum        string // -log-rate-limit-max-num
	RateLimitIntervalMs    string // -log-rate-limit-interval-ms
	RateLimitFlushPeriodMs string // -log-rate-limit-flush-period-ms
	Rotate                 string // -log-rotate
	RotateMaxSize          string // -log-rotate-max-size
	RotateMaxAge           string // -log-rotate-max-age
	RotateMaxBackups       string // -log-rotate-max-backups
	RotateLocalTime        string // -log-rotate-local-time
	RotateCompress         string // -log-rotate-compress
}

// NewOpt создаёт набор опций командной строки с параметрами логгера.
// После создания опций Opt можно использовать стандартный вызов flag.Parse()
// для заполнения полей структуры. Булевы переменные обрабатываются так же
// как и переменные окружения.
//
// Данная функция вызывает функцию NewOptSet с передачей глобального набора
// флагов flag.CommandLine в качестве первого аргумента.
//
//	prefixOpt - опциональный префикс (по умолчанию "log-")
//
// Приложения могут включить в свой usage-вывод следующий текст:
//
//	-log-level <level>                   - log level (flood/trace/debug/info/notice/warm/error/crit)
//	-log-pipe <pipe>                     - log pipe (stdout/stderr/null)
//	-log-file <file>                     - log file path
//	-log-file-mode <perm>                - log file mode (0640, 0600, 0644)
//	-log-format <format>                 - log format (json|prod/text|logfmt/tint|tinted|human/default|std)
//	-log-goid <on/off>                   - force on/off goroutine id for each record (goroutine)
//	-log-id <on/off>                     - force on/off id (UUID) for each record (logId)
//	-log-sum <on/off>                    - force on/off check sum for each record
//	-log-sum-full <on/off>               - force on/off calculate full sum for earch record
//	-log-sum-chain <on/off>              - force on/off check sum chain
//	-log-sum-alone <on/off>              - force on/off add check sum as alone atribute (logSum)
//	-log-time <on/off>                   - force on/off timestamp
//	-log-time-local <on/off>             - use local time (UTC by default)
//	-log-time-micro <on/off>             - force on/off microseconds in timestamp
//	-log-time-format <fmt>               - override tinted log time format (e.g. 15:04:05.999 or timeOnly)
//	-log-src <on/off>                    - force on/off log source file name and line number
//	-log-src-pkg <on/off>                - force on/off log source directory/file name and line number
//	-log-src-func <on/off>               - force on/off log function name
//	-log-src-ext <on/off>                - force enable/disable show ".go" extension of source file name
//	-log-color <on/off>                  - force enable/disable tinted colors (ANSI/Escape)
//	-log-level-off <true/false>          - force disable/enable level output
//	-log-rate-limit <on/off>             - force enable/disable rate limiter
//	-log-rate-limit-max-num <int>        - maximal number of rate limit messages
//	-log-rate-limit-interval-ms <ms>     - rate limiter interval [ms]
//	-log-rate-limit-flush-period-ms <ms> - rate limiter flush period [ms]
//	-log-rotate <on/off>                 - force on/off log rotate
//	-log-rotate-max-size <mb>            - rotate max size (begabytes)
//	-log-rotate-max-age <days>           - rotate max age (days)
//	-log-rotate-max-backups <num>        - rotate max backup files
//	-log-rotate-local-time <yes/no>      - use localtime (default UTC)
//	-log-rotate-compress <on/off>        - on/off compress (gzip)
func NewOpt(prefixOpt ...string) *Opt {
	return NewOptSet(flag.CommandLine, prefixOpt...)
}

// NewOptSet создаёт набор опций командной строки с параметрами логгера
// для заданного flag.FlagSet.
//
// Данная расширенная версия функции NewOpt может использоваться в
// мультифункциональных приложениях, где апплеты могут иметь
// индивидуальные опции командной строки для настройки собсвенного логгера.
func NewOptSet(fs *flag.FlagSet, prefixOpt ...string) *Opt {
	prefix := DefaultFlagPrefix
	if len(prefixOpt) != 0 {
		prefix = prefixOpt[0]
	}
	opt := &Opt{}

	fs.StringVar(&opt.Level, prefix+"level", "", "override log level (flood/trace/debug/info/notice/warm/error/crit)")
	fs.StringVar(&opt.Pipe, prefix+"pipe", "", "log pipe (stdout/stderr/null)")
	fs.StringVar(&opt.File, prefix+"file", "", "log file path")
	fs.StringVar(&opt.FileMode, prefix+"file-mode", "", "log file mode (0640, 0600, 0644)")
	fs.StringVar(&opt.Format, prefix+"format", "", "log format (json|prod/text|logfmt/tint|tinted|human/std|default)")
	fs.StringVar(&opt.GoId, prefix+"goid", "", "force on/off goroutine id for each record (goroutine)")
	fs.StringVar(&opt.Id, prefix+"id", "", "force on/off id (UUID) for each record (logId)")
	fs.StringVar(&opt.Sum, prefix+"sum", "", "force on/off check sum for each record")
	fs.StringVar(&opt.SumFull, prefix+"sum-full", "", "force on/off calculate full check sum for each record")
	fs.StringVar(&opt.SumChain, prefix+"sum-chain", "", "force on/off check sum chain")
	fs.StringVar(&opt.SumAlone, prefix+"sum-alone", "", "force on/off add check sum as alone atribute (logSum)")
	fs.StringVar(&opt.Time, prefix+"time", "", "force on/off timestamp")
	fs.StringVar(&opt.TimeLocal, prefix+"time-local", "", "use local time (UTC by default)")
	fs.StringVar(&opt.TimeMicro, prefix+"time-micro", "", "force on/off microseconds in timestamp")
	fs.StringVar(&opt.TimeFormat, prefix+"time-format", "", "override tinted log time format (e.g. 15:04:05.999 or TimeOnly)")
	fs.StringVar(&opt.Src, prefix+"src", "", "force on/off log source file name and line number")
	fs.StringVar(&opt.SrcPkg, prefix+"src-pkg", "", "force on/off log source directory/file name and line number")
	fs.StringVar(&opt.SrcFunc, prefix+"src-func", "", "force enable/disable functions name")
	fs.StringVar(&opt.SrcExt, prefix+"src-ext", "", "force enable/disable show '.go' extension of source file name")
	fs.BoolVar(&opt.Source, prefix+"source", false, "force log source info (package/file/function)")
	fs.StringVar(&opt.Color, prefix+"color", "", "force enable/disable tinted colors")
	fs.StringVar(&opt.LevelOff, prefix+"level-off", "", "force disable/enable level output")
	fs.StringVar(&opt.RateLimit, prefix+"rate-limit", "", "force enable/disable rate limiter")
	fs.StringVar(&opt.RateLimitMaxNum, prefix+"rate-limit-max-num", "", "maximal number of rate limit messages")
	fs.StringVar(&opt.RateLimitIntervalMs, prefix+"rate-limit-interval-ms", "", "rate limiter interval [ms]")
	fs.StringVar(&opt.RateLimitFlushPeriodMs, prefix+"rate-limit-flush-period-ms", "", "rate limiter flush period [ms]")
	fs.StringVar(&opt.Rotate, prefix+"rotate", "", "force enable/disable log rotate")
	fs.StringVar(&opt.RotateMaxSize, prefix+"rotate-max-size", "", "rotate max size (begabytes)")
	fs.StringVar(&opt.RotateMaxAge, prefix+"rotate-max-age", "", "rotate max age (days)")
	fs.StringVar(&opt.RotateMaxBackups, prefix+"rotate-max-backups", "", "rotate max backup files")
	fs.StringVar(&opt.RotateLocalTime, prefix+"rotate-local-time", "", "use localtime (default UTC)")
	fs.StringVar(&opt.RotateCompress, prefix+"rotate-compress", "", "compress (gzip)")

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
	if opt.Source {
		conf.Src = true
		conf.SrcPkg = true
		conf.SrcFunc = true
		conf.SrcExt = false
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
	if opt.RateLimit != "" {
		conf.RateLimit.Disable = !StringToBool(opt.RateLimit)
	}
	if opt.RateLimitMaxNum != "" {
		conf.RateLimit.MaxNum = StringToInt(opt.RateLimitMaxNum)
	}
	if opt.RateLimitIntervalMs != "" {
		conf.RateLimit.IntervalMs = StringToInt(opt.RateLimitIntervalMs)
	}
	if opt.RateLimitFlushPeriodMs != "" {
		conf.RateLimit.FlushPeriodMs = StringToInt(opt.RateLimitFlushPeriodMs)
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
