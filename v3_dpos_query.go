package gosdk

import (
	"github.com/axengine/httpc"
	"github.com/pkg/errors"
	"strconv"
)

func (p *APISDK) V3DPosPool() (*DPOSPool, error) {
	resp := struct {
		IsSuccess bool     `json:"isSuccess"`
		Message   string   `json:"message"`
		Result    DPOSPool `json:"result"`
	}{}
	err := httpc.New(p.baseUrl).Path("/v3/dpos/pool").Get(&resp)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess {
		return nil, errors.New(resp.Message)
	}

	return &resp.Result, nil
}

func (p *APISDK) V3DPosPoolLogs(cursor, limit, height int64, order string) ([]DPOSPoolLog, error) {
	resp := struct {
		IsSuccess bool          `json:"isSuccess"`
		Message   string        `json:"message"`
		Result    []DPOSPoolLog `json:"result"`
	}{}
	err := httpc.New(p.baseUrl).Path("/v3/dpos/poollogs").
		Query("cursor", strconv.FormatInt(cursor, 10)).
		Query("limit", strconv.FormatInt(limit, 10)).
		Query("order", order).
		Query("height", strconv.FormatInt(height, 10)).Get(&resp)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess {
		return nil, errors.New(resp.Message)
	}

	return resp.Result, nil
}

func (p *APISDK) V3DPosTcnLogs(cursor, limit, height int64, address string, order string) ([]DPOSTcnLog, error) {
	resp := struct {
		IsSuccess bool         `json:"isSuccess"`
		Message   string       `json:"message"`
		Result    []DPOSTcnLog `json:"result"`
	}{}
	err := httpc.New(p.baseUrl).Path("/v3/dpos/tcnlogs").
		Query("cursor", strconv.FormatInt(cursor, 10)).
		Query("limit", strconv.FormatInt(limit, 10)).
		Query("order", order).
		Query("height", strconv.FormatInt(height, 10)).
		Query("address", address).Get(&resp)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess {
		return nil, errors.New(resp.Message)
	}

	return resp.Result, nil
}

func (p *APISDK) V3DPosTinLogs(cursor, limit, height int64, address string, order string) ([]DPOSTinLog, error) {
	resp := struct {
		IsSuccess bool         `json:"isSuccess"`
		Message   string       `json:"message"`
		Result    []DPOSTinLog `json:"result"`
	}{}
	err := httpc.New(p.baseUrl).Path("/v3/dpos/tinlogs").
		Query("cursor", strconv.FormatInt(cursor, 10)).
		Query("limit", strconv.FormatInt(limit, 10)).
		Query("order", order).
		Query("height", strconv.FormatInt(height, 10)).
		Query("address", address).Get(&resp)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess {
		return nil, errors.New(resp.Message)
	}

	return resp.Result, nil
}

func (p *APISDK) V3DPosRankLogs(cursor, limit, height int64, address string, order string) ([]DPOSRankLog, error) {
	resp := struct {
		IsSuccess bool          `json:"isSuccess"`
		Message   string        `json:"message"`
		Result    []DPOSRankLog `json:"result"`
	}{}
	err := httpc.New(p.baseUrl).Path("/v3/dpos/ranklogs").
		Query("cursor", strconv.FormatInt(cursor, 10)).
		Query("limit", strconv.FormatInt(limit, 10)).
		Query("order", order).
		Query("height", strconv.FormatInt(height, 10)).
		Query("address", address).Get(&resp)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess {
		return nil, errors.New(resp.Message)
	}

	return resp.Result, nil
}

func (p *APISDK) V3DPosNodes() ([]BeanValidator, error) {
	resp := struct {
		IsSuccess bool            `json:"isSuccess"`
		Message   string          `json:"message"`
		Result    []BeanValidator `json:"result"`
	}{}
	err := httpc.New(p.baseUrl).Path("/v3/nodes").Get(&resp)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess {
		return nil, errors.New(resp.Message)
	}

	return resp.Result, nil
}

func (p *APISDK) V3DPosNode(address string) (*BeanValidator, error) {
	resp := struct {
		IsSuccess bool          `json:"isSuccess"`
		Message   string        `json:"message"`
		Result    BeanValidator `json:"result"`
	}{}
	err := httpc.New(p.baseUrl).Path("/v3/nodes/" + address).Get(&resp)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess {
		return nil, errors.New(resp.Message)
	}

	return &resp.Result, nil
}

func (p *APISDK) V3DPosNodeVoters(address string) (*BeanVots, error) {
	resp := struct {
		IsSuccess bool     `json:"isSuccess"`
		Message   string   `json:"message"`
		Result    BeanVots `json:"result"`
	}{}
	err := httpc.New(p.baseUrl).Path("/v3/nodes/" + address + "/voters").Get(&resp)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess {
		return nil, errors.New(resp.Message)
	}

	return &resp.Result, nil
}

func (p *APISDK) V3DPosVoter(address string) (*BeanVoter, error) {
	resp := struct {
		IsSuccess bool      `json:"isSuccess"`
		Message   string    `json:"message"`
		Result    BeanVoter `json:"result"`
	}{}
	err := httpc.New(p.baseUrl).Path("/v3/voters/" + address).Get(&resp)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess {
		return nil, errors.New(resp.Message)
	}

	return &resp.Result, nil
}
