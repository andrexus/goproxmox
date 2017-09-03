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
	optionsMap map[string]string
}

func NewNetworkDevice(cardModel NetworkCardModel) *networkDevice {
	options := &networkDevice{
		optionsMap: make(map[string]string),
	}
	options.SetModel(cardModel)
	return options
}

// Network Card Model. The virtio model provides the best performance with very low CPU overhead.
// If your guest does not support this driver, it is usually best to use e1000.
func (c *networkDevice) SetModel(value NetworkCardModel) {
	c.optionsMap["model"] = value.String()
}

// Bridge to attach the network device to. The Proxmox VE standard bridge is called vmbr0.
func (c *networkDevice) SetBridge(value string) {
	c.optionsMap["bridge"] = value
}

// Whether this interface should be protected by the firewall.
func (c *networkDevice) SetFirewall(value bool) {
	c.optionsMap["firewall"] = boolToString(value)
}

// Whether this interface should be disconnected (like pulling the plug).
func (c *networkDevice) SetLinkDown(value bool) {
	c.optionsMap["link_down"] = boolToString(value)
}

// <XX:XX:XX:XX:XX:XX> MAC address.
// That address must be unique withing your network. This is automatically generated if not specified.
func (c *networkDevice) SetMacAddr(value string) {
	c.optionsMap["macaddr"] = value
}

// (0 - 16) Number of packet queues to be used on the device.
func (c *networkDevice) SetQueues(value int) {
	c.optionsMap["queues"] = strconv.Itoa(value)
}

// (0 - N) Rate limit in mbps (megabytes per second) as floating point number.
func (c *networkDevice) SetRate(value float32) {
	c.optionsMap["rate"] = fmt.Sprintf("%.2f", value)
}

// (1 - 4094) VLAN tag to apply to packets on this interface.
func (c *networkDevice) SetTag(value int) {
	c.optionsMap["tag"] = strconv.Itoa(value)
}

// <vlanid[;vlanid...]> VLAN trunks to pass through this interface.
func (c *networkDevice) SetTrunks(value string) {
	c.optionsMap["trunks"] = value
}

func (c *networkDevice) GetQMOptionValue() string {
	v := make([]string, 0, len(c.optionsMap))
	for key, value := range c.optionsMap {
		v = append(v, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(v, ",")
}
