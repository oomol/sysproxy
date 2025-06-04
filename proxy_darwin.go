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
	"unsafe"
)

func GetProxyInfo() (*HttpProxyInfo, *HttpsProxyInfo) {
	httpProxyInfo := &HttpProxyInfo{}

	settings := C.SCDynamicStoreCopyProxies(C.SCDynamicStoreRef(unsafe.Pointer(nil)))
	if unsafe.Pointer(settings) == nil {
		return nil, nil
	}

	defer C.CFRelease(C.CFTypeRef(settings))

	httpInfo := C.getHttpProxyInfo(settings)
	if httpInfo.enabled != 0 {
		httpProxyInfo.Host = C.GoString(&httpInfo.host[0])
		httpProxyInfo.Port = uint16(httpInfo.port)
	}

	httpsProxyInfo := &HttpsProxyInfo{}
	httpsInfo := C.getHttpsProxyInfo(settings)
	if httpsInfo.enabled != 0 {
		httpsProxyInfo.Host = C.GoString(&httpsInfo.host[0])
		httpsProxyInfo.Port = uint16(httpsInfo.port)
	}

	return httpProxyInfo, httpsProxyInfo
}
