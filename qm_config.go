package goproxmox

import (
	"strconv"
	"fmt"
	"strings"
)

func NewVMConfig() *vmConfigMap {
	c := &vmConfigMap{
		configMap: make(map[string]string),
	}
	return c
}

func NewVMConfigFromMap(data map[string]interface{}) *vmConfigMap {
	c := &vmConfigMap{
		configMap: make(map[string]string),
	}
	for k, v := range data {
		switch value := v.(type) {
		case string:
			c.configMap[k] = value
		case int:
			c.configMap[k] = strconv.Itoa(value)
		case float64:
			c.configMap[k] = fmt.Sprintf("%.0f", value)
		default:
			fmt.Printf("[WARN] Unknown type for key %s. Value: %v", k, v)
		}
	}
	return c
}

type VMConfig interface {
	GetOptionsMap() map[string]string

	SetACPI(bool)
	SetQemuAgent(bool)
	SetArchive(string) // Create only
	SetArgs(string)
	SetAutoStart(bool)
	// background_delay. Only post
	SetBalloon(int) error
	SetBios(Bios)
	SetBootOrder(...BootDevice) error
	SetBootDisk(string)
	SetCDROM(string)
	SetCores(cores int) error
	SetCPU(CPUType)
	SetCPULimit(int) error
	SetCPUUnits(int) error
	// delete. Only post/put
	SetDescription(string)
	// digest. Only post/put
	SetForce(bool)
	SetFreeze(bool)
	SetHostPCI(string)
	SetHotPlug(string)
	SetHugePages(HugePages)
	AddIDEDevice(int, string) error
	SetKeyboardLayout(KeyboardLayout)
	SetKVMHardwareVirtualization(bool)
	SetLocalTime(bool)
	SetLock(Lock)
	SetMachineType(string)
	SetMemory(int) error
	SetMigrateDowntime(int) error
	SetMigrateSpeed(int) error
	SetName(name string)
	AddNetworkDevice(int, *networkDevice) error
	GetNetworkDevices() []networkDevice
	SetNUMA(bool)
	AddNUMA(string)
	SetStartAtBoot(bool)
	SetOSType(OSType)
	AddParallelDevice(number int, value string) error
	SetPool(string) // Create only
	SetProtection(bool)
	SetReboot(bool)
	// revert. Only post/put
	AddSATADevice(int, string) error
	AddSCSIDevice(int, string) error
	SetSCSIControllerType(SCSIControllerType)
	AddSerialDevice(int, string) error
	SetMemoryShares(int) error
	// skiplock. Only post/put
	SetSMBIOS1(string)
	SetSMP(int) error
	SetSockets(int) error
	SetStartDate(string)
	SetStartup(string)
	SetStorage(string) // Create only
	SetTablet(bool)
	SetTDF(bool)
	SetTemplate(bool)
	SetUnique(bool) // Create only
	AddUSBDevice(int, string) error
	SetVCPUs(int)
	SetVGA(VGAType)
	AddVirtIODevice(int, string) error
	SetVMID(vmID int) error // Create only
	SetWatchdog(string)
}

type vmConfigMap struct {
	configMap map[string]string
}

func boolToString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func (c *vmConfigMap) GetOptionsMap() map[string]string {
	return c.configMap
}

// Enable/disable ACPI.
// default = 1
func (c *vmConfigMap) SetACPI(value bool) {
	c.configMap["acpi"] = boolToString(value)
}

// Enable/disable Qemu GuestAgent.
// default = 0
func (c *vmConfigMap) SetQemuAgent(value bool) {
	c.configMap["agent"] = boolToString(value)
}

// The backup file.
func (c *vmConfigMap) SetArchive(value string) {
	c.configMap["archive"] = value
}

// Arbitrary arguments passed to kvm, for example:
// args: -no-reboot -no-hpet
func (c *vmConfigMap) SetArgs(value string) {
	c.configMap["args"] = value
}

// Automatic restart after crash
// default = 0
func (c *vmConfigMap) SetAutoStart(value bool) {
	c.configMap["autostart"] = boolToString(value)
}

