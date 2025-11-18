package virt

import (
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type EngineArch int

const (
	Qemu_system_x86_64 EngineArch = iota
	Qemu_system_arm
	Qemu_system_m68k
	Qemu_system_i386
	Qemu_system_mips
	Qemu_system_s390x
	Qemu_system_ppc32
	Qemu_system_ppc64
	Qemu_system_riscv32
	Qemu_system_riscv64
	Qemu_system_aarch64
)

func (n EngineArch) String() string {
	switch n {
	default:
		return "qemu-system-x86_64"
	case Qemu_system_arm:
		return "qemu-system-arm"
	case Qemu_system_aarch64:
		return "qemu-system-aarch64"
	case Qemu_system_x86_64:
		return "qemu-system-x86_64"
	case Qemu_system_i386:
		return "qemu-system-i386"
	case Qemu_system_m68k:
		return "qemu-system-m68k"
	case Qemu_system_mips:
		return "qemu-system-mips"
	case Qemu_system_ppc32:
		return "qemu-system-ppc32"
	case Qemu_system_ppc64:
		return "qemu-system-ppc64"
	case Qemu_system_riscv32:
		return "qemu-system-riscv32"
	case Qemu_system_riscv64:
		return "qemu-system-riscv64"
	case Qemu_system_s390x:
		return "qemu-system-s390x"
	}
}
func (n EngineArch) MarshalYAML() (any, error) {
	return n.String(), nil
}

func (n EngineArch) UnmarshalYAML(value *yaml.Node) error {
	switch strings.ToLower(value.Value) {
	default:
		return fmt.Errorf("status inválido: %s", value.Value)
	case "qemu-system-arm":
		n = Qemu_system_arm
	case "qemu-system-aarch64":
		n = Qemu_system_aarch64
	case "qemu-system-x86_64":
		n = Qemu_system_x86_64
	case "qemu-system-i386":
		n = Qemu_system_i386
	case "qemu-system-m68k":
		n = Qemu_system_m68k
	case "qemu-system-mips":
		n = Qemu_system_mips
	case "qemu-system-ppc32":
		n = Qemu_system_ppc32
	case "qemu-system-ppc64":
		n = Qemu_system_ppc64
	case "qemu-system-riscv32":
		n = Qemu_system_riscv32
	case "qemu-system-riscv64":
		n = Qemu_system_riscv64
	case "qemu-system-s390x":
		n = Qemu_system_s390x
	}
	return nil
}
