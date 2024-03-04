package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fabulias/metrics-api/pkg/agent"
)

func main() {
	var interval time.Duration = 20 * time.Second
	for range time.Tick(interval) {
		agent := agent.NewAgent()

		metrics, err := agent.CollectMetrics()
		if err != nil {
			fmt.Println("Error collecting metrics:", err)
			continue
		}

		jsonData, err := json.Marshal(metrics)
		if err != nil {
			fmt.Println("Error converting metrics to JSON:", err)
			continue
		}

		// hardcoded device ID because it's not possible to get from the agent at this point of development
		req, err := http.NewRequest(http.MethodPost, "http://localhost:9000/devices/1/metrics", bytes.NewReader(jsonData))
		if err != nil {
			fmt.Println("Error creating HTTP request:", err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		// The use of default client it's only for testing purposes, in production you should use your own client.
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error sending HTTP request:", err)
			continue
		}

		resp.Body.Close()
		fmt.Println("Status response", resp.Status)

	}
}
