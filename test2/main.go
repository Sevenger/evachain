package main

import "os/exec"

func main() {
	for i := 0; i < 100; i++ {
		go func() {
			for {
				exec.Command("curl", "-X", "POST", "http://localhost:8080/newBlock", "-d", "data=123456").Run()
			}
		}()
	}
	for {

	}
}
