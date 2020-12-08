package api

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"strings"
	"time"

	"github.com/WAY29/cake/utils"
)

var (
	TimeoutSecond  int = 0
	RecvNumber     int = 0
	RealRecvNumber int = 0
	RecvUntilBytes []byte
	RecvFlag       int = 0
)

// ? send data
func (conn Connection) Send(data []byte) (n int, err error) {
	return conn.StdinPipe.Write(data)
}

// ? send data with '\n'
func (conn Connection) Sendline(data []byte) (n int, err error) {
	data = append(data, 0x0a)
	return conn.StdinPipe.Write(data)
}

//
func read(reader *bufio.Reader, buffer []byte) ([]byte, error) {
	var err error
	switch RecvFlag {
	case 0:
		_, err = reader.Read(buffer)
	case 1:
		buffer, _, err = reader.ReadLine()
	case 2:
		buffer, err = ioutil.ReadAll(reader)
	case 3:
		for {
			s := ""
			s, err = reader.ReadString(RecvUntilBytes[len(RecvUntilBytes)-1])
			if err != nil {
				return nil, err
			}

			buffer = append(buffer, []byte(s)...)
			if bytes.HasSuffix(buffer, RecvUntilBytes) {
				return buffer, nil
			}
		}
	}
	RealRecvNumber = len(buffer)
	return buffer, err
}

func recv(conn Connection) []byte {
	deadline := time.Now().Add(time.Duration(TimeoutSecond) * time.Second)
	var buffer []byte
	var err error

	if RecvNumber > 0 {
		buffer = make([]byte, RecvNumber)
	}

	stdoutReader := conn.StdoutReader
	if conn.Connect != nil { // ? socket
		conn := conn.Connect
		conn.SetDeadline(deadline) // ? timeout
		buffer, err = read(stdoutReader, buffer)

		if err != nil {
			if strings.Contains(err.Error(), "timeout") {
				utils.Fatal("Read timeout")
			} else {
				utils.Fatal(err.Error())
			}
		}

		return buffer
	} else { // ? process
		result := make(chan bool, 1)
		go func() {
			buffer, err = read(stdoutReader, buffer)

			if err == nil {
				result <- true
			}
			result <- false
		}()

		select {
		case v, ok := <-result:
			if !v || !ok {
				return nil
			}
			return buffer
		case <-time.After(time.Duration(TimeoutSecond) * time.Second):
			utils.Fatal("Read timeout")
		}
	}
	return nil
}

// ? Recv {number} length data, timetout is {timeoutSecond} second.
func (conn Connection) Recv(number int, timeoutSecond ...int) (int, []byte) {
	RecvNumber = number
	if len(timeoutSecond) > 0 {
		TimeoutSecond = timeoutSecond[0]
	} else {
		TimeoutSecond = 300
	}
	RecvFlag = 0
	buffer := recv(conn)
	return RealRecvNumber, buffer
}

// ? Recv a line data.
func (conn Connection) Recvline() []byte {
	TimeoutSecond = 300
	RecvFlag = 1
	buffer := recv(conn)
	return buffer
}

// ? Recv all data.
func (conn Connection) Recvall() []byte {
	TimeoutSecond = 300
	RecvFlag = 2
	buffer := recv(conn)
	return buffer
}

// ? Recv data until {data}, set true to drop {data}.
func (conn Connection) Recvuntil(data []byte, drop ...bool) []byte {
	TimeoutSecond = 300
	RecvFlag = 3
	RecvUntilBytes = data
	buffer := recv(conn)
	if len(drop) > 0 && drop[0] {
		buffer = buffer[:len(buffer)-len(data)]
	}
	return buffer
}
