package experiment

import "context"

type Reporter interface {
	// Success reports that an experiment was successful
	Success(ctx context.Context, named string, uuid string)

	// Failure reports an experiment failure (in which the validation failed)
	Failure(ctx context.Context, named string, uuid string, err error)

	// Error reports an error occurring during the experiment
	Error(ctx context.Context, named string, uuid string, err error)
}
