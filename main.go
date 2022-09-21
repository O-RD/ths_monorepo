package main

import (
	"fmt"

	p2p "github.com/O-RD/ths_monorepo/p2p"
)

func main() {

	p2p := p2p.P2p_init()
	fmt.Println(p2p)
}
