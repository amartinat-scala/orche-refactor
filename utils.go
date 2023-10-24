package main

import (
    "math/rand"
    "time"	
	"log"
)

func DummyBootCluster(cluster *Cluster) {	
    log.Printf("Booting cluster with ID %d...", cluster.Id)
    time.Sleep(time.Second * time.Duration(60+rand.Intn(30)))
}

func DummySetUpEnvironment(cluster *Cluster) {
	log.Printf("Setting up environment for cluster with ID %d...", cluster.Id)
    time.Sleep(time.Second * time.Duration(rand.Intn(30)))
}

func DummyConfigureCluster(cluster *Cluster) {
	log.Printf("Configuring cluster with ID %d...", cluster.Id)
    time.Sleep(time.Second * time.Duration(rand.Intn(30)))
}

func DummyExecuteSimulation(simulation *Simulation) {
	log.Printf("Executing simulation with ID %d...", simulation.Id)
    time.Sleep(time.Second * time.Duration(10+rand.Intn(20)))
}

func DummyStartDSP(simulation *Simulation) {
	log.Printf("Starting DSP for simulation with ID %d...", simulation.Id)
    time.Sleep(time.Second * time.Duration(5+rand.Intn(10)))
}

func getSimulationFromDB(simulationId int) *Simulation {
	log.Printf("Retrieving simulation with ID %d from DB...", simulationId)
    return &Simulation{
        Id: simulationId,
        Cluster: &Cluster{
            Id: simulationId,
        },
    }
}

func getClusterFromDB(clusterId int) *Cluster {
	log.Printf("Retrieving cluster with ID %d from DB...", clusterId)
    return &Cluster{
        Id: clusterId,
    }
}