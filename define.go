package systeminfo

type ProxyInfo struct {
	Host   string
	Port   uint16
	Bypass string
}

type HttpsProxyInfo struct {
	ProxyInfo
}

type HttpProxyInfo struct {
	ProxyInfo
}
