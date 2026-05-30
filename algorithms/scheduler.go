package algorithms

import (
	"github.com/talaamm/cpu-scheduler-visualizer/core"
	"github.com/talaamm/cpu-scheduler-visualizer/simulation"
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
