package pacemaker

import (
	"github.com/Rachelgill00/simulation-for-calvin/blockchain"
	"github.com/Rachelgill00/simulation-for-calvin/crypto"
	"github.com/Rachelgill00/simulation-for-calvin/identity"
	"github.com/Rachelgill00/simulation-for-calvin/types"
)

type TMO struct {
	View   types.View
	NodeID identity.NodeID
	HighQC *blockchain.QC
}

type TC struct {
	types.View
	crypto.AggSig
	crypto.Signature
}

func NewTC(view types.View, requesters map[identity.NodeID]*TMO) *TC {
	// TODO: add crypto
	return &TC{View: view}
}
