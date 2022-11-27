package main

import (
	pm "SDR-Laboratoire1/main/server/processMutex"
	"fmt"
)

func main() {
	conf := pm.Config
	fmt.Println(conf)
	//server.LaunchNServ(conf.NServ)
}
