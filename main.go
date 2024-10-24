package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const url = "https://api.trongrid.io/v1/contracts/TTfvyrAz86hbZk5iDpKD78pqLGgi8C7AAw/events"

type Event struct {
	EventName      string `json:"event_name"`
	BlockTimestamp int    `json:"block_timestamp"`
	TransactionId  string `json:"transaction_id"`
}

type Response struct {
	Data []Event `json:"data"`
}

var processedEvents = make(map[string]bool)

func checkEvent() {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Parsing error JSON:", err)
		return
	}

	for _, event := range response.Data {
		if event.EventName == "TokenCreate" {
			if _, exist := processedEvents[event.TransactionId]; !exist {
				processedEvents[event.TransactionId] = true
				fmt.Println("Event TokenCreate", event.TransactionId)
			}
		}
	}
}

func main() {
	for {
		checkEvent()
		time.Sleep(1 * time.Second)
	}
}
