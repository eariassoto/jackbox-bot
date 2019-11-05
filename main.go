package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

const (
	roomInfoURL = "http://blobcast.jackboxgames.com/room/%s/?userId=%s"
)

type responseMap map[string]interface{}

func getRoomInfo(roomID string, userID uuid.UUID) responseMap {
	url := fmt.Sprintf(roomInfoURL, roomID, userID.String())

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	return response
}

func getWebsocketURL(roomInfo responseMap) string {
	now := time.Now().UTC()

	urlStr := fmt.Sprintf("http://%s:38202/socket.io/1/?t=%s", roomInfo["server"], url.QueryEscape(now.String()))

	resp, err := http.Get(urlStr)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	r := bufio.NewReader(resp.Body)
	bytesArr, err := r.ReadBytes(':')
	if err != nil {
		log.Fatalln(err)
	}

	unknown := string(bytesArr[:len(bytesArr)-1])
	return fmt.Sprintf("ws://%s:38202/socket.io/1/websocket/%s", roomInfo["server"], unknown)
}

func main() {
	roomID := "RNWP"
	userID := uuid.New()

	roomInfo := getRoomInfo(roomID, userID)

	webSocketURL := getWebsocketURL(roomInfo)
	log.Println(webSocketURL)
}
