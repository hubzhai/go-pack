package behaviago

import (
	"reflect"
	"strings"
)

// ============================================================================
func CheckResult(rets []reflect.Value) EBTStatus {
	if len(rets) == 0 {
		return BT_SUCCESS
	}
	for _, v := range rets {
		if v.Type().Name() == "EBTStatus" {
			return v.Interface().(EBTStatus)
		}
	}
	return BT_SUCCESS
}

// ============================================================================
func LoadMethod(val string) *Method {
	var paramPart string
	m := &Method{}
	m.AgentIntanceName, m.AgentClassName, m.MethodName, paramPart = ParseMethodNames(val)
	m.ParseParam(ParseForParams(paramPart))
	return m
}

// ============================================================================
//<property Method="Self.Enemy::setSpeed(float Self.Enemy::initSpeed)" />
func ParseMethodNames(fullName string) (agentIntanceName, agentClassName, methodName, paramPart string) {
	strArr1 := strings.Split(fullName, "(")
	if len(strArr1) == 2 {
		strArr2 := strings.Split(strArr1[0], ".")
		if len(strArr2) == 2 {
			agentIntanceName = strings.TrimSpace(strArr2[0])
			strArr3 := strings.Split(strArr2[1], "::")
			cnt := len(strArr3)
			if cnt >= 2 {
				agentClassName = strings.TrimSpace(strArr3[cnt-2])
				methodName = strings.TrimSpace(strArr3[cnt-1])
			}
		} else if len(strArr2) == 1 {
			strArr3 := strings.Split(strArr2[0], "::")
			cnt := len(strArr3)
			if cnt >= 2 {
				agentClassName = strings.TrimSpace(strArr3[cnt-2])
				methodName = strings.TrimSpace(strArr3[cnt-1])
			}
		}
		paramPart = strings.TrimSpace(strArr1[1])
		paramPart = strings.TrimRight(paramPart, ")")
	}
	return
}

// ============================================================================
func ParseForParams(src string) (params []string) {
	if src == "" {
		return nil
	}
	params = strings.Split(src, ",")
	for i := 0; i < len(params); i++ {
		params[i] = strings.TrimSpace(params[i])
	}
	return params
}

// ============================================================================
func ParseProperty(val string) *Property {
	strArrArr := strings.Split(val, "[")
	var idxStr string
	if len(strArrArr) > 1 {
		idxStr = strings.Join(strArrArr[1:], "[")
		idxStr = strings.TrimRight(idxStr, "]")
	}
	strArr := strings.Split(strArrArr[0], " ")
	cnt := len(strArr)
	if cnt >= 2 {
		switch strArr[0] {
		// const Int32 0
		case "const":
			if cnt >= 3 {
				return NewPropertyConst(strArr[1], strArr[2])
			}
			// static float Self.AgentNodeTest::s_float_type_0
			// static float Self.AgentNodeTest::s_float_type_0[int Self.AgentNodeTest::par_int_type_2]
		case "static":
			if cnt == 3 {
				return NewPropertyStatic(strArr[1], strArr[2], "", true)
			} else if cnt == 4 {
				return NewPropertyStatic(strArr[1], strArr[2], idxStr, true)
			}
		default:
			if cnt == 2 {
				return NewPropertyStatic(strArr[0], strArr[1], "", false)
			} else if cnt == 3 {
				return NewPropertyStatic(strArr[0], strArr[1], idxStr, false)
			}
		}
	} else {
		return NewPropertyConst(RuntimeTypeStr, strArr[0])
	}

	return nil
}

func ParseInstanceNameProperty(fullName string) (instanceName, agentType, propertyName string) {
	strArr := strings.Split(fullName, ".")
	if len(strArr) > 1 {
		instanceName = strArr[0]
		strArr1 := strings.Split(strArr[1], "::")
		if len(strArr1) > 1 {
			agentType = strArr1[0]
			propertyName = strArr1[1]
		} else {
			propertyName = strArr1[0]
		}
	} else {
		strArr1 := strings.Split(strArr[0], "::")
		if len(strArr1) > 1 {
			agentType = strArr1[0]
			propertyName = strArr1[1]
		} else {
			propertyName = fullName
		}
	}
	return
}
