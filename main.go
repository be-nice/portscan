package main

import (
	"fmt"
	"os"
	"portSec/pkg"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"

	"github.com/fatih/color"
)

var executeScanMap = map[int]func(pkg.ScanExec, chan pkg.ScanResult){
	pkg.Tcp: func(s pkg.ScanExec, res chan pkg.ScanResult) { s.TcpScan(res) },
	pkg.Udp: func(s pkg.ScanExec, res chan pkg.ScanResult) { s.UdpScan(res) },
}

var protocolMap = map[int]string{
	pkg.Tcp:     "TCP",
	pkg.Udp:     "UDP",
	pkg.Stealth: "Stealth TCP",
}

func main() {
	if len(os.Args) < 2 {
		pkg.DisplayHelp()
		return
	}

	if os.Args[1] == "-h" {
		pkg.DisplayHelp()
		return
	}

	s, err := pkg.ValidateArgs(os.Args[1:])
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		pkg.DisplayHelp()
		return
	}

	ChanWork := make(chan pkg.ScanExec, len(s.Port))
	ChanResult := make(chan pkg.ScanResult, len(s.Port))
	writeData := make([]string, 0, len(s.Port)+2)

	go func() {
		headerString := fmt.Sprintf("Scanning: %s | Protocol: %s | Num Ports: %d | Num threads: %d | DisplayType: %d", s.Addr, protocolMap[s.Protocol], len(s.Port), s.MaxWorkers, s.DisplayType)
		color.Blue(headerString)
		if s.DisplayType == pkg.WriteFile {
			writeData = append(writeData, "Scan created at: "+time.Now().String(), headerString)
		}

		for res := range ChanResult {
			switch s.DisplayType {
			case pkg.AllConsole, pkg.WriteFile:
				writeData = append(writeData, res.FormatResult())
			case pkg.OpenConsole:
				if res.PortStatus == pkg.Open || res.PortStatus == pkg.OpenFiltered {
					writeData = append(writeData, res.FormatResult())
				}
			}
		}
	}()

	bar := progressbar.Default(int64(len(s.Port)))

	var wg sync.WaitGroup
	for range s.MaxWorkers {
		wg.Add(1)
		go worker(ChanWork, ChanResult, executeScanMap, s.Protocol, &wg, bar)
	}

	go func() {
		for _, port := range s.Port {
			ChanWork <- pkg.ScanExec{
				Addr: s.Addr,
				Port: port,
			}
		}
		close(ChanWork)
	}()

	wg.Wait()
	close(ChanResult)

	if s.DisplayType == pkg.WriteFile {
		pkg.WriteToFile(writeData, s.FileName)
		fmt.Printf("Scan results written to %s\n", s.FileName)
	} else {
		for _, line := range writeData {
			fmt.Println(string(line))
		}
	}
}

func worker(
	ChanWork chan pkg.ScanExec,
	ChanResult chan pkg.ScanResult,
	executeScanMap map[int]func(pkg.ScanExec, chan pkg.ScanResult),
	protocol int,
	wg *sync.WaitGroup,
	bar *progressbar.ProgressBar,
) {
	defer wg.Done()

	for job := range ChanWork {
		bar.Add(1)
		if fn, ok := executeScanMap[protocol]; ok {
			fn(job, ChanResult)
		} else {
			fmt.Println("Unknown protocol:", protocol)
		}
	}
}

func workerWriting(res chan pkg.ScanResult, s pkg.ScanConfig, writeData [][]byte) {
	switch s.DisplayType {
	case pkg.AllConsole:
		return
	case pkg.OpenConsole:
		return
	case pkg.WriteFile:
		return
	}
}
