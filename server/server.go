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

		log.Println("New connection to the server: ", conn.LocalAddr().String())
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
			err := s.tryRead(conn, decoder)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

// Паттерн try/must - держит в себе всю функциональность получения и отправке ответа и возвращает ошибку в ее случае
func (s *Server) tryRead(conn net.Conn, decoder *json.Decoder) error {
	var msg client.Message
	err := decoder.Decode(&msg)
	if err != nil {
		if err == io.EOF {
			// Пользователь закрыл соединение
			return nil
		}
		return fmt.Errorf("decode error: %v", err)
	}
	s.msgChan <- msg
	_, err = fmt.Fprintf(conn, respProcess(msg.Message))
	if err != nil {
		return fmt.Errorf("error sending response to client: %v", err)
	}
	if msg.Message == "exit" {
		return nil
	}
	return nil
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
