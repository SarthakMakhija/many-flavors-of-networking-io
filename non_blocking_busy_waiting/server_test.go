package non_blocking_busy_waiting

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net"
	"non_blocking_busy_waiting/conn"
	"non_blocking_busy_waiting/proto"
	"testing"
)

func TestSendsAPutOrUpdateAndGetOverAConnection(t *testing.T) {
	port := randomPort()
	server, err := NewTCPServer("127.0.0.1", port)
	assert.Nil(t, err)

	go func() {
		server.Start()
	}()

	defer func() {
		server.Stop()
	}()

	connection, err := net.Dial("tcp", fmt.Sprintf("%v:%v", "127.0.0.1", port))
	assert.Nil(t, err)

	buffer, _ := proto.NewPutOrUpdateKeyValueMessage("DiskType", "NVMe SSD").Serialize()
	_, _ = connection.Write(buffer)

	buffer, _ = proto.NewPutOrUpdateKeyValueMessage("KV", "Distributed").Serialize()
	_, _ = connection.Write(buffer)

	buffer, _ = proto.NewGetValueMessage("KV").Serialize()
	_, _ = connection.Write(buffer)

	connectionReader := conn.NewConnectionReader(connection)
	_, _ = connectionReader.AttemptReadOrErrorOut()
	_, _ = connectionReader.AttemptReadOrErrorOut()

	message, err := connectionReader.AttemptReadOrErrorOut()

	assert.Nil(t, err)
	assert.Equal(t, "Distributed", message.Value)
}

func TestSendsMultiplePutOrUpdateAndAGetOverAConnection(t *testing.T) {
	port := randomPort()
	server, err := NewTCPServer("127.0.0.1", port)
	assert.Nil(t, err)

	go func() {
		server.Start()
	}()

	defer func() {
		server.Stop()
	}()

	connection, err := net.Dial("tcp", fmt.Sprintf("%v:%v", "127.0.0.1", port))
	assert.Nil(t, err)

	sendMultiplePutOrUpdates := func() {
		buffer, _ := proto.NewPutOrUpdateKeyValueMessage("DiskType", "NVMe SSD").Serialize()
		_, _ = connection.Write(buffer)

		buffer, _ = proto.NewPutOrUpdateKeyValueMessage("Storage", "LSM").Serialize()
		_, _ = connection.Write(buffer)

		buffer, _ = proto.NewPutOrUpdateKeyValueMessage("System", "Distributed").Serialize()
		_, _ = connection.Write(buffer)

		buffer, _ = proto.NewGetValueMessage("System").Serialize()
		_, _ = connection.Write(buffer)
	}
	attemptLastRead := func() (*proto.KeyValueMessage, error) {
		connectionReader := conn.NewConnectionReader(connection)
		_, _ = connectionReader.AttemptReadOrErrorOut()
		_, _ = connectionReader.AttemptReadOrErrorOut()
		_, _ = connectionReader.AttemptReadOrErrorOut()

		return connectionReader.AttemptReadOrErrorOut()
	}

	sendMultiplePutOrUpdates()
	//time.Sleep(100 * time.Second)
	message, err := attemptLastRead()

	assert.Nil(t, err)
	assert.Equal(t, "Distributed", message.Value)
}

func randomPort() uint16 {
	port := 0
	for port = rand.Intn(10000); port < 2000; port = rand.Intn(10000) {
		continue
	}
	return uint16(port)
}
