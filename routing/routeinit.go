package routing

import (
	"github.com/fnuerpod/odb-obsauth/authdatabase"
	"github.com/fnuerpod/odb-obsauth/config"
	"github.com/fnuerpod/odb-obsauth/logging"

	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Route struct {
	Type string
	Exec func(http.ResponseWriter, *http.Request)
}

type RoutingMemory struct {
	database *authdatabase.MCAuthDB_sqlite3

	logger             *logging.Logger
	site_configuration *config.Configuration
}

func Routing_Init(logger *logging.Logger) (routes map[string]Route, muxxer *mux.Router) {
	// basic config initialisation
	config.InitialiseConfig()

	// initialise everything needed for routing here, shaves down lines in main func.
	database := authdatabase.InitSQLite3DB(logger)

	muxxer = mux.NewRouter()
	site_configuration := config.InitialiseConfig()

	procmem := &RoutingMemory{
		database: database,

		logger:             logger,
		site_configuration: site_configuration,
	}

	// Runtime generate hashmap of routes based on function names
	// _0_ defines a new directory
	// _1_ defines a period
	// Prefix with GET to define a GET route, POST for a POST route
	// GET_0_ = GET '/'
	// GET_0_data_1_jpg = GET '/data.jpg'
	// POST_0_user = POST '/user'

	routes = map[string]Route{}

	typ := reflect.TypeOf(procmem)
	ref := reflect.ValueOf(procmem)

	methods := ref.NumMethod()
	syntaxreplacer := strings.NewReplacer("_0_", "/", "_1_", ".")

	// Valid routes are defined here as a hashmap
	// Add more as needed
	validroutes := map[string]struct{}{
		"GET":  {},
		"POST": {},
	}

	for i := 0; i < methods; i++ {
		// Get method name
		f := typ.Method(i).Name

		potroute := strings.Split(f, "_0_")

		// might panic if len(f) == 0, but a function name should never be len 0 so whatever
		if _, ok := validroutes[potroute[0]]; ok {
			routeTyp := potroute[0]

			// Replace syntax of _n_ values
			routedir := syntaxreplacer.Replace(strings.TrimPrefix(f, routeTyp))

			// Test if function is valid args
			fp, ok := ref.MethodByName(f).Interface().(func(http.ResponseWriter, *http.Request))

			if ok {
				routes[routedir] = Route{routeTyp, func(a http.ResponseWriter, b *http.Request) {
					//logger.Log.Println("GET " + b.URL.String())
					fp(a, b)
				}}

			} else {
				panic("RoutingMemory." + f + " Not valid function")
			}
		}
	}

	for key, element := range routes {
		muxxer.HandleFunc(key, element.Exec)
	}

	return
}
