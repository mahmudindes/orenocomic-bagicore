package logger

import "context"

type ctxFields struct{}

func ContextWith(ctx context.Context, keysAndValues ...any) context.Context {
	keysAndValues0, ok := ctx.Value(ctxFields{}).([]any)
	if ok {
		keysAndValues = append(keysAndValues0, keysAndValues...)
	}
	return context.WithValue(ctx, ctxFields{}, keysAndValues)
}

func (l Logger) WithContext(ctx context.Context) Logger {
	keysAndValues, _ := ctx.Value(ctxFields{}).([]any)
	return l.With(keysAndValues...)
}
