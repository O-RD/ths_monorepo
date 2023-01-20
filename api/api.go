package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)


func api_get_group_request(group_pub_key string) (random_message string) {

	resp, err := http.Get(os.Getenv("SERVICE_URL") + "request_to_save_group/" + group_pub_key)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println(resp)
	if resp.StatusCode == http.StatusOK {
		//bodyBytes, err := io.ReadAll(resp.Body)
		return random_message
	}
	return
}

func api_post_group_request(key *group_sign_key) {
	//func api_post_group_request(p *P2P) {

	url := os.Getenv("SERVICE_URL") +  "save_group_key"
	fmt.Println("URL:>", url)

	postBody, _ := json.Marshal(group_sign_key{
		Group_pub_key: key.Group_pub_key,
		Partial_key:   key.Partial_key,
		Id:            key.Id,
		Sign:          key.Sign,
		N:             key.N,
		T:             key.T,
	})
	responseBody := bytes.NewBuffer(postBody)
	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post(os.Getenv("SERVICE_URL") + "save_group_key", "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	fmt.Println(resp.StatusCode)
}



func get_request_to_save_random(group_pub_key string) (random_message string) {

	resp, err := http.Get(os.Getenv("SERVICE_URL") + "request_to_save_random/" + group_pub_key)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println(resp)
	if resp.StatusCode == http.StatusOK {
		//bodyBytes, err := io.ReadAll(resp.Body)
		return random_message
	}
	return
}

func post_save_random_values(key *group_random_data) {

	url := os.Getenv("SERVICE_URL") + "save_random_values"
	fmt.Println("URL:>", url)

	postBody, _ := json.Marshal(group_random_data{
		Group_pub_key: key.Group_pub_key,
		Partial_key:   key.Partial_key,
		Id:            key.Id,
		Sign:          key.Sign,
		N:             key.N,
		T:             key.T,
		Random_tag:    key.Random_tag,
		Random_val:    key.Random_val,
	})
	responseBody := bytes.NewBuffer(postBody)
	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post(os.Getenv("SERVICE_URL") + "save_random_values", "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	fmt.Println(resp.StatusCode)

}

func get_Get_signatures(group_pub_key string) (random_message string) {

	resp, err := http.Get(os.Getenv("SERVICE_URL")+ "get_signatures/" + group_pub_key)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println(resp)
	if resp.StatusCode == http.StatusOK {
		//bodyBytes, err := io.ReadAll(resp.Body)
		return random_message
	}
	return
}

func get_all_signatures(group_pub_key string) (random_message string) {

	resp, err := http.Get(os.Getenv("SERVICE_URL") + "request_to_save_random/" + group_pub_key)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println(resp)
	if resp.StatusCode == http.StatusOK {
		//bodyBytes, err := io.ReadAll(resp.Body)
		return random_message
	}
	return
}

func post_Post_signatures(key *Part_signature) {
	url := os.Getenv("SERVICE_URL") +  "post_signature"
	fmt.Println("URL:>", url)

	postBody, _ := json.Marshal(key)
	responseBody := bytes.NewBuffer(postBody)
	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post(os.Getenv("SERVICE_URL") + "post_signature", "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	fmt.Println(resp.StatusCode)
}


func post_get_signatures(key *sign_req_struct) (result []signed_msgs) {

	url := os.Getenv("SERVICE_URL") + "get_signatures"
	fmt.Println("URL:>", url)

	postBody, _ := json.Marshal(key)
	responseBody := bytes.NewBuffer(postBody)

	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post(os.Getenv("SERVICE_URL") + "get_signatures", "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	fmt.Println(resp.StatusCode)

	return result
}

