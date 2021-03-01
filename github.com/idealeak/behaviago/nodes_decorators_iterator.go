package behaviago

import (
	"strings"
)

type DecoratorIterator struct {
	DecoratorNode
	opl  *Property
	opr  *Property
	oprm *Method
}

func (d *DecoratorIterator) Load(version int, agentType string, properties []property_t) {
	d.DecoratorNode.Load(version, agentType, properties)
	for i := 0; i < len(properties); i++ {
		switch properties[i].name {
		case "Opl":
			if !strings.ContainsAny(properties[i].value, "(") {
				BTGLog.Warnf("DecoratorIterator.Load Opl contains '('")
			} else {
				d.opl = ParseProperty(properties[i].value)
			}
		case "Opr":
			if !strings.ContainsAny(properties[i].value, "(") {
				d.oprm = LoadMethod(properties[i].value)
			} else {
				d.opr = ParseProperty(properties[i].value)
			}
		}
	}
}

func (n *DecoratorIterator) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*DecoratorIterator); !ok {
		return false
	}
	return true
}

func (n *DecoratorIterator) Decompose(node BehaviorNode, seqTask *PlannerTaskComplex, depth int, planner *Planner) bool {
	//TODO DecoratorIterator.Decompose
	return true
}

func (n *DecoratorIterator) IterateIt(a *Agent, index int) bool {
	if n.oprm != nil {
		lp := n.opl.GetParentAgent(a)
		rp := n.oprm.GetParentAgent(a)
		return n.opl.SetFromByMethodByIndex(rp, n.oprm, lp, index)
	} else if n.opr != nil {
		lp := n.opl.GetParentAgent(a)
		rp := n.opr.GetParentAgent(a)
		return n.opl.SetFromByProperty(rp, n.opr, lp)
	}
	return false
}
