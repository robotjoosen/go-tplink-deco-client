# Feature: DHCP Dial

## Status
**Phase 2** - Fully discovered (from firmware analysis)

## Purpose
Get and configure DHCP dial settings for WAN connection.

## API Details
| Property | Value |
|----------|-------|
| Endpoint | `/admin/network` |
| Form | `dhcp_dial` |
| Auth Required | Yes |
| Operations | `read`, `write` |

## Request (Read)
```json
{"operation": "read"}
```

## Request (Write)
```json
{
  "operation": "write",
  "params": {
    "dial_type": "DHCP",
    "username": "",
    "password": ""
  }
}
```

## Response Structure
```json
{
  "error_code": 0,
  "dial_type": "DHCP",
  "username": "",
  "password": "",
  "host_name": "Deco",
  "service_name": ""
}
```

## Data Structures

### Internal Model
```go
type DHCPDialResponse struct {
    ErrorCode     int    `json:"error_code"`
    DialType      string `json:"dial_type"`
    Username      string `json:"username"`
    Password      string `json:"password"`
    HostName      string `json:"host_name"`
    ServiceName   string `json:"service_name"`
}
```

### Exported Model
```go
type DHCPDialSettings struct {
    DialType    string // "DHCP", "PPPoE", "Static", "L2TP", "PPTP"
    Username    string
    Password    string
    HostName    string
    ServiceName string
}
```

## Supported Dial Types
| Type | Description |
|------|-------------|
| `DHCP` | Dynamic IP |
| `PPPoE` | PPPoE dial-up |
| `Static` | Static IP |
| `L2TP` | L2TP VPN |
| `PPTP` | PPTP VPN |

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
