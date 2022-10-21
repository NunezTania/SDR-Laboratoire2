package main

import (
	"os"
	"os/exec"
	"strings"
)

func main() {
	p := exec.Command("go", "run", "./main/server/server.go")
	p.Stdin = strings.NewReader("LISTM")
	p.Stdout = os.Stdout
}
