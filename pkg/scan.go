package pkg

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func (scan ScanExec) TcpScan(ch chan ScanResult) {
	address := fmt.Sprintf("%s:%d", scan.Addr, scan.Port)
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)

	scanResult := ScanResult{
		Port: scan.Port,
	}

	if service, ok := serviceMap[scan.Port]; ok {
		scanResult.Service = service
	} else {
		scanResult.Service = "unknown"
	}

	if err != nil {
		scanResult.PortStatus = Closed
		ch <- scanResult
		return
	} else {
		scanResult.PortStatus = Open
	}
	defer conn.Close()

	ch <- scanResult
}

func (scan ScanExec) UdpScan(ch chan ScanResult) {
	addr := fmt.Sprintf("%s:%d", scan.Addr, scan.Port)
	udpAddr, _ := net.ResolveUDPAddr("udp", addr)

	scanResult := ScanResult{
		Port: scan.Port,
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		scanResult.PortStatus = Closed
		ch <- scanResult
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte{0x0})
	if err != nil {
		scanResult.PortStatus = Closed
		ch <- scanResult
		return
	}

	err = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		scanResult.PortStatus = Closed
		ch <- scanResult
		return
	}

	buf := make([]byte, 1024)
	_, _, err = conn.ReadFromUDP(buf)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			scanResult.PortStatus = OpenFiltered
			ch <- scanResult
			return

		} else if strings.Contains(err.Error(), "connection refused") {
			scanResult.PortStatus = Closed
			ch <- scanResult
			return

		} else {
			scanResult.PortStatus = Closed
			ch <- scanResult
			return
		}
	}
	scanResult.PortStatus = Open
	ch <- scanResult
	return
}
