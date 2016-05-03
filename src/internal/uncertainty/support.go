package uncertainty

import (
	"errors"
	"math"

	"github.com/ready-steady/linear/matrix"
)

var (
	infinity = math.Inf(1.0)
)

func invert(U, Λ []float64, m uint) ([]float64, error) {
	T := make([]float64, m*m)
	for i := uint(0); i < m; i++ {
		if Λ[i] == 0.0 {
			return nil, errors.New("the matrix is not invertible")
		}
		λ := 1.0 / Λ[i]
		for j := uint(0); j < m; j++ {
			T[j*m+i] = λ * U[i*m+j]
		}
	}

	I := make([]float64, m*m)
	matrix.Multiply(U, T, I, m, m, m)

	return I, nil
}

func inspect(x []float64, m uint) (bool, []float64) {
	ok, signs := true, make([]float64, m)
	for i := uint(0); i < m; i++ {
		switch x[i] {
		case -infinity:
			ok, signs[i] = false, -1.0
		case infinity:
			ok, signs[i] = false, +1.0
		}
	}
	return ok, signs
}

func multiply(A, x []float64, m, n uint) []float64 {
	y := make([]float64, m)
	ok, s := inspect(x, n)
	if ok {
		matrix.Multiply(A, x, y, m, n, 1)
		return y
	}
	for i := uint(0); i < m; i++ {
		fin, inf := 0.0, 0.0
		for j := uint(0); j < n; j++ {
			a := A[j*m+i]
			if a == 0.0 {
				continue
			}
			if s[j] == 0.0 {
				fin += a * x[j]
			} else {
				inf += a * s[j]
			}
		}
		if inf != 0.0 {
			y[i] = inf * infinity
		} else {
			y[i] = fin
		}
	}
	return y
}

func quadratic(A, x []float64, m uint) float64 {
	ok, s := inspect(x, m)
	if ok {
		y := make([]float64, m)
		matrix.Multiply(A, x, y, m, m, 1)
		return matrix.Dot(x, y, m)
	}
	Fin, Inf, INF := 0.0, 0.0, 0.0
	for i := uint(0); i < m; i++ {
		fin, inf := 0.0, 0.0
		for j := uint(0); j < m; j++ {
			a := A[j*m+i]
			if a == 0.0 {
				continue
			}
			if s[j] == 0.0 {
				fin += a * x[j]
			} else {
				inf += a * s[j]
			}
		}
		if s[i] == 0.0 {
			Fin += x[i] * fin
			Inf += x[i] * inf
		} else {
			Inf += s[i] * fin
			INF += s[i] * inf
		}
	}
	if INF != 0.0 {
		return INF * infinity
	} else if Inf != 0.0 {
		return Inf * infinity
	} else {
		return Fin
	}
}
