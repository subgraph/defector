package networkmanager

import (
	"fmt"
	"encoding/binary"
	"net"

	"github.com/godbus/dbus"
)

const (
	NM_DEVICE_STATE_UNKNOWN = 0
	NM_DEVICE_STATE_UNAVAILABLE = 20
	NM_DEVICE_STATE_DISCONNECTED = 30
	NM_DEVICE_STATE_PREPARE = 40
	NM_DEVICE_STATE_CONFIG = 50
	NM_DEVICE_STATE_NEED_AUTH = 60
	NM_DEVICE_STATE_IP_CONFIG = 70
	NM_DEVICE_STATE_SECONDARIES = 90
	NM_DEVICE_STATE_ACTIVATED = 100
	NM_DEVICE_STATE_DEACTIVATING = 110
	NM_DEVICE_STATE_FAILED = 120
)

const (
	NM_DEVICE_TYPE_ETHERNET = 1
	NM_DEVICE_TYPE_WIFI = 2
	NM_DEVICE_TYPE_UNUSED1 = 3
	NM_DEVICE_TYPE_UNUSED2 = 4
	NM_DEVICE_TYPE_BT = 5
	NM_DEVICE_TYPE_OLPC_MESH = 6
	NM_DEVICE_TYPE_WIMAX = 7
	NM_DEVICE_TYPE_MODEM = 8
	NM_DEVICE_TYPE_INFINIBAND = 9
	NM_DEVICE_TYPE_BOND = 10
	NM_DEVICE_TYPE_VLAN = 11
	NM_DEVICE_TYPE_ADSL = 12
	NM_DEVICE_TYPE_BRIDGE = 13
	NM_DEVICE_TYPE_GENERIC = 14
	NM_DEVICE_TYPE_TEAM = 15
)

const (
	NM_DEVICE_STATE_REASON_NONE = 1
	NM_DEVICE_STATE_REASON_NOW_MANAGED = 2
	NM_DEVICE_STATE_REASON_NOW_UNMANAGED = 3
	NM_DEVICE_STATE_REASON_CONFIG_FAILED = 4
	NM_DEVICE_STATE_REASON_CONFIG_UNAVAILABLE = 5
	NM_DEVICE_STATE_REASON_CONFIG_EXPIRED = 6
	NM_DEVICE_STATE_REASON_NO_SECRETS = 7
	NM_DEVICE_STATE_REASON_SUPPLICANT_DISCONNECT = 8
	NM_DEVICE_STATE_REASON_SUPPLICANT_CONFIG_FAILED = 9
	NM_DEVICE_STATE_REASON_SUPPLICANT_FAILED = 10
	NM_DEVICE_STATE_REASON_SUPPLICANT_TIMEOUT = 11
	NM_DEVICE_STATE_REASON_PPP_START_FAILED = 12
	NM_DEVICE_STATE_REASON_PPP_DISCONNECT = 13
	NM_DEVICE_STATE_REASON_PPP_FAILED = 14
	NM_DEVICE_STATE_REASON_DHCP_START_FAILED = 15
	NM_DEVICE_STATE_REASON_DHCP_ERROR = 16
	NM_DEVICE_STATE_REASON_DHCP_FAILED = 17
	NM_DEVICE_STATE_REASON_SHARED_START_FAILED = 18
	NM_DEVICE_STATE_REASON_SHARED_FAILED = 19
	NM_DEVICE_STATE_REASON_AUTOIP_START_FAILED = 20
	NM_DEVICE_STATE_REASON_AUTOIP_ERROR = 21
	NM_DEVICE_STATE_REASON_AUTOIP_FAILED = 22
	NM_DEVICE_STATE_REASON_MODEM_BUSY = 23
	NM_DEVICE_STATE_REASON_MODEM_NO_DIAL_TONE = 24
	NM_DEVICE_STATE_REASON_MODEM_NO_CARRIER = 25
	NM_DEVICE_STATE_REASON_MODEM_DIAL_TIMEOUT = 26
	NM_DEVICE_STATE_REASON_MODEM_DIAL_FAILED = 27
	NM_DEVICE_STATE_REASON_MODEM_INIT_FAILED = 28
	NM_DEVICE_STATE_REASON_GSM_APN_FAILED = 29
	NM_DEVICE_STATE_REASON_GSM_REGISTRATION_NOT_SEARCHING = 30
	NM_DEVICE_STATE_REASON_GSM_REGISTRATION_DENIED = 31
	NM_DEVICE_STATE_REASON_GSM_REGISTRATION_TIMEOUT = 32
	NM_DEVICE_STATE_REASON_GSM_REGISTRATION_FAILED = 33
	NM_DEVICE_STATE_REASON_GSM_PIN_CHECK_FAILED = 34
	NM_DEVICE_STATE_REASON_FIRMWARE_MISSING = 35
	NM_DEVICE_STATE_REASON_REMOVED = 36
	NM_DEVICE_STATE_REASON_SLEEPING = 37
	NM_DEVICE_STATE_REASON_CONNECTION_REMOVED = 38
	NM_DEVICE_STATE_REASON_USER_REQUESTED = 39
	NM_DEVICE_STATE_REASON_CARRIER = 40
	NM_DEVICE_STATE_REASON_CONNECTION_ASSUMED = 41
	NM_DEVICE_STATE_REASON_SUPPLICANT_AVAILABLE = 42
	NM_DEVICE_STATE_REASON_MODEM_NOT_FOUND = 43
	NM_DEVICE_STATE_REASON_BT_FAILED = 44
	NM_DEVICE_STATE_REASON_GSM_SIM_NOT_INSERTED = 45
	NM_DEVICE_STATE_REASON_GSM_SIM_PIN_REQUIRED = 46
	NM_DEVICE_STATE_REASON_GSM_SIM_PUK_REQUIRED = 47
	NM_DEVICE_STATE_REASON_GSM_SIM_WRONG = 48
	NM_DEVICE_STATE_REASON_INFINIBAND_MODE = 49
	NM_DEVICE_STATE_REASON_DEPENDENCY_FAILED = 50
	NM_DEVICE_STATE_REASON_BR2684_FAILED = 51
	NM_DEVICE_STATE_REASON_MODEM_MANAGER_UNAVAILABLE = 52
	NM_DEVICE_STATE_REASON_SSID_NOT_FOUND = 53
	NM_DEVICE_STATE_REASON_SECONDARY_CONNECTION_FAILED = 54
	NM_DEVICE_STATE_REASON_DCB_FCOE_FAILED = 55
	NM_DEVICE_STATE_REASON_TEAMD_CONTROL_FAILED = 56
	NM_DEVICE_STATE_REASON_MODEM_FAILED = 57
	NM_DEVICE_STATE_REASON_MODEM_AVAILABLE = 58
	NM_DEVICE_STATE_REASON_SIM_PIN_INCORRECT = 59
	NM_DEVICE_STATE_REASON_NEW_ACTIVATION = 60
	NM_DEVICE_STATE_REASON_PARENT_CHANGED = 61
	NM_DEVICE_STATE_REASON_PARENT_MANAGED_CHANGED = 62
)


