//go:build !windows

package main

import (
	"os"
	"syscall"
)

// restartSelf mengganti proses ini dengan binary cipincode (versi terbaru
// yang baru dipasang). Hanya untuk Unix; Windows memakai fallback manual.
func restartSelf() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	return syscall.Exec(exe, os.Args, os.Environ())
}

// restartSelfFn bisa diganti di tes agar tidak benar-benar exec.
var restartSelfFn = restartSelf
