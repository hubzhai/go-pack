package behaviago

/**
 * Return values of exec/update and valid states for behaviors.
 */
type EBTStatus int

const (
	BT_INVALID EBTStatus = iota
	BT_SUCCESS
	BT_FAILURE
	BT_RUNNING
)

// ============================================================================
type EBTPhase int

const (
	EBT_SUCCESS EBTPhase = iota
	EBT_FAILURE
	EBT_BOTH
)

// ============================================================================
type BehaviorNode interface {
	//get
	GetId() int16
	GetClassName() string
	GetAgentType() string
	HasEvents() bool
	HasLoadAttachment() bool
	IsManagingChildrenAsSubTrees() bool
	GetParent() BehaviorNode
	GetSelf() BehaviorNode
	GetChilds() []BehaviorNode
	GetEffectors() []BehaviorNode
	GetEvents() []BehaviorNode
	GetPreCondition() []BehaviorNode
	GetCustomCondition() BehaviorNode
	GetPars() []*Property
	GetChildById(int16) BehaviorNode
	GetChild(int) BehaviorNode
	//set
	SetId(int16)
	SetClassName(string)
	SetAgentType(string)
	SetHasEvents(bool)
	SetLoadAttachment(bool)
	SetParent(BehaviorNode)
	SetSelf(BehaviorNode)
	SetCustomCondition(BehaviorNode)
	//op
	Clear()
	CreateTask() BehaviorTask
	AddChild(c BehaviorNode)
	AddPar(agentType, typeStr, nameStr, valueStr string)
	Attach(pAttachment BehaviorNode, bIsPrecondition, bIsEffector, bIsTransition bool)
	CheckPreconditions(a *Agent, bIsAlive bool) bool
	Load(version int, agentType string, properties []property_t)
	UpdateImpl(a *Agent, childStatus EBTStatus) EBTStatus
	Execute(a *Agent, childStatus EBTStatus) EBTStatus
	IsValid(a *Agent, task BehaviorTask) bool
	EvaluteCustomCondition(a *Agent) bool
	Evaluate(a *Agent, result EBTStatus) bool
	Decompose(node BehaviorNode, seqTask *PlannerTaskComplex, depth int, planner *Planner) bool
	ApplyEffects(a *Agent, phase EBTPhase)
	CheckEvents(eventName string, a *Agent) bool
	//
	IsDecorator() bool
}

// ============================================================================
func CombineResults(firstValidPrecond, lastCombineValue, taskBoolean bool, pi PreconditionInterface) (bool, bool) {
	if firstValidPrecond {
		firstValidPrecond = false
		lastCombineValue = taskBoolean
	} else {
		andOp := pi.IsAnd()
		if andOp {
			lastCombineValue = lastCombineValue && taskBoolean
		} else {
			lastCombineValue = lastCombineValue || taskBoolean
		}
	}
	return firstValidPrecond, lastCombineValue
}

// ============================================================================
var BehaviorNodeFactories = make(map[string]BehaviorNodeCreator)

type BehaviorNodeCreator interface {
	CreateBehaviorNode() BehaviorNode
}

type BehaviorNodeCreatorWrapper func() BehaviorNode

func (btcw BehaviorNodeCreatorWrapper) CreateBehaviorNode() BehaviorNode {
	return btcw()
}

func RegisteNodeCreator(className string, creator BehaviorNodeCreator) {
	BehaviorNodeFactories[className] = creator
}

func GetNodeCreator(className string) BehaviorNodeCreator {
	if creator, exist := BehaviorNodeFactories[className]; exist {
		return creator
	}
	return nil
}

/**
 * Base class for BehaviorTree Nodes. This is the static part
 */
type BehaviorNodeBase struct {
	className        string
	agentType        string
	id               int16
	enterPreCond     int8
	updatePreCond    int8
	bothPreCond      int8
	successEffectors int8
	failureEffectors int8
	bothEffectors    int8
	hasEvents        bool
	loadAttachment   bool
	self             BehaviorNode
	parent           BehaviorNode
	childs           []BehaviorNode
	effectors        []BehaviorNode
	events           []BehaviorNode
	preCondition     []BehaviorNode
	customCondition  BehaviorNode
	pars             []*Property
}

func NewBehaviorNode(className string) BehaviorNode {
	creator := GetNodeCreator(className)
	if creator != nil {
		return creator.CreateBehaviorNode()
	} else {
		BTGLog.Warnf("!!![%v] not registe creator function", className)
	}
	return nil
}

func (bn *BehaviorNodeBase) GetId() int16 {
	return bn.id
}
func (bn *BehaviorNodeBase) GetClassName() string {
	return bn.className
}
func (bn *BehaviorNodeBase) GetAgentType() string {
	return bn.agentType
}
func (bn *BehaviorNodeBase) HasEvents() bool {
	return bn.hasEvents
}
func (bn *BehaviorNodeBase) HasLoadAttachment() bool {
	return bn.loadAttachment
}

