package behaviago

import (
	"reflect"
	"strings"
)

// ============================================================================
var dummyCompareable Compareable
var CompareableInterfaceType = reflect.TypeOf(dummyCompareable)

type Compareable interface {
	Equal(interface{}) bool
	Greater(interface{}) bool
	GreaterEqual(interface{}) bool
}

// ============================================================================
type E_VariableComparisonType int

const (
	VariableComparisonType_Assignment   E_VariableComparisonType = iota //( "Assignment (=)" )
	VariableComparisonType_Equal                                        //( "Equal (==)" )
	VariableComparisonType_NotEqual                                     //( "Not Equal (!=)" )
	VariableComparisonType_Greater                                      //( "Greater (>)"  )
	VariableComparisonType_GreaterEqual                                 //( "Greater Or Equal (>=)" )
	VariableComparisonType_Less                                         //( "Lower (<)"  )
	VariableComparisonType_LessEqual                                    //( "Lower Or Equal (<=)" )
	VariableComparisonType_And                                          //( "Lower Or Equal (&&)" )
	VariableComparisonType_Or                                           //( "Lower Or Equal (||)" )
)

func ParseComparisonType(comparionOperator string) E_VariableComparisonType {
	switch comparionOperator {
	case "Assignment":
		return VariableComparisonType_Assignment
	case "Equal":
		return VariableComparisonType_Equal
	case "NotEqual":
		return VariableComparisonType_NotEqual
	case "Greater":
		return VariableComparisonType_Greater
	case "GreaterEqual":
		return VariableComparisonType_GreaterEqual
	case "Less":
		return VariableComparisonType_Less
	case "LessEqual":
		return VariableComparisonType_LessEqual
	default:
		BTGLog.Warnf("cannot ParseComparisonType(%v)", comparionOperator)
	}
	return VariableComparisonType_Equal
}

// ============================================================================
func EvaluteCompare(a *Agent, opl *Property, oplm *Method, comparatorType E_VariableComparisonType, opr *Property, oprm *Method) bool {
	var lhs reflect.Value
	if opl != nil {
		p := opl.GetParentAgent(a)
		if p != nil {
			lhs = opl.GetValue(p)
		}
	} else {
		p := oplm.GetParentAgent(a)
		if p != nil {
			vals := p.InvokeMethod(oplm)
			if len(vals) != 0 {
				lhs = vals[0]
			}
		}
	}
	var rhs reflect.Value
	if opr != nil {
		p := opr.GetParentAgent(a)
		if p != nil {
			rhs = opr.GetValue(p)
		}
	} else {
		p := oprm.GetParentAgent(a)
		if p != nil {
			vals := p.InvokeMethod(oprm)
			if len(vals) != 0 {
				rhs = vals[0]
			}
		}
	}
	switch comparatorType {
	case VariableComparisonType_Assignment:
		if opl != nil && oplm == nil {
			p := opl.GetParentAgent(a)
			if p != nil {
				return opl.SetValue(p, rhs)
			}
		}
	case VariableComparisonType_Equal:
		return OperatorEqual(lhs, rhs)
	case VariableComparisonType_NotEqual:
		return !OperatorEqual(lhs, rhs)
	case VariableComparisonType_Greater:
		return OperatorGreater(lhs, rhs)
	case VariableComparisonType_GreaterEqual:
		return OperatorGreaterEqual(lhs, rhs)
	case VariableComparisonType_Less:
		return !OperatorGreaterEqual(lhs, rhs)
	case VariableComparisonType_LessEqual:
		return !OperatorGreater(lhs, rhs)
	}
	return false
}

