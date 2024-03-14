package ns

import (
	"fmt"
	"mydocker/config"
	"os"
	"syscall"
)

type OverlayFs struct {
	LowerDir string
	UpperDir string
	WorkDir  string
	MergeDir string
}

// GetOverlayFs 获取overlay文件系统
func GetOverlayFs(imageName, containerName string) *OverlayFs {
	return &OverlayFs{
		LowerDir: config.Cfg.Images.ImagePath + "/" + imageName,
		UpperDir: config.Cfg.RootFs.UpperLayerPath + "/" + containerName,
		WorkDir:  config.Cfg.RootFs.WorkLayerPath + "/" + containerName,
		MergeDir: config.Cfg.RootFs.MntPath + "/" + containerName,
	}
}

// InitMntNameSpace 初始化mnt namespace
func InitMntNameSpace(imageName, containerName string) error {
	fs := GetOverlayFs(imageName, containerName)
	if err := fs.create(); err != nil {
		return err
	}
	if err := fs.chroot(); err != nil {
		return err
	}
	return nil
}

// DeleteMntNameSpace 删除mnt namespace
func DeleteMntNameSpace(containerName string) error {
	fs := GetOverlayFs("", containerName)
	if err := fs.destroy(); err != nil {
		return err
	}
	return nil
}

// create 创建overlay文件系统
func (fs *OverlayFs) create() error {
	if err := os.MkdirAll(fs.UpperDir, 0777); err != nil {
		return err
	}
	if err := os.MkdirAll(fs.WorkDir, 0777); err != nil {
		return err
	}
	if err := os.MkdirAll(fs.MergeDir, 0777); err != nil {
		return err
	}
	mntInfo := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", fs.LowerDir, fs.UpperDir, fs.WorkDir)
	fmt.Println(mntInfo)
	if err := syscall.Mount("overlay", fs.MergeDir, "overlay", 0, mntInfo); err != nil {
		return fmt.Errorf("mount overlay fail err=%s", err)
	}
	if err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("reclare rootfs private fail err=%s", err)
	}
	if err := syscall.Mount(fs.MergeDir, fs.MergeDir, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("mount rootfs in new mnt space fail err=%s", err)
	}
	return nil
}

// chroot 切换根目录
func (fs *OverlayFs) chroot() error {
	if err := syscall.Chroot(fs.MergeDir); err != nil {
		return fmt.Errorf("chroot fail err=%s", err)
	}
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir fail err=%s", err)
	}
	return nil
}

// destroy 销毁overlay文件系统
func (fs *OverlayFs) destroy() error {
	if err := syscall.Unmount(fs.MergeDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount overlay fail err=%s", err)
	}
	if err := os.RemoveAll(fs.UpperDir); err != nil {
		return err
	}
	if err := os.RemoveAll(fs.WorkDir); err != nil {
		return err
	}
	if err := os.RemoveAll(fs.MergeDir); err != nil {
		return err
	}
	return nil
}
