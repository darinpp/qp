// Copyright 2020 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package qp

import (
	"fmt"
	"math"
)

const FailureDistEpsilon = 1e-10

// Represents the distribution of probabilities of failure of 1,2,...,n nodes
// For example - evaluating the risk of a 3 replica range placement, requires
// that we calculate the P(zero node failures) and P(1 node failure)
// These are stored in Prob[0], Prob[1]
// In this case the probability of 2+ node failure is (1 - Prob[0] - Prob[1])
type FailureDist struct {
	prob []float64
}

func MakeFailureDist(prob []float64) (FailureDist, error) {
	result := FailureDist{}
	total := .0
	for i, f := range prob {
		if f < 0 {
			return result, fmt.Errorf("element %d value %f not non-negative", i, f)
		} else if f > 1.0 && !(math.Abs(f-1.0) < FailureDistEpsilon) {
			return result, fmt.Errorf("element %d value %f not 1 or less", i, f)
		}
		total += f
	}
	if total > 1.0 && !(math.Abs(total-1.0) < FailureDistEpsilon) {
		return result, fmt.Errorf("total expected to be 1 or less but got %f", total)
	}
	return FailureDist{
		prob: append([]float64(nil), prob...),
	}, nil
}

func (fd FailureDist) IsEmpty() bool {
	return len(fd.prob) == 0
}

func (fd FailureDist) Prob() []float64 {
	return append([]float64(nil), fd.prob...)
}

// Compute the convolution of the current failure distribution and the other.
// Returns the resulting failure distribution.
// The current and the other failure distributions should have the same size.
// The resulting distribution will be compressed to the size of current and
// other
func (fd FailureDist) Convolve(other FailureDist) (FailureDist, error) {
	n := len(fd.prob)
	if n != len(other.prob) {
		return FailureDist{},
			fmt.Errorf("size mismatch %d != %d", n, len(other.prob))
	}
	result := make([]float64, n)

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			k := i + j
			if k < n {
				result[k] += fd.prob[i] * other.prob[j]
			}
		}
	}

	return MakeFailureDist(result)
}
