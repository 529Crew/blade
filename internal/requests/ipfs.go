package requests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/529Crew/blade/internal/types"
	"github.com/valyala/fasthttp"
)

func IpfsData(uri string) (*types.IpfsResponse, error) {
	pinataUri := strings.Replace(uri, "https://cf-ipfs.com/ipfs/", "https://pump.mypinata.cloud/ipfs/", 1)

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(pinataUri)
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set("Accept", "application/json")

	resp := fasthttp.AcquireResponse()
	timeoutErr := FastHttpClient.DoTimeout(req, resp, time.Second*5)

	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	if timeoutErr != nil {
		return nil, fmt.Errorf("connection error: %v", timeoutErr)
	}

	statusCode := resp.StatusCode()
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", statusCode)
	}

	respBody := resp.Body()

	var pfData types.IpfsResponse
	err := json.Unmarshal(respBody, &pfData)
	if err != nil {
		return nil, fmt.Errorf("unmarshal json error: %v", err)
	}

	return &pfData, nil
}
