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
	"math"
)

type Cluster struct {
	Nodes []Node
}

func (c Cluster) OptimizeToMean() {
	// min (m - X0)^2 + (m - X1)^2 +...+ (m - Xn)^2
	// where X0,X1,...,Xn are the number of replicas on each node
	// m is the mean number of replicas per node
	// subject to
	// X0 + X1 +...+ Xn=n*m
	// Xn>=0
	// -Xn + Limitn >= 0
	// Transform the objective function
	// min -2mX0 - 2mX1 -...-2mXn + X0^2 + X1^2 +...+ Xn^2
	n := len(c.Nodes)

	// Compute m(mean) and t(total)
	m := float64(0)
	t := float64(0)
	for i,node := range c.Nodes {
		t += float64(node.CurrentReplicaCount)
		m = (m*float64(i) + float64(node.CurrentReplicaCount)) / float64(i+1)
	}

	// Setup G
	G := NewMatrix()
	valG := make([]float64, n*n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				valG[i*n+j] = 2
			} else {
				valG[i*n+j] = 0
			}
		}
	}
	G.Set(valG,uint(n),uint(n))

	// Setup g0
	g0 := NewVector()
	valg0 := make([]float64, n)
	for i:= 0; i < n; i++ {
		valg0[i] = -2*m
	}
	g0.Set(valg0, uint(n))

	// Setup CE
	CE := NewMatrix()
	valCE := make([]float64, n)
	for i:= 0; i < n; i++ {
		valCE[i] = 1
	}
	CE.Set(valCE,uint(n),1)

	// Setup ce0
	ce0 := NewVector()
	ce0.Set([]float64{-t},1)

	// CI and ci0 will be set in case if limits
	CI := NewMatrix()
	ci0 := NewVector()
	valCI := make([]float64, n*n)
	valci0 := make([]float64, n)
	for i := 0; i < n; i++ {
		limit :=  float64(c.Nodes[i].ReplicaCountLimit)
		if limit >= 0 {
			valci0[i] = limit
		} else {
			valci0[i] = 0
		}
		for j := 0; j < n; j++ {
			if i == j && limit >= 0 {
				valCI[i*n+j] = -1
			} else {
				valCI[i*n+j] = 0
			}
		}
	}

	CI.Set(valCI,uint(n),uint(n))
	ci0.Set(valci0, uint(n))

	// Optimize
	x := NewVector()

	_ = Solve_quadprog(G, g0, CE, ce0, CI, ci0, x)

	for i := 0; i < n; i++ {
		c.Nodes[i].OptimalReplicaCount = int(math.Round(x.At(uint(i))))
	}
}