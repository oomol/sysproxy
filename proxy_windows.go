package systeminfo

import (
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

type WinHttpProxyInfo struct {
	Host   string
	Port   uint32
	Bypass string
}

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

func GetProxyInfo() (*WinHttpProxyInfo, error) {
	config, err := getIEProxyConfig()
	if err != nil {
		return nil, err
	}

	proxyUrl := syscall.UTF16ToString((*[1 << 16]uint16)(unsafe.Pointer(config.AutoConfigURL))[:])
	proxyBaypass := syscall.UTF16ToString((*[1 << 16]uint16)(unsafe.Pointer(config.ProxyBypass))[:])
	proxyHost := strings.Split(proxyUrl, ":")[0]
	proxyPort, err := strconv.ParseUint(strings.Split(proxyUrl, ":")[1], 10, 32)
	if err != nil {
		return nil, err
	}

	return &WinHttpProxyInfo{
		Host:   proxyHost,
		Port:   uint32(proxyPort),
		Bypass: proxyBaypass,
	}, nil
}
