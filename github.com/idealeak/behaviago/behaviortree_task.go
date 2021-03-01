package behaviago

type TriggerMode int

/**
  trigger mode to control the bt switching and back
*/
const (
	TM_Transfer TriggerMode = iota
	TM_Return
)

// ============================================================================

type NodeHandler func(BehaviorTask, *Agent, interface{}) bool

func abortHandler(node BehaviorTask, a *Agent, userData interface{}) bool {
	if node.GetStatue() == BT_RUNNING {
		node.OnExitAction(a, BT_FAILURE)
		node.SetStatue(BT_FAILURE)
		node.SetCurrentTask(nil)
	}

	return true
}

func resetHandler(node BehaviorTask, a *Agent, userData interface{}) bool {
	node.SetStatue(BT_INVALID)
	node.SetCurrentTask(nil)
	node.OnReset(a)
	return true
}

type tempRunningNodes struct {
	nodes []BehaviorTask
}

func getRunningNodesHandler(node BehaviorTask, a *Agent, userData interface{}) bool {
	if node.GetStatue() == BT_RUNNING {
		if n, ok := userData.(*tempRunningNodes); ok {
			n.nodes = append(n.nodes, node)
		}
	}
	return true
}

// ============================================================================
type BehaviorTask interface {
	//get
	GetId() int16
	HasManagingParent() bool
	GetStatue() EBTStatus
	GetNode() BehaviorNode
	GetSelf() BehaviorTask
	GetParent() BehaviorTask
	GetAttachments() []BehaviorTask
	GetNextStateId() int
	GetCurrentTask() BehaviorTask
	GetRunningNodes(bool) []BehaviorTask
	GetTopManageBranchTask() BehaviorTask
	//set
	SetId(int16)
	SetHasManagingParent(bool)
	SetStatue(EBTStatus)
	SetNode(BehaviorNode)
	SetSelf(BehaviorTask)
	SetParent(BehaviorTask)
	SetCurrentTask(BehaviorTask)
	//op
	Init(self BehaviorTask, node BehaviorNode)
	Attach(at BehaviorTask)
	GetTaskById(id int16) BehaviorTask
	GetClassNameString() string
	OnReset(a *Agent)
	OnEnter(a *Agent) bool
	OnEnterAction(a *Agent) bool
	CheckEvents(eventName string, a *Agent) bool
	OnEvent(a *Agent, event string) bool
	Exec(a *Agent, childStatus EBTStatus) EBTStatus
	Update(a *Agent, childStatus EBTStatus) EBTStatus
	UpdateCurrent(a *Agent, childStatus EBTStatus) EBTStatus
	OnExit(a *Agent, status EBTStatus)
	CheckPreconditions(a *Agent, bIsAlive bool) bool
	OnExitAction(a *Agent, status EBTStatus)
	Abort(a *Agent)
	Reset(a *Agent)
	traverse(childFirst bool, handler NodeHandler, a *Agent, userData interface{})
	CheckParentUpdatePreconditions(a *Agent) bool
	Decorate(status EBTStatus) EBTStatus

	//dynamic
	IsTree() bool
	IsBranchTask() bool
	IsLeafTask() bool
}

/**
  Base class for the BehaviorTreeTask's runtime execution management.
*/
type BehaviorTaskBase struct {
	id                int16
	hasManagingParent bool
	status            EBTStatus
	node              BehaviorNode
	self              BehaviorTask
	parent            BehaviorTask
	attachments       []BehaviorTask
}

func (bt *BehaviorTaskBase) GetId() int16 {
	return bt.id
}
func (bt *BehaviorTaskBase) HasManagingParent() bool {
	return bt.hasManagingParent
}
func (bt *BehaviorTaskBase) GetStatue() EBTStatus {
	return bt.status
}
func (bt *BehaviorTaskBase) GetNode() BehaviorNode {
	return bt.node
}
func (bt *BehaviorTaskBase) GetSelf() BehaviorTask {
	return bt.self
}
func (bt *BehaviorTaskBase) GetParent() BehaviorTask {
	return bt.parent
}
func (bt *BehaviorTaskBase) GetAttachments() []BehaviorTask {
	return bt.attachments
}
func (bt *BehaviorTaskBase) GetNextStateId() int {
	return -1
}
func (bt *BehaviorTaskBase) GetCurrentTask() BehaviorTask {
	return nil
}

