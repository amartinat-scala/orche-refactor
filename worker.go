package main

import (
    // "math/rand"
    "sync"
    // "time"
	"log"
)


func ClusterWorker() {
	log.Println("Starting ClusterWorker...")
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()
    for {
        cluster := popFromClusterQueue()
        if cluster == nil {
			// log.Println("Cluster is nil. Skipping iteration and checking queue again...")
            continue
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
		// for each simulation associated with the cluster, add it to the sim queue
		simulation := getSimulationFromDB(cluster.Id) // filter by sims waiting to be run

		// if err != nil {
		// 	log.Printf("Error: Invalid simulationId: %s", simulationId)
		// 	return
		// }
		addToSimulationQueue(simulation)
		go SimulationWorker(cluster.Id)
        cluster.Mutex.Unlock()
    }
}

func SimulationWorker(clusterId int) {
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()
    for {
        simulation := popFromSimulationQueue(clusterId)
        if simulation == nil {
            continue
        }

        simulation.Cluster.Mutex.Lock()

        DummyExecuteSimulation(simulation)

        DummyStartDSP(simulation)

		log.Printf("Sim %s complete!", simulation.Id)
        simulation.Cluster.Mutex.Unlock()
    }
}
