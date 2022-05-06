package apptrace

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

type mapLog sync.Map

func (m *mapLog) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	syncMap := (*sync.Map)(m)
	syncMap.Range(func(key, value interface{}) bool {
		keyS, ok := key.(string)
		if !ok {
			return true
		}
		_ = encoder.AddReflected(keyS, value)
		return true
	})
	return nil
}

var _ zapcore.ObjectMarshaler = (*mapLog)(nil)

func InlineSyncMap(m *sync.Map) zap.Field {
	log := (*mapLog)(m)
	return zap.Field{
		Type:      zapcore.InlineMarshalerType,
		Interface: log,
	}
}
