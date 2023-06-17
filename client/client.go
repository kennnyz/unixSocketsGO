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
	SocketPath string
}

type Message struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

func NewClient(socketAddr string) *Client {
	return &Client{
		socketAddr,
	}
}

func (c *Client) Run() {
	// Установка соединения с сервером через UNIX-сокет
	conn, err := net.Dial("unix", c.SocketPath)
	if err != nil {
		log.Fatal("Failed to connect to server:", err)
	}
	defer conn.Close()

	// Чтение ввода пользователя
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter service name: ")
	from, _ := reader.ReadString('\n')

	for {
		// Запрос ввода отправителя и сообщения
		fmt.Print("Enter your message: ")
		message, _ := reader.ReadString('\n')

		// Создание экземпляра структуры Message
		msg := Message{
			From:    strings.TrimSpace(from),
			Message: strings.TrimSpace(message),
		}

		// Кодирование структуры в JSON
		jsonData, err := json.Marshal(msg)
		if err != nil {
			log.Println("Failed to encode message:", err)
			continue
		}

		// Отправка JSON-сообщения на сервер
		_, err = conn.Write(jsonData)
		if err != nil {
			log.Println("Failed to send message to server:", err)
			continue
		}

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
}
