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
func NewExperiment(refLoader LoadFunc, expLoader LoadFunc, validator Validation) Experiment {
	return &experiment{uniqueID: uuid.New().String(), validator: validator, refFunc: refLoader, expFunc: expLoader}
}

// Experiment is an interface used for tracking an experiment
type Experiment interface {
	// Run runs the experiment
	Run(ctx context.Context) (interface{}, error)
}

type experiment struct {
	uniqueID    string
	refFunc     LoadFunc
	expFunc     LoadFunc
	validator   Validation
	start       time.Time
	expDuration time.Duration
	refDuration time.Duration
	refData     interface{}
	expData     interface{}
}

func (e *experiment) Run(ctx context.Context) (interface{}, error) {
	e.start = time.Now()
	go func(ctx context.Context) {
		data, err := e.expFunc(ctx)
		e.recordExperimentResult(ctx, data, err)
	}(ctx)

	data, err := e.refFunc(ctx)
	return e.recordReferencetResult(ctx, data, err)
}

func (e *experiment) recordExperimentResult(ctx context.Context, data interface{}, err error) {
	e.expDuration = time.Since(e.start)
	if err != nil {
		e.expData = err
	} else {
		e.expData = data
	}
	if e.expData != nil && e.refData != nil {
		go e.validateExperiment(ctx)
	}
}

func (e *experiment) recordReferencetResult(ctx context.Context, data interface{}, err error) (interface{}, error) {
	e.refDuration = time.Since(e.start)
	if err != nil {
		e.refData = err
	} else {
		e.refData = data
	}
	if e.expData != nil && e.refData != nil {
		go e.validateExperiment(ctx)
	}

	return data, err
}

func (e *experiment) validateExperiment(ctx context.Context) {

}
