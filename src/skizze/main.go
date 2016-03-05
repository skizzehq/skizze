package main

import (
	"os"

	_ "net/http/pprof"

	"github.com/codegangsta/cli"
	"github.com/njpatel/loggo"
	"golang.org/x/crypto/ssh/terminal"

	"config"
	"manager"
	"server"
)

var (
	host    string
	logger  = loggo.GetLogger("skizze")
	port    int
	version string
)

func init() {
	setupLoggers()
}

func main() {
	app := cli.NewApp()
	app.Name = "Skizze"
	app.Usage = "A sketch data store for counting and sketching using probabilistic data-structures"
	app.Version = version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "host",
			Value:       "localhost",
			Usage:       "the host interface to bind to",
			Destination: &host,
			EnvVar:      "SKIZZE_HOST",
		},
		cli.IntFlag{
			Name:        "port, p",
			Value:       3596,
			Usage:       "the port to bind to",
			Destination: &port,
			EnvVar:      "SKIZZE_PORT",
		},
	}

	app.Action = func(*cli.Context) {
		// FIXME: Allow specifying datadir and infodir
		logger.Infof("Starting Skizze...")
		logger.Infof("Using data dir: %s", config.GetConfig().DataDir)

		mngr := manager.NewManager()
		server.Run(mngr, host, uint(port))
	}

	if err := app.Run(os.Args); err != nil {
		logger.Criticalf(err.Error())
	}
}

func setupLoggers() {
	loggerSpec := os.Getenv("SKIZZE_DEBUG")

	// Setup logging and such things if we're running in a term
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		if loggerSpec == "" {
			loggerSpec = "<root>=DEBUG"
		}
		// As we're in a terminal, let's make the output a little nicer
		_, _ = loggo.ReplaceDefaultWriter(loggo.NewSimpleWriter(os.Stderr, &loggo.ColorFormatter{}))
	} else {
		if loggerSpec == "" {
			loggerSpec = "<root>=WARNING"
		}
	}

	_ = loggo.ConfigureLoggers(loggerSpec)
}
