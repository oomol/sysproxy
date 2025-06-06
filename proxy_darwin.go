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

ProxyInfo getHttpProxyInfo(CFDictionaryRef settings) {
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

ProxyInfo getHttpsProxyInfo(CFDictionaryRef settings) {
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
	"unsafe"
)

func GetInfo() (*ProxyInfo, *ProxyInfo, error) {
	settings := C.SCDynamicStoreCopyProxies(C.SCDynamicStoreRef(unsafe.Pointer(nil)))
	if unsafe.Pointer(settings) == nil {
		return nil, nil, fmt.Errorf("cannot get proxy info, get SCDynamicStoreCopyProxies error")
	}

	defer C.CFRelease(C.CFTypeRef(settings))

	httpProxy, err := GetHttpProxy(C.CFDictionaryRef(settings))
	if err != nil {
		return nil, nil, err
	}

	httpsProxy, err := GetHttpsProxy(C.CFDictionaryRef(settings))
	if err != nil {
		return nil, nil, err
	}

	return httpProxy, httpsProxy, nil
}

func GetHttpProxy(settings C.CFDictionaryRef) (*ProxyInfo, error) {
	info := &ProxyInfo{}
	httpInfo := C.getHttpProxyInfo(settings)

	if httpInfo.enabled != 0 {
		info.Host = C.GoString(&httpInfo.host[0])
		info.Port = uint16(httpInfo.port)
	}

	return info, nil
}

func GetHttpsProxy(settings C.CFDictionaryRef) (*ProxyInfo, error) {
	info := &ProxyInfo{}
	httpsInfo := C.getHttpsProxyInfo(settings)
	if httpsInfo.enabled != 0 {
		info.Host = C.GoString(&httpsInfo.host[0])
		info.Port = uint16(httpsInfo.port)
	}

	return info, nil
}
