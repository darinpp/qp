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
	"testing"

	"github.com/draffensperger/golp/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestMakeFailureDist(t *testing.T) {
	tests := []struct {
		dist, res []float64
		err       string
	}{
		{[]float64{.95, .05}, []float64{.95, .05}, ""},
		{[]float64{}, []float64{}, ""},
		{[]float64{-.1}, []float64{-.1}, "element 0 value -0.100000 not non-negative"},
		{[]float64{1.1}, []float64{1.1}, "element 0 value 1.100000 not 1 or less"},
	}
	for _, test := range tests {
		dist, err := MakeFailureDist(test.dist)
		if len(test.err) > 0 {
			assert.EqualError(t, err, test.err)
			continue
		}
		res := dist.Prob()
		n := len(res)
		assert.Equal(t, n, len(test.res))
		for i := 0; i < n; i++ {
			assert.InEpsilon(t, test.res[i], res[i], FailureDistEpsilon)
		}
	}
}

func TestConvolve(t *testing.T) {
	tests := []struct {
		dist1, dist2, res []float64
		err               string
	}{
		{[]float64{.95, .05}, []float64{.95, .05}, []float64{.9025, 0.095}, ""},
		{[]float64{.9025, 0.095}, []float64{.95, .05}, []float64{.857375, 0.135375}, ""},
		{[]float64{}, []float64{}, []float64{}, ""},
		{[]float64{.5}, []float64{.5}, []float64{.25}, ""},
		{[]float64{.1}, []float64{.5}, []float64{.05}, ""},
		{[]float64{.1, .2}, []float64{.5}, []float64{.05}, "size mismatch 2 != 1"},
	}

	for _, test := range tests {
		dist1, err1 := MakeFailureDist(test.dist1)
		dist2, err2 := MakeFailureDist(test.dist2)
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		resDist, err := dist1.Convolve(dist2)
		if len(test.err) > 0 {
			assert.EqualError(t, err, test.err)
			continue
		}
		res := resDist.Prob()
		n := len(res)
		assert.Equal(t, n, len(test.res))
		for i := 0; i < n; i++ {
			assert.InEpsilon(t, test.res[i], res[i], FailureDistEpsilon)
		}
	}
}
