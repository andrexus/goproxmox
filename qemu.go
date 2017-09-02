package goproxmox

import (
	"fmt"
	"strconv"
)

type QemuService interface {
	GetVMs(node string) ([]VM, error)
	GetVMCurrentStatus(node string, vmID int) (*VMStatus, error)
	StartVM(node string, vmID int) error
	StopVM(node string, vmID int) error
	ShutdownVM(node string, vmID int) error
	ResetVM(node string, vmID int) error
	SuspendVM(node string, vmID int) error
	ResumeVM(node string, vmID int) error
	CreateVM(node string, createRequest VMCreateRequest) error
	DeleteVM(node string, vmID int) error
}

type QemuServiceOp struct {
	client *Client
}

var _ QemuService = &QemuServiceOp{}

type VM struct {
	VMID      int     `json:"vmid"`
	Name      string  `json:"name"`
	Status    string  `json:"status"`
	Pid       string  `json:"pid"`
	Template  string  `json:"template"`
	Cpu       float64 `json:"cpu"`
	Mem       int     `json:"mem"`
	Maxmem    int     `json:"maxmem"`
	Cpus      int     `json:"cpus"`
	Disk      int     `json:"disk"`
	Maxdisk   int     `json:"maxdisk"`
	Diskread  int     `json:"diskread"`
	Diskwrite int     `json:"diskwrite"`
	Uptime    int     `json:"uptime"`
	Netout    int     `json:"netout"`
	Netin     int     `json:"netin"`
}

type VMStatus struct {
	Name      string      `json:"name"`
	Status    string      `json:"status"`
	Qmpstatus string      `json:"qmpstatus"`
	Template  string      `json:"template"`
	Cpus      int         `json:"cpus"`
	Diskread  int         `json:"diskread"`
	Cpu       int         `json:"cpu"`
	Netin     int         `json:"netin"`
	Netout    int         `json:"netout"`
	Disk      int         `json:"disk"`
	Diskwrite int         `json:"diskwrite"`
	Maxdisk   int         `json:"maxdisk"`
	Maxmem    int         `json:"maxmem"`
	Ha        interface{} `json:"ha"`
	Uptime    int         `json:"uptime"`
	Pid       int         `json:"pid"`
	Mem       int         `json:"mem"`
}

type VMCreateRequest interface {
	GetOptionsMap() map[string]string
	SetACPI(bool)
	SetQemuAgent(bool)
	SetArchive(string)
	SetArgs(string)
	SetAutoStart(bool)
	SetBalloon(int) error
	SetBios(Bios)
	SetBootOrder(...BootDevice) error
	SetBootDisk(string)
	SetCDROM(string)
	SetCores(cores int) error
	SetCPU(CPUType)
	SetCPULimit(int) error
	SetCPUUnits(int) error
	SetDescription(string)
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
	AddNetworkDevice(int, string) error
	SetNUMA(bool)
	AddNUMA(string)
	SetStartAtBoot(bool)
	SetOSType(OSType)
	AddParallelDevice(number int, value string) error
	SetPool(string)
	SetProtection(bool)
	SetReboot(bool)
	AddSATADevice(int, string) error
	AddSCSIDevice(int, string) error
	SetSCSIControllerType(SCSIControllerType)
	AddSerialDevice(int, string) error
	SetMemoryShares(int) error
	SetSMBIOS1(string)
	SetSMP(int) error
	SetCPUSockets(int) error
	SetStartDate(string)
	SetStartup(string)
	SetStorage(string)
	SetTablet(bool)
	SetTDF(bool)
	SetTemplate(bool)
	SetUnique(bool)
	AddUSBDevice(int, string) error
	SetVCPUs(int)
	SetVGA(VGAType)
	AddVirtIODevice(int, string) error
	SetVMID(vmID int) error
	SetWatchdog(string)
}

type vmCreateOptions struct {
	optionsMap map[string]string
}

func NewVMCreateRequest(vmID int) (*vmCreateOptions, error) {
	createOptions := &vmCreateOptions{
		optionsMap: make(map[string]string),
	}
	if err := createOptions.SetVMID(vmID); err != nil {
		return nil, err
	}
	return createOptions, nil
}

