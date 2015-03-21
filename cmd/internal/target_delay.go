package internal

type delayTarget struct {
	problem *Problem
	config  *TargetConfig
}

func newDelayTarget(p *Problem, c *TargetConfig) *delayTarget {
	return &delayTarget{
		problem: p,
		config:  c,
	}
}

func (t *delayTarget) String() string {
	return CommonTarget{t}.String()
}

func (t *delayTarget) Dimensions() (uint, uint) {
	return t.problem.nz, 2
}

func (t *delayTarget) Compute(node []float64, value []float64) {
	p := t.problem
	value[0] = p.time.Recompute(p.schedule, p.transform(node)).Span
	value[1] = value[0] * value[0]
}

func (t *delayTarget) Refine(node, surplus []float64, volume float64) float64 {
	Δ := CommonTarget{t}.Refine(node, surplus, volume)
	if Δ <= t.config.Tolerance {
		Δ = 0
	}
	return Δ
}

func (t *delayTarget) Monitor(level, np, na uint) {
	if t.config.Verbose {
		CommonTarget{t}.Monitor(level, np, na)
	}
}

func (t *delayTarget) Generate(ns uint) []float64 {
	return CommonTarget{t}.Generate(ns)
}
