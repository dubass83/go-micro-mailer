package main

import (
	"net/http"
	"os"

	"github.com/dubass83/go-micro-mailer/cmd/api"
	"github.com/dubass83/go-micro-mailer/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	conf, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot load configuration")
	}
	if conf.Enviroment == "devel" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		// log.Debug().Msgf("config values: %+v", conf)
	}

	s := api.CreateNewServer(conf)
	s.ConfigureCORS()
	s.AddMiddleware()
	s.MountHandlers()
	log.Info().
		Msgf("start listening on the port %s\n", conf.HTTPAddressString)
	err = http.ListenAndServe(conf.HTTPAddressString, s.Router)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("canot start http service")
	}
}
