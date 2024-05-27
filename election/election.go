package election

import (
	"github.com/Rachelgill00/simulation-for-calvin/identity"
	"github.com/Rachelgill00/simulation-for-calvin/types"
)

type Election interface {
	IsLeader(id identity.NodeID, view types.View) bool
	FindLeaderFor(view types.View) identity.NodeID
}
