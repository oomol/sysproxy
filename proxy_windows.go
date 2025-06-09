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
	// Ref: https://learn.microsoft.com/en-us/windows/win32/api/winhttp/nf-winhttp-winhttpgetieproxyconfigforcurrentuser
	procGetProxy = syscall.NewLazyDLL("winhttp.dll").NewProc("WinHttpGetIEProxyConfigForCurrentUser")
}

// Ref: https://learn.microsoft.com/en-us/windows/win32/api/winhttp/ns-winhttp-winhttp_current_user_ie_proxy_config
type rawProxyConfig struct {
	autoDetect    bool
	autoConfigUrl *uint16
	proxy         *uint16
	proxyBypass   *uint16
}

func GetHTTP() (*Info, error) {
	var c rawProxyConfig
	r1, _, err := procGetProxy.Call(uintptr(unsafe.Pointer(&c)))
	if r1 == 0 {
		return nil, fmt.Errorf("cannot get IE proxy config: %w", err)
	}

	proxyURL := windows.UTF16PtrToString(c.proxy)
	if proxyURL == "" {
		return nil, nil
	}

	part := strings.SplitN(proxyURL, ":", 2)
	if len(part) != 2 {
		return nil, fmt.Errorf("invalid proxy URL format: %s", proxyURL)
	}

	host := part[0]
	port, err := strconv.ParseUint(part[1], 10, 32)
	if err != nil {
		return nil, err
	}

	return &Info{
		Host: host,
		Port: uint16(port),
	}, nil
}

func GetHTTPS() (*Info, error) {
	return nil, nil
}

// GetAll Get Windows proxy information. Windows proxy settings only support http proxy.
func GetAll() (*Info, *Info, error) {
	http, err := GetHTTP()
	return http, nil, err
}
