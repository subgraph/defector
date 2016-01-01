package networkmanager

import ( 
	"testing"
)

func TestGetActiveConnections1(t *testing.T) {
	_, err := GetActiveConnections()
	if err != nil {
		t.Errorf("GetActiveConnection1 test failed: %s", err)
	}
}

func TestGetDhcpNameservers1(t *testing.T) {
	ac, err := GetActiveConnections()
	if err != nil {
		t.Errorf(
			"TestGetDhcpNameservers1 test failed on GetActiveConnections: %s", 
				err)
	}
	for _, c := range ac {
		_, err = GetDhcpNameservers(c)
		if err != nil {
			t.Errorf("TestGetDhcpNameservers1 test failed: %s",
				err)
			}
	}
}
