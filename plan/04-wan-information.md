# Feature: WAN IPv4 Information

## Status
**Phase 2** - Fully discovered (from firmware analysis)

## Purpose
Get Wide Area Network (WAN) IPv4 information including IP address, gateway, DNS, MAC, and connection type.

## API Details
| Property | Value |
|----------|-------|
| Endpoint | `/admin/network` |
| Form | `wan_ipv` |
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
  "wan": {
    "ip_info": {
      "ip": "192.168.1.100",
      "mask": "255.255.255.0",
      "gateway": "192.168.1.1",
      "mac": "AA:BB:CC:DD:EE:FF",
      "dns1": "8.8.8.8",
      "dns2": "8.8.4.4"
    },
    "dial_type": "DHCP",
    "link_status": "up"
  },
  "lan": {
    "ip_info": {
      "ip": "192.168.0.1",
      "mask": "255.255.255.0",
      "mac": "AA:BB:CC:DD:EE:FF"
    }
  }
}
```

## Data Structures

### Internal Model
```go
type WANIPv4Response struct {
    ErrorCode int        `json:"error_code"`
    WAN       WANSection `json:"wan"`
    LAN       LANSection `json:"lan"`
}

type WANSection struct {
    IPInfo     WANIPInfo `json:"ip_info"`
    DialType   string    `json:"dial_type"`
    LinkStatus string    `json:"link_status,omitempty"`
}

type WANIPInfo struct {
    IP      string `json:"ip"`
    Mask    string `json:"mask"`
    Gateway string `json:"gateway"`
    MAC     string `json:"mac"`
    DNS1    string `json:"dns1"`
    DNS2    string `json:"dns2"`
}

type LANSection struct {
    IPInfo LANIPInfo `json:"ip_info"`
}

type LANIPInfo struct {
    IP   string `json:"ip"`
    Mask string `json:"mask"`
    MAC  string `json:"mac"`
}
```

### Exported Model
```go
type WANInfo struct {
    WANMAC     net.HardwareAddr
    WANIP      net.IP
    WANGateway net.IP
    WANDNS1    net.IP
    WANDNS2    net.IP
    WANSubnet  net.IPMask
    WANType    string // DHCP, PPPoE, Static
    LinkStatus string
    LANMAC     net.HardwareAddr
    LANIP      net.IP
    LANSubnet  net.IPMask
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
wan, _ := client.GetWANIPv4Info(ctx)
fmt.Printf("WAN IP: %s, Gateway: %s, Type: %s\n", wan.WANIP, wan.WANGateway, wan.WANType)
```
