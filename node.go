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

type Node struct {
	CurrentReplicaCount int // The current number of replicas
	OptimalReplicaCount int // The optimal number of replicas is populated by the optimizer.
	ReplicaCountLimit int // A limit on the maximum number of replicas that can be on this node (-1 for no limit)
}