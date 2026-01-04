package tgo

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"
)

// Transport handles communication with the TGO host via Unix Socket or TCP.
type Transport struct {
	network string
	address string
	conn    net.Conn
	mu      sync.Mutex
}

func NewUnixTransport(path string) *Transport {
	return &Transport{network: "unix", address: path}
}

func NewTCPTransport(addr string) *Transport {
	return &Transport{network: "tcp", address: addr}
}

// Connect establishes a connection to the TGO host.
func (t *Transport) Connect() error {
	conn, err := net.Dial(t.network, t.address)
	if err != nil {
		return fmt.Errorf("failed to connect to TGO (%s) %s: %w", t.network, t.address, err)
	}
	t.conn = conn
	return nil
}


// Close closes the connection.
func (t *Transport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.conn != nil {
		err := t.conn.Close()
		t.conn = nil
		return err
	}
	return nil
}

// SendMessage sends a JSON-RPC message with a 4-byte big-endian length prefix.
func (t *Transport) SendMessage(msg any) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.conn == nil {
		return fmt.Errorf("not connected")
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Write 4-byte length prefix
	length := uint32(len(data))
	if err := binary.Write(t.conn, binary.BigEndian, length); err != nil {
		return fmt.Errorf("failed to write length prefix: %w", err)
	}

	// Write JSON data
	if _, err := t.conn.Write(data); err != nil {
		return fmt.Errorf("failed to write message data: %w", err)
	}

	return nil
}

// RecvMessage receives a JSON-RPC message with a 4-byte big-endian length prefix.
func (t *Transport) RecvMessage() (map[string]any, error) {
	if t.conn == nil {
		return nil, fmt.Errorf("not connected")
	}

	// Read 4-byte length prefix
	var length uint32
	if err := binary.Read(t.conn, binary.BigEndian, &length); err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, fmt.Errorf("failed to read length prefix: %w", err)
	}

	// Read JSON data
	data := make([]byte, length)
	if _, err := io.ReadFull(t.conn, data); err != nil {
		return nil, fmt.Errorf("failed to read message data: %w", err)
	}

	var msg map[string]any
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return msg, nil
}

