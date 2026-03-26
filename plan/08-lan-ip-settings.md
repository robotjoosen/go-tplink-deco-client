# Feature: LAN IP Settings

## Status
**Phase 2** - Fully discovered (from firmware analysis)

## Purpose
Get and configure LAN IP settings for the Deco device.

## API Details
| Property | Value |
|----------|-------|
| Endpoint | `/admin/network` |
| Form | `lan_ip` |
| Auth Required | Yes |
| Operations | `read`, `write` |

## Request
```json
{"operation": "read"}
```

## Response Structure
```json
{
  "error_code": 0,
  "lan": {
    "ip": "192.168.0.1",
    "mask": "255.255.255.0",
    "dhcp_enable": true,
    "dhcp_type": "server"
  }
}
```

## Data Structures

### Internal Model
```go
type LANIPResponse struct {
    ErrorCode int    `json:"error_code"`
    LAN       LANIPConfig `json:"lan"`
}

type LANIPConfig struct {
    IP         string `json:"ip"`
    Mask       string `json:"mask"`
    DHCPEnable bool   `json:"dhcp_enable"`
    DHCPType   string `json:"dhcp_type"` // "server", "relay"
}
```

### Exported Model
```go
type LANIPSettings struct {
    IP         net.IP
    SubnetMask net.IPMask
    DHCPEnable bool
    DHCPType   string
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
