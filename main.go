package main

func main() {
    InitializeRabbitMQ() // Initialize queues
    go ClusterWorker()
    StartAPI() 
}
