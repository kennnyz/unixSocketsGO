package main

import client2 "github.com/kennnyz/unixGo/client"

const (
	socketPath = "D:\\unixGoWB\\unixSocketsGO\\sock\\server.sock"
)

func main() {
	client := client2.NewClient(socketPath, "Telegram Alerts")
	err := client.ConnectAndWriteToServer()
	if err != nil {
		return
	}
}
