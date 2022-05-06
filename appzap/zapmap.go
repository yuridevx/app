package appzap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type mapLog map[string]interface{}

func (m mapLog) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	for k, v := range m {
		err := encoder.AddReflected(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

var _ zapcore.ObjectMarshaler = (*mapLog)(nil)

func InlineMap(m map[string]interface{}) zap.Field {
	if m == nil {
		return zap.Skip()
	}
	log := (mapLog)(m)
	return zap.Field{
		Type:      zapcore.InlineMarshalerType,
		Interface: log,
	}
}
