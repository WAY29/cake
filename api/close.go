package api

import (
	"os"
)

var (
	_Stdin  = os.Stdin
	_Stdout = os.Stdout
	_Stderr = os.Stderr
)

// ? Close connect
func (conn Connection) Close() {
	command := conn.Command
	os.Stdout = _Stdout
	os.Stderr = _Stderr
	conn.StdinPipe.Close()
	conn.StdoutPipe.Close()
	conn.StderrPipe.Close()
	if command != nil {
		command.Process.Kill()
	} else {
		conn.Connect.Close()
	}
}
