package main

import (
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/ftqo/kirby/api"
	"github.com/ftqo/kirby/assets"
	"github.com/ftqo/kirby/database"
	"github.com/ftqo/kirby/discord"
	"github.com/ftqo/kirby/logger"
)

func main() {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	err := godotenv.Load(d + "/.env")
	if err != nil {
		logger.L.Panic().Err(err).Msg("Failed to load .env variables")
	}
	logger.L.Info().Msg("Loaded environment variables")
	database.Open(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_DATABASE"))
	assets.Load()

	go api.Start(os.Getenv("API_PORT"))
	go discord.Start(os.Getenv("DISCORD_TOKEN"), os.Getenv("DISCORD_TEST_GUILD"), os.Getenv("DISCORD_RMCMDS"))
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	logger.L.Info().Msg("Initiated shutdown process")
	discord.Stop()
	database.Close()
}
