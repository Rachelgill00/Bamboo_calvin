package main

import (
	"github.com/gitferry/bamboo"
	"github.com/gitferry/bamboo/benchmark"
	"github.com/gitferry/bamboo/db"
)

// Database implements bamboo.DB interface for benchmarking
type Database struct {
	bamboo.Client
}

func (d *Database) Init() error {
	return nil
}

func (d *Database) Stop() error {
	return nil
}

func (d *Database) Write(k int, v []byte) error {
	key := db.Key(k)
	err := d.Put(key, v)
	return err
}

func main() {
	bamboo.Init()

	d := new(Database)							//创建一个database实例d
	d.Client = bamboo.NewHTTPClient()			//创建一个httpclient实例
	b := benchmark.NewBenchmark(d)				//使用d创建一个新的benchmark实例b
	b.Run()										//运行基准测试
}
