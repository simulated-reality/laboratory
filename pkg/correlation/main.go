package correlation

import (
	"math"

	"github.com/ready-steady/simulation/system"
)

func Compute(application *system.Application, length float64) []float64 {
	tc := uint32(len(application.Tasks))

	distance := measure(application)
	C := make([]float64, tc*tc)

	for i := uint32(0); i < tc; i++ {
		C[i*tc+i] = 1
		for j := i + 1; j < tc; j++ {
			d := distance[i*tc+j]
			C[j*tc+i] = math.Exp(-d * d / (length * length))
			C[i*tc+j] = C[j*tc+i]
		}
	}

	return C
}

func measure(application *system.Application) []float64 {
	tc := uint32(len(application.Tasks))

	depth := explore(application)

	index := make([]uint16, tc)
	count := make([]uint16, tc)
	for i, d := range depth {
		index[i] = count[d]
		count[d]++
	}

	distance := make([]float64, tc*tc)

	for i := uint32(0); i < tc; i++ {
		for j := i + 1; j < tc; j++ {
			xi := float64(index[i]) - float64(count[depth[i]])/2.0
			yi := float64(depth[i])

			xj := float64(index[j]) - float64(count[depth[j]])/2.0
			yj := float64(depth[j])

			distance[j*tc+i] = math.Sqrt((xi-xj)*(xi-xj) + (yi-yj)*(yi-yj))
			distance[i*tc+j] = distance[j*tc+i]
		}
	}

	return distance
}

func explore(application *system.Application) []uint16 {
	tc := uint16(len(application.Tasks))
	depth := make([]uint16, tc)

	for _, l := range application.Leafs() {
		ascend(application, depth, l)
	}

	return depth
}

func ascend(application *system.Application, depth []uint16, f uint16) {
	max := uint16(0)

	for _, p := range application.Tasks[f].Parents {
		if depth[p] == 0 {
			ascend(application, depth, p)
		}
		if max < depth[p]+1 {
			max = depth[p] + 1
		}
	}

	depth[f] = max
}
