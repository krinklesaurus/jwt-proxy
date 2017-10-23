# stage 1
FROM golang:1.8.3

WORKDIR /go/src/github.com/krinklesaurus/jwt_proxy

COPY . .

RUN go get -u github.com/golang/dep/cmd/dep \
    && dep ensure

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o jwt_proxy ./cmd


# stage 2
FROM alpine:3.6 

RUN apk --no-cache add ca-certificates

WORKDIR /

# copy binary from previous stage
COPY --from=0 /go/src/github.com/krinklesaurus/jwt_proxy/jwt_proxy .

# copy static web content of jwt_proxy
COPY www /www

# copy the config to be used to /config.yml
COPY ./config.yml /config.yml
COPY config.yml /config.ym

ENTRYPOINT ["/jwt_proxy"]

EXPOSE 8080
