package main

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"time"
	"log"
	"strings"
	"os"
	"bufio"
	"io"
	"github.com/nettyrnp/go-fs/models"
	"github.com/nettyrnp/go-fs/storage"
	"os/signal"
	"syscall"
)

const (
	limit = 40 // Batch size
	fname1 = "data/name1.log"
	fname2 = "data/name2.log"
)

func main() {
	// Handle termination signals
	errors := make(chan error, 1)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, syscall.SIGTERM)

	// New file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("Error", err)
	}
	defer watcher.Close()
	recordsChan := make(chan []models.LogRecord, 1)

	// Save records to db
	go func(<-chan []models.LogRecord) {
		for {
			for records := range recordsChan {
				storage.Save(records)
			}
		}
	}(recordsChan)

	go listenToFile(recordsChan, watcher, fname1)
	go listenToFile(recordsChan, watcher, fname2)

	if err := watcher.Add(fname1); err != nil {
		fmt.Println("Error", err)
	}
	if err := watcher.Add(fname2); err != nil {
		fmt.Println("Error", err)
	}

	// Wait until the application fails or it is terminated.
	select {
	case err := <-errors:
		// Handle the error from the application
		log.Printf("Error from the application: %v\n", err)
		break
	case sig := <-signals:
		// Handle shutdown signals
		log.Printf("Signal: %v\n", sig)
		break
	}

	// Terminate gracefully
	i := 3
	for i > 0 {
		log.Printf("Terminating the application in %d s\n", i)
		time.Sleep(1 * time.Second)
		i = i - 1
	}
}

func listenToFile(ch chan<- []models.LogRecord, watcher *fsnotify.Watcher, fname string) {
	var offset = 0

	logFile, err := os.Open(fname)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer logFile.Close()
	reader := bufio.NewReader(logFile)

	// Initial reading of the file
	records, offset, err := readLines(reader, offset, fname)
	ch <- records

	for {
		select {
		// Watch for events
		case event := <-watcher.Events:
			if event.Op & fsnotify.Write == fsnotify.Write {
				fname0 := normalize(event.Name)
				if fname0 == fname {
					log.Println("New writing in file:", fname)

					// Wait till writing finishes
					time.Sleep(100 * time.Millisecond)
					records, offset, err = readLines(reader, offset, fname)
					if len(records)>0 {
						ch <- records
					}
				}
			}
		// Watch for errors
		case err := <-watcher.Errors:
			if err != nil {
				log.Println("Error:", err)
			}
		}
	}
}

func normalize(s string) string {
	s = strings.Replace(s, "\\", "/", -1)
	return s
}

func readLines(reader *bufio.Reader, offset int, fname string) ([]models.LogRecord, int, error) {
	var records []models.LogRecord
	for {
		// Read lines from file
		records0 := readLines0(reader, offset, fname)
		if len(records0) == 0 {
			return records, offset, nil
		}
		records = append(records, records0...)
		log.Printf("Reading: loaded %d lines from file '%s'", len(records0), fname)

		offset = offset+len(records0)
	}
	panic("Something went wrong. Should never reach this line")
}


func readLines0(reader *bufio.Reader, offset int, fname string) []models.LogRecord{
	var records = []models.LogRecord{}
	var count = offset
	for count <= offset+limit {
		var error error
		var buf, buf0 []byte
		var hasMore = true
		for hasMore {
			buf0, hasMore, error = reader.ReadLine()
			buf = append(buf, buf0...)
		}
		line := string(buf)

		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		count++
		if len(line) == 0 {
			continue
		}
		p := models.NewRecord(line, fname)
		records = append(records, p)
	}
	return records
}
