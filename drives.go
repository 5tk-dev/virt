package virt

import (
	"fmt"
	"strings"
)

// -cdrom file     use 'file' as CD-ROM image
type CdromOptions struct {
	File string `yaml:",omitempty"`
}

func (c *CdromOptions) ToArgs() []string { return []string{"-cdrom", c.File} }

/*
-drive [file=file][,if=type][,bus=n][,unit=m][,media=d][,index=i]

	[,cache=writethrough|writeback|none|directsync|unsafe][,format=f]
	[,snapshot=on|off][,rerror=ignore|stop|report]
	[,werror=ignore|stop|report|enospc][,id=name]
	[,aio=threads|native|io_uring]
	[,readonly=on|off][,copy-on-read=on|off]
	[,discard=ignore|unmap][,detect-zeroes=on|off|unmap]
	[[,bps=b]|[[,bps_rd=r][,bps_wr=w]]]
	[[,iops=i]|[[,iops_rd=r][,iops_wr=w]]]
	[[,bps_max=bm]|[[,bps_rd_max=rm][,bps_wr_max=wm]]]
	[[,iops_max=im]|[[,iops_rd_max=irm][,iops_wr_max=iwm]]]
	[[,iops_size=is]]
	[[,group=g]]
*/
type DriveOptions struct {
	File   string `yaml:",omitempty"` // [file=file]
	If     string `yaml:",omitempty"` // [,if=type]
	Bus    string `yaml:",omitempty"` // [,bus=n]
	Unit   string `yaml:",omitempty"` // [,unit=m]
	Media  string `yaml:",omitempty"` // [,media=d]
	Index  string `yaml:",omitempty"` // [,index=i]
	Format string `yaml:",omitempty"` // [,format=f]
}

func (d *DriveOptions) ToArgs() []string {
	args := []string{fmt.Sprintf("file=%s", d.File)}
	if d.If != "" {
		args = append(args, fmt.Sprintf("if=%s", d.If))
	}
	if d.Bus != "" {
		args = append(args, fmt.Sprintf("bus=%s", d.Bus))
	}
	if d.Unit != "" {
		args = append(args, fmt.Sprintf("unit=%s", d.Unit))
	}
	if d.Index != "" {
		args = append(args, fmt.Sprintf("index=%s", d.Index))
	}
	if d.Media != "" {
		args = append(args, fmt.Sprintf("media=%s", d.Media))
	}
	if d.Format != "" {
		args = append(args, fmt.Sprintf("format=%s", d.Format))
	}

	return []string{"-drive", strings.Join(args, ",")}
}

/*
original command:

	-blockdev [driver=]driver[,node-name=N][,discard=ignore|unmap]
		[,cache.direct=on|off]
		[,cache.no-flush=on|off]
		[,read-only=on|off]
		[,auto-read-only=on|off]
		[,force-share=on|off]
		[,detect-zeroes=on|off|unmap]
		[,driver specific parameters...]
	    		configure a block backend
*/
type BlockDev struct {
	Driver       string `yaml:",omitempty"` // [driver=]driver
	Discard      string `yaml:",omitempty"` // [,discard=ignore|unmap]
	ReadOnly     string `yaml:",omitempty"` //[,read-only=on|off]
	NodeName     string `yaml:",omitempty"` // [,node-name=N]
	CacheDirect  string `yaml:",omitempty"` // [,cache.direct=on|off]
	CacheNoFlush string `yaml:",omitempty"` //[,cache.no-flush=on|off]
	AutoReadOnly string `yaml:",omitempty"` //[,auto-read-only=on|off]
	ForceShare   string `yaml:",omitempty"` //[,force-share=on|off]
	DetectZeroes string `yaml:",omitempty"` //[,detect-zeroes=on|off|unmap]
}

func (b *BlockDev) ToArgs() []string {
	args := []string{}
	if b.Driver != "" {
		args = append(args, fmt.Sprintf("driver=%s", b.Driver))
	}
	if b.Discard != "" {
		args = append(args, fmt.Sprintf("discard=%s", b.Discard))
	}
	if b.ReadOnly != "" {
		args = append(args, fmt.Sprintf("read-only=%s", b.ReadOnly))
	}
	if b.NodeName != "" {
		args = append(args, fmt.Sprintf("node-name=%s", b.NodeName))
	}
	if b.CacheDirect != "" {
		args = append(args, fmt.Sprintf("cache.direct=%s", b.CacheDirect))
	}
	if b.CacheNoFlush != "" {
		args = append(args, fmt.Sprintf("cache.no-flush=%s", b.CacheNoFlush))
	}
	if b.AutoReadOnly != "" {
		args = append(args, fmt.Sprintf("auto-read-only=%s", b.AutoReadOnly))
	}
	if b.ForceShare != "" {
		args = append(args, fmt.Sprintf("force-share=%s", b.ForceShare))
	}
	if b.DetectZeroes != "" {
		args = append(args, fmt.Sprintf("detect-zeroes=%s", b.DetectZeroes))
	}

	return []string{"-blockdev", strings.Join(args, ",")}
}

type BlockDevicesOptions struct {
	Drive    []*DriveOptions `yaml:",omitempty"` // use 'file' as a drive image
	Cdrom    []*CdromOptions `yaml:",omitempty"` // use 'file' as CD-ROM image
	BlockDev *BlockDev       `yaml:",omitempty"` //
	Fda      string          `yaml:",omitempty"` // use 'file' as floppy disk 0 image
	Fdb      string          `yaml:",omitempty"` // use 'file' as floppy disk 1 image
	Hda      string          `yaml:",omitempty"` // use 'file' as hard disk 0 image
	Hdb      string          `yaml:",omitempty"` // use 'file' as hard disk 1 image
	Hdc      string          `yaml:",omitempty"` // use 'file' as hard disk 2 image
	Hdd      string          `yaml:",omitempty"` // use 'file' as hard disk 3 image
}

func (b *BlockDevicesOptions) ToArgs() []string {
	args := []string{}
	if b.Fda != "" {
		args = append(args, "-fda", b.Fda)
	}
	if b.Fdb != "" {
		args = append(args, "-fdb", b.Fdb)
	}
	if b.Hda != "" {
		args = append(args, "-hda", b.Hda)
	}
	if b.Hdb != "" {
		args = append(args, "-hdb", b.Hdb)
	}
	if b.Hdc != "" {
		args = append(args, "-hdc", b.Hdc)
	}
	if b.Hdd != "" {
		args = append(args, "-hdd", b.Hdd)
	}
	if b.BlockDev != nil {
		args = append(args, b.BlockDev.ToArgs()...)
	}
	for _, c := range b.Cdrom {
		args = append(args, c.ToArgs()...)
	}
	for _, d := range b.Drive {
		args = append(args, d.ToArgs()...)
	}
	return args
}
