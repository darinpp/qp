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
	"github.com/darinpp/qp"
	"github.com/draffensperger/golp/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"testing"
)

func TestOptimize(t *testing.T) {
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