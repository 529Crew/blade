package requests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/529Crew/blade/internal/types"
	"github.com/valyala/fasthttp"
)

func Coins(creator string) (*types.Coins, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(fmt.Sprintf("https://client-api-2-74b1891ee9f9.herokuapp.com/coins?offset=0&limit=100&sort=created_timestamp&order=desc&includeNsfw=false&creator=%s", creator))
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set("Accept", "application/json")

	resp := fasthttp.AcquireResponse()
	timeoutErr := FastHttpClient.DoTimeout(req, resp, time.Second*5)

	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	if timeoutErr != nil {
		return nil, fmt.Errorf("connection error: %s", timeoutErr)
	}

	statusCode := resp.StatusCode()
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", statusCode)
	}

	respBody := resp.Body()

	var coins types.Coins
	err := json.Unmarshal(respBody, &coins)
	if err != nil {
		return nil, fmt.Errorf("unmarshal json error: %s", err)
	}

	return &coins, nil
}
