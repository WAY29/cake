package api

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/WAY29/cake/utils"

	"emperror.dev/errors"
)

var (
	Exit    = false
	Restart = false
	Reader  = bufio.NewReader(os.Stdin)
)

// ? recv from stdin and send
func InputAndSend(conn Connection) {
	Restart = false
	for {
		cmdline, err := Reader.ReadString('\n')
		n := len(cmdline)
		if err != nil {
			utils.Fatal(errors.New("Keyboard interrupt or stdin error").Error())
			break
		}
		if strings.HasPrefix(cmdline, ":") {
			cmdline := strings.TrimSpace(cmdline[1:])
			DealCommand(cmdline, conn)
			if Restart {
				return
			}
			if Exit { // ? exit if input exit
				utils.Exit()
			}
			continue
		}
		wn, err := conn.Send([]byte(cmdline))
		if err != nil {
			utils.Fatal(errors.New("Interactive error").Error())
			break
		}
		if n != wn {
			utils.Warn(fmt.Sprintf("%d length data recv but %d length data send", n, wn))
		}
	}
}

// ? Connect process or socket and interactive it
func (conn Connection) Interactive() {
	fd := conn.Fd
	switch fd.(type) {
	case *exec.Cmd: // ? process
	case net.Conn: // ? socket
		conn.StdinPipe = conn.Connect
	}
	go InputAndSend(conn)

	if conn.Command != nil {
		utils.CopyPipe(os.Stdout, conn.StdoutPipe)
		utils.CopyPipe(os.Stderr, conn.StderrPipe)
		conn.Command.Wait()
	} else {
		utils.CopyPipe(os.Stderr, conn.StderrPipe)
		waitNet(conn)
	}

}

// ? Connect process or socket and interactive it, support restart
func (conn Connection) InteractiveR() {
	utils.Set("Restart", false)
	conn.Interactive()
	for {
		if !utils.Get("Restart").(bool) {
			break
		}
		utils.Set("Restart", false)
		ConnectInteractive(conn.FdStr)
	}
}
