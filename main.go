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
	"errors"
)

const (
	limit = 40 // max lines to load in single round
	fname1 = "data/name1.log"
	fname2 = "data/name2.log"
)

func main() {
	// Handle termination signals
	errors := make(chan error, 1)

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

	// Wait until the service fails or it is terminated.
	select {
	case err := <-errors:
		// Handle the error from ingestor or persistor
		log.Printf("Error from ingestor or persistor: %v\n", err)
		break
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
		// watch for events
		case event := <-watcher.Events:
			if event.Op & fsnotify.Write == fsnotify.Write {
				fname0 := normalize(event.Name)
				//log.Println("New write event to file:", fname)
				// New lines appeared in the log file, so read them
				//println("\t fname0 == fname:", fname0 == fname)
				//println("\t fname0:", fname0)
				//println("\t fname:", fname)
				if fname0 == fname {
					records, offset, err = readLines(reader, offset, fname)
					println("\t len(records):", len(records), "; offset:", offset, "; fname0:", fname0)
					if len(records)>0 {
						ch <- records
					}
				}
			}
		// watch for errors
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
	var records0 []models.LogRecord
	for {
		// Read lines from file
		records := readLines0(reader, offset, fname)
		if len(records) == 0 {
			return records0, offset, nil
		}
		records0 = append(records0, records...)
		log.Printf("Reading: loaded %d lines from file '%s'", len(records), fname)

		offset = offset+len(records)

		// Slow down to better see the effect of <Ctrl+C>
		time.Sleep(200 * time.Millisecond)
	}
	return nil, -1, errors.New("sd")
}


func readLines0(reader *bufio.Reader, offset int, fname string) []models.LogRecord{
	var records = []models.LogRecord{}
	var count = offset
	for count <= offset+limit {
		var error error
		var buf, buf2 []byte
		var hasMore = true
		for hasMore {
			buf2, hasMore, error = reader.ReadLine()
			buf = append(buf, buf2...)
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
		//log.Println("Reading: read line: '" + line + "'")
		p := models.NewRecord(line, fname)
		records = append(records, p)
	}
	return records
}
