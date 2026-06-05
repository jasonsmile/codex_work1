package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

const requestIDKey = "requestID"

var (
	accessWriter  *rotatingWriter
	warningWriter *rotatingWriter
	errorWriter   *rotatingWriter
	requestSeq    uint64
)

type rotatingWriter struct {
	baseDir string
	name    string
	mu      sync.Mutex
	hourKey string
	file    *os.File
}

type Field struct {
	Key   string
	Value any
}

func Init(baseDir string) error {
	if strings.TrimSpace(baseDir) == "" {
		baseDir = "log"
	}
	accessWriter = newRotatingWriter(baseDir, "access")
	warningWriter = newRotatingWriter(baseDir, "warning")
	errorWriter = newRotatingWriter(baseDir, "error")
	return os.MkdirAll(baseDir, 0755)
}

func Close() {
	for _, writer := range []*rotatingWriter{accessWriter, warningWriter, errorWriter} {
		if writer != nil {
			writer.close()
		}
	}
}

func Access(message string, fields ...Field) {
	write(accessWriter, "access", message, fields...)
}

func Warning(message string, fields ...Field) {
	write(warningWriter, "warning", message, fields...)
}

func Error(message string, fields ...Field) {
	write(errorWriter, "error", message, fields...)
}

func RequestID(c *gin.Context) string {
	if value, ok := c.Get(requestIDKey); ok {
		if requestID, ok := value.(string); ok {
			return requestID
		}
	}
	return ""
}

func AccessMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := strings.TrimSpace(c.GetHeader("X-Request-ID"))
		if requestID == "" {
			requestID = newRequestID()
		}
		c.Set(requestIDKey, requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()

		Access("request completed",
			Field{"request_id", requestID},
			Field{"method", c.Request.Method},
			Field{"path", c.Request.URL.Path},
			Field{"query", c.Request.URL.RawQuery},
			Field{"status", c.Writer.Status()},
			Field{"duration_ms", time.Since(start).Milliseconds()},
			Field{"client_ip", c.ClientIP()},
			Field{"username", contextString(c, "username")},
			Field{"user_agent", c.Request.UserAgent()},
		)
	}
}

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				Error("panic recovered",
					Field{"request_id", RequestID(c)},
					Field{"method", c.Request.Method},
					Field{"path", c.Request.URL.Path},
					Field{"client_ip", c.ClientIP()},
					Field{"panic", fmt.Sprint(recovered)},
					Field{"stack", string(debug.Stack())},
				)
				c.AbortWithStatusJSON(500, gin.H{
					"code":         500,
					"message":      "服务内部错误",
					"errorMessage": "服务内部错误",
					"data":         nil,
				})
			}
		}()
		c.Next()
	}
}

func newRotatingWriter(baseDir string, name string) *rotatingWriter {
	return &rotatingWriter{baseDir: baseDir, name: name}
}

func (w *rotatingWriter) Write(p []byte) (int, error) {
	if w == nil {
		return os.Stdout.Write(p)
	}

	now := time.Now()
	hourKey := now.Format("2006-01-02-15")
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file == nil || w.hourKey != hourKey {
		if err := w.rotate(now, hourKey); err != nil {
			return 0, err
		}
	}
	n, err := w.file.Write(p)
	if err != nil {
		return n, err
	}
	if syncErr := w.file.Sync(); syncErr != nil {
		return n, syncErr
	}
	return n, nil
}

func (w *rotatingWriter) close() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		_ = w.file.Sync()
		_ = w.file.Close()
		w.file = nil
	}
}

func (w *rotatingWriter) rotate(now time.Time, hourKey string) error {
	if w.file != nil {
		_ = w.file.Close()
		w.file = nil
	}

	dir := filepath.Join(w.baseDir, now.Format("2006-01-02"))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(dir, fmt.Sprintf("%s_%s.log", w.name, now.Format("15")))
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	w.file = file
	w.hourKey = hourKey
	return nil
}

func write(writer io.Writer, level string, message string, fields ...Field) {
	line := strings.Builder{}
	line.WriteString(time.Now().Format(time.RFC3339))
	line.WriteString(" level=")
	line.WriteString(level)
	line.WriteString(" message=")
	line.WriteString(quoteValue(message))

	for _, field := range fields {
		line.WriteByte(' ')
		line.WriteString(field.Key)
		line.WriteByte('=')
		line.WriteString(quoteValue(fmt.Sprint(field.Value)))
	}
	line.WriteByte('\n')

	if writer == nil {
		writer = os.Stdout
	}
	if _, err := writer.Write([]byte(line.String())); err != nil && writer != os.Stderr {
		_, _ = os.Stderr.Write([]byte(line.String()))
	}
}

func quoteValue(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "\n", "\\n")
	value = strings.ReplaceAll(value, "\r", "\\r")
	value = strings.ReplaceAll(value, `"`, `\"`)
	return `"` + value + `"`
}

func contextString(c *gin.Context, key string) string {
	value, ok := c.Get(key)
	if !ok {
		return ""
	}
	return fmt.Sprint(value)
}

func newRequestID() string {
	seq := atomic.AddUint64(&requestSeq, 1)
	return fmt.Sprintf("%d-%06d", time.Now().UnixNano(), seq)
}
