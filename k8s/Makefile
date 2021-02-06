VALUES_FILE=.jwt-proxy.yaml
SECRETS_FILE=.secrets.yaml

.PHONY: clean
clean:
	@rm -rf ${VALUES_FILE} &&\
		@rm -rf ${SECRETS_FILE}

.PHONY: test
test:
	env $$(cat .credentials | xargs) envsubst < "secrets.yaml" > ${SECRETS_FILE} &&\
		env $$(cat .credentials | xargs) envsubst < "jwt-proxy.yaml" > ${VALUES_FILE}


.PHONY: inject-secrets
inject-secrets:
	env $$(cat .credentials | xargs) envsubst < "secrets.yaml" > ${SECRETS_FILE} &&\
		kubectl apply -f ${SECRETS_FILE} &&\
	env $$(cat .credentials | xargs) envsubst < "jwt-proxy.yaml" > ${VALUES_FILE} &&\
		kubectl apply -f ${VALUES_FILE}

.PHONY: setup
setup: inject-secrets
	kubectl apply -f namespace.yaml &&\
		kubectl apply -f ${SECRETS_FILE} &&\
		kubectl apply -f ingress.yaml &&\
		kubectl apply -f helloservice.yaml &&\
		kubectl apply -f ${VALUES_FILE}


.PHONY: setup-jwt-proxy
setup-jwt-proxy:
	env $$(cat .credentials | xargs) envsubst < "secrets.yaml" > ${SECRETS_FILE} &&\
		env $$(cat .credentials | xargs) envsubst < "jwt-proxy.yaml" > ${VALUES_FILE}