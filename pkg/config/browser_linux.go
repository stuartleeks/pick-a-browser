//go:build linux
// +build linux

package config

import (
	"fmt"
)

func (b *Browser) Launch(url string) error {
	return fmt.Errorf("Linux Launch implementation only for testing")
}
