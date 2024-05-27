package benchmark

import (
	"math/rand"
	"sync"
	"time"

	"github.com/Rachelgill00/simulation-for-calvin/config"
	"github.com/Rachelgill00/simulation-for-calvin/log"
)

var count uint64

// DB is general interface implemented by client to call client library
// DB 是客户端实现的通用接口，用于调用客户端库
// 定义事务参与者接口
type DB interface {
	Init() error
	Write(key int, value []byte) error
	Stop() error
}

// DefaultBConfig returns a default benchmark config
// 返回默认基准配置
func DefaultBConfig() config.Bconfig {
	return config.Bconfig{
		T:           60, // total number of running time in seconds
		N:           0,  //total number of requests
		Throttle:    0,  // requests per second throttle, unused if 0
		Concurrency: 1,  // number of simulated clients
	}
}

// Benchmark is benchmarking tool that generates workload and collects operation history and latency
type Benchmark struct {
	db DB // read/write operation interface
	config.Bconfig
	*History

	rate      *Limiter
	latency   []time.Duration // latency per operation
	startTime time.Time
	counter   int

	wait sync.WaitGroup // waiting for all generated keys to complete
}

// NewBenchmark returns new Benchmark object given implementation of DB interface
// NewBenchmark 返回给定 DB 接口实现的新基准对象
func NewBenchmark(db DB) *Benchmark {
	b := new(Benchmark)
	b.db = db
	b.Bconfig = config.Configuration.Benchmark
	b.History = NewHistory()
	if b.Throttle > 0 {
		b.rate = NewLimiter(b.Throttle)
	}
	return b
}

// 执行事务
// Run starts the main logic of benchmarking
func (b *Benchmark) Run() {

	var genCount, sendCount, confirmCount uint64

	b.latency = make([]time.Duration, 0)
	keys := make(chan int, b.Concurrency)       //transfer keys
	latencies := make(chan time.Duration, 1000) //transfer task latencies
	defer close(latencies)
	//start a goroutine collect to collect the task latencies of every worker thread
	go b.collect(latencies)

	//Concurrency  int    // number of simulated clients
	for i := 0; i < b.Concurrency; i++ {
		//make thread
		//receive the keys from chennel keys
		go b.worker(keys, latencies)
	}
	//init the database
	b.db.Init()
	//启动一个计时器，超过指定时间后，停止基准测试
	//
	b.startTime = time.Now()
	if b.T > 0 {
		timer := time.NewTimer(time.Second * time.Duration(b.T))
	loop:
		for {
			select {
			case <-timer.C:
				log.Infof("Benchmark stops")
				break loop
			default:
				b.wait.Add(1)
				//log.Debugf("is generating key No.%v", j)
				//通过发送任务键来触发工作线程执行写入操作
				//并执行生成、送法和确认的技术
				k := b.next()
				genCount++
				keys <- k
				sendCount++
				//log.Debugf("generated key No.%v", j-1)
			}
		}
	} else {
		for i := 0; i < b.N; i++ {
			b.wait.Add(1)
			keys <- b.next()
		}
		b.wait.Wait()
	}

	t := time.Now().Sub(b.startTime)
	//执行
	b.db.Stop()
	close(keys)
	//当基准测试结束后，停止数据库操作，并统计基准测试的结果
	//包括并发度，持续时间、吞吐量和任务技术
	stat := Statistic(b.latency)
	confirmCount = uint64(len(b.latency))
	log.Infof("Concurrency = %d", b.Concurrency)
	log.Infof("Benchmark Time = %v\n", t)
	log.Infof("Throughput = %f\n", float64(len(b.latency))/t.Seconds())
	log.Infof("genCount: %d, sendCount: %d, confirmCount: %d", genCount, sendCount, confirmCount)
	log.Info(stat)

	//stat.WriteFile("latency")
	//b.History.WriteFile("history")
}

<<<<<<< HEAD
// 한 txn transfer하는 게
=======
// 한 txn 전성하는 게
// 从输入通道接受键，并将操作耗时发送到结果通道
>>>>>>> 88f168de99ce2d740298a494cb48afea10dcfd10
func (b *Benchmark) worker(keys <-chan int, result chan<- time.Duration) {
	//var s time.Time
	//var e time.Time
	//var v int
	//var err error

	for k := range keys {
		// tx := []byte{1, 2, 3}

		//op := new(operation)
		//v = rand.Int()
		//s = time.Now()

		//创建一个长度等于全局配置中指定的 PayloadSize 字节切片 value。
		value := make([]byte, config.GetConfig().PayloadSize)
		//使用 rand.Read 函数填充 value 切片，生成随机数据
		rand.Read(value)

		//rand.Read(value)
		//使用 Benchmark 结构体实例 b 的 db 字段（假设是一个数据库接口）的 Write 方法，
		//将整数键 k 与生成的随机数据 value 写入数据库。
		//执行
		_ = b.db.Write(k, value)
		//res, err := strconv.Atoi(r)
		//log.Debugf("latency is %v", time.Duration(res)*time.Nanosecond)
		//e = time.Now()
		//op.input = v
		//op.start = s.Sub(b.startTime).Nanoseconds()
		//if err == nil {
		//op.end = e.Sub(b.startTime).Nanoseconds()
		//result <- e.Sub(s)
		//result <- time.Duration(res) * time.Nanosecond
		//} else {
		//op.end = math.MaxInt64
		//log.Error(err)
		//}
		//b.History.AddOperation(k, op)
	}
}

// generates key based on distribution
func (b *Benchmark) next() int {
	var key int
	switch b.Distribution {
	case "uniform":
		key = int(count)
		count += uint64(config.GetConfig().N() - config.GetConfig().ByzNo)
	default:
		log.Fatalf("unknown distribution %s", b.Distribution)
	}

	if b.Throttle > 0 {
		b.rate.Wait()
	}

	return key
}

func (b *Benchmark) collect(latencies <-chan time.Duration) {
	for t := range latencies {
		b.latency = append(b.latency, t)
		b.wait.Done()
	}
}
