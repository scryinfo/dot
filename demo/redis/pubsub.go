package main

import (
    "fmt"
    "time"
)

func pubsubDemo() {
    pubsub := dbClient.Subscribe(dbClient.Context(), "mychannel1")

    // Wait for confirmation that subscription is created before publishing anything.
    _, err := pubsub.Receive(dbClient.Context())
    if err != nil {
        panic(err)
    }

    // Go channel which receives messages.
    ch := pubsub.Channel()

    // Publish a message.
    err = dbClient.Publish(dbClient.Context(), "mychannel1", "hello").Err()
    if err != nil {
        panic(err)
    }

    time.AfterFunc(time.Second, func() {
        // When pubsub is closed channel is closed too.
        //_ = pubsub.Close()
    })

    // Consume messages. 阻塞
    for msg := range ch {
        fmt.Println(msg.Channel, msg.Payload)
    }
}
