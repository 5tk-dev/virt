package virt

import (
	"fmt"
	"strings"
)

type SmpOptions struct {
	Cpus     int `yaml:",omitempty"` // cpus=<num>
	Dies     int `yaml:",omitempty"` // dies=<num>
	Cores    int `yaml:",omitempty"` // cores=<num>
	Books    int `yaml:",omitempty"` // books=<num>
	Drawers  int `yaml:",omitempty"` // drawers=<num>
	Maxcpus  int `yaml:",omitempty"` // maxcpus=<num>
	Modules  int `yaml:",omitempty"` // modules=<num>
	Sockets  int `yaml:",omitempty"` // sockets=<num>
	Threads  int `yaml:",omitempty"` // threads=<num>
	Clusters int `yaml:",omitempty"` // clusters=<num>
}

func (s *SmpOptions) ToArgs() []string {
	args := []string{}
	if s.Cpus > 0 {
		args = append(args, fmt.Sprintf("cpus=%d", s.Cpus))
	}
	if s.Dies > 0 {
		args = append(args, fmt.Sprintf("dies=%d", s.Dies))
	}
	if s.Cores > 0 {
		args = append(args, fmt.Sprintf("cores=%d", s.Cores))
	}
	if s.Books > 0 {
		args = append(args, fmt.Sprintf("books=%d", s.Books))
	}
	if s.Drawers > 0 {
		args = append(args, fmt.Sprintf("drawers=%d", s.Drawers))
	}
	if s.Maxcpus > 0 {
		args = append(args, fmt.Sprintf("maxcpus=%d", s.Maxcpus))
	}
	if s.Modules > 0 {
		args = append(args, fmt.Sprintf("modules=%d", s.Modules))
	}
	if s.Sockets > 0 {
		args = append(args, fmt.Sprintf("sockets=%d", s.Sockets))
	}
	if s.Threads > 0 {
		args = append(args, fmt.Sprintf("threads=%d", s.Threads))
	}
	if s.Clusters > 0 {
		args = append(args, fmt.Sprintf("clusters=%d", s.Clusters))
	}
	return []string{"-smp", strings.Join(args, ",")}
}
