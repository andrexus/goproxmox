package goproxmox

import "fmt"

type Bios int

const (
	BIOS_SeaBIOS Bios = 1 + iota
	BIOS_OVMF
)

var biosValues = [...]string{
	"seabios",
	"ovmf",
}

// String returns the name of the Bios.
func (m Bios) String() string { return biosValues[m-1] }

func BiosFromString(s string) (Bios, error) {
	for i, v := range biosValues {
		if s == v {
			return Bios(i + 1), nil
		}
	}
	return 0, fmt.Errorf("%s does not belong to Bios values", s)
}

type BootDevice int

const (
	BOOT_Floppy BootDevice = 1 + iota
	BOOT_HardDisk
	BOOT_CDROM
	BOOT_Network
)

var bootDeviceValues = [...]string{
	"a",
	"c",
	"d",
	"n",
}

// String returns the name of the BootDevice.
func (m BootDevice) String() string { return bootDeviceValues[m-1] }

func BootDeviceFromString(s string) (BootDevice, error) {
	for i, v := range bootDeviceValues {
		if s == v {
			return BootDevice(i + 1), nil
		}
	}
	return 0, fmt.Errorf("%s does not belong to BootDevice values", s)
}

type CPUType int

const (
	CPU_486 CPUType = 1 + iota
	CPU_Broadwell
	CPU_Broadwell_noTSX
	CPU_Conroe
	CPU_Haswell
	CPU_Haswell_noTSX
	CPU_IvyBridge
	CPU_Nehalem
	CPU_Opteron_G1
	CPU_Opteron_G2
	CPU_Opteron_G3
	CPU_Opteron_G4
	CPU_Opteron_G5
	CPU_Penryn
	CPU_SandyBridge
	CPU_Skylake_Client
	CPU_Westmere
	CPU_Athlon
	CPU_Core2duo
	CPU_CoreDuo
	CPU_HOST
	CPU_KVM32
	CPU_KVM64
	CPU_Pentium
	CPU_Pentium2
	CPU_Pentium3
	CPU_Phenom
	CPU_Qemu32
	CPU_Qemu64
)

var cpuTypeValues = [...]string{
	"486",
	"Broadwell",
	"Broadwell-noTSX",
	"Conroe",
	"Haswell",
	"Haswell-noTSX",
	"IvyBridge",
	"Nehalem",
	"Opteron_G1",
	"Opteron_G2",
	"Opteron_G3",
	"Opteron_G4",
	"Opteron_G5",
	"Penryn",
	"SandyBridge",
	"Skylake-Client",
	"Westmere",
	"athlon",
	"core2duo",
	"coreduo",
	"host",
	"kvm32",
	"kvm64",
	"pentium",
	"pentium2",
	"pentium3",
	"phenom",
	"qemu32",
	"qemu64",
}

// String returns the name of the CPUType.
func (m CPUType) String() string { return cpuTypeValues[m-1] }

func CPUTypeFromString(s string) (CPUType, error) {
	for i, v := range bootDeviceValues {
		if s == v {
			return CPUType(i + 1), nil
		}
	}
	return 0, fmt.Errorf("%s does not belong to CPUType values", s)
}

type HugePages int

const (
	HugePages_1024 HugePages = 1 + iota
	HugePages_2
	HugePages_ANY
)

var hugePagesValues = [...]string{
	"1024",
	"2",
	"any",
}

// String returns the name of the HugePages.
func (m HugePages) String() string { return hugePagesValues[m-1] }

func HugePagesFromString(s string) (HugePages, error) {
	for i, v := range bootDeviceValues {
		if s == v {
			return HugePages(i + 1), nil
		}
	}
	return 0, fmt.Errorf("%s does not belong to HugePages values", s)
}

type KeyboardLayout int

const (
	KeyboardLayout_DA KeyboardLayout = 1 + iota
	KeyboardLayout_DE
	KeyboardLayout_DE_CH
	KeyboardLayout_EN_GB
	KeyboardLayout_EN_US
	KeyboardLayout_ES
	KeyboardLayout_FI
	KeyboardLayout_FR
	KeyboardLayout_FR_BE
	KeyboardLayout_FR_CA
	KeyboardLayout_FR_CH
	KeyboardLayout_HU
	KeyboardLayout_IS
	KeyboardLayout_IT
	KeyboardLayout_JA
	KeyboardLayout_LT
	KeyboardLayout_MK
	KeyboardLayout_NL
	KeyboardLayout_NO
	KeyboardLayout_PL
	KeyboardLayout_PT
	KeyboardLayout_PT_BR
	KeyboardLayout_SL
	KeyboardLayout_SV
	KeyboardLayout_TR
)

var keyboardLayoutValues = [...]string{
	"da",
	"de",
	"de-ch",
	"en-gb",
	"en-us",
	"es",
	"fi",
	"fr",
	"fr-be",
	"fr-ca",
	"fr-ch",
	"hu",
	"is",
	"it",
	"ja",
	"lt",
	"mk",
	"nl",
	"no",
	"pl",
	"pt",
	"pt-br",
	"sl",
	"sv",
	"tr",
}

