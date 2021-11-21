package main

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

func goRunCommand(cmdStr string) (pid string) {
	log.Println("Run command [go run", cmdStr, "]")
	cmd := exec.Command("go", "run", cmdStr)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	pid = strconv.Itoa(cmd.Process.Pid)
	go func() {
		err = cmd.Wait()
		if err != nil {
			log.Println("Command finished with error: ", err)
		}
	}()
	log.Println("Pid =", pid)

	return
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

	done := make(chan bool)
	var writeEventHandler sync.Mutex
	lastModifiedTime := getFileModifiedTime(watchedFile)
	go func() {
		goPid := goRunCommand(watchedFile)
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
					writeEventHandler.Lock()
					pidChan := make(chan string)
					go func() {
						pidToKill := <-pidChan
						log.Println("pKill process", pidToKill)
						out, err := exec.Command("pkill", "-P", "-9", pidToKill).Output()
						if err != nil {
							log.Fatal(err)
						} else {
							log.Println(out)
						}

						pidChan <- goRunCommand(watchedFile)
					}()
					pidChan <- goPid
					goPid = <-pidChan
					lastModifiedTime = modifiedTime
					writeEventHandler.Unlock()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					done <- true
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
	<-done
}
