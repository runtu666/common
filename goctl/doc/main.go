package main

import (
	"fmt"
	"go-zero-demo/goctl/doc/action"
	"os"
	"runtime"

	"github.com/urfave/cli/v2"
)

const (
	version = "20220303"
)

func main() {
	app := cli.NewApp()
	app.Usage = "a plugin of goctl to generate doc files"
	app.Version = fmt.Sprintf("%s %s/%s", version, runtime.GOOS, runtime.GOARCH)
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "api",
			Required: false,
			Usage:    "the api file",
		},
		&cli.StringFlag{
			Name:     "o",
			Required: false,
			Usage:    "the output markdown file",
		},
		&cli.StringFlag{
			Name:     "mainTemplate",
			Required: false,
			Usage:    "the output metadata template file",
		},
		&cli.StringFlag{
			Name:     "routesTemplate",
			Required: false,
			Usage:    "the output routes template file",
		},
	}
	app.Action = action.DocAction

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("goctl-doc: %+v\n", err)
	}
}
