package goproxmox

import (
	"fmt"
	"strconv"
)

type QemuService interface {
	GetVMs(node string) ([]VM, error)
	GetVMCurrentStatus(node string, vmID string) (*VMStatus, error)
	StartVM(node string, vmID string) error
	StopVM(node string, vmID string) error
	ShutdownVM(node string, vmID string) error
	ResetVM(node string, vmID string) error
	SuspendVM(node string, vmID string) error
	ResumeVM(node string, vmID string) error
	CreateVM(node string, createRequest VMCreateRequest) error
	DeleteVM(node string, vmID string) error
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
	SetACPI(value bool)
	SetAgent(value bool)
	SetArchive(value string)
	SetArgs(value string)
	SetAutostart(value bool)
	SetBalloon(value int) error
	SetBios(value string)
	SetBoot(value string)
	SetBootdisk(value string)
	SetCdrom(value string)
	SetCores(cores int) error
	SetCpu(value string)
	SetCpulimit(value int)
	SetCpuunits(value int)
	SetDescription(value string)
	SetForce(value bool)
	SetFreeze(value bool)
	SetHostpci(value string)
	SetHotplug(value string)
	SetHugepages(value string)
	AddIDE(value string) error
	SetKeyboard(value string)
	SetKvm(value bool)
	SetLocaltime(value bool)
	SetLock(value string)
	SetMachine(value string)
	SetMemory(value int)
	SetMigrateDowntime(value int)
	SetMigrateSpeed(value int)
	SetName(name string)
	AddNet(value string)
	SetNuma(value bool)
	AddNuma(value string)
	SetOnboot(value bool)
	SetOstype(value string)
	AddParallel(value string)
	SetPool(value string)
	SetProtection(value bool)
	SetReboot(value bool)
	AddSATA(value string)
	AddSCSI(value string)
	SetScsihw(value string)
	AddSerial(value string)
	SetShares(value int)
	SetSmbios1(value string)
	SetSmp(value int)
	SetSockets(value int)
	SetStartdate(value string)
	SetStartup(value string)
	SetStorage(value string)
	SetTablet(value bool)
	SetTdf(value bool)
	SetTemplate(value bool)
	SetUnique(value bool)
	AddUnused(value string)
	AddUSB(value string)
	SetVcpus(value int)
	SetVga(value string)
	AddVirtio(value string)
	SetVmID(vmID int)
	SetWatchdog(value string)
}

type vmCreateOptions struct {
	optionsMap map[string]string
}

func NewVMCreateRequest(vmID int) *vmCreateOptions {
	optionsMap := map[string]string{
		"vmid": strconv.Itoa(vmID),
	}
	return &vmCreateOptions{optionsMap: optionsMap}
}

func (c *vmCreateOptions) GetOptionsMap() map[string]string {
	return c.optionsMap
}

// Enable/disable ACPI.
func (c *vmCreateOptions) SetACPI(value bool) {
	c.optionsMap["acpi"] = strconv.FormatBool(value)
}

// Enable/disable Qemu GuestAgent.
func (c *vmCreateOptions) SetAgent(value bool) {
	c.optionsMap["agent"] = strconv.FormatBool(value)
}

// The backup file.
func (c *vmCreateOptions) SetArchive(value string) {
	c.optionsMap["archive"] = value
}

// Arbitrary arguments passed to kvm, for example: args: -no-reboot -no-hpet
func (c *vmCreateOptions) SetArgs(value string) {
	c.optionsMap["args"] = value
}

// Automatic restart after crash
func (c *vmCreateOptions) SetAutostart(value bool) {
	c.optionsMap["autostart"] = strconv.FormatBool(value)
}

// (0 - N) Amount of target RAM for the VM in MB. Using zero disables the balloon driver.
func (c *vmCreateOptions) SetBalloon(value int) error {
	if value < 0 {
		return NewArgError("balloon", "it can't be negative number")
	}
	c.optionsMap["balloon"] = strconv.Itoa(value)
	return nil
}

// enum: seabios | ovmf. Select BIOS implementation.
func (c *vmCreateOptions) SetBios(value string) {
	c.optionsMap["bios"] = value
}

// [acdn]{1,4} Boot on floppy (a), hard disk (c), CD-ROM (d), or network (n).
func (c *vmCreateOptions) SetBoot(value string) {
	c.optionsMap["boot"] = value
}

// pve-qm-bootdisk Enable booting from specified disk.
func (c *vmCreateOptions) SetBootdisk(value string) {
	c.optionsMap["bootdisk"] = value
}

// <volume> This is an alias for option -ide2
func (c *vmCreateOptions) SetCdrom(value string) {
	c.optionsMap["cdrom"] = value
}

