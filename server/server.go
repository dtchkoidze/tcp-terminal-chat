package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

var (
	clients   = make(map[net.Conn]string)
	clientsMu sync.Mutex
)

func main() {
	var addr string = ":9000"
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		panic(err)
	}

	defer listener.Close()
	fmt.Println("chat server started on", addr)

	id := 1

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("err: ", err)
			continue
		}

		name := fmt.Sprintf("Client %d", id)

		id++

		clientsMu.Lock()
		clients[conn] = name
		clientsMu.Unlock()

		go handleClient(conn, name)
	}

}

func handleClient(conn net.Conn, name string) {
	defer func() {
		clientsMu.Lock()
		delete(clients, conn)
		clientsMu.Unlock()
		conn.Close()
		broadcast(fmt.Sprintf("%s has left the chat.\n", name), conn)
	}()

	broadcast(fmt.Sprintf("%s joined the chat.\n", name), conn)
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		msg := scanner.Text()
		fmt.Println("read by bufio scanner Text: ", msg)
		broadcast(fmt.Sprintf("%s: %s\n", name, msg), conn)
	}
}

func broadcast(message string, ignoreConn net.Conn) {
	fmt.Print(message)
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for conn := range clients {
		if conn != ignoreConn {
			fmt.Fprint(conn, message)
		}
	}
}
