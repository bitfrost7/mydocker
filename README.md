# 手写Docker系列

Docker是一个开源的应用容器引擎，基于Go语言并遵从Apache2.0协议开源。Docker可以让开发者打包他们的应用以及依赖包到一个可移植的容器中，然后发布到任何流行的Linux机器上，也可以实现虚拟化。 容器是完全使用沙箱机制，相互之间不会有任何接口。

### 1.1 LXC介绍

在Docker之前，linux就已经出现了LXC(Linux Containers)虚拟化技术，它可以将应用打包成成一个软件容器（Container),内含应用软件本身的代码，以及所需要的操作系统核心和库。通过统一的命名空间(Namespace)和共享API来分配不同软件容器的可用硬件资源，创造出应用程序的独立沙箱运行环境，使得Linux用户可以容易的创建和管理系统或应用容器。
作为一个开源容器平台，Linux 容器项目（LXC）提供了一组工具、模板、库和语言绑定。LXC 采用简单的命令行界面，可改善容器启动时的用户体验。LXC 提供了一个操作系统级的虚拟化环境，可在许多基于 Linux 的系统上安装。在 Linux 发行版中，可能会通过其软件包存储库来提供 LXC。
在Linux内核中，container技术的核心还是cgroup+chroot+namespace技术：其提供了cgroups功能来达成资源的隔离；它同时也提供了命名空间隔离的功能，使应用程序看到的操作系统环境被区隔成独立区间，包括进程树、网络、用户id、以及挂载的文件系统；但是cgroups并不一定需要启动任何虚拟机。
LXC利用cgroups与命名空间的功能，为应用软件提供一个相对独立的操作系统环境。LXC不需要Hypervisor这个软件层，软件容器（Container）本身极为轻量化，提升了创建虚拟机的速度。
可以说Docker只是LXC的一个封装，它提供了更简单易用的接口，使得LXC更容易使用。

#### 1.1.1 chroot

chroot技术可以让你在不影响系统其他部分的情况下，将进程的根目录改变为指定的目录，这样进程就无法访问到系统的其他目录了，这样就可以达到隔离的目的。

#### 1.1.2 Namespace

Namespace是Linux内核提供的一种资源隔离机制，它可以将一系列的全局系统资源封装到一个抽象的命名空间中，使得其他的命名空间看不到这些资源，从而实现资源的隔离。
Linux内核提供了7种Namespace：

- Mount Namespace：隔离文件系统挂载点
- UTS Namespace：隔离主机名和域名
- IPC Namespace：隔离System V IPC和POSIX消息队列
- PID Namespace：隔离进程ID
- Network Namespace：隔离网络设备、网络栈、端口等网络资源
- User Namespace：隔离用户和用户组
- Cgroup Namespace：隔离cgroup根目录

#### 1.1.3 Cgroup

Cgroup(Control Group)是Linux内核提供的一种资源限制机制，它可以限制进程组能够使用的资源上限，包括CPU、内存、磁盘、网络带宽等等。

## 1. Docker 相关技术简介

CGroup 为每种可以控制的资源定义了一个子系统。典型的子系统介绍如下:

- cpuset：这个子系统可以分配一个或者多个CPU及内存节点给进程组使用。
- cpu：这个子系统使用调度程序来分配CPU时间给进程组。
- cpuacct：这个子系统在cpu子系统的基础上增加了对资源使用情况统计的功能。
- blkio：这个子系统主要用于限制进程组对块设备的访问。
- memory：这个子系统用于限制进程组对内存资源的使用。
- devices：这个子系统用于限制进程组对设备的访问。
- freezer：这个子系统用于挂起或者恢复进程组。
- net_cls：这个子系统用于标记进程组产生的网络数据包，方便tc命令使用。
- ns：可以使不同cgroups下面的进程使用不同的namespace。

## 2. 环境准备

采用的是CentOS7.6版本，内核版本为3.10.0-1160.92.1.el7.x86_64。Go语言版本为1.19.10。

## 3. 开始编写

### 3.1 NameSpace隔离实现

```go
package main
func main() {
    cmd := exec.Command("sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
	    // 隔离 uts,ipc,pid,mount,user,network
	    Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWNET,
	    // 设置容器的UID和GID
            UidMappings: []syscall.SysProcIDMap{
			{
				// 容器的UID
				ContainerID: 1,
				// 宿主机的UID
				HostID: 0,
				Size:   1,
			},
	    },
	    GidMappings: []syscall.SysProcIDMap{
			{
				// 容器的GID
				ContainerID: 1,
				// 宿主机的GID
				HostID: 0,
				Size:   1,
			},
	    },
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
```
