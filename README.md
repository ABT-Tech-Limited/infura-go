# infura-go

Golang client for Infura Gas API

## 简介

`infura-go` 是一个用于访问 Infura Gas API 的 Go 语言客户端库。它提供了简单易用的接口来获取以太坊网络的 Gas 费用建议。

## 功能特性

- 支持 Infura Gas API 的所有端点：
  - `suggestedGasFees` - 获取 Gas 费用建议
  - `baseFeeHistory` - 获取基础费用历史
  - `baseFeePercentile` - 获取基础费用百分位数
  - `busyThreshold` - 获取网络繁忙阈值
- 支持两种认证方式：
  - 仅使用 API Key（API Key 放在 URL 路径中）
  - 使用 API Key + Secret（Basic Authentication）
- 自动检测认证方式（如果 Secret 为空，自动使用仅 API Key 方式）
- 支持 context.Context，便于控制请求的取消和超时
- 支持调试模式（WithDebug），打印详细的 HTTP 请求和响应信息
- 支持自定义 HTTP 客户端和超时设置
- 完整的测试覆盖

## 安装

```bash
go get github.com/ABT-Tech-Limited/infura-go
```

## 认证方式

Infura Gas API 支持两种认证方式：

1. **仅使用 API Key**：将 API Key 放在 URL 路径中（`/v3/{apiKey}/networks/{chainId}/suggestedGasFees`）
2. **使用 API Key + Secret**：使用 Basic Authentication（`/networks/{chainId}/suggestedGasFees`）

## 使用方法

### 方式一：仅使用 API Key

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/ABT-Tech-Limited/infura-go"
)

func main() {
    ctx := context.Background()
    
    // 创建客户端（仅使用 API Key）
    client := infura.NewClientWithAPIKey("your-api-key")
    
    // 获取以太坊主网（chain ID: 1）的 Gas 费用建议
    gasFees, err := client.GetSuggestedGasFees(ctx, 1)
    if err != nil {
        log.Fatal(err)
    }
    
    // 打印结果
    fmt.Printf("Low: %s gwei\n", gasFees.Low.SuggestedMaxFeePerGas)
    fmt.Printf("Medium: %s gwei\n", gasFees.Medium.SuggestedMaxFeePerGas)
    fmt.Printf("High: %s gwei\n", gasFees.High.SuggestedMaxFeePerGas)
    fmt.Printf("Estimated Base Fee: %s gwei\n", gasFees.EstimatedBaseFee)
    fmt.Printf("Network Congestion: %.2f%%\n", gasFees.NetworkCongestion*100)
    
    // 获取基础费用历史
    baseFeeHistory, err := client.GetBaseFeeHistory(ctx, 1)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Base Fee History: %v\n", baseFeeHistory)
    
    // 获取基础费用百分位数
    baseFeePercentile, err := client.GetBaseFeePercentile(ctx, 1)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Base Fee Percentile: %s\n", baseFeePercentile.BaseFeePercentile)
    
    // 获取网络繁忙阈值
    busyThreshold, err := client.GetBusyThreshold(ctx, 1)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Busy Threshold: %s\n", busyThreshold.BusyThreshold)
}
```

### 方式二：使用 API Key + Secret（Basic Auth）

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/ABT-Tech-Limited/infura-go"
)

func main() {
    ctx := context.Background()
    
    // 创建客户端（使用 API Key + Secret）
    client := infura.NewClient("your-api-key", "your-api-secret")
    
    // 获取以太坊主网（chain ID: 1）的 Gas 费用建议
    gasFees, err := client.GetSuggestedGasFees(ctx, 1)
    if err != nil {
        log.Fatal(err)
    }
    
    // 打印结果
    fmt.Printf("Low: %s gwei\n", gasFees.Low.SuggestedMaxFeePerGas)
    fmt.Printf("Medium: %s gwei\n", gasFees.Medium.SuggestedMaxFeePerGas)
    fmt.Printf("High: %s gwei\n", gasFees.High.SuggestedMaxFeePerGas)
    fmt.Printf("Estimated Base Fee: %s gwei\n", gasFees.EstimatedBaseFee)
    fmt.Printf("Network Congestion: %.2f%%\n", gasFees.NetworkCongestion*100)
}
```

### 使用空 Secret（自动使用 API Key 方式）

如果传入空的 Secret，客户端会自动使用仅 API Key 的认证方式：

```go
// 传入空字符串作为 Secret，会自动使用 API Key 方式
ctx := context.Background()
client := infura.NewClient("your-api-key", "")
gasFees, err := client.GetSuggestedGasFees(ctx, 1)
```

### 启用调试模式

