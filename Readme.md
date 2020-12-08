# Cake
***A cake for for connect local file or socket***

## Depends
  - emperror.dev/errors v0.8.0
  - github.com/jawher/mow.cli v1.2.0

## Usage
```go
package main

import (
	"fmt"

	"github.com/WAY29/cake/api"
)

func main() {
	c := api.Connect("./test.exe")  // local file
	// c := api.Connect("1.1.1.1:2333")  // remote socket
	n, data := c.Recv(100, 10)  // recvnumber, timeout
	if n > 0 {
		fmt.Print(string(data))
  }
  //n, err := c.Sendline([]byte("qwe"))  // send data
	n, err := c.Sendline([]byte("qwe"))  // send data with '\n'
	if err != nil {
		fmt.Println("error", err)
	}

	data = c.Recvuntil([]byte("we"), true)  // recv data until bytes, drops utilsbytes if set true
    fmt.Print(string(data))
    c.InteractiveR()  // interactive with it, use :r to reconnect, :exit to exit
	c.Close()  // close connect
}

```