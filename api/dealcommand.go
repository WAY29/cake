package api

import "github.com/WAY29/cake/utils"

// ? deal command which start with :
func DealCommand(cmdline string, conn Connection) {
	switch cmdline {
	case "e":
		fallthrough
	case "q":
		fallthrough
	case "exit":
		fallthrough
	case "quit": // ? exit cake
		Exit = true
		return
	case "r":
		fallthrough
	case "restart": // ? restart connect
		Restart = true             // ? for InputAndSend return
		utils.Set("Restart", true) // ? for cycle
		conn.Close()
		return
	}
	return
}
