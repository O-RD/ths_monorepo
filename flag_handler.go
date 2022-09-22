package main

import (
	"flag"
)

func process_flags() (string, string, int) {

	var port = *flag.String("p", "41803", "channel id for connecting")
	var version = *flag.String("t", "t1", "Type of run")
	var threshold = *flag.Int("thresh", 2, "Type of run")

	return port, version, threshold
}
