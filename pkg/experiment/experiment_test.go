package experiment

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExperimentRun(t *testing.T) {
	t.Run("success experiment", func(t *testing.T) {
		testName := "test-success"
		wg := sync.WaitGroup{}
		wg.Add(1)
		ref := func(ctx context.Context) (interface{}, error) {
			return "success", nil
		}
		exp := func(ctx context.Context) (interface{}, error) {
			wg.Done()
			return "success", nil
		}
		wg.Add(1)
		validator := MockValidator{ValidateMock: func(ref interface{}, experiment interface{}) error {
			wg.Done()
			return nil
		}}
		reporter := NewMockReporter()

		experiment := NewExperiment(testName, ref, exp, validator, reporter)
		_, err := experiment.Run(context.Background())

		assert.NoError(t, err)

		wg.Wait()

		assert.Len(t, reporter.failures, 0)
		assert.Len(t, reporter.errs, 0)
		assert.Len(t, reporter.successes, 1)
	})

	t.Run("validation fails", func(t *testing.T) {
		testName := "test-validation-fails"
		wg := sync.WaitGroup{}
		wg.Add(1)
		ref := func(ctx context.Context) (interface{}, error) {
			return "success", nil
		}
		exp := func(ctx context.Context) (interface{}, error) {
			wg.Done()
			return "success", nil
		}
		wg.Add(1)
		validator := MockValidator{ValidateMock: func(ref interface{}, experiment interface{}) error {
			wg.Done()
			return errors.New("no watch")
		}}
		reporter := NewMockReporter()

		experiment := NewExperiment(testName, ref, exp, validator, reporter)
		_, err := experiment.Run(context.Background())

		assert.NoError(t, err)

		wg.Wait()

		assert.Len(t, reporter.failures, 1)
		assert.Len(t, reporter.errs, 0)
		assert.Len(t, reporter.successes, 0)
	})

	t.Run("experiment error", func(t *testing.T) {
		testName := "test-experiment-error"
		wg := sync.WaitGroup{}
		wg.Add(1)
		ref := func(ctx context.Context) (interface{}, error) {
			return "success", nil
		}
		exp := func(ctx context.Context) (interface{}, error) {
			wg.Done()
			return nil, errors.New("failed")
		}
		reporter := NewMockReporter()

		experiment := NewExperiment(testName, ref, exp, nil, reporter)
		_, err := experiment.Run(context.Background())

		assert.NoError(t, err)

		wg.Wait()

		assert.Len(t, reporter.failures, 0)
		assert.Len(t, reporter.errs, 1)
		assert.Len(t, reporter.successes, 0)
	})
	t.Run("ref error", func(t *testing.T) {
		testName := "test-ref-error"
		wg := sync.WaitGroup{}
		wg.Add(1)
		ref := func(ctx context.Context) (interface{}, error) {
			return nil, errors.New("failed")
		}
		exp := func(ctx context.Context) (interface{}, error) {
			wg.Done()
			return "success", nil
		}
		reporter := NewMockReporter()

		experiment := NewExperiment(testName, ref, exp, nil, reporter)
		_, err := experiment.Run(context.Background())

		assert.Error(t, err)

		wg.Wait()

		assert.Len(t, reporter.failures, 0)
		assert.Len(t, reporter.errs, 1)
		assert.Len(t, reporter.successes, 0)
	})
}

type MockValidator struct {
	ValidateMock func(ref interface{}, experiment interface{}) error
}

func (m MockValidator) Validate(ref interface{}, experiment interface{}) error {
	return m.ValidateMock(ref, experiment)
}

func NewMockReporter() MockReporter {
	return MockReporter{
		errs:      make(map[string]int),
		successes: make(map[string]int),
		failures:  make(map[string]int),
	}
}

type MockReporter struct {
	successes map[string]int
	failures  map[string]int
	errs      map[string]int
}

func (t MockReporter) Success(_ context.Context, named string, _ string) {
	t.successes[named] = t.successes[named] + 1
}

// Failure reports an experiment failure (in which the validation failed)
func (t MockReporter) Failure(_ context.Context, named string, _ string, _ error) {
	t.failures[named] = t.failures[named] + 1
}

// Error reports an error occurring during the experiment
func (t MockReporter) Error(_ context.Context, named string, _ string, _ error) {
	t.errs[named] = t.errs[named] + 1
}