func (bt *BehaviorTaskBase) GetRunningNodes(onlyLeaves bool) []BehaviorTask {
	var temp tempRunningNodes
	bt.traverse(true, getRunningNodesHandler, nil, &temp)
	if onlyLeaves && len(temp.nodes) > 0 {
		var ret []BehaviorTask
		cnt := len(temp.nodes)
		for i := 0; i < cnt; i++ {
			if temp.nodes[i].IsLeafTask() {
				ret = append(ret, temp.nodes[i])
			}
		}
		return ret
	} else {
		return temp.nodes
	}
}

func (bt *BehaviorTaskBase) SetId(id int16) {
	bt.id = id
}
func (bt *BehaviorTaskBase) SetHasManagingParent(has bool) {
	bt.hasManagingParent = has
}
func (bt *BehaviorTaskBase) SetStatue(s EBTStatus) {
	bt.status = s
}
func (bt *BehaviorTaskBase) SetNode(n BehaviorNode) {
	bt.node = n
}
func (bt *BehaviorTaskBase) SetSelf(n BehaviorTask) {
	bt.self = n
}
func (bt *BehaviorTaskBase) SetParent(n BehaviorTask) {
	bt.parent = n
}
func (bt *BehaviorTaskBase) SetCurrentTask(BehaviorTask) {
}
func (bt *BehaviorTaskBase) Init(self BehaviorTask, node BehaviorNode) {
	bt.self = self
	bt.node = node
	bt.id = node.GetId()
}

func (bt *BehaviorTaskBase) Attach(at BehaviorTask) {
	bt.attachments = append(bt.attachments, at)
}

func (bt *BehaviorTaskBase) GetTaskById(id int16) BehaviorTask {
	if bt.id == id {
		return bt
	}
	return nil
}

func (bt *BehaviorTaskBase) GetClassNameString() string {
	if bt.node != nil {
		return bt.node.GetClassName()
	}
	return ""
}

func (bt *BehaviorTaskBase) OnReset(a *Agent) {
	BTGLog.Tracef("BehaviorTaskBase.OnReset not impletmented(%v)", bt.GetClassNameString())
}

func (bt *BehaviorTaskBase) OnEnter(a *Agent) bool {
	BTGLog.Tracef("BehaviorTaskBase.OnEnter not impletmented(%v)", bt.GetClassNameString())
	return true
}

func (bt *BehaviorTaskBase) OnEnterAction(a *Agent) bool {
	BTGLog.Tracef("(bt *BehaviorTaskBase) OnEnterAction enter(%v)", bt.GetClassNameString())
	if bt.self.CheckPreconditions(a, false) {
		bt.hasManagingParent = false
		bt.self.SetCurrentTask(nil)
		return bt.self.OnEnter(a)
	}
	return false
}

/**
  return false if the event handling needs to be stopped

  an event can be configured to stop being checked if triggered
*/
func (bt *BehaviorTaskBase) CheckEvents(eventName string, a *Agent) bool {
	BTGLog.Tracef("(btt *BehaviorTaskBase) CheckEvents enter(%v)", bt.GetClassNameString())
	if bt.node != nil {
		return bt.node.CheckEvents(eventName, a)
	}
	return false
}

/**
  return false if the event handling  needs to be stopped
  return true, the event hanlding will be checked furtherly
*/
func (bt *BehaviorTaskBase) OnEvent(a *Agent, event string) bool {
	BTGLog.Tracef("BehaviorTaskBase.OnEvent not impletmented(%v)", bt.GetClassNameString())
	return true
}

func (bt *BehaviorTaskBase) Exec(a *Agent, childStatus EBTStatus) EBTStatus {
	BTGLog.Tracef("(btt *BehaviorTaskBase) Exec enter(%v)", bt.GetClassNameString())
	bEnterResult := false
	if bt.status == BT_RUNNING {
		bEnterResult = true
	} else {
		//reset it to invalid when it was success/failure
		bt.status = BT_INVALID
		bEnterResult = bt.self.OnEnterAction(a)
	}
	if bEnterResult {
		if bt.self.CheckParentUpdatePreconditions(a) {
			BTGLog.Tracef("(btt *BehaviorTaskBase) Exec UpdateCurrent before status=%v (%v)", bt.status, bt.GetClassNameString())
			bt.status = bt.self.UpdateCurrent(a, childStatus)
			BTGLog.Tracef("(btt *BehaviorTaskBase) Exec UpdateCurrent after status=%v (%v)", bt.status, bt.GetClassNameString())
		} else {
			bt.status = BT_FAILURE
			if bt.self.GetCurrentTask() != nil {
				BTGLog.Tracef("(btt *BehaviorTaskBase) Exec UpdateCurrent1 before status=%v (%v)", bt.status, bt.GetClassNameString())
				bt.self.UpdateCurrent(a, BT_FAILURE)
				BTGLog.Tracef("(btt *BehaviorTaskBase) Exec UpdateCurrent1 after status=%v (%v)", bt.status, bt.GetClassNameString())
			}
		}

		if bt.status != BT_RUNNING {
			//clear it
			bt.self.OnExitAction(a, bt.status)
			//this node is possibly ticked by its parent or by the topBranch who records it as currrent node
			//so, we can't here reset the topBranch's current node
		} else {
			tree := bt.self.GetTopManageBranchTask()
			if tree != nil {
				tree.SetCurrentTask(bt)
			}
		}
	} else {
		bt.status = BT_FAILURE
	}
	return bt.status
}

