package main

import "context"

func getValue(key string) (string, error) {
	return dbClient.Get(context.Background(), key).Result()
}
