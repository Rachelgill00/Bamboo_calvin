package bcb

import (
	"sync"
	"time"

	"github.com/gitferry/zeitgeber/log"

	"github.com/gitferry/zeitgeber"
)

type bcb struct {
	zeitgeber.Node
	zeitgeber.Election

	curView      zeitgeber.View
	quorum       *zeitgeber.Quorum
	highCert     *zeitgeber.TC
	lastViewTime time.Time
	viewDuration map[zeitgeber.View]time.Duration

	newViewChan chan zeitgeber.View

	mu sync.Mutex
}

func NewBcb(n zeitgeber.Node, election zeitgeber.Election) zeitgeber.Pacemaker {
	bcb := new(bcb)
	bcb.Node = n
	bcb.Election = election
	bcb.newViewChan = make(chan zeitgeber.View)
	bcb.quorum = zeitgeber.NewQuorum()
	bcb.Register(zeitgeber.TCMsg{}, bcb.HandleTC)
	bcb.Register(zeitgeber.TmoMsg{}, bcb.HandleTmo)
	bcb.highCert = zeitgeber.NewTC(0)
	bcb.lastViewTime = time.Now()
	bcb.viewDuration = make(map[zeitgeber.View]time.Duration)
	return bcb
}

func (b *bcb) HandleTmo(tmo zeitgeber.TmoMsg) {
	b.mu.Lock()
	if tmo.View < b.curView {
		//log.Warningf("[%v] received timeout msg with view %v lower than the current view %v", b.NodeID(), tmo.View, b.curView)
		b.mu.Unlock()
		return
	}
	b.quorum.ACK(tmo.View, tmo.NodeID)
	if b.quorum.SuperMajority(tmo.View) {
		//log.Infof("[%v] a time certificate for view %v is generated", b.NodeID(), tmo.View)
		b.Send(b.FindLeaderFor(tmo.View), zeitgeber.TCMsg{View: tmo.View})
		b.mu.Unlock()
		b.NewView(tmo.View)
		return
	}
	if tmo.HighTC.View >= b.curView {
		b.mu.Unlock()
		b.NewView(tmo.HighTC.View)
		return
	}
	b.mu.Unlock()
}

func (b *bcb) HandleTC(tc zeitgeber.TCMsg) {
	//log.Infof("[%v] is processing tc for view %v", b.NodeID(), tc.View)
	b.mu.Lock()
	if tc.View < b.curView {
		//log.Warningf("[%s] received tc's view %v is lower than current view %v", b.NodeID(), tc.View, b.curView)
		b.mu.Unlock()
		return
	}
	if tc.View > b.highCert.View {
		b.highCert = zeitgeber.NewTC(tc.View)
	}
	b.mu.Unlock()
	b.NewView(tc.View)
}

// TimeoutFor broadcasts the timeout msg for the view when it timeouts
func (b *bcb) TimeoutFor(view zeitgeber.View) {
	tmoMsg := zeitgeber.TmoMsg{
		View:   view,
		NodeID: b.ID(),
		HighTC: zeitgeber.NewTC(view - 1),
	}
	//log.Debugf("[%s] is timeout for view %v", b.NodeID(), view)
	if b.IsByz() {
		b.MulticastQuorum(zeitgeber.GetConfig().ByzNo, tmoMsg)
		return
	}
	b.Broadcast(tmoMsg)
	b.HandleTmo(tmoMsg)
}

func (b *bcb) NewView(view zeitgeber.View) {
	b.mu.Lock()
	if view < b.curView {
		//log.Warningf("the view %v is lower than current view %v", view, b.curView)
		b.mu.Unlock()
		return
	}
	b.viewDuration[b.curView] = time.Now().Sub(b.lastViewTime)
	b.curView = view + 1
	b.lastViewTime = time.Now()
	b.mu.Unlock()
	if view == 100 {
		b.printViewTime()
	}
	b.newViewChan <- view + 1 // reset timer for the next view
}

func (b *bcb) printViewTime() {
	//log.Infof("[%v] is printing view duration", b.NodeID())
	for view, duration := range b.viewDuration {
		log.Infof("view %v duration: %v seconds", view, duration.Seconds())
	}
}

func (b *bcb) EnteringViewEvent() chan zeitgeber.View {
	return b.newViewChan
}

func (b *bcb) GetCurView() zeitgeber.View {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.curView
}

func (b *bcb) GetHighCert() *zeitgeber.TC {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.highCert
}
