package embeddednats

import (
	"fmt"
	"time"

	natsserver "github.com/nats-io/nats-server/v2/server"
	natsgo "github.com/nats-io/nats.go"
)

const readyTimeout = 5 * time.Second

const IssueCreatedSubject = "issues.created"

// StartEmbeddedServer starts an in-process NATS server on a random available port
// and returns it once it is ready to accept connections.
// The caller is responsible for calling server.Shutdown() on shutdown.
func StartEmbeddedServer() (*natsserver.Server, error) {
	options := &natsserver.Options{
		Host:   "127.0.0.1",
		Port:   natsserver.RANDOM_PORT,
		NoLog:  true,
		NoSigs: true,
	}
	server, err := natsserver.NewServer(options)
	if err != nil {
		return nil, fmt.Errorf("new embedded nats server: %w", err)
	}
	server.Start()
	if !server.ReadyForConnections(readyTimeout) {
		return nil, fmt.Errorf("embedded nats server did not become ready within %s", readyTimeout)
	}
	return server, nil
}

// Connect returns a NATS client connection to the given embedded server.
// The caller is responsible for calling connection.Close() on shutdown.
func Connect(server *natsserver.Server) (*natsgo.Conn, error) {
	connection, err := natsgo.Connect(server.ClientURL())
	if err != nil {
		return nil, fmt.Errorf("connect to embedded nats: %w", err)
	}
	return connection, nil
}
