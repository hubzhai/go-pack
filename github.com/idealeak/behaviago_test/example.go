package main

import "github.com/idealeak/behaviago"

type CBTPlayer struct {
	*behaviago.Agent
	m_iBaseSpeed uint
	m_Frames     int
}

func (this *CBTPlayer) Condition() bool {
	behaviago.BTGLog.Trace("(this *CBTPlayer) Condition() enter")
	defer behaviago.BTGLog.Trace("(this *CBTPlayer) Condition() leave")
	this.m_Frames = 0
	return true
}

func (this *CBTPlayer) Action1() behaviago.EBTStatus {
	behaviago.BTGLog.Trace("(this *CBTPlayer) Action1() enter")
	defer behaviago.BTGLog.Trace("(this *CBTPlayer) Action1() leave")
	return behaviago.BT_SUCCESS
}

func (this *CBTPlayer) Action3() behaviago.EBTStatus {
	behaviago.BTGLog.Trace("(this *CBTPlayer) Action3() enter")
	defer behaviago.BTGLog.Trace("(this *CBTPlayer) Action3() leave")
	this.m_Frames++

	if this.m_Frames == 3 {
		return behaviago.BT_SUCCESS
	}

	return behaviago.BT_RUNNING
}

func main() {
	behaviago.BTGWorkspace.SetFilePath("behaviac/exported")
	player := &CBTPlayer{}
	player.Agent = behaviago.NewAgent(1, 0, "CBTPlayer")
	player.SetClient(player)
	player.BTGLoad("demo_running", true)
	player.BTGSetCurrent("demo_running", behaviago.TM_Transfer, false)

	frames := 0
	status := behaviago.BT_RUNNING
	for status == behaviago.BT_RUNNING {
		frames++
		behaviago.BTGLog.Infof("frame %v", frames)
		status = player.BTGExec()
	}
}
