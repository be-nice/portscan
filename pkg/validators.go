package pkg

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

func ValidateArgs(args []string) (ScanConfig, error) {
	var err error
	scan := ScanConfig{
		Addr:        "",
		Port:        nil,
		Protocol:    -1,
		DisplayType: -1,
		FileName:    "",
		MaxWorkers:  -1,
	}

	err = scan.validateIp(args[0])
	if err != nil {
		return scan, err
	}

	for i := 1; i < len(args); i++ {
		switch strings.ToLower(args[i]) {
		case "-p":
			if i+1 > len(args)-1 {
				return scan, errors.New("Expected value after -p")
			}
			err = scan.validatePort(args[i+1])
			if err != nil {
				return scan, err
			}
			i++
		case "-tcp", "-udp", "-s":
			err := scan.validateProtocol(strings.TrimSpace(strings.ToLower(args[i])))
			if err != nil {
				return scan, err
			}
		case "-o", "-a":
			err := scan.validateOutputType(strings.TrimSpace(strings.ToLower(args[i])))
			if err != nil {
				return scan, err
			}
		case "-f":
			if i+1 > len(args)-1 {
				return scan, errors.New("Expected filename after -f")
			}
			err := scan.validateOutputType(args[i], strings.TrimSpace(strings.ToLower(args[i+1])))
			if err != nil {
				return scan, err
			}
			i++
		case "-w":
			if i+1 > len(args)-1 {
				return scan, errors.New("Expected value after -w")
			}
			err = scan.SetWorkerCount(args[i+1])
			if err != nil {
				return scan, err
			}
			i++
		default:
			return scan, errors.New("Unknown flag")
		}
	}

	err = scan.scanConfigSanityCheck()
	if err != nil {
		return scan, err
	}

	return scan, nil
}

func (scan *ScanConfig) scanConfigSanityCheck() error {
	if scan.Addr == "" {
		return errors.New("internal error, sanity check Addr")
	}

	if scan.Port == nil {
		err := scan.validatePort("0-65535")
		if err != nil {
			return errors.New("internal error, sanity check Port")
		}
	}

	if scan.Protocol == -1 {
		scan.Protocol = Tcp
	}

	if scan.DisplayType == -1 {
		scan.DisplayType = AllConsole
	}

	if scan.MaxWorkers <= 0 {
		scan.MaxWorkers = 10
	}

	if scan.MaxWorkers >= 101 {
		scan.MaxWorkers = 100
	}

	return nil
}

func (scan *ScanConfig) validateIp(ip string) error {
	validIP := func(s string) bool {
		return net.ParseIP(s) != nil
	}(ip)

	if !validIP {
		return errors.New("invalid IP address")
	}

	scan.Addr = ip
	return nil
}

func (scan *ScanConfig) validatePort(ports string) error {
	if scan.Port != nil {
		return errors.New("Invalid syntax, port/s allready selected")
	}

	validPortRange := func(n int) bool {
		return n >= 0 && n <= 65535
	}

	makePortRange := func(s *ScanConfig, parts []string) error {
		startPort, err := strconv.Atoi(parts[0])
		if err != nil || !validPortRange(startPort) {
			return errors.New("invalid port value")
		}
		endPort, err := strconv.Atoi(parts[1])
		if err != nil || !validPortRange(endPort) {
			return errors.New("invalid port value")
		}

		if startPort > endPort {
			return errors.New("Invalid syntax, start port is larger than end port")
		}

		portList := make([]int, 0, endPort-startPort)

		for ; startPort <= endPort; startPort++ {
			portList = append(portList, startPort)
		}

		s.Port = portList
		return nil
	}

	makePortList := func(s *ScanConfig, parts []string) error {
		dupeMap := make(map[int]struct{})
		portList := make([]int, 0, len(parts))

		for _, val := range parts {
			port, err := strconv.Atoi(val)
			if err != nil || !validPortRange(port) {
				return errors.New("invalid port value")
			}

			if _, ok := dupeMap[port]; ok {
				return errors.New("duplicate port in the list")
			}
			dupeMap[port] = struct{}{}
			portList = append(portList, port)
		}

		if len(portList) == 0 {
			return errors.New("Invalid syntax, no ports provided")
		}

		s.Port = portList
		return nil
	}

	if strings.Contains(ports, "-") {
		parts := strings.Split(ports, "-")
		if len(parts) != 2 {
			return errors.New("Invalid syntax for port range")
		}

		return makePortRange(scan, parts)
	} else {
		parts := strings.Split(ports, ",")
		if len(parts) == 0 {
			return errors.New("invalid syntax for ports")
		}

		return makePortList(scan, parts)
	}
}

func (scan *ScanConfig) validateProtocol(prot string) error {
	if scan.Protocol != -1 {
		return errors.New("invalid syntax, multiple protocol flags")
	}

	switch prot {
	case "-tcp":
		scan.Protocol = Tcp
	case "-udp":
		scan.Protocol = Udp
	case "-s":
		scan.Protocol = Stealth
	}

	return nil
}

func (scan *ScanConfig) validateOutputType(flag string, filename ...string) error {
	if scan.DisplayType != -1 {
		return errors.New("invalid syntax, multiple output flags")
	}

	switch flag {
	case "-o":
		scan.DisplayType = OpenConsole
	case "-a":
		scan.DisplayType = AllConsole
	case "-f":
		scan.DisplayType = WriteFile
		if len(filename) == 0 || filename[0] == "" {
			return errors.New("invalid syntax, missing output file name")
		}
		scan.FileName = filename[0]
	}

	return nil
}

func (scan *ScanConfig) SetWorkerCount(worker string) error {
	if scan.MaxWorkers != -1 {
		return errors.New("invalid syntax, multiple worker flags")
	}

	n, err := strconv.Atoi(worker)
	if err != nil {
		return errors.New("invalid worker count value, expected int")
	}

	scan.MaxWorkers = n
	return nil
}
