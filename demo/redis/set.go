package main

import "time"

func setValue(key, value string, expiration ...time.Duration) error {
    var tim time.Duration
    if len(expiration) != 0 {
        tim = expiration[0]
    }

    return dbClient.Set(dbClient.Context(), key, value, tim).Err()
}

func setExpire(key string, expiration time.Duration) error {
    return dbClient.Expire(dbClient.Context(), key, expiration).Err()
}
