package behaviago

import (
	"reflect"
	"strings"
)

const (
	ECO_INVALID EComputeOperator = iota
	ECO_ADD
	ECO_SUB
	ECO_MUL
	ECO_DIV
)

var IntType = reflect.ValueOf(1).Type()

///Compute
/**
  Compute the result of Operand1 and Operand2 and assign it to the Left Operand.
  Compute node can perform Add, Sub, Mul and Div operations. a left and right Operand
  can be a agent property or a par value.
*/
const (
	ComputeNodeName = "Compute"
)

func init() {
	RegisteNodeCreator(ComputeNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		c := &Compute{}
		return c
	}))
}

func newCompute() *Compute {
	c := &Compute{}
	c.SetClassName(ComputeNodeName)
	c.SetSelf(c)
	return c
}

// ============================================================================
var dummyComputeable Computeable
var ComputeableInterfaceType = reflect.TypeOf(dummyComputeable)

type Computeable interface {
	Add(r interface{}) interface{}
	Sub(r interface{}) interface{}
	Mul(r interface{}) interface{}
	Div(r interface{}) interface{}
}

// ============================================================================
type EComputeOperator int

// ============================================================================
type Compute struct {
	BehaviorNodeBase
	opl      *Property
	opr1     *Property
	opr1m    *Method
	opr2     *Property
	opr2m    *Method
	operator EComputeOperator
}

func (c *Compute) Load(version int, agentType string, properties []property_t) {
	c.BehaviorNodeBase.Load(version, agentType, properties)
	for i := 0; i < len(properties); i++ {
		switch properties[i].name {
		case "Opl":
			c.opl = ParseProperty(properties[i].value)
		case "Operator":
			switch properties[i].value {
			case "Add":
				c.operator = ECO_ADD
			case "Sub":
				c.operator = ECO_SUB
			case "Mul":
				c.operator = ECO_MUL
			case "Div":
				c.operator = ECO_DIV
			}
		case "Opr1":
			if !strings.ContainsAny(properties[i].value, "(") {
				c.opr1 = ParseProperty(properties[i].value)
			} else {
				c.opr1m = LoadMethod(properties[i].value)
			}
		case "Opr2":
			if !strings.ContainsAny(properties[i].value, "(") {
				c.opr2 = ParseProperty(properties[i].value)
			} else {
				c.opr2m = LoadMethod(properties[i].value)
			}
		}
	}
}

func (c *Compute) CreateTask() BehaviorTask {
	BTGLog.Trace("(c *Compute) CreateTask()")
	return &ComputeTask{}
}

func (c *Compute) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*Compute); !ok {
		return false
	}
	return true
}

// ============================================================================
type ComputeTask struct {
	LeafTask
}

func (c *ComputeTask) OnEnter(a *Agent) bool {
	return c.LeafTask.OnEnter(a)
}

func (c *ComputeTask) OnExit(a *Agent, status EBTStatus) {
	c.LeafTask.OnExit(a, status)
}

func (c *ComputeTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	if node, ok := c.GetNode().(*Compute); ok {
		if EvaluteCompute(a, node.opl, node.opr1, node.opr1m, node.operator, node.opr2, node.opr2m) {
			return BT_SUCCESS
		} else {
			return node.UpdateImpl(a, childStatus)
		}
	}
	return BT_SUCCESS
}

