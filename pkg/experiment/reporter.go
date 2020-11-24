package experiment

import (
	"context"
	"github.com/jfbramlett/go-experiment/pkg/logging"
	"github.com/sirupsen/logrus"
	"time"
)

type Reporter interface {
	// Success reports that an experiment was successful
	Success(ctx context.Context, named string, uuid string, refDuration time.Duration, expDuration time.Duration)

	// Failure reports an experiment failure (in which the validation failed)
	Failure(ctx context.Context, named string, uuid string, err error, refDuration time.Duration, expDuration time.Duration)

	// Error reports an error occurring during the experiment
	Error(ctx context.Context, named string, uuid string, err error, refDuration time.Duration, expDuration time.Duration)
}


type LoggingReporter struct {

}

func (l *LoggingReporter) Success(ctx context.Context, named string, uuid string, refDuration time.Duration, expDuration time.Duration) {
	logger, _ := logging.FromContext(ctx)
	logger.WithFields(logrus.Fields{"refDuration": refDuration.Milliseconds(), "expDuration": expDuration.Milliseconds(), "status": "passed"}).Infof("experiment %s with uuid %s passed", named, uuid)
}

func (l *LoggingReporter) Failure(ctx context.Context, named string, uuid string, err error, refDuration time.Duration, expDuration time.Duration) {
	logger, _ := logging.FromContext(ctx)
	logger.WithError(err).WithFields(logrus.Fields{"refDuration": refDuration.Milliseconds(), "expDuration": expDuration.Milliseconds(), "status": "failed"}).Infof("experiment %s with uuid %s failed", named, uuid)
}

func (l *LoggingReporter) Error(ctx context.Context, named string, uuid string, err error, refDuration time.Duration, expDuration time.Duration) {
	logger, _ := logging.FromContext(ctx)
	logger.WithError(err).WithFields(logrus.Fields{"refDuration": refDuration.Milliseconds(), "expDuration": expDuration.Milliseconds(), "status": "error"}).Errorf("experiment %s with uuid %s errored", named, uuid)
}
