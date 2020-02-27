package ws

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"testing"
)

type tMessage struct {
	Msg chan []byte
}

func TestGetMessage(t *testing.T) {
	log.Fatal()
	fmt.Println("sss")
	w := &tMessage{
		Msg: make(chan []byte),
	}
	go func() {
		w.Msg <- []byte("hello")
	}()

	for {
		select {
		case rMsg := <-w.Msg:
			fmt.Println(rMsg)
			TTT(rMsg)
			break
		}
	}

	//var f = make(chan bool)
	// <- f
}

func TTT(m []byte) {
	for i := 0; i < 9999; i++ {
		fmt.Println(string(m[:]))
	}
}

var addr = flag.String("addr", "localhost:8080", "http service address")

func TestServer(t *testing.T) {
	flag.Parse()
	log.SetFlags(0)
	log.Fatal(http.ListenAndServe(*addr, Wsp))
}
