# Feature: Log Export

## Status
**Phase 2** - Fully discovered (from firmware analysis)

## Purpose
Export and save logs from the Deco device.

## API Details
| Property | Value |
|----------|-------|
| Endpoint | `/admin/log_export` |
| Form | `save_log` |
| Auth Required | Yes |
| Operations | `read`, `write` |

## Request (Save Log)
```json
{"operation": "save_log"}
```

## Response Structure
```json
{
  "error_code": 0,
  "filename": "/tmp/log_export_20240101.tar.gz"
}
```

## Data Structures

### Internal Model
```go
type LogExportResponse struct {
    ErrorCode int    `json:"error_code"`
    Filename  string `json:"filename"`
}
```

### Exported Model
```go
type LogExportResult struct {
    Filename string
}
```

## Additional Forms
| Form | Description |
|------|-------------|
| `feedback_log` | Submit feedback logs |
| `types` | Get available log types |

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
