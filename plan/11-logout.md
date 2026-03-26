# Feature: Logout

## Status
**Phase 2** - Fully discovered (from firmware analysis)

## Purpose
End the current authenticated session properly.

## API Details
| Property | Value |
|----------|-------|
| Endpoint | `/admin/system` |
| Form | `logout` |
| Auth Required | Yes |
| Operation | `logout` |

## Request
```json
{"operation": "logout"}
```

## Response Structure
```json
{
  "error_code": 0
}
```

## Data Structures

### Internal Model
```go
type LogoutResponse struct {
    ErrorCode int `json:"error_code"`
}
```

## Implementation Steps

### 1. Add internal model
**File:** `internal/model/model.go`

### 2. Update interface
**File:** `client.go`

### 3. Implement in httpclient
**File:** `internal/httpclient/client.go`

### 4. Add wrapper in main client
**File:** `client.go`

## Files Modified
- `internal/model/model.go`
- `internal/httpclient/client.go`
- `client.go`

## Notes
- Should be called when done using the client to clean up session
- Sets `authenticated = false` and clears `stok`

## Testing
```bash
err := client.Logout(ctx)
if err == nil {
    fmt.Println("Logged out successfully")
}
```
