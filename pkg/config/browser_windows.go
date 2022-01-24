//go:build windows
// +build windows

package config

import (
	"log"
	"syscall"
)

func (b *Browser) Launch(url string) error {
	// NOTE: I tried using os/exec.Command but hit issues with args :-(

	var sI syscall.StartupInfo
	var pI syscall.ProcessInformation

	exe := b.Exe
	if exe[0] != '"' {
		exe = "\"" + exe + "\""
	}
	runCommand := exe
	if b.Args != nil {
		runCommand += " " + *b.Args
	}
	if url != "" {
		runCommand += " " + url
	}

	log.Println(runCommand)
	argv := syscall.StringToUTF16Ptr(runCommand) //nolint:staticcheck
	err := syscall.CreateProcess(
		nil,
		argv,
		nil,
		nil,
		true,
		0,
		nil,
		nil,
		&sI,
		&pI)
	return err
}
