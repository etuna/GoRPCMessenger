package main

import (
	"bufio"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strings"
)

var participants []string

type MSGService struct {
}

type Message struct {
	Transcript string
	SID        string
}

func send(line string) {
	var msg MSGService
	msg2Send := Message{line, addr + "/" + port}
	msg.Multicast(msg2Send)
}

func (msg *MSGService) Multicast(msg2Send Message) error {
	for _, element := range participants {
		partAddr := strings.Split(element, "/")[0] + ":" + strings.Split(element, "/")[1]
		client, err := rpc.Dial("tcp", partAddr)
		if err != nil {
			fmt.Println("Error occured on multicast :", err)
		}
		var reply int
		err = client.Call("MSGService.MessagePost", msg2Send, &reply)
		if err != nil {
			fmt.Println("Error occured on multicast, calling :", err)
		}
		client.Close()
	}
	return nil

}

func (msg *MSGService) MessagePost(msg2Send *Message, reply *int) error {

	fmt.Println(strings.TrimRight(msg2Send.SID, "\r\n") + ":" + msg2Send.Transcript)
	*reply = 1
	return nil
}

var sid string
var addr string
var port string

func main() {

	rpc.Register(new(MSGService))
	//path, _ := os.Getwd()
	//fmt.Println("[PEER] Working directory :" + path)
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "port")
		os.Exit(1)
	}
	sid = strings.TrimRight(os.Args[1], "\r\n")
	addr = strings.Split(sid, "/")[0]
	_init()
	var line string
	reader := bufio.NewReader(os.Stdin)
	go handleClientRequest()
	fmt.Println("\n\n###   Welcome to DSMessenger   ###\n\n")
	for {
		line, _ = reader.ReadString('\n')
		send(line)
	}
}

func _init() {
	file, _ := os.Open("room.txt")
	reader := bufio.NewReader(file)
	var part string
	//var err error

	for {
		part, _ = reader.ReadString('\n')
		if part == "" {
			break
		}
		maddr := strings.Split(part, "/")[0]
		mport := strings.Split(part, "/")[1]
		part = strings.TrimRight(part, "\r\n")
		if sid == part {
			addr = maddr
			port = mport
		} else {
			participants = append(participants, part)
		}
	}
	file.Close()
}

func handleClientRequest() {
	// sarting server
	tcpAddr, er1 := net.ResolveTCPAddr("tcp", ":"+strings.TrimRight(port, "\r\n"))
	if er1 != nil {
		fmt.Println("[LISTEN] client request error, resolve tcp addr error: ", er1)
	}
	listener, er2 := net.ListenTCP("tcp", tcpAddr)
	if er2 != nil {
		fmt.Println("[LISTEN] listener, listen tcp error: ", er2)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Print("error listening...")
			continue
		}
		rpc.ServeConn(conn)
	}
}
