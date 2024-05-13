package requests

import (
	"time"

	"github.com/valyala/fasthttp"
)

var (
	FastHttpClient *fasthttp.Client = &fasthttp.Client{
		MaxIdleConnDuration:           time.Minute,
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
		MaxConnsPerHost:               100000,
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      0,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}
)

var (
	headerContentTypeJson = []byte("application/json")
)
