package internal

import (
	"github.com/ready-steady/adapt"
	"github.com/simulated-reality/laboratory/internal/config"
	"github.com/simulated-reality/laboratory/internal/problem"
)

type energyTarget struct {
	problem *problem.Problem
	config  *config.Target
}

func newEnergyTarget(p *problem.Problem, c *config.Target) *energyTarget {
	return &energyTarget{
		problem: p,
		config:  c,
	}
}

func (t *energyTarget) String() string {
	return String(t)
}

func (t *energyTarget) Dimensions() (uint, uint) {
	return uint(t.problem.Model.Len()), 2
}

func (t *energyTarget) Compute(node, value []float64) {
	s, m := t.problem.System, t.problem.Model

	schedule := s.ComputeSchedule(m.Transform(node))
	time, power := s.ComputeTime(schedule), s.DistributePower(schedule)

	value[0] = 0
	for i := range time {
		value[0] += time[i] * power[i]
	}

	value[1] = value[0] * value[0]
}

func (t *energyTarget) Monitor(progress *adapt.Progress) {
	if t.config.Verbose {
		Monitor(t, progress)
	}
}

func (t *energyTarget) Score(location *adapt.Location, progress *adapt.Progress) float64 {
	return Score(t, t.config, location, progress)
}
