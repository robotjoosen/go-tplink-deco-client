# Feature: Internet Status

## Status
**Phase 2** - Fully discovered (from firmware analysis)

## Purpose
Get internet connection status and diagnostics.

## API Details
| Property | Value |
|----------|-------|
| Endpoint | `/admin/network` |
| Form | `internet` |
| Auth Required | Yes |
| Operation | `read` |

## Request
```json
{"operation": "read"}
```

## Response Structure (from firmware analysis)
```json
{
  "error_code": 0,
  "inet_status": "up",
  "inet_error_msg": "",
  "speed": 1000,
  "duplex": 1,
  "WAN": {...}
}
```

## Data Structures

### Internal Model
```go
type InternetStatusResponse struct {
    ErrorCode     int    `json:"error_code"`
    InetStatus    string `json:"inet_status"`
    InetErrorMsg  string `json:"inet_error_msg"`
    Speed         int    `json:"speed"`
    Duplex        int    `json:"duplex"`
}
```

### Exported Model
```go
type InternetStatus struct {
    Status    string // "up", "down"
    ErrorMsg  string
    Speed     int    // Mbps
    Duplex    int    // 1=full, 0=half
}
```

## Implementation Steps

### 1. Add internal model
**File:** `internal/model/model.go`

### 2. Add exported model
**File:** `pkg/model/model.go`

### 3. Update interface
**File:** `client.go`

### 4. Implement in httpclient
**File:** `internal/httpclient/client.go`

### 5. Add wrapper in main client
**File:** `client.go`

## Files Modified
- `internal/model/model.go`
- `pkg/model/model.go`
- `internal/httpclient/client.go`
- `client.go`

## Testing
```bash
status, _ := client.GetInternetStatus(ctx)
fmt.Printf("Internet: %s, Speed: %d Mbps\n", status.Status, status.Speed)
```
