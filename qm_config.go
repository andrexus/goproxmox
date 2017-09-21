package goproxmox

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"log"
)

const (
	parameterACPI                      = "acpi"
	parameterQemuAgent                 = "agent"
	parameterArchive                   = "archive"
	parameterArgs                      = "args"
	parameterAutoStart                 = "autostart"
	parameterBalloon                   = "balloon"
	parameterBios                      = "bios"
	parameterBootOrder                 = "boot"
	parameterBootDisk                  = "bootdisk"
	parameterCDROM                     = "cdrom"
	parameterCores                     = "cores"
	parameterCPU                       = "cpu"
	parameterCPULimit                  = "cpulimit"
	parameterCPUUnits                  = "cpuunits"
	parameterDescription               = "description"
	parameterForce                     = "force"
	parameterFreeze                    = "freeze"
	parameterHostPCI                   = "hostpci"
	parameterHotPlug                   = "hotplug"
	parameterHugePages                 = "hugepages"
	parameterIDEDevices                = "ide"
	parameterKeyboardLayout            = "keyboard"
	parameterKVMHardwareVirtualization = "kvm"
	parameterLocalTime                 = "localtime"
	parameterLock                      = "lock"
	parameterMachineType               = "machine"
	parameterMemory                    = "memory"
	parameterMigrateDowntime           = "migrate_downtime"
	parameterMigrateSpeed              = "migrate_speed"
	parameterName                      = "name"
	parameterNetworkDevices            = "net"
	parameterNUMA                      = "numa"
	parameterNUMATopologies            = "numa"
	parameterStartAtBoot               = "onboot"
	parameterOSType                    = "ostype"
	parameterParallelDevices           = "parallel"
	parameterPool                      = "pool"
	parameterProtection                = "protection"
	parameterReboot                    = "reboot"
	parameterSATADevices               = "sata"
	parameterSCSIDevices               = "scsi"
	parameterSCSIControllerType        = "scsihw"
	parameterSerialDevices             = "serial"
	parameterMemoryShares              = "shares"
	parameterSMBIOS1                   = "smbios1"
	parameterSMP                       = "smp"
	parameterSockets                   = "sockets"
	parameterStartDate                 = "startdate"
	parameterStartup                   = "startup"
	parameterStorage                   = "storage"
	parameterTablet                    = "tablet"
	parameterTDF                       = "tdf"
	parameterTemplate                  = "template"
	parameterUnique                    = "unique"
	parameterUSBDevices                = "usb"
	parameterVCPUs                     = "vcpus"
	parameterVGA                       = "vga"
	parameterVirtIODevices             = "virtio"
	parameterVMID                      = "vmid"
	parameterWatchdog                  = "watchdog"
)

