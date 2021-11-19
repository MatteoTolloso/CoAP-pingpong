package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func pings(conn net.Conn, cont int, delay int) (error, []int, int) {

	var err error
	buff := make([]byte, 2048) // buffer di lettura

	rtts := make([]int, 0) // array di round-trip time
	var lost int = 0

	for ; cont > 0; cont-- {

		start := time.Now()

		_, err = conn.Write([]byte("ping")) // scrivo (invio un pacchetto) sulla connessione

		if err != nil {
			fmt.Print("impossibile inviare messaggio")
			return err, nil, 0
		}

		fmt.Println("inviato: ping")

		// setto una deadline per la lettura di 2 secondi, altrimenti considero il pacchetto perso
		conn.SetReadDeadline(time.Now().Add(time.Second * 2))
		_, err = conn.Read(buff)                    // lettura del pacchetto nel buffer
		elapsed := time.Since(start).Milliseconds() // calcolo tempo trascorso

		if err == nil {
			rtts = append(rtts, int(elapsed))
			fmt.Println("ricevuto: " + string(buff))
		} else if os.IsTimeout(err) {
			lost++
			fmt.Println("timeout scaduto")
		} else {
			fmt.Printf("\nErrore:  %v\n", err)
			//return err, nil, 0
		}

		err = nil
		time.Sleep(time.Duration(delay) * 1000000000)
	}

	return nil, rtts, lost

}

func main() {

	// controllo argomenti da linea di comando

	help := "Uso: <nome server remoto> <porta remota> <num of pings> <delay>"

	if len(os.Args) != 5 {
		fmt.Println(help)
		return
	}

	remoteName := os.Args[1]
	remotePort := os.Args[2]
	nOfPings, err1 := strconv.Atoi(os.Args[3])
	delay, err2 := strconv.Atoi(os.Args[4])

	if err1 != nil || err2 != nil {
		fmt.Println(help)
		return
	}

	// connessione con il server remoto, in realtà utilizzo udp quindi non viene realmente instaurata una connessione
	// se il client viene eseguito in un container, "remoteName" può essere anche il nome del sever, che poi viene
	// risolto tramite un servizio DNS del network di Docker

	conn, err := net.Dial("udp", remoteName+":"+remotePort)
	defer conn.Close()

	if err != nil {
		fmt.Printf("\nErrore Dial:  %v", err)
		return
	}

	err, rtts, losts := pings(conn, nOfPings, delay)

	if err != nil {
		fmt.Println(err)
		return
	}

	// statistiche finali su round-trip time e packet loss

	var avg, min, max int = 0, int(^uint(0) >> 1), 0

	for i := 0; i < len(rtts); i++ {
		avg += rtts[i]
		if rtts[i] > max {
			max = rtts[i]
		}
		if rtts[i] < min {
			min = rtts[i]
		}
	}
	avg = avg / (len(rtts) + 1)

	fmt.Printf("\n%d ping inviati, di cui %d ricevuti correttamente e %d persi\n", nOfPings, nOfPings-losts, losts)
	fmt.Printf("RTT medio: %d msec\nRTT massimo: %d msec\nRTT minimo: %d msec\n\n", avg, max, min)

	return

}
