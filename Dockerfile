FROM alpine:3.4

# install certs
RUN apk --update add ca-certificates

# copy the certs of jwt_proxy
COPY certs /certs

# copy static web content of jwt_proxy
COPY www /www

# copy the config to be used to /config.yml
COPY config.yml /config.yml

# start /cmd/main when running this image
ADD jwt_proxy /jwt_proxy

# Start with defined config
CMD ["/jwt_proxy", "--config=/config.yml"]

EXPOSE 8080
