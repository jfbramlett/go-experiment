package experiment

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type LoadFunc func(ctx context.Context) (interface{}, error)

// Validation is an interface for a validator used to check the results of the experiment
type Validation interface {
	Validate(ref interface{}, experiment interface{}) error
}

// NewExperiment starts a new experiment
func NewExperiment(named string, refLoader LoadFunc, expLoader LoadFunc, validator Validation, reporter Reporter) Experiment {
	return &experiment{named: named, uniqueID: uuid.New().String(), validator: validator, reporter: reporter,
		refFunc: refLoader, expFunc: expLoader}
}

// Experiment is an interface used for tracking an experiment
type Experiment interface {
	// Run runs the experiment
	Run(ctx context.Context) (interface{}, error)
}

type experimentResult struct {
	data interface{}
	err  error
	dur  time.Duration
}

func (e *experimentResult) HasError() bool {
	return e.err != nil
}

func newExperimentResult(data interface{}, err error, dur time.Duration) *experimentResult {
	return &experimentResult{data: data, err: err, dur: dur}
}

type experiment struct {
	uniqueID  string
	named     string
	reporter  Reporter
	refFunc   LoadFunc
	expFunc   LoadFunc
	validator Validation
	expResult *experimentResult
	refResult *experimentResult
}

func (e *experiment) Run(ctx context.Context) (interface{}, error) {
	go func(ctx context.Context) {
		start := time.Now()
		data, err := e.expFunc(ctx)
		e.expResult = newExperimentResult(data, err, time.Since(start))
		e.validateExperiment(ctx)
	}(ctx)

	start := time.Now()
	data, err := e.refFunc(ctx)
	e.refResult = newExperimentResult(data, err, time.Since(start))
	e.validateExperiment(ctx)
	return data, err
}

func (e *experiment) validateExperiment(ctx context.Context) {
	if e.refResult != nil && e.expResult != nil {
		go func() {
			if e.expResult.HasError() || e.refResult.HasError() {
				var err error
				if e.expResult.HasError() {
					err = e.expResult.err
				}
				if e.refResult.HasError() {
					err = e.refResult.err
				}
				e.reporter.Error(ctx, e.named, e.uniqueID, err)
				return
			}
			if err := e.validator.Validate(e.refResult.data, e.expResult.data); err != nil {
				e.reporter.Failure(ctx, e.named, e.uniqueID, err)
				return
			}
			e.reporter.Success(ctx, e.named, e.uniqueID)
		}()
	}
}
