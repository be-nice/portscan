package pkg

// Protocol type enums, Tcp is default if no protocol selected
const (
	Tcp int = iota
	Udp
	Stealth
)

// Port status enums, no default- always set by response
const (
	Open int = iota
	Filtered
	Closed
	OpenFiltered
)

// Output type enums, AllConsole is default if no type selected
const (
	AllConsole int = iota
	OpenConsole
	WriteFile
)

type ScanConfig struct {
	Addr        string
	Port        []int
	Protocol    int
	DisplayType int
	MaxWorkers  int
	FileName    string
}

type ScanExec struct {
	Addr string
	Port int
}

type ScanResult struct {
	Port       int
	PortStatus int
	Service    string
}