func boolToString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func (c *vmCreateOptions) GetOptionsMap() map[string]string {
	return c.optionsMap
}

// Enable/disable ACPI.
// default = 1
func (c *vmCreateOptions) SetACPI(value bool) {
	c.optionsMap["acpi"] = boolToString(value)
}

// Enable/disable Qemu GuestAgent.
// default = 0
func (c *vmCreateOptions) SetQemuAgent(value bool) {
	c.optionsMap["agent"] = boolToString(value)
}

// The backup file.
func (c *vmCreateOptions) SetArchive(value string) {
	c.optionsMap["archive"] = value
}

// Arbitrary arguments passed to kvm, for example:
// args: -no-reboot -no-hpet
func (c *vmCreateOptions) SetArgs(value string) {
	c.optionsMap["args"] = value
}

// Automatic restart after crash
// default = 0
func (c *vmCreateOptions) SetAutoStart(value bool) {
	c.optionsMap["autostart"] = boolToString(value)
}

// Amount of target RAM for the VM in MB. Using zero disables the balloon driver.
func (c *vmCreateOptions) SetBalloon(value int) error {
	if value < 0 {
		return NewArgError("balloon", "it can't be < 0")
	}
	c.optionsMap["balloon"] = strconv.Itoa(value)
	return nil
}

// Select BIOS implementation.
// default = seabios
func (c *vmCreateOptions) SetBios(value Bios) {
	c.optionsMap["bios"] = value.String()
}

// Boot on floppy (a), hard disk (c), CD-ROM (d), or network (n).
// default = cdn
func (c *vmCreateOptions) SetBootOrder(values ...BootDevice) error {
	value := ""
	if len(values) > 4 {
		return NewArgError("boot", "there are too many boot devices specified")
	}
	for _, bootDevice := range values {
		value += bootDevice.String()
	}
	c.optionsMap["boot"] = value
	return nil
}

// Enable booting from specified disk.
// (ide|sata|scsi|virtio)\d+
func (c *vmCreateOptions) SetBootDisk(value string) {
	c.optionsMap["bootdisk"] = value
}

// <volume> This is an alias for option -ide2
func (c *vmCreateOptions) SetCDROM(value string) {
	c.optionsMap["cdrom"] = value
}

// The number of cores per socket.
// default = 1
func (c *vmCreateOptions) SetCores(value int) error {
	if value < 1 {
		return NewArgError("cores", "it must be > 0")
	}
	c.optionsMap["cores"] = strconv.Itoa(value)
	return nil
}

// [cputype=]<enum> [,hidden=<1|0>] Emulated CPU type.
// default = kvm64
// hidden=<boolean> Do not identify as a KVM virtual machine.
// default = 0
func (c *vmCreateOptions) SetCPU(value CPUType) {
	c.optionsMap["cpu"] = value.String()
}

// (0 - 128) Limit of CPU usage. NOTE: If the computer has 2 CPUs, it has total of '2' CPU time. Value '0' indicates no CPU limit.
// default = 0
func (c *vmCreateOptions) SetCPULimit(value int) error {
	if value < 0 || value > 128 {
		return NewArgError("cpulimit", "it must be 0 to 128")
	}
	c.optionsMap["cpulimit"] = strconv.Itoa(value)
	return nil
}

// (0 - 500000) CPU weight for a VM. Argument is used in the kernel fair scheduler.
// The larger the number is, the more CPU time this VM gets.
// Number is relative to weights of all the other running VMs.
// You can disable fair-scheduler configuration by setting this to 0.
// default = 1024
func (c *vmCreateOptions) SetCPUUnits(value int) error {
	if value < 0 || value > 500000 {
		return NewArgError("cpulimit", "it must be 0 to 500000")
	}
	c.optionsMap["cpuunits"] = strconv.Itoa(value)
	return nil
}

// Description for the VM. Only used on the configuration web interface.
// This is saved as comment inside the configuration file.
func (c *vmCreateOptions) SetDescription(value string) {
	c.optionsMap["description"] = value
}

// Allow to overwrite existing VM.
func (c *vmCreateOptions) SetForce(value bool) {
	c.optionsMap["force"] = boolToString(value)
}

