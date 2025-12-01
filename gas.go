package infura

import (
	"context"
	"fmt"
)

// GetSuggestedGasFees retrieves suggested gas fees for a given chain ID
// If API Key Secret is provided, uses Basic Auth: /networks/{chainId}/suggestedGasFees
// If only API Key is provided, uses URL path auth: /v3/{apiKey}/networks/{chainId}/suggestedGasFees
func (c *Client) GetSuggestedGasFees(ctx context.Context, chainID int64) (*SuggestedGasFees, error) {
	var endpoint string
	if c.hasSecret() {
		// Basic Auth: API Key + Secret
		endpoint = fmt.Sprintf("/networks/%d/suggestedGasFees", chainID)
	} else {
		// URL path auth: API Key only
		endpoint = fmt.Sprintf("/v3/%s/networks/%d/suggestedGasFees", c.apiKey, chainID)
	}

	var result SuggestedGasFees
	if err := c.doJSONRequest(ctx, "GET", endpoint, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetBaseFeeHistory retrieves base fee history for a given chain ID
// If API Key Secret is provided, uses Basic Auth: /networks/{chainId}/baseFeeHistory
// If only API Key is provided, uses URL path auth: /v3/{apiKey}/networks/{chainId}/baseFeeHistory
// The API returns an array of strings directly
func (c *Client) GetBaseFeeHistory(ctx context.Context, chainID int64) (BaseFeeHistory, error) {
	var endpoint string
	if c.hasSecret() {
		// Basic Auth: API Key + Secret
		endpoint = fmt.Sprintf("/networks/%d/baseFeeHistory", chainID)
	} else {
		// URL path auth: API Key only
		endpoint = fmt.Sprintf("/v3/%s/networks/%d/baseFeeHistory", c.apiKey, chainID)
	}

	var result BaseFeeHistory
	if err := c.doJSONRequest(ctx, "GET", endpoint, nil, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetBaseFeePercentile retrieves base fee percentile for a given chain ID
// If API Key Secret is provided, uses Basic Auth: /networks/{chainId}/baseFeePercentile
// If only API Key is provided, uses URL path auth: /v3/{apiKey}/networks/{chainId}/baseFeePercentile
func (c *Client) GetBaseFeePercentile(ctx context.Context, chainID int64) (*BaseFeePercentile, error) {
	var endpoint string
	if c.hasSecret() {
		// Basic Auth: API Key + Secret
		endpoint = fmt.Sprintf("/networks/%d/baseFeePercentile", chainID)
	} else {
		// URL path auth: API Key only
		endpoint = fmt.Sprintf("/v3/%s/networks/%d/baseFeePercentile", c.apiKey, chainID)
	}

	var result BaseFeePercentile
	if err := c.doJSONRequest(ctx, "GET", endpoint, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetBusyThreshold retrieves busy threshold for a given chain ID
// If API Key Secret is provided, uses Basic Auth: /networks/{chainId}/busyThreshold
// If only API Key is provided, uses URL path auth: /v3/{apiKey}/networks/{chainId}/busyThreshold
func (c *Client) GetBusyThreshold(ctx context.Context, chainID int64) (*BusyThreshold, error) {
	var endpoint string
	if c.hasSecret() {
		// Basic Auth: API Key + Secret
		endpoint = fmt.Sprintf("/networks/%d/busyThreshold", chainID)
	} else {
		// URL path auth: API Key only
		endpoint = fmt.Sprintf("/v3/%s/networks/%d/busyThreshold", c.apiKey, chainID)
	}

	var result BusyThreshold
	if err := c.doJSONRequest(ctx, "GET", endpoint, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
