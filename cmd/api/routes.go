package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healtcheck", app.healtcheck)

	router.HandlerFunc(http.MethodPost, "/v1/invoices", app.createInvoiceHandler)

	return router
}
