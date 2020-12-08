package api

import (
	"bufio"
	"io"
	"net"
	"os/exec"
	"sync"
)

type Connection struct {
	FdStr        string
	Fd           interface{}
	IsStart      bool
	Lock         *sync.Mutex
	StdinPipe    io.WriteCloser
	StdoutPipe   io.ReadCloser
	StderrPipe   io.ReadCloser
	StdoutReader *bufio.Reader
	StderrReader *bufio.Reader
	Command      *exec.Cmd
	Connect      net.Conn
}
