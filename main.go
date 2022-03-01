package main

import (
	"flag"
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

var (
	e      map[string]string
	ep     string
	debug  bool
	resume bool
)

func init() { // initialize logger, parse flags and env variables
	var err error // prevent shadowing with :=
	logger.Initialize()

	flag.BoolVar(&debug, "d", false, "enable debugging logs")
	flag.BoolVar(&resume, "r", false, "resume discord events if possible")
	flag.Parse()
	if !debug {
		logger.NoDebug()
	}
	logger.L.Info().Msg("Parsed command-line flags")

	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	ep = path.Join(d, "/.env")
	e, err = godotenv.Read(ep)
	if err != nil {
		logger.L.Panic().Err(err).Msg("Failed to load .env variables")
	}
	if !resume {
		e["DISCORD_SESSIONID"] = ""
		e["DISCORD_SEQUENCE"] = ""
	}
	logger.L.Info().Msg("Parsed environment variables")
}

func main() {
	database.Open(e["DB_HOST"], e["DB_PORT"], e["DB_USER"], e["DB_PASS"], e["DB_DATABASE"])
	assets.Load()
	go api.Start(e["API_PORT"])
	go discord.Start(e["DISCORD_TOKEN"], e["DISCORD_SESSIONID"], e["DISCORD_SEQUENCE"], e["DISCORD_TEST_GUILD"])
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	logger.L.Info().Msg("Initiated shutdown process")
	sid, seq := discord.Stop()
	e["DISCORD_SESSIONID"] = sid
	e["DISCORD_SEQUENCE"] = seq
	godotenv.Write(e, ep)
	database.Close()
	logger.Close()
}
