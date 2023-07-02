package main

import (
	"fmt"
	"os"
	"strings"
	"time"

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
			Name:  "d",
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
			Tty:    ctx.Bool("it"),
			Remove: ctx.Bool("rm"),
			Detach: ctx.Bool("d"),
			Volume: ctx.String("v"),
		}
		log.Info("runningConfig: %v", runningConfig)

		info := container.ContainerInfo{
			ID:          id,
			Name:        name,
			InitCmd:     strings.Join(cmds, " "),
			CreatedTime: time.Now().Format("2006-01-02 15:04:05"),
		}
		container := &container.Container{
			ContainerInfo:  info,
			ImagePath:      imagePath,
			ResourceConfig: res,
			RunningConfig:  runningConfig,
		}

		log.Info("container: %v", container)
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

var listCommand = cli.Command{
	Name:  "ps",
	Usage: "list all containers",
	Action: func(ctx *cli.Context) error {
		container.ListContainers()
		return nil
	},
}

var execCommand = cli.Command{
	Name:  "exec",
	Usage: "enter a container and exec command",
	Action: func(ctx *cli.Context) error {
		log.Info("got env %s, %s", os.Getenv(container.ENV_EXEC_PID), os.Getenv(container.ENV_EXEC_CMD))
		if os.Getenv(container.ENV_EXEC_PID) != "" {
			log.Info("pid callback pid %s", os.Getgid())
			return nil
		}
		if len(ctx.Args()) < 2 {
			return fmt.Errorf("missing con name or cmd")
		}
		containerName := ctx.Args().Get(0)
		var cmds []string
		cmds = append(cmds, ctx.Args()[1:]...)
		container.ExecContainer(containerName, cmds)
		return nil
	},
}
