package oidc_storage

import (
	"time"

	"github.com/zitadel/oidc/v4/pkg/oidc"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
)

type EnumTo[TProto any, TGo any] interface {
	ToGoType(pEnum TProto) TGo
	ToProto(gEnum TGo) TProto
	EnumsGo() []TGo
	EnumsProto() []TProto
}

const (
	QueryAuthRequestID = "authRequestID"
	LoginEndpoint      = "/login"
)

var (
	_scopes = ScopeEx.EnumsGo()
)

func PromptToInternal(oidcPrompt oidc.SpaceDelimitedArray) []string {
	prompts := make([]string, 0, len(oidcPrompt))
	for _, oidcPrompt := range oidcPrompt {
		switch oidcPrompt {
		case oidc.PromptNone,
			oidc.PromptLogin,
			oidc.PromptConsent,
			oidc.PromptSelectAccount:
			prompts = append(prompts, oidcPrompt)
		}
	}
	return prompts
}

func MaxAgeToInternal(maxAge *uint) *durationpb.Duration {
	if maxAge == nil {
		return nil
	}
	dur := time.Duration(*maxAge) * time.Second
	return durationpb.New(dur)
}
