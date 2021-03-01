package behaviago

import "strings"

// ============================================================================
type EOperatorType int

const (
	E_INVALID      EOperatorType = iota
	E_ASSIGN                     // =
	E_ADD                        // +
	E_SUB                        // -
	E_MUL                        // *
	E_DIV                        // /
	E_EQUAL                      // ==
	E_NOTEQUAL                   // !=
	E_GREATER                    // >
	E_LESS                       // <
	E_GREATEREQUAL               // >=
	E_LESSEQUAL                  // <=
)

// ============================================================================
type ActionConfig struct {
	opl      *Property
	oplm     *Method
	opr1     *Property
	oprm1    *Method
	operator EOperatorType
	opr2     *Property
	oprm2    *Method
}

func (ac *ActionConfig) Load(properties []property_t) {
	for i := 0; i < len(properties); i++ {
		switch properties[i].name {
		case "Opl":
			if !strings.ContainsAny(properties[i].value, "(") {
				ac.opl = ParseProperty(properties[i].value)
			} else {
				ac.oplm = LoadMethod(properties[i].value)
			}
		case "Opr1":
			if !strings.ContainsAny(properties[i].value, "(") {
				ac.opr1 = ParseProperty(properties[i].value)
			} else {
				ac.oprm1 = LoadMethod(properties[i].value)
			}
		case "Operator":
			switch properties[i].value {
			case "Invalid":
				ac.operator = E_INVALID
			case "Assign":
				ac.operator = E_ASSIGN
			case "Add":
				ac.operator = E_ADD
			case "Sub":
				ac.operator = E_SUB
			case "Mul":
				ac.operator = E_MUL
			case "Div":
				ac.operator = E_DIV
			case "Equal":
				ac.operator = E_EQUAL
			case "NotEqual":
				ac.operator = E_NOTEQUAL
			case "Greater":
				ac.operator = E_GREATER
			case "Less":
				ac.operator = E_LESS
			case "GreaterEqual":
				ac.operator = E_GREATEREQUAL
			case "LessEqual":
				ac.operator = E_LESSEQUAL
			}
		case "Opr2":
			if !strings.ContainsAny(properties[i].value, "(") {
				ac.opr2 = ParseProperty(properties[i].value)
			} else {
				ac.oprm2 = LoadMethod(properties[i].value)
			}
		}
	}
}

func (ac *ActionConfig) Execute(a *Agent) bool {
	BTGLog.Trace("(ac *ActionConfig) Execute")
	var bValid bool
	if ac.oplm != nil && ac.operator == E_INVALID { // action
		bValid = true
		ac.oplm.Invoke(a)
	} else if ac.operator == E_ASSIGN { // assign
		bValid = EvaluteAssignment(a, ac.opl, ac.opr2, ac.oprm2)
	} else if ac.operator >= E_ADD && ac.operator <= E_DIV {
		computeOperator := ECO_ADD + EComputeOperator(ac.operator-E_ADD)
		bValid = EvaluteCompute(a, ac.opl, ac.opr1, ac.oprm1, computeOperator, ac.opr2, ac.oprm2)
	} else if ac.operator >= E_EQUAL && ac.operator <= E_LESSEQUAL {
		var compareOperator E_VariableComparisonType
		switch ac.operator {
		case E_EQUAL: // ==
			compareOperator = VariableComparisonType_Equal
		case E_NOTEQUAL: // !=
			compareOperator = VariableComparisonType_NotEqual
		case E_GREATER: // >
			compareOperator = VariableComparisonType_Greater
		case E_LESS: // <
			compareOperator = VariableComparisonType_Less
		case E_GREATEREQUAL: // >=
			compareOperator = VariableComparisonType_GreaterEqual
		case E_LESSEQUAL: // <=
			compareOperator = VariableComparisonType_LessEqual
		default:
			compareOperator = VariableComparisonType_Assignment
		}
		bValid = EvaluteCompare(a, ac.opl, ac.oplm, compareOperator, ac.opr2, ac.oprm2)
	}
	return bValid
}

// ============================================================================
type AttachAction struct {
	BehaviorNodeBase
	ac *ActionConfig
}

func NewAttachAction() *AttachAction {
	return &AttachAction{ac: &ActionConfig{}}
}

func (aa *AttachAction) Load(version int, agentType string, properties []property_t) {
	aa.BehaviorNodeBase.Load(version, agentType, properties)
	aa.ac.Load(properties)
}

func (aa *AttachAction) Evaluate(a *Agent, result EBTStatus) bool {
	BTGLog.Tracef("AttachAction.Evaluate enter(%v)", aa.GetClassName())
	if aa.ac != nil {
		if aa.ac.Execute(a) {
			return aa.UpdateImpl(a, BT_INVALID) == BT_SUCCESS
		}
	}
	return false
}

func (aa *AttachAction) CreateTask() BehaviorTask {
	BTGLog.Warn("AttachAction.CreateTask unsupport!!!")
	return nil
}
