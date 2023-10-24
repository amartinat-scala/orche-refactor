package main

import (
	// "math/rand"
	"strconv"
	"sync"

	// "time"
	"log"
)

func ClusterWorker() {
	log.Println("Starting ClusterWorker...")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		"ClusterQueue", // queue name
		"",             // consumer tag
		false,          // auto-ack
		false,          // exclusive
		false,          // no local
		false,          // no wait
		nil,            // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	for msg := range msgs {
		clusterId, err := strconv.Atoi(string(msg.Body))
		if err != nil {
			log.Printf("Failed to convert clusterId: %v", err)
			continue
		}
		log.Printf("Retrieved cluster with ID %d from ClusterQueue.", clusterId)

		cluster := getClusterFromDB(clusterId)
		if cluster == nil {
			continue
		}

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
		// Acknowledge the message once processing is complete.
		msg.Ack(false)
	}
}

func SimulationWorker(clusterId int) {
	log.Println("Starting SimulationWorker for cluster", clusterId)

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	queueName := "SimulationQueue_" + strconv.Itoa(clusterId)

	msgs, err := ch.Consume(
		queueName, // queue name
		"",        // consumer tag
		false,     // auto-ack
		false,     // exclusive
		false,     // no local
		false,     // no wait
		nil,       // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	for msg := range msgs {
		simulationId, err := strconv.Atoi(string(msg.Body))
		if err != nil {
			log.Printf("Failed to convert simulationId: %v", err)
			continue
		}
		log.Printf("Retrieved simulation with ID %d from SimulationQueue.", simulationId)

		simulation := getSimulationFromDB(simulationId)
		if simulation == nil {
			continue
		}

		simulation.Cluster.Mutex.Lock()

		DummyExecuteSimulation(simulation)

		DummyStartDSP(simulation)

		log.Printf("Sim %s complete!", simulation.Id)

		simulation.Cluster.Mutex.Unlock()

		// Acknowledge the message once processing is complete.
		msg.Ack(false)
	}
}
