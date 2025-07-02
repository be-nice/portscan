package main

import (
	"fmt"
	"io"
	"os"
	"portSec/pkg"
	"sync"
	"sync/atomic"
	"time"

	"fortio.org/progressbar"
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
	bar := progressbar.NewBar()
	writer := bar.Writer()
	var counter int32 = 0

	var wgWrite sync.WaitGroup
	wgWrite.Add(1)
	go func(w io.Writer) {
		defer wgWrite.Done()
		headerString := fmt.Sprintf("Scanning: %s | Protocol: %s | Num Ports: %d | Num threads: %d | DisplayType: %d", s.Addr, protocolMap[s.Protocol], len(s.Port), s.MaxWorkers, s.DisplayType)
		color.Blue(headerString)
		if s.DisplayType == pkg.WriteFile {
			writeData = append(writeData, "Scan created at: "+time.Now().String(), headerString)
		}

		for res := range ChanResult {
			switch s.DisplayType {
			case pkg.WriteFile:
				writeData = append(writeData, res.FormatResult())
			case pkg.OpenConsole:
				if res.PortStatus == pkg.Open || res.PortStatus == pkg.OpenFiltered {
					fmt.Fprintln(w, res.FormatResult())
				}
			default:
				fmt.Fprintln(w, res.FormatResult())
			}
		}
	}(writer)

	var wg sync.WaitGroup
	for range s.MaxWorkers {
		wg.Add(1)
		go worker(ChanWork, ChanResult, executeScanMap, s.Protocol, len(s.Port), &counter, &wg, bar)
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
	bar.End()
	close(ChanResult)
	wgWrite.Wait()

	if s.DisplayType == pkg.WriteFile {
		pkg.WriteToFile(writeData, s.FileName)
	}
}

func worker(
	ChanWork chan pkg.ScanExec,
	ChanResult chan pkg.ScanResult,
	executeScanMap map[int]func(pkg.ScanExec, chan pkg.ScanResult),
	protocol int,
	elemCount int,
	counter *int32,
	wg *sync.WaitGroup,
	bar *progressbar.Bar,
) {
	defer wg.Done()

	for job := range ChanWork {
		atomic.AddInt32(counter, 1)
		bar.Progress(100. * float64(atomic.LoadInt32(counter)) / float64(elemCount))
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
