package ths

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func api_get(send_topic string) []postgres {

	resp, err := http.Get("http://43.205.198.157:3000/get_peers/" + send_topic)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println(resp)
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		data := []postgres{}
		_ = json.Unmarshal([]byte(bodyBytes), &data)

		if err != nil {
			log.Fatal(err)
		}
		// bodyString := string(bodyBytes)
		return data
	}
	return nil
}

type postgres struct {
	Item string   `json:"topic"`
	Ip   []string `json:"ip"`
	Peer string   `json:"peer"`
}

func api_send(p *P2P) {
	url := "http://43.205.198.157:3000/advertise/"
	fmt.Println("URL:>", url)

	postBody, _ := json.Marshal(postgres{
		Item: p.Topic,
		Ip:   p.Host_ips,
		Peer: p.Host_id,
	})
	responseBody := bytes.NewBuffer(postBody)
	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post("http://43.205.198.157:3000/advertise", "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	fmt.Println(resp.StatusCode)
}
