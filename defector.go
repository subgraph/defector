package main

import (
	"time"
	"fmt"
	"os/exec"
	"bytes"

	nm "github.com/subgraph/defector/networkmanager"
	tc "github.com/subgraph/defector/torcontrol"
	"github.com/godbus/dbus"
	"github.com/TheCreeper/go-notify"
)

func MonitorDeviceStateChanged() {
	conn, err := dbus.SystemBus()
	if err != nil {
		panic(err)
	}
	conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
	                "type='signal',interface='org.freedesktop.NetworkManager',member='StateChanged'")
	fmt.Printf("Monitoring org.freedesktop.NetworkManager for StateChanged signals.\n")
	c := make(chan *dbus.Signal, 10)
	conn.Signal(c)
        for v := range c {
		HandleDeviceStateChange(v.Body[0].(uint32))
        }
}

func HandleDeviceStateChange(state uint32) () {
	cmd := exec.Command("./captivebrowser/captivebrowser", "-d")
	if state == nm.NM_DEVICE_STATE_IP_CONFIG {
		time.Sleep(1 * time.Second)
		connections, err := nm.GetActiveConnections()
		if err != nil {
			panic(err)
		}
		fmt.Println("Checking for open wifi connections")
		if len(nm.GetOpenWifiConnections(connections)) > 0 {
			fmt.Println("Running detection")
			runDetection := detectCaptivePortalNotification()
			if runDetection {
				var out bytes.Buffer
				cmd.Stdout = &out
				err := cmd.Run()
				if err != nil {
					// TODO: actual error handling
					_ = failedCaptivePortalNotification()
				}
				fmt.Println("Checking Tor connection.")
				if !tc.IsTorConnected() {
					authenticatedCaptivePortalNotification()
				} else {
					fmt.Println("Tor is connected.")
				}
			}
		} else {
			fmt.Println("No open wifi connections found.")
		}
	}
}

func detectCaptivePortalNotification() (bool) {
	hints := make(map[string]interface{})
	hints["notify.Persistence"] = true
	notification := notify.Notification{AppName: "Tor cannot connect to the internet",
				Summary: "Tor is having problems connecting. Detect captive portal?",
				Actions: []string{"No", "No", "Yes", "Yes"},
				AppIcon: "network-wireless",
				Hints: hints}
	id, err := notification.Show()
	if err != nil {
		panic(err)
	}
	return HandleNotificationAction(id)
}

func failedCaptivePortalNotification() (bool) {
	hints := make(map[string]interface{})
	hints["notify.Persistence"] = true
	notification := notify.Notification{AppName: "Tor cannot connect to the internet",
				Summary: "Did not detect a captive portal",
				Actions: []string{"Cancel", "Cancel", "Retry", "Retry"},
				AppIcon: "network-wireless",
				Hints: hints}
	_ , err := notification.Show()
	if err != nil {
		panic(err)
	}
	return true
}

func authenticatedCaptivePortalNotification() {
	hints := make(map[string]interface{})
	hints["notify.Persistence"] = true
	notification := notify.Notification{AppName: "Tor cannot connect to the internet",
				Summary: "Captive portal authentication successful",
				AppIcon: "network-wireless",
				Hints: hints}
	_ , err := notification.Show()
	if err != nil {
		panic(err)
	}
}

func HandleNotificationAction(id uint32) (bool) {
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}
	conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, 
		"type='signal',interface='org.freedesktop.Notifications',member='ActionInvoked'")
	c := make(chan *dbus.Signal, 10)
	conn.Signal(c)
	for v := range c {
		actionId := v.Body[0].(uint32)
		actionString := v.Body[1].(string)
		if actionString == "Yes" && actionId == id {
			return true
		}
	}
	return false
}

func main() {
	MonitorDeviceStateChanged()
}
