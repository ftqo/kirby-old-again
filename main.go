package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/fittsqo/kirby/discord"
)

var (
	BotToken = flag.String("token", "", "Bot access token")
)

func main() {
	flag.Parse()
	discord.Start(*BotToken)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutting down")
	discord.Stop()
}
