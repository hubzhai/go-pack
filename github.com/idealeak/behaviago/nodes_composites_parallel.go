package behaviago

/// Enumerates the options for when a parallel node is considered to have failed.
/**
  - FAIL_ON_ONE indicates that the node will return failure as soon as one of its children fails.
  - FAIL_ON_ALL indicates that all of the node's children must fail before it returns failure.

  If FAIL_ON_ONE and SUCEED_ON_ONE are both active and are both trigerred in the same time step, failure will take precedence.
*/
type FAILURE_POLICY int

const (
	FAIL_ON_ONE FAILURE_POLICY = iota
	FAIL_ON_ALL
)

/// Enumerates the options for when a parallel node is considered to have succeeded.
/**
  - SUCCEED_ON_ONE indicates that the node will return success as soon as one of its children succeeds.
  - SUCCEED_ON_ALL indicates that all of the node's children must succeed before it returns success.
*/
type SUCCESS_POLICY int

const (
	SUCCEED_ON_ONE SUCCESS_POLICY = iota
	SUCCEED_ON_ALL
)

/// Enumerates the options when a parallel node is exited
/**
  - EXIT_NONE indicates that the parallel node just exit.
  - EXIT_ABORT_RUNNINGSIBLINGS indicates that the parallel node abort all other running siblings.
*/
type EXIT_POLICY int

const (
	EXIT_NONE EXIT_POLICY = iota
	EXIT_ABORT_RUNNINGSIBLINGS
)

/// Enumerates the options of what to do when a child finishes
/**
  - CHILDFINISH_ONCE indicates that the child node just executes once.
  - CHILDFINISH_LOOP indicates that the child node run again and again.
*/
type CHILDFINISH_POLICY int

const (
	CHILDFINISH_ONCE CHILDFINISH_POLICY = iota
	CHILDFINISH_LOOP
)

// ============================================================================
type Parallel struct {
	BehaviorNodeBase
	failPolicy        FAILURE_POLICY
	succeedPolicy     SUCCESS_POLICY
	exitPolicy        EXIT_POLICY
	childFinishPolicy CHILDFINISH_POLICY
}

func NewParallel() *Parallel {
	return &Parallel{
		failPolicy:        FAIL_ON_ONE,
		succeedPolicy:     SUCCEED_ON_ALL,
		exitPolicy:        EXIT_NONE,
		childFinishPolicy: CHILDFINISH_LOOP,
	}
}

func (p *Parallel) Load(version int, agentType string, properties []property_t) {
	p.BehaviorNodeBase.Load(version, agentType, properties)
	for i := 0; i < len(properties); i++ {
		switch properties[i].name {
		case "FailurePolicy":
			switch properties[i].value {
			case "FAIL_ON_ONE":
				p.failPolicy = FAIL_ON_ONE
			case "FAIL_ON_ALL":
				p.failPolicy = FAIL_ON_ALL
			}
		case "SuccessPolicy":
			switch properties[i].value {
			case "SUCCEED_ON_ONE":
				p.succeedPolicy = SUCCEED_ON_ONE
			case "SUCCEED_ON_ALL":
				p.succeedPolicy = SUCCEED_ON_ALL
			}
		case "ExitPolicy":
			switch properties[i].value {
			case "EXIT_NONE":
				p.exitPolicy = EXIT_NONE
			case "EXIT_ABORT_RUNNINGSIBLINGS":
				p.exitPolicy = EXIT_ABORT_RUNNINGSIBLINGS
			}
		case "ChildFinishPolicy":
			switch properties[i].value {
			case "CHILDFINISH_ONCE":
				p.childFinishPolicy = CHILDFINISH_ONCE
			case "CHILDFINISH_LOOP":
				p.childFinishPolicy = CHILDFINISH_LOOP
			}
		}
	}
}

func (p *Parallel) CreateTask() BehaviorTask {
	BTGLog.Trace("(p *Parallel) CreateTask()")
	return NewParallelTask()
}

func (p *Parallel) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*Parallel); !ok {
		return false
	}
	return true
}

func (p *Parallel) Decompose(node BehaviorNode, seqTask *PlannerTaskComplex, depth int, planner *Planner) bool {
	childs := node.GetChilds()
	for i := 0; i < len(childs); i++ {
		c := childs[i]
		if c != nil {
			childTask := planner.DecomposeNode(c, depth)
			if childTask != nil {
				seqTask.AddChild(childTask)
			} else {
				break
			}
		}
	}
	return len(seqTask.childs) == len(childs)
}

func (p *Parallel) ParallelUpdate(a *Agent, childs []BehaviorTask) EBTStatus {
	sawSuccess := false
	sawFail := false
	sawRunning := false
	sawAllFails := false
	sawAllSuccess := false
	bLoop := p.childFinishPolicy == CHILDFINISH_LOOP
	for _, c := range childs {
		s := c.GetStatue()
		if bLoop || s == BT_RUNNING || s == BT_INVALID {
			s = c.Exec(a, s)
			switch s {
			case BT_FAILURE:
				sawFail = true
				sawAllSuccess = false
			case BT_SUCCESS:
				sawSuccess = true
				sawAllFails = false
			case BT_RUNNING:
				sawRunning = true
				sawAllFails = false
				sawAllSuccess = false
			}
		} else if s == BT_SUCCESS {
			sawSuccess = true
			sawAllFails = false
		} else {
			sawFail = true
			sawAllSuccess = false
		}
	}

	var result EBTStatus
	if sawRunning {
		result = BT_RUNNING
	} else {
		result = BT_FAILURE
	}

	if (p.failPolicy == FAIL_ON_ALL && sawAllFails) || (p.failPolicy == FAIL_ON_ONE && sawFail) {
		result = BT_FAILURE
	} else if (p.succeedPolicy == SUCCEED_ON_ALL && sawAllSuccess) || (p.succeedPolicy == SUCCEED_ON_ONE && sawSuccess) {
		result = BT_SUCCESS
	}
	if p.exitPolicy == EXIT_ABORT_RUNNINGSIBLINGS && (result == BT_FAILURE || result == BT_SUCCESS) {
		for _, c := range childs {
			if c.GetStatue() == BT_RUNNING {
				c.Abort(a)
			}
		}
	}
	return result
}

///Execute behaviors in parallel
/** There are two policies that control the flow of execution. The first is the policy for failure,
  and the second is the policy for success.

  For failure, the options are "fail when one child fails" and "fail when all children fail".
  For success, the options are similarly "complete when one child completes", and "complete when all children complete".
*/

type ParallelTask struct {
	CompositeTask
}

func NewParallelTask() *ParallelTask {
	return &ParallelTask{CompositeTask: CompositeTask{activeChildIndex: -1}}
}

func (pt *ParallelTask) OnEnter(a *Agent) bool {
	return pt.BehaviorTaskBase.OnEnter(a)
}

func (pt *ParallelTask) OnExit(a *Agent, status EBTStatus) {
	pt.BehaviorTaskBase.OnExit(a, status)
}

func (pt *ParallelTask) UpdateCurrent(a *Agent, childStatus EBTStatus) EBTStatus {
	return pt.Update(a, childStatus)
}

func (pt *ParallelTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	if node, ok := pt.GetNode().(*Parallel); ok {
		return node.ParallelUpdate(a, pt.childs)
	}
	return BT_FAILURE
}
