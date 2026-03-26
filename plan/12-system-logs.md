# Feature: System Logs

## Status
**Phase 2** - Fully discovered (from firmware analysis)

## Purpose
Retrieve system logs from the Deco device.

## API Details
| Property | Value |
|----------|-------|
| Endpoint | `/admin/syslog` |
| Form | `log` |
| Auth Required | Yes |
| Operation | `read` |

## Request
```json
{"operation": "read"}
```

## Response Structure
```json
{
  "error_code": 0,
  "log": "syslog content here..."
}
```

## Data Structures

### Internal Model
```go
type SysLogResponse struct {
    ErrorCode int    `json:"error_code"`
    Log       string `json:"log"`
}
```

### Exported Model
```go
type SystemLog struct {
    Content string
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
log, _ := client.GetSystemLog(ctx)
fmt.Println(log.Content)
```