// Freeze CPU at startup (use 'c' monitor command to start execution).
func (c *vmCreateOptions) SetFreeze(value bool) {
	c.optionsMap["freeze"] = boolToString(value)
}

// Map host PCI devices into guest.
// NOTE: This option allows direct access to host hardware.
// So it is no longer possible to migrate such machines - use with special care.
func (c *vmCreateOptions) SetHostPCI(value string) {
	c.optionsMap["hostpci"] = value
}

// Selectively enable hotplug features.
// This is a comma separated list of hotplug features: 'network', 'disk', 'cpu', 'memory' and 'usb'.
// Use '0' to disable hotplug completely.
// Value '1' is an alias for the default 'network,disk,usb'.
func (c *vmCreateOptions) SetHotPlug(value string) {
	c.optionsMap["hotplug"] = value
}

// Enable/disable hugepages memory.
func (c *vmCreateOptions) SetHugePages(value HugePages) {
	c.optionsMap["hugepages"] = value.String()
}

// Use volume as IDE hard disk or CD-ROM
func (c *vmCreateOptions) AddIDEDevice(number int, value string) error {
	if number < 0 || number > 3 {
		return NewArgError("ide[n]", "it must be 0 to 3")
	}

	key := fmt.Sprintf("ide%d", number)
	if _, ok := c.optionsMap[key]; ok {
		return NewArgError("ide[n]", fmt.Sprintf("IDE device %s already exists", key))
	}

	c.optionsMap[key] = value

	return nil
}

// Keyboard layout for vnc server. Default is read from the '/etc/pve/datacenter.conf' configuration file.
// default = en-us
func (c *vmCreateOptions) SetKeyboardLayout(value KeyboardLayout) {
	c.optionsMap["keyboard"] = value.String()
}

// Enable/disable KVM hardware virtualization.
// default = 1
func (c *vmCreateOptions) SetKVMHardwareVirtualization(value bool) {
	c.optionsMap["kvm"] = boolToString(value)
}

// Set the real time clock to local time. This is enabled by default if ostype indicates a Microsoft OS.
func (c *vmCreateOptions) SetLocalTime(value bool) {
	c.optionsMap["localtime"] = boolToString(value)
}

// Lock/unlock the VM.
func (c *vmCreateOptions) SetLock(value Lock) {
	c.optionsMap["lock"] = value.String()
}

// Specify the Qemu machine type.
// (pc|pc(-i440fx)?-\d+\.\d+(\.pxe)?|q35|pc-q35-\d+\.\d+(\.pxe)?)
func (c *vmCreateOptions) SetMachineType(value string) {
	c.optionsMap["machine"] = value
}

// Amount of RAM for the VM in MB. This is the maximum available memory when you use the balloon device.
// default = 512
func (c *vmCreateOptions) SetMemory(value int) error {
	if value < 16 {
		return NewArgError("memory", "it must be >= 16")
	}
	c.optionsMap["memory"] = strconv.Itoa(value)
	return nil
}

// (0 - N) Set maximum tolerated downtime (in seconds) for migrations.
// default = 0.1
func (c *vmCreateOptions) SetMigrateDowntime(value int) error {
	if value < 1 {
		return NewArgError("migrate_downtime", "it must be >= 0")
	}
	c.optionsMap["migrate_downtime"] = strconv.Itoa(value)
	return nil
}

// Set maximum speed (in MB/s) for migrations. Value 0 is no limit.
// default = 0
func (c *vmCreateOptions) SetMigrateSpeed(value int) error {
	if value < 1 {
		return NewArgError("migrate_speed", "it must be >= 0")
	}
	c.optionsMap["migrate_speed"] = strconv.Itoa(value)
	return nil
}

// Set a name for the VM. Only used on the configuration web interface.
func (c *vmCreateOptions) SetName(name string) {
	c.optionsMap["name"] = name
}

// Specify network devices.
func (c *vmCreateOptions) AddNetworkDevice(number int, value string) error {
	key := fmt.Sprintf("net%d", number)
	if _, ok := c.optionsMap[key]; ok {
		return NewArgError("net[n]", fmt.Sprintf("Network device %s already exists", key))
	}

	c.optionsMap[key] = value
	return nil
}

