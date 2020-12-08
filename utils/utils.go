package utils

import (
	"fmt"
	"io"
	"os"
)

var (
	GLOBALMAP map[string]interface{} = make(map[string]interface{})
)

func IsExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

//
func CopyPipe(dst io.Writer, src io.Reader) {
	go func() {
		if _, err := io.Copy(dst, src); err != nil {
			// TODO: print error
		}
	}()
}

func Fatal(message string) {
	fmt.Println("[-] " + message)
	os.Exit(1)
}

func Warn(message string) {
	fmt.Println("[!] " + message)
}

func Success(message string) {
	fmt.Println("[+] " + message)
}

func Exit() {
	fmt.Println("[*] Exit")
	os.Exit(0)
}

//
func Get(name string) interface{} {
	if v, ok := GLOBALMAP[name]; !ok {
		return nil
	} else {
		return v
	}
}

func Set(name string, value interface{}) {
	GLOBALMAP[name] = value
}
