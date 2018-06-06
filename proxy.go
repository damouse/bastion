package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"github.com/gliderlabs/ssh"
	"github.com/kr/pty"
)

func setWinsize(f *os.File, w, h int) {
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ), uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
}

func directProxy() {
	ssh.Handle(func(s ssh.Session) {
		log.Printf("New connection: %v\n", s.User())
		// I wonder if its possible to read a command from here first?

		io.WriteString(s, "pwd")

		cmd := exec.Command("ssh", "lhr-vpn")
		ptyReq, winCh, isPty := s.Pty()

		if isPty {
			cmd.Env = append(cmd.Env, fmt.Sprintf("TERM=%s", ptyReq.Term))
			f, err := pty.Start(cmd)
			if err != nil {
				panic(err)
			}

			go func() {
				for win := range winCh {
					setWinsize(f, win.Width, win.Height)
				}
			}()

			go io.Copy(f, s) // stdin
			io.Copy(s, f)    // stdout
		} else {
			io.WriteString(s, "No PTY requested.\n")
			s.Exit(1)
		}
	})

	log.Println("starting ssh server on port 2222...")
	log.Fatal(ssh.ListenAndServe(":2222", nil))
}

func forwarding(ctx ssh.Context, destinationHost string, destinationPort uint32) bool {
	log.Println("Hello")
	// log.Printf("Forwarding request: context %v, host: %v, port: %v\n", ctx, destinationHost, destinationPort)
	return true
}

func paswordHandler(ctx ssh.Context, password string) bool {
	log.Printf("Password request. Context: %v, password: %v\n", ctx, password)
	return true
}

func publicKeyHandler(ctx ssh.Context, key ssh.PublicKey) bool {
	log.Println("Public key handler")
	return true
}

func handler(s ssh.Session) {
	log.Printf("New connection: %v\n", s.User())
	// I wonder if its possible to read a command from here first?

	io.WriteString(s, "pwd")

	cmd := exec.Command("ssh", "lhr-vpn")
	// cmd := exec.Command("pwd")

	ptyReq, winCh, isPty := s.Pty()

	if isPty {
		cmd.Env = append(cmd.Env, fmt.Sprintf("TERM=%s", ptyReq.Term))
		f, err := pty.Start(cmd)
		if err != nil {
			panic(err)
		}

		go func() {
			for win := range winCh {
				setWinsize(f, win.Width, win.Height)
			}
		}()

		go io.Copy(f, s) // stdin
		io.Copy(s, f)    // stdout
	} else {
		io.WriteString(s, "No PTY requested.\n")
		s.Exit(1)
	}
}

func fullServer() {
	s := ssh.Server{
		Addr:                        ":2222",
		Handler:                     handler,
		PublicKeyHandler:            publicKeyHandler,
		LocalPortForwardingCallback: forwarding,
	}

	log.Fatal(s.ListenAndServe())
}

func main() {
	// directProxy()
	fullServer()
}
