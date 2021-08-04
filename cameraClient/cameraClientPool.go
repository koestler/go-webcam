package cameraClient

import "sync"

type ClientPool struct {
	clients      map[string]*Client
	clientsMutex sync.RWMutex
}

func RunPool() (pool *ClientPool) {
	pool = &ClientPool{
		clients: make(map[string]*Client),
	}
	return
}

func (p *ClientPool) Shutdown() {
	p.clientsMutex.RLock()
	defer p.clientsMutex.RUnlock()
	for _, c := range p.clients {
		c.Shutdown()
	}
}

func (p *ClientPool) AddClient(client *Client) {
	p.clientsMutex.Lock()
	defer p.clientsMutex.Unlock()
	p.clients[client.Name()] = client
}

func (p *ClientPool) RemoveClient(client *Client) {
	p.clientsMutex.Lock()
	defer p.clientsMutex.Unlock()
	delete(p.clients, client.Name())
}

func (p *ClientPool) GetClient(clientName string) *Client {
	p.clientsMutex.RLock()
	defer p.clientsMutex.RUnlock()
	return p.clients[clientName]
}
