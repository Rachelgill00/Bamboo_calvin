package replica

import (
	"github.com/gitferry/bamboo/blockchain"
	"github.com/gitferry/bamboo/message"
	"github.com/gitferry/bamboo/pacemaker"
	"github.com/gitferry/bamboo/types"
)

type Safety interface {
	ProcessBlock(block *blockchain.Block) error
	//接受一个指向 blockchain.Block 结构体的指针作为参数；。返回的错误表示区块处理是否成功或遇到问题。

	ProcessVote(vote *blockchain.Vote)
	//处理与区块链共识机制相关的投票。

	ProcessRemoteTmo(tmo *pacemaker.TMO)
	//接收一个指向 pacemaker.TMO 结构体的指针。
	//TMO 通常向当前领导者或节拍器指示发送者的时间超时预期或其对网络进展的看法。
	//处理此类消息可能涉及调整本地定时参数、触发重传，或参与视图变更程序。

	ProcessLocalTmo(view types.View)
	//处理与节点参与共识协议相关的本地超时事件。此方法可能触发诸如提出新区块、广播投票或响应超时而过渡到新视图等操作。

	MakeProposal(view types.View, payload []*message.Transaction) *blockchain.Block
	//此方法为给定的 view 和交易负载构建一个新的区块提议。

	GetChainStatus() string
	//此方法检索并返回区块链的状态描述。
}
