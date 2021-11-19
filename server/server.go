package main

import (
	"fmt"
	"net"
	"os"
)

func sendResponse(conn net.PacketConn, remoteAddr net.Addr, s string) {

	_, err := conn.WriteTo([]byte(s), remoteAddr) // invio la risposta

	if err != nil {
		fmt.Printf("\nErrore invio risposta %v", err)
	}
}

func main() {

	// controllo argomenti linea di comando

	if cap(os.Args) != 2 {
		fmt.Println("Uso: <porta> ")
		return
	}

	port := os.Args[1]

	conn, err := net.ListenPacket("udp", ":"+port) // mi metto in ascolto sulla porta (tutti gli ip della macchina)
	defer conn.Close()

	if err != nil {
		fmt.Printf("Error %v\n", err)
		return
	}

	buf := make([]byte, 2048) // buffer per lettura/scrittura

	fmt.Println("Server pronto")

	for {
		_, remoteAddr, err := conn.ReadFrom(buf) // chiamata bloccante per la lettura

		if err != nil {
			fmt.Printf("\nErrore nella lettura:  %v", err)
			continue
		}

		fmt.Printf("Recived: \"%s\" from: %v \n", buf, remoteAddr)

		var resp string

		if string(buf)[:4] == "ping" { // preparo la risposta
			resp = "pong"
		} else {
			resp = "bad request"
		}

		go sendResponse(conn, remoteAddr, resp)

	}
}
