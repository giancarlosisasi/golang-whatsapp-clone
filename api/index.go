package handler

import (
	"golang-whatsapp-clone/server"
	"net/http"
)

// Handler is the main entry point of the application. Think of it like the main() method
func Handler(w http.ResponseWriter, r *http.Request) {
	app, _, rootHandler := server.NewServer()
	defer app.DBpool.Close()

	rootHandler.ServeHTTP(w, r)
}
