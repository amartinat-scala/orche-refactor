package main

import (
    "math/rand"
    "time"	
	"log"
)

func DummyBootCluster(cluster *Cluster) {	
    log.Printf("Booting cluster with ID %d...", cluster.Id)
    // Simulate time taken to boot cluster
    time.Sleep(time.Second * time.Duration(60+rand.Intn(30)))
}

func DummySetUpEnvironment(cluster *Cluster) {
	log.Printf("Setting up environment for cluster with ID %d...", cluster.Id)
    // Simulate time taken to set up environment
    time.Sleep(time.Second * time.Duration(rand.Intn(30)))
}

func DummyConfigureCluster(cluster *Cluster) {
	log.Printf("Configuring cluster with ID %d...", cluster.Id)
    // Simulate time taken to configure cluster
    time.Sleep(time.Second * time.Duration(rand.Intn(30)))
}

func DummyExecuteSimulation(simulation *Simulation) {
	log.Printf("Executing simulation with ID %d...", simulation.Id)
    // Simulate time taken to execute simulation
    time.Sleep(time.Second * time.Duration(10+rand.Intn(20)))
}

func DummyStartDSP(simulation *Simulation) {
	log.Printf("Starting DSP for simulation with ID %d...", simulation.Id)
    // Simulate time taken to start DSP
    time.Sleep(time.Second * time.Duration(5+rand.Intn(10)))
}

func getSimulationFromDB(simulationId int) *Simulation {
	log.Printf("Retrieving simulation with ID %d from DB...", simulationId)
    // Placeholder: Return a dummy simulation for the provided ID
    return &Simulation{
        Id: simulationId,
        Cluster: &Cluster{
            Id: 1, // Dummy cluster ID
            // ... other fields initialized as needed
        },
        // ... other fields initialized as needed
    }
}

func getClusterFromDB(clusterId int) *Cluster {
	log.Printf("Retrieving cluster with ID %d from DB...", clusterId)
    // Placeholder: Return a dummy cluster for the provided ID
    return &Cluster{
        Id: clusterId,
        // ... other fields initialized as needed
    }
}