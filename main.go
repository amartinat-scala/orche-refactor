package main

func main() {
    InitializeRabbitMQ() 
    go ClusterWorker()
    StartAPI() 
}
