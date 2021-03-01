package behaviago

import (
	"errors"
	"reflect"
	"strconv"
)

const (
	BoolTypeStr       string = "bool"
	CharTypeStr              = "char"
	UbyteTypeStr             = "ubyte"
	SbyteTypeStr             = "sbyte"
	UshortTypeStr            = "ushort"
	ShortTypeStr             = "short"
	UintTypeStr              = "uint"
	IntTypeStr               = "int"
	UlongTypeStr             = "ulong"
	LongTypeStr              = "long"
	UllongTypeStr            = "ullong"
	LlongTypeStr             = "llong"
	FloatTypeStr             = "float"
	DoubleTypeStr            = "double"
	VoidTypeStr              = "void*"
	StringTypeStr            = "string"
	WStringTypeStr           = "behaviac::wstring"
	StdStringTypeStr         = "std::string"
	StdWStringTypeStr        = "std::wstring"
	//特殊
	RuntimeTypeStr = "runtime"
)

func editorRefType(gotype string) string {
	switch gotype {
	case "bool":
		return BoolTypeStr
	case "string":
		return StringTypeStr
	case "float32":
		return FloatTypeStr
	case "float64":
		return DoubleTypeStr
	case "int8":
		return SbyteTypeStr
	case "int16":
		return ShortTypeStr
	case "int32", "int":
		return IntTypeStr
	case "int64":
		return LlongTypeStr
	case "byte", "uint8":
		return UbyteTypeStr
	case "uint16":
		return UshortTypeStr
	case "uint32":
		return UintTypeStr
	case "uint64":
		return UllongTypeStr
	default:
		return gotype
	}
	return ""
}

// ============================================================================
type TypeCreator interface {
	Create() interface{}
}

type TypeCreatorWrapper func() interface{}

func (tcw TypeCreatorWrapper) Create() interface{} {
	return tcw()
}

// ============================================================================
var typeCreators = make(map[string]TypeCreator)

func RegisteTypeCreator(name string, creator TypeCreator) {
	typeCreators[name] = creator
}
func UnRegisteTypeCreator(name string) {
	delete(typeCreators, name)
}
func GetTypeCreator(name string) TypeCreator {
	return typeCreators[name]
}

// ============================================================================
type Property struct {
	VariableName        string
	VariableId          uint32
	RefParName          string
	RefParNameId        uint32
	MethodName          string
	InstanceName        string
	ValueStr            string
	TypeName            string
	Value               reflect.Value
	IsValidDefaultValue bool
	IsConst             bool
	IsStatic            bool
	IsLocal             bool
	IsMember            bool
	IsRuntimeType       bool
	IdxProperty         *Property
}

func NewProperty(typeName, instanceName, propertyName, valueStr string) *Property {
	p := &Property{
		TypeName:     typeName,
		InstanceName: instanceName,
		VariableName: propertyName,
		ValueStr:     valueStr,
	}
	if valueStr != "" {
		p.Value = p.convertStr2Value()
	}
	return p
}

func NewPropertyConst(typeName, valueStr string) *Property {
	p := &Property{
		TypeName: typeName,
		ValueStr: valueStr,
	}
	p.IsConst = true
	p.IsRuntimeType = typeName == RuntimeTypeStr
	if valueStr != "" {
		p.Value = p.convertStr2Value()
	}
	return p
}

func NewPropertyStatic(typeName, fullName, arrayIndexStr string, bIsStatic bool) *Property {
	instanceName, agentType, propertyName := ParseInstanceNameProperty(fullName)
	pProperty := AgentProperty_GetPropertyByVarName(agentType, propertyName)
	if pProperty != nil {
		if pProperty.InstanceName != instanceName {
			pProperty = pProperty.Clone()
			pProperty.InstanceName = instanceName
			AgentProperty_AddInstance(agentType, pProperty)
		}
	} else {
		pProperty = AgentProperty_AddLocal(agentType, typeName, propertyName, "")
	}

	if len(arrayIndexStr) != 0 {
		pProperty.IdxProperty = ParseProperty(arrayIndexStr)
	}
	pProperty.IsStatic = bIsStatic
	return pProperty
}