var parameterMapping = map[*regexp.Regexp]string{
	regexp.MustCompile(parameterACPI):                      "ACPI",
	regexp.MustCompile(parameterQemuAgent):                 "QemuAgent",
	regexp.MustCompile(parameterArchive):                   "Archive",
	regexp.MustCompile(parameterArgs):                      "Args",
	regexp.MustCompile(parameterAutoStart):                 "AutoStart",
	regexp.MustCompile(parameterBalloon):                   "Balloon",
	regexp.MustCompile(parameterBios):                      "Bios",
	regexp.MustCompile(parameterBootOrder):                 "BootOrder",
	regexp.MustCompile(parameterBootDisk):                  "BootDisk",
	regexp.MustCompile(parameterCDROM):                     "CDROM",
	regexp.MustCompile(parameterCores):                     "Cores",
	regexp.MustCompile(parameterCPU):                       "CPU",
	regexp.MustCompile(parameterCPULimit):                  "CPULimit",
	regexp.MustCompile(parameterCPUUnits):                  "CPUUnits",
	regexp.MustCompile(parameterDescription):               "Description",
	regexp.MustCompile(parameterForce):                     "Force",
	regexp.MustCompile(parameterFreeze):                    "Freeze",
	regexp.MustCompile(parameterHostPCI):                   "HostPCI",
	regexp.MustCompile(parameterHotPlug):                   "HotPlug",
	regexp.MustCompile(parameterHugePages):                 "HugePages",
	regexp.MustCompile(`ide(\d+)`):                         "IDEDevices",
	regexp.MustCompile(parameterKeyboardLayout):            "KeyboardLayout",
	regexp.MustCompile(parameterKVMHardwareVirtualization): "KVMHardwareVirtualization",
	regexp.MustCompile(parameterLocalTime):                 "LocalTime",
	regexp.MustCompile(parameterLock):                      "Lock",
	regexp.MustCompile(parameterMachineType):               "MachineType",
	regexp.MustCompile(parameterMemory):                    "Memory",
	regexp.MustCompile(parameterMigrateDowntime):           "MigrateDowntime",
	regexp.MustCompile(parameterMigrateSpeed):              "MigrateSpeed",
	regexp.MustCompile(parameterName):                      "Name",
	regexp.MustCompile(`net(\d+)`):                         "NetworkDevices",
	regexp.MustCompile(parameterNUMA):                      "NUMA",
	regexp.MustCompile(`numa(\d+)`):                        "NUMATopologies",
	regexp.MustCompile(parameterStartAtBoot):               "StartAtBoot",
	regexp.MustCompile(parameterOSType):                    "OSType",
	regexp.MustCompile(`parallel(\d+)`):                    "ParallelDevices",
	regexp.MustCompile(parameterPool):                      "Pool",
	regexp.MustCompile(parameterProtection):                "Protection",
	regexp.MustCompile(parameterReboot):                    "Reboot",
	regexp.MustCompile(`sata(\d+)`):                        "SATADevices",
	regexp.MustCompile(`scsi(\d+)`):                        "SCSIDevices",
	regexp.MustCompile(parameterSCSIControllerType):        "SCSIControllerType",
	regexp.MustCompile(`serial(\d+)`):                      "SerialDevices",
	regexp.MustCompile(parameterMemoryShares):              "MemoryShares",
	regexp.MustCompile(parameterSMBIOS1):                   "SMBIOS1",
	regexp.MustCompile(parameterSMP):                       "SMP",
	regexp.MustCompile(parameterSockets):                   "Sockets",
	regexp.MustCompile(parameterStartDate):                 "StartDate",
	regexp.MustCompile(parameterStartup):                   "Startup",
	regexp.MustCompile(parameterStorage):                   "Storage",
	regexp.MustCompile(parameterTablet):                    "Tablet",
	regexp.MustCompile(parameterTDF):                       "TDF",
	regexp.MustCompile(parameterTemplate):                  "Template",
	regexp.MustCompile(parameterUnique):                    "Unique",
	regexp.MustCompile(`usb(\d+)`):                         "USBDevices",
	regexp.MustCompile(parameterVCPUs):                     "VCPUs",
	regexp.MustCompile(parameterVGA):                       "VGA",
	regexp.MustCompile(`virtio(\d+)`):                      "VirtIODevices",
	regexp.MustCompile(parameterVMID):                      "VMID",
	regexp.MustCompile(parameterWatchdog):                  "Watchdog",
}

