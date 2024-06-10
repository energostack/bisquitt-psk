package tests

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/IBM/sarama"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}
	assert := assert.New(t)
	require := require.New(t)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		client := mqtt.NewClient(mqtt.NewClientOptions().AddBroker(os.Getenv("MQTT_BROKER")))
		token := client.Connect()
		token.Wait()
		require.NoError(token.Error())

		client.Subscribe(os.Getenv("MQTT_TOPIC"), 0, func(client mqtt.Client, msg mqtt.Message) {
			require.Equal(os.Getenv("MQTT_TOPIC"), msg.Topic())
			require.Equal("Hello, World!", string(msg.Payload()))
			wg.Done()
		})
	}()

	producer, err := sarama.NewSyncProducer([]string{os.Getenv("KAFKA_BROKER")}, nil)
	require.NoError(err)

	producer.SendMessage(&sarama.ProducerMessage{
		Topic: os.Getenv("KAFKA_TOPIC"),
		Key:   sarama.StringEncoder("bisquitt"),
		Value: sarama.StringEncoder("psk"),
	})

	producer.Close()

	time.Sleep(15 * time.Second)

	httpClient := &http.Client{}

	requestURL, _ := url.Parse(fmt.Sprintf("%s%s", os.Getenv("SERVICE_URL"), "/clients/bisquitt"))
	resp, err := httpClient.Do(&http.Request{
		Method: http.MethodGet,
		URL:    requestURL,
		Header: http.Header{
			"Authorization": []string{fmt.Sprintf("Basic %s",
				base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(
					"%s:%s", os.Getenv("BASIC_AUTH_USERNAME"), os.Getenv("BASIC_AUTH_PASSWORD")),
				),
				))},
		},
	})

	require.NoError(err)
	require.Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(err)
	resp.Body.Close()
	assert.Equal(`{"client":"bisquitt","psk":"cHNr"}`, string(body))
	wg.Wait()
}
