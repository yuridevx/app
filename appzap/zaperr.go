package appzap

import (
	"errors"
	"github.com/yuridevx/app/apperr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapErr struct {
	err    error
	appErr *apperr.AppError
}

func (e *zapErr) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("error", e.err.Error())
	for k, v := range e.appErr.Attributes {
		err := encoder.AddReflected(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func AppErr(err error) zap.Field {
	if err == nil {
		return zap.Skip()
	}
	var appErr *apperr.AppError
	if errors.As(err, &appErr) {
		return zap.Field{
			Type:      zapcore.InlineMarshalerType,
			Interface: &zapErr{err: err, appErr: appErr},
		}
	}
	return zap.Error(err)
}
