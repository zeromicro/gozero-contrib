package slogx

import (
	"fmt"
	"log/slog"

	"github.com/zeromicro/go-zero/core/logx"
)

type SlogWriter struct {
	logger *slog.Logger
}

func NewSlogWriter(handler slog.Handler) logx.Writer {
	logger := slog.New(handler)

	return &SlogWriter{
		logger: logger,
	}
}

func (w *SlogWriter) Alert(v interface{}) {
	w.logger.Error(fmt.Sprint(v))
}

func (w *SlogWriter) Close() error {
	return nil
}

func (w *SlogWriter) Debug(v interface{}, fields ...logx.LogField) {
	w.logger.Debug(fmt.Sprint(v), toSlogFields(fields...)...)
}

func (w *SlogWriter) Error(v interface{}, fields ...logx.LogField) {
	w.logger.Error(fmt.Sprint(v), toSlogFields(fields...)...)
}

func (w *SlogWriter) Info(v interface{}, fields ...logx.LogField) {
	w.logger.Info(fmt.Sprint(v), toSlogFields(fields...)...)
}

func (w *SlogWriter) Severe(v interface{}) {
	w.logger.Error(fmt.Sprint(v))
}

func (w *SlogWriter) Slow(v interface{}, fields ...logx.LogField) {
	w.logger.Warn(fmt.Sprint(v), toSlogFields(fields...)...)
}

func (w *SlogWriter) Stack(v interface{}) {
	w.logger.Error(fmt.Sprint(v))
}

func (w *SlogWriter) Stat(v interface{}, fields ...logx.LogField) {
	w.logger.Info(fmt.Sprint(v), toSlogFields(fields...)...)
}

func toSlogFields(fields ...logx.LogField) []interface{} {
	slogFields := make([]interface{}, 0, len(fields) * 2)
	for _, field := range fields {
		slogFields = append(slogFields, field.Key, field.Value)
	}

	return slogFields
}