type ConnectionSetting struct {
	Setting		map[string]map[string]dbus.Variant
	Ssid		string
	IsWifi		bool
	IsOpenWifi	bool
}

func uint32ToIp(ipInt uint32) (string) {
     var ipByte [4]byte
     ip := ipByte[:]
     binary.LittleEndian.PutUint32(ip, ipInt)
     return net.IPv4(ip[0], ip[1], ip[2], ip[3]).String()
}

func GetActiveConnections() ([]dbus.ObjectPath, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		panic(err)
	}
	obj := conn.Object("org.freedesktop.NetworkManager",
		"/org/freedesktop/NetworkManager")
	props, err := obj.GetProperty("org.freedesktop.NetworkManager.ActiveConnections")
	if err != nil {
		panic(err)
	}
	return props.Value().([]dbus.ObjectPath), nil
}

func GetConnection(activeConnection dbus.ObjectPath) (dbus.ObjectPath, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		panic(err)
	}
	obj := conn.Object("org.freedesktop.NetworkManager", activeConnection)
	props, err := obj.GetProperty("org.freedesktop.NetworkManager.Connection.Active.Connection")
	if err != nil {
		panic(err)
	}
	return props.Value().(dbus.ObjectPath), nil
}


func GetConnectionSetting(connection dbus.ObjectPath) (ConnectionSetting) {
	conn, err := dbus.SystemBus()
	if err != nil {
		panic(err)
	}
	obj := conn.Object("org.freedesktop.NetworkManager", connection)
	call :=	obj.Call("org.freedesktop.NetworkManager.Settings.Connection.GetSettings", 0)
	if call.Err != nil {
		panic(call.Err)
	}
	connectionSetting := ConnectionSetting{}
	connectionSetting.Setting = call.Body[0].(map[string]map[string]dbus.Variant)
	connectionSetting.Ssid = fmt.Sprintf("%s",
		connectionSetting.Setting["802-11-wireless"]["ssid"].Value())
	settingWifi := connectionSetting.Setting["802-11-wireless"]
	if len(settingWifi) > 0 {
		connectionSetting.IsWifi = true
	}
	if connectionSetting.IsWifi {
		if len(connectionSetting.Setting["802-11-wireless-security"]) == 0 {
			connectionSetting.IsOpenWifi = true
		}
	}
	return connectionSetting
}

func GetOpenWifiConnections(activeConnections []dbus.ObjectPath) ([]ConnectionSetting) {
	var openWifiConnections []ConnectionSetting
	for _, activeConnection := range activeConnections {
		connection, err := GetConnection(activeConnection)
		if err != nil {
			panic(err)
		}
		setting := GetConnectionSetting(connection)
		if setting.IsOpenWifi {
			openWifiConnections = append(openWifiConnections, setting)
		}
	}
	return openWifiConnections
}

func GetDhcpNameservers(activeConnection dbus.ObjectPath) ([]string, error) {
	var nameservers []string
	conn, err := dbus.SystemBus()
	if err != nil {
		return nameservers, err
	}
	obj := conn.Object("org.freedesktop.NetworkManager", activeConnection)
	configProps, err := obj.GetProperty("org.freedesktop.NetworkManager.Connection.Active.Ip4Config")
	if err != nil {
		return nameservers, err
	}
	obj2 := conn.Object("org.freedesktop.NetworkManager", configProps.Value().(dbus.ObjectPath))
	nsProps, err := obj2.GetProperty("org.freedesktop.NetworkManager.IP4Config.Nameservers")
	for _, ns := range nsProps.Value().([]uint32) {
		nameservers = append(nameservers, uint32ToIp(ns))
	}
	return nameservers, nil
}
