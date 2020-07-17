package main

var dbClient = conn("192.168.1.65:6379")

func main() {
    //basicDemo()
    //expireDemo()
    //updateDemo()
    pubsubDemo()
}
