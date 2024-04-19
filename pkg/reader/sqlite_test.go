package reader_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"bisquitt-psk/pkg/reader"

	_ "github.com/glebarez/go-sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestDB() {
	db, err := sql.Open("sqlite", "test.db")
	if err != nil {
		panic(err)
	}
	db.Exec("CREATE TABLE clients (client_id TEXT PRIMARY KEY, psk TEXT)")
	db.Exec("INSERT INTO clients (client_id, psk) VALUES ('client1', 'psk1')")
	db.Exec("INSERT INTO clients (client_id, psk) VALUES ('client2', 'psk2')")
	db.Close()
}

func TestSQLiteReaderReadsValidData(t *testing.T) {
	createTestDB()
	sqliteReader, err := reader.NewSQLiteReader("test.db", context.Background())
	require.NoError(t, err)

	clientMap, _ := sqliteReader.Read()

	value, _ := clientMap.Load("client1")
	assert.Equal(t, []byte("psk1"), value)

	value, _ = clientMap.Load("client2")
	assert.Equal(t, []byte("psk2"), value)

	os.Remove("test.db")
}

func TestSQLiteReaderUpdatesOnFileChange(t *testing.T) {
	createTestDB()
	sqliteReader, err := reader.NewSQLiteReader("test.db", context.Background())
	require.NoError(t, err)

	db, err := sql.Open("sqlite", "test.db")
	if err != nil {
		assert.NoError(t, err)
	}
	db.Exec("INSERT INTO clients (client_id, psk) VALUES ('client3', 'psk3')")
	db.Close()

	<-sqliteReader.Updates()

	clientMap, _ := sqliteReader.Read()

	value, _ := clientMap.Load("client1")
	assert.Equal(t, []byte("psk1"), value)

	value, _ = clientMap.Load("client2")
	assert.Equal(t, []byte("psk2"), value)

	value, _ = clientMap.Load("client3")
	assert.Equal(t, []byte("psk3"), value)
	os.Remove("test.db")
}
