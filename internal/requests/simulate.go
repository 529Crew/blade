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

type SimulatedTransactionResponse struct {
	Jsonrpc string
	Result  struct {
		Context struct {
			Slot int
		}
		Value struct {
			Err           interface{}
			Logs          []string
			Accounts      []*rpc.Account
			UnitsConsumed uint64
		}
	}
	Id int
}

func SimulateTxFast(rpcURL string, tx *solana.Transaction) (*SimulatedTransactionResponse, error) {
	encodedTx, err := tx.ToBase64()
	if err != nil {
		return nil, err
	}

	payload := `{"jsonrpc":"2.0","id":1,"method":"simulateTransaction","params":["` + encodedTx + `",{"replaceRecentBlockhash":true,"commitment":"processed","encoding":"base64"}]}`

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

	var simulationResponse SimulatedTransactionResponse
	if err := json.Unmarshal(respBody, &simulationResponse); err != nil {
		return nil, err
	}

	return &simulationResponse, nil
}
