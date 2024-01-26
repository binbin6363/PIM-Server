package log

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
)

const (
	LoggerTag     = "logger"
	LoggerTraceID = "traceid"
)

var defaultLogger *zap.SugaredLogger

func GetLogger() *zap.SugaredLogger {
	return defaultLogger
}

func InitLogger(fileName string, maxSize, maxBackups, maxAge, level, callerSkip int) {
	writeSyncer := getLogWriter(fileName, maxSize, maxBackups, maxAge)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.Level(level))

	l := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(callerSkip))
	defaultLogger = l.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(fileName string, maxSize, maxBackups, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    maxSize, // mb
		MaxBackups: maxBackups,
		MaxAge:     maxAge, // day
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// Debug .
func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

// Debugf .
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

// Info .
func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

// Infof .
func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

// Warn .
func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

// Warnf .
func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

// Error .
func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

// Errorf .
func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

// Fatal .
func Fatal(args ...interface{}) {
	defaultLogger.Error(args...)
}

// Fatalf .
func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args...)
}

// ======== log with context ========
// getContextLogger .
func getContextLogger(ctx context.Context) *zap.SugaredLogger {
	log := defaultLogger
	if c, ok := ctx.(*gin.Context); ok {
		if l, ok := c.Get(LoggerTag); ok {
			log = l.(*zap.SugaredLogger)
		}
	}
	return log
}

// DebugContext .
func DebugContext(ctx context.Context, args ...interface{}) {
	log := getContextLogger(ctx)
	log.Debug(args...)
}

// DebugContextf .
func DebugContextf(ctx context.Context, format string, args ...interface{}) {
	log := getContextLogger(ctx)
	log.Debugf(format, args...)
}

// InfoContext .
func InfoContext(ctx context.Context, args ...interface{}) {
	log := getContextLogger(ctx)
	log.Info(args...)
}

// InfoContextf .
func InfoContextf(ctx context.Context, format string, args ...interface{}) {
	log := getContextLogger(ctx)
	log.Infof(format, args...)
}

// WarnContext .
func WarnContext(ctx context.Context, args ...interface{}) {
	log := getContextLogger(ctx)
	log.Warn(args...)
}

// WarnContextf .
func WarnContextf(ctx context.Context, format string, args ...interface{}) {
	log := getContextLogger(ctx)
	log.Warnf(format, args...)
}

// ErrorContext .
func ErrorContext(ctx context.Context, args ...interface{}) {
	log := getContextLogger(ctx)
	log.Error(args...)
}

// ErrorContextf .
func ErrorContextf(ctx context.Context, format string, args ...interface{}) {
	log := getContextLogger(ctx)
	log.Errorf(format, args...)
}

// FatalContext .
func FatalContext(ctx context.Context, args ...interface{}) {
	log := getContextLogger(ctx)
	log.Error(args...)
}

// FatalContextf .
func FatalContextf(ctx context.Context, format string, args ...interface{}) {
	log := getContextLogger(ctx)
	log.Fatalf(format, args...)
}

// ======== log with context ========

// ZapGormLogger 给gorm适配使用的logger
type ZapGormLogger struct {
	logger *zap.Logger
}

// LogMode .
func (l ZapGormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

// Info .
func (l ZapGormLogger) Info(ctx context.Context, s string, i ...interface{}) {
	l.logger.Info(s, append([]zap.Field{zap.Any("arguments", i)}, getContextFields(ctx)...)...)
}

// Warn .
func (l ZapGormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	l.logger.Warn(s, append([]zap.Field{zap.Any("arguments", i)}, getContextFields(ctx)...)...)
}

// Error .
func (l ZapGormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	l.logger.Error(s, append([]zap.Field{zap.Any("arguments", i)}, getContextFields(ctx)...)...)
}

// Trace .
func (l ZapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	fields := getContextFields(ctx)
	if err != nil {
		sql, rows := fc()
		l.logger.Error(err.Error(), append(fields, zap.String("sql", sql), zap.Int64("rows", rows), zap.Duration("elapsed", elapsed))...)
	} else {
		sql, rows := fc()
		l.logger.Debug("", append(fields, zap.String("sql", sql), zap.Int64("rows", rows), zap.Duration("elapsed", elapsed))...)
	}
}

// getContextFields .
func getContextFields(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0)
	// 将 ctx 中的字段添加到 fields 切片中，根据需要自定义
	// 例如，你可以将请求 ID、用户 ID 等上下文信息添加到日志中
	// fields = append(fields, zap.String("requestID", getRequestIDFromContext(ctx)))
	return fields
}

// NewZapGormLogger New一个logger实例
func NewZapGormLogger() logger.Interface {
	return ZapGormLogger{logger: GetLogger().Desugar()}
}
