package goproxmox

type NodesService interface {
	GetNodes() ([]Node, error)
}

type NodesServiceOp struct {
	client *Client
}

var _ NodesService = &NodesServiceOp{}

type Node struct {
	Id      string  `json:"id"`
	Node    string  `json:"node"`
	Type    string  `json:"type"`
	Cpu     float64 `json:"cpu"`
	Mem     uint64  `json:"mem"`
	Disk    uint64  `json:"disk"`
	Maxcpu  int     `json:"maxcpu"`
	Maxmem  uint64  `json:"maxmem"`
	Maxdisk uint64  `json:"maxdisk"`
	Uptime  uint    `json:"uptime"`
	Level   string  `json:"level"`
}

type nodesRoot struct {
	Nodes []Node `json:"data"`
}

func (s *NodesServiceOp) GetNodes() ([]Node, error) {
	path := "nodes"

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	root := new(nodesRoot)
	if _, err = s.client.Do(req, root); err != nil {
		return nil, err
	}

	return root.Nodes, err
}
