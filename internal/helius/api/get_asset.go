package helius_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/529Crew/blade/internal/config"
	"github.com/529Crew/blade/internal/requests"
	"github.com/529Crew/blade/internal/types"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

func GetAsset(mint string) (*types.GetAssetResponse, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(fmt.Sprintf("https://mainnet.helius-rpc.com/?api-key=%s", config.Get().HeliusApiKey))
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.Set("Accept", "application/json")
	req.SetBodyString(fmt.Sprintf(`{"jsonrpc":"2.0","id":"%s","method":"getAsset","params":{"id":"%s","displayOptions":{"showFungible":true}}}`, uuid.New().String(), mint))

	resp := fasthttp.AcquireResponse()
	timeoutErr := requests.FastHttpClient.DoTimeout(req, resp, time.Second*5)

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

	var getAssetResponse types.GetAssetResponse
	err := json.Unmarshal(respBody, &getAssetResponse)
	if err != nil {
		return nil, fmt.Errorf("unmarshal json error: %v", err)
	}

	return &getAssetResponse, nil
}
