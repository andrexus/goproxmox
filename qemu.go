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

type vmsRoot struct {
	VMs []VM `json:"data"`
}

type vmStatusRoot struct {
	VMStatus VMStatus `json:"data"`
}

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

func (s *QemuServiceOp) StartVM(node string, vmID string) error {
	path := fmt.Sprintf("nodes/%s/qemu/%s/status/start", node, vmID)

	req, err := s.client.NewRequest("POST", path, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

func (s *QemuServiceOp) StopVM(node string, vmID string) error {
	path := fmt.Sprintf("nodes/%s/qemu/%s/status/stop", node, vmID)

	req, err := s.client.NewRequest("POST", path, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

func (s *QemuServiceOp) ShutdownVM(node string, vmID string) error {
	path := fmt.Sprintf("nodes/%s/qemu/%s/status/shutdown", node, vmID)

	req, err := s.client.NewRequest("POST", path, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

func (s *QemuServiceOp) ResetVM(node string, vmID string) error {
	path := fmt.Sprintf("nodes/%s/qemu/%s/status/reset", node, vmID)

	req, err := s.client.NewRequest("POST", path, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

func (s *QemuServiceOp) SuspendVM(node string, vmID string) error {
	path := fmt.Sprintf("nodes/%s/qemu/%s/status/suspend", node, vmID)

	req, err := s.client.NewRequest("POST", path, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

func (s *QemuServiceOp) ResumeVM(node string, vmID string) error {
	path := fmt.Sprintf("nodes/%s/qemu/%s/status/resume", node, vmID)

	req, err := s.client.NewRequest("POST", path, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}
