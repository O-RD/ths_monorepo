package main

import (
	"flag"
)

func Process_flags() (string, string, int, int, string) {

	var port = *flag.String("p", "41803", "channel id for connecting")
	var version = *flag.String("v", "t1", "Type of run") //
	var threshold = *flag.Int("t", 2, "Threshold")
	var party_size = *flag.Int("n", 3, "N value")
	var moniker = *flag.String("m", "Default_Name", "Moniker")
	return port, version, threshold, party_size, moniker
}
