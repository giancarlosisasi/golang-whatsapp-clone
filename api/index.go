package handler

import (
	"errors"
	"golang-whatsapp-clone/server"
	"net/http"
)

// Handler is the main entry point of the application. Think of it like the main() method
func Handler(w http.ResponseWriter, r *http.Request) {
	// This is needed to set the proper request path in `*fiber.Ctx`
	r.RequestURI = r.URL.String()

	app, server := server.NewServer()
	defer app.DBpool.Close()

	err := server.ListenAndServe()
	if err != nil {
		msg := "error to run the application in vercel api/index.go file"
		app.Logger.Error().Msgf("%s: %s", msg, err)

		app.Handler.ServerErrorResponse(w, r, errors.New(msg))
	}
}
