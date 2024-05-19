package conn

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"many-flavors-of-nwing-io/single_threaded_blocking_io/proto"
	"many-flavors-of-nwing-io/single_threaded_blocking_io/store"
	"net"
	"sync"
	"testing"
	"time"
)

func TestIncomingConnection(t *testing.T) {
	inMemoryStore := store.NewInMemoryStore()

	putOrUpdate := func() {
		var wg sync.WaitGroup
		wg.Add(2)

		source, incoming := net.Pipe()
		defer func() {
			_ = source.Close()
			_ = incoming.Close()
		}()

		incomingConnectionForPutOrUpdate := NewIncomingTCPConnection(incoming, inMemoryStore)

		go func() {
			defer wg.Done()
			incomingConnectionForPutOrUpdate.Handle()
		}()
		go func() {
			defer wg.Done()
			_ = source.SetReadDeadline(time.Now().Add(5 * time.Second))
			_, err := proto.DeserializeFrom(bufio.NewReader(source))

			assert.Nil(t, err)
		}()

		buffer, _ := proto.NewPutOrUpdateKeyValueMessage("DiskType", "NVMe SSD").Serialize()
		_, _ = source.Write(buffer)

		incomingConnectionForPutOrUpdate.Close()
		wg.Wait()
	}
	get := func() {
		var wg sync.WaitGroup
		wg.Add(2)

		source, incoming := net.Pipe()
		defer func() {
			_ = source.Close()
			_ = incoming.Close()
		}()

		incomingConnectionForGet := NewIncomingTCPConnection(incoming, inMemoryStore)

		go func() {
			defer wg.Done()
			incomingConnectionForGet.Handle()
		}()
		go func() {
			defer wg.Done()
			_ = source.SetReadDeadline(time.Now().Add(5 * time.Second))
			message, err := proto.DeserializeFrom(bufio.NewReader(source))

			assert.Nil(t, err)
			assert.Equal(t, "NVMe SSD", message.Value)
		}()

		buffer, _ := proto.NewGetValueMessage("DiskType").Serialize()
		_, _ = source.Write(buffer)

		incomingConnectionForGet.Close()
		wg.Wait()
	}
	putOrUpdate()
	get()
}
