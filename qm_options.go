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

// Network device

type NetworkDevice struct {
	// Network Card Model. The virtio model provides the best performance with very low CPU overhead.
	// If your guest does not support this driver, it is usually best to use e1000.
	Model *NetworkCardModel

	// Bridge to attach the network device to. The Proxmox VE standard bridge is called vmbr0.
	Bridge *string

	// Whether this interface should be protected by the firewall.
	Firewall *bool

	// Whether this interface should be disconnected (like pulling the plug).
	LinkDown *bool

	// <XX:XX:XX:XX:XX:XX> MAC address.
	// That address must be unique withing your network. This is automatically generated if not specified.
	MacAddr *string

	// (0 - 16) Number of packet queues to be used on the device.
	Queues *int

	// (0 - N) Rate limit in mbps (megabytes per second) as floating point number.
	Rate *float64

	// (1 - 4094) VLAN tag to apply to packets on this interface.
	Tag *int

	// <vlanid[;vlanid...]> VLAN trunks to pass through this interface.
	Trunks *string
}

func NewNetworkDeviceFromString(value string) *NetworkDevice {
	d := &NetworkDevice{}
	options := strings.Split(value, ",")
	for _, option := range options {
		optionParts := strings.Split(option, "=")
		if len(optionParts) == 2 {
			k := optionParts[0]
			v := optionParts[1]
			switch k {
			case "bridge":
				d.Bridge = String(v)
			case "firewall":
				d.Firewall = Bool(stringToBool(v))
			case "link_down":
				d.LinkDown = Bool(stringToBool(v))
			case "queues":
				val, _ := strconv.Atoi(v)
				d.Queues = Int(val)
			case "rate":
				val, _ := strconv.ParseFloat(v, 10)
				d.Rate = Float64(val)
			case "tag":
				val, _ := strconv.Atoi(v)
				d.Tag = Int(val)
			case "trunks":
				d.Trunks = String(v)
			}
		}
	}
	optionParts := strings.Split(options[0], "=")
	v, _ := NetworkCardModelFromString(optionParts[0])
	d.Model = &v
	d.MacAddr = String(optionParts[1])

	return d
}

func (c *NetworkDevice) GetQMOptionValue() string {
	v := make([]string, 0, 1)
	if c.Model != nil {
		v = append(v, fmt.Sprintf("%s=%s", "model", c.Model.String()))
	}
	if c.Bridge != nil {
		v = append(v, fmt.Sprintf("%s=%s", "bridge", *c.Bridge))
	}
	if c.Firewall != nil {
		v = append(v, fmt.Sprintf("%s=%s", "firewall", boolToString(*c.Firewall)))
	}
	if c.LinkDown != nil {
		v = append(v, fmt.Sprintf("%s=%s", "link_down", boolToString(*c.LinkDown)))
	}
	if c.MacAddr != nil {
		v = append(v, fmt.Sprintf("%s=%s", "macaddr", *c.MacAddr))
	}
	if c.Queues != nil {
		v = append(v, fmt.Sprintf("%s=%s", "queues", *c.Queues))
	}
	if c.Rate != nil {
		v = append(v, fmt.Sprintf("%s=%s", "rate", *c.Rate))
	}
	if c.Tag != nil {
		v = append(v, fmt.Sprintf("%s=%s", "tag", *c.Tag))
	}
	if c.Trunks != nil {
		v = append(v, fmt.Sprintf("%s=%s", "trunks", *c.Trunks))
	}
	return strings.Join(v, ",")
}

// VirtIO device

type VirtIODevice struct {
	File     *string
	Format   *VolumeFormat
	Backup   *bool
	IOThread *bool
	Size     *string
	Snapshot *bool
}

func NewVirtIODeviceFromString(value string) *VirtIODevice {
	d := &VirtIODevice{}
	options := strings.Split(value, ",")
	for _, option := range options {
		optionParts := strings.Split(option, "=")
		if len(optionParts) == 2 {
			k := optionParts[0]
			v := optionParts[1]
			switch k {
			case "file":
				d.File = String(v)
			case "format":
				format, _ := VolumeFormatFromString(v)
				d.Format = &format
			case "backup":
				d.Backup = Bool(stringToBool(v))
			case "iothread":
				d.IOThread = Bool(stringToBool(v))
			case "size":
				d.Size = String(v)
			case "snapshot":
				d.Snapshot = Bool(stringToBool(v))
			}
		}
	}
	d.File = String(options[0])

	return d
}

func (c *VirtIODevice) GetQMOptionValue() string {
	v := make([]string, 0, 1)
	if c.File != nil {
		v = append(v, fmt.Sprintf("%s=%s", "file", *c.File))
	}
	if c.Format != nil {
		v = append(v, fmt.Sprintf("%s=%s", "format", c.Format.String()))
	}
	if c.Backup != nil {
		v = append(v, fmt.Sprintf("%s=%s", "backup", boolToString(*c.Backup)))
	}
	if c.IOThread != nil {
		v = append(v, fmt.Sprintf("%s=%s", "iothread", boolToString(*c.IOThread)))
	}
	if c.Size != nil {
		v = append(v, fmt.Sprintf("%s=%s", "size", *c.Size))
	}
	if c.Snapshot != nil {
		v = append(v, fmt.Sprintf("%s=%s", "snapshot", boolToString(*c.Snapshot)))
	}
	return strings.Join(v, ",")
}

// IDE device
type IDEDevice struct {
	File  *string
	Media *MediaType
	Size  *string
}

func NewIDEDeviceFromString(value string) *IDEDevice {
	d := &IDEDevice{}
	options := strings.Split(value, ",")
	for _, option := range options {
		optionParts := strings.Split(option, "=")
		if len(optionParts) == 2 {
			k := optionParts[0]
			v := optionParts[1]
			switch k {
			case "file":
				d.File = String(v)
			case "media":
				media, _ := MediaTypeFromString(v)
				d.Media = &media
			case "size":
				d.Size = String(v)
			}
		}
	}
	d.File = String(options[0])

	return d
}

func (c *IDEDevice) GetQMOptionValue() string {
	v := make([]string, 0, 1)
	if c.File != nil {
		v = append(v, fmt.Sprintf("%s=%s", "file", *c.File))
	}
	if c.Media != nil {
		v = append(v, fmt.Sprintf("%s=%s", "media", c.Media.String()))
	}
	if c.Size != nil {
		v = append(v, fmt.Sprintf("%s=%s", "size", *c.Size))
	}
	return strings.Join(v, ",")
}

// Serial device
type SerialDevice struct {
	Value  string
}

func NewSerialDeviceFromString(value string) *SerialDevice {
	return &SerialDevice{value}
}

func (c *SerialDevice) GetQMOptionValue() string {
	return c.Value
}
