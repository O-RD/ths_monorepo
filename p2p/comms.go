package p2p

import (
	"fmt"

	ths "github.com/O-RD/ths_monorepo/ths"
)

func Send(message_channel chan ths.Message) {

	for {
		var values ths.Message
		values = <-message_channel

		fmt.Println("Woohooo", values)
	}
}