//return true for Parallel, SelectorLoop, etc., which is responsible to update all its children just like sub trees
//so that they are treated as a return-running node and the next update will continue them.
func (bn *BehaviorNodeBase) IsManagingChildrenAsSubTrees() bool {
	return false
}
func (bn *BehaviorNodeBase) GetSelf() BehaviorNode {
	return bn.self
}
func (bn *BehaviorNodeBase) GetParent() BehaviorNode {
	return bn.parent
}
func (bn *BehaviorNodeBase) GetChilds() []BehaviorNode {
	return bn.childs
}
func (bn *BehaviorNodeBase) GetEffectors() []BehaviorNode {
	return bn.effectors
}
func (bn *BehaviorNodeBase) GetEvents() []BehaviorNode {
	return bn.events
}
func (bn *BehaviorNodeBase) GetPreCondition() []BehaviorNode {
	return bn.preCondition
}
func (bn *BehaviorNodeBase) GetCustomCondition() BehaviorNode {
	return bn.customCondition
}
func (bn *BehaviorNodeBase) GetPars() []*Property {
	return bn.pars
}

func (bn *BehaviorNodeBase) GetChildById(nodeId int16) BehaviorNode {
	for _, c := range bn.childs {
		if c.GetId() == nodeId {
			return c
		}
	}
	return nil
}

func (bn *BehaviorNodeBase) GetChild(idx int) BehaviorNode {
	childsCnt := len(bn.childs)
	if idx >= 0 && idx < childsCnt {
		return bn.childs[idx]
	}
	return nil
}

func (bn *BehaviorNodeBase) SetId(id int16) {
	bn.id = id
}
func (bn *BehaviorNodeBase) SetClassName(cn string) {
	bn.className = cn
}
func (bn *BehaviorNodeBase) SetAgentType(at string) {
	bn.agentType = at
}
func (bn *BehaviorNodeBase) SetHasEvents(has bool) {
	bn.hasEvents = has
}
func (bn *BehaviorNodeBase) SetLoadAttachment(has bool) {
	bn.loadAttachment = has
}
func (bn *BehaviorNodeBase) SetSelf(node BehaviorNode) {
	bn.self = node
}
func (bn *BehaviorNodeBase) SetParent(node BehaviorNode) {
	bn.parent = node
}
func (bn *BehaviorNodeBase) SetCustomCondition(node BehaviorNode) {
	bn.customCondition = node
}
func (bn *BehaviorNodeBase) Clear() {
}

func (bn *BehaviorNodeBase) AddChild(c BehaviorNode) {
	BTGLog.Tracef("(bn *BehaviorNodeBase) AddChild(%v) enter (%v)", c.GetClassName(), bn.GetClassName())
	c.SetParent(bn)
	bn.childs = append(bn.childs, c)
}

func (bn *BehaviorNodeBase) AddPar(agentType, typeStr, nameStr, valueStr string) {
	p := AgentProperty_GetPropertyByVarName(agentType, nameStr)
	if p == nil {
		p = AgentProperty_AddLocal(agentType, typeStr, nameStr, valueStr)
	}
	if p != nil {
		bn.pars = append(bn.pars, p)
	}
}

func (bn *BehaviorNodeBase) Attach(pAttachment BehaviorNode, bIsPrecondition, bIsEffector, bIsTransition bool) {
	if bIsPrecondition {
		bn.preCondition = append(bn.preCondition, pAttachment)
		if predicate, ok := pAttachment.(PreconditionInterface); ok {
			phase := predicate.GetPhase()
			switch phase {
			case E_PRECOND_ENTER:
				bn.enterPreCond++
			case E_PRECOND_UPDATE:
				bn.updatePreCond++
			case E_PRECOND_BOTH:
				bn.bothPreCond++
			}
		}
	} else if bIsEffector {
		bn.effectors = append(bn.effectors, pAttachment)
		if effector, ok := pAttachment.(EffectorInterface); ok {
			phase := effector.GetPhace()
			switch phase {
			case EBT_SUCCESS:
				bn.successEffectors++
			case EBT_FAILURE:
				bn.failureEffectors++
			case EBT_BOTH:
				bn.bothEffectors++
			}
		}
	} else {
		bn.events = append(bn.events, pAttachment)
	}
}

func (bn *BehaviorNodeBase) CheckPreconditions(a *Agent, bIsAlive bool) bool {
	//satisfied if there is no preconditions
	if len(bn.preCondition) == 0 {
		return true
	}

	var phase EPreCondPhase
	if bIsAlive {
		phase = E_PRECOND_UPDATE
	} else {
		phase = E_PRECOND_ENTER
	}
	if bn.bothPreCond == 0 {
		if phase == E_PRECOND_ENTER && bn.enterPreCond == 0 {
			return true
		}
		if phase == E_PRECOND_UPDATE && bn.updatePreCond == 0 {
			return true
		}
	}

	firstValidPrecond := true
	lastCombineValue := false
	for _, c := range bn.preCondition {
		if p, ok := c.(PreconditionInterface); ok {
			ph := p.GetPhase()
			if phase == E_PRECOND_BOTH || ph == E_PRECOND_BOTH || phase == ph {
				taskBoolean := c.Evaluate(a, BT_INVALID)
				firstValidPrecond, lastCombineValue = CombineResults(firstValidPrecond, lastCombineValue, taskBoolean, p)
			}
		}
	}
	return lastCombineValue
}

