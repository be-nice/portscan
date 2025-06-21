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

	ChanWork := make(chan pkg.ScanExec)
	ChanResult := make(chan pkg.ScanResult)

	var wg sync.WaitGroup
	var wgPrint sync.WaitGroup
	for range s.MaxWorkers {
		wg.Add(1)
		go worker(ChanWork, ChanResult, executeScanMap, s.Protocol, &wg)
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

	wgPrint.Add(1)
	go func() {
		result := make([]pkg.ScanResult, 0, len(s.Port))
		var filedata [][]byte
		headerString := fmt.Sprintf("Scanning: %s | Protocol: %s | Num Ports: %d | Num threads: %d", s.Addr, protocolMap[s.Protocol], len(s.Port), s.MaxWorkers)
		color.Blue(headerString)

		if s.DisplayType == pkg.WriteFile {
			filedata = make([][]byte, 0, len(s.Port)+1)
			filedata = append(filedata, []byte("Scan created at: "+time.Now().String()), []byte(headerString))
		}

		bar := progressbar.Default(int64(len(s.Port)))
		defer wgPrint.Done()

		for res := range ChanResult {
			bar.Add(1)
			result = append(result, res)
		}

		for _, res := range result {
			switch s.DisplayType {
			case pkg.AllConsole:
				res.WriteConsole()
			case pkg.OpenConsole:
				if res.PortStatus == pkg.Open || res.PortStatus == pkg.OpenFiltered {
					res.WriteConsole()
				}
			case pkg.WriteFile:
				filedata = append(filedata, res.CreateFileData())
			}
		}

		if s.DisplayType == pkg.WriteFile {
			err := pkg.WriteToFile(filedata, s.FileName)
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
			fmt.Printf("Scan results written to %s\n", s.FileName)
		}
	}()

	wg.Wait()
	close(ChanResult)
	wgPrint.Wait()
}

func worker(
	ChanWork chan pkg.ScanExec,
	ChanResult chan pkg.ScanResult,
	executeScanMap map[int]func(pkg.ScanExec, chan pkg.ScanResult),
	protocol int,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for job := range ChanWork {
		if fn, ok := executeScanMap[protocol]; ok {
			fn(job, ChanResult)
		} else {
			fmt.Println("Unknown protocol:", protocol)
		}
	}
}
