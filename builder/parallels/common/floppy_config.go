package common

import (
	"github.com/mitchellh/packer/template/interpolate"
)

// FloppyConfig is configuration related to created floppy disks and attaching
// them to a Parallels virtual machine.
type FloppyConfig struct {
	FloppyFiles    []string `mapstructure:"floppy_files"`
	FloppyContents []string `mapstructure:"floppy_contents"`
}

func (c *FloppyConfig) Prepare(ctx *interpolate.Context) []error {
	if c.FloppyFiles == nil {
		c.FloppyFiles = make([]string, 0)
	}

	if c.FloppyContents == nil {
		c.FloppyContents = make([]string, 0)
	}

	return nil
}