type VMConfig struct {
	// Enable/disable ACPI.
	// default = 1
	ACPI *bool

	// Enable/disable Qemu GuestAgent.
	// default = 0
	QemuAgent *bool

	//
	// The backup file.
	Archive *string // TODO Create only

	//
	// Arbitrary arguments passed to kvm, for example:
	// args: -no-reboot -no-hpet
	Args *string

	// Automatic restart after crash
	// default = 0
	AutoStart *bool

	// background_delay. Only post

	//
	// Amount of target RAM for the VM in MB. Using zero disables the balloon driver.
	Balloon *int

	//
	// Select BIOS implementation.
	// default = seabios
	Bios *Bios

	//
	// Boot on floppy (a), hard disk (c), CD-ROM (d), or network (n).
	// default = cdn
	BootOrder []BootDevice

	//
	// Enable booting from specified disk.
	// (ide|sata|scsi|virtio)\d+
	BootDisk *string

	//
	// <volume> This is an alias for option -ide2
	CDROM *string

	//
	// The number of cores per socket.
	// default = 1
	Cores *int

	//
	// [cputype=]<enum> [,hidden=<1|0>] Emulated CPU type.
	// default = kvm64
	// hidden=<boolean> Do not identify as a KVM virtual machine.
	// default = 0
	CPU *CPUType

	//
	// (0 - 128) Limit of CPU usage. NOTE: If the computer has 2 CPUs, it has total of '2' CPU time. Value '0' indicates no CPU limit.
	// default = 0
	CPULimit *int

	//
	// (0 - 500000) CPU weight for a VM. Argument is used in the kernel fair scheduler.
	// The larger the number is, the more CPU time this VM gets.
	// Number is relative to weights of all the other running VMs.
	// You can disable fair-scheduler configuration by setting this to 0.
	// default = 1024
	CPUUnits *int
	//

	// TODO delete. Only post/put

	//
	// Description for the VM. Only used on the configuration web interface.
	// This is saved as comment inside the configuration file.
	Description *string
	//
	// TODO digest. Only post/put
	//
	// Allow to overwrite existing VM.
	Force *bool
	//

	//
	// Freeze CPU at startup (use 'c' monitor command to start execution).
	Freeze *bool

	//
	// Map host PCI devices into guest.
	// NOTE: This option allows direct access to host hardware.
	// So it is no longer possible to migrate such machines - use with special care.
	HostPCI *string

	//
	// Selectively enable hotplug features.
	// This is a comma separated list of hotplug features: 'network', 'disk', 'cpu', 'memory' and 'usb'.
	// Use '0' to disable hotplug completely.
	// Value '1' is an alias for the default 'network,disk,usb'.
	HotPlug *string

	//
	// Enable/disable hugepages memory.
	HugePages *HugePages

	//
	// Use volume as IDE hard disk or CD-ROM
	IDEDevices map[int]*IDEDevice

	//
	// Keyboard layout for vnc server. Default is read from the '/etc/pve/datacenter.conf' configuration file.
	// default = en-us
	KeyboardLayout *KeyboardLayout

	//
	// Enable/disable KVM hardware virtualization.
	// default = 1
	KVMHardwareVirtualization *bool

	//
	// Set the real time clock to local time. This is enabled by default if ostype indicates a Microsoft OS.
	LocalTime *bool

	//
	// Lock/unlock the VM.
	Lock *Lock

	//
	// Specify the Qemu machine type.
	// (pc|pc(-i440fx)?-\d+\.\d+(\.pxe)?|q35|pc-q35-\d+\.\d+(\.pxe)?)
	MachineType *string

	//
	// Amount of RAM for the VM in MB. This is the maximum available memory when you use the balloon device.
	// default = 512
	Memory *int

	//
	// (0 - N) Set maximum tolerated downtime (in seconds) for migrations.
	// default = 0.1
	MigrateDowntime *int

	//
	// Set maximum speed (in MB/s) for migrations. Value 0 is no limit.
	// default = 0
	MigrateSpeed *int

	//
	// Set a name for the VM. Only used on the configuration web interface.
	Name *string

	//
	// Specify network devices.
	NetworkDevices map[int]*NetworkDevice

	//
	// Enable/disable NUMA.
	// default = 0
	NUMA *bool

	//
	// NUMA topology.
	NUMATopologies map[int]QMOption

	//
	// Specifies whether a VM will be started during system bootup.
	// default = 0
	StartAtBoot *bool

	//
	// Specify guest operating system. This is used to enable special optimization/features for specific operating systems.
	OSType *OSType

	//
	// Map host parallel devices (n is 0 to 2).
	// NOTE: This option allows direct access to host hardware.
	// So it is no longer possible to migrate such machines - use with special care.
	ParallelDevices map[int]QMOption

	//
	// Add theVM to the specified pool.
	Pool *string // TODO Create only

	//
	// Sets the protection flag of the VM. This will disable the remove VM and remove disk operations.
	Protection *bool

	//
	// Allow reboot. If set to '0' the VM exit on reboot.
	Reboot *bool

	//
	// TODO revert. Only post/put

	//
	// Use volume as SATA hard disk or CD-ROM (n is 0 to 5).
	SATADevices map[int]QMOption

	//
	// Use volume as SCSI hard disk or CD-ROM (n is 0 to 13).
	SCSIDevices map[int]QMOption

	//
	// SCSI controller model
	// default = lsi
	SCSIControllerType *SCSIControllerType

	//
	// Create a serial device inside the VM (n is 0 to 3), and pass through a host serial device (i.e. /dev/ttyS0),
	// or create a unix socket on the host side (use 'qm terminal' to open a terminal connection).
	// NOTE: If you pass through a host serial device, it is no longer possible to migrate such machines - use with special care.
	SerialDevices map[int]QMOption

	//
	// Amount of memory shares for auto-ballooning.
	// The larger the number is, the more memory this VM gets.
	// Number is relative to weights of all other running VMs. Using zero disables auto-ballooning
	// default = 1000
	MemoryShares *int

	//
	// TODO skiplock. Only post/put

	//
	// Specify SMBIOS type 1 fields.
	SMBIOS1 *string

	//
	// The number of CPUs. Please use option -sockets instead.
	// default = 1
	SMP *int

	//
	// The number of CPU sockets.
	// default = 1
	Sockets *int

	//
	// Set the initial date of the real time clock. Valid format for date are: 'now' or '2006-06-17T16:01:21' or'2006-06-17'.
	// (now |YYYY-MM-DD | YYYY-MM-DDTHH:MM:SS)
	// default = now
	StartDate *string

	//
	// Startup and shutdown behavior.
	// Order is a non-negative number defining the general startup order. Shutdown is done with reverse ordering.
	// Additionally you can set the 'up' or 'down' delay in seconds, which specifies a delay to wait before the next VM is started or stopped.
	// [[order=]\d+] [,up=\d+] [,down=\d+]
	Startup *string

	//
	// Default storage.
	Storage *string // TODO Create only

	//
	// Enable/disable the USB tablet device. This device is usually needed to allow absolute mouse positioning with VNC.
	// Else the mouse runs out of sync with normal VNC clients. If you're running lots of console-only guests on one host,
	// you may consider disabling this to save some context switches. This is turned off by default if you use spice (-vga=qxl).
	// default = 1
	Tablet *bool

	//
	// Enable/disable time drift fix.
	// default = 0
	TDF *bool

	//
	// Enable/disable Template.
	// default = 0
	Template *bool

	//
	// Assign a unique random ethernet address.
	Unique *bool // TODO Create only

	//
	// Configure an USB device (n is 0 to 4).
	USBDevices map[int]QMOption

	//
	// Number of hotplugged vCPUs.
	// default = 0
	VCPUs *int

	//
	// Select the VGA type. If you want to use high resolution modes (>= 1280x1024x16) then you should use the options std or vmware.
	// Default is std for win8/win7/w2k8, and cirrus for other OS types. The qxl option enables the SPICE display sever.
	// For win* OS you can select how many independent displays you want, Linux guests can add displays them self.
	// You can also run without any graphic card, using a serial device as terminal.
	VGAType *VGAType

	//
	// Use volume as VirtIO hard disk (n is 0 to 15).
	VirtIODevices map[int]*VirtIODevice

	//
	// The (unique) ID of the VM.
	VMID *int // TODO Create only

	//
	// Create a virtual hardware watchdog device. Once enabled (by a guest action),
	// the watchdog must be periodically polled by an agent inside the guest or else the watchdog will reset the guest
	// (or execute the respective action specified)
	Watchdog *string
}

