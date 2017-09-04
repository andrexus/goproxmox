package goproxmox

import (
	"fmt"
	"errors"
)

type QemuService interface {
	GetVMList(node string) ([]VM, error)
	GetVMCurrentStatus(node string, vmID int) (*VMStatus, error)
	StartVM(node string, vmID int) error
	StopVM(node string, vmID int) error
	ShutdownVM(node string, vmID int) error
	ResetVM(node string, vmID int) error
	SuspendVM(node string, vmID int) error
	ResumeVM(node string, vmID int) error
	GetVMConfig(node string, vmID int) (VMConfig, error)
	CreateVM(node string, vmID int, config VMConfig) error
	UpdateVM(node string, vmID int, config VMConfig, async bool) error
	DeleteVM(node string, vmID int) error
}

type QemuServiceOp struct {
	client *Client
}

var _ QemuService = &QemuServiceOp{}

type responseRoot struct {
	Data map[string]interface{} `json:"data"`
}

type vmsRoot struct {
	VMs []VM `json:"data"`
}

type vmStatusRoot struct {
	VMStatus VMStatus `json:"data"`
}

type VM struct {
	VMID      int     `json:"vmid"`
	Name      string  `json:"name"`
	Status    string  `json:"status"`
	PID       string  `json:"pid"`
	Template  string  `json:"template"`
	CPU       float64 `json:"cpu"`
	CPUs      int     `json:"cpus"`
	Memory    int     `json:"mem"`
	MaxMemory int     `json:"maxmem"`
	Disk      int     `json:"disk"`
	DiskRead  int     `json:"diskread"`
	DiskWrite int     `json:"diskwrite"`
	MaxDisk   int     `json:"maxdisk"`
	NetIn     int     `json:"netin"`
	NetOut    int     `json:"netout"`
	Uptime    int     `json:"uptime"`
}

type VMStatus struct {
	Name      string      `json:"name"`
	Status    string      `json:"status"`
	QMPstatus string      `json:"qmpstatus"`
	PID       int         `json:"pid"`
	Template  string      `json:"template"`
	CPU       int         `json:"cpu"`
	CPUs      int         `json:"cpus"`
	Memory    int         `json:"mem"`
	MaxMemory int         `json:"maxmem"`
	Disk      int         `json:"disk"`
	DiskRead  int         `json:"diskread"`
	DiskWrite int         `json:"diskwrite"`
	MaxDisk   int         `json:"maxdisk"`
	NetIn     int         `json:"netin"`
	NetOut    int         `json:"netout"`
	Uptime    int         `json:"uptime"`
	HA        interface{} `json:"ha"`
}

// Virtual machine index (per node).
func (s *QemuServiceOp) GetVMList(node string) ([]VM, error) {
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

// Get config for the virtual machine
func (s *QemuServiceOp) GetVMConfig(node string, vmID int) (VMConfig, error) {
	path := fmt.Sprintf("nodes/%s/qemu/%d/config", node, vmID)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	root := new(responseRoot)
	if _, err = s.client.Do(req, root); err != nil {
		return nil, err
	}
	config := NewVMConfigFromMap(root.Data)
	return config, nil
}

// Create virtual machine.
func (s *QemuServiceOp) CreateVM(node string, vmID int, config VMConfig) error {
	if config == nil {
		config = NewVMConfig()
	}
	if err := config.SetVMID(vmID); err != nil {
		return err
	}

	if vms, err := s.GetVMList(node); err != nil {
		return err
	} else {
		for _, vm := range vms {
			if vmID == vm.VMID {
				return errors.New(fmt.Sprintf("VM with id %d already exists", vmID))
			}
		}
	}

	path := fmt.Sprintf("nodes/%s/qemu", node)
	req, err := s.client.NewRequest("POST", path, config.GetOptionsMap())
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Update virtual machine.
func (s *QemuServiceOp) UpdateVM(node string, vmID int, config VMConfig, async bool) error {
	path := fmt.Sprintf("nodes/%s/qemu/%d/config", node, vmID)
	method := "PUT" // synchronous API
	if async == true {
		method = "POST" // asynchronous API
	}
	req, err := s.client.NewRequest(method, path, config.GetOptionsMap())
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
