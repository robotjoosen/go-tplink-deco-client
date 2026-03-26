# Feature: Custom Request

## Status
**Phase 1** - Ready to implement (ported from MrMarble/deco)

## Purpose
Provides a generic escape hatch for making arbitrary requests to the Deco API, enabling exploration of undocumented features and future extensibility without code changes.

## API Details
| Property | Value |
|----------|-------|
| Endpoint | `{stok}{path}` |
| Form | User-provided |
| Auth Required | Yes |
| Method | `POST` |

## Signature
```go
func (c *Client) Custom(ctx context.Context, path string, form string, operation string, params map[string]interface{}) (interface{}, error)
```

## Implementation Steps

### 1. Update interface
**File:** `client.go`

Add to `ClientAware` interface:
```go
Custom(ctx context.Context, path string, form string, body []byte) (interface{}, error)
```

### 2. Implement in httpclient
**File:** `internal/httpclient/client.go`

Add method to `HTTPClient`:
```go
func (c *HTTPClient) Custom(ctx context.Context, path string, form string, body []byte) (interface{}, error) {
    var result interface{}
    
    args := url.Values{}
    args.Add("form", form)

    err := c.encryptPost(ctx, fmt.Sprintf(";stok=%s%s", c.stok, path), args, body, false, &result)
    return result, err
}
```

### 3. Add wrapper in main client
**File:** `client.go`

Add public method:
```go
func (c *Client) Custom(ctx context.Context, path string, form string, operation string, params map[string]interface{}) (interface{}, error) {
    if !c.authenticated {
        return nil, errors.New("not authenticated")
    }

    req := model.OperationRequest{
        Operation: operation,
        Params:    params,
    }

    jsonBody, err := json.Marshal(req)
    if err != nil {
        return nil, err
    }

    return c.client.Custom(ctx, path, form, jsonBody)
}
```

## Usage Examples

### Discover available forms
```go
// Try common forms with read operation
result, err := client.Custom(ctx, "/admin/network", "wireless", "read", nil)
```

### Check system logs
```go
result, err := client.Custom(ctx, "/admin/log", "log", "read", nil)
```

### Test parental controls
```go
result, err := client.Custom(ctx, "/admin/parental", "parental", "read", map[string]interface{}{
    "device_mac": "default",
})
```

## Files Modified
- `internal/httpclient/client.go`
- `client.go`

## Notes
- This is a low-level escape hatch; most users should use typed methods
- Response is `interface{}` - caller must type assert
- Use for exploration and prototyping new features
- Once a feature is stable, consider adding a typed method
