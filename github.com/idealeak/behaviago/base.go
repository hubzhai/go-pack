package behaviago

var bTGStartCnt int

func IsStarted() bool {
	return bTGStartCnt > 0
}

func GetVerStr() string {
	return "1.0"
}

func BaseStartup() error {
	bTGStartCnt++
	if bTGStartCnt == 1 {
		//todo:start
	}
	return nil
}

func BaseCleanup() {
	bTGStartCnt--
	if bTGStartCnt == 0 {
		//todo:stop
	}
}

func TryStartup() error {
	if !IsStarted() {
		return BaseStartup()
	}
	return nil
}
