package main

import (
	"fmt"
	"golang-whatsapp-clone/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	server := server.NewServer()

	go func() {

		if err := server.App.Listen(":" + server.AppConfig.Port); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	<-c // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")
	_ = server.App.Shutdown()

	fmt.Println("Running cleanup tasks...")

	// cleanup tasks
	server.DBpool.Close()
	fmt.Println("Fiber was successful shutdown.")
}
