package httpx

import (
	"context"
	"encoding/json"
	"fmt"
)

type rpcReq struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int64       `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type rpcResp struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int64           `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *rpcErr         `json:"error,omitempty"`
}

type rpcErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func CallRPC(ctx context.Context, d Doer, url string, id int64, method string, params interface{}, out interface{}) error {
	req := rpcReq{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}
	body, _ := json.Marshal(req)

	raw, err := DoJSON(ctx, d, "POST", url, body, nil)
	if err != nil {
		return err
	}

	var resp rpcResp
	if err := json.Unmarshal(raw, &resp); err != nil {
		return err
	}
	if resp.Error != nil {
		return fmt.Errorf("rpc error %d: %s", resp.Error.Code, resp.Error.Message)
	}
	if out == nil {
		return nil
	}
	return json.Unmarshal(resp.Result, out)
}
