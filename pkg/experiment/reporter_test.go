package experiment

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestLoggingReporter_Error(t *testing.T) {
	reporter := LoggingReporter{}
	reporter.Error(context.Background(), "my-test-error",  uuid.New().String(), errors.New("failed"), time.Duration(100), time.Duration(200))
}

func TestLoggingReporter_Failure(t *testing.T) {
	reporter := LoggingReporter{}
	reporter.Failure(context.Background(), "my-test-failure",  uuid.New().String(), errors.New("failed"), time.Duration(100), time.Duration(200))
}

func TestLoggingReporter_Success(t *testing.T) {
	reporter := LoggingReporter{}
	reporter.Success(context.Background(), "my-test-success",  uuid.New().String(), time.Duration(100), time.Duration(200))
}


