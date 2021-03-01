package behaviago

import (
	"hash/crc32"
	"io/ioutil"
	"path/filepath"
)

var BTGAgentTypeBlackboards = make(map[string]*AgentProperties)

type AgentProperties struct {
	properties         map[uint32]*Property
	locals             map[uint32]*Property
	propertiesInstance []*Property
}

func AgentProperty_AddLocal(agentType, typeName, variableName, valueStr string) *Property {
	ap := AgentProperty_Get(agentType)
	if ap == nil {
		ap = NewAgentProperty()
		BTGAgentTypeBlackboards[agentType] = ap
	}
	return ap.AddLocal(typeName, variableName, valueStr)
}

func AgentProperty_AddInstance(agentType string, instance *Property) {
	ap := AgentProperty_Get(agentType)
	if ap == nil {
		ap = NewAgentProperty()
		BTGAgentTypeBlackboards[agentType] = ap
	}
	ap.AddPropertyInstance(instance)
}

func AgentProperty_Get(agentType string) *AgentProperties {
	if ap, exist := BTGAgentTypeBlackboards[agentType]; exist {
		return ap
	}
	return nil
}

func AgentProperty_GetPropertyByVarName(agentType, variableName string) *Property {
	if ap, exist := BTGAgentTypeBlackboards[agentType]; exist {
		return ap.GetPropertyByVarName(variableName)
	}
	return nil
}

func AgentProperty_GetPropertyByVarId(agentType string, variableId uint32) *Property {
	if ap, exist := BTGAgentTypeBlackboards[agentType]; exist {
		return ap.GetPropertyByVarId(variableId)
	}
	return nil
}

func AgentProperty_UnloadLocals() {
	for _, ap := range BTGAgentTypeBlackboards {
		ap.ClearLocals()
	}
}

func AgentProperty_Cleanup() {
	for _, ap := range BTGAgentTypeBlackboards {
		ap.Cleanup()
	}
}

func MakeVariableId(variableName string) uint32 {
	return crc32.ChecksumIEEE([]byte(variableName))
}

func NewAgentProperty() *AgentProperties {
	ap := &AgentProperties{
		properties: make(map[uint32]*Property),
		locals:     make(map[uint32]*Property),
	}
	return ap
}

func (ap *AgentProperties) GetPropertyByVarId(variableId uint32) *Property {
	BTGLog.Tracef("(ap *AgentProperties) GetPropertyByVarId(variableId=%v)", variableId)
	if p, exist := ap.properties[variableId]; exist {
		return p
	}
	if p, exist := ap.locals[variableId]; exist {
		return p
	}
	return nil
}

func (ap *AgentProperties) GetPropertyByVarName(variableName string) *Property {
	return ap.GetPropertyByVarId(MakeVariableId(variableName))
}

func (ap *AgentProperties) AddPropertyInstance(p *Property) {
	ap.propertiesInstance = append(ap.propertiesInstance, p)
}

func (ap *AgentProperties) AddProperty(typeName string, bIsStatic bool, variableName, valueStr, agentType string) *Property {
	p := NewProperty(typeName, "Self", variableName, valueStr)
	p.IsStatic = bIsStatic
	p.IsMember = GetAgentMemberMeta(agentType, variableName) != nil
	varId := MakeVariableId(variableName)
	BTGLog.Tracef("AgentProperty.AddProperty varName=%v->varId=%v", variableName, varId)
	ap.properties[varId] = p
	return p
}

func (ap *AgentProperties) AddLocal(typeName, variableName, valueStr string) *Property {
	p := NewProperty(typeName, "Self", variableName, valueStr)
	p.IsLocal = true
	varId := MakeVariableId(variableName)
	BTGLog.Tracef("AgentProperty.AddLocal varName=%v->varId=%v", variableName, varId)
	ap.locals[varId] = p
	return p
}

func (ap *AgentProperties) ClearLocals() {
	ap.locals = make(map[uint32]*Property)
}

func (ap *AgentProperties) Cleanup() {
	ap.properties = make(map[uint32]*Property)
	ap.locals = make(map[uint32]*Property)
	ap.propertiesInstance = nil
}

func (ap *AgentProperties) GetLocal(variableName string) *Property {
	if len(ap.locals) != 0 {
		varId := MakeVariableId(variableName)
		if p, exist := ap.locals[varId]; exist {
			return p
		}
	}
	return nil
}

func LoadAgentProperties() error {
	relativePath := "behaviac.bb"
	fullPath := filepath.Join(BTGWorkspace.ExportPath, relativePath)
	switch BTGWorkspace.FileFormat {
	case EFF_xml:
		fullPath = fullPath + ".xml"
	}
	buf, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return err
	}
	BTGLog.Tracef("AgentProperties.Load fullPath=%v", fullPath)

	btl := GetLoader(BTGWorkspace.FileFormat)
	if btl != nil {
		err = btl.LoadAgents(buf)
		if err != nil {
			BTGLog.Warnf("'%v' is not loaded!(error:%v)", fullPath, err)
			return err
		}
	} else {
		BTGLog.Warnf("'%v' is not registe fit Loader(%v)!", fullPath, BTGWorkspace.FileFormat)
	}
	return nil
}

func RegisterCustomizedTypes() error {
	return nil
}