// Enable/disable NUMA.
// default = 0
func (c *vmCreateOptions) SetNUMA(value bool) {
	c.optionsMap["numa"] = boolToString(value)
}

// NUMA topology.
func (c *vmCreateOptions) AddNUMA(value string) {
	c.optionsMap["numa1"] = value
}

// Specifies whether a VM will be started during system bootup.
// default = 0
func (c *vmCreateOptions) SetStartAtBoot(value bool) {
	c.optionsMap["onboot"] = boolToString(value)
}

// Specify guest operating system. This is used to enable special optimization/features for specific operating systems.
func (c *vmCreateOptions) SetOSType(value OSType) {
	c.optionsMap["ostype"] = value.String()
}

// Map host parallel devices (n is 0 to 2).
// NOTE: This option allows direct access to host hardware.
// So it is no longer possible to migrate such machines - use with special care.
func (c *vmCreateOptions) AddParallelDevice(number int, value string) error {
	if number < 0 || number > 2 {
		return NewArgError("parallel[n]", "it must be 0 to 2")
	}

	key := fmt.Sprintf("parallel%d", number)
	if _, ok := c.optionsMap[key]; ok {
		return NewArgError("parallel[n]", fmt.Sprintf("Parallel device %s already exists", key))
	}

	c.optionsMap[key] = value

	return nil
}

// Add theVM to the specified pool.
func (c *vmCreateOptions) SetPool(value string) {
	c.optionsMap["pool"] = value
}

// Sets the protection flag of the VM. This will disable the remove VM and remove disk operations.
func (c *vmCreateOptions) SetProtection(value bool) {
	c.optionsMap["protection"] = boolToString(value)
}

// Allow reboot. If set to '0' the VM exit on reboot.
func (c *vmCreateOptions) SetReboot(value bool) {
	c.optionsMap["reboot"] = boolToString(value)
}

// Use volume as SATA hard disk or CD-ROM (n is 0 to 5).
func (c *vmCreateOptions) AddSATADevice(number int, value string) error {
	if number < 0 || number > 5 {
		return NewArgError("sata[n]", "it must be 0 to 5")
	}

	key := fmt.Sprintf("sata%d", number)
	if _, ok := c.optionsMap[key]; ok {
		return NewArgError("sata[n]", fmt.Sprintf("SATA device %s already exists", key))
	}

	c.optionsMap[key] = value

	return nil
}

// Use volume as SCSI hard disk or CD-ROM (n is 0 to 13).
func (c *vmCreateOptions) AddSCSIDevice(number int, value string) error {
	if number < 0 || number > 13 {
		return NewArgError("scsi[n]", "it must be 0 to 13")
	}

	key := fmt.Sprintf("scsi%d", number)
	if _, ok := c.optionsMap[key]; ok {
		return NewArgError("scsi[n]", fmt.Sprintf("SCSI device %s already exists", key))
	}

	c.optionsMap[key] = value

	return nil
}

// SCSI controller model
// default = lsi
func (c *vmCreateOptions) SetSCSIControllerType(value SCSIControllerType) {
	c.optionsMap["scsihw"] = value.String()
}

// Create a serial device inside the VM (n is 0 to 3), and pass through a host serial device (i.e. /dev/ttyS0),
// or create a unix socket on the host side (use 'qm terminal' to open a terminal connection).
// NOTE: If you pass through a host serial device, it is no longer possible to migrate such machines - use with special care.
func (c *vmCreateOptions) AddSerialDevice(number int, value string) error {
	if number < 0 || number > 3 {
		return NewArgError("serial[n]", "it must be 0 to 3")
	}

	key := fmt.Sprintf("serial%d", number)
	if _, ok := c.optionsMap[key]; ok {
		return NewArgError("serial[n]", fmt.Sprintf("Serial device %s already exists", key))
	}

	c.optionsMap[key] = value

	return nil
}

