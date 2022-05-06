package apptrace

import (
	"github.com/newrelic/go-agent/v3/integrations/logcontext"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

func NewRelicEncoderConfig(c *zapcore.EncoderConfig) {
	c.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendInt64(t.Unix())
	}
	c.TimeKey = logcontext.KeyTimestamp
	c.LevelKey = logcontext.KeyLevel
	c.MessageKey = logcontext.KeyMessage
}

func NewRelicOptions(options ...zap.Option) []zap.Option {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	base := []zap.Option{
		zap.WithCaller(false),
		zap.AddStacktrace(zapcore.PanicLevel),
		zap.Fields(
			zap.String(logcontext.KeyHostname, hostname),
		),
	}
	return append(base, options...)
}
