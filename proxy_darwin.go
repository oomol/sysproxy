package sysproxy

/*
#cgo CFLAGS: -mmacosx-version-min=10.10
#cgo LDFLAGS: -framework CoreFoundation -framework SystemConfiguration

#include <CoreFoundation/CoreFoundation.h>
#include <SystemConfiguration/SystemConfiguration.h>
#include <stdlib.h>

typedef struct {
    int enabled;
    char host[256];
    int port;
} ProxyInfo;

ProxyInfo getHTTP(CFDictionaryRef settings) {
    ProxyInfo info = {0};
    if (!settings) return info;

    CFNumberRef enabledVal = CFDictionaryGetValue(settings, kSCPropNetProxiesHTTPEnable);
    if (enabledVal && CFGetTypeID(enabledVal) == CFNumberGetTypeID()) {
        CFNumberGetValue(enabledVal, kCFNumberIntType, &info.enabled);
    }

    if (info.enabled) {
        CFStringRef hostVal = CFDictionaryGetValue(settings, kSCPropNetProxiesHTTPProxy);
        CFNumberRef portVal = CFDictionaryGetValue(settings, kSCPropNetProxiesHTTPPort);

        if (hostVal && CFGetTypeID(hostVal) == CFStringGetTypeID()) {
            CFStringGetCString(hostVal, info.host, sizeof(info.host), kCFStringEncodingUTF8);
        }
        if (portVal && CFGetTypeID(portVal) == CFNumberGetTypeID()) {
            CFNumberGetValue(portVal, kCFNumberIntType, &info.port);
        }
    }

    return info;
}

ProxyInfo getHTTPS(CFDictionaryRef settings) {
    ProxyInfo info = {0};
    if (!settings) return info;

    CFNumberRef enabledVal = CFDictionaryGetValue(settings, kSCPropNetProxiesHTTPSEnable);
    if (enabledVal && CFGetTypeID(enabledVal) == CFNumberGetTypeID()) {
        CFNumberGetValue(enabledVal, kCFNumberIntType, &info.enabled);
    }

    if (info.enabled) {
        CFStringRef hostVal = CFDictionaryGetValue(settings, kSCPropNetProxiesHTTPSProxy);
        CFNumberRef portVal = CFDictionaryGetValue(settings, kSCPropNetProxiesHTTPSPort);

        if (hostVal && CFGetTypeID(hostVal) == CFStringGetTypeID()) {
            CFStringGetCString(hostVal, info.host, sizeof(info.host), kCFStringEncodingUTF8);
        }
        if (portVal && CFGetTypeID(portVal) == CFNumberGetTypeID()) {
            CFNumberGetValue(portVal, kCFNumberIntType, &info.port);
        }
    }

    return info;
}
*/
import "C"
import (
	"fmt"
	"io"
	"unsafe"
)

type setting struct {
	ref C.CFDictionaryRef

	io.Closer
}

func (s *setting) Close() error {
	if unsafe.Pointer(s.ref) != nil {
		C.CFRelease(C.CFTypeRef(s.ref))
	}

	return nil
}

func createSetting() (*setting, error) {
	setRef := C.SCDynamicStoreCopyProxies(C.SCDynamicStoreRef(unsafe.Pointer(nil)))
	if unsafe.Pointer(setRef) == nil {
		return nil, fmt.Errorf("failed to get system proxy settings")
	}

	return &setting{
		ref: C.CFDictionaryRef(setRef),
	}, nil
}

func GetAll() (*Info, *Info, error) {
	s, err := createSetting()
	if err != nil {
		return nil, nil, err
	}
	defer s.Close()

	httpProxy, err := getHTTP(s.ref)
	if err != nil {
		return nil, nil, err
	}

	httpsProxy, err := getHTTPS(s.ref)
	if err != nil {
		return nil, nil, err
	}

	return httpProxy, httpsProxy, nil
}

func getHTTP(settings C.CFDictionaryRef) (*Info, error) {
	raw := C.getHTTP(settings)
	if raw.enabled == 0 {
		return nil, nil
	}

	return &Info{
		Host: C.GoString(&raw.host[0]),
		Port: uint16(raw.port),
	}, nil
}

func GetHTTP() (*Info, error) {
	s, err := createSetting()
	if err != nil {
		return nil, err
	}
	defer s.Close()

	return getHTTP(s.ref)
}

func getHTTPS(settings C.CFDictionaryRef) (*Info, error) {
	raw := C.getHTTPS(settings)
	if raw.enabled == 0 {
		return nil, nil
	}

	return &Info{
		Host: C.GoString(&raw.host[0]),
		Port: uint16(raw.port),
	}, nil
}

func GetHTTPS() (*Info, error) {
	s, err := createSetting()
	if err != nil {
		return nil, err
	}
	defer s.Close()

	return getHTTPS(s.ref)
}
