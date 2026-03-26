# Feature: WiFi Settings

## Status
**Phase 2** - Fully discovered (from firmware analysis)

## Purpose
Read and modify WiFi settings including SSID, password, and enable/disable per band (2.4GHz/5GHz/6GHz for host/guest/IoT).

## API Details
| Property | Value |
|----------|-------|
| Endpoint | `/admin/wireless` |
| Form | `wlan` |
| Auth Required | Yes |
| Operations | `read`, `write` |

## Request Structure (Read)
```json
{"operation": "read"}
```

## Response Structure (from firmware analysis)
```json
{
  "error_code": 0,
  "band2_4": {
    "host": {"enable": true, "ssid": "...", "key": "..."},
    "guest": {"enable": false, "ssid": "...", "key": "..."},
    "iot": {"enable": false}
  },
  "band5_1": {
    "host": {"enable": true, "ssid": "...", "key": "..."},
    "guest": {"enable": false, "ssid": "...", "key": "..."},
    "iot": {"enable": false}
  },
  "band6": {
    "host": {"enable": true, "ssid": "...", "key": "..."},
    "guest": {"enable": false},
    "iot": {"enable": false}
  }
}
```

## Request Structure (Write)
```json
{
  "operation": "write",
  "params": {
    "band2_4": {"host": {"enable": false}},
    "band5_1": {"host": {"enable": true}}
  }
}
```

## Data Structures

### Internal Model
```go
type WiFiResponse struct {
    ErrorCode int           `json:"error_code"`
    Band24    WiFiBand      `json:"band2_4"`
    Band5     WiFiBand      `json:"band5_1"`
    Band6     WiFiBand      `json:"band6"`
}

type WiFiBand struct {
    Host  WiFiNetwork `json:"host"`
    Guest WiFiNetwork `json:"guest"`
    IoT   WiFiNetwork `json:"iot"`
}

type WiFiNetwork struct {
    Enable bool   `json:"enable"`
    SSID   string `json:"ssid,omitempty"`
    Key    string `json:"key,omitempty"`
}
```

### Exported Model
```go
type WiFiSettings struct {
    Band24 WiFiBandSettings `json:"band24"`
    Band5  WiFiBandSettings `json:"band5"`
    Band6  WiFiBandSettings `json:"band6"`
}

type WiFiBandSettings struct {
    Host  WiFiNetworkSettings
    Guest WiFiNetworkSettings
    IoT   WiFiNetworkSettings
}

type WiFiNetworkSettings struct {
    Enabled  bool
    SSID     string
    Password string
}
```

### Connection Types
| Type | Description |
|------|-------------|
| `host` | Main network |
| `guest` | Guest network |
| `iot` | IoT network |

## Capabilities
| Method | Description |
|--------|-------------|
| `GetWiFiSettings(ctx)` | Get all WiFi settings |
| `SetWiFiEnabled(ctx, band, network, enabled)` | Enable/disable WiFi type |
| `SetWiFiSettings(ctx, settings)` | Update full WiFi configuration |

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

## Notes
- 6GHz band support depends on hardware (Wi-Fi 6E/7 devices like Deco X50)
- IoT network is for smart home devices
- Guest network supports portal authentication (firmware v1.4.4+)
