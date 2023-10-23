package main

import (
    "log"
    "strconv" // Add this import
    amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection
var ch *amqp.Channel

func InitializeRabbitMQ() {	
    log.Println("Initializing RabbitMQ connection...")
    var err error
    conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }

    ch, err = conn.Channel()
    if err != nil {
        log.Fatalf("Failed to open a channel: %v", err)
    }

    // Initialize queues
    _, err = ch.QueueDeclare("ClusterQueue", true, false, false, false, nil)
    if err != nil {
        log.Fatalf("Failed to declare a queue: %v", err)
    }	
    log.Println("Successfully initialized RabbitMQ connection.")
}

func addToClusterQueue(cluster *Cluster) {
	log.Printf("Adding cluster with ID %d to ClusterQueue...", cluster.Id)
    body := strconv.Itoa(cluster.Id) // Convert int to string
    err := ch.Publish("", "ClusterQueue", false, false, amqp.Publishing{
        ContentType: "text/plain",
        Body:        []byte(body), // Convert string to []byte
    })
    if err != nil {
        log.Printf("Failed to publish a message: %v", err)
    }
	log.Println("Cluster added to ClusterQueue successfully.")
}

func addToSimulationQueue(simulation *Simulation) {
	log.Printf("Adding simulation with ID %d to SimulationQueue...", simulation.Id)
    queueName := "SimulationQueue_" + strconv.Itoa(simulation.Cluster.Id) // Convert int to string
    body := strconv.Itoa(simulation.Id)                                    // Convert int to string
    err := ch.Publish("", queueName, false, false, amqp.Publishing{        // Initialize err here
        ContentType: "text/plain",
        Body:        []byte(body), // Convert string to []byte
    })
    if err != nil {
        log.Printf("Failed to publish a message: %v", err)
    }
	log.Println("Simulation added to SimulationQueue successfully.")
}

func popFromClusterQueue() *Cluster {
    msg, ok, err := ch.Get("ClusterQueue", false)
    if err != nil {
        log.Printf("Failed to fetch a message: %v", err)
        return nil
    }
    if !ok {
        return nil // No message to process
    }
    msg.Ack(false) // Acknowledge the message
    clusterId, err := strconv.Atoi(string(msg.Body)) // Convert []byte to string and then to int
    if err != nil {
        log.Printf("Failed to convert clusterId: %v", err)
        return nil
    }
	log.Printf("Cluster with ID %d found in ClusterQueue.", clusterId)
    return getClusterFromDB(clusterId) // Assume the function expects int

}

func popFromSimulationQueue(clusterId int) *Simulation {
    queueName := "SimulationQueue_" + strconv.Itoa(clusterId) // Convert int to string correctly
    msg, ok, err := ch.Get(queueName, false)
    if err != nil {
        log.Printf("Failed to fetch a message: %v", err)
        return nil
    }
    if !ok {
        return nil // No message to process
    }
    msg.Ack(false) // Acknowledge the message
    simulationId, err := strconv.Atoi(string(msg.Body)) // Convert []byte to string and then to int
    if err != nil {
        log.Printf("Failed to convert simulationId: %v", err)
        return nil
    }
	log.Printf("Simulation with ID %d found in SimulationQueue.", simulationId)
    return getSimulationFromDB(simulationId) // Assume the function expects int

}
