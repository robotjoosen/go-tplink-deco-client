# Feature: QoS / Bandwidth Control

## Status
**Phase 2** - Fully discovered (from firmware analysis)

## Purpose
Quality of Service settings - control flow and bandwidth management.

## API Details
| Property | Value |
|----------|-------|
| Endpoint | `/admin/network` |
| Forms | `flow_control`, `flow_control_lan_wan` |
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
  "flow_control_enable": true,
  "upload_bandwidth": 1000,
  "download_bandwidth": 1000
}
```

## Data Structures

### Internal Model
```go
type FlowControlResponse struct {
    ErrorCode          int  `json:"error_code"`
    FlowControlEnable  bool `json:"flow_control_enable"`
    UploadBandwidth    int  `json:"upload_bandwidth"`
    DownloadBandwidth  int  `json:"download_bandwidth"`
}

type FlowControlLANWANResponse struct {
    ErrorCode   int  `json:"error_code"`
    LANEnable   bool `json:"lan_enable"`
    WANEnable   bool `json:"wan_enable"`
}
```

### Exported Model
```go
type QoSSettings struct {
    Enabled          bool
    UploadBandwidth  int // Mbps
    DownloadBandwidth int // Mbps
}

type FlowControlLANWAN struct {
    LANEnabled bool
    WANEnabled bool
}
```

## Notes
- `flow_control` - Main QoS settings (bandwidth limits)
- `flow_control_lan_wan` - LAN/WAN specific flow control

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
