package processMutex

import (
	"SDR-Laboratoire1/main/dataRW"
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type Conf struct {
	NServ      int            `yaml:"nServ"`
	PortServ   int            `yaml:"portServ"`
	PortClient int            `yaml:"portClient"`
	Host       string         `yaml:"host"`
	Type       string         `yaml:"type"`
	Users      []dataRW.User  `yaml:"users"`
	Events     []dataRW.Event `yaml:"events"`
	Debug      int            `yaml:"debug"`
}

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

var Config = ReadConfigFile(basepath + "/config.yaml")

func RunBtwServer(id int, clock *Lamport, inSC *bool, ChannelSC *chan string, DataChannel *chan chan []byte, done chan bool, listenConn net.Listener) {
	for {
		conn, err := listenConn.Accept()
		if err != nil {
			log.Fatal(err)
		}
		buf := make([]byte, 1024)
		_, err = conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		go handleMessage(buf, id, clock, inSC, ChannelSC, DataChannel)
	}
	done <- true
}

func handleMessage(buf []byte, id int, clock *Lamport, inSC *bool, ChannelSC *chan string, DataChannel *chan chan []byte) {

	var res = strings.Split(string(buf), " ")

	if res[0] == "data" { // message is a data sync
		clientChannel := make(chan []byte)
		*DataChannel <- clientChannel
		clientChannel <- []byte(strings.Join(res[1:], " "))
		res := <-clientChannel
		fmt.Println(string(res))

	} else { // message is a SC message
		var msg = strToMessage(string(buf))
		if msg.rType == "req" {
			*clock = clock.Update(msg.time)
			NoteNewMessage(msg, msg.id, id, inSC, ChannelSC)
			var r = Message{"ack", *clock, id}
			SendMessageTo(msg.id, r)

		} else if msg.rType == "rel" {
			*clock = clock.Update(msg.time)
			NoteNewMessage(msg, msg.id, id, inSC, ChannelSC)

		} else if msg.rType == "ack" {
			*clock = clock.Update(msg.time)
			NoteNewMessage(msg, msg.id, id, inSC, ChannelSC)
		}
	}
}

func sendRequests(clock *Lamport, id int) {
	var request = Message{"req", *clock, id}
	SendToAll(request, id)
}

func sendReleases(clock *Lamport, id int) {
	var request = Message{"rel", *clock, id}
	SendToAll(request, id)
}

func ReadConfigFile(path string) Conf {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	var c Conf
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func strToMessage(str string) Message {
	var request Message
	var args = strings.Split(str, " ")
	request.rType = args[0]
	request.time = strToLamport(args[1])
	request.id, _ = strconv.Atoi(string(bytes.Trim([]byte(args[2]), "\x00")))
	return request
}

func MessageToStr(request Message) string {
	return request.rType + " " + strconv.Itoa(request.time.counterTime) + " " + strconv.Itoa(request.id)
}

func SendMessageTo(id int, request Message) {
	msg := MessageToStr(request)
	currConn, err := net.Dial("tcp", Config.Host+":"+strconv.Itoa(Config.PortServ+id))
	if err != nil {
		log.Fatal(err)
	}
	_, err = currConn.Write([]byte(msg))
	errClose := currConn.Close()
	if errClose != nil {
		log.Fatal(err)
	}
}

func SendDataSyncTo(id int, data []byte) {
	currConn, err := net.Dial(Config.Type, Config.Host+":"+strconv.Itoa(Config.PortServ+id))
	if err != nil {
		log.Fatal(err)
	}
	_, err = currConn.Write(data)
	errClose := currConn.Close()
	if errClose != nil {
		log.Fatal(err)
	}
}

func SendToAll(request Message, id int) {
	for i := 0; i < Config.NServ; i++ {
		if i != id {
			SendMessageTo(i, request)
		}
	}
}

func SendDataSyncToAll(command []byte, id int) {
	msg := append([]byte("data "), command...)
	for i := 0; i < Config.NServ; i++ {
		if i != id {
			SendDataSyncTo(i, msg)
		}
	}
}

func WaitForEveryBody(id int, listenConn net.Listener) {
	fmt.Println("I'm id = ", id, " and im Waiting for every body to be ready")
	msg := "ready"
	waitReady := make(chan bool)
	for i := 0; i < Config.NServ; i++ {
		if i != id {
			conn, err := net.Dial(Config.Type, Config.Host+":"+strconv.Itoa(Config.PortServ+i))
			for err != nil {
				fmt.Println("Could not dial server ", i, " trying again...")
				conn, err = net.Dial(Config.Type, Config.Host+":"+strconv.Itoa(Config.PortServ+i))
			}
			_, err = conn.Write([]byte(msg))
			errClose := conn.Close()
			if errClose != nil {
				log.Fatal(err)
			}
		}
	}

	go func() {
		for i := 0; i < Config.NServ-1; i++ {
			conn, err := listenConn.Accept()
			if err != nil {
				log.Fatal(err)
			}
			buf := make([]byte, 1024)
			_, err = conn.Read(buf)
		}
		waitReady <- true
	}()
	<-waitReady
	fmt.Println("I'm id = ", id, " and everybody told me they are ready")
}
