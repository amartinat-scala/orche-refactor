package main

import "sync"

type Simulation struct {
    Id     int
    Cluster *Cluster
}

type Cluster struct {
    Id     int
    Mutex  sync.Mutex
    IsClusterReady bool  
}

func (c *Cluster) IsReady() bool {
    return c.IsClusterReady
}

func (c *Cluster) SetReady() {
    c.IsClusterReady = true
}
