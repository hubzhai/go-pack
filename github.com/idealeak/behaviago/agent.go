package behaviago

import (
	"errors"
	"reflect"
	"strings"
)

var (
	idMask         uint
	agentAutoIndex int
)

var (
	BTGAgents = make(map[string]*Agent)
)

// ============================================================================
type BehaviorTreeStackItem_t struct {
	btt            *BehaviorTreeTask
	tm             TriggerMode
	triggerByEvent bool
}

// ============================================================================
func GetAgentInstance(a *Agent, agentInstanceName string) *Agent {
	if len(agentInstanceName) != 0 && !strings.HasPrefix(agentInstanceName, "Self") {
		var contextId int
		if a != nil {
			contextId = a.ContextId
		}
		p := GetAgentInstanceInContext(agentInstanceName, contextId)
		if p == nil && a != nil {
			v := a.GetVariable(agentInstanceName)
			if v.IsValid() {
				ref := v
				kind := v.Type().Kind()
				if kind == reflect.Ptr || kind == reflect.Interface {
					ref = v.Elem()
					if ref.Type().Name() == "Agent" {
						if reta, ok := v.Interface().(*Agent); ok {
							return reta
						}
					}
				}
			}
		}
		return p
	}
	return a
}

func GetAgentInstanceInContext(agentInstanceName string, contextId int) *Agent {
	ctx := GetContext(contextId)
	if ctx == nil {
		return nil
	}
	return ctx.GetInstance(agentInstanceName)
}

func NewAgent(contextId, priority int, agentInstanceName string) *Agent {
	a := &Agent{
		Id:                agentAutoIndex,
		Name:              agentInstanceName,
		Priority:          priority,
		ContextId:         contextId,
		IdFlag:            0xffffffff,
		IsActive:          true,
		IsReferenceTree:   false,
		IsBlackboardBound: false,
	}
	c := GetContext(contextId)
	if c != nil {
		c.AddAgent(a)
	}
	agentAutoIndex++
	return a
}

/// The Agent class is the base class to manage(load/unload/exec) behaviors.
// ============================================================================
type Agent struct {
	Id                int
	Name              string
	Priority          int
	ContextId         int
	PlanningTop       int
	IdFlag            uint
	IsActive          bool
	IsReferenceTree   bool
	IsBlackboardBound bool
	client            interface{} //委托人
	bttStack          []BehaviorTreeStackItem_t
	behaviorTreeTasks []*BehaviorTreeTask
	bttCurrent        *BehaviorTreeTask
}

func SetIdMast(mask uint) {
	idMask = mask
}

func GetIdMask() uint {
	return idMask
}

func (a *Agent) IsMasked() bool {
	return a.IdFlag&GetIdMask() != 0
}

func (a *Agent) GetClient() interface{} {
	return a.client
}

func (a *Agent) SetClient(c interface{}) error {
	if c == nil {
		return errors.New("agent of client not nil")
	}
	k := reflect.TypeOf(c).Kind()
	if k != reflect.Ptr {
		return errors.New("agent of client must be pointer")
	}
	a.client = c
	return nil
}

func (a *Agent) GetName() string {
	return a.Name
}

func (a *Agent) GetAgentName() string {
	client := reflect.ValueOf(a.client)
	if client.IsValid() && client.Kind() == reflect.Ptr {
		return client.Elem().Type().Name()
	}
	return "unknow"
}

/**
  @param relativePath, relativePath is relative to the workspace exported path. relativePath should not include extension.
  the file format(xml/bson) is specified by Init.
  @param bForce, the loaded BT is kept in the cache so the subsequent loading will just return it from the cache.
  if bForce is true, it will not check the cache and force to load it.

  @return
  return nil if successfully loaded
*/
func (a *Agent) BTGLoad(relativePath string, bForce bool) error {
	err := BTGWorkspace.Load(relativePath, bForce)
	if err == nil {
		BTGWorkspace.RecordBTGAgentMapping(relativePath, a)
	}
	return err
}

/**
  exec the BT specified by 'btName'. if 'btName' is null, exec the current behavior tree specified by 'btsetcurrent'.
*/
func (a *Agent) BTGExec() EBTStatus {
	BTGLog.Trace("(a *Agent) BTGExec()")
	if a.IsActive {
		s := a.bTGExec_()
		for a.IsReferenceTree && s == BT_RUNNING {
			a.IsReferenceTree = false
			s = a.bTGExec_()
		}
		if a.IsMasked() {
			a.LogClientVariables(false)
		}
		return s
	}
	return BT_INVALID
}

func (a *Agent) bTGExec_() EBTStatus {
	BTGLog.Trace("(a *Agent) bTGExec_()")
	if a.bttCurrent != nil {
		bttOld := a.bttCurrent
		s := a.bttCurrent.Exec(a, BT_RUNNING)
		for s != BT_RUNNING {
			stackSize := len(a.bttStack)
			if stackSize != 0 {
				lastOne := a.bttStack[stackSize-1]
				a.bttStack = a.bttStack[:stackSize-1]
				a.bttCurrent = lastOne.btt
				var bExecCurrent bool
				if lastOne.tm == TM_Return {
					if !lastOne.triggerByEvent {
						if a.bttCurrent != bttOld {
							s = a.bttCurrent.Resume(a, s)
						}
					} else {
						bExecCurrent = true
					}
				} else {
					bExecCurrent = true
				}
				if bExecCurrent {
					bttOld = a.bttCurrent
					s = a.bttCurrent.Resume(a, s)
					break
				}
			} else {
				break
			}
		}
		return s
	} else {
		BTGLog.Warnf("NO ACTIVE BT!")
	}
	return BT_INVALID
}

