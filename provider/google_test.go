package provider

import "fmt"

func ExampleGoogleAuthCodeURL() {
	f := NewGoogle("http://localhost:8080", "client-id", "client-secret", []string{"scope-1", "scope-2"})
	authCodeURL := f.AuthCodeURL("state")
	fmt.Println(authCodeURL)

	// Output:
	// https://accounts.google.com/o/oauth2/auth?client_id=client-id&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fjwt-proxy%2Fcallback%2Fgoogle&response_type=code&scope=scope-1+scope-2&state=state
}
