package main

import (
	"time"

	"github.com/idealeak/behaviago"
)

type CBTPlayer struct {
	*behaviago.Agent
	m_iX         int
	m_iY         int
	m_iBaseSpeed uint
	m_Frames     int
}

func (this *CBTPlayer) Condition() bool {
	behaviago.BTGLog.Trace("============(this *CBTPlayer) Condition() enter")
	defer behaviago.BTGLog.Trace("============(this *CBTPlayer) Condition() leave")
	this.m_Frames = 0
	return true
}

func (this *CBTPlayer) Action1() behaviago.EBTStatus {
	behaviago.BTGLog.Trace("============(this *CBTPlayer) Action1() enter")
	defer behaviago.BTGLog.Trace("============(this *CBTPlayer) Action1() leave")
	return behaviago.BT_SUCCESS
}

func (this *CBTPlayer) Action3() behaviago.EBTStatus {
	behaviago.BTGLog.Trace("============(this *CBTPlayer) Action3() enter")
	defer behaviago.BTGLog.Trace("============(this *CBTPlayer) Action3() leave")
	this.m_Frames++

	if this.m_Frames == 3 {
		return behaviago.BT_SUCCESS
	}

	return behaviago.BT_RUNNING
}

func (this *CBTPlayer) MoveAhead(speed int) {
	this.m_iX += (int(this.m_iBaseSpeed) + speed)
	this.SetVariable("CurStep", this.m_iX)
	val := this.GetVariableVal("CurStep")
	behaviago.BTGLog.Tracef("============after MoveAhead [%v] Name:[%v] MoveAhead CurStep:[%v] CurSpeed[%v]", time.Now(), this.GetName(), val, speed)
}

func (this *CBTPlayer) MoveBack(speed int) {
	this.m_iX -= (int(this.m_iBaseSpeed) + speed)
	this.SetVariable("CurStep", this.m_iX)
	val := this.GetVariableVal("CurStep")
	behaviago.BTGLog.Tracef("============after MoveBack [%v] Name:[%v] MoveAhead CurStep:[%v] CurSpeed[%v]", time.Now(), this.GetName(), val, speed)
}

func main() {
	behaviago.BTGWorkspace.SetFilePath("behaviac/exported")
	player := &CBTPlayer{}
	player.Agent = behaviago.NewAgent(1, 0, "CBTPlayer")
	player.SetClient(player)
	player.BTGLoad("player", true)
	player.BTGSetCurrent("player", behaviago.TM_Transfer, false)

	frames := 0
	for frames < 1000 {
		frames++
		behaviago.BTGLog.Infof("==========frame %v", frames)
		behaviago.BTGWorkspace.SetTimeSinceStartup(behaviago.BTGWorkspace.GetTimeSinceStartup() + int64(time.Millisecond*200))
		behaviago.BTGWorkspace.Update()
	}
}
