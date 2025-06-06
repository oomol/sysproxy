package sysproxy

import "testing"

func TestGetProxyInfo(t *testing.T) {
	httpInfo, httpsInfo, err := GetAll()
	if err != nil {
		t.Fatalf("%v", err)
	}

	if httpInfo != nil {
		t.Logf("HTTP Proxy Host: %v", httpInfo.Host)
		t.Logf("HTTP Proxy Port: %v", httpInfo.Port)
	}

	if httpsInfo != nil {
		t.Logf("HTTPS Proxy Host: %v", httpsInfo.Host)
		t.Logf("HTTPS Proxy Port: %v", httpsInfo.Port)
	}

}