func (bt *BehaviorTaskBase) CheckParentUpdatePreconditions(a *Agent) bool {
	BTGLog.Tracef("(btt *BehaviorTaskBase) CheckParentUpdatePreconditions enter(%v)", bt.GetClassNameString())
	bValid := true
	if bt.hasManagingParent {
		bHasManagingParent := false
		parents := make([]BehaviorTask, 0, 512)
		parents = append(parents, bt)
		parentBranch := bt.GetParent()

		//back track the parents until the managing branch
		for parentBranch != nil {
			parents = append(parents, parentBranch)
			if parentBranch.GetCurrentTask() == bt.GetSelf() {
				bHasManagingParent = true
				break
			}
			parentBranch = parentBranch.GetParent()
		}

		if bHasManagingParent {
			for i := len(parents) - 1; i >= 0; i-- {
				bValid = parents[i].CheckPreconditions(a, true)
				if !bValid {
					break
				}
			}
		}
	} else {
		bValid = bt.self.CheckPreconditions(a, true)
	}
	return bValid
}
func (bt *BehaviorTaskBase) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	BTGLog.Tracef("BehaviorTaskBase.Update not impletmented(%v)", bt.GetClassNameString())
	return BT_SUCCESS
}

func (bt *BehaviorTaskBase) UpdateCurrent(a *Agent, childStatus EBTStatus) EBTStatus {
	BTGLog.Tracef("BehaviorTaskBase.UpdateCurrent (%v)", bt.GetClassNameString())
	return bt.self.Update(a, childStatus)
}

func (bt *BehaviorTaskBase) OnExit(a *Agent, status EBTStatus) {
	BTGLog.Tracef("BehaviorTaskBase.OnExit not impletmented(%v)", bt.GetClassNameString())
}

func (bt *BehaviorTaskBase) CheckPreconditions(a *Agent, bIsAlive bool) bool {
	BTGLog.Tracef("(bt *BehaviorTaskBase) CheckPreconditions enter(%v)", bt.GetClassNameString())
	bResult := true
	if bt.node != nil {
		if len(bt.node.GetPreCondition()) != 0 {
			bResult = bt.node.CheckPreconditions(a, bIsAlive)
		}
	}
	return bResult
}

func (bt *BehaviorTaskBase) OnExitAction(a *Agent, status EBTStatus) {
	BTGLog.Tracef("(bt *BehaviorTaskBase) OnExitAction enter(%v)", bt.GetClassNameString())
	bt.self.OnExit(a, status)
	if bt.node != nil {
		phase := EBT_SUCCESS
		if status == BT_FAILURE {
			phase = EBT_FAILURE
		}
		bt.node.ApplyEffects(a, phase)
	}
}

func (bt *BehaviorTaskBase) Abort(a *Agent) {
	bt.self.traverse(true, abortHandler, a, nil)
}

func (bt *BehaviorTaskBase) Reset(a *Agent) {
	bt.self.traverse(true, resetHandler, a, nil)
}

func (bt *BehaviorTaskBase) traverse(childFirst bool, handler NodeHandler, a *Agent, userData interface{}) {
}

func (bt *BehaviorTaskBase) GetTopManageBranchTask() BehaviorTask {
	BTGLog.Tracef("(btt *BehaviorTaskBase) GetTopManageBranchTask enter(%v)", bt.GetClassNameString())
	var tree BehaviorTask = nil
	task := bt.GetParent()
	for task != nil {
		if task.IsTree() {
			//to overwrite the child branch
			tree = task
			break
		} else if task.GetNode() != nil && task.GetNode().IsManagingChildrenAsSubTrees() {
			//until it is Parallel/SelectorLoop, it's child is used as tree to store current task
			break
		} else if task.IsBranchTask() {
			//this if must be after BehaviorTreeTask and IsManagingChildrenAsSubTrees
			tree = task
		}
		task = task.GetParent()
	}
	return tree
}

