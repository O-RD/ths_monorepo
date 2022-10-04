package p2p

import (
	"encoding/json"
	"fmt"
	"time"

	ths "github.com/O-RD/ths_monorepo/ths"
)

func Send(message_channel chan ths.Message) {

	for {
		var message_data ths.Message
		message_data = <-message_channel
		send_stream, err := message_data.From.Host.NewStream(message_data.From.Ctx, message_data.To, "ths_stream_keygen")
		if err != nil {
			fmt.Println(err)
		}
		message := ths.Payload{
			Sender:       message_data.From.Host.ID(),
			Type:         message_data.Type,
			Payload_name: message_data.Payload_name,
			Payload:      message_data.Payload,
		}
		b_message, _ := json.Marshal(message)
	Inner:
		for {
			_, err := send_stream.Write(append(b_message, '\n'))

			if err != nil {
				fmt.Println(err)
				time.Sleep(time.Millisecond)
			} else {
				break Inner
			}
		}
		// fmt.Println(err)
		// if message_data.End == 1 {
		// 	return
		// }

	}
}