// Amount of target RAM for the VM in MB. Using zero disables the balloon driver.
func (c *vmConfigMap) SetBalloon(value int) error {
	if value < 0 {
		return NewArgError("balloon", "it can't be < 0")
	}
	c.configMap["balloon"] = strconv.Itoa(value)
	return nil
}

// Select BIOS implementation.
// default = seabios
func (c *vmConfigMap) SetBios(value Bios) {
	c.configMap["bios"] = value.String()
}

// Boot on floppy (a), hard disk (c), CD-ROM (d), or network (n).
// default = cdn
func (c *vmConfigMap) SetBootOrder(values ...BootDevice) error {
	value := ""
	if len(values) > 4 {
		return NewArgError("boot", "there are too many boot devices specified")
	}
	for _, bootDevice := range values {
		value += bootDevice.String()
	}
	c.configMap["boot"] = value
	return nil
}

// Enable booting from specified disk.
// (ide|sata|scsi|virtio)\d+
func (c *vmConfigMap) SetBootDisk(value string) {
	c.configMap["bootdisk"] = value
}

// <volume> This is an alias for option -ide2
func (c *vmConfigMap) SetCDROM(value string) {
	c.configMap["cdrom"] = value
}

// The number of cores per socket.
// default = 1
func (c *vmConfigMap) SetCores(value int) error {
	if value < 1 {
		return NewArgError("cores", "it must be > 0")
	}
	c.configMap["cores"] = strconv.Itoa(value)
	return nil
}

// [cputype=]<enum> [,hidden=<1|0>] Emulated CPU type.
// default = kvm64
// hidden=<boolean> Do not identify as a KVM virtual machine.
// default = 0
func (c *vmConfigMap) SetCPU(value CPUType) {
	c.configMap["cpu"] = value.String()
}

// (0 - 128) Limit of CPU usage. NOTE: If the computer has 2 CPUs, it has total of '2' CPU time. Value '0' indicates no CPU limit.
// default = 0
func (c *vmConfigMap) SetCPULimit(value int) error {
	if value < 0 || value > 128 {
		return NewArgError("cpulimit", "it must be 0 to 128")
	}
	c.configMap["cpulimit"] = strconv.Itoa(value)
	return nil
}

// (0 - 500000) CPU weight for a VM. Argument is used in the kernel fair scheduler.
// The larger the number is, the more CPU time this VM gets.
// Number is relative to weights of all the other running VMs.
// You can disable fair-scheduler configuration by setting this to 0.
// default = 1024
func (c *vmConfigMap) SetCPUUnits(value int) error {
	if value < 0 || value > 500000 {
		return NewArgError("cpulimit", "it must be 0 to 500000")
	}
	c.configMap["cpuunits"] = strconv.Itoa(value)
	return nil
}

// Description for the VM. Only used on the configuration web interface.
// This is saved as comment inside the configuration file.
func (c *vmConfigMap) SetDescription(value string) {
	c.configMap["description"] = value
}

// Allow to overwrite existing VM.
func (c *vmConfigMap) SetForce(value bool) {
	c.configMap["force"] = boolToString(value)
}

// Freeze CPU at startup (use 'c' monitor command to start execution).
func (c *vmConfigMap) SetFreeze(value bool) {
	c.configMap["freeze"] = boolToString(value)
}

// Map host PCI devices into guest.
// NOTE: This option allows direct access to host hardware.
// So it is no longer possible to migrate such machines - use with special care.
func (c *vmConfigMap) SetHostPCI(value string) {
	c.configMap["hostpci"] = value
}

// Selectively enable hotplug features.
// This is a comma separated list of hotplug features: 'network', 'disk', 'cpu', 'memory' and 'usb'.
// Use '0' to disable hotplug completely.
// Value '1' is an alias for the default 'network,disk,usb'.
func (c *vmConfigMap) SetHotPlug(value string) {
	c.configMap["hotplug"] = value
}

// Enable/disable hugepages memory.
func (c *vmConfigMap) SetHugePages(value HugePages) {
	c.configMap["hugepages"] = value.String()
}

