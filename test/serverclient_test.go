// The test package contains the tests of the program.
// It acts as a client and sends commands to the server.
// The server will then process the commands and send back the result.
package test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func Test(t *testing.T) {
	fmt.Println("Test")

	go func() {
		p := exec.Command("go", "run", "../main/server/server.go")
		p.Stdout = os.Stdout
		p.Stderr = os.Stderr
		err := p.Run()
		if err != nil {
			fmt.Println(err)
		}
	}()

	p := exec.Command("go", "run", "../main/client/client.go")
	buffer := bytes.Buffer{}
	buffer.Write([]byte("LISTM\n"))
	buffer.Write([]byte("QUIT\n"))
	p.Stdin = &buffer
	p.Stdout = os.Stdout
	p.Stderr = os.Stderr
	err := p.Run()
	if err != nil {
		fmt.Println(err)
	}
}
