package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kennnyz/unixGo/client"
	"io"
	"log"
	"net"
	"strings"
)

type Server struct {
	listenAddress string
	listener      net.Listener
	msgChan       chan client.Message
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddress: listenAddr,
		msgChan:       make(chan client.Message, 10),
	}
}

func (s *Server) Start(ctx context.Context) error {
	listener, err := net.Listen("unix", s.listenAddress)
	if err != nil {
		return err
	}
	log.Println("Unix server is ALIVE!")
	s.listener = listener

	go s.acceptLoop(ctx)

	go func() {
		for msg := range s.msgChan {
			log.Printf("%s: %s\n", msg.From, msg.Message)
		}
	}()

	return nil
}

func (s *Server) acceptLoop(ctx context.Context) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("accept error ", err)
			continue
		}

		log.Println("New connection to the server: ", conn.RemoteAddr().String())
		go s.readLoop(ctx, conn)
	}
}

func (s *Server) readLoop(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)
	for {
		select {
		case <-ctx.Done():
			// Контекст завершен, выходим из цикла
			return
		default:
			var msg client.Message
			err := decoder.Decode(&msg)
			if err != nil {
				if err == io.EOF {
					// Пользователь закрыл соединение
					return
				}
				log.Println("decode error:", err)
				return
			}
			s.msgChan <- msg
			_, err = fmt.Fprintf(conn, respProcess(msg.Message))
			if err != nil {
				log.Println("Error sending response to client: ", err)
				return
			}
			if msg.Message == "exit" {
				return
			}
		}
	}
}

func (s *Server) Close(ctx context.Context) error {
	// check if chan already closed
	select {
	case _, ok := <-s.msgChan:
		if ok {
			close(s.msgChan)
		}
	default:
	}

	err := s.listener.Close()
	if err != nil {
		return err
	}
	ctx.Done()
	return nil
}

func respProcess(msg string) string {
	// TODO implement request process between services
	return fmt.Sprintf("Your message successfully has been received. Processed message is:  %s", strings.ToUpper(msg))
}
