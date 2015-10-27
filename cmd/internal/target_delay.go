package internal

import (
	"github.com/ready-steady/adapt"
	"github.com/simulated-reality/laboratory/internal/config"
	"github.com/simulated-reality/laboratory/internal/problem"
)

type delayTarget struct {
	problem *problem.Problem
	config  *config.Target
}

func newDelayTarget(p *problem.Problem, c *config.Target) *delayTarget {
	return &delayTarget{
		problem: p,
		config:  c,
	}
}

func (t *delayTarget) String() string {
	return String(t)
}

func (t *delayTarget) Dimensions() (uint, uint) {
	return uint(t.problem.Model.Len()), 2
}

func (t *delayTarget) Compute(node []float64, value []float64) {
	s, m := t.problem.System, t.problem.Model

	value[0] = s.ComputeSchedule(m.Transform(node)).Span
	value[1] = value[0] * value[0]
}

func (t *delayTarget) Monitor(progress *adapt.Progress) {
	if t.config.Verbose {
		Monitor(t, progress)
	}
}

func (t *delayTarget) Score(location *adapt.Location, progress *adapt.Progress) float64 {
	return Score(t, t.config, location, progress)
}
