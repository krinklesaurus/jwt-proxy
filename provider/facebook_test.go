package provider

import "fmt"

func ExampleFacebookAuthCodeURL() {
	f := NewFacebook("http://localhost:8080", "client-id", "client-secret", []string{"scope-1", "scope-2"})
	authCodeURL := f.AuthCodeURL("state")
	fmt.Println(authCodeURL)

	// Output:
	// https://www.facebook.com/v3.2/dialog/oauth?client_id=client-id&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fjwt-proxy%2Fcallback%2Ffacebook&response_type=code&scope=scope-1+scope-2&state=state
}
