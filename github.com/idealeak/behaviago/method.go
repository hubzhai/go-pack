package behaviago

import "reflect"

type Method struct {
	AgentIntanceName string
	AgentClassName   string
	MethodName       string
	MethodParams     []*Property
}

func (m *Method) ParseParam(params []string) {
	paramsCnt := len(params)
	if paramsCnt != 0 {
		m.MethodParams = make([]*Property, paramsCnt, paramsCnt)
		for i := 0; i < paramsCnt; i++ {
			m.MethodParams[i] = ParseProperty(params[i])
		}
	}
}

func (m *Method) GetParentAgent(a *Agent) *Agent {
	return GetAgentInstance(a, m.AgentIntanceName)
}

func (m *Method) Invoke(a *Agent) []reflect.Value {
	p := m.GetParentAgent(a)
	if p != nil {
		return p.InvokeMethod(m)
	}
	return nil
}

func (m *Method) InvokeAndGetFloat(a *Agent) float64 {
	vals := m.Invoke(a)
	if len(vals) != 0 {
		v := vals[0]
		k := v.Kind()
		if k == reflect.Float32 || k == reflect.Float64 {
			return v.Float()
		}
	}
	return .0
}

func (m *Method) InvokeAndGetInt(a *Agent) int64 {
	vals := m.Invoke(a)
	if len(vals) != 0 {
		v := vals[0]
		k := v.Kind()
		if k == reflect.Int || k == reflect.Int8 || k == reflect.Int16 || k == reflect.Int32 || k == reflect.Int64 {
			return v.Int()
		}
	}
	return 0
}

func (m *Method) InvokeAndGetUint(a *Agent) uint64 {
	vals := m.Invoke(a)
	if len(vals) != 0 {
		v := vals[0]
		k := v.Kind()
		if k == reflect.Uint || k == reflect.Uint8 || k == reflect.Uint16 || k == reflect.Uint32 || k == reflect.Uint64 {
			return v.Uint()
		}
	}
	return 0
}

func (m *Method) InvokeAndGetStr(a *Agent) string {
	vals := m.Invoke(a)
	if len(vals) != 0 {
		v := vals[0]
		if v.Kind() == reflect.String {
			return v.String()
		}
	}
	return ""
}
