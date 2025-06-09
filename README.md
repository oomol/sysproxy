# sysproxy is a Golang library for retrieving system proxy settings

### Supported platforms and proxy type:
- windows (http)
- MacOS (http/https)

### Retrieve the system HTTP proxy:
```go
httpInfo, err := GetHTTP()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v:", err)
	}

	if httpInfo != nil {
		fmt.Printf("HTTP Proxy Host: %v\n", httpInfo.Host)
		fmt.Printf("HTTP Proxy Port: %v\n", httpInfo.Port)
	}
```

### Retrieve the system HTTPS proxy:
```go
	httpsInfo, err := GetHTTPS()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v:", err)
	}

	if httpsInfo != nil {
		fmt.Printf("HTTPS Proxy Host: %v\n", httpsInfo.Host)
		fmt.Printf("HTTPS Proxy Port: %v\n", httpsInfo.Port)
	}
```

You can also use `GetAll` to efficiently retrieve all HTTP/HTTPS proxy settings.
```go
    httpInfo, httpsInfo, err := GetAll()
```

**Note:** Only HTTP proxy supported on the Windows platform, to ensure accurate results, please use `GetHTTP`, the `GetHTTPS` function just return nil with no error.
