package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ftqo/kirby/discord"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env files")
	}

	discord.Start(os.Getenv("TOKEN"))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	println()
	log.Print("gracefully shutting down !")
	discord.Stop()
}