// Amount of memory shares for auto-ballooning.
// The larger the number is, the more memory this VM gets.
// Number is relative to weights of all other running VMs. Using zero disables auto-ballooning
// default = 1000
func (c *vmCreateOptions) SetMemoryShares(value int) error {
	if value < 0 || value > 50000 {
		return NewArgError("shares", "it must be 0 to 50000")
	}
	c.optionsMap["shares"] = strconv.Itoa(value)
	return nil
}

// Specify SMBIOS type 1 fields.
func (c *vmCreateOptions) SetSMBIOS1(value string) {
	c.optionsMap["smbios1"] = value
}

// The number of CPUs. Please use option -sockets instead.
// default = 1
func (c *vmCreateOptions) SetSMP(value int) error {
	if value < 1 {
		return NewArgError("smp", "it must be > 0")
	}
	c.optionsMap["smp"] = strconv.Itoa(value)
	return nil
}

// The number of CPU sockets.
// default = 1
func (c *vmCreateOptions) SetCPUSockets(value int) error {
	if value < 1 {
		return NewArgError("sockets", "it must be > 0")
	}
	c.optionsMap["sockets"] = strconv.Itoa(value)
	return nil
}

// Set the initial date of the real time clock. Valid format for date are: 'now' or '2006-06-17T16:01:21' or'2006-06-17'.
// (now |YYYY-MM-DD | YYYY-MM-DDTHH:MM:SS)
// default = now
func (c *vmCreateOptions) SetStartDate(value string) {
	c.optionsMap["startdate"] = value
}

// Startup and shutdown behavior.
// Order is a non-negative number defining the general startup order. Shutdown is done with reverse ordering.
// Additionally you can set the 'up' or 'down' delay in seconds, which specifies a delay to wait before the next VM is started or stopped.
// [[order=]\d+] [,up=\d+] [,down=\d+]
func (c *vmCreateOptions) SetStartup(value string) {
	c.optionsMap["startup"] = value
}

// Default storage.
func (c *vmCreateOptions) SetStorage(value string) {
	c.optionsMap["storage"] = value
}

// Enable/disable the USB tablet device. This device is usually needed to allow absolute mouse positioning with VNC.
// Else the mouse runs out of sync with normal VNC clients. If you're running lots of console-only guests on one host,
// you may consider disabling this to save some context switches. This is turned off by default if you use spice (-vga=qxl).
// default = 1
func (c *vmCreateOptions) SetTablet(value bool) {
	c.optionsMap["tablet"] = boolToString(value)
}

// Enable/disable time drift fix.
// default = 0
func (c *vmCreateOptions) SetTDF(value bool) {
	c.optionsMap["tdf"] = boolToString(value)
}

// Enable/disable Template.
// default = 0
func (c *vmCreateOptions) SetTemplate(value bool) {
	c.optionsMap["template"] = boolToString(value)
}

// Assign a unique random ethernet address.
func (c *vmCreateOptions) SetUnique(value bool) {
	c.optionsMap["unique"] = boolToString(value)
}

// Configure an USB device (n is 0 to 4).
func (c *vmCreateOptions) AddUSBDevice(number int, value string) error {
	if number < 0 || number > 4 {
		return NewArgError("usb[n]", "it must be 0 to 4")
	}

	key := fmt.Sprintf("usb%d", number)
	if _, ok := c.optionsMap[key]; ok {
		return NewArgError("usb[n]", fmt.Sprintf("USB device %s already exists", key))
	}

	c.optionsMap[key] = value

	return nil
}

// Number of hotplugged vCPUs.
// default = 0
func (c *vmCreateOptions) SetVCPUs(value int) {
	c.optionsMap["vcpus"] = strconv.Itoa(value)
}

// Select the VGA type. If you want to use high resolution modes (>= 1280x1024x16) then you should use the options std or vmware.
// Default is std for win8/win7/w2k8, and cirrus for other OS types. The qxl option enables the SPICE display sever.
// For win* OS you can select how many independent displays you want, Linux guests can add displays them self.
// You can also run without any graphic card, using a serial device as terminal.
func (c *vmCreateOptions) SetVGA(value VGAType) {
	c.optionsMap["vga"] = value.String()
}

