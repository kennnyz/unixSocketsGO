package main

import (
	"bufio"
	"fmt"
	client2 "github.com/kennnyz/unixGo/client"
	"log"
	"os"
	"strings"
)

const (
	socketPath = "D:\\unixGoWB\\unixSocketsGO\\sock\\server.sock"
)

func main() {
	client, err := client2.NewClient(socketPath, "Telegram Alerts")
	if err != nil {
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		// Запрос ввода отправителя и сообщения
		fmt.Print("Enter your message: ")
		message, _ := reader.ReadString('\n')

		// Создание экземпляра структуры Message
		msg := client2.Message{
			From:    strings.TrimSpace(client.ServiceName),
			Message: strings.TrimSpace(message),
		}

		err, resp := client.WriteToServer(msg)
		if err != nil {
			log.Println(err)
		}

		// Вывод ответа от сервера
		fmt.Println("Server response:", resp)

		// Проверка условия завершения
		if strings.TrimSpace(message) == "exit" {
			break
		}
	}
	err = client.Close()
	if err != nil {
		log.Println(err)
	}
}
