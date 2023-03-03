package promote

import (
	"sync"
)

var allowed = [9]string{"invoice", "payment", "application", "historicalUsed", "historicalOrder", "historicalReceivable", "historicalSettle", "poolPlan", "poolUsed"}

type Pools struct {
	slowPool *pendingPool
	fastPool *encryptedPool
}

type encryptedPool struct {
	encryptedMessage      map[string][]interface{}
	encryptedMessageMutex sync.Mutex
}

type pendingPool struct {
	pendingMessage   map[string][]interface{}
	pendingPoolMutex sync.Mutex
}

type packedMessage struct {
	header        string
	cipher        []byte
	encryptionKey []byte
	signed        []byte
}

func NewPools() *Pools {
	fast := newEncryptedPool()
	slow := newPendingPool()
	return &Pools{
		slowPool: slow,
		fastPool: fast,
	}
}
func (p *Pools) Insert(packed packedMessage, name string, poolType string) {
	if !verify(name) {
		panic("指定方法名称错误")
	}
	if poolType == "fast" {
		p.fastPool.insertMessage(packed, name)

	} else if poolType == "slow" {
		p.slowPool.insertMessage(packed, name)

	} else {
		panic("池子类型错误")
	}
}
func (p *Pools) Delete(name string, poolType string) {
	if !verify(name) {
		panic("指定方法名称错误")
	}
	if poolType == "fast" {
		p.fastPool.deleteMessage(name)

	} else if poolType == "slow" {
		p.slowPool.deleteMessage(name)

	} else {
		panic("池子类型错误")
	}
}
func (p *Pools) GetPoolLength(name string, poolType string) int {
	if poolType == "fast" {
		length := p.fastPool.getLength(name)
		return length

	} else if poolType == "slow" {
		length := p.slowPool.getLength(name)
		return length

	} else {
		panic("池子类型错误")
	}
}
func (p *Pools) QueryMessages(name string, poolType string) []interface{} {
	if poolType == "fast" {
		res := p.fastPool.queryMessage(name)
		return res
	} else if poolType == "slow" {
		res := p.slowPool.queryMessage(name)
		return res
	} else {
		panic("池子类型错误")
	}

}
func verify(name string) bool {
	if name == "" {
		return false
	}
	for _, str := range allowed {
		if name == str {
			return true
		}
	}
	return false
}
func newEncryptedPool() *encryptedPool {
	encrypted := make(map[string][]interface{})
	return &encryptedPool{
		encryptedMessage: encrypted,
	}
}

func newPendingPool() *pendingPool {
	pending := make(map[string][]interface{})
	return &pendingPool{
		pendingMessage: pending,
	}
}
func (e *encryptedPool) insertMessage(packed packedMessage, name string) {
	e.encryptedMessageMutex.Lock()
	e.encryptedMessage[name] = append(e.encryptedMessage[name], packed)
	e.encryptedMessageMutex.Unlock()
}
func (e *encryptedPool) deleteMessage(name string) {
	e.encryptedMessageMutex.Lock()
	e.encryptedMessage[name] = nil
	e.encryptedMessageMutex.Unlock()
}
func (e *encryptedPool) getLength(name string) int {
	e.encryptedMessageMutex.Lock()
	length := len(e.encryptedMessage[name])
	e.encryptedMessageMutex.Unlock()
	return length
}
func (e *encryptedPool) queryMessage(name string) []interface{} {
	e.encryptedMessageMutex.Lock()
	temp := e.encryptedMessage[name]
	e.encryptedMessage[name] = nil
	e.encryptedMessageMutex.Unlock()
	return temp
}
func (p *pendingPool) insertMessage(packed packedMessage, name string) {
	p.pendingPoolMutex.Lock()
	p.pendingMessage[name] = append(p.pendingMessage[name], packed)
	p.pendingPoolMutex.Unlock()
}
func (p *pendingPool) deleteMessage(name string) {
	p.pendingPoolMutex.Lock()
	p.pendingMessage[name] = nil
	p.pendingPoolMutex.Unlock()
}
func (p *pendingPool) getLength(name string) int {
	p.pendingPoolMutex.Lock()
	length := len(p.pendingMessage[name])
	p.pendingPoolMutex.Unlock()
	return length
}
func (p *pendingPool) queryMessage(name string) []interface{} {
	p.pendingPoolMutex.Lock()
	temp := p.pendingMessage[name]
	p.pendingMessage[name] = nil
	p.pendingPoolMutex.Unlock()
	return temp
}