// String returns the name of the KeyboardLayout.
func (m KeyboardLayout) String() string { return keyboardLayoutValues[m-1] }

func KeyboardLayoutFromString(s string) (KeyboardLayout, error) {
	for i, v := range bootDeviceValues {
		if s == v {
			return KeyboardLayout(i + 1), nil
		}
	}
	return 0, fmt.Errorf("%s does not belong to KeyboardLayout values", s)
}

type Lock int

const (
	Lock_Migrate Lock = 1 + iota
	Lock_Backup
	Lock_Snapshot
	Lock_Rollback
)

var lockValues = [...]string{
	"migrate",
	"backup",
	"snapshot",
	"rollback",
}

// String returns the name of the Lock.
func (m Lock) String() string { return lockValues[m-1] }

func LockFromString(s string) (Lock, error) {
	for i, v := range bootDeviceValues {
		if s == v {
			return Lock(i + 1), nil
		}
	}
	return 0, fmt.Errorf("%s does not belong to Lock values", s)
}

type OSType int

const (
	OS_Unspecified OSType = 1 + iota
	OS_WindowsXP
	OS_Windows2000
	OS_Windows2003
	OS_Windows2008
	OS_WindowsVista
	OS_Windows7
	OS_Windows8_2012
	OS_Linux24
	OS_Linux26_3X
	OS_Solaris
)

var osTypeValues = [...]string{
	"other",
	"wxp",
	"w2k",
	"w2k3",
	"w2k8",
	"wvista",
	"win7",
	"win8",
	"l24",
	"l26",
	"solaris",
}

// String returns the name of the OSType.
func (m OSType) String() string { return osTypeValues[m-1] }

func OSTypeFromString(s string) (OSType, error) {
	for i, v := range bootDeviceValues {
		if s == v {
			return OSType(i + 1), nil
		}
	}
	return 0, fmt.Errorf("%s does not belong to OSType values", s)
}

type SCSIControllerType int

const (
	SCSI_LSI SCSIControllerType = 1 + iota
	SCSI_LSI53C810
	SCSI_VirtIO_SCSI_PCI
	SCSI_VirtIO_SCSI_SINGLE
	SCSI_MEGASAS
	SCSI_PVSCSI
)

var scsiControllerTypeValues = [...]string{
	"lsi",
	"lsi53c810",
	"virtio-scsi-pci",
	"virtio-scsi-single",
	"megasas",
	"pvscsi",
}

// String returns the name of the SCSIControllerType.
func (m SCSIControllerType) String() string { return scsiControllerTypeValues[m-1] }

func SCSIControllerTypeFromString(s string) (SCSIControllerType, error) {
	for i, v := range bootDeviceValues {
		if s == v {
			return SCSIControllerType(i + 1), nil
		}
	}
	return 0, fmt.Errorf("%s does not belong to SCSIControllerType values", s)
}

type VGAType int

const (
	VGA_Cirrus VGAType = 1 + iota
	VGA_QXL
	VGA_QXL2
	VGA_QXL3
	VGA_QXL4
	VGA_Serial0
	VGA_Serial1
	VGA_Serial2
	VGA_Serial3
	VGA_std
	VGA_VMWare
)

var vgaTypeValues = [...]string{
	"cirrus",
	"qxl",
	"qxl2",
	"qxl3",
	"qxl4",
	"serial0",
	"serial1",
	"serial2",
	"serial3",
	"std",
	"vmware",
}

// String returns the name of the VGAType.
func (m VGAType) String() string { return vgaTypeValues[m-1] }

func VGATypeFromString(s string) (VGAType, error) {
	for i, v := range bootDeviceValues {
		if s == v {
			return VGAType(i + 1), nil
		}
	}
	return 0, fmt.Errorf("%s does not belong to VGAType values", s)
}

type NetworkCardModel int

const (
	NetworkCard_E1000 NetworkCardModel = 1 + iota
	NetworkCard_E1000_82540em
	NetworkCard_E1000_82544gc
	NetworkCard_E1000_82545em
	NetworkCard_I82551
	NetworkCard_I82557b
	NetworkCard_I82559er
	NetworkCard_NE2K_ISA
	NetworkCard_NE2K_PCI
	NetworkCard_PCNET
	NetworkCard_RTL8139
	NetworkCard_VIRTIO
	NetworkCard_VMXNET3
)

var networkCardModelValues = [...]string{
	"e1000",
	"e1000-82540em",
	"e1000-82544gc",
	"e1000-82545em",
	"i82551",
	"i82557b",
	"i82559er",
	"ne2k_isa",
	"ne2k_pci",
	"pcnet",
	"rtl8139",
	"virtio",
	"vmxnet3",
}

// String returns the name of the NetworkCardModel.
func (m NetworkCardModel) String() string { return networkCardModelValues[m-1] }
