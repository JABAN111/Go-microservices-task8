package managers

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"yadro.com/course/api/core"
)

type ClientManager struct {
	clients map[string]core.GrpcClient
	mu      sync.RWMutex

	log *slog.Logger
}

func NewClientManager(log *slog.Logger) *ClientManager {

	return &ClientManager{
		clients: make(map[string]core.GrpcClient),
		log:     log,
	}
}

func (cm *ClientManager) Register(name string, client core.GrpcClient) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.clients[name] = client
}

// GetClient deprecated
// По изначальной задумке было нужно, сейчас осталось как часть legacy
func (cm *ClientManager) GetClient(name string) (core.GrpcClient, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	client, exists := cm.clients[name]
	if !exists {
		return nil, errors.New("client not found: " + name)
	}

	return client, nil
}

// CloseAll нужен для gracefullshutdown
func (cm *ClientManager) CloseAll(ctx context.Context) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	var wg sync.WaitGroup
	for _, client := range cm.clients {
		wg.Add(1)
		go func(client core.GrpcClient) {
			defer wg.Done()

			done := make(chan error, 1)

			go func() {
				done <- client.Close()
			}()

			select {
			case <-ctx.Done():
				cm.log.Warn("Time-out of client disconnecting")
				return
			case err := <-done:
				if err != nil {
					cm.log.Error("Error while closing the client", "error", err)
				}
				return
			}
		}(client)
	}

	wg.Wait()
}

func (cm *ClientManager) PingAll(ctx context.Context) map[string]string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	status := make(map[string]string)

	for name, client := range cm.clients {
		if err := client.Ping(ctx); err != nil {
			status[name] = "unavailable"
		} else {
			status[name] = "ok"
		}
	}

	return status
}
