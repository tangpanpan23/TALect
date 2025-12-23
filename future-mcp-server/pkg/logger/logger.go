package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 全局日志实例
var Logger *zap.Logger

// 日志级别常量
const (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
	FatalLevel = zapcore.FatalLevel
)

// 初始化日志
func Init() error {
	config := zap.NewProductionConfig()

	// 设置日志级别
	level := viper.GetString("log.level")
	if level != "" {
		var zapLevel zapcore.Level
		if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
			return fmt.Errorf("invalid log level: %s", level)
		}
		config.Level = zap.NewAtomicLevelAt(zapLevel)
	}

	// 设置日志格式
	format := viper.GetString("log.format")
	if format == "json" {
		config.Encoding = "json"
	} else {
		config.Encoding = "console"
	}

	// 设置输出目标
	output := viper.GetString("log.output")
	if output == "file" {
		// 创建日志目录
		logPath := viper.GetString("log.file_path")
		if logPath == "" {
			logPath = "./logs/future-mcp.log"
		}

		if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}

		config.OutputPaths = []string{logPath}
		config.ErrorOutputPaths = []string{logPath}
	} else {
		config.OutputPaths = []string{"stdout"}
		config.ErrorOutputPaths = []string{"stderr"}
	}

	// 开发环境优化
	if viper.GetString("server.mode") == "debug" {
		config.Development = true
		config.DisableCaller = false
		config.DisableStacktrace = false
	}

	// 创建日志实例
	var err error
	Logger, err = config.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(ErrorLevel),
	)
	if err != nil {
		return fmt.Errorf("failed to build logger: %w", err)
	}

	// 替换全局logger
	zap.ReplaceGlobals(Logger)

	return nil
}

// Debug 调试日志
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// Info 信息日志
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Warn 警告日志
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

// Error 错误日志
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// Fatal 致命错误日志
func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

// WithContext 创建带有上下文的logger
func WithContext(ctx interface{}) *zap.Logger {
	return Logger.With(Any("context", ctx))
}

// WithUser 创建带有用户信息的logger
func WithUser(userID interface{}) *zap.Logger {
	return Logger.With(Any("user_id", userID))
}

// WithRequest 创建带有请求信息的logger
func WithRequest(requestID string) *zap.Logger {
	return Logger.With(Any("request_id", requestID))
}

// WithSession 创建带有会话信息的logger
func WithSession(sessionID string) *zap.Logger {
	return Logger.With(Any("session_id", sessionID))
}

// RequestLogger 请求日志记录器
type RequestLogger struct {
	logger *zap.Logger
}

// NewRequestLogger 创建请求日志记录器
func NewRequestLogger(requestID string) *RequestLogger {
	return &RequestLogger{
		logger: WithRequest(requestID),
	}
}

// LogRequest 记录请求
func (rl *RequestLogger) LogRequest(method, path, userAgent, ip string) {
	rl.logger.Info("request started",
		Any("method", method),
		Any("path", path),
		Any("user_agent", userAgent),
		Any("ip", ip),
	)
}

// LogResponse 记录响应
func (rl *RequestLogger) LogResponse(statusCode int, duration time.Duration, size int64) {
	level := InfoLevel
	if statusCode >= 400 {
		level = WarnLevel
	}
	if statusCode >= 500 {
		level = ErrorLevel
	}

	rl.logger.Log(level, "request completed",
		Any("status_code", statusCode),
		Any("duration_ms", duration.Milliseconds()),
		Any("response_size", size),
	)
}

// MCPLogger MCP协议日志记录器
type MCPLogger struct {
	logger *zap.Logger
}

// NewMCPLogger 创建MCP日志记录器
func NewMCPLogger(requestID string) *MCPLogger {
	return &MCPLogger{
		logger: WithRequest(requestID),
	}
}

// LogMCPRequest 记录MCP请求
func (ml *MCPLogger) LogMCPRequest(method string, params interface{}) {
	ml.logger.Info("mcp request",
		Any("method", method),
		Any("params", fmt.Sprintf("%+v", params)),
	)
}

// LogMCPResponse 记录MCP响应
func (ml *MCPLogger) LogMCPResponse(result interface{}, duration time.Duration) {
	ml.logger.Info("mcp response",
		Any("duration_ms", duration.Milliseconds()),
		Any("result", fmt.Sprintf("%+v", result)),
	)
}

// LogMCPError 记录MCP错误
func (ml *MCPLogger) LogMCPError(err error, code int) {
	ml.logger.Error("mcp error",
		Any("error", err.Error()),
		Any("code", code),
	)
}

// ToolLogger 工具日志记录器
type ToolLogger struct {
	logger *zap.Logger
}

// NewToolLogger 创建工具日志记录器
func NewToolLogger(toolName, requestID string) *ToolLogger {
	return &ToolLogger{
		logger: Logger.With(
			Any("tool", toolName),
			Any("request_id", requestID),
		),
	}
}

// LogToolExecution 记录工具执行
func (tl *ToolLogger) LogToolExecution(args interface{}, startTime time.Time) {
	tl.logger.Info("tool execution started",
		Any("args", fmt.Sprintf("%+v", args)),
		Any("start_time", startTime),
	)
}

// LogToolResult 记录工具结果
func (tl *ToolLogger) LogToolResult(result interface{}, duration time.Duration) {
	tl.logger.Info("tool execution completed",
		Any("duration_ms", duration.Milliseconds()),
		Any("result", fmt.Sprintf("%+v", result)),
	)
}

// LogToolError 记录工具错误
func (tl *ToolLogger) LogToolError(err error, duration time.Duration) {
	tl.logger.Error("tool execution failed",
		Any("error", err.Error()),
		Any("duration_ms", duration.Milliseconds()),
	)
}

// 便捷字段构造器
func Any(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

func String(key, value string) zap.Field {
	return zap.String(key, value)
}

func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}

func Int64(key string, value int64) zap.Field {
	return zap.Int64(key, value)
}

func Float64(key string, value float64) zap.Field {
	return zap.Float64(key, value)
}

func Bool(key string, value bool) zap.Field {
	return zap.Bool(key, value)
}

func Err(err error) zap.Field {
	return zap.Error(err)
}

func Duration(key string, value time.Duration) zap.Field {
	return zap.Duration(key, value)
}

func Time(key string, value time.Time) zap.Field {
	return zap.Time(key, value)
}

// 同步日志输出
func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}
