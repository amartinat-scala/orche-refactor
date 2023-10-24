package main

func main() {
    InitializeRabbitMQConnection() 
    go ClusterWorker()
    StartAPI() 
}
