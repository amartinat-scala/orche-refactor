package main

import (
    "log"
    "strconv" 
    amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection

func InitializeRabbitMQConnection() {
    var err error
    conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }
    log.Println("Successfully initialized RabbitMQ connection.")
}

func getChannel() (*amqp.Channel, error) {
    return conn.Channel()
}

func addToClusterQueue(cluster *Cluster) {
    ch, err := getChannel()
    if err != nil {
        log.Fatalf("Failed to open a channel: %v", err)
        return 
    }
    defer ch.Close()
	log.Printf("Adding cluster with ID %d to ClusterQueue...", cluster.Id)
    body := strconv.Itoa(cluster.Id) 
	_, err = ch.QueueDeclare(
		"ClusterQueue",  // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
    if err != nil {
        log.Printf("Failed to publish a message: %v", err)
    }
    err = ch.Publish("", "ClusterQueue", false, false, amqp.Publishing{
        ContentType: "text/plain",
        Body:        []byte(body), 
		DeliveryMode: amqp.Persistent,
    })
	log.Println("Cluster added to ClusterQueue successfully.")
}

func addToSimulationQueue(simulation *Simulation) {
    ch, err := getChannel()
    if err != nil {
        log.Fatalf("Failed to open a channel: %v", err)
        return 
    }
    defer ch.Close()

    log.Printf("Adding simulation with ID %d to SimulationQueue...", simulation.Id)

    queueName := "SimulationQueue_" + strconv.Itoa(simulation.Cluster.Id)

    // Declare the queue (ensure it exists before publishing)
    _, err = ch.QueueDeclare(
        queueName,  // name
        true,       // durable
        false,      // delete when unused
        false,      // exclusive
        false,      // no-wait
        nil,        // arguments
    )
    if err != nil {
        log.Fatalf("Failed to declare the queue: %v", err)
        return
    }

    body := strconv.Itoa(simulation.Id)                                  
    err = ch.Publish("", queueName, false, false, amqp.Publishing{    
        ContentType: "text/plain",
        Body:        []byte(body), 
		DeliveryMode: amqp.Persistent,
    })
    if err != nil {
        log.Printf("Failed to publish a message: %v", err)
    }

    log.Println("Simulation added to SimulationQueue successfully.")
}


func popFromClusterQueue() *Cluster {
    ch, err := getChannel()
    if err != nil {
        log.Fatalf("Failed to open a channel: %v", err)
        return nil
    }
    defer ch.Close()
    msg, ok, err := ch.Get("ClusterQueue", false)
    if err != nil {
        log.Printf("Failed to fetch a message: %v", err)
        return nil
    }
    if !ok {
        return nil
    }
    msg.Ack(false) 
    clusterId, err := strconv.Atoi(string(msg.Body)) 
    if err != nil {
        log.Printf("Failed to convert clusterId: %v", err)
        return nil
    }
	log.Printf("Cluster with ID %d found in ClusterQueue.", clusterId)
    return getClusterFromDB(clusterId) 

}

func popFromSimulationQueue(clusterId int) *Simulation {
    ch, err := getChannel()
    if err != nil {
        log.Fatalf("Failed to open a channel: %v", err)
        return nil
    }
    defer ch.Close()
    queueName := "SimulationQueue_" + strconv.Itoa(clusterId)
    msg, ok, err := ch.Get(queueName, false)
    if err != nil {
        log.Printf("Failed to fetch a message: %v", err)
        return nil
    }
    if !ok {
        return nil
    }
    msg.Ack(false) 
    simulationId, err := strconv.Atoi(string(msg.Body)) 
    if err != nil {
        log.Printf("Failed to convert simulationId: %v", err)
        return nil
    }
	log.Printf("Simulation with ID %d found in SimulationQueue.", simulationId)
    return getSimulationFromDB(simulationId) 

}
