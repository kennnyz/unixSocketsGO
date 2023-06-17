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
	listenAddress string
	listener      net.Listener
	MsgChan       chan client.Message
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddress: listenAddr,
		MsgChan:       make(chan client.Message, 10),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("unix", s.listenAddress)
	if err != nil {
		return err
	}
	log.Println("Unix server is ALIVE!")
	s.listener = listener

	go s.AcceptLoop()

	return nil
}

func (s *Server) AcceptLoop() {
	for {
		conn, err := s.listener.Accept()
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
			fmt.Fprintf(conn, respProcess(msg.Message))
			return
		}
	}
}

func (s *Server) Close() error {
	// check if chan already closed
	select {
	case _, ok := <-s.MsgChan:
		if ok {
			close(s.MsgChan)
		}
	default:
	}

	err := s.listener.Close()
	if err != nil {
		return err
	}
	return nil
}

func respProcess(msg string) string {
	// TODO implement request process between services
	return fmt.Sprintf("Your message successfully has been received. Processed message is:  %s", strings.ToUpper(msg))
}
