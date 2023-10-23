package main

import (
    // "math/rand"
    "sync"
    // "time"
	"log"
)


func ClusterWorker() {
	log.Println("Starting ClusterWorker...")
    for {
        cluster := popFromClusterQueue()
        if cluster == nil {
			log.Println("Cluster is nil. Skipping iteration and checking queue again...")
            continue // If the cluster is nil, skip the iteration and check the queue again
        }
		log.Printf("Retrieved cluster with ID %d from ClusterQueue.", cluster.Id)
       

        cluster.Mutex.Lock()

        // Use WaitGroup for parallelizable tasks
        var wg sync.WaitGroup

        // Example parallel tasks: BootCluster, SetUpEnvironment
        wg.Add(2)
        go func() {
            DummyBootCluster(cluster)
            wg.Done()
        }()
        go func() {
            DummySetUpEnvironment(cluster)
            wg.Done()
        }()

        wg.Wait()

        // Example serial tasks after parallel tasks are done
        DummyConfigureCluster(cluster)

        cluster.SetReady()
		log.Printf("Cluster with ID %d is ready.", cluster.Id)
        cluster.Mutex.Unlock()
    }
}

func SimulationWorker(clusterId int) {
    for {
        simulation := popFromSimulationQueue(clusterId)
        if simulation == nil {
            continue // If the simulation is nil, skip the iteration and check the queue again
        }

        simulation.Cluster.Mutex.Lock()

        DummyExecuteSimulation(simulation)

        // After simulation is complete, kick off DSP
        DummyStartDSP(simulation)

        simulation.Cluster.Mutex.Unlock()
    }
}
