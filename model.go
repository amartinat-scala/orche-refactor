package main

import "sync"

type Simulation struct {
    Id     int
    Cluster *Cluster
    // ... other fields
}

type Cluster struct {
    Id     int
    Mutex  sync.Mutex
    IsClusterReady bool  // Added this field to track the cluster's readiness state
    // ... other fields
}

func (c *Cluster) IsReady() bool {
    // Placeholder: Returns the readiness state of the cluster.
    // Update this with actual logic if needed.
    return c.IsClusterReady
}

func (c *Cluster) SetReady() {
    // Placeholder: Sets the cluster as ready.
    // Update this with actual logic if needed.
    c.IsClusterReady = true
}
