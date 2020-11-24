package logging

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
)

type requestKey struct{}

var reqLoggerKey = &requestKey{}

// FromContext returns the logger contained in the given context. A second
// return value indicates if a logger was found in the context. If no logger is
// found, a new logger is returned.
func FromContext(ctx context.Context) (*logrus.Entry, bool) {
	v := ctx.Value(reqLoggerKey)
	if s, ok := v.(*logrus.Entry); ok {
		return s, true
	}

	lg := logrus.New()
	lg.Out = os.Stdout
	lg.SetFormatter(&logrus.JSONFormatter{})
	return lg.WithContext(ctx), false
}

// ContextWithLogger returns a copy of the given context with the logger
func ContextWithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, reqLoggerKey, logger)
}

// UpdateInContext returns a copy of the context with a logger that has the
// passed `fields`. If the original context included a logger, it will be used
// to add fields to.
func UpdateInContext(ctx context.Context, fields logrus.Fields) (context.Context, *logrus.Entry) {
	log, _ := FromContext(ctx)
	log = log.WithFields(fields)
	return ContextWithLogger(ctx, log), log
}
