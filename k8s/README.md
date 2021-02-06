# Kubernetes setup

1. Run `make create-certs` from Makefile in root directy. New certificates are created and put into folder `certs`
2. Put your secrets into a file `.credentials` within this directory.

Example:
```
PROVIDERS_GITHUB_CLIENTID=28745b23hb23z8
PROVIDERS_GITHUB_CLIENTSECRET=345ib4398f3b4k43bf834
PROVIDERS_GITHUB_SCOPES=user

PRIVATE_PEM=LS0tLS1C...
PUBLIC_PEM=LS0tLS1CR...
```

Note that you must encode your certificates with Base64 running `cat private.pem | base64` and `cat public.pem | base64` first

3. Run `make test` and check if a file `.jwt-proxy.yaml` and `.secrets.yaml` have been created. Within the files, you should find Kube configs with the secrets from the `.credentials` file injected to them