package systeminfo

import (
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

type winHttpCurrentUserIeProxyConfig struct {
	AutoDetect    int32
	AutoConfigURL *uint16
	Proxy         *uint16
	ProxyBypass   *uint16
}

func getIEProxyConfig() (*winHttpCurrentUserIeProxyConfig, error) {
	winhttp := syscall.NewLazyDLL("winhttp.dll")
	procGetProxy := winhttp.NewProc("WinHttpGetIEProxyConfigForCurrentUser")

	var config winHttpCurrentUserIeProxyConfig

	r1, _, err := procGetProxy.Call(uintptr(unsafe.Pointer(&config)))
	if r1 == 0 {
		return nil, err
	}

	return &config, nil
}

func GetProxyInfo() (*HttpProxyInfo, *HttpsProxyInfo, error) {
	config, err := getIEProxyConfig()
	if err != nil {
		return nil, nil, err
	}

	proxyUrl := syscall.UTF16ToString((*[1 << 16]uint16)(unsafe.Pointer(config.Proxy))[:])
	proxyBypass := syscall.UTF16ToString((*[1 << 16]uint16)(unsafe.Pointer(config.ProxyBypass))[:])
	proxyHost := strings.Split(proxyUrl, ":")[0]
	proxyPort, err := strconv.ParseUint(strings.Split(proxyUrl, ":")[1], 10, 32)
	if err != nil {
		return nil, nil, err
	}

	httpProxyInfo := &HttpProxyInfo{
		ProxyInfo: ProxyInfo{
			Host:   proxyHost,
			Port:   uint16(proxyPort),
			Bypass: proxyBypass,
		},
	}

	httpsProxyInfo := &HttpsProxyInfo{
		ProxyInfo: ProxyInfo{
			Host:   proxyHost,
			Port:   uint16(proxyPort),
			Bypass: proxyBypass,
		},
	}

	return httpProxyInfo, httpsProxyInfo, nil
}
