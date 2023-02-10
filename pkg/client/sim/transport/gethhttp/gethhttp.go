package gethhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/blocknative/dreamboat/pkg/client/sim/types"
	"github.com/lthibault/log"
)

type Client struct {
	namespace string
	client    *http.Client
	address   string
	l         log.Logger
}

func NewClient(address string, namespace string, l log.Logger) *Client {
	return &Client{
		namespace: namespace,
		address:   address,
		l:         l,
		client:    &http.Client{},
	}
}

func (c *Client) ValidateBlock(ctx context.Context, params []byte) (rrr types.RpcRawResponse, err error) {
	return justsend(ctx, c.client, c.address, io.MultiReader(strings.NewReader(`{"id": 1, "method": "`+c.namespace+`_validateBuilderSubmissionV1","params": `), bytes.NewReader(params), strings.NewReader(`}`)))
}

func justsend(ctx context.Context, client *http.Client, url string, body io.Reader) (rrr types.RpcRawResponse, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return rrr, err
	}
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	res, err := client.Do(req)
	if err != nil {
		return rrr, err
	}
	defer res.Body.Close()

	rrr = types.RpcRawResponse{}
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&rrr)
	return rrr, err
}