// ============================================================================
func OperatorEqual(lhs, rhs reflect.Value) bool {
	if !lhs.IsValid() || !rhs.IsValid() {
		return false
	}
	lhsKind := lhs.Type().Kind()
	rhsKind := rhs.Type().Kind()
	if lhsKind != rhsKind {
		BTGLog.Warnf("OperatorEqual type fit warning!!! lhs kind of type is %v , rhs kind of type is %v", lhsKind, rhsKind)
	}
	switch lhsKind {
	case reflect.Bool:
		if lhsKind == rhsKind {
			return lhs.Bool() == rhs.Bool()
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if rhs.Type().ConvertibleTo(lhs.Type()) {
			return lhs.Int() == rhs.Int()
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if rhs.Type().ConvertibleTo(lhs.Type()) {
			return lhs.Uint() == rhs.Uint()
		}
	case reflect.Float32, reflect.Float64:
		if rhs.Type().ConvertibleTo(lhs.Type()) {
			return lhs.Float() == rhs.Float()
		}
	case reflect.Array, reflect.Slice, reflect.String:
		if lhsKind == rhsKind {
			if lhs.Len() != rhs.Len() {
				return false
			}
			for i := 0; i < lhs.Len(); i++ {
				if !OperatorEqual(lhs.Index(i), rhs.Index(i)) {
					return false
				}
			}
			return true
		}
	case reflect.Map:
		if lhsKind == rhsKind {
			if lhs.Len() != rhs.Len() {
				return false
			}
			keys := lhs.MapKeys()
			for i := 0; i < len(keys); i++ {
				vl := lhs.MapIndex(keys[i])
				vr := rhs.MapIndex(keys[i])
				if !vl.IsValid() || vr.IsValid() {
					return false
				}
				if !OperatorEqual(vl, vr) {
					return false
				}
			}
			return true
		}
	case reflect.Ptr, reflect.Interface:
		if lhsKind == rhsKind {
			lhs = lhs.Elem()
			rhs = rhs.Elem()
			return OperatorEqual(lhs, rhs)
		}
	case reflect.Struct:
		if lhsKind == rhsKind {
			if c, ok := lhs.Interface().(Compareable); ok {
				return c.Equal(rhs.Interface())
			}
		}
	case reflect.Chan, reflect.UnsafePointer, reflect.Func:
		return lhs.Pointer() == rhs.Pointer()
	}
	return false
}

// ============================================================================
func OperatorGreater(lhs, rhs reflect.Value) bool {
	if !lhs.IsValid() || !rhs.IsValid() {
		return false
	}
	lhsKind := lhs.Type().Kind()
	rhsKind := rhs.Type().Kind()
	if lhsKind != rhsKind {
		BTGLog.Warnf("OperatorGreater type fit warning!!! lhs kind of type is %v , rhs kind of type is %v", lhsKind, rhsKind)
	}
	switch lhsKind {
	case reflect.Bool:
		return false
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if rhs.Type().ConvertibleTo(lhs.Type()) {
			return lhs.Int() > rhs.Int()
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if rhs.Type().ConvertibleTo(lhs.Type()) {
			return lhs.Uint() > rhs.Uint()
		}
	case reflect.Float32, reflect.Float64:
		if rhs.Type().ConvertibleTo(lhs.Type()) {
			return lhs.Float() > rhs.Float()
		}
	case reflect.Array, reflect.Slice, reflect.String:
		if lhsKind == rhsKind {
			var min int
			if lhs.Len() > rhs.Len() {
				min = rhs.Len()
			} else {
				min = lhs.Len()
			}
			for i := 0; i < min; i++ {
				if OperatorGreater(lhs.Index(i), rhs.Index(i)) {
					return true
				} else {
					return false
				}
			}
			return lhs.Len() > rhs.Len()
		}
	case reflect.Map:
		if lhsKind == rhsKind {
			if lhs.Len() != rhs.Len() {
				return false
			}
			keys := lhs.MapKeys()
			for i := 0; i < len(keys); i++ {
				vl := lhs.MapIndex(keys[i])
				vr := rhs.MapIndex(keys[i])
				if !vl.IsValid() || vr.IsValid() {
					return false
				}
				if OperatorGreater(vl, vr) {
					return true
				} else {
					return false
				}
			}
			return true
		}
	case reflect.Ptr, reflect.Interface:
		if lhsKind == rhsKind {
			lhs = lhs.Elem()
			rhs = rhs.Elem()
			return OperatorGreater(lhs, rhs)
		}
	case reflect.Struct:
		if lhsKind == rhsKind {
			if c, ok := lhs.Interface().(Compareable); ok {
				return c.Greater(rhs.Interface())
			}
		}
	}
	return false
}

// ============================================================================
func OperatorGreaterEqual(lhs, rhs reflect.Value) bool {
	if !lhs.IsValid() || !rhs.IsValid() {
		return false
	}
	lhsKind := lhs.Type().Kind()
	rhsKind := rhs.Type().Kind()
	if lhsKind != rhsKind {
		BTGLog.Warnf("OperatorGreaterEqual type fit warning!!! lhs kind of type is %v , rhs kind of type is %v", lhsKind, rhsKind)
	}
	switch lhsKind {
	case reflect.Bool:
		return false
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if rhs.Type().ConvertibleTo(lhs.Type()) {
			return lhs.Int() >= rhs.Int()
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if rhs.Type().ConvertibleTo(lhs.Type()) {
			return lhs.Uint() >= rhs.Uint()
		}
	case reflect.Float32, reflect.Float64:
		if rhs.Type().ConvertibleTo(lhs.Type()) {
			return lhs.Float() >= rhs.Float()
		}
	case reflect.Array, reflect.Slice, reflect.String:
		if lhsKind == rhsKind {
			var min int
			if lhs.Len() > rhs.Len() {
				min = rhs.Len()
			} else {
				min = lhs.Len()
			}
			for i := 0; i < min; i++ {
				if OperatorGreaterEqual(lhs.Index(i), rhs.Index(i)) {
					return true
				} else {
					return false
				}
			}
			return lhs.Len() >= rhs.Len()
		}
	case reflect.Map:
		if lhsKind == rhsKind {
			if lhs.Len() != rhs.Len() {
				return false
			}
			keys := lhs.MapKeys()
			for i := 0; i < len(keys); i++ {
				vl := lhs.MapIndex(keys[i])
				vr := rhs.MapIndex(keys[i])
				if !vl.IsValid() || vr.IsValid() {
					return false
				}
				if OperatorGreaterEqual(vl, vr) {
					return true
				} else {
					return false
				}
			}
			return true
		}
	case reflect.Ptr, reflect.Interface:
		if lhsKind == rhsKind {
			lhs = lhs.Elem()
			rhs = rhs.Elem()
			return OperatorGreaterEqual(lhs, rhs)
		}
	case reflect.Struct:
		if lhsKind == rhsKind {
			if c, ok := lhs.Interface().(Compareable); ok {
				return c.GreaterEqual(rhs.Interface())
			}
		}
	}
	return false
}

/**
  Condition node can compare the value of right and left. return Failure or Success
*/
const (
	ConditionNodeName = "Condition"
)

func init() {
	RegisteNodeCreator(ConditionNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		c := newConditionNode()
		return c
	}))
}

func newConditionNode() *Condition {
	c := &Condition{}
	c.SetClassName(ConditionNodeName)
	c.SetSelf(c)
	return c
}

// ============================================================================
type Condition struct {
	BehaviorNodeBase
	opl            *Property
	oplm           *Method
	opr            *Property
	oprm           *Method
	comparatorType E_VariableComparisonType
}

func (c *Condition) Load(version int, agentType string, properties []property_t) {
	c.BehaviorNodeBase.Load(version, agentType, properties)
	var comparatorName string
	for i := 0; i < len(properties); i++ {
		switch properties[i].name {
		case "Operator":
			comparatorName = properties[i].value
		case "Opl":
			if !strings.ContainsAny(properties[i].value, "(") {
				c.opl = ParseProperty(properties[i].value)
			} else {
				c.oplm = LoadMethod(properties[i].value)
			}
		case "Opr":
			if !strings.ContainsAny(properties[i].value, "(") {
				c.opr = ParseProperty(properties[i].value)
			} else {
				c.oprm = LoadMethod(properties[i].value)
			}
		}
	}
	c.comparatorType = ParseComparisonType(comparatorName)
}

func (c *Condition) CreateTask() BehaviorTask {
	BTGLog.Trace("(c *Condition) CreateTask()")
	return &ConditionTask{}
}

func (c *Condition) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*Condition); !ok {
		return false
	}
	return true
}

func (c *Condition) Evaluate(a *Agent, result EBTStatus) bool {
	BTGLog.Trace("(c *Condition) Evaluate() enter")
	return EvaluteCompare(a, c.opl, c.oplm, c.comparatorType, c.opr, c.oprm)
}

// ============================================================================
type ConditionTask struct {
	LeafTask
}

func (ct *ConditionTask) OnEnter(a *Agent) bool {
	BTGLog.Trace("(ct *ConditionTask) OnEnter() enter")
	return true
}

func (ct *ConditionTask) OnExit(a *Agent, status EBTStatus) {
	BTGLog.Trace("(ct *ConditionTask) OnExit() enter")
}

func (ct *ConditionTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	BTGLog.Trace("(ct *ConditionTask) Update() enter")
	if c, ok := ct.GetNode().(*Condition); ok {
		if c.Evaluate(a, BT_INVALID) {
			return BT_SUCCESS
		} else {
			return BT_FAILURE
		}
	}
	return BT_FAILURE
}