func NewVMConfigFromMap(data map[string]interface{}) *VMConfig {
	config := new(VMConfig)
	for k, v := range data {
		fieldName, matchResults, err := findFieldName(k)
		if err != nil {
			continue
		}
		switch fieldName {
		case "ACPI":
			config.ACPI = Bool(intToBool(int(v.(float64))))
		case "QemuAgent":
			config.QemuAgent = Bool(intToBool(int(v.(float64))))
		case "AutoStart":
			config.AutoStart = Bool(intToBool(int(v.(float64))))
		case "Bios":
			v, err := BiosFromString(v.(string))
			if err == nil {
				config.Bios = &v
			}
		case "BootOrder":
			log.Printf("[DEBUG] Field %s is not supported yet", fieldName)
		case "CPU":
			v, err := CPUTypeFromString(v.(string))
			if err == nil {
				config.CPU = &v
			}
		case "Force":
			config.Force = Bool(intToBool(int(v.(float64))))
		case "Freeze":
			config.Freeze = Bool(intToBool(int(v.(float64))))
		case "HugePages":
			v, err := HugePagesFromString(v.(string))
			if err == nil {
				config.HugePages = &v
			}
		case "IDEDevices":
			number, _ := strconv.Atoi(matchResults[1])
			config.AddIDEDevice(number, NewIDEDeviceFromString(v.(string)))
		case "KeyboardLayout":
			v, err := KeyboardLayoutFromString(v.(string))
			if err == nil {
				config.KeyboardLayout = &v
			}
		case "KVMHardwareVirtualization":
			config.KVMHardwareVirtualization = Bool(intToBool(int(v.(float64))))
		case "LocalTime":
			config.LocalTime = Bool(intToBool(int(v.(float64))))
		case "Lock":
			v, err := LockFromString(v.(string))
			if err == nil {
				config.Lock = &v
			}
		case "NetworkDevices":
			number, _ := strconv.Atoi(matchResults[1])
			config.AddNetworkDevice(number, NewNetworkDeviceFromString(v.(string)))
		case "NUMA":
			config.NUMA = Bool(intToBool(int(v.(float64))))
		case "NUMATopologies":
			log.Printf("[DEBUG] Field %s is not supported yet", fieldName)
		case "StartAtBoot":
			config.StartAtBoot = Bool(intToBool(int(v.(float64))))
		case "ParallelDevices":
			log.Printf("[DEBUG] Field %s is not supported yet", fieldName)
		case "OSType":
			v, err := OSTypeFromString(v.(string))
			if err == nil {
				config.OSType = &v
			}
		case "Protection":
			config.Protection = Bool(intToBool(int(v.(float64))))
		case "Reboot":
			config.Reboot = Bool(intToBool(int(v.(float64))))
		case "SATADevices":
			log.Printf("[DEBUG] Field %s is not supported yet", fieldName)
		case "SCSIDevices":
			log.Printf("[DEBUG] Field %s is not supported yet", fieldName)
		case "SCSIControllerType":
			v, err := SCSIControllerTypeFromString(v.(string))
			if err == nil {
				config.SCSIControllerType = &v
			}
		case "SerialDevices":
			log.Printf("[DEBUG] Field %s is not supported yet", fieldName)
		case "Tablet":
			config.Tablet = Bool(intToBool(int(v.(float64))))
		case "TDF":
			config.TDF = Bool(intToBool(int(v.(float64))))
		case "Template":
			config.Template = Bool(intToBool(int(v.(float64))))
		case "Unique":
			config.Unique = Bool(intToBool(int(v.(float64))))
		case "USBDevices":
			log.Printf("[DEBUG] Field %s is not supported yet", fieldName)
		case "VGAType":
			v, err := VGATypeFromString(v.(string))
			if err == nil {
				config.VGAType = &v
			}
		case "VirtIODevices":
			number, _ := strconv.Atoi(matchResults[1])
			config.AddVirtIODevice(number, NewVirtIODeviceFromString(v.(string)))
		default:
			s := reflect.Indirect(reflect.ValueOf(config))
			field := s.FieldByName(fieldName)
			log.Printf("[DEBUG] Field %s: %v\n", fieldName, v)
			switch v.(type) {
			case string:
				val := v.(string)
				field.Set(reflect.ValueOf(&val))
			case float64:
				val := int(v.(float64))
				field.Set(reflect.ValueOf(&val))
			default:
				log.Printf("[DEBUG] Field %s: %v\n", fieldName, v)
			}
		}
	}

	return config
}

