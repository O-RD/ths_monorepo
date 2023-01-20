package api



type group_sign_key struct {
	Group_pub_key string `json:"group_pub_key"`
	Partial_key   string `json:"partial_key"`
	Id            string `json:"id"`
	Sign          string `json:"sign"`
	N             int    `json:"n"`
	T             int    `json:"t"`
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

type Part_signature struct {
	Key_group_name          string `json:"key_group"`
	Partial_group_pub_key   string `json:"partial_group_pub_key"`
	Peer_Id                 string `json:"peer_id"`
	Message                 string `json:"message"`
	Part_signature          string `json:"part_signature"`
	Randomness_tag          string `json:"randomness_tag"`
	Nonce                   int    `json:"nonce"`
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