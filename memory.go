package virt

import (
	"fmt"
	"strings"
)

type MemoryOptions struct {
	Size   int    `yaml:",omitempty"` // megabyte
	Slots  string `yaml:",omitempty"`
	Maxmen int    `yaml:",omitempty"`
}

// original command:
//
//	-m [size=]megs[,slots=n,maxmem=size]
//		configure guest RAM
//		size: initial amount of guest memory
//		slots: number of hotplug slots (default: none)
//		maxmem: maximum amount of guest memory (default: none)
func (m *MemoryOptions) ToArgs() []string {
	args := []string{
		fmt.Sprintf("size=%d", m.Size),
	}
	if m.Slots != "" {
		args = append(args, fmt.Sprintf("slots=%s", m.Slots))
	}
	if m.Maxmen != 0 {
		args = append(args, fmt.Sprintf("maxmen=%d", m.Maxmen))
	}
	return []string{"-m", strings.Join(args, ",")}
}