func findFieldName(parameter string) (string, []string, error) {
	for parameterRegexp, fieldName := range parameterMapping {
		matchResults := parameterRegexp.FindStringSubmatch(parameter)
		//matches := parameterRegexp.MatchString(parameter)
		if len(matchResults) > 0 {
			return fieldName, matchResults, nil
		}
	}
	return "", nil, errors.New("Can't find fieldName for parameter " + parameter)
}

func (c *VMConfig) AddIDEDevice(number int, value *IDEDevice) {
	if c.IDEDevices == nil {
		c.IDEDevices = make(map[int]*IDEDevice)
	}
	c.IDEDevices[number] = value
}

func (c *VMConfig) AddNetworkDevice(number int, value *NetworkDevice) {
	if c.NetworkDevices == nil {
		c.NetworkDevices = make(map[int]*NetworkDevice)
	}
	c.NetworkDevices[number] = value
}

func (c *VMConfig) AddNUMATopology(number int, value QMOption) {
	if c.NUMATopologies == nil {
		c.NUMATopologies = make(map[int]QMOption)
	}
	c.NUMATopologies[number] = value
}

func (c *VMConfig) AddParallelDevice(number int, value QMOption) {
	if c.ParallelDevices == nil {
		c.ParallelDevices = make(map[int]QMOption)
	}
	c.ParallelDevices[number] = value
}

func (c *VMConfig) AddSATADevice(number int, value QMOption) {
	if c.SATADevices == nil {
		c.SATADevices = make(map[int]QMOption)
	}
	c.SATADevices[number] = value
}

func (c *VMConfig) AddSCSIDevice(number int, value QMOption) {
	if c.SCSIDevices == nil {
		c.SCSIDevices = make(map[int]QMOption)
	}
	c.SCSIDevices[number] = value
}

func (c *VMConfig) AddSerialDevice(number int, value QMOption) {
	if c.SerialDevices == nil {
		c.SerialDevices = make(map[int]QMOption)
	}
	c.SerialDevices[number] = value
}

