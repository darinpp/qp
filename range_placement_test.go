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

func TestFailureDist(t *testing.T) {
	n0 := &Node{Locality: []string{"aws", "west", "az1"}}
	n1 := &Node{Locality: []string{"aws", "west", "az1"}}
	n2 := &Node{Locality: []string{"aws", "west", "az2"}}
	n3 := &Node{Locality: []string{"aws", "west", "az3"}}
	n4 := &Node{Locality: []string{"aws", "east", "az1"}}
	n5 := &Node{Locality: []string{"gcp", "east", "az1"}}
	n6 := &Node{Locality: []string{"gcp", "east", "az1"}}
	n7 := &Node{Locality: []string{"gcp", "east", "az1"}}
	n8 := &Node{Locality: []string{"azure", "east", "az1"}}
	n9 := &Node{Locality: []string{"azure", "west", "az1"}}

	tests := []struct {
		Nodes       []*Node
		FailureProb float64
	}{
		{[]*Node{n0}, 0.087434},

		{[]*Node{n0, n1}, 0.133062}, // Same Cloud, DC and AZ
		{[]*Node{n0, n2}, 0.141731}, // Same Cloud, DC but different AZs
		{[]*Node{n0, n9}, 0.158811}, // Different Clouds

		{[]*Node{n0, n1, n2}, 0.047235}, // Same Cloud, Same DC with two different AZs
		{[]*Node{n1, n2, n3}, 0.039598}, // Same Cloud, Same DC with three different AZs
		{[]*Node{n1, n2, n4}, 0.040619}, // Same Cloud, Two different DCs, three different AZs
		{[]*Node{n5, n6, n7}, 0.046368}, // Same Cloud, Same DC, same AZ
		{[]*Node{n1, n4, n5}, 0.034390}, // Two Clouds with three different DCs
		{[]*Node{n1, n5, n8}, 0.027222}, // Three different Clouds

		{[]*Node{n0, n1, n2, n3}, 0.055411}, // Same Cloud, DC and AZ
		{[]*Node{n1, n4, n5, n6}, 0.067484}, // Same Cloud, DC but different AZs
		{[]*Node{n1, n5, n8, n9}, 0.048655}, // Different Clouds

		{[]*Node{n0, n1, n5, n6, n8}, 0.021900}, // Three different Clouds, three different DCs
		{[]*Node{n0, n1, n5, n8, n9}, 0.019888}, // Three different Clouds, four different DCs
	}

	for _, test := range tests {
		rp := RangePlacement{Nodes: test.Nodes}
		rp.Recompute()
		assert.InEpsilon(t, test.FailureProb, rp.FailureProb, 1e-4)
	}
}
