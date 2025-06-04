package sysproxy

import (
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

var (
	procGetProxy *syscall.LazyProc
)

func init() {
	procGetProxy = syscall.NewLazyDLL("winhttp.dll").NewProc("WinHttpGetIEProxyConfigForCurrentUser")
}

type winHttpCurrentUserIeProxyConfig struct {
	AutoDetect    int32
	AutoConfigURL *uint16
	Proxy         *uint16
	ProxyBypass   *uint16
}

func getIEProxyConfig() (*winHttpCurrentUserIeProxyConfig, error) {
	var config winHttpCurrentUserIeProxyConfig

	r1, _, err := procGetProxy.Call(uintptr(unsafe.Pointer(&config)))
	if r1 == 0 {
		return nil, err
	}

	return &config, nil
}

// GetInfo Get Windows proxy information. Windows proxy settings only support http proxy.
func GetInfo() (*ProxyInfo, *ProxyInfo, error) {
	config, err := getIEProxyConfig()
	if err != nil {
		return nil, nil, err
	}

	proxyUrl := syscall.UTF16ToString((*[1 << 16]uint16)(unsafe.Pointer(config.Proxy))[:])
	proxyHost := strings.Split(proxyUrl, ":")[0]
	proxyPort, err := strconv.ParseUint(strings.Split(proxyUrl, ":")[1], 10, 32)
	if err != nil {
		return nil, nil, err
	}

	info := &ProxyInfo{
		Host: proxyHost,
		Port: uint16(proxyPort),
	}

	return info, nil, nil
}
