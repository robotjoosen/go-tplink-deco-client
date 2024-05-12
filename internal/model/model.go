package model

type OperationRequest struct {
	Operation string `json:"operation"`
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
			SpaceId        string `json:"space_id"`
			Ip             string `json:"ip"`
			ClientMesh     bool   `json:"client_mesh"`
			Online         bool   `json:"online"`
			Name           string `json:"name"`
			EnablePriority bool   `json:"enable_priority"`
			RemainTime     int    `json:"remain_time"`
			OwnerId        string `json:"owner_id"`
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
			ParentDeviceId    string   `json:"parent_device_id,omitempty"`
			Role              string   `json:"role"`
			BssidSta5G        string   `json:"bssid_sta_5g"`
			SupportPlc        bool     `json:"support_plc"`
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
			DeviceIp          string `json:"device_ip"`
			SoftwareVer       string `json:"software_ver"`
			DeviceId          string `json:"device_id,omitempty"`
			BssidSta2G        string `json:"bssid_sta_2g"`
			OemId             string `json:"oem_id"`
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

type LoginKeyResponse struct {
	Result struct {
		Username string   `json:"username"`
		Password []string `json:"password"`
	} `json:"result"`
	ErrorCode int `json:"error_code"`
}

type LoginAuthResponse struct {
	Result struct {
		Key []string `json:"key"`
		Seq int      `json:"seq"`
	} `json:"result"`
	ErrorCode int `json:"error_code"`
}
