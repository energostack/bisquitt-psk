package reader

import (
	"context"
	"database/sql"

	"bisquitt-psk/pkg/clientmap"

	"github.com/fsnotify/fsnotify"
	_ "github.com/glebarez/go-sqlite"
	"github.com/rs/zerolog/log"
)

// SQLiteReader is a reader that reads clients from a SQLite database.
type SQLiteReader struct {
	FilePath  string
	updatesCh chan bool
}

// NewSQLiteReader creates a new SQLiteReader with the specified file path.
func NewSQLiteReader(filePath string, ctx context.Context) (*SQLiteReader, error) {
	reader := SQLiteReader{
		FilePath:  filePath,
		updatesCh: make(chan bool),
	}
	err := reader.watch(ctx)
	if err != nil {
		return nil, err

	}
	return &reader, nil
}

// Read reads clients from a SQLite database.
func (cr *SQLiteReader) Read() (*clientmap.Map, error) {
	db, err := sql.Open("sqlite", cr.FilePath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT client_id, psk FROM clients")
	if err != nil {
		return nil, err
	}

	clientMap := clientmap.New()

	for rows.Next() {
		var clientID string
		var psk []byte

		err = rows.Scan(&clientID, &psk)
		if err != nil {
			return nil, err
		}

		clientMap.Store(clientID, psk)
	}

	return clientMap, nil
}

// Updates returns a channel that receives updates when the SQLite database changes.
func (cr *SQLiteReader) Updates() <-chan bool {
	return cr.updatesCh
}

func (cr *SQLiteReader) watch(ctx context.Context) error {
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
