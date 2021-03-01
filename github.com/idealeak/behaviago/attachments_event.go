package behaviago

const (
	EventNodeName = "Event"
)

func init() {
	RegisteNodeCreator(EventNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		e := newEventNode()
		return e
	}))
}

func newEventNode() *Event {
	e := &Event{}
	e.SetClassName(EventNodeName)
	e.SetSelf(e)
	return e
}

// ============================================================================
type Event struct {
	BehaviorNodeBase
	event                  *Method
	referencedBehaviorPath string
	triggerMode            TriggerMode
	bTriggeredOnce         bool //an event can be configured to stop being checked if triggered
}

func NewEvent() *Event {
	e := &Event{
		triggerMode:    TM_Transfer,
		bTriggeredOnce: false,
	}
	return e
}

func (e *Event) GetEventName() string {
	if e.event != nil {
		return e.event.MethodName
	}
	return ""
}

func (e *Event) TriggeredOnce() bool {
	return e.bTriggeredOnce
}

func (e *Event) GetTriggerMode() TriggerMode {
	return e.triggerMode
}

func (e *Event) Load(version int, agentType string, properties []property_t) {
	e.BehaviorNodeBase.Load(version, agentType, properties)
	for i := 0; i < len(properties); i++ {
		switch properties[i].name {
		case "Task":
			e.event = LoadMethod(properties[i].value)
		case "ReferenceFilename":
			e.referencedBehaviorPath = properties[i].value
		case "TriggeredOnce":
			if properties[i].value == "true" {
				e.bTriggeredOnce = true
			}
		case "TriggerMode":
			switch properties[i].value {
			case "Transfer":
				e.triggerMode = TM_Transfer
			case "Return":
				e.triggerMode = TM_Return
			}
		}
	}
}

func (e *Event) CreateTask() BehaviorTask {
	return &EventTask{}
}

func (e *Event) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*Event); !ok {
		return false
	}
	return true
}

func (e *Event) switchTo(a *Agent) EBTStatus {
	if e.referencedBehaviorPath != "" {
		if a != nil {
			a.BTGSetCurrent(e.referencedBehaviorPath, e.triggerMode, true)
			return a.BTGExec()
		}
	}
	return BT_SUCCESS
}

// ============================================================================
type EventTask struct {
	AttachmentTask
}

func (e *EventTask) GetEventName() string {
	if node, ok := e.GetNode().(*Event); ok {
		return node.GetEventName()
	}
	return ""
}

func (e *EventTask) TriggeredOnce() bool {
	if node, ok := e.GetNode().(*Event); ok {
		return node.TriggeredOnce()
	}
	return false
}

func (e *EventTask) GetTriggerMode() TriggerMode {
	if node, ok := e.GetNode().(*Event); ok {
		return node.GetTriggerMode()
	}
	return TM_Transfer
}

func (e *EventTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	if node, ok := e.GetNode().(*Event); ok {
		return node.switchTo(a)
	}
	return BT_SUCCESS
}
