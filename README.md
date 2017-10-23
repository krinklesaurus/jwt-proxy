# jwt_proxy

<p align="center">
    <img src="https://img.shields.io/docker/stars/krinklesaurus/jwt_proxy.svg"
          alt="Docker Hub Stars">
    <img src="https://img.shields.io/docker/pulls/krinklesaurus/jwt_proxy.svg"
          alt="Docker Hub Pulls version">
    <img src="https://img.shields.io/docker/automated/krinklesaurus/jwt_proxy.svg"
          alt="Docker Hub Automated">
    <img src="https://img.shields.io/docker/build/krinklesaurus/jwt_proxy.svg"
          alt="Docker Hub Build status">
</p>

---

**PLEASE NOTE THAT jwt_proxy SHOULD NOT BE USED IN PRODUCTION ENVIRONMENTS YET!**

---

## What is jwt_proxy?

jwt_proxy is a small OAuth2 proxy service that returns JWT tokens that can be used for offline authentication in distributed environments like microservices, where several different services need a fast and safe way to process authentication.

The JWT token returned by jwt_proxy cannot only be used to check wether a user is allowed to call your api, it also contains the OAuth2 provider's access token and a unique user id based on the provider and the user's id from the provider.

## How does jwt_proxy work?

jwt_proxy provides a simple login page that contains links to all supported third party OAuth2 providers (currently these are Google, Facebook and Github). As soon as a user selects one of the provider, jwt_proxy starts the OAuth2 `Authorization Code Grant` to obtain an access token from the selected OAuth2 provider. This token is then signed with a custom private signing key and returned as a JWT token which can be used by further calls to your Web API.

This figure illustrates the basic flow jwt_proxy provides:

```
+----------------+                 +-----------------+
|                |                 |                 |
|      User      |-------(9)------>|     Your API    |
|                |                 |                 |
|                |                 |                 |
+----------------+                 +-----------------+
  |       ^                                 |
 (1)      |     _ _ _ _ _ _ (10) _ _ _ _ _ _|
  |      (8)   |         
  v       |    v
+----------------+                 +-----------------+
|                |-------(2)------>|                 |
|    jwt_proxy   |<------(3)-------| OAuth2 Provider |
|                |                 |  (e.g. Google)  |
|                |-------(4)------>|                 |
|                |<------(5)-------|                 |
|                |                 |                 |
|                |-------(6)------>|                 |
|                |<------(7)-------|                 |
+----------------+                 +-----------------+                        

```
1. The user opens `http://myjwt_proxy/login` and selects one of the OAuth2 provider login buttons, e.g. `Sign in with Google`
2. jwt_proxy redirects the user to the OAuth2 provider's authorization endpoint, e.g. `https://accounts.google.com/o/oauth2/auth` with all required parameters to perform the OAuth2 [Authorization Code Grant](https://tools.ietf.org/html/rfc6749#section-4.1)
3. The Oauth2 Provider redirects back to jwt_proxy to `/callback/{provider}`, e.g. `/callback/google`.
4. jwt_proxy requests an access token from the OAuth2 provider by calling its token endpoint with the provided authorization code from 3.
5. An OAuth2 Token is returned from the OAuth2 provider to jwt_proxy.
6. jwt_proxy puts the OAuth2 provider`s access token to obtain the user's id from the OAuth2 provider
7. The user's id is returned back to jwt_proxy
8. jwt_proxy marshalls the access token, the selected provider and a hashed user id into a JWT token and signs it with a custom private signing key. This JWT token is returned to your client, e.g. to his mobile app.
9. The user makes request to **your** API with the JWT token.
10. Your API calls `http://myjwt_proxy/pubkey` to obtain the public key from jwt_proxy and checks wether the JWT token is valid. You now
  1. know your user is allowed to call your API
  2. have a unique user id you can work with
  3. could make additional calls to the OAuth2 provider with the provider's access token in the JWT token

 ## How can I test/lint/build jwt_proxy?
 
 jwt_proxy comes with a Makefile that defines target for clean, test, lint, build and dockerbuild.

 ### Running jwt_proxy

 Run the `cmd/main.go` with Go tools, e.g.

 ```
 go run cmd/main.go --config=config.yml
 ```

 In order to run jwt_proxy as a Docker container, you can run the provided `build.sh` script or create the Docker image yourself. Note that the `Dockerfile` copies the `www`, `certs` folder and the `config.yml` into the image, so in case you want to use different names you need to adapt the Dockerfile correspondingly.

 ### Helping tools

 jwt_proxy comes with a little help in the `certs` directory. It contains an additional `install.go` with which you can simply create some certs for jwt_proxy and create a test token with a currently expiration of one year. The created certs within this directory are also the ones jwt_proxy uses as a default.

 #### Create certs

 ```
 go run certs/install.go --config=config.yml --certs=true
 ```

 #### Create a test token

 ```
 go run certs/install.go --config=config.yml --token=true
 ```

 ### Configuring jwt_proxy

