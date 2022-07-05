package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func handleReceiving(conn net.Conn, c chan Message) {
	// handle all incomming messages for a connection
	for {
		buff := make([]byte, 1024)
		n, _ := conn.Read(buff)
		m, success := decodeMessage(string(buff[:n]))
		if success {
			c <- m
		} else {
			fmt.Println("[+] Couldn't process message - server down? ")
			conn.Close()
			os.Exit(1)
		}
	}
}

func handleSending(conn net.Conn, c chan Message) {
	for {
		msg := <-c
		conn.Write([]byte(msg.Encode()))
	}
}

func client(addr string) (chan Message, chan Message) {

	// channels for sending and receiving
	receiving := make(chan Message)
	sending := make(chan Message)

	// make a connection
	conn, err := net.Dial("tcp", addr) //Connect Server port at 1000
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("[+] Connected to %s\n", addr)

	// initialize receiving and sending
	go handleReceiving(conn, receiving) //run the goroutine or thread what accept data
	go handleSending(conn, sending)     // run other thread what send data

	return receiving, sending
}
