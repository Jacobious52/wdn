package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const title = "Watch Do Notify"

var logfile = os.DevNull
var display bool

// notify shows an macOS notification
func notify(subtitle string, text string) {
	cmd := fmt.Sprintf("display notification \"%s\" with title \"%s\" subtitle \"%s\"", text, title, subtitle)

	err := exec.Command("osascript", "-e", cmd).Run()
	if err != nil {
		log.Println("Failed to show notification.", err)
	}
}

// do runs the cmd given by the user
func do(script string) string {
	out, err := exec.Command("bash", "-c", script).Output()
	if err != nil {
		log.Println("Failed to run commmand.", err)
	}
	return string(out)
}

// getMod returns last modified time of file
func getMod(file string) time.Time {
	info, err := os.Stat(file)
	if err != nil {
		log.Printf("%s failed to access. %s\n", file, err)
		return time.Time{}
	}
	return info.ModTime()
}

// watch starts watching the files given in a new thread
func watch(files []string, update chan bool) {
	ticker := time.NewTicker(time.Second)

	lastMod := make([]time.Time, len(files))
	for i, f := range files {
		lastMod[i] = getMod(f)
	}

	go func() {
		// every second
		for range ticker.C {
			updated := false
			for i, f := range files {
				mod := getMod(f)
				if mod != lastMod[i] {
					lastMod[i] = mod
					log.Println(f, "modified.")

					// we only want to run the script once each round.
					// to stop a "save-all" from running too many times
					if !updated {
						log.Println("running script\n", f)
						updated = true
						update <- true
					}
				}
			}
		}
	}()
}

func main() {
	// command line flags
	name := flag.String("name", "Saved", "string: name of command")
	script := flag.String("cmd", "echo hello", "string: shell script to run")
	logFilename := flag.String("log", "/dev/null", "string: logging output file")
	display := flag.Bool("notify", false, "bool: should push macOS notification?")

	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("need at least one file to watch")
		return
	}

	files := flag.Args()

	// set logging file
	if *logFilename == "/dev/null" {
		logfile = os.DevNull
	}

	logfile, err := os.Create(*logFilename)
	log.SetOutput(logfile)
	if err != nil {
		logfile = os.Stderr
		log.SetOutput(logfile)
		log.Println("failed to open log file.", err)
	}

	log.Printf("running on files {%s}\n", strings.Join(files, ", "))

	// do main loop waiting on concurrent watch channel
	update := make(chan bool)
	watch(files, update)
	for {
		<-update // wait for an update
		output := do(*script)
		if *display {
			notify(*name, output)
		}
	}
}
