package api

import (
	"net/http"
	"strconv"

	"github.com/jmcfarlane/notable/flags"
	"github.com/julienschmidt/httprouter"

	log "github.com/Sirupsen/logrus"
)

// Start the service
func Start(router *httprouter.Router) {
	log.Infof("Listening on %s:%v", *flags.Bind, *flags.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*flags.Port), router))
}
