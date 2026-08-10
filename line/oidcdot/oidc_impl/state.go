package oidc_impl

import "github.com/google/uuid"

func NewState() string {
	return uuid.New().String()
}
