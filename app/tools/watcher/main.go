package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

func runCommand(cmdStr string, restart <-chan bool) {
	for {
		cmd := exec.Command("go", "run", cmdStr)

		// create a pipe for the output of the script
		cmdReader, err := cmd.StdoutPipe()
		if err != nil {
			log.Println("Error creating StdoutPipe for Cmd. ", err)
			return
		}

		scanner := bufio.NewScanner(cmdReader)
		go func() {
			for scanner.Scan() {
				log.Println("\t > ", scanner.Text())
			}
		}()

		err = cmd.Start()
		if err != nil {
			log.Fatal(err)
		}

		pid := strconv.Itoa(cmd.Process.Pid)
		log.Println("Started command [go run", cmdStr, "] pid =", pid)

		// Block until told restart
		<-restart

		if err = cmd.Process.Kill(); err != nil {
			log.Println("Failed to kill process ", pid, ", error ", err)
		}
		log.Println("Terminate command [go run", cmdStr, "] pid =", pid)

		// Wait until process is terminated. Ignore error.
		cmd.Wait()
	}
}

func getFileModifiedTime(fileName string) (lastModifiedTime time.Time) {
	file, err := os.Stat(fileName)
	if err != nil {
		log.Fatal(err)
	}
	lastModifiedTime = file.ModTime()
	log.Println("File", fileName, "modified at", lastModifiedTime)
	return
}

func main() {

	if len(os.Args) != 2 {
		log.Fatal(("Usage: gowatcher <file>"))
	}
	watchedFile := os.Args[1]
	log.Println("Watching", watchedFile)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	var watcherWaitGroup sync.WaitGroup
	lastModifiedTime := getFileModifiedTime(watchedFile)
	watcherWaitGroup.Add(1)
	go func() {
		defer watcherWaitGroup.Done()
		processRestart := make(chan bool)
		go runCommand(watchedFile, processRestart)
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {

					modifiedTime := getFileModifiedTime(watchedFile)
					diff := modifiedTime.Sub(lastModifiedTime)
					if diff.Seconds() < 5 {
						log.Println("Ignore event", event.Name)
						continue
					}

					log.Println("modified file:", event.Name)
					processRestart <- true
					lastModifiedTime = modifiedTime
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(watchedFile)
	if err != nil {
		log.Fatal(err)
	}

	watcherWaitGroup.Wait()
}
