package behaviago

import "math/rand"

/**
  Base class of Stochastic Nodes
*/
const (
	CompositeStochasticNodeName = "CompositeStochastic"
)

func init() {
	RegisteNodeCreator(CompositeStochasticNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		c := newCompositeStochasticNode()
		return c
	}))
}

func newCompositeStochasticNode() *CompositeStochastic {
	c := &CompositeStochastic{}
	c.SetClassName(CompositeStochasticNodeName)
	c.SetSelf(c)
	return c
}

// ============================================================================
func GetRandomValue(m *Method, a *Agent) float64 {
	BTGLog.Trace("GetRandomValue enter")
	if m == nil || a == nil {
		return rand.Float64()
	}
	vals := m.Invoke(a)
	if len(vals) != 0 {
		return vals[0].Float()
	}
	return rand.Float64()
}

// ============================================================================
type CompositeStochastic struct {
	BehaviorNodeBase
	method *Method
}

func (c *CompositeStochastic) Load(version int, agentType string, properties []property_t) {
	c.BehaviorNodeBase.Load(version, agentType, properties)
	for i := 0; i < len(properties); i++ {
		if properties[i].name == "RandomGenerator" {
			c.method = LoadMethod(properties[i].value)
		}
	}
}

func (c *CompositeStochastic) CreateTask() BehaviorTask {
	BTGLog.Trace("(c *CompositeStochastic) CreateTask()")
	return &CompositeStochasticTask{}
}

func (c *CompositeStochastic) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*CompositeStochastic); !ok {
		return false
	}
	return true
}

// ============================================================================
type CompositeStochasticTask struct {
	CompositeTask
	set []int
}

func (cst *CompositeStochasticTask) randomChild(a *Agent) {
	if n, ok := cst.GetNode().(*CompositeStochastic); ok {
		cst.set = nil
		for i := 0; i < len(cst.childs); i++ {
			cst.set = append(cst.set, i)
		}
		for i := 0; i < len(cst.set); i++ {
			j := int(GetRandomValue(n.method, a) * float64(i+1))
			cst.set[i], cst.set[j] = cst.set[j], cst.set[i]
		}
	}
}

func (cst *CompositeStochasticTask) OnEnter(a *Agent) bool {
	cst.randomChild(a)
	cst.activeChildIndex = 0
	return cst.CompositeTask.OnEnter(a)
}

func (cst *CompositeStochasticTask) OnExit(a *Agent, status EBTStatus) {
	cst.CompositeTask.OnExit(a, status)
}

func (cst *CompositeStochasticTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	bFirst := true
	if cst.activeChildIndex != -1 && len(cst.set) != 0 {
		// Keep going until a child behavior says its running.
		for {
			s := childStatus
			if !bFirst || s == BT_RUNNING {
				childIdx := cst.set[cst.activeChildIndex]
				c := cst.childs[childIdx]
				if c != nil {
					s = c.Exec(a, s)
				}
				bFirst = false
				// If the child succeeds, or keeps running, do the same.
				if s != BT_FAILURE {
					return s
				}
				// Hit the end of the array, job done!
				cst.activeChildIndex++
				if cst.activeChildIndex >= len(cst.childs) {
					return BT_FAILURE
				}
			}
		}
	}
	return BT_FAILURE
}
