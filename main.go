package main

import (
	"context"
	"fmt"
	"golang-whatsapp-clone/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	app, server := server.NewServer()

	go func() {

		if err := server.ListenAndServe(); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	<-c // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")
	_ = server.Shutdown(context.Background())

	fmt.Println("Running cleanup tasks...")

	// cleanup tasks
	app.DBpool.Close()
	fmt.Println("Server was successful shutdown.")
}
