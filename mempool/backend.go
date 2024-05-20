package mempool

import (
	"container/list"
	"sync"

	"github.com/gitferry/bamboo/message"
)

// Backend 结构体定义了一个内存池后端实现
type Backend struct {
	txns          *list.List  //txns（用于存放交易）是一个指向list.List（标准库中一个双向链表数据结构）类型的指针。
	limit         int         //limit（用于限制交易池中交易的数量）是一个int类型的变量。
	totalReceived int64       //	totalReceived（用于记录接收到的交易数量）是一个int64类型的变量。
	*BloomFilter              //一种概率型数据结构，用于快速判断一个元素是否可能属于某个集合
	mu            *sync.Mutex //mu 是一个指向 sync.Mutex 类型的指针。sync.Mutex 是Go语言中实现互斥锁（mutex）的结构体
}

func NewBackend(limit int) *Backend {
	var mu sync.Mutex
	return &Backend{
		txns:        list.New(),
		BloomFilter: NewBloomFilter(),
		mu:          &mu,
		limit:       limit,
	}
}

func (b *Backend) insertBack(txn *message.Transaction) {
	if txn == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.size() > b.limit {
		return
	}
	b.totalReceived++
	b.txns.PushBack(txn)
}

func (b *Backend) insertFront(txn *message.Transaction) {
	if txn == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.txns.PushFront(txn)
}

func (b *Backend) size() int {
	return b.txns.Len()
}

func (b *Backend) front() *message.Transaction {
	if b.size() == 0 {
		return nil
	}
	ele := b.txns.Front()
	val, ok := ele.Value.(*message.Transaction)
	if !ok {
		return nil
	}
	b.txns.Remove(ele)
	return val
}

func (b *Backend) some(n int) []*message.Transaction {
	// n means the number of TXN to be acquired
	var batchSize int
	b.mu.Lock()
	defer b.mu.Unlock()
	batchSize = b.size() //to get the size of the storage TXN
	if batchSize >= n {
		batchSize = n
	}
	//make a slice to storage the acquired txn
	batch := make([]*message.Transaction, 0, batchSize)
	//to acquire the txn in the backend and append it to the batch
	for i := 0; i < batchSize; i++ {
		tx := b.front()
		batch = append(batch, tx)
	}
	return batch
}

func (b *Backend) addToBloom(id string) {
	b.Add(id)
}
