package behaviago

// ============================================================================
type DecoratorCountInterface interface {
	GetCount(a *Agent) int64
}

// ============================================================================
type DecoratorCount struct {
	DecoratorNode
	countVar *Property
}

func (d *DecoratorCount) Load(version int, agentType string, properties []property_t) {
	d.DecoratorNode.Load(version, agentType, properties)
	for i := 0; i < len(properties); i++ {
		if properties[i].name == "Count" {
			d.countVar = ParseProperty(properties[i].value)
		}
	}
}

func (d *DecoratorCount) GetCount(a *Agent) int64 {
	BTGLog.Trace("(dt *DecoratorCount) GetCount enter")
	if d.countVar != nil {
		n := d.countVar.GetInt(a)
		BTGLog.Tracef("(dt *DecoratorCount) GetCount enter %v", n)
		return n
	}
	return 0
}

// ============================================================================
type DecoratorCountTask struct {
	DecoratorTask
	count int64
}

func (dt *DecoratorCountTask) GetCount(a *Agent) int64 {
	BTGLog.Trace("(dt *DecoratorCountTask) GetCount enter")
	if node, ok := dt.GetNode().(DecoratorCountInterface); ok {
		return node.GetCount(a)
	}
	return 0
}

func (dt *DecoratorCountTask) OnEnter(a *Agent) bool {
	BTGLog.Trace("(dt *DecoratorCountTask) OnEnter enter")
	count := dt.GetCount(a)
	if count == 0 {
		return false
	}
	dt.count = count
	return dt.DecoratorTask.OnEnter(a)
}
