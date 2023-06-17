package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// Client Получает и записывает сообщения
type Client struct {
	serviceName string
	socketPath  string
	conn        net.Conn
}

type Message struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

func NewClient(socketAddr, serviceName string) *Client {
	return &Client{
		socketPath:  socketAddr,
		serviceName: serviceName,
	}
}

func (c *Client) writeToServer(msg Message) error {
	// Кодирование структуры в JSON
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to encode message: %v", err)

	}
	// Отправка JSON-сообщения на сервер
	_, err = c.conn.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to send message to server: %v", err)
	}
	return nil
}

// TODO может это вообще не надо? убрать и писать в main (от туда где создается Client)

func (c *Client) ConnectAndWriteToServer() error {
	// Установка соединения с сервером через UNIX-сокет
	conn, err := net.Dial("unix", c.socketPath)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	c.conn = conn
	defer conn.Close()

	// Чтение ввода пользователя
	reader := bufio.NewReader(os.Stdin)

	for {
		// Запрос ввода отправителя и сообщения
		fmt.Print("Enter your message: ")
		message, _ := reader.ReadString('\n')

		// Создание экземпляра структуры Message
		msg := Message{
			From:    strings.TrimSpace(c.serviceName),
			Message: strings.TrimSpace(message),
		}

		err := c.writeToServer(msg)

		// Чтение ответа от сервера
		buffer := make([]byte, 2048)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println("Failed to receive response from server:", err)
			continue
		}

		// Вывод ответа от сервера
		fmt.Println("Server response:", string(buffer[:n]))

		// Проверка условия завершения
		if strings.TrimSpace(message) == "exit" {
			break
		}
	}
	return nil
}
