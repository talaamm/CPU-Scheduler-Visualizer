package algorithms

import (
	"../core"
	"../simulation"
)

// all algos will implement same interface

type Scheduler interface {
	Run(processes []core.Process) simulation.SimulationResult
	Name() string
}

/*interface is better and cleaner so that
frontend api can call any algo without caring 
about implementation details or what algo it is
*/