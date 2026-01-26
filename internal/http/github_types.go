package httpapi

import "golang.org/x/oauth2"

type GitHubAuth struct {
	OAuthConfig oauth2.Config
	WebOrigin   string
}
