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

func createOverlayFs(fs *OverlayFs) error {
	if err := os.MkdirAll(fs.UpperDir, 0777); err != nil {
		return err
	}
	if err := os.MkdirAll(fs.WorkDir, 0777); err != nil {
		return err
	}
	if err := os.MkdirAll(fs.MergeDir, 0777); err != nil {
		return err
	}
	data := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", fs.LowerDir, fs.UpperDir, fs.WorkDir)
	if err := syscall.Mount("overlay", fs.MergeDir, "overlay", 0, data); err != nil {
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

func Chroot(path string) error {
	if err := syscall.Chroot(path); err != nil {
		return fmt.Errorf("chroot fail err=%s", err)
	}
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir fail err=%s", err)
	}
	return nil
}

func InitMntNameSpace(imageName, containerName string) error {
	overlayFs := &OverlayFs{
		LowerDir: config.Cfg.Images.ImagePath + "/" + imageName,
		UpperDir: config.Cfg.RootFs.UpperLayerPath + "/" + containerName,
		WorkDir:  config.Cfg.RootFs.WorkLayerPath + "/" + containerName,
		MergeDir: config.Cfg.RootFs.MntPath + "/" + containerName,
	}
	if err := createOverlayFs(overlayFs); err != nil {
		return err
	}
	if err := Chroot(overlayFs.MergeDir); err != nil {
		return err
	}
	return nil
}

func delOverlayFs(fs *OverlayFs) error {
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

func DeleteMntNameSpace(containerName string) error {
	overlayFs := &OverlayFs{
		UpperDir: config.Cfg.RootFs.UpperLayerPath + "/" + containerName,
		WorkDir:  config.Cfg.RootFs.WorkLayerPath + "/" + containerName,
		MergeDir: config.Cfg.RootFs.MntPath + "/" + containerName,
	}
	if err := delOverlayFs(overlayFs); err != nil {
		return err
	}
	return nil
}
