# Feature: IPv6 Settings

## Status
**Phase 2** - Fully discovered (from firmware analysis)

## Purpose
Get and configure IPv6 internet connection settings.

## API Details
| Property | Value |
|----------|-------|
| Endpoint | `/admin/network` |
| Form | `ipv` |
| Auth Required | Yes |
| Operations | `read`, `write` |

## Request
```json
{"operation": "read"}
```

## Response Structure (from firmware analysis)
```json
{
  "error_code": 0,
  "enable_ipv6": true,
  "wan": {
    "dial_type": "native",
    "ip": "...",
    "prefix": "...",
    "dns1": "...",
    "dns2": "..."
  },
  "lan": {...}
}
```

## Data Structures

### Internal Model
```go
type IPv6Response struct {
    ErrorCode   int         `json:"error_code"`
    EnableIPv6  bool        `json:"enable_ipv6"`
    WAN         IPv6WAN     `json:"wan"`
}

type IPv6WAN struct {
    DialType string `json:"dial_type"`
    IP       string `json:"ip"`
    Prefix   string `json:"prefix"`
    DNS1     string `json:"dns1"`
    DNS2     string `json:"dns2"`
}
```

### Exported Model
```go
type IPv6Settings struct {
    Enabled  bool
    DialType string // "native", "passthrough", "static"
    IP       net.IP
    Prefix   string
    DNS1     net.IP
    DNS2     net.IP
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
