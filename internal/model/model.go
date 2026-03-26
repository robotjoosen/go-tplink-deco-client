package model

type OperationRequest struct {
	Operation string                 `json:"operation,omitempty"`
	Params    map[string]interface{} `json:"params,omitempty"`
}

type ErrorResponse struct {
	ErrorCode int `json:"error_code"`
}

type ClientListResponse struct {
	Result struct {
		ClientList []struct {
			Mac            string `json:"mac"`
			UpSpeed        int    `json:"up_speed"`
			DownSpeed      int    `json:"down_speed"`
			WireType       string `json:"wire_type"`
			AccessHost     string `json:"access_host"`
			ConnectionType string `json:"connection_type"`
			SpaceID        string `json:"space_id"`
			IP             string `json:"ip"`
			ClientMesh     bool   `json:"client_mesh"`
			Online         bool   `json:"online"`
			Name           string `json:"name"`
			EnablePriority bool   `json:"enable_priority"`
			RemainTime     int    `json:"remain_time"`
			OwnerID        string `json:"owner_id"`
			ClientType     string `json:"client_type"`
			Interface      string `json:"interface"`
		} `json:"client_list"`
	} `json:"result"`
	ErrorCode int `json:"error_code"`
}

type DeviceListResponse struct {
	Result struct {
		DeviceList []struct {
			NandFlash         bool     `json:"nand_flash"`
			OwnerTransfer     bool     `json:"owner_transfer,omitempty"`
			Previous          string   `json:"previous"`
			ParentDeviceID    string   `json:"parent_device_id,omitempty"`
			Role              string   `json:"role"`
			BssidSta5G        string   `json:"bssid_sta_5g"`
			SupportPLC        bool     `json:"support_plc"`
			SetGatewaySupport bool     `json:"set_gateway_support"`
			GroupStatus       string   `json:"group_status"`
			WiredPortList     struct{} `json:"wired_port_list,omitempty"`
			PortCount         int      `json:"port_count,omitempty"`
			SignalLevel       struct {
				Band24 string `json:"band2_4"`
				Band5  string `json:"band5"`
			} `json:"signal_level"`
			DeviceModel       string `json:"device_model"`
			Bssid5G           string `json:"bssid_5g"`
			SpeedGetSupport   bool   `json:"speed_get_support,omitempty"`
			DeviceType        string `json:"device_type"`
			HardwareVer       string `json:"hardware_ver"`
			Bssid2G           string `json:"bssid_2g"`
			InetStatus        string `json:"inet_status"`
			Mac               string `json:"mac"`
			InetErrorMsg      string `json:"inet_error_msg"`
			DeviceIP          string `json:"device_ip"`
			SoftwareVersion   string `json:"software_ver"`
			DeviceID          string `json:"device_id,omitempty"`
			BssidSta2G        string `json:"bssid_sta_2g"`
			OemID             string `json:"oem_id"`
			Nickname          string `json:"nickname"`
			ProductLevel      int    `json:"product_level"`
			OversizedFirmware bool   `json:"oversized_firmware"`
			Topology          struct {
				Auto     bool   `json:"auto"`
				DeviceId string `json:"device_id"`
			} `json:"topology,omitempty"`
			HwId           string   `json:"hw_id"`
			ConnectionType []string `json:"connection_type,omitempty"`
		} `json:"device_list"`
	} `json:"result"`
	ErrorCode int `json:"error_code"`
}

type LoginResponse struct {
	Result struct {
		Stok string `json:"stok"`
	} `json:"result"`
	ErrCode int `json:"error_code"`
}

type LoginKeyResponse struct {
	Result struct {
		Username string   `json:"username"`
		Password []string `json:"password"`
	} `json:"result"`
	ErrCode int `json:"error_code"`
}

type SessionKeyResponse struct {
	Result struct {
		Key []string `json:"key"`
		Seq uint     `json:"seq"`
	} `json:"result"`
	ErrorCode int `json:"error_code"`
}

type PerformanceResponse struct {
	ErrorCode int `json:"error_code"`
	Result    struct {
		CPU float32 `json:"cpu_usage"`
		MEM float32 `json:"mem_usage"`
	} `json:"result"`
}

type RebootRequest struct {
	Operation string `json:"operation"`
	Params    struct {
		MACList []map[string]string `json:"mac_list"`
	} `json:"params"`
}

type RebootResponse struct {
	ErrorCode int                    `json:"error_code"`
	Result    map[string]interface{} `json:"result"`
}

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

type WiFiResponse struct {
	ErrorCode int      `json:"error_code"`
	Band24    WiFiBand `json:"band2_4"`
	Band5     WiFiBand `json:"band5_1"`
	Band6     WiFiBand `json:"band6"`
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

type InternetStatusResponse struct {
	ErrorCode    int    `json:"error_code"`
	InetStatus   string `json:"inet_status"`
	InetErrorMsg string `json:"inet_error_msg"`
	Speed        int    `json:"speed"`
	Duplex       int    `json:"duplex"`
}

type IPv6Response struct {
	ErrorCode  int     `json:"error_code"`
	EnableIPv6 bool    `json:"enable_ipv6"`
	WAN        IPv6WAN `json:"wan"`
}

type IPv6WAN struct {
	DialType string `json:"dial_type"`
	IP       string `json:"ip"`
	Prefix   string `json:"prefix"`
	DNS1     string `json:"dns1"`
	DNS2     string `json:"dns2"`
}

type LANIPResponse struct {
	ErrorCode int         `json:"error_code"`
	LAN       LANIPConfig `json:"lan"`
}

type LANIPConfig struct {
	IP         string `json:"ip"`
	Mask       string `json:"mask"`
	DHCPEnable bool   `json:"dhcp_enable"`
	DHCPType   string `json:"dhcp_type"`
}
