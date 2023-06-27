package main

import (
	"fmt"

	"wdocker/cgroups/subsystems"
	"wdocker/container"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name:  "run",
	Usage: "create a container with namespace and cgroups limit \n my docker run -ti [command]",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name: "mem",
			Usage: "memory",
		},
	},
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("missing container command")
		}

		var cmds []string
		cmds = append(cmds, ctx.Args()...)

		tty := ctx.Bool("ti")
		res := &subsystems.ResourceConfig{
			MemoryLimit: ctx.String("mem"),
		}
		Run(tty, cmds, res)
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	Action: func(ctx *cli.Context) error {
		log.Infof("init come on")
		err := container.RunContainerInitProcess()
		return err
	},
}
