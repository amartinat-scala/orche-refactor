package main

import (
    "net/http"
    "strconv"
	"log"
)

func StartAPI() {
    log.Println("Starting the API on port 8080...")
    http.HandleFunc("/runSimulation", RunSimulationEndpoint)
    http.ListenAndServe(":8080", nil)
}


func RunSimulationEndpoint(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request for /runSimulation")

    simulationIdStr := r.URL.Query().Get("simulationId")
    simulationId, err := strconv.Atoi(simulationIdStr)
    if err != nil {
        log.Printf("Error: Invalid simulationId: %s", simulationIdStr)
        return
    }
    
    simulation := getSimulationFromDB(simulationId)
    
    // Log the retrieved simulation details (assuming the simulation has a String() method or similar)
    log.Printf("Retrieved simulation: %s", simulation)
    
    simulation.Cluster.Mutex.Lock()
    defer simulation.Cluster.Mutex.Unlock()
    
    if simulation.Cluster.IsReady() {
        log.Println("Cluster is ready. Adding simulation to SimulationQueue.")
        addToSimulationQueue(simulation)
    } else {
        log.Println("Cluster is not ready. Adding cluster to ClusterQueue.")
        addToClusterQueue(simulation.Cluster)
    }
}

