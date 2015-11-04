package target

import (
	"github.com/ready-steady/adapt"
	"github.com/turing-complete/laboratory/src/internal/config"
	"github.com/turing-complete/laboratory/src/internal/support"
	"github.com/turing-complete/laboratory/src/internal/system"
	"github.com/turing-complete/laboratory/src/internal/uncertainty"
)

type profile struct {
	system *system.System
	config *config.Target

	coreIndex   []uint
	timeIndex   []float64
	uncertainty uncertainty.Uncertainty
}

func newProfile(system *system.System, config *config.Target) (*profile, error) {
	coreIndex, err := support.ParseNaturalIndex(config.CoreIndex, 0, uint(system.Platform.Len())-1)
	if err != nil {
		return nil, err
	}

	timeIndex, err := support.ParseRealIndex(config.TimeIndex, 0, 1)
	if err != nil {
		return nil, err
	}
	if timeIndex[0] == 0 {
		timeIndex = timeIndex[1:]
	}
	for i := range timeIndex {
		timeIndex[i] *= system.Span()
	}

	uncertainty, err := uncertainty.New(system, &config.Uncertainty)
	if err != nil {
		return nil, err
	}

	return &profile{
		system: system,
		config: config,

		coreIndex:   coreIndex,
		timeIndex:   timeIndex,
		uncertainty: uncertainty,
	}, nil
}

func (t *profile) Dimensions() (uint, uint) {
	nci, nsi := uint(len(t.coreIndex)), uint(len(t.timeIndex))
	return uint(t.uncertainty.Len()), nsi * nci * 2
}

func (t *profile) Compute(node, value []float64) {
	const (
		ε = 1e-10
	)

	schedule := t.system.ComputeSchedule(t.uncertainty.Transform(node))
	P, ΔT, timeIndex := t.system.PartitionPower(schedule, t.timeIndex, ε)
	for i := range timeIndex {
		if timeIndex[i] == 0 {
			panic("the timeline of interest should not contain time 0")
		}
		timeIndex[i]--
	}

	Q := t.system.ComputeTemperature(P, ΔT)

	coreIndex := t.coreIndex
	nc := uint(t.system.Platform.Len())
	nci, nsi := uint(len(coreIndex)), uint(len(timeIndex))

	for i, k := uint(0), uint(0); i < nsi; i++ {
		for j := uint(0); j < nci; j++ {
			value[k] = Q[timeIndex[i]*nc+coreIndex[j]]
			value[k+1] = value[k] * value[k]
			k += 2
		}
	}
}

func (t *profile) Monitor(progress *adapt.Progress) {
	if t.config.Verbose {
		monitor(t, progress)
	}
}

func (t *profile) Score(location *adapt.Location, progress *adapt.Progress) float64 {
	return score(t, t.config, location, progress)
}

func (t *profile) String() string {
	return display(t)
}
