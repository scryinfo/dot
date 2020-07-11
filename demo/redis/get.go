package main

func getValue(key string) (string, error) {
    return dbClient.Get(dbClient.Context(), key).Result()
}
