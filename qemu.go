package goproxmox

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
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
	GetVMConfig(node string, vmID int) (*VMConfig, error)
	CreateVM(node string, vmID int, config *VMConfig) error
	UpdateVM(node string, vmID int, config *VMConfig, async bool) error
	DeleteVM(node string, vmID int) error
	CreateVMTemplate(node string, vmID int, disk string) error
	CloneVM(node string, vmID int, newID int, config *VMCloneConfig) error
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
	VMID      int         `json:"vmid"`
	Name      string      `json:"name"`
	Status    string      `json:"status"`
	PID       string      `json:"pid"`
	Template  interface{} `json:"template"`
	CPU       float64     `json:"cpu"`
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
}

type VMStatus struct {
	Name      string      `json:"name"`
	Status    string      `json:"status"`
	QMPstatus string      `json:"qmpstatus"`
	PID       string      `json:"pid"`
	Template  interface{} `json:"template"`
	CPU       float64     `json:"cpu"`
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

type VMCloneConfig struct {
	Name          *string // Set a name for the new VM
	Description   *string // Description for the new VM
	Full          *bool   // Create a full copy of all disk. This is always done when you clone a normal VM. For VM templates, we try to create a linked clone by default
	Pool          *string // Add the new VM to the specified pool
	SnapshotName  *string // The name of the snapshot
	Storage       *string // Target storage for full clone
	StorageFormat *string // Target format for file storage
	TargetNode    *string // Target node. Only allowed if the original VM is on shared storage
}

func (c *VMCloneConfig) getRequestBodyParameters() map[string]string {
	bodyParams := map[string]string{}
	if c.Name != nil {
		bodyParams["name"] = StringValue(c.Name)
	}
	if c.Description != nil {
		bodyParams["description"] = StringValue(c.Description)
	}
	if c.Full != nil {
		bodyParams["full"] = boolToString(BoolValue(c.Full))
	}
	if c.Pool != nil {
		bodyParams["pool"] = StringValue(c.Pool)
	}
	if c.SnapshotName != nil {
		bodyParams["snapname"] = StringValue(c.SnapshotName)
	}
	if c.Storage != nil {
		bodyParams["storage"] = StringValue(c.Storage)
	}
	if c.StorageFormat != nil {
		bodyParams["format"] = StringValue(c.StorageFormat)
	}
	if c.TargetNode != nil {
		bodyParams["target"] = StringValue(c.TargetNode)
	}
	return bodyParams
}

// Virtual machine index (per node).
func (s *QemuServiceOp) GetVMList(node string) ([]VM, error) {
	path := fmt.Sprintf("nodes/%s/qemu", node)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
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

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
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

	req, err := s.client.NewRequest(http.MethodPost, path, nil)
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

	req, err := s.client.NewRequest(http.MethodPost, path, nil)
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

	req, err := s.client.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Reset virtual machine.
func (s *QemuServiceOp) ResetVM(node string, vmID int) error {
	path := fmt.Sprintf("nodes/%s/qemu/%d/status/reset", node, vmID)

	req, err := s.client.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Suspend virtual machine.
func (s *QemuServiceOp) SuspendVM(node string, vmID int) error {
	path := fmt.Sprintf("nodes/%s/qemu/%d/status/suspend", node, vmID)

	req, err := s.client.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Resume virtual machine.
func (s *QemuServiceOp) ResumeVM(node string, vmID int) error {
	path := fmt.Sprintf("nodes/%s/qemu/%d/status/resume", node, vmID)

	req, err := s.client.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Get config for the virtual machine
func (s *QemuServiceOp) GetVMConfig(node string, vmID int) (*VMConfig, error) {
	path := fmt.Sprintf("nodes/%s/qemu/%d/config", node, vmID)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
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
func (s *QemuServiceOp) CreateVM(node string, vmID int, config *VMConfig) error {
	if config == nil {
		config = &VMConfig{}
	}
	config.VMID = Int(vmID)

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
	optionsMap, err := config.GetOptionsMap()
	if err != nil {
		return err
	}
	req, err := s.client.NewRequest(http.MethodPost, path, optionsMap)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Update virtual machine.
func (s *QemuServiceOp) UpdateVM(node string, vmID int, config *VMConfig, async bool) error {
	path := fmt.Sprintf("nodes/%s/qemu/%d/config", node, vmID)
	method := http.MethodPut // synchronous API
	if async == true {
		method = http.MethodPost // asynchronous API
	}
	optionsMap, err := config.GetOptionsMap()
	if err != nil {
		return err
	}
	req, err := s.client.NewRequest(method, path, optionsMap)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Create virtual machine.
func (s *QemuServiceOp) DeleteVM(node string, vmID int) error {
	path := fmt.Sprintf("nodes/%s/qemu/%d", node, vmID)

	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Create a template from VM.
func (s *QemuServiceOp) CreateVMTemplate(node string, vmID int, disk string) error {
	path := fmt.Sprintf("nodes/%s/qemu/%d/template", node, vmID)

	body := make(map[string]string)
	if disk != "" {
		body["disk"] = disk
	}

	req, err := s.client.NewRequest(http.MethodPost, path, body)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Clone VM.
func (s *QemuServiceOp) CloneVM(node string, vmID int, newID int, config *VMCloneConfig) error {
	path := fmt.Sprintf("nodes/%s/qemu/%d/clone", node, vmID)
	body := make(map[string]string)
	body["newid"] = strconv.Itoa(newID)
	if config != nil {
		for k, v := range config.getRequestBodyParameters() {
			body[k] = v
		}
	}
	req, err := s.client.NewRequest(http.MethodPost, path, body)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}
