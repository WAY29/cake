package main

import (
	"fmt"
	"os"

	"github.com/WAY29/cake/api"
	"github.com/WAY29/cake/cmd"

	cli "github.com/jawher/mow.cli"
)

var (
	RELEASE string
	EXD, _  = os.Executable()
	PWD, _  = os.Getwd()
)

//
func init() {
	if len(RELEASE) > 0 {

	}
}

//
func main() {
	app := cli.App("cake", "A cake for connect local file or socket")
	app.Version("v version", "cake Version: 0.2")

	app.Spec = "[-v]"

	app.Command("c connect", "Connect a socket or a process", cmd.ConnectCmd)
	app.Command("test", "", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			c := api.Connect("172.28.2.160:2333")
			n, data := c.Recv(100, 10)
			if n > 0 {
				fmt.Print(string(data))
			}
			n, err := c.Sendline([]byte("qwe"))
			if err != nil {
				fmt.Println("error", err)
			}

			data = c.Recvuntil([]byte("we"), true)
			fmt.Print(string(data))
			c.Close()
		}
	})
	app.Run(os.Args)
}
