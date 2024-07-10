package reader_test

import (
	"context"
	"os"
	"testing"

	"github.com/energostack/bisquitt-psk/pkg/reader"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCSVReaderReadsValidFile(t *testing.T) {
	filePath := "test.csv"
	file, _ := os.Create(filePath)
	file.WriteString("client1,psk1\nclient2,psk2")
	file.Close()

	csvReader, err := reader.NewCSVReader(filePath, context.Background())
	require.NoError(t, err)
	clientMap, _ := csvReader.Read()

	value, _ := clientMap.Load("client1")
	assert.Equal(t, []byte("psk1"), value)

	value, _ = clientMap.Load("client2")
	assert.Equal(t, []byte("psk2"), value)

	os.Remove(filePath)
}

func TestCSVInvalid(t *testing.T) {
	filePath := "test.csv"
	file, _ := os.Create(filePath)
	file.WriteString("client1,psk1\ninvalid\nclient2,psk2")
	file.Close()

	csvReader, err := reader.NewCSVReader(filePath, context.Background())
	require.NoError(t, err)
	_, err = csvReader.Read()

	assert.Error(t, err)
	os.Remove(filePath)
}

func TestCSVReaderReturnsErrorOnInvalidFile(t *testing.T) {
	csvReader, err := reader.NewCSVReader("nonexistent.csv", context.Background())
	require.Error(t, err)
	require.Nil(t, csvReader)
}

func TestCSVReaderUpdatesOnFileChange(t *testing.T) {
	filePath := "test.csv"
	file, _ := os.Create(filePath)
	file.WriteString("client1,psk1")
	file.Close()

	csvReader, err := reader.NewCSVReader(filePath, context.Background())
	require.NoError(t, err)

	file, _ = os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	file.WriteString("\nclient2,psk2")
	file.Close()

	<-csvReader.Updates()

	clientMap, _ := csvReader.Read()

	value, _ := clientMap.Load("client1")
	assert.Equal(t, []byte("psk1"), value)

	value, _ = clientMap.Load("client2")
	assert.Equal(t, []byte("psk2"), value)

	os.Remove(filePath)
}
