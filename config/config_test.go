package config

import (
	"fmt"
	"os"
)

func ExampleInitialize() {
	configPath := "../resources/test/config.yml"

	cfg, err := Initialize(configPath)
	if err != nil {
		fmt.Printf("error initializing config %v", err)
		return
	}

	fmt.Println(cfg.RootURI)
	fmt.Println(cfg.RedirectURI)
	fmt.Println(cfg.SigningMethod.Alg())
	fmt.Println(cfg.Audience)
	fmt.Println(cfg.Issuer)
	fmt.Println(cfg.Subject)
	fmt.Println(cfg.Providers["google"].String())
	fmt.Println(cfg.Providers["facebook"].String())
	fmt.Println(cfg.Providers["github"].String())
	// Output:
	// http://localhost:8080
	// http://localhost:8080/callback
	// RS256
	// your-audience
	// you
	// your-subject
	// {"client_id":"your-google-client-id","auth_url":"https://accounts.google.com/o/oauth2/auth","token_url":"https://accounts.google.com/o/oauth2/token","redirect_url":"http://localhost:8080/callback/google","scopes":["profile"]}
	// {"client_id":"your-facebook-client-id","auth_url":"https://www.facebook.com/dialog/oauth","token_url":"https://graph.facebook.com/oauth/access_token","redirect_url":"http://localhost:8080/callback/facebook","scopes":["public_profile"]}
	// {"client_id":"your-github-client-id","auth_url":"https://github.com/login/oauth/authorize","token_url":"https://github.com/login/oauth/access_token","redirect_url":"http://localhost:8080/callback/github","scopes":["user"]}
}

func ExampleInitialize_second() {
	configPath := "../resources/test/config.yml"

	os.Setenv("ROOT_URI", "http://envvar:8080")
	os.Setenv("REDIRECT_URI", "http://envvar:8080/callback")
	os.Setenv("SIGNINGMETHOD", "RS256")
	os.Setenv("JWT_AUDIENCE", "envar-audience")
	os.Setenv("JWT_ISSUER", "envar-issuer")
	os.Setenv("JWT_SUBJECT", "envar-subject")

	os.Setenv("GOOGLE_CLIENTID", "envvar-google-client-id")
	os.Setenv("GOOGLE_SECRET", "envvar-google-client-secret")
	os.Setenv("GOOGLE_SCOPES", "envvar-scope-1 envvar-scope-2")

	os.Setenv("FACEBOOK_CLIENTID", "envvar-facebook-client-id")
	os.Setenv("FACEBOOK_SECRET", "envvar-facebook-client-secret")
	os.Setenv("FACEBOOK_SCOPES", "envvar-scope-1 envvar-scope-2")

	os.Setenv("GITHUB_CLIENTID", "envvar-github-client-id")
	os.Setenv("GITHUB_SECRET", "envvar-github-client-secret")
	os.Setenv("GITHUB_SCOPES", "envvar-scope-1 envvar-scope-2")

	defer func() {
		os.Unsetenv("ROOT_URI")
		os.Unsetenv("REDIRECT_URI")
		os.Unsetenv("SIGNINGMETHOD")
		os.Unsetenv("JWT_AUDIENCE")
		os.Unsetenv("JWT_ISSUER")
		os.Unsetenv("JWT_SUBJECT")
		os.Unsetenv("PUBLICKEY_PATH")
		os.Unsetenv("PRIVATEKEY_PATH")

		os.Unsetenv("GOOGLE_CLIENTID")
		os.Unsetenv("GOOGLE_SECRET")
		os.Unsetenv("GOOGLE_SCOPES")

		os.Unsetenv("FACEBOOK_CLIENTID")
		os.Unsetenv("FACEBOOK_SECRET")
		os.Unsetenv("FACEBOOK_SCOPES")

		os.Unsetenv("GITHUB_CLIENTID")
		os.Unsetenv("GITHUB_SECRET")
		os.Unsetenv("GITHUB_SCOPES")
	}()

	cfg, err := Initialize(configPath)
	if err != nil {
		fmt.Printf("error initializing config %v", err)
		return
	}

	fmt.Println(cfg.RootURI)
	fmt.Println(cfg.RedirectURI)
	fmt.Println(cfg.SigningMethod.Alg())
	fmt.Println(cfg.Audience)
	fmt.Println(cfg.Issuer)
	fmt.Println(cfg.Subject)
	fmt.Println(cfg.Providers["google"].String())
	fmt.Println(cfg.Providers["facebook"].String())
	fmt.Println(cfg.Providers["github"].String())
	// Output:
	// http://envvar:8080
	// http://envvar:8080/callback
	// RS256
	// envar-audience
	// envar-issuer
	// envar-subject
	// {"client_id":"envvar-google-client-id","auth_url":"https://accounts.google.com/o/oauth2/auth","token_url":"https://accounts.google.com/o/oauth2/token","redirect_url":"http://envvar:8080/callback/google","scopes":["envvar-scope-1","envvar-scope-2"]}
	// {"client_id":"envvar-facebook-client-id","auth_url":"https://www.facebook.com/dialog/oauth","token_url":"https://graph.facebook.com/oauth/access_token","redirect_url":"http://envvar:8080/callback/facebook","scopes":["envvar-scope-1","envvar-scope-2"]}
	// {"client_id":"envvar-github-client-id","auth_url":"https://github.com/login/oauth/authorize","token_url":"https://github.com/login/oauth/access_token","redirect_url":"http://envvar:8080/callback/github","scopes":["envvar-scope-1","envvar-scope-2"]}
}
