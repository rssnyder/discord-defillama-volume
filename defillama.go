package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Volume struct {
	ID                      string      `json:"id"`
	Name                    string      `json:"name"`
	URL                     string      `json:"url"`
	Description             string      `json:"description"`
	Logo                    string      `json:"logo"`
	GeckoID                 string      `json:"gecko_id"`
	CmcID                   string      `json:"cmcId"`
	Chains                  []any       `json:"chains"`
	Twitter                 string      `json:"twitter"`
	Treasury                string      `json:"treasury"`
	GovernanceID            []string    `json:"governanceID"`
	Github                  []string    `json:"github"`
	DefillamaID             string      `json:"defillamaId"`
	DisplayName             string      `json:"displayName"`
	Total24H                float64     `json:"total24h"`
	TotalAllTime            float64     `json:"totalAllTime"`
	LatestFetchIsOk         bool        `json:"latestFetchIsOk"`
	Disabled                bool        `json:"disabled"`
	Change1D                float64     `json:"change_1d"`
	MethodologyURL          any         `json:"methodologyURL"`
	Methodology             any         `json:"methodology"`
	Module                  any         `json:"module"`
	TotalDataChart          [][]float64 `json:"totalDataChart"`
	TotalDataChartBreakdown [][]any     `json:"totalDataChartBreakdown"`
	ChildProtocols          []string    `json:"childProtocols"`
}

const (
	DefiLlamaVolumeURL = "https://api.llama.fi/summary/dexs/%s?excludeTotalDataChart=true&excludeTotalDataChartBreakdown=true"
)

type Earn struct {
	Tvl    float64 `json:"tvl"`
	MaxAPR float64 `json:"maxAPR"`
	Data   []struct {
		ContractAddress string  `json:"contractAddress"`
		AprBase         float64 `json:"aprBase"`
		AprBonus        float64 `json:"aprBonus"`
		SymbolBase      string  `json:"symbolBase"`
		SymbolBonus     string  `json:"symbolBonus"`
		Tvl             float64 `json:"tvl"`
		Network         string  `json:"network"`
		Link            string  `json:"link"`
	} `json:"data"`
}

func GetVolume(protocol string) (result Volume, err error) {

	req, err := http.NewRequest("GET", fmt.Sprintf(DefiLlamaVolumeURL, protocol), nil)
	if err != nil {
		return result, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0")
	req.Header.Add("accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	results, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(results, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
