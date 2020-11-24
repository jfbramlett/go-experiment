[![Actions Status](https://github.com/jfbramlett/go-template/workflows/Go/badge.svg)](https://github.com/jfbramlett/go-template/actions)

# Go-Experiment
This is a simple library which can be used to run an "experiment" in the code for validating an alternative an approach.
It was designed to be used to experiment with using a readstore to serve data from a cache rather than from DB reads.


To use the framework you need to first have some framework for reporting results (see the interface experiment.Reporter). 
An example might be something that reports the results to Prometheus or DataDog for graphing. The framework publishes counts
of the number of experiments that were successful (meaning the data matched between the two approaches), the number that failed,
and the number that resulted in an error. In addition it publishes the duration of each path. 

Experiments are also assigned a unique UUID which is added to the logger allowing correlation of log messages across an experiment. 

The experiment itself is pretty simple:

```
experiment := NewExperiment("catalog.asset.readstore", func(ctx context.Context) (interface{}, error) {
    // this is the reference function, the main path - it is executed in the main thread
}, 
func(ctx context.Context) (interface{}, error) {
    // this is the experiment function, it is executed in a goroutine and is expected to produce the same
    // output as the reference function using some other means
}, 
func(ref interface{}, experiment interface{}) error {
    // this is a function that is used to compare the results of the reference and experiment to determine
    // if the experiment is producing the same result. It is run in its own goroutine after both functions
    // complete
},
reporter)

// this runs the experiment returning the results of the reference function
refData, err := experiment.Run(ctx)

// convert refData to our return type since it comes out as an interface{}
...
```

By using goroutines there should be minimal overhead added to the main processing.