func (c *VMConfig) AddUSBDevice(number int, value QMOption) {
	if c.USBDevices == nil {
		c.USBDevices = make(map[int]QMOption)
	}
	c.USBDevices[number] = value
}

func (c *VMConfig) AddVirtIODevice(number int, value *VirtIODevice) {
	if c.VirtIODevices == nil {
		c.VirtIODevices = make(map[int]*VirtIODevice)
	}
	c.VirtIODevices[number] = value
}

func (c *VMConfig) GetOptionsMap() (map[string]string, error) {
	configMap := make(map[string]string)

	if c.ACPI != nil {
		configMap[parameterACPI] = boolToString(BoolValue(c.ACPI))
	}
	if c.QemuAgent != nil {
		configMap[parameterQemuAgent] = boolToString(BoolValue(c.QemuAgent))
	}
	if c.Archive != nil {
		configMap[parameterArchive] = StringValue(c.Archive)
	}
	if c.Args != nil {
		configMap[parameterArgs] = StringValue(c.Args)
	}
	if c.AutoStart != nil {
		configMap[parameterAutoStart] = boolToString(BoolValue(c.AutoStart))
	}
	if c.Balloon != nil {
		value := IntValue(c.Balloon)
		if value < 0 {
			return nil, NewArgError(parameterBalloon, "it can't be < 0")
		}
		configMap[parameterBalloon] = strconv.Itoa(value)
	}
	if c.Bios != nil {
		configMap[parameterBios] = c.Bios.String()
	}
	if c.BootOrder != nil {
		bootOrder := ""
		if len(c.BootOrder) > 4 {
			return nil, NewArgError(parameterBootOrder, "there are too many boot devices specified")
		}
		for _, bootDevice := range c.BootOrder {
			bootOrder += bootDevice.String()
		}
		configMap[parameterBootOrder] = bootOrder
	}
	if c.BootDisk != nil {
		configMap[parameterBootDisk] = StringValue(c.BootDisk)
	}
	if c.CDROM != nil {
		configMap[parameterCDROM] = StringValue(c.CDROM)
	}
	if c.Cores != nil {
		value := IntValue(c.Cores)
		if value < 1 {
			return nil, NewArgError(parameterCores, "it must be > 0")
		}
		configMap[parameterCores] = strconv.Itoa(value)
	}
	if c.CPU != nil {
		configMap[parameterCPU] = c.CPU.String()
	}
	if c.CPULimit != nil {
		value := IntValue(c.CPULimit)
		if value < 0 || value > 128 {
			return nil, NewArgError(parameterCPULimit, "it must be 0 to 128")
		}
		configMap[parameterCPULimit] = strconv.Itoa(value)
	}
	if c.CPUUnits != nil {
		value := IntValue(c.CPUUnits)
		if value < 0 || value > 500000 {
			return nil, NewArgError(parameterCPUUnits, "it must be 0 to 500000")
		}
		configMap[parameterCPUUnits] = strconv.Itoa(value)
	}
	if c.Description != nil {
		configMap[parameterDescription] = StringValue(c.Description)
	}
	if c.Force != nil {
		configMap[parameterForce] = boolToString(BoolValue(c.Force))
	}
	if c.Freeze != nil {
		configMap[parameterFreeze] = boolToString(BoolValue(c.Freeze))
	}
	if c.HostPCI != nil {
		configMap[parameterHostPCI] = StringValue(c.HostPCI)
	}
	if c.HotPlug != nil {
		configMap[parameterHotPlug] = StringValue(c.HotPlug)
	}
	if c.HugePages != nil {
		configMap[parameterHugePages] = c.HugePages.String()
	}
	if c.IDEDevices != nil {
		if len(c.IDEDevices) > 4 {
			return nil, NewArgError(fmt.Sprintf("%s[n]", parameterIDEDevices), "there are too many IDE devices specified. Max. 4")
		}
		for number, v := range c.IDEDevices {
			if number < 0 || number > 3 {
				return nil, NewArgError(fmt.Sprintf("%s[n]", parameterIDEDevices), "it must be 0 to 3")
			} else {
				key := fmt.Sprintf("%s%d", parameterIDEDevices, number)
				configMap[key] = v.GetQMOptionValue()
			}
		}
	}
	if c.KeyboardLayout != nil {
		configMap[parameterKeyboardLayout] = c.KeyboardLayout.String()
	}
	if c.KVMHardwareVirtualization != nil {
		configMap[parameterKVMHardwareVirtualization] = boolToString(BoolValue(c.KVMHardwareVirtualization))
	}
	if c.LocalTime != nil {
		configMap[parameterLocalTime] = boolToString(BoolValue(c.LocalTime))
	}
	if c.Lock != nil {
		configMap[parameterLock] = c.Lock.String()
	}
	if c.MachineType != nil {
		configMap[parameterMachineType] = StringValue(c.MachineType)
	}
	if c.Memory != nil {
		value := IntValue(c.Memory)
		if value < 16 {
			return nil, NewArgError(parameterMemory, "it must be >= 16")
		}
		configMap[parameterMemory] = strconv.Itoa(value)
	}
	if c.MigrateDowntime != nil {
		value := IntValue(c.MigrateDowntime)
		if value < 1 {
			return nil, NewArgError(parameterMigrateDowntime, "it must be >= 0")
		}
		configMap[parameterMigrateDowntime] = strconv.Itoa(value)
	}
	if c.MigrateSpeed != nil {
		value := IntValue(c.MigrateSpeed)
		if value < 1 {
			return nil, NewArgError(parameterMigrateSpeed, "it must be >= 0")
		}
		configMap[parameterMigrateSpeed] = strconv.Itoa(value)
	}
	if c.Name != nil {
		configMap[parameterName] = StringValue(c.Name)
	}
	if c.NetworkDevices != nil {
		for number, v := range c.NetworkDevices {
			if number < 0 {
				return nil, NewArgError(fmt.Sprintf("%s[n]", parameterNetworkDevices), "it must be > 0")
			} else {
				key := fmt.Sprintf("%s%d", parameterNetworkDevices, number)
				configMap[key] = v.GetQMOptionValue()
			}
		}
	}
	if c.NUMA != nil {
		configMap[parameterNUMA] = boolToString(BoolValue(c.NUMA))
	}
	if c.NUMATopologies != nil {
		for number, v := range c.NUMATopologies {
			if number < 0 {
				return nil, NewArgError(fmt.Sprintf("%s[n]", parameterNUMATopologies), "it must be > 0")
			} else {
				key := fmt.Sprintf("%s%d", parameterNUMATopologies, number)
				configMap[key] = v.GetQMOptionValue()
			}
		}
	}
	if c.StartAtBoot != nil {
		configMap[parameterStartAtBoot] = boolToString(BoolValue(c.StartAtBoot))
	}
	if c.OSType != nil {
		configMap[parameterOSType] = c.OSType.String()
	}
	if c.ParallelDevices != nil {
		if len(c.ParallelDevices) > 3 {
			return nil, NewArgError(fmt.Sprintf("%s[n]", parameterParallelDevices), "there are too many parallel devices specified. Max. 3")
		}
		for number, v := range c.ParallelDevices {
			if number < 0 || number > 2 {
				return nil, NewArgError(fmt.Sprintf("%s[n]", parameterParallelDevices), "it must be 0 to 2")
			} else {
				key := fmt.Sprintf("%s%d", parameterParallelDevices, number)
				configMap[key] = v.GetQMOptionValue()
			}
		}
	}
	if c.Pool != nil {
		configMap[parameterPool] = StringValue(c.Pool)
	}
	if c.Protection != nil {
		configMap[parameterProtection] = boolToString(BoolValue(c.Protection))
	}
	if c.Reboot != nil {
		configMap[parameterReboot] = boolToString(BoolValue(c.Reboot))
	}
	if c.SATADevices != nil {
		if len(c.SATADevices) > 6 {
			return nil, NewArgError(fmt.Sprintf("%s[n]", parameterSATADevices), "there are too many SATA devices specified. Max. 6")
		}
		for number, v := range c.SATADevices {
			if number < 0 || number > 5 {
				return nil, NewArgError(fmt.Sprintf("%s[n]", parameterSATADevices), "it must be 0 to 5")
			} else {
				key := fmt.Sprintf("%s%d", parameterSATADevices, number)
				configMap[key] = v.GetQMOptionValue()
			}
		}
	}
	if c.SCSIDevices != nil {
		if len(c.SCSIDevices) > 14 {
			return nil, NewArgError(fmt.Sprintf("%s[n]", parameterSCSIDevices), "there are too many SCSI devices specified. Max. 14")
		}
		for number, v := range c.SCSIDevices {
			if number < 0 || number > 13 {
				return nil, NewArgError(fmt.Sprintf("%s[n]", parameterSCSIDevices), "it must be 0 to 13")
			} else {
				key := fmt.Sprintf("%s%d", parameterSCSIDevices, number)
				configMap[key] = v.GetQMOptionValue()
			}
		}
	}
	if c.SCSIControllerType != nil {
		configMap[parameterSCSIControllerType] = c.SCSIControllerType.String()
	}
	if c.SerialDevices != nil {
		if len(c.SerialDevices) > 4 {
			return nil, NewArgError(fmt.Sprintf("%s[n]", parameterSerialDevices), "there are too many serial devices specified. Max. 4")
		}
		for number, v := range c.SerialDevices {
			if number < 0 || number > 3 {
				return nil, NewArgError(fmt.Sprintf("%s[n]", parameterSerialDevices), "it must be 0 to 3")
			} else {
				key := fmt.Sprintf("%s%d", parameterSerialDevices, number)
				configMap[key] = v.GetQMOptionValue()
			}
		}
	}
	if c.MemoryShares != nil {
		value := IntValue(c.MemoryShares)
		if value < 0 || value > 50000 {
			return nil, NewArgError(parameterMemoryShares, "it must be 0 to 50000")
		}
		configMap[parameterMemoryShares] = strconv.Itoa(value)
	}
	if c.SMBIOS1 != nil {
		configMap[parameterSMBIOS1] = StringValue(c.SMBIOS1)
	}
	if c.SMP != nil {
		value := IntValue(c.SMP)
		if value < 1 {
			return nil, NewArgError(parameterSMP, "it must be > 0")
		}
		configMap[parameterSMP] = strconv.Itoa(value)
	}
	if c.Sockets != nil {
		value := IntValue(c.Sockets)
		if value < 1 {
			return nil, NewArgError(parameterSockets, "it must be > 0")
		}
		configMap[parameterSockets] = strconv.Itoa(value)
	}
	if c.StartDate != nil {
		configMap[parameterStartDate] = StringValue(c.StartDate)
	}
	if c.Startup != nil {
		configMap[parameterStartup] = StringValue(c.Startup)
	}
	if c.Storage != nil {
		configMap[parameterStorage] = StringValue(c.Storage)
	}
	if c.Tablet != nil {
		configMap[parameterTablet] = boolToString(BoolValue(c.Tablet))
	}
	if c.TDF != nil {
		configMap[parameterTDF] = boolToString(BoolValue(c.TDF))
	}
	if c.Template != nil {
		configMap[parameterTemplate] = boolToString(BoolValue(c.Template))
	}
	if c.Unique != nil {
		configMap[parameterUnique] = boolToString(BoolValue(c.Unique))
	}
	if c.USBDevices != nil {
		if len(c.USBDevices) > 5 {
			return nil, NewArgError(fmt.Sprintf("%s[n]", parameterUSBDevices), "there are too many USB devices specified. Max. 5")
		}
		for number, v := range c.USBDevices {
			if number < 0 || number > 4 {
				return nil, NewArgError(fmt.Sprintf("%s[n]", parameterUSBDevices), "it must be 0 to 4")
			} else {
				key := fmt.Sprintf("%s%d", parameterUSBDevices, number)
				configMap[key] = v.GetQMOptionValue()
			}
		}
	}
	if c.VCPUs != nil {
		configMap[parameterVCPUs] = strconv.Itoa(IntValue(c.VCPUs))
	}
	if c.VGAType != nil {
		configMap[parameterVGA] = c.VGAType.String()
	}
	if c.VirtIODevices != nil {
		if len(c.VirtIODevices) > 16 {
			return nil, NewArgError(fmt.Sprintf("%s[n]", parameterVirtIODevices), "there are too many VirtIO devices specified. Max. 16")
		}
		for number, v := range c.VirtIODevices {
			if number < 0 || number > 15 {
				return nil, NewArgError(fmt.Sprintf("%s[n]", parameterVirtIODevices), "it must be 0 to 15")
			} else {
				key := fmt.Sprintf("%s%d", parameterVirtIODevices, number)
				configMap[key] = v.GetQMOptionValue()
			}
		}
	}
	if c.VMID != nil {
		value := IntValue(c.VMID)
		if value < 100 {
			return nil, NewArgError(parameterVMID, "it should be >= 100. IDs < 100 are reserved for internal purposes.")
		}
		configMap[parameterVMID] = strconv.Itoa(value)
	}
	if c.Watchdog != nil {
		configMap[parameterWatchdog] = StringValue(c.Watchdog)
	}

	return configMap, nil
}
