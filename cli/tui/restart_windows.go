//go:build windows

package main

import "errors"

// restartSelf belum didukung di Windows: pengguna menutup dan membuka lagi.
func restartSelf() error {
	return errors.New("restart otomatis belum didukung di Windows")
}
