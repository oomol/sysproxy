# sysproxy

A Golang library for retrieving system proxy settings

### Supported platforms and proxy types:

- Windows (HTTP)
- Darwin (HTTP/HTTPS)

### Get HTTP/HTTPS proxies:

```go
package main

import (
	"fmt"

	"github.com/oomol-lab/sysproxy"
)

func main() {
	httpInfo, _ := sysproxy.GetHTTP()
	if httpInfo != nil {
		fmt.Printf("http proxy host: %s, port: %d \n", httpInfo.Host, httpInfo.Port)
	}

	// Since Windows does not support setting an HTTPS proxy
	// GetHTTPS will always return nil with no error on the Windows platform
	httpsInfo, _ := sysproxy.GetHTTPS()
	if httpsInfo != nil {
		fmt.Printf("https proxy host: %s, port: %d \n", httpsInfo.Host, httpsInfo.Port)
	}

	// You can also using sysproxy.GetAll() get all proxy settings
	// httpInfo, httpsInfo, _ := sysproxy.GetAll()
}
```