func (bt *BehaviorTaskBase) Decorate(status EBTStatus) EBTStatus {
	BTGLog.Tracef("(btt *BehaviorTaskBase) Decorate enter(%v)", bt.GetClassNameString())
	return status
}

func (bt *BehaviorTaskBase) IsTree() bool {
	return false
}

func (bt *BehaviorTaskBase) IsBranchTask() bool {
	return false
}

func (bt *BehaviorTaskBase) IsLeafTask() bool {
	return false
}

// ============================================================================
type AttachmentTask struct {
	BehaviorTaskBase
}

func (at *AttachmentTask) traverse(childFirst bool, handler NodeHandler, a *Agent, userData interface{}) {
	if handler != nil {
		handler(at, a, userData)
	}
}

// ============================================================================
type LeafTask struct {
	BehaviorTaskBase
}

func (lt *LeafTask) IsLeafTask() bool {
	return true
}

func (lt *LeafTask) traverse(childFirst bool, handler NodeHandler, a *Agent, userData interface{}) {
	if handler != nil {
		handler(lt, a, userData)
	}
}

// ============================================================================
type BranchTask struct {
	BehaviorTaskBase
	currentNodeId int
	currentTask   BehaviorTask
}

func (bt *BranchTask) IsBranchTask() bool {
	return true
}

func (bt *BranchTask) GetCurrentNodeId() int {
	return bt.currentNodeId
}

func (bt *BranchTask) SetCurrentNodeId(nodeId int) {
	bt.currentNodeId = nodeId
}

func (bt *BranchTask) GetCurrentTask() BehaviorTask {
	return bt.currentTask
}

func (bt *BranchTask) oneventCurrentNode(a *Agent, event string) bool {
	BTGLog.Tracef("(btt *BranchTask) oneventCurrentNode enter(%v)", bt.GetClassNameString())
	bGoOn := bt.currentTask.OnEvent(a, event)
	//give the handling back to parents
	if bGoOn && bt.currentTask != nil {
		parent := bt.currentTask.GetParent()
		//back track the parents until the branch
		for parent != nil && parent != bt.GetSelf() {
			bGoOn = parent.OnEvent(a, event)
			if !bGoOn {
				return false
			}
			parent = parent.GetParent()
		}
	}
	return bGoOn
}

func (bt *BranchTask) OnEvent(a *Agent, event string) bool {
	BTGLog.Tracef("(btt *BranchTask) OnEvent enter(%v)", bt.GetClassNameString())
	if bt.GetNode().HasEvents() {
		bGoOn := true
		if bt.currentTask != nil {
			bGoOn = bt.oneventCurrentNode(a, event)
			if bGoOn {
				bGoOn = bt.BehaviorTaskBase.OnEvent(a, event)
			}
		}
	}
	return true
}

//
//Set the currentTask as task
//if the leaf node is runninng ,then we should set the leaf's parent node also as running
//
func (bt *BranchTask) SetCurrentTask(task BehaviorTask) {
	BTGLog.Tracef("(btt *BranchTask) SetCurrentTask enter(%v)", bt.GetClassNameString())
	if task != nil {
		//if the leaf node is running, then the leaf's parent node is also as running,
		//the leaf is set as the tree's current task instead of its parent
		if bt.currentTask == nil {
			bt.currentTask = task
			task.SetHasManagingParent(true)
		}
	} else {
		if bt.status != BT_RUNNING {
			bt.currentTask = task
		}
	}
}

func (bt *BranchTask) execCurrentTask(a *Agent, childStatus EBTStatus) EBTStatus {
	BTGLog.Tracef("(btt *BranchTask) execCurrentTask enter(%v)", bt.GetClassNameString())
	status := bt.currentTask.Exec(a, childStatus)
	//give the handling back to parents
	if status != BT_RUNNING {
		parentBranch := bt.currentTask.GetParent()

		bt.currentTask = nil

		//back track the parents until the branch
		for parentBranch != nil {
			if parentBranch == bt.GetSelf() {
				status = parentBranch.Update(a, status)
			} else {
				status = parentBranch.Exec(a, status)
			}
			if status == BT_RUNNING {
				return status
			}
			if parentBranch == bt.GetSelf() {
				break
			}
			parentBranch = parentBranch.GetParent()
		}
	}

	return status
}

