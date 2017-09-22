package goproxmox

import (
	"fmt"
	"strconv"
)

type StorageService interface {
	GetStorageList(node string) ([]Storage, error)
	GetStorageVolumes(node, storageName string) ([]StorageVolume, error)
	GetVolume(node, storageName, volumeId string) (*StorageVolume, error)
	CreateVolume(node, storageName string, vmID int, filename string, size string, format *string) error
	DeleteVolume(node, storageName, volumeId string) error
}

type StorageServiceOp struct {
	client *Client
}

var _ StorageService = &StorageServiceOp{}

type storagesRoot struct {
	Storages []Storage `json:"data"`
}

type volumeRoot struct {
	Volume StorageVolume `json:"data"`
}

type volumesRoot struct {
	Volumes []StorageVolume `json:"data"`
}

type Storage struct {
	StorageName string `json:"storage"`
	Content     string `json:"content"`
	Type        string `json:"type"`
	Active      int    `json:"active"`
	Enabled     int    `json:"enabled"`
	Shared      int    `json:"shared"`
	Used        int64  `json:"used"`
	Available   int64  `json:"avail"`
	Total       int64  `json:"total"`
}

type StorageVolume struct {
	Format   string      `json:"format"`
	VolumeId string      `json:"volid"`
	Parent   interface{} `json:"parent"`
	Size     int64       `json:"size"`
	Content  string      `json:"content"`
	Used     int64       `json:"used"`
	VMID     string      `json:"vmid"`
	Path     string      `json:"path"`
}

// Storage list (per node).
func (s *StorageServiceOp) GetStorageList(node string) ([]Storage, error) {
	path := fmt.Sprintf("nodes/%s/storage", node)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	root := new(storagesRoot)
	if _, err = s.client.Do(req, root); err != nil {
		return nil, err
	}

	return root.Storages, err
}

// Get list of volumes per node and storage
func (s *StorageServiceOp) GetStorageVolumes(node, storageName string) ([]StorageVolume, error) {
	path := fmt.Sprintf("nodes/%s/storage/%s/content", node, storageName)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	root := new(volumesRoot)
	if _, err = s.client.Do(req, root); err != nil {
		return nil, err
	}

	return root.Volumes, err
}

// Get volume attributes
func (s *StorageServiceOp) GetVolume(node, storageName, volumeId string) (*StorageVolume, error) {
	path := fmt.Sprintf("nodes/%s/storage/%s/content/%s", node, storageName, volumeId)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	root := new(volumeRoot)
	if _, err = s.client.Do(req, root); err != nil {
		return nil, err
	}

	return &root.Volume, err
}

// Create new volume.
func (s *StorageServiceOp) CreateVolume(node, storageName string, vmID int, filename string, size string, format *string) error {
	path := fmt.Sprintf("nodes/%s/storage/%s/content", node, storageName)
	optionsMap := make(map[string]string)
	optionsMap["filename"] = filename
	optionsMap["size"] = size
	optionsMap["vmid"] = strconv.Itoa(vmID)
	if format != nil {
		optionsMap["format"] = *format
	}

	req, err := s.client.NewRequest("POST", path, optionsMap)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// Delete existing volume.
func (s *StorageServiceOp) DeleteVolume(node, storageName, volumeId string) error {
	path := fmt.Sprintf("nodes/%s/storage/%s/content/%s", node, storageName, volumeId)

	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}