// Use volume as IDE hard disk or CD-ROM
func (c *vmConfigMap) AddIDEDevice(number int, value string) error {
	if number < 0 || number > 3 {
		return NewArgError("ide[n]", "it must be 0 to 3")
	}

	key := fmt.Sprintf("ide%d", number)
	if _, ok := c.configMap[key]; ok {
		return NewArgError("ide[n]", fmt.Sprintf("IDE device %s already exists", key))
	}

	c.configMap[key] = value

	return nil
}

// Keyboard layout for vnc server. Default is read from the '/etc/pve/datacenter.conf' configuration file.
// default = en-us
func (c *vmConfigMap) SetKeyboardLayout(value KeyboardLayout) {
	c.configMap["keyboard"] = value.String()
}

// Enable/disable KVM hardware virtualization.
// default = 1
func (c *vmConfigMap) SetKVMHardwareVirtualization(value bool) {
	c.configMap["kvm"] = boolToString(value)
}

// Set the real time clock to local time. This is enabled by default if ostype indicates a Microsoft OS.
func (c *vmConfigMap) SetLocalTime(value bool) {
	c.configMap["localtime"] = boolToString(value)
}

// Lock/unlock the VM.
func (c *vmConfigMap) SetLock(value Lock) {
	c.configMap["lock"] = value.String()
}

// Specify the Qemu machine type.
// (pc|pc(-i440fx)?-\d+\.\d+(\.pxe)?|q35|pc-q35-\d+\.\d+(\.pxe)?)
func (c *vmConfigMap) SetMachineType(value string) {
	c.configMap["machine"] = value
}

// Amount of RAM for the VM in MB. This is the maximum available memory when you use the balloon device.
// default = 512
func (c *vmConfigMap) SetMemory(value int) error {
	if value < 16 {
		return NewArgError("memory", "it must be >= 16")
	}
	c.configMap["memory"] = strconv.Itoa(value)
	return nil
}

// (0 - N) Set maximum tolerated downtime (in seconds) for migrations.
// default = 0.1
func (c *vmConfigMap) SetMigrateDowntime(value int) error {
	if value < 1 {
		return NewArgError("migrate_downtime", "it must be >= 0")
	}
	c.configMap["migrate_downtime"] = strconv.Itoa(value)
	return nil
}

// Set maximum speed (in MB/s) for migrations. Value 0 is no limit.
// default = 0
func (c *vmConfigMap) SetMigrateSpeed(value int) error {
	if value < 1 {
		return NewArgError("migrate_speed", "it must be >= 0")
	}
	c.configMap["migrate_speed"] = strconv.Itoa(value)
	return nil
}

// Set a name for the VM. Only used on the configuration web interface.
func (c *vmConfigMap) SetName(name string) {
	c.configMap["name"] = name
}

// Specify network devices.
func (c *vmConfigMap) AddNetworkDevice(number int, value *networkDevice) error {
	key := fmt.Sprintf("net%d", number)
	if _, ok := c.configMap[key]; ok {
		return NewArgError("net[n]", fmt.Sprintf("Network device %s already exists", key))
	}

	c.configMap[key] = value.GetQMOptionValue()
	return nil
}

func (c *vmConfigMap) GetNetworkDevices() []networkDevice {
	networkDevices := make([]networkDevice, 0, 0)
	for key, value := range c.configMap {
		if strings.HasPrefix(key, "net") {
			networkDevices = append(networkDevices, *NewNetworkDeviceFromString(value))
		}
	}
	return networkDevices
}

// Enable/disable NUMA.
// default = 0
func (c *vmConfigMap) SetNUMA(value bool) {
	c.configMap["numa"] = boolToString(value)
}

// NUMA topology.
func (c *vmConfigMap) AddNUMA(value string) {
	c.configMap["numa1"] = value
}

// Specifies whether a VM will be started during system bootup.
// default = 0
func (c *vmConfigMap) SetStartAtBoot(value bool) {
	c.configMap["onboot"] = boolToString(value)
}

// Specify guest operating system. This is used to enable special optimization/features for specific operating systems.
func (c *vmConfigMap) SetOSType(value OSType) {
	c.configMap["ostype"] = value.String()
}

