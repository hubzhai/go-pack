package behaviago

// ============================================================================
type ReferencedBehavior struct {
	BehaviorNodeBase
	referencedBehaviorPathVar    *Property
	referencedBehaviorPathMethod *Method
}

// ============================================================================
type ReferencedBehaviorTask struct {
	SingeChildTask
	nextStateId int
	subTree     *BehaviorTreeTask
}