// Use volume as VirtIO hard disk (n is 0 to 15).
func (c *vmCreateOptions) AddVirtIODevice(number int, value string) error {
	if number < 0 || number > 15 {
		return NewArgError("virtio[n]", "it must be 0 to 15")
	}

	key := fmt.Sprintf("virtio%d", number)
	if _, ok := c.optionsMap[key]; ok {
		return NewArgError("virtio[n]", fmt.Sprintf("VirtIO device %s already exists", key))
	}

	c.optionsMap[key] = value

	return nil
}

// The (unique) ID of the VM.
func (c *vmCreateOptions) SetVMID(value int) error {
	if value < 0 {
		return NewArgError("vmid", "it can't be a negative number")
	}
	c.optionsMap["vmid"] = strconv.Itoa(value)
	return nil
}

// Create a virtual hardware watchdog device. Once enabled (by a guest action),
// the watchdog must be periodically polled by an agent inside the guest or else the watchdog will reset the guest
// (or execute the respective action specified)
func (c *vmCreateOptions) SetWatchdog(value string) {
	c.optionsMap["watchdog"] = value
}

type vmsRoot struct {
	VMs []VM `json:"data"`
}

type vmStatusRoot struct {
	VMStatus VMStatus `json:"data"`
}

// Virtual machine index (per node).
func (s *QemuServiceOp) GetVMs(node string) ([]VM, error) {
	path := fmt.Sprintf("nodes/%s/qemu", node)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	root := new(vmsRoot)
	if _, err = s.client.Do(req, root); err != nil {
		return nil, err
	}

	return root.VMs, err
}

// Get virtual machine status.
func (s *QemuServiceOp) GetVMCurrentStatus(node string, vmID int) (*VMStatus, error) {
	path := fmt.Sprintf("nodes/%s/qemu/%d/status/current", node, vmID)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	root := new(vmStatusRoot)
	if _, err = s.client.Do(req, root); err != nil {
		return nil, err
	}

	return &root.VMStatus, err
}

// Start virtual machine.
func (s *QemuServiceOp) StartVM(node string, vmID int) error {
	path := fmt.Sprintf("nodes/%s/qemu/%d/status/start", node, vmID)

	req, err := s.client.NewRequest("POST", path, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Stop virtual machine. The qemu process will exit immediately.
// This is akin to pulling the power plug of a running computer and may damage the VM data.
func (s *QemuServiceOp) StopVM(node string, vmID int) error {
	path := fmt.Sprintf("nodes/%s/qemu/%d/status/stop", node, vmID)

	req, err := s.client.NewRequest("POST", path, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Shutdown virtual machine. This is similar to pressing the power button on a physical machine.
// This will send an ACPI event for the guest OS, which should then proceed to a clean shutdown.
func (s *QemuServiceOp) ShutdownVM(node string, vmID int) error {
	path := fmt.Sprintf("nodes/%s/qemu/%d/status/shutdown", node, vmID)

	req, err := s.client.NewRequest("POST", path, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Reset virtual machine.
func (s *QemuServiceOp) ResetVM(node string, vmID int) error {
	path := fmt.Sprintf("nodes/%s/qemu/%d/status/reset", node, vmID)

	req, err := s.client.NewRequest("POST", path, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Suspend virtual machine.
func (s *QemuServiceOp) SuspendVM(node string, vmID int) error {
	path := fmt.Sprintf("nodes/%s/qemu/%d/status/suspend", node, vmID)

	req, err := s.client.NewRequest("POST", path, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Resume virtual machine.
func (s *QemuServiceOp) ResumeVM(node string, vmID int) error {
	path := fmt.Sprintf("nodes/%s/qemu/%d/status/resume", node, vmID)

	req, err := s.client.NewRequest("POST", path, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Create virtual machine.
func (s *QemuServiceOp) CreateVM(node string, createRequest VMCreateRequest) error {
	path := fmt.Sprintf("nodes/%s/qemu", node)
	req, err := s.client.NewRequest("POST", path, createRequest.GetOptionsMap())
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Create virtual machine.
func (s *QemuServiceOp) DeleteVM(node string, vmID int) error {
	path := fmt.Sprintf("nodes/%s/qemu/%d", node, vmID)

	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}
