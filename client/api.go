package nearapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
			NumProducedBlocks int   `json:"num_produced_blocks"`
			NumExpectedBlocks int   `json:"num_expected_blocks"`
		} `json:"current_validators"`
		NextValidators []struct {
			Validator
			Shards []int `json:"shards"`
		} `json:"next_validators"`
		CurrentProposals []struct {
			Validator
		} `json:"current_proposals"`
		EpochStartHeight int64 `json:"epoch_start_height"`
	} `json:"result_validators"`
}

type Result struct {
	StatusResult
	ValidatorsResult
}

type Client struct {
	httpClient       *http.Client
	Endpoint         string
	lastCacheTime    int64
	cacheVals        *Result
	cacheStatus      *Result
	seatPrice        int64
	currentStake     int64
	versionNumber    string
	versionBuild     string
	epochStartHeight int64
}

func NewClient(endpoint string) *Client {
	timeout := time.Duration(5 * time.Second)
	httpClient := &http.Client{
		Timeout: timeout,
	}
	return &Client{
		Endpoint:      endpoint,
		httpClient:    httpClient,
		versionNumber: "0",
		versionBuild:  "0",
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
		log.Fatalln(err)
	}

	r, err := c.httpClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
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

func (c *Client) freshCache() bool {
	now := time.Now().Unix()
	if (now - c.lastCacheTime) <= 10 {
		return true
	}
	return false
}

// Near json api
// Get Status
func (c *Client) Status() (*Result, error) {
	if c.freshCache() {
		return c.cacheStatus, nil
	}
	var err error
	c.cacheStatus, err = c.Get("status", nil)
	c.lastCacheTime = time.Now().Unix()
	c.versionNumber = c.cacheStatus.Status.Version.Version
	c.versionBuild = c.cacheStatus.Status.Version.Build
	return c.cacheStatus, err
}

// Get Current and Next validators
func (c *Client) Validators() (*Result, error) {
	if c.freshCache() {
		return c.cacheVals, nil
	}
	r, err := c.Status()
	if err != nil {
		fmt.Println(err)
	}
	blockHeight := r.Status.SyncInfo.LatestBlockHeight
	c.cacheVals, err = c.Get("validators", []uint64{blockHeight})
	c.epochStartHeight = c.cacheVals.Validators.EpochStartHeight
	return c.cacheVals, err
}

func (c *Client) Syncing() (bool, error) {
	r, err := c.Status()
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return r.Status.SyncInfo.Syncing, nil
}

func (c *Client) LatestBlockHash() string {
	r, err := c.Status()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return r.Status.SyncInfo.LatestBlockHash
}

func (c *Client) LatestBlockHeight() (uint64, error) {
	r, err := c.Status()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return r.Status.SyncInfo.LatestBlockHeight, err
}

func stringToInt64(s string) int64 {
	v, err := strconv.ParseInt(s[0:6], 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	return int64(v)
}

func (c *Client) ProducedBlocks(accountId string) (int, error) {
	r, err := c.Validators()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	pb := 0
	c.seatPrice = 0
	for _, v := range r.Validators.CurrentValidators {
		t := stringToInt64(v.Stake)
		if c.seatPrice == 0 {
			c.seatPrice = t
		}
		if c.seatPrice > t {
			c.seatPrice = t
		}
		if v.AccountId == accountId {
			pb = v.NumProducedBlocks
			c.currentStake = t
		}
	}
	return pb, nil
}

func (c *Client) ExpectedBlocks(accountId string) (int, error) {
	r, err := c.Validators()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	for _, v := range r.Validators.CurrentValidators {
		if v.AccountId == accountId {
			return v.NumExpectedBlocks, nil
		}
	}
	return 0, errors.New("account id not found")
}

func (c *Client) SeatPrice() (int64, error) {
	return c.seatPrice, nil
}

func (c *Client) CurrentStake() (int64, error) {
	return c.currentStake, nil
}

func (c *Client) VersionNumber() (string, error) {
	return c.versionNumber, nil
}

func (c *Client) VersionBuild() (string, error) {
	return c.versionBuild, nil
}

func (c *Client) EpochStartHeight() (int64, error) {
	return c.epochStartHeight, nil
}