// ============================================================================
func EvaluteCompute(a *Agent, opl, opr1 *Property, opr1m *Method, operator EComputeOperator, opr2 *Property, opr2m *Method) bool {
	al := opl.GetParentAgent(a)
	var valr1 reflect.Value
	if opr1 != nil {
		valr1 = opr1.GetValue(a)
	} else {
		p := opr1m.GetParentAgent(a)
		if p != nil {
			vals := p.InvokeMethod(opr1m)
			if len(vals) != 0 {
				valr1 = vals[0]
			}
		}
	}
	var valr2 reflect.Value
	if opr2 != nil {
		valr2 = opr2.GetValue(a)
	} else {
		p := opr2m.GetParentAgent(a)
		if p != nil {
			vals := p.InvokeMethod(opr2m)
			if len(vals) != 0 {
				valr2 = vals[0]
			}
		}
	}
	var result interface{}
	switch operator {
	case ECO_ADD:
		result = EvaluteAdd(valr1, valr2)
	case ECO_SUB:
		result = EvaluteSub(valr1, valr2)
	case ECO_MUL:
		result = EvaluteMul(valr1, valr2)
	case ECO_DIV:
		result = EvaluteDiv(valr1, valr2)
	}
	return opl.SetValue(al, reflect.ValueOf(result))
}

