package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ftqo/kirby/api"
	"github.com/ftqo/kirby/discord"
	"github.com/ftqo/kirby/logger"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.L.Err(err).Msgf("Failed to load .env variables")
	}

	go api.Start(os.Getenv("API_PORT"))
	go discord.Start(os.Getenv("DISCORD_TOKEN"), os.Getenv("TEST_GUILD"), os.Getenv("RMCMDS")) // TEST_GUILD only set in dev environment

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	println()
	logger.L.Info().Msg("gracefully shutting down !")
	discord.Stop() // does this properly stop the goroutine too?
}
