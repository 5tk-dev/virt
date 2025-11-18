package virt

import "strings"

/*
original command:

-device driver[,prop[=value][,...]]

	                add device (based on driver)
	                prop=value,... sets driver properties
	                use '-device help' to print all possible drivers
	                use '-device driver,help' to print all possible properties
		ex:
*/
type DeviceOptions struct {
	Driver     string `yaml:",omitempty"` // driver
	Properties string `yaml:",omitempty"` // [,prop[=value][,...]]
}

func (d *DeviceOptions) ToArgs() []string {
	args := []string{d.Driver}
	if d.Properties != "" {
		args = append(args, d.Properties)
	}
	return []string{"-device", strings.Join(args, ",")}
}
