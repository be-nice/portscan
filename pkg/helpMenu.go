package pkg

import "fmt"

func DisplayHelp() {
	fmt.Printf(`
||  Usage: <ip> <optional params>
||  Optional Flags:
|| <-w> number of threads used for parallel scanning
|| <-p argument> specify port/s or ranges
||   <-p n> single port
||  <-p n,n,n> port list
||  <-p n-n> port range 
|| <-tcp> <-udp> <-s> sets scanning protocol
|| <-f filename> saves scan results to a file
|| <-a> <-o> display all scan result | display open only results
`)
}