// ============================================================================
func EvaluteAdd(lhs, rhs reflect.Value) interface{} {
	if !lhs.IsValid() {
		return rhs.Interface()
	}
	if !rhs.IsValid() {
		return lhs.Interface()
	}
	lhsKind := lhs.Type().Kind()
	rhsKind := rhs.Type().Kind()
	switch lhsKind {
	case reflect.Bool:
		if lhsKind == rhsKind {
			return lhs.Bool() || rhs.Bool()
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if rhs.Type().ConvertibleTo(lhs.Type()) {
			return lhs.Int() + rhs.Int()
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if rhs.Type().ConvertibleTo(lhs.Type()) {
			return lhs.Uint() + rhs.Uint()
		}
	case reflect.Float32, reflect.Float64:
		if rhs.Type().ConvertibleTo(lhs.Type()) {
			return lhs.Float() + rhs.Float()
		}
	case reflect.String:
		if lhsKind == rhsKind {
			return lhs.String() + rhs.String()
		}
	case reflect.Slice, reflect.Array:
		if lhsKind == rhsKind {
			ret := reflect.MakeSlice(lhs.Type(), lhs.Len()+rhs.Len(), 0)
			ret = reflect.AppendSlice(ret, lhs)
			ret = reflect.AppendSlice(ret, rhs)
			return ret.Interface()
		}
	case reflect.Map:
		if lhsKind == rhsKind {
			ret := reflect.MakeMap(lhs.Type())
			keys := lhs.MapKeys()
			for i := 0; i < len(keys); i++ {
				ret.SetMapIndex(keys[i], lhs.MapIndex(keys[i]))
			}
			keys = rhs.MapKeys()
			for i := 0; i < len(keys); i++ {
				ret.SetMapIndex(keys[i], rhs.MapIndex(keys[i]))
			}
			return ret.Interface()
		}
	case reflect.Ptr, reflect.Interface:
		if lhsKind == rhsKind {
			lhs = lhs.Elem()
			rhs = rhs.Elem()
			ret := EvaluteAdd(lhs, rhs)
			return &ret
		}
	case reflect.Struct:
		if lhsKind == rhsKind {
			if c, ok := lhs.Interface().(Computeable); ok {
				return c.Add(rhs.Interface())
			}
		}
	}
	return nil
}

// ============================================================================
func EvaluteSub(lhs, rhs reflect.Value) interface{} {
	if !lhs.IsValid() {
		return nil
	}
	if !rhs.IsValid() {
		return lhs.Interface()
	}
	lhsKind := lhs.Type().Kind()
	rhsKind := rhs.Type().Kind()
	switch lhsKind {
	case reflect.Bool:
		if lhsKind == rhsKind {
			return lhs.Bool() && rhs.Bool()
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if rhs.Type().ConvertibleTo(lhs.Type()) {
			return lhs.Int() - rhs.Int()
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if rhs.Type().ConvertibleTo(lhs.Type()) {
			return lhs.Uint() - rhs.Uint()
		}
	case reflect.Float32, reflect.Float64:
		if rhs.Type().ConvertibleTo(lhs.Type()) {
			return lhs.Float() - rhs.Float()
		}
	case reflect.Array, reflect.Slice, reflect.String:
		BTGLog.Warnf("EvaluteSub unsupport %v", lhsKind)
		return nil
	case reflect.Map:
		if lhsKind == rhsKind {
			ret := reflect.MakeMap(lhs.Type())
			keys := lhs.MapKeys()
			for i := 0; i < len(keys); i++ {
				if rhs.MapIndex(keys[i]).IsValid() {
					ret.SetMapIndex(keys[i], lhs.MapIndex(keys[i]))
				}
			}
			return ret.Interface()
		}
	case reflect.Ptr, reflect.Interface:
		if lhsKind == rhsKind {
			lhs = lhs.Elem()
			rhs = rhs.Elem()
			ret := EvaluteSub(lhs, rhs)
			return &ret
		}
	case reflect.Struct:
		if lhsKind == rhsKind {
			if c, ok := lhs.Interface().(Computeable); ok {
				return c.Sub(rhs.Interface())
			}
		}
	}
	return nil
}

// ============================================================================
func EvaluteMul(lhs, rhs reflect.Value) interface{} {
	if !lhs.IsValid() {
		return nil
	}
	if !rhs.IsValid() {
		return nil
	}
	lhsKind := lhs.Type().Kind()
	rhsKind := rhs.Type().Kind()
	switch lhsKind {
	case reflect.Bool:
		if lhsKind == rhsKind {
			return lhs.Bool() && rhs.Bool()
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if rhs.Type().ConvertibleTo(lhs.Type()) {
			return lhs.Int() * rhs.Int()
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if rhs.Type().ConvertibleTo(lhs.Type()) {
			return lhs.Uint() * rhs.Uint()
		}
	case reflect.Float32, reflect.Float64:
		if rhs.Type().ConvertibleTo(lhs.Type()) {
			return lhs.Float() * rhs.Float()
		}
	case reflect.Array, reflect.Slice, reflect.String, reflect.Map:
		BTGLog.Warnf("EvaluteMul unsupport %v", lhsKind)
		return nil
	case reflect.Ptr, reflect.Interface:
		if lhsKind == rhsKind {
			lhs = lhs.Elem()
			rhs = rhs.Elem()
			ret := EvaluteMul(lhs, rhs)
			return &ret
		}
	case reflect.Struct:
		if lhsKind == rhsKind {
			if c, ok := lhs.Interface().(Computeable); ok {
				return c.Mul(rhs.Interface())
			}
		}
	}
	return nil
}

// ============================================================================
func EvaluteDiv(lhs, rhs reflect.Value) interface{} {
	if !lhs.IsValid() {
		return nil
	}
	if !rhs.IsValid() {
		return nil
	}
	lhsKind := lhs.Type().Kind()
	rhsKind := rhs.Type().Kind()
	switch lhsKind {
	case reflect.Bool:
		if lhsKind == rhsKind {
			return lhs.Bool() && rhs.Bool()
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if rhs.Type().ConvertibleTo(lhs.Type()) {
			return lhs.Int() / rhs.Int()
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if rhs.Type().ConvertibleTo(lhs.Type()) {
			return lhs.Uint() / rhs.Uint()
		}
	case reflect.Float32, reflect.Float64:
		if rhs.Type().ConvertibleTo(lhs.Type()) {
			return lhs.Float() / rhs.Float()
		}
	case reflect.Array, reflect.Slice, reflect.String, reflect.Map:
		BTGLog.Warnf("EvaluteDiv unsupport %v", lhsKind)
		return nil
	case reflect.Ptr, reflect.Interface:
		if lhsKind == rhsKind {
			lhs = lhs.Elem()
			rhs = rhs.Elem()
			ret := EvaluteDiv(lhs, rhs)
			return &ret
		}
	case reflect.Struct:
		if lhsKind == rhsKind {
			if c, ok := lhs.Interface().(Computeable); ok {
				return c.Div(rhs.Interface())
			}
		}
	}
	return nil
}
