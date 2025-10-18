package main

import (
	"io"
	"log"
	"net"
	"net/netip"
	"os"
	"time"

	"github.com/hashicorp/raft"
)

type Config struct {
	Raft         *raft.Config
	TCPTransport *TCPTransportConfig
}

type TCPTransportConfig struct {
	Address   *net.TCPAddr
	MaxPool   int
	Timeout   time.Duration
	LogOutput io.Writer
}

func NewConfig(serverID raft.ServerID) (*Config, error) {
	raftAddrPort, err := netip.ParseAddrPort(raftAddr)
	if err != nil {
		log.Printf("error parse raft addr: %v", err)
		return nil, err
	}
	return &Config{
		Raft: &raft.Config{
			ProtocolVersion:    raft.ProtocolVersionMax,
			HeartbeatTimeout:   1000 * time.Millisecond,
			ElectionTimeout:    1000 * time.Millisecond,
			CommitTimeout:      50 * time.Millisecond,
			MaxAppendEntries:   64,
			ShutdownOnRemove:   true,
			TrailingLogs:       10240,
			SnapshotInterval:   120 * time.Second,
			SnapshotThreshold:  8192,
			LeaderLeaseTimeout: 500 * time.Millisecond,
			LogLevel:           "DEBUG",
			LocalID:            serverID,
		},
		TCPTransport: &TCPTransportConfig{
			Address:   net.TCPAddrFromAddrPort(raftAddrPort),
			MaxPool:   3,
			Timeout:   10 * time.Second,
			LogOutput: os.Stderr,
		},
	}, nil
}
