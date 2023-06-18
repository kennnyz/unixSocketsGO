package client

import (
	"encoding/json"
	"fmt"
	"net"
)

// Client Получает и записывает сообщения
type Client struct {
	ServiceName string
	socketPath  string
	conn        net.Conn
}

type Message struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

func NewClient(socketAddr, serviceName string) (*Client, error) {
	conn, err := net.Dial("unix", socketAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %v", err)
	}

	return &Client{
		conn:        conn,
		socketPath:  socketAddr,
		ServiceName: serviceName,
	}, nil
}

func (c *Client) WriteToServer(msg Message) (error, string) {
	// Кодирование структуры в JSON
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to encode message: %v", err), ""

	}
	// Отправка JSON-сообщения на сервер
	_, err = c.conn.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to send message to server: %v", err), ""
	}

	// Чтение ответа от сервера
	buffer := make([]byte, 2048)
	n, err := c.conn.Read(buffer)
	if err != nil {
		return fmt.Errorf("failed to receive response from server: %v", err), ""
	}

	return nil, string(buffer[:n])
}

func (c *Client) Close() error {
	err := c.conn.Close()
	if err != nil {
		return err
	}
	return nil
}
