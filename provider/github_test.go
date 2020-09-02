package provider

import "fmt"

func ExampleGithubAuthCodeURL() {
	f := NewGithub("http://localhost:8080", "client-id", "client-secret", []string{"scope-1", "scope-2"})
	authCodeURL := f.AuthCodeURL("state")
	fmt.Println(authCodeURL)

	// Output:
	// https://github.com/login/oauth/authorize?client_id=client-id&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fjwt-proxy%2Fcallback%2Fgithub&response_type=code&scope=scope-1+scope-2&state=state
}
