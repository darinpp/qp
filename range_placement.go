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

import "fmt"

// This is the tree representation of the localities.
type LocalityTree struct {
	nodeCount    int
	subLocations map[string]*LocalityTree
	nodes        []*Node
}

func MakeLocalityTree() LocalityTree {
	return LocalityTree{
		subLocations: make(map[string]*LocalityTree),
	}
}

// Add a node as a tree leaf. Add any nodes between the root and the leaf that
// don't exist already.
func (lt *LocalityTree) AddNode(n *Node, level int) {
	if len(n.Locality) > level {
		levelName := n.Locality[level]
		subTree, found := lt.subLocations[levelName]
		if !found {
			newTree := MakeLocalityTree()
			subTree = &newTree
			lt.subLocations[levelName] = subTree
		}
		subTree.AddNode(n, level+1)
	} else {
		lt.nodes = append(lt.nodes, n)
	}
	lt.nodeCount++
}

// Compute the failure distribution of the current locality tree.
// The maxFailNodeCount is the maximum failed node count that we want to track.
// For a 3 replica range, we will want to track the probability of nodes failed P(0) and P(1)
// The P(2+) = 1 - P(0) - P(1). So in this case the maxFailNodeCount will equal to 1.
// P(2+) in this case will give us the total probability that a range can become
// unavailable (given that the quorum size is 2 so the range is available only if
// no more than 1 node is down.
func (lt *LocalityTree) FailureDist(maxFailNodeCount int) (FailureDist, error) {
	if len(lt.nodes) == 0 && len(lt.subLocations) == 0 {
		return FailureDist{},
			fmt.Errorf("locality tree without nodes and sub locations")
	}
	result := FailureDist{}
	for _, _ = range lt.nodes {
		nodeFailureProbs := make([]float64, maxFailNodeCount+1)
		// It may make sense to have data here that depends on the type of node.
		nodeFailureProbs[0] = 0.95 // 95% chance that the node won't fail.
		if maxFailNodeCount > 0 {
			nodeFailureProbs[1] = 0.05 // 5% chance that the node will fail.
		}
		nodeFailureDist, err := MakeFailureDist(nodeFailureProbs)
		if err != nil {
			return FailureDist{}, err
		}
		if result.IsEmpty() {
			result = nodeFailureDist
		} else {
			result, err = result.Convolve(nodeFailureDist)
			if err != nil {
				return FailureDist{}, err
			}
		}
	}

	for _, lt := range lt.subLocations {
		subLocationFailureDist, err := lt.FailureDist(maxFailNodeCount)
		if err != nil {
			return FailureDist{}, err
		}
		if result.IsEmpty() {
			result = subLocationFailureDist
		} else {
			result, err = result.Convolve(subLocationFailureDist)
			if err != nil {
				return FailureDist{}, err
			}
		}
	}

	// Adjust the result to reflect the probability of failure of the current location
	prob := result.Prob()
	for i, _ := range prob {
		prob[i] *= 0.99 // Assume 1% of location failure
	}
	if lt.nodeCount < len(prob) {
		prob[lt.nodeCount] += 0.01 // Add the 1% to the relevant bucket if it exists.
	}
	return MakeFailureDist(prob)
}

type RangePlacement struct {
	Nodes []*Node

	FailureProb   float64
	StartWeight   float64
	OptimalWeight float64
}

func (rp *RangePlacement) Recompute() error {
	tree := MakeLocalityTree()
	for _, n := range rp.Nodes {
		tree.AddNode(n, 0)
	}
	dist, err := tree.FailureDist((tree.nodeCount - 1) / 2)
	if err != nil {
		return err
	}
	rp.FailureProb = 1.0
	for _, p := range dist.Prob() {
		rp.FailureProb -= p
	}
	return nil
}
