package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	cli.AppHelpTemplate = `NAME:
	{{.Name}} - {{.Usage}}
USAGE:
	fava-auth --username/-u <username> --password/-p <password> --server/-s <fava ip:port> [--port/-P <listen port>]
	`
	app := cli.NewApp()
	app.Version = "2019.4"
	app.Name = "fava-auth"
	app.Usage = "beancount fava auth proxy"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "username",
			Aliases:  []string{"u"},
			Usage:    "auth username",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "password",
			Aliases:  []string{"p"},
			Usage:    "auth password",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "server",
			Aliases:  []string{"s"},
			Usage:    "fava server address",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "port",
			Aliases:  []string{"P"},
			Usage:    "fava auth listen port",
			Required: false,
		},
	}

	app.Action = func(c *cli.Context) error {
		username := c.String("username")
		server := c.String("server")
		password := c.String("password")
		port := c.Int("port")
		if port == 0 {
			port = 9000
		}

		config = &Config{
			Port:     port,
			Username: username,
			Password: password,
			Server:   server,
		}
		// start proxy server
		startServer()
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
