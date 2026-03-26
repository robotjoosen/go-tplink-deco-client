# Feature: IGMP Settings

## Status
**Phase 2** - Fully discovered (from firmware analysis)

## Purpose
Get and configure IGMP (Internet Group Management Protocol) settings for multicast streaming.

## API Details
| Property | Value |
|----------|-------|
| Endpoint | `/admin/network` |
| Form | `igmp_setting` |
| Auth Required | Yes |
| Operations | `read`, `write` |

## Request (Read)
```json
{"operation": "read"}
```

## Response Structure
```json
{
  "error_code": 0,
  "igmp_enable": true
}
```

## Data Structures

### Internal Model
```go
type IGMPSettingResponse struct {
    ErrorCode   int  `json:"error_code"`
    IGMPEnable  bool `json:"igmp_enable"`
}
```

### Exported Model
```go
type IGMPSettings struct {
    Enabled bool
}
```

## Notes
- IGMP is used for IPTV/multicast streaming
- Common in ISP configurations requiring multicast

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
