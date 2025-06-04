package systeminfo

import "testing"

func TestGetProxyInfo(t *testing.T) {
	h, hs, err := GetProxyInfo()
	if err != nil {
		t.Fatalf("%v", err)
	}

	t.Logf("httpProxyInfo: %v", h)
	t.Logf("httpsProxyInfo: %v", hs)

}
