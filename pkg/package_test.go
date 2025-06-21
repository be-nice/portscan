package pkg

import (
	"reflect"
	"testing"
)

func TestValidateIp(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid IPv4", "192.168.0.1", false},
		{"Valid IPv6", "::1", false},
		{"Invalid IP format", "300.168.1.1", true},
		{"Random string", "not_an_ip", true},
		{"Empty input", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scan := ScanConfig{}
			err := scan.validateIp(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestValidatePort(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		expected []int
	}{
		{"Valid range", "20-22", false, []int{20, 21, 22}},
		{"Invalid range order", "100-50", true, nil},
		{"Out-of-range port", "0-70000", true, nil},
		{"Valid list", "80,443,22", false, []int{80, 443, 22}},
		{"Invalid list", "80,443,22,-4", true, nil},
		{"Duplicate in list", "22,22", true, nil},
		{"Valid single", "443", false, []int{443}},
		{"Invalid single", "65536", true, nil},
		{"Random string", "abc", true, nil},
		{"Empty input", "", true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scan := ScanConfig{Port: nil}
			err := scan.validatePort(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error: %v, got: %v", tt.wantErr, err)
			}
			if err == nil && !reflect.DeepEqual(scan.Port, tt.expected) {
				t.Errorf("Expected ports: %v, got: %v", tt.expected, scan.Port)
			}
		})
	}
}

func TestValidateProtocol(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{"TCP flag", "-tcp", Tcp, false},
		{"UDP flag", "-udp", Udp, false},
		{"Stealth flag", "-s", Stealth, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scan := ScanConfig{Protocol: -1}
			err := scan.validateProtocol(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error: %v, got: %v", tt.wantErr, err)
			}
			if scan.Protocol != tt.want {
				t.Errorf("Expected protocol: %d, got: %d", tt.want, scan.Protocol)
			}
		})
	}
}

func TestValidateOutputType(t *testing.T) {
	tests := []struct {
		name      string
		flag      string
		file      string
		wantType  int
		wantFile  string
		expectErr bool
	}{
		{"All console", "-a", "", AllConsole, "", false},
		{"Open only", "-o", "", OpenConsole, "", false},
		{"Write to file", "-f", "results.txt", WriteFile, "results.txt", false},
		{"Missing filename", "-f", "", WriteFile, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scan := ScanConfig{DisplayType: -1}
			err := scan.validateOutputType(tt.flag, tt.file)
			if (err != nil) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, err)
			}
			if scan.DisplayType != tt.wantType {
				t.Errorf("Expected DisplayType: %d, got: %d", tt.wantType, scan.DisplayType)
			}
			if scan.FileName != tt.wantFile {
				t.Errorf("Expected FileName: %s, got: %s", tt.wantFile, scan.FileName)
			}
		})
	}
}

func TestSetWorkerCount(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
		wantErr  bool
	}{
		{"Valid number", "5", 5, false},
		{"Zero workers", "0", 0, false},
		{"Negative workers", "-1", -1, false},
		{"Random string", "abc", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scan := ScanConfig{MaxWorkers: -1}
			err := scan.SetWorkerCount(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error: %v, got: %v", tt.wantErr, err)
			}
			if scan.MaxWorkers != tt.expected && !tt.wantErr {
				t.Errorf("Expected MaxWorkers: %d, got: %d", tt.expected, scan.MaxWorkers)
			}
		})
	}
}

func TestScanConfigSanityCheck(t *testing.T) {
	tests := []struct {
		name      string
		input     ScanConfig
		wantErr   bool
		wantPorts int
		wantProto int
		wantDisp  int
		wantWkrs  int
	}{
		{
			name:    "Valid full config",
			input:   ScanConfig{Addr: "127.0.0.1", Port: []int{80}, Protocol: Tcp, DisplayType: OpenConsole, MaxWorkers: 5},
			wantErr: false, wantPorts: 1, wantProto: Tcp, wantDisp: OpenConsole, wantWkrs: 5,
		},
		{
			name:    "Missing Addr",
			input:   ScanConfig{Port: []int{80}, Protocol: Tcp, DisplayType: AllConsole, MaxWorkers: 10},
			wantErr: true,
		},
		{
			name:    "Nil ports defaulted",
			input:   ScanConfig{Addr: "192.168.1.1", Protocol: Tcp, DisplayType: AllConsole, MaxWorkers: 10},
			wantErr: false, wantPorts: 65536, wantProto: Tcp, wantDisp: AllConsole, wantWkrs: 10,
		},
		{
			name:    "Protocol defaulted",
			input:   ScanConfig{Addr: "192.168.1.1", Port: []int{22}, Protocol: -1, DisplayType: AllConsole, MaxWorkers: 10},
			wantErr: false, wantProto: Tcp,
		},
		{
			name:    "DisplayType defaulted",
			input:   ScanConfig{Addr: "192.168.1.1", Port: []int{22}, Protocol: Tcp, DisplayType: -1, MaxWorkers: 10},
			wantErr: false, wantDisp: AllConsole,
		},
		{
			name:    "MaxWorkers too low, default to 10",
			input:   ScanConfig{Addr: "192.168.1.1", Port: []int{22}, Protocol: Tcp, DisplayType: AllConsole, MaxWorkers: -5},
			wantErr: false, wantWkrs: 10,
		},
		{
			name:    "MaxWorkers too high, capped to 100",
			input:   ScanConfig{Addr: "192.168.1.1", Port: []int{22}, Protocol: Tcp, DisplayType: AllConsole, MaxWorkers: 200},
			wantErr: false, wantWkrs: 100,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.scanConfigSanityCheck()
			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error: %v, got: %v", tt.wantErr, err)
			}
			if err == nil {
				if tt.wantPorts != 0 && len(tt.input.Port) != tt.wantPorts {
					t.Errorf("Expected %d ports, got %d", tt.wantPorts, len(tt.input.Port))
				}
				if tt.wantProto != 0 && tt.input.Protocol != tt.wantProto {
					t.Errorf("Expected protocol %d, got %d", tt.wantProto, tt.input.Protocol)
				}
				if tt.wantDisp != 0 && tt.input.DisplayType != tt.wantDisp {
					t.Errorf("Expected display %d, got %d", tt.wantDisp, tt.input.DisplayType)
				}
				if tt.wantWkrs != 0 && tt.input.MaxWorkers != tt.wantWkrs {
					t.Errorf("Expected MaxWorkers %d, got %d", tt.wantWkrs, tt.input.MaxWorkers)
				}
			}
		})
	}
}