jwt_proxy requires some configuration in order to be run. A basic configuration file can be found in `config.yml`. The config file contains the following keys:

<table>
  <tr>
    <th>config.yml Key</th>
    <th>environment var</th>
    <th>Meaning</th>
  </tr>
  <tr>
    <td>root_uri</td>
    <td>ROOT_URI</td>  
    <td>This is the URI that is used as the base URI for the `/callback/{provider}` endpoint that needs to be registered by an OAuth2 provider, e.g. `[root_uri]/callback/google` needs to be registred by Google as the OAuth2 callback URI</td>
  <tr>
  <tr>
    <td>redirect_uri</td>  
    <td>REDIRECT_URI</td>
    <td>The redirect URI to the client that jwt_proxy redirects to with the JWT token as a query parameter. For instance, with Ionic you need to set this to `http://localhost/callback`. jwt_proxy that redirects to `http://localhost/callback?token=[JWT-Token]`</td>
  <tr>
  <tr>
    <td>jwt.signingMethod</td>  
    <td>SIGNINGMETHOD</td>
    <td>The used method for signing the JWT token. Supported methods can be found in `app.go`. **The selected method must match the used private and public key!**</td>
  <tr>
  <tr>
    <td>jwt.public_key</td>  
    <td>PUBLICKEY_PATH</td>
    <td>Path to the public key that is exposed by jwt_proxy under `http://myjwt_proxy/pubkey` and can be used by your authentication middleware to check a user's JWT token.</td>
  <tr>
  <tr>
    <td>jwt.private_key</td>  
    <td>PRIVATEKEY_PATH</td>
    <td>Path to the private key that is used by jwt_proxy to sign the JWT token. **Make sure this private key stays absolutely secret**</td>
  <tr>  
  <tr>
    <td>providers.[name].client_id</td>  
    <td>
      GOOGLE_CLIENTID
      FACEBOOK_CLIENTID
      GITHUB_CLIENTID
    </td>
    <td>The OAuth2 client id for an OAuth2 provider.</td>
  <tr>     
  <tr>
    <td>providers.[name].client_secret</td>  
    <td>
      GOOGLE_SECRET
      FACEBOOK_SECRET
      GITHUB_SECRET
    </td>    
    <td>The OAuth2 client secret for an OAuth2 provider.</td>
  <tr>
  <tr>
    <td>providers.[name].scopes</td>
    <td>
      GOOGLE_SCOPES
      FACEBOOK_SCOPES
      GITHUB_SCOPES
    </td>        
    <td>The OAuth2 scopes for an OAuth2 provider. The selected scopes must at least contain the necessary scope to fetch the user's unique id from the provider. If you want to make additional API calls to the OAuth2 provider, add your custom scopes here.</td>
  <tr>        
</table>

 jwt_proxy can be run either as a standard application by calling `go run cmd/main.go` or as a docker container `docker run jwt_proxy:[tag]`(recommended way).

 When run as a docker container you can easily use the environemnt variables to configure jwt_proxy, e.g.
 ```
docker run
-e ROOT_URI=http://myjwt_proxy:80
-e REDIRECT_URI=http://localhost/callback
-e SIGNINGMETHOD=RS256
-e PUBLICKEY_PATH=certs/publickey.pem
-e PRIVATEKEY_PATH=certs/privatekey.pem
-e GOOGLE_CLIENTID=my-google-client-id
-e GOOGLE_SECRET=my-google-secret
-e GOOGLE_SCOPES=profile
-p 80:8080
jwt_proxy[:tag]
 ```

 ## FAQs

 #### How is the unique user id generated?

 The unique user id is generated by hashing the Oauth2 provider's name (e.g. google) along with the unique user id (e.g. your-id) from this provider, e.g.

  `google:your-id` becomes `878fadbf4add33[...]6643f442339de`

  As every user id is unique per provider, the hash is unique, too. As jwt_proxy is open source, feel free to fork and use your own way to generate a unique user id (e.g. generate a UUID and store in a database along with the provider and the user's id from the provider).
