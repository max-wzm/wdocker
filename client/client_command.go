package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/max-wzm/geerpc/xclient"

	"wdocker/cgroups/subsystems"
	"wdocker/container"
	"wdocker/log"
	"wdocker/network"
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
		cli.StringFlag{
			Name:  "net, n",
			Usage: "connect to network",
		},
		cli.StringSliceFlag{
			Name:  "e",
			Usage: "set environment",
		},
		cli.StringSliceFlag{
			Name:  "p",
			Usage: "port mapping",
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
			Env:    ctx.StringSlice("e"),
		}
		log.Info("runningConfig: %v", runningConfig)

		info := container.ContainerInfo{
			ID:          id,
			Name:        name,
			InitCmd:     strings.Join(cmds, " "),
			CreatedTime: time.Now().Format("2006-01-02 15:04:05"),
			Network:     ctx.String("net"),
			PortMapping: ctx.StringSlice("p"),
		}
		con := &container.Container{
			ContainerInfo:  info,
			ImagePath:      imagePath,
			ResourceConfig: res,
			RunningConfig:  runningConfig,
			Stdout:         os.Stdout,
			Stdin:          os.Stdin,
			Stderr:         os.Stderr,
		}
		log.Info("container: %v", con)
		b, _ := os.ReadFile("../registry")
		raddr := string(b)
		var err int
		xclient.SimpleCall(raddr, "Daemon.RunContainer", con, &err)
		fmt.Println(err)
		return nil
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
		b, _ := os.ReadFile("../registry")
		raddr := string(b)
		log.Info(raddr)
		var err error
		xclient.SimpleCall(raddr, "Daemon.ListContainers", struct{}{}, err)
		log.Info("end")
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

var stopCommand = cli.Command{
	Name:  "stop",
	Usage: "stop a container",
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("require container name!")
		}
		return container.StopContainer(ctx.Args().Get(0))
	},
}

var removeCommand = cli.Command{
	Name:  "rm",
	Usage: "remove an exited container",
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("require container name!")
		}
		return container.RemoveContainer(ctx.Args().Get(0))
	},
}

var networkCommand = cli.Command{
	Name:  "network",
	Usage: "network related",
	Subcommands: []cli.Command{
		{
			Name:  "create",
			Usage: "create network",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "subnet, n",
				},
				cli.StringFlag{
					Name: "driver, d",
				},
			},
			Action: func(ctx *cli.Context) error {
				network.Init()
				err := network.CreateNetwork(ctx.String("driver"), ctx.String("subnet"), ctx.Args().Get(0))
				return err
			},
		},
		{
			Name:  "ls",
			Usage: "list networks",
			Action: func(ctx *cli.Context) error {
				network.Init()
				network.ListNetwork()
				return nil
			},
		},
		{
			Name:  "rm",
			Usage: "remove network",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "name",
				},
			},
			Action: func(ctx *cli.Context) error {
				network.Init()
				err := network.DeleteNetwork(ctx.Args().Get(0))
				return err
			},
		},
	},
}
