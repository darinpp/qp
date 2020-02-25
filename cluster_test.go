// Copyright 2020 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.
package qp_test

import (
	"fmt"
	"github.com/darinpp/qp"
	"github.com/draffensperger/golp/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"testing"
)

// This is trying to setup an optimization problem where we are given a
// number of replicas and they have to be spread over the given nodes in a manner
// that makes the number of the replicas on each node be as close as possible to
// the mean number of replicas per node. If there is no restriction on the number
// of replicas each node can host - the problem will be a trivial problem where
// each node will get the mean. So we also add a constraint on the maximum number
// of replicas for some of the nodes. We are then checking that the excess
// replicas that are not places on the constrained nodes are being spread over
// the unconstrained nodes.
func TestOptimizeToMean(t *testing.T) {
	tests := []struct {
			existing []int // existing number of replicas per node
			optimal  []int // expected optimal number of replicas per node
			limit    []int // upper limit on the number of replicas (-1 for no limit)
	} {
		{[]int{1,3,11},[]int{5,5,5},[]int{-1,-1,-1}},
		{[]int{1,5,11},[]int{6,6,6},[]int{-1,-1,-1}},
		{[]int{1,3,11},[]int{1,7,7},[]int{1,-1,-1}},
		{[]int{1,3,11},[]int{3,3,9},[]int{3,3,-1}},
	}
	for _, test := range tests {
		nodes := make([]qp.Node, len(test.existing))
		for i := 0; i < len(nodes); i++ {
			nodes[i].CurrentReplicaCount = test.existing[i]
			nodes[i].OptimalReplicaCount = test.optimal[i]
			nodes[i].ReplicaCountLimit = test.limit[i]
		}

		c := qp.Cluster {Nodes: nodes}
		c.OptimizeToMean()

		for i := 0; i < len(nodes); i++ {
			assert.Equal(t,test.optimal[i],c.Nodes[i].OptimalReplicaCount)
		}
	}
}

// This optimization aims to compute the optimal spread of the
// available data over a number of range choices. Each range is defined
// as an array of 3 nodes. Each node contains the replicas for that range.
// It is also assumed that each node has an implicit probability of failure
// equal to 10%. The probability of failure of each node is assumed to be
// unconditional on the status of the other nodes.
func TestOptimizeAcrossRanges(t *testing.T) {
	tests := []struct {
		ranges [][]int
		limitWeight []float64
		optimalWeight []float64
	} {
		//{
		//	ranges: [][]int{{0,2,4}, {1,3,5}, {0,2,3}, {4,5,3}, {0,1,2}},
		//	optimalWeight: []float64{1,0,0,0,0},
		//},
		//{
		//	ranges: [][]int{{0,2,4}, {1,3,5}, {0,2,3}, {4,5,3}, {0,1,2}},
		//	limitWeight: []float64{.6,.6,.6,.6,.6},
		//	optimalWeight: []float64{.6,.4,0,0,0},
		//},
		{
			ranges: [][]int{{0,2}, {0,3}},
			//limitWeight: []float64{.6,.6},
			optimalWeight: []float64{1,0,0,0,0},
		},
	}

	for _, test := range tests {
		c := qp.Cluster {Ranges: test.ranges, RangeWeightLimits: test.limitWeight}
		c.OptimizeForRisk()

		fmt.Printf("result %+v\n", c.RangeWeights)
	}
}