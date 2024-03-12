package ns

import "syscall"

// MountProc 挂载/proc文件系统
func MountProc() error {
	procMountPoint := "/proc"
	source := "proc"
	fsType := "proc"
	procMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	data := ""
	if err := syscall.Mount(source, procMountPoint, fsType, uintptr(procMountFlags), data); err != nil {
		return err
	}
	return nil
}
