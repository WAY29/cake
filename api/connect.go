package api

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/WAY29/cake/utils"

	"emperror.dev/errors"
)

// ? Start process
func (conn Connection) Start() {
	conn.Lock.Lock()
	defer conn.Lock.Unlock()
	if conn.IsStart {
		return
	}

	command := conn.Command
	fdStr := conn.FdStr

	if err := command.Start(); err != nil {
		fmt.Println("err", err)
		utils.Fatal(errors.Errorf("Run %s error\n", fdStr).Error())
	}
}

// ? waitNet
func waitNet(conn Connection) {
	var data = make([]byte, 4096)
	stdoutPipe := conn.StdoutPipe
	for {
		if n, err := stdoutPipe.Read(data); err == io.EOF {
			conn.Connect.Close()
			break
		} else if err != nil {
			if !Restart {
				utils.Fatal(errors.New("Connect error").Error())
			}
			break
		} else if n > 0 {
			os.Stdout.Write(data)
		}
	}
}

// ? Connect process or socket
func Connect(fdStr string, notStartNow ...bool) (conn Connection) {
	conn.IsStart = false
	conn.FdStr = fdStr
	conn.Lock = new(sync.Mutex)

	if strings.Contains(fdStr, ":") { // ? socket
		connect, err := net.DialTimeout("tcp", fdStr, 30*time.Second)
		if err != nil {
			utils.Fatal(errors.New("Connect error").Error())
		}

		// ? set conn
		conn.Connect = connect
		conn.StdinPipe = connect
		conn.StdoutPipe = connect
		conn.StderrPipe = connect
		conn.Fd = connect
	} else { // ? process
		command := exec.Command(fdStr)
		if !utils.IsExist(fdStr) {
			utils.Fatal(errors.New("File not found").Error())
		}
		// ? fork stdin pipe
		Cstdin, err := command.StdinPipe()
		if err != nil {
			utils.Fatal(errors.New("Fork stdin pipe error").Error())
		}
		conn.StdinPipe = Cstdin

		// ? fork stdout pipe
		Cstdout, err := command.StdoutPipe()
		if err != nil {
			utils.Fatal(errors.New("Fork stdout pipe error").Error())
		}
		conn.StdoutPipe = Cstdout

		// ? fork stderr pipe
		Cstderr, err := command.StderrPipe()
		if err != nil {
			utils.Fatal(errors.New("Fork stdout pipe error").Error())
		}
		conn.StderrPipe = Cstderr

		// ? set conn
		conn.Fd = command
		conn.Command = command
		// ? start process now if not interactive
		if len(notStartNow) == 0 {
			conn.Start()
		}
	}

	conn.StdoutReader = bufio.NewReader(conn.StdoutPipe)
	conn.StderrReader = bufio.NewReader(conn.StderrPipe)
	return
}

// ? Connect process or socket and interactive it
func ConnectInteractive(fdStr string) Connection {
	conn := Connect(fdStr, true)
	conn.Interactive()
	if conn.Command != nil {
		conn.Start()
	}
	if conn.Command != nil {
		conn.Command.Wait()
	}
	return conn
}
