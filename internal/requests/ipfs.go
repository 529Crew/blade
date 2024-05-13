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
	urlSplit := strings.Split(uri, "/ipfs/")
	if len(urlSplit) < 2 {
		return nil, fmt.Errorf("uri invalid: %s", uri)
	}

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(fmt.Sprintf("https://flowgocrazy.mypinata.cloud/ipfs/%s", urlSplit[1]))
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
