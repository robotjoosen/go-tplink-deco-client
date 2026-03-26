# Feature: Performance Monitoring

## Status
**Phase 1** - Ready to implement (ported from MrMarble/deco)

## API Details
| Property | Value |
|----------|-------|
| Endpoint | `/admin/network` |
| Form | `performance` |
| Auth Required | Yes |
| Method | `POST` |
| Operation | `read` |

## Request
```json
{"operation": "read"}
```

## Response Model
```go
type PerformanceResponse struct {
    ErrorCode int `json:"error_code"`
    Result    struct {
        CPU float32 `json:"cpu_usage"`
        MEM float32 `json:"mem_usage"`
    } `json:"result"`
}
```

## Implementation Steps

### 1. Add internal model
**File:** `internal/model/model.go`

Add `PerformanceResponse` struct:
```go
type PerformanceResponse struct {
    ErrorCode int `json:"error_code"`
    Result    struct {
        CPU float32 `json:"cpu_usage"`
        MEM float32 `json:"mem_usage"`
    } `json:"result"`
}
```

### 2. Add exported model
**File:** `pkg/model/model.go`

Add exported `Performance` struct:
```go
type Performance struct {
    CPU float32 `json:"cpu"`
    MEM float32 `json:"mem"`
}
```

### 3. Update interface
**File:** `client.go`

Add to `ClientAware` interface:
```go
GetPerformance(ctx context.Context) (model.PerformanceResponse, error)
```

### 4. Implement in httpclient
**File:** `internal/httpclient/client.go`

Add method to `HTTPClient`:
```go
func (c *HTTPClient) GetPerformance(ctx context.Context) (model.PerformanceResponse, error) {
    var res model.PerformanceResponse
    readBody, err := json.Marshal(model.OperationRequest{Operation: "read"})
    if err != nil {
        return res, err
    }

    args := url.Values{}
    args.Add("form", "performance")

    err = c.encryptPost(ctx, fmt.Sprintf(";stok=%s/admin/network", c.stok), args, readBody, false, &res)
    return res, err
}
```

### 5. Add wrapper in main client
**File:** `client.go`

Add public method with auth check and data transformation:
```go
func (c *Client) GetPerformance(ctx context.Context) (exportModel.Performance, error) {
    if !c.authenticated {
        return exportModel.Performance{}, errors.New("not authenticated")
    }

    res, err := c.client.GetPerformance(ctx)
    if err != nil {
        return exportModel.Performance{}, err
    }

    return exportModel.Performance{
        CPU: res.Result.CPU,
        MEM: res.Result.MEM,
    }, nil
}
```

## Files Modified
- `internal/model/model.go`
- `pkg/model/model.go`
- `internal/httpclient/client.go`
- `client.go`

## Testing
```bash
# Manual test
client.Authenticate(ctx, password)
perf, _ := client.GetPerformance(ctx)
fmt.Printf("CPU: %.2f%%, MEM: %.2f%%\n", perf.CPU, perf.MEM)
```
