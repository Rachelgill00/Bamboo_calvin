package election

import (
	"github.com/Rachelgill00/simulation-for-calvin/identity"
	"github.com/Rachelgill00/simulation-for-calvin/types"
)

type Static struct {
	master identity.NodeID
}

func NewStatic(master identity.NodeID) *Static {
	return &Static{
		master: master,
	}
}

func (st *Static) IsLeader(id identity.NodeID, view types.View) bool {
	return id == st.master
}

func (st *Static) FindLeaderFor(view types.View) identity.NodeID {
	return st.master
}
