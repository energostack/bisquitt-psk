package reader

import (
	"context"
	"encoding/csv"
	"io"
	"os"

	"bisquitt-psk/pkg/clientmap"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
)

// CSVReader is a reader that reads clients from a CSV file.
type CSVReader struct {
	FilePath  string
	updatesCh chan bool
}

// NewCSVReader creates a new CSVReader with the specified file path.
func NewCSVReader(filePath string, ctx context.Context) (*CSVReader, error) {
	reader := CSVReader{
		FilePath:  filePath,
		updatesCh: make(chan bool),
	}
	err := reader.watch(ctx)
	if err != nil {
		return nil, err
	}
	return &reader, nil
}

// Read reads clients from a CSV file.
func (cr *CSVReader) Read() (*clientmap.Map, error) {
	f, err := os.Open(cr.FilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	clientMap := clientmap.New()

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if len(record) != 2 {
			log.Warn().Msgf("Invalid record: %v. Rows should be in format <clientID>, <PSK>", record)
			continue
		}
		clientID := record[0]
		psk := []byte(record[1])

		clientMap.Store(clientID, psk)

	}

	return clientMap, nil
}

// Updates returns a channel that receives updates when the CSV file changes.
func (cr *CSVReader) Updates() <-chan bool {
	return cr.updatesCh
}

func (cr *CSVReader) watch(ctx context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Err(err)
		return err
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					continue
				}

				if event.Has(fsnotify.Write) {
					cr.updatesCh <- true
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					log.Err(err)
					continue
				}
			case <-ctx.Done():
				watcher.Close()
				return
			}
		}

	}()

	err = watcher.Add(cr.FilePath)
	if err != nil {
		log.Err(err)
		return err
	}
	return nil
}
