package behaviago

/**
  DecoratorCountLimit can be set a integer Count limit value. DecoratorCountLimit node tick its child until
  inner count less equal than count limit value. Whether node increase inner count depend on
  the return value of its child when it updates. if DecorateChildEnds flag is true, node increase count
  only when its child node return value is Success or Failure. The inner count will never reset until
  attachment on the node evaluate true.
*/
const (
	DecoratorCountLimitNodeName = "DecoratorCountLimit"
)

func init() {
	RegisteNodeCreator(DecoratorCountLimitNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		n := &DecoratorCountLimit{}
		return n
	}))
}

func newDecoratorCountLimitNode() *DecoratorCountLimit {
	n := &DecoratorCountLimit{}
	n.SetClassName(DecoratorCountLimitNodeName)
	n.SetSelf(n)
	return n
}

// ============================================================================
type DecoratorCountLimit struct {
	DecoratorCount
}

func (d *DecoratorCountLimit) CheckIfReInit(a *Agent) bool {
	return d.EvaluteCustomCondition(a)
}

func (d *DecoratorCountLimit) CreateTask() BehaviorTask {
	BTGLog.Trace("(d *DecoratorCountLimit) CreateTask()")
	return &DecoratorCountLimitTask{}
}

func (n *DecoratorCountLimit) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*DecoratorCountLimit); !ok {
		return false
	}
	return true
}

// ============================================================================

type DecoratorCountLimitTask struct {
	DecoratorCountTask
	inited bool
}

func (dt *DecoratorCountLimitTask) OnEnter(a *Agent) bool {
	if node, ok := dt.GetNode().(*DecoratorCountLimit); ok {
		if node.CheckIfReInit(a) {
			dt.inited = false
		}
		if !dt.inited {
			dt.inited = true
			dt.count = dt.GetCount(a)
		}
		//if count is -1, it is endless
		if dt.count > 0 {
			dt.count--
			return true
		} else if dt.count == 0 {
			return false
		} else if dt.count == -1 {
			return true
		}
	}
	return false
}

func (dt *DecoratorCountLimitTask) Decorate(status EBTStatus) EBTStatus {
	return status
}
