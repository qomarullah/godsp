package lib

import (
	"math/rand"
	"net"
	"sync"
)

// ConnectionPool is a thread safe list of net.Conn instances
type ConnectionPool struct {
	mutex   sync.RWMutex
	list    map[int]net.Conn
	session map[int]string
}

// NewConnectionPool is the factory method to create new connection pool
func NewConnectionPool() *ConnectionPool {
	pool := &ConnectionPool{
		list:    make(map[int]net.Conn),
		session: make(map[int]string),
	}

	return pool
}

// Add collection to pool
func (pool *ConnectionPool) Add(connection net.Conn, session string) int {
	pool.mutex.Lock()
	nextConnectionId := len(pool.list)
	pool.list[nextConnectionId] = connection
	pool.session[nextConnectionId] = session
	pool.mutex.Unlock()
	return nextConnectionId
}

// Add by id
func (pool *ConnectionPool) AddWithId(connection net.Conn, connectionId int, session string) int {
	pool.mutex.Lock()
	pool.list[connectionId] = connection
	pool.session[connectionId] = session
	pool.mutex.Unlock()
	return connectionId
}

// Get connection by id
func (pool *ConnectionPool) Get(connectionId int) net.Conn {
	pool.mutex.RLock()
	connection := pool.list[connectionId]
	pool.mutex.RUnlock()
	return connection
}

// Get connection by random
func (pool *ConnectionPool) GetWithId() (net.Conn, int, string) {
	//rand.Seed(time.Now().Unix())
	connectionId := rand.Intn(len(pool.list))
	pool.mutex.RLock()
	connection := pool.list[connectionId]
	session := pool.session[connectionId]
	pool.mutex.RUnlock()
	return connection, connectionId, session
}

// Remove connection from pool
func (pool *ConnectionPool) Remove(connectionId int) {
	pool.mutex.Lock()
	delete(pool.list, connectionId)
	delete(pool.session, connectionId)
	pool.mutex.Unlock()

}

// Size of connections pool
func (pool *ConnectionPool) Size() int {
	return len(pool.list)
}

// Range iterates over pool
func (pool *ConnectionPool) Range(callback func(net.Conn, int)) {
	pool.mutex.RLock()
	for connectionId, connection := range pool.list {
		callback(connection, connectionId)
	}
	pool.mutex.RUnlock()
}
