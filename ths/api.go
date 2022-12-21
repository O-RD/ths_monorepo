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

type group_sign_key struct {
	Group_pub_key string `json:"group_pub_key"`
	Partial_key   string `json:"partial_key"`
	Id            string `json:"id"`
	Sign          string `json:"sign"`
	N             int    `json:"n"`
	T             int    `json:"t"`
}


func api_get_group_request(group_pub_key string) (random_message string){
    
	resp, err := http.Get("localhost:3000/request_to_save_group/" + group_pub_key)
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
    
	url := "localhost:3000/save_group_key"
	fmt.Println("URL:>", url)

	postBody, _ := json.Marshal(group_sign_key{ 
        Group_pub_key : key.Group_pub_key , 
	    Partial_key   : key.Partial_key,        
	    Id            : key.Id  ,
	    Sign          : key.Sign ,
	    N             : key.N  ,
	    T             : key.T ,
		// ?
	})
	responseBody := bytes.NewBuffer(postBody)
	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post("localhost:3000/save_group_key", "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	fmt.Println(resp.StatusCode)
} 

type group_random_data struct {
	Group_pub_key string   `json:"group_pub_key"`
	Partial_key   string   `json:"partial_key"`
	Id            string   `json:"id"`
	Sign          string   `json:"sign"`
	N             int      `json:"n"`
	T             int      `json:"t"`
	Random_tag    []string `json:"random_tag"`
	Random_val    []string `json:"random_val"`
} 

func get_request_to_save_random(group_pub_key string) (random_message string){
    
	resp, err := http.Get("localhost:3000/request_to_save_random/" + group_pub_key)
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

func post_save_random_values( key *group_random_data){ 
	
	url := "localhost:3000/save_random_values"
	fmt.Println("URL:>", url)

	postBody, _ := json.Marshal(group_random_data{ 
        Group_pub_key : key.Group_pub_key , 
	    Partial_key   : key.Partial_key,        
	    Id            : key.Id  ,
	    Sign          : key.Sign ,
	    N             : key.N  ,
	    T             : key.T ,
		Random_tag    : key.Random_tag ,
		Random_val    : key.Random_val,
	
	})
	responseBody := bytes.NewBuffer(postBody)
	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post("localhost:3000/save_random_values", "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	fmt.Println(resp.StatusCode)
	
} 



func get_Get_signatures(group_pub_key string) (random_message string){
    
	resp, err := http.Get("localhost:3000/get_signatures/" + group_pub_key)
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

func get_all_signatures(group_pub_key string) (random_message string){
    
	resp, err := http.Get("localhost:3000/request_to_save_random/" + group_pub_key)
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
/*
/sign_message POST
Input:
{
 “group_pub key”:    type String,
"Partial_key" : type String
 “id”:  peer_id ,    type strings
"txn_message":string
 "sign": string
"random_tag" string.
"nonce":int

}
Output: 200

*/
func post_Post_signatures(p *P2P) {
	url := "localhost:3000/post_signature"
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

type sign_req_struct struct {
	Group_pub_key string `json:"group_pub_key"`
	Partial_key   string `json:"partial_key"`
	Sign          string `json:"sign"`
}  

type signed_msgs struct {
	Txn_id string `json:"txn_id"`
	Msg    string `json:"msg"`
	Nonce  int    `json:"nonce"`
}

func post_get_signatures( key *sign_req_struct) ( result []signed_msgs) { 
	
	url := "localhost:3000/post_signature"
	fmt.Println("URL:>", url)

	postBody, _ := json.Marshal(key)
	responseBody := bytes.NewBuffer(postBody) 
    
	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post("localhost:3000/post_signature", "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	fmt.Println(resp.StatusCode)

	return result
} 