func (bt *BranchTask) resumeBranch(a *Agent, status EBTStatus) EBTStatus {
	BTGLog.Tracef("(btt *BranchTask) resumeBranch enter(%v)", bt.GetClassNameString())
	var parent BehaviorTask
	node := bt.currentTask.GetNode()
	if node.IsManagingChildrenAsSubTrees() {
		parent = bt.currentTask
	} else {
		parent = bt.currentTask.GetParent()
	}
	//clear it as it ends and the next exec might need to set it
	bt.currentTask = nil
	if parent != nil {
		return parent.Exec(a, status)
	}
	return status
}

func (bt *BranchTask) UpdateCurrent(a *Agent, childStatus EBTStatus) EBTStatus {
	BTGLog.Tracef("(btt *BranchTask) UpdateCurrent enter(%v)", bt.GetClassNameString())
	var status EBTStatus
	if bt.currentTask != nil {
		status = bt.execCurrentTask(a, childStatus)
	} else {
		status = bt.self.Update(a, childStatus)
	}
	return status
}

// ============================================================================
type CompositeTask struct {
	BranchTask
	childs           []BehaviorTask
	activeChildIndex int //book mark the current child
}

func (ct *CompositeTask) Init(self BehaviorTask, node BehaviorNode) {
	BTGLog.Tracef("(ct *CompositeTask) Init enter(%v)", ct.GetClassNameString())
	ct.BranchTask.Init(self, node)
	childs := node.GetChilds()
	for _, c := range childs {
		task := CreateAndInitTask(c)
		if task != nil {
			ct.AddChild(task)
		}
	}
}

func (ct *CompositeTask) AddChild(c BehaviorTask) {
	BTGLog.Tracef("(ct *CompositeTask) AddChild(%v) enter(%v)", c.GetClassNameString(), ct.GetClassNameString())
	ct.childs = append(ct.childs, c)
	c.SetParent(ct.self)
}

func (ct *CompositeTask) GetChildById(nodeId int) BehaviorTask {
	for _, c := range ct.childs {
		if c.GetId() == int16(nodeId) {
			return c
		}
	}
	return nil
}

func (ct *CompositeTask) GetTaskById(id int16) BehaviorTask {
	t := ct.BranchTask.GetTaskById(id)
	if t != nil {
		return t
	}
	for _, tt := range ct.childs {
		t := tt.GetTaskById(id)
		if t != nil {
			return t
		}
	}
	return nil
}

func (ct *CompositeTask) traverse(childFirst bool, handler NodeHandler, a *Agent, userData interface{}) {
	BTGLog.Tracef("(ct *CompositeTask) traverse enter(%v)", ct.GetClassNameString())
	if childFirst {
		for _, c := range ct.childs {
			c.traverse(childFirst, handler, a, userData)
		}
		handler(ct, a, userData)
	} else {
		if handler(ct, a, userData) {
			for _, c := range ct.childs {
				c.traverse(childFirst, handler, a, userData)
			}
		}
	}
}

// ============================================================================
type SingeChildTask struct {
	BranchTask
	root BehaviorTask
}

func (sct *SingeChildTask) Init(self BehaviorTask, node BehaviorNode) {
	BTGLog.Tracef("(sct *SingeChildTask) Init enter(%v)", sct.GetClassNameString())
	sct.BranchTask.Init(self, node)
	childs := node.GetChilds()
	if len(childs) == 1 {
		task := CreateAndInitTask(childs[0])
		if task != nil {
			sct.AddChild(task)
		}
	}
}

func (sct *SingeChildTask) AddChild(c BehaviorTask) {
	BTGLog.Tracef("(sct *SingeChildTask) AddChild(%v) enter(%v)",c.GetClassNameString(), sct.GetClassNameString())
	c.SetParent(sct.self)
	sct.root = c
}

func (sct *SingeChildTask) GetTaskById(id int16) BehaviorTask {
	BTGLog.Tracef("(sct *SingeChildTask) GetTaskById enter(%v)", sct.GetClassNameString())
	t := sct.BranchTask.GetTaskById(id)
	if t != nil {
		return t
	}
	if sct.root != nil {
		if sct.root.GetId() == id {
			return sct.root
		}
		return sct.root.GetTaskById(id)
	}
	return nil
}

func (sct *SingeChildTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	BTGLog.Tracef("(sct *SingeChildTask) Update enter(%v)", sct.GetClassNameString())
	if sct.root != nil {
		return sct.root.Exec(a, childStatus)
	}
	return BT_FAILURE
}

