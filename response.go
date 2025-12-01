package infura

// SuggestedGasFees represents the response from the suggestedGasFees endpoint
type SuggestedGasFees struct {
	Low    GasFeeLevel `json:"low"`
	Medium GasFeeLevel `json:"medium"`
	High   GasFeeLevel `json:"high"`

	EstimatedBaseFee           string   `json:"estimatedBaseFee"`
	NetworkCongestion          float64  `json:"networkCongestion"`
	LatestPriorityFeeRange     []string `json:"latestPriorityFeeRange"`
	HistoricalPriorityFeeRange []string `json:"historicalPriorityFeeRange"`
	HistoricalBaseFeeRange     []string `json:"historicalBaseFeeRange"`
	PriorityFeeTrend           string   `json:"priorityFeeTrend"`
	BaseFeeTrend               string   `json:"baseFeeTrend"`
}

// GasFeeLevel represents a gas fee level (low, medium, or high)
type GasFeeLevel struct {
	SuggestedMaxPriorityFeePerGas string `json:"suggestedMaxPriorityFeePerGas"`
	SuggestedMaxFeePerGas         string `json:"suggestedMaxFeePerGas"`
	MinWaitTimeEstimate           int64  `json:"minWaitTimeEstimate"`
	MaxWaitTimeEstimate           int64  `json:"maxWaitTimeEstimate"`
}

// BaseFeeHistory represents the response from the baseFeeHistory endpoint
// The API directly returns an array of strings
type BaseFeeHistory []string

// BaseFeePercentile represents the response from the baseFeePercentile endpoint
type BaseFeePercentile struct {
	BaseFeePercentile string `json:"baseFeePercentile"`
}

// BusyThreshold represents the response from the busyThreshold endpoint
type BusyThreshold struct {
	BusyThreshold string `json:"busyThreshold"`
}
