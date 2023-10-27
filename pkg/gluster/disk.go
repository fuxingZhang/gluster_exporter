package gluster

import "syscall"

// DiskStatus represents Disk usage
type DiskStatus struct {
	All        float64 `json:"all"`
	Used       float64 `json:"used"`
	Free       float64 `json:"free"`
	InodesAll  float64 `json:"inodesall"`
	InodesFree float64 `json:"inodesfree"`
	InodesUsed float64 `json:"inodesused"`
}

func DiskUsage(path string) (disk DiskStatus, err error) {
	fs := syscall.Statfs_t{}
	err = syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = float64(fs.Blocks * uint64(fs.Bsize))
	disk.Free = float64(fs.Bfree * uint64(fs.Bsize))
	disk.Used = disk.All - disk.Free
	disk.InodesAll = float64(fs.Files)
	disk.InodesFree = float64(fs.Ffree)
	disk.InodesUsed = disk.InodesAll - disk.InodesFree
	return
}
