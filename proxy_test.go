package sysproxy

import "testing"

func TestGetProxyInfo(t *testing.T) {
	httpInfo, httpsInfo, err := GetInfo()
	if err != nil {
		t.Fatalf("%v", err)
	}

	if httpInfo.Host != "" {
		t.Logf("HTTP Proxy Host: %v", httpInfo.Host)
		t.Logf("HTTP Proxy Port: %v", httpInfo.Port)
	}

	if httpInfo.Host != "" {
		t.Logf("HTTPS Proxy Host: %v", httpsInfo.Host)
		t.Logf("HTTPS Proxy Port: %v", httpsInfo.Port)
	}

}
