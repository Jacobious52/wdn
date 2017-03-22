package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const title = "Watch Do Notify"

func notify(subtitle string, text string) {
	cmd := fmt.Sprintf("display notification \"%s\" with title \"%s\" subtitle \"%s\"", text, title, subtitle)

	err := exec.Command("osascript", "-e", cmd).Run()
	if err != nil {
		log.Println("Failed to show notification.", err)
		log.Println(cmd)
	}
}

func do(script string) string {
	out, err := exec.Command("bash", "-c", script).Output()
	if err != nil {
		log.Println("Failed to run commmand.", err)
	}
	return string(out)
}

func getMod(file string) time.Time {
	info, err := os.Stat(file)
	if err != nil {
		log.Printf("%s failed to access. %s\n", file, err)
		return time.Time{}
	}
	return info.ModTime()
}

func watch(files []string, update chan bool) {
	ticker := time.NewTicker(time.Second)

	lastMod := make([]time.Time, len(files))
	for i, f := range files {
		lastMod[i] = getMod(f)
	}

	go func() {
		for range ticker.C {
			for i, f := range files {
				mod := getMod(f)
				if mod != lastMod[i] {
					lastMod[i] = mod
					log.Printf("%v modified. running script\n", f)
					update <- true
					break
				}
			}
		}
	}()
}

func main() {

	if len(os.Args) < 4 {
		fmt.Println("usage: wdn \"name\" \"cmd args\" files...")
		return
	}
	name := os.Args[1]
	script := os.Args[2]
	files := os.Args[3:len(os.Args)]

	fmt.Printf("running on files {%s}\n", strings.Join(files, ", "))

	update := make(chan bool)
	watch(files, update)
	for {
		<-update
		output := do(script)
		notify(name, output)
	}
}