func (p *Property) Clone() *Property {
	ret := &Property{
		VariableName:        p.VariableName,
		VariableId:          p.VariableId,
		RefParName:          p.RefParName,
		RefParNameId:        p.RefParNameId,
		MethodName:          p.MethodName,
		InstanceName:        p.InstanceName,
		ValueStr:            p.ValueStr,
		TypeName:            p.TypeName,
		IsValidDefaultValue: p.IsValidDefaultValue,
		IsConst:             p.IsConst,
		IsStatic:            p.IsStatic,
		IsLocal:             p.IsLocal,
		IsMember:            p.IsMember,
		IsRuntimeType:       p.IsRuntimeType,
		IdxProperty:         p.IdxProperty.Clone(),
	}
	ret.Value = reflect.New(p.Value.Type())
	ret.Value.Set(p.Value)
	return ret
}

func (p *Property) GetParentAgent(a *Agent) *Agent {
	return GetAgentInstance(a, p.InstanceName)
}

func (p *Property) SetFromByMethod(agentFrom *Agent, methodFrom *Method, agentTo *Agent) bool {
	if agentFrom == nil || methodFrom == nil || agentTo == nil {
		return false
	}

	vals := agentFrom.InvokeMethod(methodFrom)
	if len(vals) != 0 {
		return p.SetValue(agentTo, vals[0])
	}
	return false
}

func (p *Property) SetFromByProperty(agentFrom *Agent, propertyFrom *Property, agentTo *Agent) bool {
	if agentFrom == nil || propertyFrom == nil || agentTo == nil {
		return false
	}
	v := propertyFrom.GetValue(agentFrom)
	if v.IsValid() {
		return p.SetValue(agentTo, v)
	}
	return false
}

func (p *Property) SetFromByMethodByIndex(agentFrom *Agent, methodFrom *Method, agentTo *Agent, index int) bool {
	if agentFrom == nil || methodFrom == nil || agentTo == nil {
		return false
	}

	vals := agentFrom.InvokeMethod(methodFrom)
	if len(vals) != 0 {
		val := vals[0]
		kind := val.Type().Kind()
		if kind == reflect.Array || kind == reflect.Slice || kind == reflect.String {
			if index >= 0 && index < val.Len() {
				return p.SetValue(agentTo, val.Index(index))
			}
		}
	}
	return false
}

func (p *Property) SetValue(a *Agent, v reflect.Value) bool {
	BTGLog.Warnf("(p *Property) SetValue enter %v", *p)
	if a == nil || !v.IsValid() {
		return false
	}
	vv := p.GetValue(a)
	if vv.IsValid() && v.Type().ConvertibleTo(vv.Type()) {
		if p.IsConst {
			BTGLog.Warnf("(p *Property) SetValue found %v is const", p.VariableName)
			return false
		}
		if vv.CanSet() {
			vv.Set(v)
		} else {
			p.Value = v
		}
		return true
	}

	return false
}

func (p *Property) GetValue(a *Agent) reflect.Value {
	BTGLog.Tracef("(p *Property) GetValue enter ValueStr=%v, Value=%v", p.ValueStr, p.Value)
	if p.IsMember {
		v := reflect.ValueOf(a.client).Elem().FieldByName(p.VariableName)
		if v.IsValid() {
			if p.IdxProperty != nil {
				idx := p.IdxProperty.GetIndex(a)
				kind := v.Kind()
				if kind == reflect.String || kind == reflect.Array || kind == reflect.Slice {
					if idx >= 0 && idx < v.Len() {
						return v.Index(idx)
					}
				}
			} else {
				return v
			}
		}
	} else {
		return p.Value
	}
	return reflect.Value{}
}

func (p *Property) GetFloat(a *Agent) float64 {
	v := p.GetValue(a)
	if v.IsValid() {
		k := v.Kind()
		if k == reflect.Float32 || k == reflect.Float64 {
			return v.Float()
		}
	}
	return .0
}

func (p *Property) GetInt(a *Agent) int64 {
	v := p.GetValue(a)
	if v.IsValid() {
		k := v.Kind()
		if k == reflect.Int || k == reflect.Int8 || k == reflect.Int16 || k == reflect.Int32 || k == reflect.Int64 {
			return v.Int()
		}
	}
	return 0
}

func (p *Property) GetUint(a *Agent) uint64 {
	v := p.GetValue(a)
	if v.IsValid() {
		k := v.Kind()
		if k == reflect.Uint || k == reflect.Uint8 || k == reflect.Uint16 || k == reflect.Uint32 || k == reflect.Uint64 {
			return v.Uint()
		}
	}
	return 0
}

