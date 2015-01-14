package main

// #include <string.h>
import "C"

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/ready-steady/linal/matrix"
	"github.com/ready-steady/persim/power"
	"github.com/ready-steady/probability/gaussian"
	"github.com/ready-steady/tempan/expint"

	"../../pkg/solver"
)

type tempTarget struct {
	problem *problem

	ic uint32 // inputs
	oc uint32 // outputs
	sc uint32 // steps

	power       *power.Distributer
	temperature *expint.Solver
}

func newTempTarget(p *problem) (target, error) {
	c := &p.config

	if c.TempAnalysis.TimeStep <= 0 {
		return nil, errors.New("the time step is invalid")
	}

	power := power.New(p.platform, p.application, c.TempAnalysis.TimeStep)
	temperature, err := expint.New(expint.Config(c.TempAnalysis))
	if err != nil {
		return nil, err
	}

	target := &tempTarget{
		problem: p,

		ic: 1 + p.zc, // +1 for time
		oc: uint32(len(c.CoreIndex)),
		sc: uint32(p.schedule.Span / c.TempAnalysis.TimeStep),

		power:       power,
		temperature: temperature,
	}

	return target, nil
}

func (t *tempTarget) String() string {
	return fmt.Sprintf("Target{inputs: %d, outputs: %d}", t.ic, t.oc)
}

func (t *tempTarget) InputsOutputs() (uint32, uint32) {
	return t.ic, t.oc
}

func (t *tempTarget) Serve(jobs <-chan solver.Job) {
	p := t.problem
	c := &p.config

	cc, uc, zc, oc, sc := p.cc, p.uc, p.zc, t.oc, t.sc
	coreIndex := c.CoreIndex

	g := gaussian.New(0, 1)
	m := p.marginals

	P := make([]float64, cc*sc)
	S := make([]float64, t.temperature.Nodes*sc)

	z := make([]float64, zc)
	u := make([]float64, uc)
	d := make([]float64, p.tc)

	for job := range jobs {
		Q := job.Data

		if Q == nil {
			Q = make([]float64, cc*sc)

			// Independent uniform to independent Gaussian
			for i := uint32(0); i < zc; i++ {
				z[i] = g.InvCDF(job.Node[1+i]) // +1 for time
			}

			// Independent Gaussian to dependent Gaussian
			matrix.Multiply(p.transform, z, u, uc, zc, 1)

			// Dependent Gaussian to dependent uniform to dependent target
			for i, tid := range c.TaskIndex {
				d[tid] = m[i].InvCDF(g.CDF(u[i]))
			}

			// FIXME: Bad, bad, bad!
			C.memset(unsafe.Pointer(&P[0]), 0, C.size_t(8*cc*sc))

			t.power.Compute(p.time.Recompute(p.schedule, d), P, sc)
			t.temperature.ComputeTransient(P, Q, S, sc)
		}

		sid := uint32(job.Node[0] * float64(sc-1))
		for i := uint32(0); i < oc; i++ {
			job.Value[i] = Q[sid*cc+uint32(coreIndex[i])]
		}

		job.Done <- solver.Result{job.Key, Q}
	}
}
