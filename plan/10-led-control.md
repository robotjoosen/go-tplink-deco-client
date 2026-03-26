# Feature: LED Power Control

## Status
**Phase 2** - Fully discovered (from firmware analysis)

## Purpose
Control the LED indicator lights on Deco nodes - turn them on/off.

## API Details
| Property | Value |
|----------|-------|
| Endpoint | `/admin/wireless` |
| Form | `power` |
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
    "led_enable": true
  }
}
```

## Response Structure
```json
{
  "error_code": 0,
  "led_enable": true
}
```

## Data Structures

### Internal Model
```go
type LEDPowerResponse struct {
    ErrorCode  int  `json:"error_code"`
    LEDEnable  bool `json:"led_enable"`
}
```

### Exported Model
```go
type LEDPowerSettings struct {
    Enabled bool
}
```

## Notes
- Simple on/off control via `led_enable` parameter
- For schedule-based LED control, additional forms may exist

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
led, _ := client.GetLEDPower(ctx)
fmt.Printf("LED Enabled: %v\n", led.Enabled)
```
