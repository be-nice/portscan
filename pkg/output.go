package pkg

import (
	"fmt"
	"os"
	"regexp"

	"github.com/fatih/color"
)

func (scan ScanResult) FormatResult() string {
	formatString := ""
	switch scan.PortStatus {
	case Open:
		formatString = color.RedString(fmt.Sprintf("Port: %d | Status: Open | Service: %s | OS: %s", scan.Port, scan.Service, scan.OsType))
	case Closed:
		formatString = color.GreenString(fmt.Sprintf("Port: %d | Status: Closed | Service: %s | OS: %s", scan.Port, scan.Service, scan.OsType))
	case Filtered:
		formatString = color.YellowString(fmt.Sprintf("Port: %d | Status: Filtered | Service: %s | OS: %s", scan.Port, scan.Service, scan.OsType))
	case OpenFiltered:
		formatString = color.YellowString(fmt.Sprintf("Port: %d | Status: Open/Filtered | Service: %s | OS: %s", scan.Port, scan.Service, scan.OsType))
	}

	return formatString
}

func WriteToFile(data []string, fileName string) error {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	ansi := regexp.MustCompile(`\x1b\[[0-9;]*m`)

	for _, val := range data {
		val = ansi.ReplaceAllString(val, "")
		_, err := f.Write(append([]byte(val), '\n'))
		if err != nil {
			return err
		}
	}

	return nil
}
