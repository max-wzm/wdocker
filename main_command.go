package main

import (
	"fmt"

	"wdocker/cgroups/subsystems"
	"wdocker/container"
	"wdocker/log"
	"wdocker/utils"

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
			Name:  "mem, m",
			Usage: "memory",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "container name",
		},
	},
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 2 {
			return fmt.Errorf("missing image path or container command")
		}

		id := utils.RandomID()
		name := ctx.String("name")
		if name == "" {
			name = id
		}
		imagePath := ctx.Args().Get(0)

		var cmds []string
		cmds = append(cmds, ctx.Args()[1:]...)

		res := &subsystems.ResourceConfig{
			MemoryLimit: ctx.String("mem"),
		}
		log.Info("res: %v", res)

		container := &container.Container{
			ID:             id,
			Name:           name,
			ImagePath:      imagePath,
			ResourceConfig: res,
			InitCmds:       cmds,
		}

		log.Info("container info: %v", container)

		tty := ctx.Bool("ti")
		return Run(container, tty)
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	Action: func(ctx *cli.Context) error {
		log.Info("init come on")
		err := container.RunContainerInitProcess()
		return err
	},
}