func (bn *BehaviorNodeBase) Load(version int, agentType string, properties []property_t) {
	BTGWorkspace.BehaviorNodeLoaded(bn.GetClassName(), properties)
}

func (bn *BehaviorNodeBase) UpdateImpl(a *Agent, childStatus EBTStatus) EBTStatus {
	return BT_SUCCESS
}

func (bn *BehaviorNodeBase) Execute(a *Agent, childStatus EBTStatus) EBTStatus {
	return BT_SUCCESS
}

func CreateAndInitTask(node BehaviorNode) BehaviorTask {
	task := node.CreateTask()
	if task != nil {
		task.Init(task, node)
	} else {
		BTGLog.Warnf("(%v) not impletements CreateTask interface", node.GetClassName())
	}
	return task
}

func (bn *BehaviorNodeBase) CreateTask() BehaviorTask {
	return nil
}

func (bn *BehaviorNodeBase) IsValid(a *Agent, task BehaviorTask) bool {
	return true
}

func (bn *BehaviorNodeBase) EvaluteCustomCondition(a *Agent) bool {
	if bn.customCondition != nil {
		return bn.customCondition.Evaluate(a, BT_INVALID)
	}
	return false
}

func (bn *BehaviorNodeBase) Evaluate(a *Agent, result EBTStatus) bool {
	BTGLog.Warnf("Only Condition/Sequence/And/Or allowed (%v)", bn.GetClassName())
	return false
}

func (bn *BehaviorNodeBase) Decompose(node BehaviorNode, seqTask *PlannerTaskComplex, depth int, planner *Planner) bool {
	BTGLog.Tracef("BehaviorNodeBase.Decompose not impletements (%v)", bn.GetClassName())
	return false
}

func (bn *BehaviorNodeBase) ApplyEffects(a *Agent, phase EBTPhase) {
	if len(bn.effectors) == 0 {
		return
	}
	if bn.bothEffectors == 0 {
		if phase == EBT_SUCCESS && bn.successEffectors == 0 {
			return
		}
		if phase == EBT_FAILURE && bn.failureEffectors == 0 {
			return
		}
	}
	for _, c := range bn.effectors {
		if e, ok := c.(EffectorInterface); ok {
			ph := e.GetPhace()
			if phase == EBT_BOTH || ph == EBT_BOTH || ph == phase {
				e.Evaluate(a, BT_INVALID)
			}
		}
	}
}

func (bn *BehaviorNodeBase) CheckEvents(eventName string, a *Agent) bool {
	if len(bn.events) > 0 {
		for _, node := range bn.events {
			if e, ok := node.(*Event); ok {
				//check events only
				if e.GetEventName() == eventName {
					e.switchTo(a)
					if e.TriggeredOnce() {
						return false
					}
				}
			}
		}
	}
	return true
}

func (bn *BehaviorNodeBase) IsDecorator() bool {
	return false
}

// ============================================================================
type DecoratorNodeInterface interface {
	IsDecorateWhenChildEnds() bool
	SetIsDecorateWhenChildEnds(bool)
}

type DecoratorNode struct {
	BehaviorNodeBase
	isDecorateWhenChildEnds bool
}

func (dn *DecoratorNode) Load(version int, agentType string, properties []property_t) {
	dn.BehaviorNodeBase.Load(version, agentType, properties)
	for i := 0; i < len(properties); i++ {
		if properties[i].name == "DecorateWhenChildEnds" {
			if properties[i].value == "true" {
				dn.isDecorateWhenChildEnds = true
			}
		}
	}
}

func (dn *DecoratorNode) IsDecorator() bool {
	return true
}

func (dn *DecoratorNode) IsManagingChildrenAsSubTrees() bool {
	return true
}

func (dn *DecoratorNode) IsDecorateWhenChildEnds() bool {
	return dn.isDecorateWhenChildEnds
}
func (dn *DecoratorNode) SetIsDecorateWhenChildEnds(is bool) {
	dn.isDecorateWhenChildEnds = is
}

// ============================================================================
type BehaviorTree struct {
	BehaviorNodeBase
	name    string
	domains string
	isFSM   bool
}

func NewBehaviorTree() *BehaviorTree {
	bt := &BehaviorTree{}
	return bt
}
func (bt *BehaviorTree) GetName() string {
	return bt.name
}
func (bt *BehaviorTree) SetName(name string) {
	bt.name = name
}
func (bt *BehaviorTree) GetDomains() string {
	return bt.domains
}
func (bt *BehaviorTree) SetDomains(domains string) {
	bt.domains = domains
}
func (bt *BehaviorTree) IsFSM() bool {
	return bt.isFSM
}
func (bt *BehaviorTree) SetIsFSM(isFSM bool) {
	bt.isFSM = isFSM
}
func (bt *BehaviorTree) CreateTask() BehaviorTask {
	BTGLog.Trace("(bt *BehaviorTree) CreateTask()")
	btt := &BehaviorTreeTask{}
	return btt
}
func (bt *BehaviorTree) IsManagingChildrenAsSubTrees() bool {
	return true
}
