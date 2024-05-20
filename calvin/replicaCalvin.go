package calvin

import (
	"encoding/gob"
	"fmt"
)

type ReplicaCalvin struct {
	node.Node
	start           chan bool // signal to start the node
	isStarted       atomic.Bool
	timer           *time.Timer // timeout for each view
	eventChan       chan interface{}

	/* for monitoring node statistics 用于监控节点统计信息*/
	thrus                string
}

// NewReplica creates a new ReplicaCalvin instance （创建一个新的副本实例）
func NewReplicaCalvin(id identity.NodeID, alg string, isByz bool) *ReplicaCalvin {
	r := new(ReplicaCalvin)
	r.Node = node.NewNode(id, isByz)
	if isByz {
		log.Infof("[%v] is Byzantine", r.ID())
	}
	if config.GetConfig().Master == "0" {
		r.Election = election.NewRotation(config.GetConfig().N())
	} else {
		r.Election = election.NewStatic(config.GetConfig().Master)
	}

	gob.Register(pacemaker.TMO{})

	// Is there a better way to reduce the number of parameters?
	switch alg {
	case "hotstuff":
		r.Safety = hotstuff.NewHotStuff(r.Node, r.pm, r.Election, r.committedBlocks, r.forkedBlocks)
	case "tchs":
		r.Safety = tchs.NewTchs(r.Node, r.pm, r.Election, r.committedBlocks, r.forkedBlocks)
	case "streamlet":
		r.Safety = streamlet.NewStreamlet(r.Node, r.pm, r.Election, r.committedBlocks, r.forkedBlocks)
	case "lbft":
		r.Safety = lbft.NewLbft(r.Node, r.pm, r.Election, r.committedBlocks, r.forkedBlocks)
	case "fasthotstuff":
		r.Safety = fhs.NewFhs(r.Node, r.pm, r.Election, r.committedBlocks, r.forkedBlocks)
	default:
		r.Safety = hotstuff.NewHotStuff(r.Node, r.pm, r.Election, r.committedBlocks, r.forkedBlocks)
	}
	return r
}

func (r *ReplicaCalvin) Start() {
	go r.Run()
	// wait for the start signal
	<-r.start
	for r.isStarted.Load() {
		event := <-r.eventChan
		switch v := event.(type) {
		case types.View:
			r.processNewView(v)
		case blockchain.Block:
			startProcessTime := time.Now()
			r.totalProposeDuration += startProcessTime.Sub(v.Timestamp)
			_ = r.Safety.ProcessBlock(&v)
			r.totalProcessDuration += time.Now().Sub(startProcessTime)
			r.voteStart = time.Now()
			r.processedNo++
		case blockchain.Vote:
			startProcessTime := time.Now()
			r.Safety.ProcessVote(&v)
			processingDuration := time.Now().Sub(startProcessTime)
			r.totalVoteTime += processingDuration
			r.voteNo++
		case pacemaker.TMO:
			r.Safety.ProcessRemoteTmo(&v)
		}
	}
}

