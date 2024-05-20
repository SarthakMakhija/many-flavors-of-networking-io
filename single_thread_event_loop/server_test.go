package single_thread_event_loop

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net"
	"single_thread_eventloop/conn"
	"single_thread_eventloop/proto"
	"testing"
	"time"
)

func randomPort() int {
	port := 0
	for port = rand.Intn(10000); port < 2000; port = rand.Intn(10000) {
		continue
	}
	return port
}

func TestSendsAPutOrUpdateAndGetOverAConnection(t *testing.T) {
	port := randomPort()
	server, err := NewTCPServer("127.0.0.1", uint16(port))
	assert.Nil(t, err)

	go func() {
		server.Start()
	}()

	connection, err := net.Dial("tcp", fmt.Sprintf("localhost:%v", port))
	assert.Nil(t, err)

	defer func() {
		server.Stop()
		if connection != nil {
			_ = connection.Close()
		}
	}()

	buffer, _ := proto.NewPutOrUpdateKeyValueMessage("DiskType", "NVMe SSD").Serialize()
	_, _ = connection.Write(buffer)

	time.Sleep(20 * time.Millisecond)

	buffer, _ = proto.NewGetValueMessage("DiskType").Serialize()
	_, _ = connection.Write(buffer)

	connectionReader := conn.NewConnectionReader(connection)
	_, _ = connectionReader.AttemptReadOrErrorOut()

	message, err := connectionReader.AttemptReadOrErrorOut()

	assert.Nil(t, err)
	assert.Equal(t, "NVMe SSD", message.Value)
}

func TestSendsMultiplePutOrUpdateAndAGetOverAConnection(t *testing.T) {
	port := randomPort()
	server, err := NewTCPServer("127.0.0.1", uint16(port))
	assert.Nil(t, err)

	go func() {
		server.Start()
	}()

	defer func() {
		server.Stop()
	}()

	connection, err := net.Dial("tcp", fmt.Sprintf("localhost:%v", port))
	assert.Nil(t, err)

	sendMultiplePutOrUpdates := func() {
		buffer, _ := proto.NewPutOrUpdateKeyValueMessage("DiskType", "NVMe SSD").Serialize()
		_, _ = connection.Write(buffer)
		time.Sleep(2 * time.Millisecond)

		buffer, _ = proto.NewPutOrUpdateKeyValueMessage("Storage", "LSM").Serialize()
		_, _ = connection.Write(buffer)
		time.Sleep(2 * time.Millisecond)

		buffer, _ = proto.NewPutOrUpdateKeyValueMessage("System", "Distributed").Serialize()
		_, _ = connection.Write(buffer)
		time.Sleep(2 * time.Millisecond)
	}
	attemptLastRead := func() (*proto.KeyValueMessage, error) {
		connectionReader := conn.NewConnectionReader(connection)
		_, _ = connectionReader.AttemptReadOrErrorOut()
		_, _ = connectionReader.AttemptReadOrErrorOut()
		_, _ = connectionReader.AttemptReadOrErrorOut()

		return connectionReader.AttemptReadOrErrorOut()
	}

	sendMultiplePutOrUpdates()

	buffer, _ := proto.NewGetValueMessage("System").Serialize()
	_, _ = connection.Write(buffer)

	message, err := attemptLastRead()

	assert.Nil(t, err)
	assert.Equal(t, "Distributed", message.Value)
}
