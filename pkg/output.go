package pkg

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func (scan ScanResult) WriteConsole() {
	switch scan.PortStatus {
	case Open:
		formatString := fmt.Sprintf("Port: %d | Status: Open | Service: %s | OS: %s", scan.Port, scan.Service, scan.OsType)
		color.Green(formatString)
	case Closed:
		formatString := fmt.Sprintf("Port: %d | Status: Closed | Service: %s | OS: %s", scan.Port, scan.Service, scan.OsType)
		color.Yellow(formatString)
	case Filtered:
		formatString := fmt.Sprintf("Port: %d | Status: Filtered | Service: %s | OS: %s", scan.Port, scan.Service, scan.OsType)
		color.Red(formatString)
	case OpenFiltered:
		formatString := fmt.Sprintf("Port: %d | Status: Open/Filtered | Service: %s | OS: %s", scan.Port, scan.Service, scan.OsType)
		color.Magenta(formatString)
	}
}

func (scan ScanResult) CreateFileData() []byte {
	formatString := ""
	switch scan.PortStatus {
	case Open:
		formatString = fmt.Sprintf("Port: %d | Status: Open | Service: %s | OS: %s", scan.Port, scan.Service, scan.OsType)
	case Closed:
		formatString = fmt.Sprintf("Port: %d | Status: Closed | Service: %s | OS: %s", scan.Port, scan.Service, scan.OsType)
	case Filtered:
		formatString = fmt.Sprintf("Port: %d | Status: Filtered | Service: %s | OS: %s", scan.Port, scan.Service, scan.OsType)
	case OpenFiltered:
		formatString = fmt.Sprintf("Port: %d | Status: Open/Filtered | Service: %s | OS: %s", scan.Port, scan.Service, scan.OsType)
	}

	return []byte(formatString)
}

func WriteToFile(data [][]byte, fileName string) error {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, val := range data {
		_, err := f.Write(append(val, '\n'))
		if err != nil {
			return err
		}
	}

	return nil
}
