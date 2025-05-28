package logger

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger é a estrutura que encapsula o logger Zap.
type Logger struct {
	zap   *zap.Logger
	sugar *zap.SugaredLogger
}

// Config contém as opções de configuração para o logger.
type Config struct {
	//Level         string              // Níveis: "debug", "info", "warn", "error", "fatal", "panic"
	IsDevelopment bool                // Define se o output é para desenvolvimento (console) ou produção (JSON)
	Output        zapcore.WriteSyncer // Onde o log será escrito (padrão: os.Stdout)
	CallerSkip    int                 // Quantos frames de chamada pular para o log do chamador (padrão: 1)
}

// NewLogger cria e retorna uma nova instância do Logger.
// Ele usa a variável de ambiente "DEV" para definir o modo de desenvolvimento se `Config.IsDevelopment`
// não for explicitamente definido (ou seja, se for o valor zero/false de bool).
// Se a variável "DEV" estiver definida como "true" (case-insensitive),
// o logger será configurado para o modo de desenvolvimento.
func NewLogger(cfg Config) (*Logger, error) {
	var logLevel zapcore.Level
	var envLogLevel string = os.Getenv("LOG_LEVEL")
	// Configura o nível de log
	switch strings.ToLower(envLogLevel) {
	case "debug":
		logLevel = zapcore.DebugLevel
	case "info":
		logLevel = zapcore.InfoLevel
	case "warn":
		logLevel = zapcore.WarnLevel
	case "error":
		logLevel = zapcore.ErrorLevel
	case "fatal":
		logLevel = zapcore.FatalLevel
	case "panic":
		logLevel = zapcore.PanicLevel
	default:
		logLevel = zapcore.InfoLevel // Padrão para Info se nível inválido
		fmt.Printf("Nível de log inválido, usando 'info' como padrão.\n")
	}

	// Determina o modo de desenvolvimento
	// Se cfg.IsDevelopment não foi definido, verifica a variável de ambiente
	isDevelopmentMode := cfg.IsDevelopment
	if !isDevelopmentMode { // Se o usuário não especificou, tentamos via ENV
		devEnv := strings.ToLower(os.Getenv("DEV_MODE"))
		if devEnv == "true" {
			isDevelopmentMode = true
		}
	}

	output := cfg.Output
	if output == nil {
		output = zapcore.AddSync(os.Stdout) // Loga para stdout como padrão
	}

	var encoder zapcore.Encoder
	if isDevelopmentMode {
		// Para desenvolvimento, um formato mais legível no console
		devEncoderConfig := zap.NewDevelopmentEncoderConfig()
		devEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Níveis com cor
		devEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        // Mantém o formato de tempo
		devEncoderConfig.CallerKey = "caller"                           // Mostra o chamador
		encoder = zapcore.NewConsoleEncoder(devEncoderConfig)
	} else {
		// Para produção, JSON
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // Formato de tempo padrão
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // Nível em maiúsculas
		encoderConfig.TimeKey = "timestamp"                     // Nome do campo para o timestamp
		encoderConfig.MessageKey = "message"                    // Nome do campo para a mensagem
		encoderConfig.LevelKey = "level"                        // Nome do campo para o nível
		encoderConfig.CallerKey = "caller"                      // Nome do campo para o chamador
		encoderConfig.StacktraceKey = "stacktrace"              // Nome do campo para o stacktrace
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	core := zapcore.NewCore(encoder, output, logLevel)

	callerSkip := cfg.CallerSkip
	if callerSkip <= 0 {
		callerSkip = 1 // Padrão para pular os métodos do próprio logger
	}

	loggerOpts := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(callerSkip),
		zap.AddStacktrace(zapcore.ErrorLevel), // Inclui stack trace em erros ou níveis mais altos
	}
	zapLogger := zap.New(core, loggerOpts...)
	sugaredLogger := zapLogger.Sugar()

	l := &Logger{
		zap:   zapLogger,
		sugar: sugaredLogger,
	}

	l.sugar.Info("Logger successfully initialized.")
	if isDevelopmentMode {
		l.sugar.Debug("Development mode ENABLED (console output).")
	} else {
		l.sugar.Debug("Production mode ENABLED (JSON output).")
	}
	l.sugar.Debugf("Log level: %s", logLevel)

	return l, nil
}

// Sync garante que todos os logs em buffer sejam escritos.
// Chame antes de sair da aplicação.
func (l *Logger) Sync() error {
	if l == nil || l.zap == nil {
		return nil
	}
	return l.zap.Sync()
}

// GetZap retorna a instância do logger Zap principal.
func (l *Logger) GetZap() *zap.Logger {
	return l.zap
}

// GetSugar retorna a instância do SugaredLogger.
func (l *Logger) GetSugar() *zap.SugaredLogger {
	return l.sugar
}

// Debug logs a message at DebugLevel.
func (l *Logger) Debug(args ...interface{}) {
	if l == nil || l.sugar == nil {
		return
	}
	l.sugar.Debug(args...)
}

// Debugf logs a formatted message at DebugLevel.
func (l *Logger) Debugf(template string, args ...interface{}) {
	if l == nil || l.sugar == nil {
		return
	}
	l.sugar.Debugf(template, args...)
}

// Info logs a message at InfoLevel.
func (l *Logger) Info(args ...interface{}) {
	if l == nil || l.sugar == nil {
		return
	}
	l.sugar.Info(args...)
}

// Infof logs a formatted message at InfoLevel.
func (l *Logger) Infof(template string, args ...interface{}) {
	if l == nil || l.sugar == nil {
		return
	}
	l.sugar.Infof(template, args...)
}

// Warn logs a message at WarnLevel.
func (l *Logger) Warn(args ...interface{}) {
	if l == nil || l.sugar == nil {
		return
	}
	l.sugar.Warn(args...)
}

// Warnf logs a formatted message at WarnLevel.
func (l *Logger) Warnf(template string, args ...interface{}) {
	if l == nil || l.sugar == nil {
		return
	}
	l.sugar.Warnf(template, args...)
}

// Error logs a message at ErrorLevel.
func (l *Logger) Error(args ...interface{}) {
	if l == nil || l.sugar == nil {
		return
	}
	l.sugar.Error(args...)
}

// Errorf logs a formatted message at ErrorLevel.
func (l *Logger) Errorf(template string, args ...interface{}) {
	if l == nil || l.sugar == nil {
		return
	}
	l.sugar.Errorf(template, args...)
}

// Fatal logs a message at FatalLevel then calls os.Exit(1).
func (l *Logger) Fatal(args ...interface{}) {
	if l == nil || l.sugar == nil {
		return
	}
	l.sugar.Fatal(args...)
}

// Fatalf logs a formatted message at FatalLevel then calls os.Exit(1).
func (l *Logger) Fatalf(template string, args ...interface{}) {
	if l == nil || l.sugar == nil {
		return
	}
	l.sugar.Fatalf(template, args...)
}

// Panic logs a message at PanicLevel then panics.
func (l *Logger) Panic(args ...interface{}) {
	if l == nil || l.sugar == nil {
		return
	}
	l.sugar.Panic(args...)
}

// Panicf logs a formatted message at PanicLevel then panics.
func (l *Logger) Panicf(template string, args ...interface{}) {
	if l == nil || l.sugar == nil {
		return
	}
	l.sugar.Panicf(template, args...)
}