func (p *Property) GetStr(a *Agent) string {
	v := p.GetValue(a)
	if v.IsValid() {
		if v.Kind() == reflect.String {
			return v.String()
		}
	}
	return ""
}

func (p *Property) GetIndex(a *Agent) int {
	if p.IdxProperty == nil {
		return 0
	}

	pa := p.GetParentAgent(a)
	v := p.GetValue(pa)
	if v.IsValid() && v.Type().ConvertibleTo(IntType) {
		return int(v.Int())
	}
	return 0
}

func (p *Property) GetRuntimeValue(t reflect.Type) (val reflect.Value, err error) {
	BTGLog.Tracef("(p *Property) GetRuntimeValue(%v)", t.Name())
	if p.TypeName != RuntimeTypeStr {
		return reflect.Value{}, errors.New("not runtime type")
	}
	switch t.Kind() {
	case reflect.Bool:
		var b bool
		b, err = strconv.ParseBool(p.ValueStr)
		if err == nil {
			val = reflect.ValueOf(b)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var i int64
		i, err = strconv.ParseInt(p.ValueStr, 10, 64)
		if err == nil {
			val = reflect.ValueOf(i)
			val = val.Convert(t)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var u uint64
		u, err = strconv.ParseUint(p.ValueStr, 10, 64)
		if err == nil {
			val = reflect.ValueOf(u)
			val = val.Convert(t)
		}
	case reflect.Float32, reflect.Float64:
		var f float64
		f, err = strconv.ParseFloat(p.ValueStr, 64)
		if err == nil {
			val = reflect.ValueOf(f)
			val = val.Convert(t)
		}
	case reflect.String:
		val = reflect.ValueOf(p.ValueStr)
	default:
		err = errors.New("unsupport type " + t.Name())
	}
	return
}

func (p *Property) convertStr2Value() reflect.Value {
	switch p.TypeName {
	case BoolTypeStr:
		val, err := strconv.ParseBool(p.ValueStr)
		if err == nil {
			return reflect.ValueOf(val)
		}
	case CharTypeStr:
		if len(p.ValueStr) != 0 {
			return reflect.ValueOf(p.ValueStr[0])
		}
	case UbyteTypeStr:
		val, err := strconv.ParseUint(p.ValueStr, 10, 0)
		if err == nil {
			return reflect.ValueOf(uint8(val))
		}
	case SbyteTypeStr:
		val, err := strconv.ParseInt(p.ValueStr, 10, 0)
		if err == nil {
			return reflect.ValueOf(int8(val))
		}
	case UshortTypeStr:
		val, err := strconv.ParseUint(p.ValueStr, 10, 0)
		if err == nil {
			return reflect.ValueOf(uint16(val))
		}
	case ShortTypeStr:
		val, err := strconv.ParseInt(p.ValueStr, 10, 0)
		if err == nil {
			return reflect.ValueOf(int16(val))
		}
	case UintTypeStr:
		val, err := strconv.ParseUint(p.ValueStr, 10, 0)
		if err == nil {
			return reflect.ValueOf(uint32(val))
		}
	case IntTypeStr:
		val, err := strconv.ParseInt(p.ValueStr, 10, 0)
		if err == nil {
			return reflect.ValueOf(int32(val))
		}
	case UlongTypeStr, UllongTypeStr:
		val, err := strconv.ParseUint(p.ValueStr, 10, 64)
		if err == nil {
			return reflect.ValueOf(val)
		}
	case LongTypeStr, LlongTypeStr:
		val, err := strconv.ParseInt(p.ValueStr, 10, 64)
		if err == nil {
			return reflect.ValueOf(val)
		}
	case FloatTypeStr:
		val, err := strconv.ParseFloat(p.ValueStr, 32)
		if err == nil {
			return reflect.ValueOf(val)
		}
	case DoubleTypeStr:
		val, err := strconv.ParseFloat(p.ValueStr, 64)
		if err == nil {
			return reflect.ValueOf(val)
		}
	case StringTypeStr, WStringTypeStr, StdStringTypeStr, StdWStringTypeStr, RuntimeTypeStr:
		return reflect.ValueOf(p.ValueStr)
	default: //自定义类型
		BTGLog.Warn("(p *Property) convertStr2Value unsupport ", p.TypeName)
	}
	return reflect.Value{}
}
