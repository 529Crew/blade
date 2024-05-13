package requests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/valyala/fasthttp"
)

type GetAccountInfoResponse struct {
	JsonRpc string
	Result  rpc.GetAccountInfoResult
}

func GetAccountDataFast(rpcURL string, account solana.PublicKey) (*rpc.GetAccountInfoResult, error) {
	payload := `{"method":"getAccountInfo","jsonrpc":"2.0","params":["` + account.String() + `",{"encoding":"base64","commitment":"processed"}],"id":"1"}`
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(rpcURL)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentTypeBytes(headerContentTypeJson)
	req.Header.Set("Accept", "application/json")
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

	var accountInfo GetAccountInfoResponse
	jsonErr := json.Unmarshal(respBody, &accountInfo)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &accountInfo.Result, nil
}
