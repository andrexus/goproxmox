package goproxmox

import "fmt"

type QemuService interface {
	GetVMs(node string) ([]VM, error)
	GetVMCurrentStatus(node string, vmID string) (*VMStatus, error)
	StartVM(node string, vmID string) error
	StopVM(node string, vmID string) error
	ShutdownVM(node string, vmID string) error
	ResetVM(node string, vmID string) error
	SuspendVM(node string, vmID string) error
	ResumeVM(node string, vmID string) error
	CreateVM(node string, vmID string, options *VMCreateOptions) error
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

type VMCreateOptions struct {
	VMID            string // The (unique) ID of the VM.
	Acpi            bool   // Enable/disable ACPI.
	Agent           bool   // Enable/disable Qemu GuestAgent.
	Archive         string // The backup file.
	Args            string // Arbitrary arguments passed to kvm, for example: args: -no-reboot -no-hpet
	Autostart       bool   // Automatic restart after crash
	Balloon         int    // (0 - N) Amount of target RAM for the VM in MB. Using zero disables the balloon driver.
	Bios            string // enum: seabios | ovmf. Select BIOS implementation.
	Boot            string // [acdn]{1,4} Boot on floppy (a), hard disk (c), CD-ROM (d), or network (n).
	Bootdisk        string // pve-qm-bootdisk Enable booting from specified disk.
	Cdrom           string // <volume> This is an alias for option -ide2
	Cores           int    // (1 - N) The number of cores per socket.
	Cpu             string // [cputype=]<enum> [,hidden=<1|0>] Emulated CPU type.
	Cpulimit        int    // (0 - 128) Limit of CPU usage. NOTE: If the computer has 2 CPUs, it has total of '2' CPU time. Value '0' indicates no CPU limit.
	Cpuunits        int    // (0 - 500000) Limit of CPU usage. NOTE: If the computer has 2 CPUs, it has total of '2' CPU time. Value '0' indicates no CPU limit.
	Description     string // Description for the VM. Only used on the configuration web interface. This is saved as comment inside the configuration file.
	Force           bool   // Allow to overwrite existing VM.
	Freeze          bool   // Freeze CPU at startup (use 'c' monitor command to start execution).
	Hostpci         string // Map host PCI devices into guest. NOTE: This option allows direct access to host hardware. So it is no longer possible to migrate such machines - use with special care.
	Hotplug         string // Selectively enable hotplug features. This is a comma separated list of hotplug features: 'network', 'disk', 'cpu', 'memory' and 'usb'. Use '0' to disable hotplug completely. Value '1' is an alias for the default 'network,disk,usb'.
	Hugepages       string // enum: any | 2 | 1024. Enable/disable hugepages memory.
	Ide0            string // Use volume as IDE hard disk or CD-ROM
	Ide1            string // Use volume as IDE hard disk or CD-ROM
	Ide2            string // Use volume as IDE hard disk or CD-ROM
	Ide3            string // Use volume as IDE hard disk or CD-ROM
	Keyboard        string // enum: de | de-ch| da | en-gb | en-us | es | fi | fr | fr-be | fr-ca | fr-ch | hu | is | it | ja | lt | mk | nl | no| pl | pt | pt-br | sv | sl | tr. Keybordlayout for vnc server. Default is read from the '/etc/pve/datacenter.conf' configuration file.
	Kvm             bool   // Enable/disable KVM hardware virtualization.
	Localtime       bool   // Set the real time clock to local time. This is enabled by default if ostype indicates a Microsoft OS.
	Lock            string // enum: migrate |backup | snapshot | rollback. Lock/unlock the VM.
	Machine         string // Specificthe Qemu machine type. (pc|pc(-i440fx)?-\d+\.\d+(\.pxe)?|q35|pc-q35-\d+\.\d+(\.pxe)?)
	Memory          int    // (16 - N) Amount of RAM for the VM in MB. This is the maximum available memory when you use the balloon device.
	MigrateDowntime int    // (0 - N) Set maximum tolerated downtime (in seconds) for migrations.
	MigrateSpeed    int    // (0 - N) Set maximum speed (in MB/s) for migrations. Value 0 is no limit.
	Name            string // Set aname for the VM. Only used on the configuration web interface.
	Net1            string // Specify network devices.
	Numa            bool   // Enable/disable NUMA.
	Numa1           string // NUMA topology.
	Onboot          bool   // Specifies whether a VM will be started during system bootup.
	Ostype          string // enum: other | wxp | w2k | w2k3 | w2k8 | wvista | win7 | win8 | win10 | l24 | l26 | solaris. Specify guest operating system. This is used to enable special optimization/features for specific operating systems.
	Parallel1       string // Map host parallel devices (n is 0 to 2). NOTE: This option allows direct access to host hardware. So it is no longer possible to migrate such machines - use with special care.
	Pool            string // Add theVM to the specified pool.
	Protection      bool   // Sets the protection flag of the VM. This will disable the remove VM and remove disk operations.
	Reboot          bool   // Allow reboot. If set to '0' the VM exit on reboot.
	Sata1           string // Use volume as SATA hard disk or CD-ROM (n is 0 to 5).
	Scsi1           string // Use volume as SCSI hard disk or CD-ROM (n is 0 to 13).
	Scsihw          string // enum: lsi |lsi53c810 | virtio-scsi-pci | virtio-scsi-single | megasas | pvscsi. SCSI controller model
	Serial1         string // (/dev/.+|socket) Create a serial device inside the VM (n is 0 to 3), and pass through a host serial device (i.e. /dev/ttyS0), or create a unix socket on the host side (use 'qm terminal' to open a terminal connection). NOTE: If you pass through a host serial device, it is no longer possible to migrate such machines - use with special care.
	Shares          int    // Amount of memory shares for auto-ballooning. The larger the number is, the more memory this VM gets. Number is relative to weights of all other running VMs. Using zero disables auto-ballooning
	Smbios1         string // Specify SMBIOS type 1 fields.
	Smp             int    // Then umber of CPUs. Please use option -sockets instead.
	Sockets         int    // The number of CPU sockets.
	Startdate       string // (now |YYYY-MM-DD | YYYY-MM-DDTHH:MM:SS) Set the initial date of the real time clock. Valid format for date are: 'now' or '2006-06-17T16:01:21' or'2006-06-17'.
	Startup         string // [[order=]\d+] [,up=\d+] [,down=\d+] Startup and shutdown behavior. Order is a non-negative number defining the general startup order. Shutdown is done with reverse ordering. Additionally you can set the 'up' or 'down' delay in seconds, which specifies a delay to wait before the next VM is started or stopped.
	Storage         string // Default storage.
	Tablet          bool   // Enable/disable the USB tablet device. This device is usually needed to allow absolute mouse positioning with VNC. Else the mouse runs out of sync with normal VNC clients. If you're running lots of console-only guests on one host, you may consider disabling this to save some context switches. This is turned off by default if you use spice (-vga=qxl).
	Tdf             bool   // Enable/disable time drift fix.
	Template        bool   // Enable/disable Template.
	Unique          bool   // Assign a unique random ethernet address.
	Unused1         string // Reference to unused volumes. This is used internally, and should not be modified manually.
	Usb1            string // Configure an USB device (n is 0 to 4).
	Vcpus           int    // Number of hotplugged vcpus.
	Vga             string // enum: std |cirrus | vmware | qxl | serial0 | serial1 | serial2 | serial3 | qxl2 | qxl3 | qxl4 Select the VGA type. If you want to use high resolution modes (&gt;= 1280x1024x16) then you should use the options 'std' or 'vmware'. Default is 'std' for win8/win7/w2k8, and 'cirrus' for other OS types. The'qxl' option enables the SPICE display sever. For win* OS you can select how many independent displays you want, Linux guests can add displays them self. You can also run without any graphic card, using a serial device as terminal.
	Virtio1         string // Use volume as VIRTIO hard disk (n is 0 to 15).
	Watchdog        string // Create a virtual hardware watchdog device. Once enabled (by a guest action), the watchdog must be periodically polled by an agent inside the guest or else the watchdog will reset the guest(or execute the respective action specified)
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
func (s *QemuServiceOp) CreateVM(node string, vmID string, options *VMCreateOptions) error {
	path := fmt.Sprintf("nodes/%s/qemu", node)

	options.VMID = vmID
	req, err := s.client.NewRequest("POST", path, options)
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