// (1 - N) The number of cores per socket.
func (c *vmCreateOptions) SetCores(value int) error {
	if value < 1 {
		return NewArgError("cores", "it must be > 0")
	}
	c.optionsMap["cores"] = strconv.Itoa(value)
	return nil
}

// [cputype=]<enum> [,hidden=<1|0>] Emulated CPU type.
func (c *vmCreateOptions) SetCpu(value string) {
	c.optionsMap["cpu"] = value
}

// (0 - 128) Limit of CPU usage. NOTE: If the computer has 2 CPUs, it has total of '2' CPU time. Value '0' indicates no CPU limit.
func (c *vmCreateOptions) SetCpulimit(value int) {
	c.optionsMap["cpulimit"] = strconv.Itoa(value)
}

// (0 - 500000) Limit of CPU usage. NOTE: If the computer has 2 CPUs, it has total of '2' CPU time. Value '0' indicates no CPU limit.
func (c *vmCreateOptions) SetCpuunits(value int) {
	c.optionsMap["cpuunits"] = strconv.Itoa(value)
}

// Description for the VM. Only used on the configuration web interface. This is saved as comment inside the configuration file.
func (c *vmCreateOptions) SetDescription(value string) {
	c.optionsMap["description"] = value
}

// Allow to overwrite existing VM.
func (c *vmCreateOptions) SetForce(value bool) {
	c.optionsMap["force"] = strconv.FormatBool(value)
}

// Freeze CPU at startup (use 'c' monitor command to start execution).
func (c *vmCreateOptions) SetFreeze(value bool) {
	c.optionsMap["freeze"] = strconv.FormatBool(value)
}

// Map host PCI devices into guest. NOTE: This option allows direct access to host hardware. So it is no longer possible to migrate such machines - use with special care.
func (c *vmCreateOptions) SetHostpci(value string) {
	c.optionsMap["hostpci"] = value
}

// Selectively enable hotplug features. This is a comma separated list of hotplug features: 'network', 'disk', 'cpu', 'memory' and 'usb'. Use '0' to disable hotplug completely. Value '1' is an alias for the default 'network,disk,usb'.
func (c *vmCreateOptions) SetHotplug(value string) {
	c.optionsMap["hotplug"] = value
}

// enum: any | 2 | 1024. Enable/disable hugepages memory.
func (c *vmCreateOptions) SetHugepages(value string) {
	c.optionsMap["hugepages"] = value
}

// Use volume as IDE hard disk or CD-ROM
func (c *vmCreateOptions) AddIDE(value string) error {
	return nil
}

// enum: de | de-ch| da | en-gb | en-us | es | fi | fr | fr-be | fr-ca | fr-ch | hu | is | it | ja | lt | mk | nl | no| pl | pt | pt-br | sv | sl | tr. Keybordlayout for vnc server. Default is read from the '/etc/pve/datacenter.conf' configuration file.
func (c *vmCreateOptions) SetKeyboard(value string) {
	c.optionsMap["keyboard"] = value
}

// Enable/disable KVM hardware virtualization.
func (c *vmCreateOptions) SetKvm(value bool) {
	c.optionsMap["kvm"] = strconv.FormatBool(value)
}

// Set the real time clock to local time. This is enabled by default if ostype indicates a Microsoft OS.
func (c *vmCreateOptions) SetLocaltime(value bool) {
	c.optionsMap["localtime"] = strconv.FormatBool(value)
}

// enum: migrate |backup | snapshot | rollback. Lock/unlock the VM.
func (c *vmCreateOptions) SetLock(value string) {
	c.optionsMap["lock"] = value
}

// Specificthe Qemu machine type. (pc|pc(-i440fx)?-\d+\.\d+(\.pxe)?|q35|pc-q35-\d+\.\d+(\.pxe)?)
func (c *vmCreateOptions) SetMachine(value string) {
	c.optionsMap["machine"] = value
}

// (16 - N) Amount of RAM for the VM in MB. This is the maximum available memory when you use the balloon device.
func (c *vmCreateOptions) SetMemory(value int) {
	c.optionsMap["memory"] = strconv.Itoa(value)
}

// (0 - N) Set maximum tolerated downtime (in seconds) for migrations.
func (c *vmCreateOptions) SetMigrateDowntime(value int) {
	c.optionsMap["migratedowntime"] = strconv.Itoa(value)
}

// (0 - N) Set maximum speed (in MB/s) for migrations. Value 0 is no limit.
func (c *vmCreateOptions) SetMigrateSpeed(value int) {
	c.optionsMap["migratespeed"] = strconv.Itoa(value)
}

