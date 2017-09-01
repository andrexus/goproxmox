package goproxmox

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
