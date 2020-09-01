package provider

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

// Provider is the interface every OAuth provider has to fulfill for being
// used within jwt-proxy
type Provider interface {
	AuthCodeURL(state string) string
	User() (string, error)
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
	String() string
	Name() string
	ClientID() string
}
