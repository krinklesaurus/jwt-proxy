package main

import (
	"flag"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/krinklesaurus/jwt_proxy/config"
	"github.com/krinklesaurus/jwt_proxy/core"
	"github.com/krinklesaurus/jwt_proxy/handler"
	"github.com/krinklesaurus/jwt_proxy/log"
	"github.com/krinklesaurus/jwt_proxy/user"
	"github.com/meatballhat/negroni-logrus"
	"github.com/urfave/negroni"
)

func main() {
	configPtr := flag.String("config", "", "configuration file")

	flag.Parse()

	config, err := config.Initialize(*configPtr)
	if err != nil {
		panic(err)
	}

	log.Infof("Config initialized: %s", config.String())

	userService := &user.HashUserService{}

	tokenizer := core.NewRSATokenizer(config.SigningMethod, config.PrivateRSAKey)

	core := core.New(config, tokenizer, userService)
	store, err := handler.NewHttpSessionStore()
	if err != nil {
		panic(err)
	}
	handler, err := handler.New(core, store)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handler.HomeHandler)
	r.HandleFunc("/login", handler.LoginHandler)
	r.HandleFunc("/login/{provider}", handler.ProviderLoginHandler)
	r.HandleFunc("/callback/{provider}", handler.CallbackHandler)
	r.HandleFunc("/pubkey", handler.PublicKeyHandler)

	// r.HandleFunc("/auth", handler.AuthHandler).Methods("POST")

	n := negroni.New()
	n.Use(negronilogrus.NewMiddleware())
	n.UseHandler(r)

	http.ListenAndServe(":8080", n)
}
