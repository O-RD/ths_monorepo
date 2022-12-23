package main

import (
	"flag"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/O-RD/ths_monorepo/ths"
)

func temp_topic() string {
	f, _ := ioutil.ReadFile("extras.txt")
	f2 := string(f)
	fi, _ := strconv.Atoi(f2)
	var fs string
	if fi%2 == 0 {
		fs = strconv.Itoa(fi - 1)
	} else {
		fs = strconv.Itoa(fi)
	}
	fi++

	f3, _ := os.Create("extras.txt")
	f3.WriteString(strconv.Itoa(fi))
	return "john2" + fs
}

func Process_flags() (string, string, int, int, string) {

	// strconv.Itoa(rand.Intn(99999-20000) + 20000)
	name := temp_topic()
	// fmt.Println("NAME:", name)
	var topic = *flag.String("topic", name, "channel id for connecting")
	var version = *flag.String("v", "t1", "Type of run") //
	var threshold = *flag.Int("t", 2, "Threshold")
	var party_size = *flag.Int("n", 2, "N value")
	var moniker = *flag.String("m", "Default_Name", "Moniker")
	flag.Parse()

	return topic, version, threshold, party_size, moniker
}
func Sort_Peers(p *ths.P2P) {

}
