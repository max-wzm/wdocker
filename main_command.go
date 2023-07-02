package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"text/tabwriter"
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
		ListContainers()
		return nil
	},
}

func ListContainers() {
	dirEntries, _ := os.ReadDir("/wdocker")
	containers := make([]*container.Container, 0)
	for _, e := range dirEntries {
		fi, _ := e.Info()
		configURL := path.Join("/wdocker", fi.Name(), container.ConfigName)
		b, err := os.ReadFile(configURL)
		if err != nil {
			continue
		}
		var tmpCon container.Container
		json.Unmarshal(b, &tmpCon)
		containers = append(containers, &tmpCon)
	}

	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	for _, con := range containers {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t\n", con.ID, con.Name, con.PID, con.Status, con.InitCmd, con.CreatedTime)
	}
	w.Flush()
}