// Set a name for the VM. Only used on the configuration web interface.
func (c *vmCreateOptions) SetName(name string) {
	c.optionsMap["name"] = name
}

// Specify network devices.
func (c *vmCreateOptions) AddNet(value string) {
	c.optionsMap["net1"] = value
}

// Enable/disable NUMA.
func (c *vmCreateOptions) SetNuma(value bool) {
	c.optionsMap["numa"] = strconv.FormatBool(value)
}

// NUMA topology.
func (c *vmCreateOptions) AddNuma(value string) {
	c.optionsMap["numa1"] = value
}

// Specifies whether a VM will be started during system bootup.
func (c *vmCreateOptions) SetOnboot(value bool) {
	c.optionsMap["onboot"] = strconv.FormatBool(value)
}

// enum: other | wxp | w2k | w2k3 | w2k8 | wvista | win7 | win8 | win10 | l24 | l26 | solaris. Specify guest operating system. This is used to enable special optimization/features for specific operating systems.
func (c *vmCreateOptions) SetOstype(value string) {
	c.optionsMap["ostype"] = value
}

// Map host parallel devices (n is 0 to 2). NOTE: This option allows direct access to host hardware. So it is no longer possible to migrate such machines - use with special care.
func (c *vmCreateOptions) AddParallel(value string) {
	c.optionsMap["parallel1"] = value
}

// Add theVM to the specified pool.
func (c *vmCreateOptions) SetPool(value string) {
	c.optionsMap["pool"] = value
}

// Sets the protection flag of the VM. This will disable the remove VM and remove disk operations.
func (c *vmCreateOptions) SetProtection(value bool) {
	c.optionsMap["protection"] = strconv.FormatBool(value)
}

// Allow reboot. If set to '0' the VM exit on reboot.
func (c *vmCreateOptions) SetReboot(value bool) {
	c.optionsMap["reboot"] = strconv.FormatBool(value)
}

// Use volume as SATA hard disk or CD-ROM (n is 0 to 5).
func (c *vmCreateOptions) AddSATA(value string) {
	c.optionsMap["sata1"] = value
}

// Use volume as SCSI hard disk or CD-ROM (n is 0 to 13).
func (c *vmCreateOptions) AddSCSI(value string) {
	c.optionsMap["scsi1"] = value
}

// enum: lsi |lsi53c810 | virtio-scsi-pci | virtio-scsi-single | megasas | pvscsi. SCSI controller model
func (c *vmCreateOptions) SetScsihw(value string) {
	c.optionsMap["scsihw"] = value
}

// (/dev/.+|socket) Create a serial device inside the VM (n is 0 to 3), and pass through a host serial device (i.e. /dev/ttyS0), or create a unix socket on the host side (use 'qm terminal' to open a terminal connection). NOTE: If you pass through a host serial device, it is no longer possible to migrate such machines - use with special care.
func (c *vmCreateOptions) AddSerial(value string) {
	c.optionsMap["serial1"] = value
}

// Amount of memory shares for auto-ballooning. The larger the number is, the more memory this VM gets. Number is relative to weights of all other running VMs. Using zero disables auto-ballooning
func (c *vmCreateOptions) SetShares(value int) {
	c.optionsMap["shares"] = strconv.Itoa(value)
}

// Specify SMBIOS type 1 fields.
func (c *vmCreateOptions) SetSmbios1(value string) {
	c.optionsMap["smbios1"] = value
}

// Then umber of CPUs. Please use option -sockets instead.
func (c *vmCreateOptions) SetSmp(value int) {
	c.optionsMap["smp"] = strconv.Itoa(value)
}

// The number of CPU sockets.
func (c *vmCreateOptions) SetSockets(value int) {
	c.optionsMap["sockets"] = strconv.Itoa(value)
}

// (now |YYYY-MM-DD | YYYY-MM-DDTHH:MM:SS) Set the initial date of the real time clock. Valid format for date are: 'now' or '2006-06-17T16:01:21' or'2006-06-17'.
func (c *vmCreateOptions) SetStartdate(value string) {
	c.optionsMap["startdate"] = value
}

// [[order=]\d+] [,up=\d+] [,down=\d+] Startup and shutdown behavior. Order is a non-negative number defining the general startup order. Shutdown is done with reverse ordering. Additionally you can set the 'up' or 'down' delay in seconds, which specifies a delay to wait before the next VM is started or stopped.
func (c *vmCreateOptions) SetStartup(value string) {
	c.optionsMap["startup"] = value
}

// Default storage.
func (c *vmCreateOptions) SetStorage(value string) {
	c.optionsMap["storage"] = value
}

