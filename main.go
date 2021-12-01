package main

// imports
import (
	"net/http"

	"github.com/fnuerpod/odb-obsauth/config"
	"github.com/fnuerpod/odb-obsauth/logging"
	"github.com/fnuerpod/odb-obsauth/routing"
)

func main() {
	// instantise functions for site.
	// initialises database, oauth2 config, otp handler, session storage and the fiber app.
	// routing initialisation will also handle the addition of middlewares to the base app.

	site_configuration := config.InitialiseConfig()

	logger := logging.New()

	_, routing_mux := routing.Routing_Init(logger)

	// start listening
	logger.Log.Println("ODB Stream Authenticator starting on port " + site_configuration.BindPort + "...")
	//app.Listen(":" + constants.BindPort)

	err := http.ListenAndServe(":"+site_configuration.BindPort, routing_mux)
	logger.Fatal.Fatalln(err)

}

//EOF
