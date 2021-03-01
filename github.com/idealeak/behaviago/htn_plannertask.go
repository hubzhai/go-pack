package behaviago

// ============================================================================
type PlannerTask interface {
	BehaviorTask
	IsHigherPriority(other PlannerTask) bool
	CannNotInterruptable() bool
}

func CreatePlannerTask(node BehaviorNode, a *Agent) PlannerTask {
	return nil
}

// ============================================================================
type PlannerTaskBase struct {
	BehaviorTaskBase
	NotInterruptable bool
}

func (ptb *PlannerTaskBase) IsHigherPriority(other PlannerTask) bool {
	return true
}

func (ptb *PlannerTaskBase) CannNotInterruptable() bool {
	return ptb.NotInterruptable
}

// ============================================================================
type PlannerTaskAction struct {
	PlannerTaskBase
}

// ============================================================================
type PlannerTaskComplex struct {
	PlannerTaskBase
	childs []PlannerTask
}

func (ptc *PlannerTaskComplex) AddChild(task PlannerTask) {
	ptc.childs = append(ptc.childs, task)
	task.SetParent(ptc.self)
}

// ============================================================================
type PlannerTaskSequence struct {
	PlannerTaskComplex
}

// ============================================================================
type PlannerTaskParallel struct {
	PlannerTaskComplex
}

// ============================================================================
type PlannerTaskLoop struct {
	PlannerTaskComplex
	n int
}

// ============================================================================
type PlannerTaskIterator struct {
	PlannerTaskComplex
	index int
}

// ============================================================================
type PlannerTaskReference struct {
	PlannerTaskComplex
	subTree *BehaviorTreeTask
}

// ============================================================================
type PlannerTaskTask struct {
	PlannerTaskComplex
}

// ============================================================================
type PlannerTaskMethod struct {
	PlannerTaskComplex
}
