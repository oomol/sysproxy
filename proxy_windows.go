package sysproxy

import (
	"fmt"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	procGetProxy *syscall.LazyProc
)

func init() {
	procGetProxy = syscall.NewLazyDLL("winhttp.dll").NewProc("WinHttpGetIEProxyConfigForCurrentUser")
}

type rawProxyConfig struct {
	autoDetect    bool
	autoConfigUrl *uint16
	proxy         *uint16
	proxyBypass   *uint16
}

func GetHttpProxy() (*Info, error) {
	var rawConfig rawProxyConfig
	r1, _, err := procGetProxy.Call(uintptr(unsafe.Pointer(&rawConfig)))
	if r1 == 0 {
		return nil, fmt.Errorf("WinHttpGetIEProxyConfigForCurrentUser error: %v", err)
	}
	proxyURL := convertUTF16Ptr(rawConfig.proxy)

	if proxyURL == "" {
		return nil, nil
	}

	host := strings.Split(proxyURL, ":")[0]
	port, err := strconv.ParseUint(strings.Split(proxyURL, ":")[1], 10, 32)
	if err != nil {
		return nil, err
	}

	info := &Info{
		Host: host,
		Port: uint16(port),
	}

	return info, nil
}

func GetHttpsProxy() (*Info, error) {
	return nil, nil
}

// GetAll Get Windows proxy information. Windows proxy settings only support http proxy.
func GetAll() (*Info, *Info, error) {
	httpProxyInfo, err := GetHttpProxy()
	if err != nil {
		return nil, nil, err
	}

	httpsProxyInfo, err := GetHttpsProxy()
	if err != nil {
		return nil, nil, err
	}

	return httpProxyInfo, httpsProxyInfo, nil
}

// convertUTF16Ptr safely converts a pointer to a UTF16 string.
func convertUTF16Ptr(ptr *uint16) string {
	if ptr == nil {
		return ""
	}
	return windows.UTF16PtrToString(ptr)
}
