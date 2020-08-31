package main

import (
	"flag"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/krinklesaurus/jwt-proxy/config"
	"github.com/krinklesaurus/jwt-proxy/core"
	"github.com/krinklesaurus/jwt-proxy/handler"
	"github.com/krinklesaurus/jwt-proxy/log"
	"github.com/krinklesaurus/jwt-proxy/user"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

func main() {
	log.WithLevel(logrus.InfoLevel)

	configPtr := flag.String("config", "config.yml", "configuration file")

	flag.Parse()

	config, err := config.Initialize(*configPtr)
	if err != nil {
		log.Errorf("error initializing config from %s, %v", *configPtr, err)
		return
	}

	log.Infof("Config initialized: %s", config.String())

	userService := &user.PlainUserService{}

	tokenizer := core.NewRSATokenizer(config.SigningMethod, config.PrivateRSAKey)

	core := core.New(config, tokenizer, userService)
	store, err := handler.NewHTTPSessionStore()
	if err != nil {
		log.Errorf("error initializing session store %v", err)
		return
	}
	handler, err := handler.New(config, core, store)
	if err != nil {
		log.Errorf("error initializing handler store %v", err)
		return
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handler.HomeHandler)
	r.HandleFunc("/login", handler.LoginHandler)
	r.HandleFunc("/login/{provider}", handler.ProviderLoginHandler)
	r.HandleFunc("/callback/{provider}", handler.CallbackHandler)
	r.HandleFunc("/pubkey", handler.PublicKeyHandler)
	r.HandleFunc("/token", handler.VerifyToken)

	n := negroni.New()
	n.Use(negronilogrus.NewMiddleware())
	n.UseHandler(r)

	http.ListenAndServe(":8080", n)
}