// Map host parallel devices (n is 0 to 2).
// NOTE: This option allows direct access to host hardware.
// So it is no longer possible to migrate such machines - use with special care.
func (c *vmConfigMap) AddParallelDevice(number int, value string) error {
	if number < 0 || number > 2 {
		return NewArgError("parallel[n]", "it must be 0 to 2")
	}

	key := fmt.Sprintf("parallel%d", number)
	if _, ok := c.configMap[key]; ok {
		return NewArgError("parallel[n]", fmt.Sprintf("Parallel device %s already exists", key))
	}

	c.configMap[key] = value

	return nil
}

// Add theVM to the specified pool.
func (c *vmConfigMap) SetPool(value string) {
	c.configMap["pool"] = value
}

// Sets the protection flag of the VM. This will disable the remove VM and remove disk operations.
func (c *vmConfigMap) SetProtection(value bool) {
	c.configMap["protection"] = boolToString(value)
}

// Allow reboot. If set to '0' the VM exit on reboot.
func (c *vmConfigMap) SetReboot(value bool) {
	c.configMap["reboot"] = boolToString(value)
}

// Use volume as SATA hard disk or CD-ROM (n is 0 to 5).
func (c *vmConfigMap) AddSATADevice(number int, value string) error {
	if number < 0 || number > 5 {
		return NewArgError("sata[n]", "it must be 0 to 5")
	}

	key := fmt.Sprintf("sata%d", number)
	if _, ok := c.configMap[key]; ok {
		return NewArgError("sata[n]", fmt.Sprintf("SATA device %s already exists", key))
	}

	c.configMap[key] = value

	return nil
}

// Use volume as SCSI hard disk or CD-ROM (n is 0 to 13).
func (c *vmConfigMap) AddSCSIDevice(number int, value string) error {
	if number < 0 || number > 13 {
		return NewArgError("scsi[n]", "it must be 0 to 13")
	}

	key := fmt.Sprintf("scsi%d", number)
	if _, ok := c.configMap[key]; ok {
		return NewArgError("scsi[n]", fmt.Sprintf("SCSI device %s already exists", key))
	}

	c.configMap[key] = value

	return nil
}

// SCSI controller model
// default = lsi
func (c *vmConfigMap) SetSCSIControllerType(value SCSIControllerType) {
	c.configMap["scsihw"] = value.String()
}

// Create a serial device inside the VM (n is 0 to 3), and pass through a host serial device (i.e. /dev/ttyS0),
// or create a unix socket on the host side (use 'qm terminal' to open a terminal connection).
// NOTE: If you pass through a host serial device, it is no longer possible to migrate such machines - use with special care.
func (c *vmConfigMap) AddSerialDevice(number int, value string) error {
	if number < 0 || number > 3 {
		return NewArgError("serial[n]", "it must be 0 to 3")
	}

	key := fmt.Sprintf("serial%d", number)
	if _, ok := c.configMap[key]; ok {
		return NewArgError("serial[n]", fmt.Sprintf("Serial device %s already exists", key))
	}

	c.configMap[key] = value

	return nil
}

// Amount of memory shares for auto-ballooning.
// The larger the number is, the more memory this VM gets.
// Number is relative to weights of all other running VMs. Using zero disables auto-ballooning
// default = 1000
func (c *vmConfigMap) SetMemoryShares(value int) error {
	if value < 0 || value > 50000 {
		return NewArgError("shares", "it must be 0 to 50000")
	}
	c.configMap["shares"] = strconv.Itoa(value)
	return nil
}

// Specify SMBIOS type 1 fields.
func (c *vmConfigMap) SetSMBIOS1(value string) {
	c.configMap["smbios1"] = value
}

// The number of CPUs. Please use option -sockets instead.
// default = 1
func (c *vmConfigMap) SetSMP(value int) error {
	if value < 1 {
		return NewArgError("smp", "it must be > 0")
	}
	c.configMap["smp"] = strconv.Itoa(value)
	return nil
}

// The number of CPU sockets.
// default = 1
func (c *vmConfigMap) SetSockets(value int) error {
	if value < 1 {
		return NewArgError("sockets", "it must be > 0")
	}
	c.configMap["sockets"] = strconv.Itoa(value)
	return nil
}

