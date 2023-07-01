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
			Name:  "it",
			Usage: "enable interactive terminal",
		},
		cli.StringFlag{
			Name:  "mem, m",
			Usage: "memory",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "container name",
		},
		cli.StringFlag{
			Name:  "v",
			Usage: "volume",
		},
		cli.BoolFlag{
			Name:  "rm",
			Usage: "remove after exit",
		},
		cli.BoolFlag{
			Name: "d",
			Usage: "detach container",
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

		runningConfig := &container.RunningConfig{
			Tty: ctx.Bool("it"),
			Remove: ctx.Bool("rm"),
			Volume: ctx.String("v"),
		}
		log.Info("runningConfig: %v", runningConfig)

		container := &container.Container{
			ID:             id,
			Name:           name,
			ImagePath:      imagePath,
			ResourceConfig: res,
			RunningConfig:  runningConfig,
			InitCmds:       cmds,
		}

		log.Info("container info: %v", container)

		return Run(container)
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
