package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

func main() {
	t, err := strconv.Atoi(os.Getenv("PING_INTERVAL_MS"))
	if err != nil {
		log.Println("Error parsing PING_INTERVAL_MS, assuming interval of 1000 ms")
		t = 1000
	}

	for {
		query()
		time.Sleep(time.Millisecond * time.Duration(t))
	}
}

func query() {
	log.Println("Pinger working")
	log.Printf("Docker: %s\n", os.Getenv("DOCKER_SOCKET"))

	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", os.Getenv("DOCKER_SOCKET"))
			},
		},
	}
	resp, err := httpc.Get("http://unix/containers/json")
	if err != nil {
		log.Println(err)
		panic(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	decoded := make([]any, 0)
	err = json.Unmarshal(body, &decoded)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	ipList := make([]string, 0)
	for _, e := range decoded {
		outer := e.(map[string]any)
		network := outer["NetworkSettings"].(map[string]any)
		networks := network["Networks"].(map[string]any)
		for _, net := range networks {
			ip := net.(map[string]any)["IPAddress"].(string)

      pinger, err := probing.NewPinger(ip)
      if err != nil {
        continue
      }
      pinger.Timeout = time.Duration(30) * time.Millisecond
      pinger.Count = 1
      err = pinger.Run()
      if err != nil {
        continue
      }

			log.Println("IP: ", ip)
			ipList = append(ipList, ip)
		}
	}

	listToSend, err := json.Marshal(ipList)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	r := bytes.NewReader(listToSend)
	http.Post("http://conmon-backend/ping-upload", "JSON", r)
}
