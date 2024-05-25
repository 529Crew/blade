package requests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/529Crew/blade/internal/config"
	"github.com/valyala/fasthttp"
)

type SubmitBloxrouteTxResponse struct {
	Signature string `json:"signature"`
	UUID      string `json:"uuid"`
}

func SubmitBloxrouteTx(tx string) (*SubmitBloxrouteTxResponse, error) {
	cfg := config.Get()

	payload := fmt.Sprintf(`{"transaction":{"content":"%s"},"frontRunningProtection":false,"useStakedRPCs":true}`, tx)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("https://ny.solana.dex.blxrbdn.com/api/v2/submit")
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentTypeBytes(headerContentTypeJson)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", cfg.BloxrouteAuthHeader)
	req.SetBodyString(payload)

	resp := fasthttp.AcquireResponse()
	timeoutErr := FastHttpClient.DoTimeout(req, resp, time.Second*5)
	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	if timeoutErr != nil {
		return nil, fmt.Errorf("connection error: %s", timeoutErr)
	}

	statusCode := resp.StatusCode()
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("received status code: %d", statusCode)
	}

	respBody := resp.Body()

	var submitBloxrouteTxResponse SubmitBloxrouteTxResponse
	jsonErr := json.Unmarshal(respBody, &submitBloxrouteTxResponse)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &submitBloxrouteTxResponse, nil
}
