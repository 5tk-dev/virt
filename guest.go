package virt

import (
	"strings"

	"github.com/google/uuid"
)

type Guest struct {
	// Guest Identity
	Name string // domain name
	UUID string // domain uuid

	//
	Engine EngineArch // default qemu-system-x86_64

	// Memory Options
	Memory      *MemoryOptions `yaml:",omitempty"` // memory spec
	MemPath     string         `yaml:",omitempty"` // -mem-path FILE  provide backing storage for guest RAM
	MemPrealloc string         `yaml:",omitempty"` // -mem-prealloc   preallocate guest memory (use with -mem-path)

	// Process
	Smp *SmpOptions `yaml:",omitempty"` // vcpu spec

	// Storage
	BlockDevices *BlockDevicesOptions `yaml:",omitempty"`

	// NETWORK
	Nic               *NicOptions               `yaml:",omitempty"`
	NetDev_Bridge     *NetDev_BridgeOptions     `yaml:",omitempty"`
	Netdev_Hubport    *Netdev_HubportOptions    `yaml:",omitempty"`
	Netdev_Passt      *Netdev_PasstOptions      `yaml:",omitempty"`
	Netdev_Tap        *Netdev_TapOptions        `yaml:",omitempty"`
	Netdev_User       *Netdev_UserOptions       `yaml:",omitempty"`
	Netdev_Vde        *Netdev_VdeOptions        `yaml:",omitempty"`
	Netdev_Vhost_user *Netdev_Vhost_userOptions `yaml:",omitempty"`
	Netdev_Vhost_vdpa *Netdev_Vhost_vdpaOptions `yaml:",omitempty"`

	//

	Devices []*DeviceOptions

	//
	Qmp       *QmpOptions `yaml:",omitempty"`
	Daemonize bool        `yaml:",omitempty"`

	// GENERAL OPTIONS
	K         string `yaml:",omitempty"` // -k language - use keyboard layout (for example 'fr' for French)
	NoGraphic bool   `yaml:",omitempty"` // -nographic  - disable graphical output and redirect serial I/Os to console

}

func (g *Guest) ToArgs() []string {
	if g.UUID == "" {
		g.UUID = uuid.NewString()
	}
	args := []string{
		g.Engine.String(),
		"-name", g.Name,
		"-uuid", g.UUID,
	}

	if g.Memory != nil {
		args = append(args, g.Memory.ToArgs()...)
	}
	if g.Smp != nil {
		args = append(args, g.Smp.ToArgs()...)
	}
	if g.BlockDevices != nil {
		args = append(args, g.BlockDevices.ToArgs()...)
	}
	if g.Qmp != nil {
		args = append(args, g.Qmp.ToArgs()...)
	}
	if g.Nic != nil {
		args = append(args, g.Nic.ToArgs()...)
	}
	if g.NetDev_Bridge != nil {
		args = append(args, g.NetDev_Bridge.ToArgs()...)
	}
	if g.Netdev_Hubport != nil {
		args = append(args, g.Netdev_Hubport.ToArgs()...)
	}
	if g.Netdev_Passt != nil {
		args = append(args, g.Netdev_Passt.ToArgs()...)
	}
	if g.Netdev_Tap != nil {
		args = append(args, g.Netdev_Tap.ToArgs()...)
	}
	if g.Netdev_User != nil {
		args = append(args, g.Netdev_User.ToArgs()...)
	}
	if g.Netdev_Vde != nil {
		args = append(args, g.Netdev_Vde.ToArgs()...)
	}
	if g.Netdev_Vhost_user != nil {
		args = append(args, g.Netdev_Vhost_user.ToArgs()...)
	}
	if g.Netdev_Vhost_vdpa != nil {
		args = append(args, g.Netdev_Vhost_vdpa.ToArgs()...)
	}
	for _, d := range g.Devices {
		args = append(args, d.ToArgs()...)
	}
	if g.Daemonize {
		args = append(args, "-daemonize")
	}

	return args
}

// -qmp stdio
// -qmp tcp:127.0.0.1:4444,server,nowait
// -qmp unix:/tmp/qmp-sock,server,nowait
type QmpOptions struct {
	ProtoPath string // unix:path | tcp:host:port | stdio
	Wait      bool   //
	Serve     bool
}

func (q *QmpOptions) ToArgs() []string {
	args := []string{}
	if p := strings.ToLower(q.ProtoPath); p == "stdio" {
		args = append(args, p)
	} else {
		args = append(args, q.ProtoPath)
		if q.Serve {
			args = append(args, "server")
		}
		if q.Wait {
			args = append(args, "wait")
		} else {
			args = append(args, "nowait")
		}
	}
	return []string{"-qmp", strings.Join(args, ",")}
}
