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
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize/convex/lp"
	"math"
)

type Cluster struct {
	Nodes []Node
	Ranges [][]int
	RangeFailureProbs []float64
	RangeWeightLimits []float64
	RangeWeights []float64
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

func (c *Cluster) ComputeRangeFailureProbs() {
	rangeFailProb := 0.028
	jointFailProbForNSharedNodes := []float64{0.000784, 0.0037, 0.0118, rangeFailProb}

	n := len(c.Ranges)
	c.RangeFailureProbs = make([]float64, n*(n+1)/2)
	// Prepare an array with the coefficients for each variable
	// First populate the weights for each range count variable
	for i := 0; i < n; i++ {
		c.RangeFailureProbs[i] = rangeFailProb
	}
	// Second populate the weights for each variable representing a pair of ranges
	i := n
	for k := 0; k < n; k++ {
		for l := k+1; l < n; l++ {
			sharedNodeCount := 0
			for _, kNode := range c.Ranges[k] {
				for _, lNode := range c.Ranges[l] {
					if kNode == lNode {
						sharedNodeCount++
					}
				}
			}
			c.RangeFailureProbs[i] = -jointFailProbForNSharedNodes[sharedNodeCount]
			i++
		}
	}
}
// The general form of an LP is:
//  minimize cᵀ * x
//  s.t      G * x <= h
//           A * x = b
func (cluster *Cluster) OptimizeForRisk() error {
	cluster.ComputeRangeFailureProbs()

	n := len(cluster.RangeFailureProbs)
	r := len(cluster.Ranges)
	
	c := cluster.RangeFailureProbs

	a := mat.NewDense(1, n, nil)
	for i:= 0; i < r; i++ {
		a.Set(0,i,1)
	}

	b := []float64{1}

	rows := n+r*(r-1)
	if cluster.RangeWeightLimits != nil {
		rows += r
	}
	g := mat.NewDense(rows, n, nil)
	for i := 0; i < n; i++ {
		g.Set(i,i,-1)
	}

	i := 0
	for k := 0; k < r; k++ {
		for l := k+1; l < r; l++ {
			//fmt.Printf("k=%d,l=%d,i=%d\n", k,l,i)
			g.Set(n+2*i, k, -1)
			g.Set(n+2*i, r + i, 1)
			g.Set(n+2*i+1, l, -1)
			g.Set(n+2*i+1, r + i, 1)

			i++
		}
	}
	h := make([]float64, rows)

	if cluster.RangeWeightLimits != nil {
		firstLimitRow := n+2*i
		for i := 0; i < r; i++ {
			g.Set(firstLimitRow + i, i, 1)
			h[firstLimitRow+i] = cluster.RangeWeightLimits[i]
		}
	}
	
	cNew, aNew, bNew := lp.Convert(c,g,h,a,b)
	optF, optX, err := lp.Simplex(cNew, aNew, bNew, 1e-8, nil)
	fmt.Printf("Got optF %f for %+v\n", optF, optX)
	if err != nil {
		return err
	}
	cluster.RangeWeights = make([]float64, r)
	for i := 0; i < r; i++ {
		cluster.RangeWeights[i] = optX[i]
	}
	return nil
}

//  min 0.5 * x G x + g0 x
//  s.t.
//  CE^T x + ce0 = 0
//  CI^T x + ci0 >= 0
func (c *Cluster) OptimizeForRiskUsingQPOptimizer() {
	c.ComputeRangeFailureProbs()

	n := len(c.RangeFailureProbs)
	r := len(c.Ranges)

	// Setup G
	G := NewMatrix()
	G.Set(make([]float64, n*n), uint(n), uint(n))

	// Setup g0
	g0 := NewVector()
	g0.Set(c.RangeFailureProbs, uint(n))

	// Setup CE
	CE := NewMatrix()
	valCE := make([]float64, n)
	for i:= 0; i < r; i++ {
		valCE[i] = 1
	}
	CE.Set(valCE,uint(n),1)

	// Setup ce0
	ce0 := NewVector()
	ce0.Set([]float64{-1},1)

	CI := NewMatrix()
	CI.Resize(0, uint(n), uint(n+r*(r-1)))
	for i := 0; i < n; i++ {
		rowVec := NewVector()
		valRowVec := make([]float64, n)
		valRowVec[i] = 1
		rowVec.Set(valRowVec, uint(n))
		CI.SetColumn(uint(i), rowVec)
	}

	i := 0
	for k := 0; k < r; k++ {
		for l := k+1; l < r; l++ {
			//fmt.Printf("k=%d,l=%d,i=%d\n", k,l,i)
			rowVec1 := NewVector()
			valRowVec1 := make([]float64, n)
			valRowVec1[k] = 1
			valRowVec1[r + i] = -1
			rowVec1.Set(valRowVec1, uint(n))
			CI.SetColumn(uint(n+2*i),rowVec1)

			rowVec2 := NewVector()
			valRowVec2 := make([]float64, n)
			valRowVec2[l] = 1
			valRowVec2[r + i] = -1
			rowVec2.Set(valRowVec2, uint(n))
			CI.SetColumn(uint(n+2*i+1),rowVec2)

			i++
		}
	}

	ci0 := NewVector()
	ci0.Resize(0, uint(n+r*(r-1)))

	// Optimize
	x := NewVector()

	// This doesn't work currently as G needs to be positive definite and in this
	// case it is semi-definite
	res := Solve_quadprog(G, g0, CE, ce0, CI, ci0, x)
	fmt.Printf("return code is %+v and the result size is %d\n", res, x.Size())
	c.RangeWeights = make([]float64, n)

	for i := 0; i < n; i++ {
		c.RangeWeights[i] = math.Round(x.At(uint(i)))
	}
}

//  min 0.5 * x G x + g0 x
//  s.t.
//  CE^T x + ce0 = 0
//  CI^T x + ci0 >= 0
func (c *Cluster) OptimizeForRiskUsingQPOptimizer2() {
	c.ComputeRangeFailureProbs()

	n := len(c.RangeFailureProbs)
	r := len(c.Ranges)

	// Setup G
	G := NewMatrix()
	valG := make([]float64, r*r)

	for i := 0; i < r; i++ {
		for j := 0; j < r; j++ {
		}
	}
	G.Set(valG, uint(r), uint(r))

	// Setup g0
	g0 := NewVector()
	g0.Set(c.RangeFailureProbs, uint(r))

	// Setup CE
	CE := NewMatrix()
	valCE := make([]float64, n)
	for i:= 0; i < r; i++ {
		valCE[i] = 1
	}
	CE.Set(valCE,uint(n),1)

	// Setup ce0
	ce0 := NewVector()
	ce0.Set([]float64{-1},1)

	CI := NewMatrix()
	CI.Resize(0, uint(n), uint(n+r*(r-1)))
	for i := 0; i < n; i++ {
		rowVec := NewVector()
		valRowVec := make([]float64, n)
		valRowVec[i] = 1
		rowVec.Set(valRowVec, uint(n))
		CI.SetColumn(uint(i), rowVec)
	}

	i := 0
	for k := 0; k < r; k++ {
		for l := k+1; l < r; l++ {
			//fmt.Printf("k=%d,l=%d,i=%d\n", k,l,i)
			rowVec1 := NewVector()
			valRowVec1 := make([]float64, n)
			valRowVec1[k] = 1
			valRowVec1[r + i] = -1
			rowVec1.Set(valRowVec1, uint(n))
			CI.SetColumn(uint(n+2*i),rowVec1)

			rowVec2 := NewVector()
			valRowVec2 := make([]float64, n)
			valRowVec2[l] = 1
			valRowVec2[r + i] = -1
			rowVec2.Set(valRowVec2, uint(n))
			CI.SetColumn(uint(n+2*i+1),rowVec2)

			i++
		}
	}

	ci0 := NewVector()
	ci0.Resize(0, uint(n+r*(r-1)))

	// Optimize
	x := NewVector()

	// This doesn't work currently as G needs to be positive definite and in this
	// case it is semi-definite
	res := Solve_quadprog(G, g0, CE, ce0, CI, ci0, x)
	fmt.Printf("return code is %+v and the result size is %d\n", res, x.Size())
	c.RangeWeights = make([]float64, n)

	for i := 0; i < n; i++ {
		c.RangeWeights[i] = math.Round(x.At(uint(i)))
	}
}