使用 `WithDebug(true)` 选项可以启用调试模式，打印详细的 HTTP 请求和响应信息：

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/ABT-Tech-Limited/infura-go"
)

func main() {
    ctx := context.Background()
    
    // 创建客户端并启用调试模式
    client := infura.NewClientWithOptions(
        "your-api-key",
        "your-api-secret",
        infura.WithDebug(true), // 启用调试模式
    )
    
    // 调用 API 时，会在控制台打印详细的请求和响应信息
    gasFees, err := client.GetSuggestedGasFees(ctx, 1)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Low: %s gwei\n", gasFees.Low.SuggestedMaxFeePerGas)
}
```
```

调试模式会打印以下信息：
- **请求信息**：HTTP 方法、URL、协议版本、Host、所有请求头（Authorization 头会被部分掩码以保护安全）
- **响应头信息**：状态码、协议版本、所有响应头
- **响应体**：格式化的 JSON 响应体
- **解析后的对象**：解析后的 Go 结构体

### 高级用法

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "time"
    
    "github.com/ABT-Tech-Limited/infura-go"
)

func main() {
    // 使用自定义选项创建客户端
    client := infura.NewClientWithOptions(
        "your-api-key",
        "your-api-secret",
        infura.WithBaseURL("https://custom-gas-api.url"), // 自定义 API 地址
        infura.WithTimeout(60*time.Second),               // 自定义超时时间
        infura.WithDebug(true),                           // 启用调试模式，打印 HTTP 请求和响应
        infura.WithHTTPClient(&http.Client{              // 自定义 HTTP 客户端
            Timeout: 30 * time.Second,
        }),
    )
    
    ctx := context.Background()
    
    // 获取不同链的 Gas 费用
    // 以太坊主网
    ethGasFees, _ := client.GetSuggestedGasFees(ctx, 1)
    
    // Polygon (chain ID: 137)
    polygonGasFees, _ := client.GetSuggestedGasFees(ctx, 137)
    
    fmt.Println("Ethereum:", ethGasFees.EstimatedBaseFee)
    fmt.Println("Polygon:", polygonGasFees.EstimatedBaseFee)
}
```

## API 参考

### Client

#### NewClient

创建新的 Infura Gas API 客户端（使用 API Key + Secret，Basic Auth）。

```go
func NewClient(apiKey, apiKeySecret string) *Client
```

**注意**：如果 `apiKeySecret` 为空字符串，将自动使用仅 API Key 的认证方式。

#### NewClientWithAPIKey

创建仅使用 API Key 的客户端（API Key 放在 URL 路径中）。

```go
func NewClientWithAPIKey(apiKey string) *Client
```

#### NewClientWithOptions

使用自定义选项创建客户端（使用 API Key + Secret）。

```go
func NewClientWithOptions(apiKey, apiKeySecret string, opts ...ClientOption) *Client
```

**注意**：如果 `apiKeySecret` 为空字符串，将自动使用仅 API Key 的认证方式。

#### NewClientWithAPIKeyAndOptions

使用自定义选项创建仅使用 API Key 的客户端。

```go
func NewClientWithAPIKeyAndOptions(apiKey string, opts ...ClientOption) *Client
```

可用的选项：
- `WithBaseURL(baseURL string)` - 设置自定义基础 URL
- `WithTimeout(timeout time.Duration)` - 设置 HTTP 请求超时时间
- `WithHTTPClient(httpClient *http.Client)` - 设置自定义 HTTP 客户端
- `WithDebug(debug bool)` - 启用调试模式，打印详细的 HTTP 请求和响应信息（包括 headers、body 等）

### Gas API

#### GetSuggestedGasFees

获取指定链的 Gas 费用建议。

```go
func (c *Client) GetSuggestedGasFees(ctx context.Context, chainID int64) (*SuggestedGasFees, error)
```

**参数：**
- `ctx` - 上下文，用于控制请求的取消和超时
- `chainID` - 链 ID（例如：1 表示以太坊主网）

**返回：**
- `*SuggestedGasFees` - Gas 费用建议数据
- `error` - 错误信息

**认证方式：**
- 如果客户端使用 API Key + Secret，会使用 Basic Auth：`/networks/{chainId}/suggestedGasFees`
- 如果客户端仅使用 API Key，会将 API Key 放在 URL 路径中：`/v3/{apiKey}/networks/{chainId}/suggestedGasFees`

#### GetBaseFeeHistory

获取指定链的基础费用历史。

```go
func (c *Client) GetBaseFeeHistory(ctx context.Context, chainID int64) (*BaseFeeHistory, error)
```

**参数：**
- `ctx` - 上下文，用于控制请求的取消和超时
- `chainID` - 链 ID（例如：1 表示以太坊主网）

**返回：**
- `BaseFeeHistory` - 基础费用历史数据（字符串数组）
- `error` - 错误信息

**注意**：API 直接返回字符串数组，而不是包含 `baseFeeHistory` 字段的对象。

**认证方式：**
- 如果客户端使用 API Key + Secret，会使用 Basic Auth：`/networks/{chainId}/baseFeeHistory`
- 如果客户端仅使用 API Key，会将 API Key 放在 URL 路径中：`/v3/{apiKey}/networks/{chainId}/baseFeeHistory`

#### GetBaseFeePercentile

获取指定链的基础费用百分位数。

```go
func (c *Client) GetBaseFeePercentile(ctx context.Context, chainID int64) (*BaseFeePercentile, error)
```

**参数：**
- `ctx` - 上下文，用于控制请求的取消和超时
- `chainID` - 链 ID（例如：1 表示以太坊主网）

**返回：**
- `*BaseFeePercentile` - 基础费用百分位数数据
- `error` - 错误信息

**认证方式：**
- 如果客户端使用 API Key + Secret，会使用 Basic Auth：`/networks/{chainId}/baseFeePercentile`
- 如果客户端仅使用 API Key，会将 API Key 放在 URL 路径中：`/v3/{apiKey}/networks/{chainId}/baseFeePercentile`

#### GetBusyThreshold

获取指定链的网络繁忙阈值。

```go
func (c *Client) GetBusyThreshold(ctx context.Context, chainID int64) (*BusyThreshold, error)
```

**参数：**
- `ctx` - 上下文，用于控制请求的取消和超时
- `chainID` - 链 ID（例如：1 表示以太坊主网）

**返回：**
- `*BusyThreshold` - 网络繁忙阈值数据
- `error` - 错误信息

**认证方式：**
- 如果客户端使用 API Key + Secret，会使用 Basic Auth：`/networks/{chainId}/busyThreshold`
- 如果客户端仅使用 API Key，会将 API Key 放在 URL 路径中：`/v3/{apiKey}/networks/{chainId}/busyThreshold`

### 响应结构

#### SuggestedGasFees

```go
type SuggestedGasFees struct {
    Low    GasFeeLevel `json:"low"`
    Medium GasFeeLevel `json:"medium"`
    High   GasFeeLevel `json:"high"`
    
    EstimatedBaseFee          string   `json:"estimatedBaseFee"`
    NetworkCongestion         float64  `json:"networkCongestion"`
    LatestPriorityFeeRange    []string `json:"latestPriorityFeeRange"`
    HistoricalPriorityFeeRange []string `json:"historicalPriorityFeeRange"`
    HistoricalBaseFeeRange    []string `json:"historicalBaseFeeRange"`
    PriorityFeeTrend          string   `json:"priorityFeeTrend"`
    BaseFeeTrend              string   `json:"baseFeeTrend"`
}
```

#### GasFeeLevel

```go
type GasFeeLevel struct {
    SuggestedMaxPriorityFeePerGas string `json:"suggestedMaxPriorityFeePerGas"`
    SuggestedMaxFeePerGas         string `json:"suggestedMaxFeePerGas"`
    MinWaitTimeEstimate           int64  `json:"minWaitTimeEstimate"`
    MaxWaitTimeEstimate           int64  `json:"maxWaitTimeEstimate"`
}
```

#### BaseFeeHistory

```go
// BaseFeeHistory is a type alias for []string
// The API directly returns an array of strings
type BaseFeeHistory []string
```

#### BaseFeePercentile

```go
type BaseFeePercentile struct {
    BaseFeePercentile string `json:"baseFeePercentile"`
}
```

#### BusyThreshold

```go
type BusyThreshold struct {
    BusyThreshold string `json:"busyThreshold"`
}
```

## 测试

运行测试：

```bash
go test ./...
```

运行测试并查看覆盖率：

```bash
go test -cover ./...
```

## 支持的链 ID

常见的链 ID：
- `1` - 以太坊主网 (Ethereum Mainnet)
- `5` - Goerli 测试网
- `137` - Polygon
- `42161` - Arbitrum One
- `10` - Optimism
- `8453` - Base

更多支持的链 ID 请参考 [Infura Gas API 文档](https://docs.metamask.io/services/reference/gas-api/api-reference/)。

## 参考文档

- [Infura Gas API 快速开始](https://docs.metamask.io/services/reference/gas-api/quickstart/)
- [Infura Gas API 参考文档](https://docs.metamask.io/services/reference/gas-api/api-reference/)

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

