package common_test

import (
	"testing"
	"virtualhostingserver/common"
)

func TestReadConfig(t *testing.T) {
	expectedIP := "127.0.0.1"
	expectedPort := 8000
	expectedVirtualHost := common.VirtualHost{
		Name:    "akali",
		Path:    "templates/akali",
		Default: "404.html",
	}
	config, err := common.ReadConfig("config_test.yml")
	if err != nil {
		t.Fail()
	}

	if config.IP != expectedIP {
		t.Fatalf("TestReadConfig(): Expected \"%s\", got \"%s\"", expectedIP, config.IP)
	}

	if config.Port != expectedPort {
		t.Fatalf("TestReadConfig(): Expected \"%d\", got \"%d\"", expectedPort, config.Port)
	}

	if config.VHosts[0] != expectedVirtualHost {
		t.Fatalf("TestReadConfig(): Expected \"%v\", got \"%v\"", expectedVirtualHost, config.VHosts[0])
	}
}