// Set the initial date of the real time clock. Valid format for date are: 'now' or '2006-06-17T16:01:21' or'2006-06-17'.
// (now |YYYY-MM-DD | YYYY-MM-DDTHH:MM:SS)
// default = now
func (c *vmConfigMap) SetStartDate(value string) {
	c.configMap["startdate"] = value
}

// Startup and shutdown behavior.
// Order is a non-negative number defining the general startup order. Shutdown is done with reverse ordering.
// Additionally you can set the 'up' or 'down' delay in seconds, which specifies a delay to wait before the next VM is started or stopped.
// [[order=]\d+] [,up=\d+] [,down=\d+]
func (c *vmConfigMap) SetStartup(value string) {
	c.configMap["startup"] = value
}

// Default storage.
func (c *vmConfigMap) SetStorage(value string) {
	c.configMap["storage"] = value
}

// Enable/disable the USB tablet device. This device is usually needed to allow absolute mouse positioning with VNC.
// Else the mouse runs out of sync with normal VNC clients. If you're running lots of console-only guests on one host,
// you may consider disabling this to save some context switches. This is turned off by default if you use spice (-vga=qxl).
// default = 1
func (c *vmConfigMap) SetTablet(value bool) {
	c.configMap["tablet"] = boolToString(value)
}

// Enable/disable time drift fix.
// default = 0
func (c *vmConfigMap) SetTDF(value bool) {
	c.configMap["tdf"] = boolToString(value)
}

// Enable/disable Template.
// default = 0
func (c *vmConfigMap) SetTemplate(value bool) {
	c.configMap["template"] = boolToString(value)
}

// Assign a unique random ethernet address.
func (c *vmConfigMap) SetUnique(value bool) {
	c.configMap["unique"] = boolToString(value)
}

// Configure an USB device (n is 0 to 4).
func (c *vmConfigMap) AddUSBDevice(number int, value string) error {
	if number < 0 || number > 4 {
		return NewArgError("usb[n]", "it must be 0 to 4")
	}

	key := fmt.Sprintf("usb%d", number)
	if _, ok := c.configMap[key]; ok {
		return NewArgError("usb[n]", fmt.Sprintf("USB device %s already exists", key))
	}

	c.configMap[key] = value

	return nil
}

// Number of hotplugged vCPUs.
// default = 0
func (c *vmConfigMap) SetVCPUs(value int) {
	c.configMap["vcpus"] = strconv.Itoa(value)
}

// Select the VGA type. If you want to use high resolution modes (>= 1280x1024x16) then you should use the options std or vmware.
// Default is std for win8/win7/w2k8, and cirrus for other OS types. The qxl option enables the SPICE display sever.
// For win* OS you can select how many independent displays you want, Linux guests can add displays them self.
// You can also run without any graphic card, using a serial device as terminal.
func (c *vmConfigMap) SetVGA(value VGAType) {
	c.configMap["vga"] = value.String()
}

// Use volume as VirtIO hard disk (n is 0 to 15).
func (c *vmConfigMap) AddVirtIODevice(number int, value string) error {
	if number < 0 || number > 15 {
		return NewArgError("virtio[n]", "it must be 0 to 15")
	}

	key := fmt.Sprintf("virtio%d", number)
	if _, ok := c.configMap[key]; ok {
		return NewArgError("virtio[n]", fmt.Sprintf("VirtIO device %s already exists", key))
	}

	c.configMap[key] = value

	return nil
}

// The (unique) ID of the VM.
func (c *vmConfigMap) SetVMID(value int) error {
	if value < 100 {
		return NewArgError("vmid", "it should be >= 100. IDs < 100 are reserved for internal purposes.")
	}
	c.configMap["vmid"] = strconv.Itoa(value)
	return nil
}

// Create a virtual hardware watchdog device. Once enabled (by a guest action),
// the watchdog must be periodically polled by an agent inside the guest or else the watchdog will reset the guest
// (or execute the respective action specified)
func (c *vmConfigMap) SetWatchdog(value string) {
	c.configMap["watchdog"] = value
}