func (sct *SingeChildTask) traverse(childFirst bool, handler NodeHandler, a *Agent, userData interface{}) {
	BTGLog.Tracef("(sct *SingeChildTask) traverse enter(%v)", sct.GetClassNameString())
	if childFirst {
		if sct.root != nil {
			sct.root.traverse(childFirst, handler, a, userData)
		}
		handler(sct, a, userData)
	} else {
		if handler(sct, a, userData) {
			if sct.root != nil {
				sct.root.traverse(childFirst, handler, a, userData)
			}
		}
	}
}

// ============================================================================
type DecoratorTask struct {
	SingeChildTask
	DecorateWhenChildEnds bool
}

func (dt *DecoratorTask) Init(self BehaviorTask, node BehaviorNode) {
	BTGLog.Tracef("(dt *DecoratorTask) Init enter(%v)", dt.GetClassNameString())
	dt.SingeChildTask.Init(self, node)
	if node.IsDecorator() {
		if dti, ok := node.(DecoratorNodeInterface); ok {
			dt.DecorateWhenChildEnds = dti.IsDecorateWhenChildEnds()
		}
	}
}

func (dt *DecoratorTask) Decorate(status EBTStatus) EBTStatus {
	BTGLog.Tracef("(dt *DecoratorTask) Decorate enter(%v)", dt.GetClassNameString())
	return status
}

func (dt *DecoratorTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	BTGLog.Tracef("(dt *DecoratorTask) Update enter(%v)", dt.GetClassNameString())
	status := BT_INVALID
	if node, ok := dt.GetNode().(DecoratorNodeInterface); ok {
		if childStatus != BT_RUNNING {
			status = childStatus
			if !node.IsDecorateWhenChildEnds() || status != BT_RUNNING {
				status = dt.self.Decorate(status)
				if status != BT_RUNNING {
					return status
				}
				return BT_RUNNING
			}
		}

		status = dt.SingeChildTask.Update(a, childStatus)
		if !node.IsDecorateWhenChildEnds() || status != BT_RUNNING {
			BTGLog.Tracef("(dt *DecoratorTask) Update enter(%v) status=%v", dt.GetClassNameString(), status)
			return dt.self.Decorate(status)
		}
	}
	return BT_RUNNING
}

// ============================================================================
type BehaviorTreeTask struct {
	SingeChildTask
}

func (btt *BehaviorTreeTask) IsTree() bool {
	return true
}

func (btt *BehaviorTreeTask) GetName() string {
	if bt, ok := btt.GetNode().(*BehaviorTree); ok {
		return bt.GetName()
	}
	return ""
}

func (btt *BehaviorTreeTask) UpdateCurrent(a *Agent, childStatus EBTStatus) EBTStatus {
	BTGLog.Trace("(btt *BehaviorTreeTask) UpdateCurrent enter")
	status := BT_RUNNING
	if tree, ok := btt.GetNode().(*BehaviorTree); ok {
		if tree.IsFSM() {
			status = btt.Update(a, childStatus)
		} else {
			status = btt.SingeChildTask.UpdateCurrent(a, childStatus)
		}
	}
	return status
}

func (btt *BehaviorTreeTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	BTGLog.Trace("(btt *BehaviorTreeTask) Update enter")
	if childStatus != BT_RUNNING {
		return childStatus
	}
	return btt.SingeChildTask.Update(a, childStatus)
}

func (btt *BehaviorTreeTask) SetRootTask(root BehaviorTask) {
	BTGLog.Trace("(btt *BehaviorTreeTask) SetRootTask enter")
	btt.AddChild(root)
}

func (btt *BehaviorTreeTask) OnEnter(a *Agent) bool {
	BTGLog.Tracef("===(btt *BehaviorTreeTask)OnEnter #%v %v ===", a.GetAgentName(), btt.GetName())
	return true
}

func (btt *BehaviorTreeTask) OnExit(a *Agent, status EBTStatus) {
	BTGLog.Tracef("===(btt *BehaviorTreeTask)OnExit #%v %v ===", a.GetAgentName(), btt.GetName())
	btt.SingeChildTask.OnExit(a, status)
}

func (btt *BehaviorTreeTask) Resume(a *Agent, status EBTStatus) EBTStatus {
	BTGLog.Tracef("===(btt *BehaviorTreeTask)Resume #%v %v ===", a.GetAgentName(), btt.GetName())
	return btt.SingeChildTask.resumeBranch(a, status)
}
