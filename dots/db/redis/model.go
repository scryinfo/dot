package redis

import "github.com/albrow/zoom"

type Model interface {
    zoom.Model // zoom.RandomID implements this interface
    zoom.MarshalerUnmarshaler
    ModelName() string
    SetModelName(string)
}
