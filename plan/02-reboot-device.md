# Feature: Reboot Device(s)

## Status
**Phase 1** - Ready to implement (ported from MrMarble/deco)

## API Details
| Property | Value |
|----------|-------|
| Endpoint | `/admin/device` |
| Form | `system` |
| Auth Required | Yes |
| Method | `POST` |
| Operation | `reboot` |

## Request
```json
{
  "operation": "reboot",
  "params": {
    "mac_list": [
      {"mac": "AA:BB:CC:DD:EE:FF"}
    ]
  }
}
```

## Response Model
```go
// Generic map response - API returns success indicator
type RebootResponse map[string]interface{}
```

## Implementation Steps

### 1. Add internal model
**File:** `internal/model/model.go`

Add `RebootRequest` struct:
```go
type RebootRequest struct {
    Operation string `json:"operation"`
    Params    struct {
        MACList []map[string]string `json:"mac_list"`
    } `json:"params"`
}

type RebootResponse struct {
    ErrorCode int `json:"error_code"`
    Result    map[string]interface{} `json:"result"`
}
```

### 2. Update interface
**File:** `client.go`

Add to `ClientAware` interface:
```go
RebootDevice(ctx context.Context, macAddrs []string) error
```

### 3. Implement in httpclient
**File:** `internal/httpclient/client.go`

Add method to `HTTPClient`:
```go
func (c *HTTPClient) RebootDevice(ctx context.Context, macAddrs []string) error {
    var macList []map[string]string
    for _, mac := range macAddrs {
        macList = append(macList, map[string]string{
            "mac": strings.ToUpper(mac),
        })
    }

    rebootReq := model.OperationRequest{
        Operation: "reboot",
        Params: map[string]interface{}{
            "mac_list": macList,
        },
    }

    jsonBody, err := json.Marshal(rebootReq)
    if err != nil {
        return err
    }

    args := url.Values{}
    args.Add("form", "system")

    var res model.RebootResponse
    return c.encryptPost(ctx, fmt.Sprintf(";stok=%s/admin/device", c.stok), args, jsonBody, false, &res)
}
```

### 4. Add wrapper in main client
**File:** `client.go`

Add public method with auth check:
```go
func (c *Client) RebootDevice(ctx context.Context, macAddrs ...string) error {
    if !c.authenticated {
        return errors.New("not authenticated")
    }

    return c.client.RebootDevice(ctx, macAddrs)
}
```

## Files Modified
- `internal/model/model.go`
- `internal/httpclient/client.go`
- `client.go`

## Error Handling
| Error | Cause |
|-------|-------|
| `not authenticated` | Client not authenticated |
| API returns error_code != 0 | Reboot command failed |

## Testing
```bash
# Reboot single device
client.RebootDevice(ctx, "AA:BB:CC:DD:EE:FF")

# Reboot multiple devices
client.RebootDevice(ctx, "AA:BB:CC:DD:EE:FF", "11:22:33:44:55:66")
```
