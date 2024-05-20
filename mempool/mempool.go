package mempool

import (
	"github.com/gitferry/bamboo/config"
	"github.com/gitferry/bamboo/message"
	"time"
)

//简单结构体，包含一个指向backend的指针
//意味着Mempool本质上是对Backend对象的封装，继承了Backend的所有方法和属性
type MemPool struct {
	*Backend
}

// NewTransactions creates a new memory pool for transactions.
//Mempool的构造函数，用于创建一个Mempool实例。
func NewMemPool() *MemPool {
	mp := &MemPool{
		Backend: NewBackend(config.GetConfig().MemSize),
	}

	return mp
}

//这个方法接受一个指向message.Transaction的指针tx作为参数，表示待添加的新交易。
func (mp *MemPool) addNew(tx *message.Transaction) {
	//为交易添加一个时间戳。
	tx.Timestamp = time.Now()
	//将交易添加到Backend的后部。
	mp.Backend.insertBack(tx)
	//将交易的ID添加到布隆过滤器（Bloom filter）中，这通常用于快速查询是否存在某个交易，以提高效率。
	mp.Backend.addToBloom(tx.ID)
}
//处理的是“旧”交易。
func (mp *MemPool) addOld(tx *message.Transaction) {
	mp.Backend.insertFront(tx)
}
