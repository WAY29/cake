package cmd

import (
	"github.com/WAY29/cake/api"
	"github.com/WAY29/cake/utils"

	cli "github.com/jawher/mow.cli"
)

//
func ConnectCmd(cmd *cli.Cmd) {
	cmd.Spec = "FD"

	var (
		fdStr = cmd.StringArg("FD", "", "Process or socket to connect. Ex. ./pwn or 1.1.1.1:80")
	)
	cmd.Action = func() {
		utils.Set("Restart", true)
		for {
			utils.Set("Restart", false)
			api.ConnectInteractive(*fdStr)
			if !utils.Get("Restart").(bool) {
				break
			}
		}
		utils.Success("finish")
	}
}
