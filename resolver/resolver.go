package resolver

import (
	"errors"
	"math"

	"gonum.org/v1/gonum/mat"
)

// RealRoots calculates the real roots of a polynomial.
func RealRoots(degree int, variables []float64) ([]float64, error) {
	// Error if the number of coefficients doesn't match the degree + 1
	if degree+1 != len(variables) {
		return nil, errors.New("the number of variables must be degree + 1")
	}

	roots, err := findPolynomialRoots(variables)
	if err != nil {
		return nil, err
	}

	var realRoots []float64
	for _, root := range roots {
		if math.Abs(imag(root)) == 0 {
			realRoots = append(realRoots, real(root))
		}
	}

	return realRoots, nil
}

// findPolynomialRoots finds the roots of a polynomial using Eigenvalue decomposition.
func findPolynomialRoots(coeffs []float64) ([]complex128, error) {
	n := len(coeffs) - 1
	companion := mat.NewDense(n, n, nil)
	for i := 1; i < n; i++ {
		companion.Set(i, i-1, 1)
	}
	for i := 0; i < n; i++ {
		companion.Set(i, n-1, -coeffs[i]/coeffs[n])
	}

	var eig mat.Eigen
	ok := eig.Factorize(companion, mat.EigenRight)

	// If Eigenvalue decomposition fails, return an error
	if !ok {
		return nil, errors.New("eigenvalue decomposition failed")
	}

	return eig.Values(nil), nil
}

// GetIRR calculates the Internal Rate of Return.
func GetIRR(v float64) float64 {
	if v == 0 {
		return math.NaN()
	}
	d := 1 - v
	irr := d / (1 - d)
	return irr * 100
}
