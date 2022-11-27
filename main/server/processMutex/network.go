package processMutex

// the network.go file is the one that handles the communication between the servers
// it contains the functions that are used to send messages to other servers
// and to handle the messages received from other servers
import (
	"SDR-Laboratoire1/main/dataRW"
	"bytes"
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

// RunBtwServer allow a server to listen to other servers
func RunBtwServer(id int, clock *Lamport, inSC *bool, ChannelSC *chan string, DataChannel *chan chan []byte, done chan bool, listenConn net.Listener, messages *[]Message) {
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
		go handleMessage(buf, id, clock, inSC, ChannelSC, DataChannel, messages)
	}
	done <- true
}

// handleMessage handles the messages received from other servers
func handleMessage(buf []byte, id int, clock *Lamport, inSC *bool, ChannelSC *chan string, DataChannel *chan chan []byte, messages *[]Message) {

	var res = strings.Split(string(buf), " ")

	if res[0] == "data" { // message is a data sync
		clientChannel := make(chan []byte)
		*DataChannel <- clientChannel
		clientChannel <- []byte(strings.Join(res[1:], " "))
		<-clientChannel

	} else { // message is a SC message
		var msg = strToMessage(string(buf))
		if msg.rType == "req" {
			*clock = clock.Update(msg.time)
			NoteNewMessage(msg, msg.id, id, inSC, ChannelSC, messages)
			var r = Message{"ack", *clock, id}
			SendMessageTo(msg.id, r)

		} else if msg.rType == "rel" || msg.rType == "ack" {
			*clock = clock.Update(msg.time)
			NoteNewMessage(msg, msg.id, id, inSC, ChannelSC, messages)

		}
	}
}

// sendRequests sends a request to all servers
func sendRequests(clock *Lamport, id int) {
	var request = Message{"req", *clock, id}
	SendToAll(request, id)
}

// sendRelease sends a release message to all servers
func sendReleases(clock *Lamport, id int) {
	var request = Message{"rel", *clock, id}
	SendToAll(request, id)
}

// ReadConfigFile reads the config file and returns a Conf struct
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

// strToMessage converts a string to a Message struct
func strToMessage(str string) Message {
	var request Message
	var args = strings.Split(str, " ")
	request.rType = args[0]
	request.time = strToLamport(args[1])
	request.id, _ = strconv.Atoi(string(bytes.Trim([]byte(args[2]), "\x00")))
	return request
}

// MessageToStr converts a Message struct to a string
func MessageToStr(request Message) string {
	return request.rType + " " + strconv.Itoa(request.time.counterTime) + " " + strconv.Itoa(request.id)
}

// SendMessageTo sends a message to a specific server
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

// SendDataSyncTo sends a data sync message to a specific server
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

// SendToAll sends a message to all servers
func SendToAll(request Message, id int) {
	for i := 0; i < Config.NServ; i++ {
		if i != id {
			SendMessageTo(i, request)
		}
	}
}

// SendDataSyncToAll sends a data sync message to all servers
func SendDataSyncToAll(command []byte, id int) {
	msg := append([]byte("data "), command...)
	for i := 0; i < Config.NServ; i++ {
		if i != id {
			SendDataSyncTo(i, msg)
		}
	}
}

// WaitForEveryBody waits for all servers to be ready before listening to clients
func WaitForEveryBody(id int, listenConn net.Listener) {
	msg := "ready"
	waitReady := make(chan bool)
	for i := 0; i < Config.NServ; i++ {
		if i != id {
			conn, err := net.Dial(Config.Type, Config.Host+":"+strconv.Itoa(Config.PortServ+i))
			for err != nil {
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
}
