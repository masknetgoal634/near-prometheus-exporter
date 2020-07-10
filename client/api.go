package nearapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type StatusResult struct {
	Status struct {
		Version struct {
			Version string `json:"version"`
			Build   string `json:"build"`
		} `json:"version"`
		ChainId string `json:"chain_id"`
		RpcAddr string `json:"rpc_addr"`
		//Validators []string `json:"validators"`
		SyncInfo struct {
			LatestBlockHash   string `json:"latest_block_hash"`
			LatestBlockHeight uint64 `json:"latest_block_height"`
			LatestStateRoot   string `json:"latest_state_root"`
			LatestBlockTime   string `json:"latest_block_time"`
			Syncing           bool   `json:"syncing"`
		} `json:"sync_info"`
	} `json:"result_status"`
}

type Validator struct {
	AccountId string `json:"account_id"`
	PublicKey string `json:"public_key"`
	Stake     string `json:"stake"`
}

type ValidatorsResult struct {
	Validators struct {
		CurrentValidators []struct {
			Validator
			IsSlashed         bool  `json:"is_slashed"`
			Shards            []int `json:"shards"`
			NumProducedBlocks int64 `json:"num_produced_blocks"`
			NumExpectedBlocks int64 `json:"num_expected_blocks"`
		} `json:"current_validators"`
		NextValidators []struct {
			Validator
			Shards []int `json:"shards"`
		} `json:"next_validators"`
		CurrentProposals []struct {
			Validator
		} `json:"current_proposals"`
		EpochStartHeight int64 `json:"epoch_start_height"`
		PrevEpochKickOut []struct {
			AccountId string                            `json:"account_id"`
			Reason    map[string]map[string]interface{} `json:"reason"`
		} `json:"prev_epoch_kickout"`
	} `json:"result_validators"`
}

type Result struct {
	StatusResult
	ValidatorsResult
}

type Client struct {
	httpClient *http.Client
	Endpoint   string
}

func NewClient(endpoint string) *Client {
	timeout := time.Duration(10 * time.Second)
	httpClient := &http.Client{
		Timeout: timeout,
	}
	return &Client{
		Endpoint:   endpoint,
		httpClient: httpClient,
	}
}

func NewClientWith(client *http.Client, endpoint string) *Client {
	return &Client{
		Endpoint:   endpoint,
		httpClient: client,
	}
}

func (c *Client) do(method string, params interface{}) (string, error) {
	payload, err := json.Marshal(map[string]string{
		"query": method,
	})

	if params != "" {
		type Payload struct {
			JsonRPC string      `json:"jsonrpc"`
			Id      string      `json:"id"`
			Method  string      `json:"method"`
			Params  interface{} `json:"params"`
		}
		p := Payload{
			JsonRPC: "2.0",
			Id:      "dontcare",
			Method:  method,
			Params:  params,
		}
		payload, err = json.Marshal(p)
		if err != nil {
			log.Println(err)
		}
	}
	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return "", err
	}

	r, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (c *Client) Get(method string, variables interface{}) (*Result, error) {
	res, err := c.do(method, variables)
	if err != nil {
		return nil, err
	}
	var d Result
	res = strings.Replace(res, "result", fmt.Sprintf("%s_%s", "result", method), -1)
	r := bytes.NewReader([]byte(res))
	err2 := json.NewDecoder(r).Decode(&d)
	if err2 != nil {
		log.Println(err2)
		return nil, err2
	}
	return &d, nil
}
