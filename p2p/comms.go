package p2p

import (
	"encoding/json"

	ths "github.com/O-RD/ths_monorepo/ths"
)

func Send(message_channel chan ths.Message) {

	for {
		var message_data ths.Message
		message_data = <-message_channel

		send_stream, _ := message_data.From.Host.NewStream(message_data.From.Ctx, message_data.To, "ths_stream")
		message := ths.Payload{
			Sender:  message_data.From.Host.ID(),
			Type:    message_data.Type,
			Payload: message_data.Payload,
		}
		b_message, _ := json.Marshal(message)
		_, _ = send_stream.Write(append(b_message, '\n'))

		// fmt.Println(err)
		// if message_data.End == 1 {
		// 	return
		// }

	}
}
