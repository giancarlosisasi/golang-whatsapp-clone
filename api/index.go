package handler

import (
	"net/http"

	"golang-whatsapp-clone/server"

	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

// Handler is the main entry point of the application. Think of it like the main() method
func Handler(w http.ResponseWriter, r *http.Request) {
	// This is needed to set the proper request path in `*fiber.Ctx`
	r.RequestURI = r.URL.String()

	server := server.NewServer()
	defer server.DBpool.Close()

	h := adaptor.FiberApp(server.App)
	h.ServeHTTP(w, r)

}
