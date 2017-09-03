package goproxmox

import (
	"fmt"
	"strconv"
	"strings"
)

type QMOption interface {
	GetQMOptionValue() string
}

type rawQMOption struct {
	value string
}

func NewRawQMOption(value string) *rawQMOption {
	return &rawQMOption{value: value}
}

func (c *rawQMOption) GetQMOptionValue() string {
	return c.value
}

type networkDevice struct {
	configMap map[string]string
}

func NewNetworkDevice(cardModel NetworkCardModel) *networkDevice {
	config := &networkDevice{
		configMap: make(map[string]string),
	}
	config.SetModel(cardModel)
	return config
}

func NewNetworkDeviceFromString(value string) *networkDevice {
	config := &networkDevice{
		configMap: make(map[string]string),
	}
	for _, option := range strings.Split(value, ",") {
		optionParts := strings.Split(option, "=")
		if len(optionParts) == 2 {
			config.configMap[optionParts[0]] = optionParts[1]
		}
	}
	return config
}

// Network Card Model. The virtio model provides the best performance with very low CPU overhead.
// If your guest does not support this driver, it is usually best to use e1000.
func (c *networkDevice) SetModel(value NetworkCardModel) {
	c.configMap["model"] = value.String()
}

// Bridge to attach the network device to. The Proxmox VE standard bridge is called vmbr0.
func (c *networkDevice) SetBridge(value string) {
	c.configMap["bridge"] = value
}

// Whether this interface should be protected by the firewall.
func (c *networkDevice) SetFirewall(value bool) {
	c.configMap["firewall"] = boolToString(value)
}

// Whether this interface should be disconnected (like pulling the plug).
func (c *networkDevice) SetLinkDown(value bool) {
	c.configMap["link_down"] = boolToString(value)
}

// <XX:XX:XX:XX:XX:XX> MAC address.
// That address must be unique withing your network. This is automatically generated if not specified.
func (c *networkDevice) SetMacAddr(value string) {
	c.configMap["macaddr"] = value
}

// (0 - 16) Number of packet queues to be used on the device.
func (c *networkDevice) SetQueues(value int) {
	c.configMap["queues"] = strconv.Itoa(value)
}

// (0 - N) Rate limit in mbps (megabytes per second) as floating point number.
func (c *networkDevice) SetRate(value float32) {
	c.configMap["rate"] = fmt.Sprintf("%.2f", value)
}

// (1 - 4094) VLAN tag to apply to packets on this interface.
func (c *networkDevice) SetTag(value int) {
	c.configMap["tag"] = strconv.Itoa(value)
}

// <vlanid[;vlanid...]> VLAN trunks to pass through this interface.
func (c *networkDevice) SetTrunks(value string) {
	c.configMap["trunks"] = value
}

func (c *networkDevice) GetQMOptionValue() string {
	v := make([]string, 0, len(c.configMap))
	for key, value := range c.configMap {
		v = append(v, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(v, ",")
}
