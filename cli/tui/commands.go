package main

import "strings"

type command struct {
	name string
	desc string
}

// commands adalah isi menu yang muncul saat mengetik "/".
var commands = []command{
	{"help", "Tampilkan bantuan"},
	{"chat", "Pindah ke mode Chat"},
	{"coding", "Pindah ke mode Coding (auto-approve)"},
	{"debug", "Pindah ke mode Debugging"},
	{"new", "Mulai percakapan baru"},
	{"clear", "Bersihkan layar"},
	{"status", "Lihat status akun, agent, dan mode"},
	{"usage", "Lihat sisa jatah sesi 5 jam dan mingguan"},
	{"effort", "Atur effort: low, medium, high, xhigh, ultra"},
	{"agent", "Pilih agent buatan Anda"},
	{"whoami", "Lihat akun yang sedang login"},
	{"login", "Masuk atau ganti akun"},
	{"logout", "Keluar dari akun"},
	{"dashboard", "Buka dashboard di browser"},
	{"keys", "Buka halaman API Keys"},
	{"setup", "Buka panduan instalasi"},
	{"version", "Tampilkan versi CLI"},
	{"exit", "Keluar dari CipinCode"},
}

// matchCommands menyaring perintah dari teks input. Awalan nama diutamakan,
// lalu yang sekadar memuat teks itu. Mengembalikan nil jika input bukan "/perintah".
func matchCommands(input string) []command {
	if !strings.HasPrefix(input, "/") || strings.ContainsAny(input, " \t") {
		return nil
	}
	query := strings.ToLower(strings.TrimPrefix(input, "/"))

	var prefix, contains []command
	for _, c := range commands {
		switch {
		case strings.HasPrefix(c.name, query):
			prefix = append(prefix, c)
		case strings.Contains(c.name, query):
			contains = append(contains, c)
		}
	}
	return append(prefix, contains...)
}
