package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ftqo/kirby/discord"
	"github.com/joho/godotenv"
)

var (
	BotToken = flag.String("token", "", "Bot access token")
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env files")
	}

	flag.Parse()
	discord.Start(*BotToken)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	println()
	log.Print("gracefully shutting down !")
	discord.Stop()
}