func (a *Agent) LogClientVariables(bForce bool) {
	if BTGConfig.IsLogging || BTGConfig.IsSocketing {
		BTGLog.Tracef("Agent(Id:%v Name:%v) Client:%v", a.Id, a.Name, a.client)
	}
}

func (a *Agent) InvokeMethod(m *Method) []reflect.Value {
	BTGLog.Warnf("(a *Agent) InvokeMethod enter %v:%v", a.GetAgentName(), *m)
	if m == nil {
		return nil
	}
	mv := reflect.ValueOf(a.client).MethodByName(m.MethodName)
	if mv.IsValid() {
		mvt := mv.Type()
		if !mvt.IsVariadic() && len(m.MethodParams) != mvt.NumIn() {
			BTGLog.Warnf("Agent.InvokeMethod(%v.%v.%v) NumIn not fit,expect %v, but get %v ,%v", m.AgentClassName, m.AgentIntanceName, m.MethodName, mvt.NumIn(), len(m.MethodParams), m.MethodParams)
			return nil
		}

		params := make([]reflect.Value, mvt.NumIn(), mvt.NumIn())
		for i := 0; i < mvt.NumIn(); i++ {
			p := m.MethodParams[i]
			if p != nil {
				if p.IsRuntimeType {
					pv, err := p.GetRuntimeValue(mvt.In(i))
					if err == nil {
						params[i] = pv
					} else {
						BTGLog.Warnf("(a *Agent) InvokeMethod(m *Method) methodName=%v %vst param not fit", m.MethodName, (i + 1))
						return nil
					}
				} else {
					tmp := p.GetValue(a)
					if tmp.IsValid() && tmp.Type().ConvertibleTo(mvt.In(i)) {
						params[i] = tmp.Convert(mvt.In(i))
					} else {
						BTGLog.Warnf("(a *Agent) InvokeMethod(m *Method) methodName=%v %vst param not fit2", m.MethodName, (i + 1))
						return nil
					}
				}
			}
		}
		return mv.Call(params)
	} else {
		BTGLog.Warnf("Agent.InvokeMethod %v.%v not fount", a.GetAgentName(), m.MethodName)
	}
	return nil
}

func (a *Agent) SetVariable(variableName string, val interface{}) {
	p := AgentProperty_GetPropertyByVarId(a.GetAgentName(), MakeVariableId(variableName))
	if p != nil {
		BTGLog.Tracef("(a *Agent) SetVariable(%v,%v)", variableName, val)
		p.SetValue(a, reflect.ValueOf(val))
	} else {
		valMem := reflect.ValueOf(a.client).Elem().FieldByName(variableName)
		if valMem.IsValid() {
			valMem.Set(reflect.ValueOf(val))
		}
	}
}

func (a *Agent) GetVariableVal(variableName string) interface{} {
	val := a.GetVariable(variableName)
	if val.IsValid() {
		return val.Interface()
	}
	return nil
}

func (a *Agent) GetVariable(variableName string) reflect.Value {
	p := AgentProperty_GetPropertyByVarId(a.GetAgentName(), MakeVariableId(variableName))
	if p != nil {
		return p.GetValue(a)
	}
	return reflect.ValueOf(a.client).Elem().FieldByName(variableName)
}

func (a *Agent) BTGSetCurrent(relativePath string, triggerMode TriggerMode, bByEvent bool) {
	if relativePath == "" {
		return
	}

	err := BTGWorkspace.Load(relativePath, false)
	if err != nil {
		BTGLog.Warnf("%v is not a valid loaded behavior tree of %v", relativePath, a.GetAgentName())
		return
	} else {
		BTGWorkspace.RecordBTGAgentMapping(relativePath, a)
		if a.bttCurrent != nil {
			//if trigger mode is 'return', just push the current bt 'oldBt' on the stack and do nothing more
			//'oldBt' will be restored when the new triggered one ends
			if triggerMode == TM_Return {
				item := BehaviorTreeStackItem_t{
					btt:            a.bttCurrent,
					tm:             triggerMode,
					triggerByEvent: bByEvent,
				}
				a.bttStack = append(a.bttStack, item)
			} else if triggerMode == TM_Transfer {
				//don't use the bt stack to restore, we just abort the current one.
				//as the bt node has onenter/onexit, the abort can make them paired
				//BEHAVIAC_ASSERT(bttCurrent.GetName() != relativePath)
				a.bttCurrent.Abort(a)
				a.bttCurrent.Reset(a)
			}
		}
		var task *BehaviorTreeTask
		for _, btt := range a.behaviorTreeTasks {
			if btt.GetName() == relativePath {
				task = btt
				break
			}
		}
		var bRecursive bool
		if task != nil {
			for i := 0; i < len(a.bttStack); i++ {
				if a.bttStack[i].btt.GetName() == relativePath {
					bRecursive = true
					break
				}
			}
			task.Reset(a)
		}

		if task == nil || bRecursive {
			task = BTGWorkspace.CreateBehaviorTreeTask(relativePath)
			if task != nil {
				a.behaviorTreeTasks = append(a.behaviorTreeTasks, task)
			}
			a.bttCurrent = task
		}
	}
}

func (a *Agent) BTGOnEvent(event string) {
	if a.bttCurrent != nil {
		a.bttCurrent.OnEvent(a, event)
	}
}
