package main

import (
	"sync/atomic"

	"github.com/idealeak/goserver.v3/core/netlib"
)

var (
	DupFilterName = "session-filter-dup"
)

type DupFilter struct {
}

func (df *DupFilter) GetName() string {
	return DupFilterName
}

func (df *DupFilter) GetInterestOps() uint {
	return 1 << netlib.InterestOps_Received
}

func (df *DupFilter) OnSessionOpened(s *netlib.Session) bool {
	return true
}

func (df *DupFilter) OnSessionClosed(s *netlib.Session) bool {
	return true
}

func (df *DupFilter) OnSessionIdle(s *netlib.Session) bool {
	return true
}

func (df *DupFilter) OnPacketReceived(s *netlib.Session, packetid int, logicNo uint32, packet interface{}) bool {
	if s.GroupId != 0 {
		bs := BundleMgrSington.GetBundleSession(uint16(s.GroupId))
		if bs != nil {
			if atomic.CompareAndSwapUint32(&bs.rcvLogicNo, logicNo-1, logicNo) {
				return true
			}
		}
		return false
	}
	return true
}

func (df *DupFilter) OnPacketSent(s *netlib.Session, data []byte) bool {
	return true
}

func init() {
	netlib.RegisteSessionFilterCreator(DupFilterName, func() netlib.SessionFilter {
		return &DupFilter{}
	})
}