// Enable/disable the USB tablet device. This device is usually needed to allow absolute mouse positioning with VNC. Else the mouse runs out of sync with normal VNC clients. If you're running lots of console-only guests on one host, you may consider disabling this to save some context switches. This is turned off by default if you use spice (-vga=qxl).
func (c *vmCreateOptions) SetTablet(value bool) {
	c.optionsMap["tablet"] = strconv.FormatBool(value)
}

// Enable/disable time drift fix.
func (c *vmCreateOptions) SetTdf(value bool) {
	c.optionsMap["tdf"] = strconv.FormatBool(value)
}

// Enable/disable Template.
func (c *vmCreateOptions) SetTemplate(value bool) {
	c.optionsMap["template"] = strconv.FormatBool(value)
}

// Assign a unique random ethernet address.
func (c *vmCreateOptions) SetUnique(value bool) {
	c.optionsMap["unique"] = strconv.FormatBool(value)
}

// Reference to unused volumes. This is used internally, and should not be modified manually.
func (c *vmCreateOptions) AddUnused(value string) {
	c.optionsMap["unused1"] = value
}

// Configure an USB device (n is 0 to 4).
func (c *vmCreateOptions) AddUSB(value string) {
	c.optionsMap["usb1"] = value
}

// Number of hotplugged vcpus.
func (c *vmCreateOptions) SetVcpus(value int) {
	c.optionsMap["vcpus"] = strconv.Itoa(value)
}

// enum: std |cirrus | vmware | qxl | serial0 | serial1 | serial2 | serial3 | qxl2 | qxl3 | qxl4 Select the VGA type. If you want to use high resolution modes (&gt;= 1280x1024x16) then you should use the options 'std' or 'vmware'. Default is 'std' for win8/win7/w2k8, and 'cirrus' for other OS types. The'qxl' option enables the SPICE display sever. For win* OS you can select how many independent displays you want, Linux guests can add displays them self. You can also run without any graphic card, using a serial device as terminal.
func (c *vmCreateOptions) SetVga(value string) {
	c.optionsMap["vga"] = value
}

// Use volume as VIRTIO hard disk (n is 0 to 15).
func (c *vmCreateOptions) AddVirtio(value string) {
	c.optionsMap["virtio1"] = value
}

// The (unique) ID of the VM.
func (c *vmCreateOptions) SetVmID(vmID int) {
	c.optionsMap["vmid"] = strconv.Itoa(vmID)
}

// Create a virtual hardware watchdog device. Once enabled (by a guest action), the watchdog must be periodically polled by an agent inside the guest or else the watchdog will reset the guest(or execute the respective action specified)
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
func (s *QemuServiceOp) GetVMCurrentStatus(node string, vmID string) (*VMStatus, error) {
	path := fmt.Sprintf("nodes/%s/qemu/%s/status/current", node, vmID)

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
func (s *QemuServiceOp) StartVM(node string, vmID string) error {
	path := fmt.Sprintf("nodes/%s/qemu/%s/status/start", node, vmID)

	req, err := s.client.NewRequest("POST", path, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Stop virtual machine. The qemu process will exit immediately.
// This is akin to pulling the power plug of a running computer and may damage the VM data.
func (s *QemuServiceOp) StopVM(node string, vmID string) error {
	path := fmt.Sprintf("nodes/%s/qemu/%s/status/stop", node, vmID)

	req, err := s.client.NewRequest("POST", path, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Shutdown virtual machine. This is similar to pressing the power button on a physical machine.
// This will send an ACPI event for the guest OS, which should then proceed to a clean shutdown.
func (s *QemuServiceOp) ShutdownVM(node string, vmID string) error {
	path := fmt.Sprintf("nodes/%s/qemu/%s/status/shutdown", node, vmID)

	req, err := s.client.NewRequest("POST", path, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Reset virtual machine.
func (s *QemuServiceOp) ResetVM(node string, vmID string) error {
	path := fmt.Sprintf("nodes/%s/qemu/%s/status/reset", node, vmID)

	req, err := s.client.NewRequest("POST", path, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Suspend virtual machine.
func (s *QemuServiceOp) SuspendVM(node string, vmID string) error {
	path := fmt.Sprintf("nodes/%s/qemu/%s/status/suspend", node, vmID)

	req, err := s.client.NewRequest("POST", path, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Resume virtual machine.
func (s *QemuServiceOp) ResumeVM(node string, vmID string) error {
	path := fmt.Sprintf("nodes/%s/qemu/%s/status/resume", node, vmID)

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
func (s *QemuServiceOp) DeleteVM(node string, vmID string) error {
	path := fmt.Sprintf("nodes/%s/qemu/%s", node, vmID)

	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}
