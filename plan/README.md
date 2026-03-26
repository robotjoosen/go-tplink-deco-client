# Implementation Plans Index

## Phase 1: Ready to Implement (Ported from MrMarble/deco)

| # | Feature | File | Status |
|---|---------|------|--------|
| 1 | Performance Monitoring | [01-performance-monitoring.md](01-performance-monitoring.md) | Ready |
| 2 | Reboot Device(s) | [02-reboot-device.md](02-reboot-device.md) | Ready |
| 3 | Custom Request | [03-custom-request.md](03-custom-request.md) | Ready |

## Phase 2: Fully Discovered (From Firmware Analysis)

**Source:** Deco X50 v3.0 Firmware 1.5.0 (20250619)

| # | Feature | File | Endpoint | Status |
|---|---------|------|----------|--------|
| 4 | WAN IPv4 Info | [04-wan-information.md](04-wan-information.md) | `/admin/network?form=wan_ipv` | Discovered |
| 5 | WiFi Settings | [05-wifi-settings.md](05-wifi-settings.md) | `/admin/wireless?form=wlan` | Discovered |
| 6 | Internet Status | [06-internet-status.md](06-internet-status.md) | `/admin/network?form=internet` | Discovered |
| 7 | IPv6 Settings | [07-ipv6-settings.md](07-ipv6-settings.md) | `/admin/network?form=ipv` | Discovered |
| 8 | LAN IP Settings | [08-lan-ip-settings.md](08-lan-ip-settings.md) | `/admin/network?form=lan_ip` | Discovered |
| 9 | DHCP Dial | [09-dhcp-dial.md](09-dhcp-dial.md) | `/admin/network?form=dhcp_dial` | Discovered |
| 10 | LED Control | [10-led-control.md](10-led-control.md) | `/admin/wireless?form=power` | Discovered |
| 11 | Logout | [11-logout.md](11-logout.md) | `/admin/system?form=logout` | Discovered |
| 12 | System Logs | [12-system-logs.md](12-system-logs.md) | `/admin/syslog?form=log` | Discovered |
| 13 | Log Export | [13-log-export.md](13-log-export.md) | `/admin/log_export?form=save_log` | Discovered |
| 14 | IGMP Settings | [14-igmp-setting.md](14-igmp-setting.md) | `/admin/network?form=igmp_setting` | Discovered |
| 15 | Fast Transmit | [15-fast-xmit-setting.md](15-fast-xmit-setting.md) | `/admin/network?form=fast_xmit_setting` | Discovered |
| 16 | QoS / Bandwidth | [16-qos-bandwidth.md](16-qos-bandwidth.md) | `/admin/network?form=flow_control` | Discovered |
| 17 | Guest Network | (part of WiFi Settings) | `/admin/wireless?form=wlan` | Discovered |

## Phase 3: Requires Further Research

| # | Feature | Notes |
|---|---------|-------|
| - | Port Forwarding | Needs discovery |
| - | Parental Controls | Needs discovery |
| - | VPN Server/Client | Needs discovery |
| - | Wi-Fi Access Control | Needs discovery |
| - | Internet Backup | Needs discovery |
| - | Eco Mode / Wi-Fi Schedule | Needs discovery |
| - | Auto Firmware Update | `/admin/cloud?form=firmware_status` - partial |
| - | Firmware Upgrade | `/admin/firmware?form=upgrade` - partial |

## Complete API Endpoints Discovered (Firmware 1.5.0)

### Auth
- `POST /admin/system?form=logout` - End session

### Device
- `POST /admin/device?form=device_list` - List Deco nodes
- `POST /admin/device?form=system` - Reboot/system operations
- `POST /admin/device?form=timesetting` - Time settings

### Client
- `POST /admin/client?form=client_list` - List connected clients

### Network
- `POST /admin/network?form=internet` - Internet connection status
- `POST /admin/network?form=wan_ipv` - WAN IPv4 info
- `POST /admin/network?form=lan_ip` - LAN IP settings
- `POST /admin/network?form=dhcp_dial` - DHCP configuration
- `POST /admin/network?form=performance` - CPU/Memory usage
- `POST /admin/network?form=ipv` - IPv6 settings
- `POST /admin/network?form=igmp_setting` - IGMP settings
- `POST /admin/network?form=flow_control` - Flow/QoS control
- `POST /admin/network?form=flow_control_lan_wan` - LAN/WAN flow control
- `POST /admin/network?form=erp_setting` - ERP settings
- `POST /admin/network?form=dsl_status` - DSL status
- `POST /admin/network?form=fast_xmit_setting` - Fast transmit settings
- `POST /admin/network?form=wifi_network` - WiFi network settings

### Wireless
- `POST /admin/wireless?form=wlan` - WiFi settings (SSID, security, enable/disable)
- `POST /admin/wireless?form=power` - LED power settings

### Firmware
- `POST /admin/firmware?form=upgrade` - Firmware upgrade
- `POST /admin/firmware?form=config` - Firmware config
- `POST /admin/firmware?form=config_multipart` - Firmware multipart config

### Cloud
- `POST /admin/cloud?form=firmware` - Cloud firmware operations
- `POST /admin/cloud?form=firmware_status` - Firmware upgrade status

### System
- `POST /admin/syslog?form=log` - System logs
- `POST /admin/syslog?form=mail` - Mail settings

### Log
- `POST /admin/log_export?form=save_log` - Export logs
- `POST /admin/log_export?form=feedback_log` - Feedback logs
- `POST /admin/log_export?form=types` - Log types

### Other
- `POST /admin/folder_sharing?form=tree` - Network map tree
- `POST /admin/isp?form=isp_upgrade` - ISP upgrade
- `POST /admin/cwmp?form=cwmp_info` - CWMP info
- `POST /admin/cloud_account?form=get_device` - Get device info
- `POST /admin/cloud_account?form=check_internet` - Check internet

## Discovery Methods Used

1. **Firmware Analysis** - Extracted SquashFS from Deco X50 v3.0 firmware 1.5.0
2. **Web UI JavaScript** - Found all API endpoints in `www/webpages/modules/*/models.js`

## Implementation Pattern

All features follow this pattern:

1. Add internal model in `internal/model/model.go`
2. Add exported model in `pkg/model/model.go`
3. Add method to `ClientAware` interface
4. Implement in `internal/httpclient/client.go`
5. Add wrapper with auth check in `client.go`

## Recommended Implementation Order

1. **Phase 1:** Custom Request (03) - enables all other features
2. **Phase 1:** Performance (01) - simple read operation
3. **Phase 1:** Reboot (02) - simple write operation
4. **Phase 2:** WAN IPv4 Info (04) - high value, well understood
5. **Phase 2:** WiFi Settings (05) - high value, complex
6. **Phase 2:** Internet Status (06) - complements WAN
7. **Phase 2:** LAN IP Settings (08) - common configuration
8. **Phase 2:** Logout (11) - proper session management
9. **Phase 2:** Continue with remaining features
