package httpapi

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type GoogleAuth struct {
	OAuthConfig oauth2.Config
	Verifier    *oidc.IDTokenVerifier
	WebOrigin   string
}
