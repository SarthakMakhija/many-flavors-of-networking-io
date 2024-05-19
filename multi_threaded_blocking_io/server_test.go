package single_threaded_blocking_io

import (
	"github.com/stretchr/testify/assert"
	"many-flavors-of-nwing-io/multi_threaded_blocking_io/conn"
	proto2 "many-flavors-of-nwing-io/multi_threaded_blocking_io/proto"
	"net"
	"testing"
)

func TestSendsAPutOrUpdateAndGetOverAConnection(t *testing.T) {
	server, err := NewTCPServer("localhost", 9090)
	assert.Nil(t, err)

	go func() {
		server.Start()
	}()

	defer func() {
		server.Stop()
	}()

	connection, err := net.Dial("tcp", "localhost:9090")
	assert.Nil(t, err)

	buffer, _ := proto2.NewPutOrUpdateKeyValueMessage("DiskType", "NVMe SSD").Serialize()
	_, _ = connection.Write(buffer)

	buffer, _ = proto2.NewGetValueMessage("DiskType").Serialize()
	_, _ = connection.Write(buffer)

	connectionReader := conn.NewConnectionReader(connection)
	_, _ = connectionReader.AttemptReadOrErrorOut()

	message, err := connectionReader.AttemptReadOrErrorOut()

	assert.Nil(t, err)
	assert.Equal(t, "NVMe SSD", message.Value)
}

func TestSendsMultiplePutOrUpdateAndAGetOverAConnection(t *testing.T) {
	server, err := NewTCPServer("localhost", 8888)
	assert.Nil(t, err)

	go func() {
		server.Start()
	}()

	defer func() {
		server.Stop()
	}()

	connection, err := net.Dial("tcp", "localhost:8888")
	assert.Nil(t, err)

	sendMultiplePutOrUpdates := func() {
		buffer, _ := proto2.NewPutOrUpdateKeyValueMessage("DiskType", "NVMe SSD").Serialize()
		_, _ = connection.Write(buffer)

		buffer, _ = proto2.NewPutOrUpdateKeyValueMessage("Storage", "LSM").Serialize()
		_, _ = connection.Write(buffer)

		buffer, _ = proto2.NewPutOrUpdateKeyValueMessage("System", "Distributed").Serialize()
		_, _ = connection.Write(buffer)

		buffer, _ = proto2.NewGetValueMessage("System").Serialize()
		_, _ = connection.Write(buffer)
	}
	attemptLastRead := func() (*proto2.KeyValueMessage, error) {
		connectionReader := conn.NewConnectionReader(connection)
		_, _ = connectionReader.AttemptReadOrErrorOut()
		_, _ = connectionReader.AttemptReadOrErrorOut()
		_, _ = connectionReader.AttemptReadOrErrorOut()

		return connectionReader.AttemptReadOrErrorOut()
	}

	sendMultiplePutOrUpdates()
	message, err := attemptLastRead()

	assert.Nil(t, err)
	assert.Equal(t, "Distributed", message.Value)
}
