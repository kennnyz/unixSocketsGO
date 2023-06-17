package server

import (
	"encoding/json"
	"fmt"
	"github.com/kennnyz/unixGo/client"
	"log"
	"net"
	"strings"
)

type Server struct {
	ListenAddress string
	Listener      net.Listener
	MsgChan       chan client.Message
}

func NewServer(listenAddr string) *Server {
	return &Server{
		ListenAddress: listenAddr,
		MsgChan:       make(chan client.Message, 10),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("unix", s.ListenAddress)
	if err != nil {
		return err
	}
	log.Println("Unix server is ALIVE!")
	s.Listener = listener

	go s.AcceptLoop()

	return nil
}

func (s *Server) AcceptLoop() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			log.Println("accept error ", err)
			continue
		}

		log.Println("New connection to the server: ", conn.RemoteAddr().String())
		go s.ReadLoop(conn)
	}
}

func (s *Server) ReadLoop(conn net.Conn) {
	defer conn.Close()
	// buf := make([]byte, 2048) // Можно получать так.
	decoder := json.NewDecoder(conn)
	for {
		var msg client.Message
		err := decoder.Decode(&msg)
		if err != nil {
			log.Println("decode error:", err)
			continue
		}
		s.MsgChan <- msg
		_, err = conn.Write([]byte(respProcess(msg.Message)))
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (s *Server) Close() {
	close(s.MsgChan)
	s.Listener.Close()
}

func respProcess(msg string) string {
	return fmt.Sprintf("Your message successfully has been received. Processed message is:  %s", strings.ToUpper(msg))
}
