package main

import (
	"flag"

	"github.com/O-RD/ths_monorepo/ths"
)

func Process_flags() (string, string, int, int, string) {

	// strconv.Itoa(rand.Intn(99999-20000) + 20000)
	var topic = *flag.String("topic", "john2020", "channel id for connecting")
	var version = *flag.String("v", "t1", "Type of run") //
	var threshold = *flag.Int("t", 2, "Threshold")
	var party_size = *flag.Int("n", 2, "N value")
	var moniker = *flag.String("m", "Default_Name", "Moniker")
	flag.Parse()

	return topic, version, threshold, party_size, moniker
}
func Sort_Peers(p *ths.P2P) {

}
