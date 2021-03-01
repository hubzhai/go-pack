package behaviago

import (
	"errors"
	"strconv"

	"github.com/donnie4w/dom4g"
)

type XmlBehaviorTreeLoader struct {
}

// ============================================================================
const (
	kStrBehavior   = "behavior"
	kStrAgentType  = "agenttype"
	kStrId         = "id"
	kStrPars       = "pars"
	kStrPar        = "par"
	kStrNode       = "node"
	kStrCustom     = "custom"
	kStrProperty   = "property"
	kStrAttachment = "attachment"
	kStrClass      = "class"
	kStrName       = "name"
	kStrType       = "type"
	kStrValue      = "value"
	kStrVersion    = "version"
	kStrFSM        = "fsm"
	kStrFlag       = "flag"
)

func (xbtl *XmlBehaviorTreeLoader) LoadBehaviorTree(buf []byte, bt *BehaviorTree) error {
	el, err := dom4g.LoadByXml(string(buf[:]))
	if err != nil {
		BTGLog.Tracef("dom4g.LoadByXml err(%v)", err)
		return err
	}

	behaviorNode := el.Root()
	if behaviorNode.Name() != kStrBehavior {
		return errors.New("not is behavior")
	}
	name, _ := behaviorNode.AttrValue(kStrName)
	agentType, _ := behaviorNode.AttrValue(kStrAgentType)
	version, _ := behaviorNode.AttrValue(kStrVersion)
	if str, exist := behaviorNode.AttrValue(kStrFSM); exist {
		isFSM, err := strconv.ParseBool(str)
		if err != nil {
			return err
		}
		bt.SetIsFSM(isFSM)
	}
	bt.SetName(name)
	bt.SetClassName("BehaviorTree")
	bt.SetId(int16(-1))
	ver, _ := strconv.Atoi(version)
	err = xbtl.load_properties_pars_attachments_children(true, ver, agentType, behaviorNode, bt)
	return err
}

func (xbtl *XmlBehaviorTreeLoader) load_node(version int, agentType string, node *dom4g.Element) BehaviorNode {
	if node.Name() != kStrNode {
		return nil
	}
	if className, exist := node.AttrValue(kStrClass); exist {
		bn := NewBehaviorNode(className)
		if bn != nil {
			if strId, exist := node.AttrValue(kStrId); exist {
				id, err := strconv.Atoi(strId)
				if err == nil {
					bn.SetId(int16(id))
				}
				err = xbtl.load_properties_pars_attachments_children(true, version, agentType, node, bn)
				if err == nil {
					return bn
				}
			}
		}
	}
	return nil
}

func (xbtl *XmlBehaviorTreeLoader) load_properties_pars_attachments_children(bNode bool, version int, agentType string, node *dom4g.Element, bn BehaviorNode) error {
	bn.SetAgentType(agentType)
	hasEvents := bn.HasEvents()
	childs := node.AllNodes()
	var err error
	var properties []property_t
	for _, c := range childs {
		properties, err = xbtl.load_property_pars(properties, version, agentType, c, bn)
		if err != nil {
			return err
		}
		if bNode {
			switch c.Name() {
			case kStrAttachment:
				if xbtl.load_attachment(version, agentType, c, bn) {
					hasEvents = true
				}
			case kStrCustom:
				customNode := c.Node(kStrNode)
				if customNode != nil {
					child := xbtl.load_node(version, agentType, customNode)
					if child != nil {
						if child.HasEvents() {
							hasEvents = true
						}
						bn.SetCustomCondition(child)
					}
				}
			case kStrNode:
				child := xbtl.load_node(version, agentType, c)
				if child != nil {
					if child.HasEvents() {
						hasEvents = true
					}
					bn.AddChild(child)
				}
			}
		} else {
			if c.Name() == kStrAttachment {
				if xbtl.load_attachment(version, agentType, c, bn) {
					hasEvents = true
				}
			}
		}
	}
	if len(properties) != 0 {
		bn.Load(version, agentType, properties)
	}

	bn.SetHasEvents(hasEvents)
	return nil
}

func (xbtl *XmlBehaviorTreeLoader) load_property_pars(inProperties []property_t, version int, agentType string, node *dom4g.Element, bn BehaviorNode) (outProperties []property_t, err error) {
	switch node.Name() {
	case kStrProperty:
		for _, a := range node.Attrs {
			inProperties = append(inProperties, property_t{name: a.Name(), value: a.Value})
		}
	case kStrPars:
		for _, c := range node.AllNodes() {
			if c.Name() == kStrPar {
				xbtl.load_par(version, agentType, c, bn)
			}
		}
	}
	return inProperties, nil
}

func (xbtl *XmlBehaviorTreeLoader) load_attachment(version int, agentType string, node *dom4g.Element, bn BehaviorNode) (hasEvents bool) {
	className, exist := node.AttrValue(kStrClass)
	if !exist {
		xbtl.load_attachment_transition_effectors(version, agentType, node, bn)
		return true
	}
	pAttachment := NewBehaviorNode(className)
	if pAttachment != nil {
		if strId, exist := node.AttrValue(kStrId); exist {
			id, err := strconv.Atoi(strId)
			if err == nil {
				pAttachment.SetId(int16(id))
			}
		}
		var (
			bIsPrecondition bool
			bIsEffector     bool
			bIsTransition   bool
		)
		flagStr, _ := node.AttrValue(kStrFlag)
		switch flagStr {
		case "precondition":
			bIsPrecondition = true
		case "effector":
			bIsEffector = true
		case "transition":
			bIsTransition = true
		}
		xbtl.load_properties_pars_attachments_children(false, version, agentType, node, pAttachment)
		bn.Attach(pAttachment, bIsPrecondition, bIsEffector, bIsTransition)
		//todo:检测pAttachment是否实现了event接口
		//hasEvents=
		return
	}
	return false
}

func (xbtl *XmlBehaviorTreeLoader) load_attachment_transition_effectors(version int, agentType string, node *dom4g.Element, bn BehaviorNode) {
	bn.SetLoadAttachment(true)
	xbtl.load_properties_pars_attachments_children(false, version, agentType, node, bn)
	bn.SetLoadAttachment(false)
}

func (xbtl *XmlBehaviorTreeLoader) load_par(version int, agentType string, node *dom4g.Element, bn BehaviorNode) {
	attrName, _ := node.AttrValue(kStrName)
	attrType, _ := node.AttrValue(kStrType)
	attrValue, _ := node.AttrValue(kStrValue)
	bn.AddPar(agentType, attrType, attrName, attrValue)
}

// ============================================================================
const (
	kStrAgents                           = "agents"
	kStrAgents_Agent                     = "agent"
	kStrAgents_AgentType                 = "type"
	kStrAgents_AgentProperties           = "properties"
	kStrAgents_AgentProperty             = "property"
	kStrAgents_AgentPropertyName         = "name"
	kStrAgents_AgentPropertyType         = "type"
	kStrAgents_AgentPropertyMember       = "member"
	kStrAgents_AgentPropertyStatic       = "static"
	kStrAgents_AgentPropertyDefaultvalue = "defaultvalue"
	kStrAgents_AgentPropertyAgent        = "agent"
	kStrAgents_AgentMethods              = "methods"
	kStrAgents_AgentMethod               = "method"
	kStrAgents_AgentParameter            = "parameter"
	kStrAgents_AgentParameterType        = "type"
)

func (xbtl *XmlBehaviorTreeLoader) LoadAgents(buf []byte) error {
	el, err := dom4g.LoadByXml(string(buf[:]))
	if err != nil {
		return err
	}
	rootNode := el.Root()
	if rootNode == nil || rootNode.Name() != kStrAgents {
		return errors.New("XmlBehaviorTreeLoader.LoadAgents incorrect xml")
	}

	allNodes := rootNode.Nodes(kStrAgents_Agent)
	for _, node := range allNodes {
		if node.Name() == kStrAgents_Agent {
			typeStr, _ := node.AttrValue(kStrAgents_AgentType)
			ap := NewAgentProperty()
			if ap != nil {
				BTGAgentTypeBlackboards[typeStr] = ap
				allSubNodes := node.AllNodes()
				for _, subNode := range allSubNodes {
					switch subNode.Name() {
					case kStrAgents_AgentProperties:
						propertiesNodes := subNode.Nodes(kStrAgents_AgentProperty)
						for _, pnode := range propertiesNodes {
							nameStr, _ := pnode.AttrValue(kStrAgents_AgentPropertyName)
							typeStr, _ := pnode.AttrValue(kStrAgents_AgentPropertyType)
							memberStr, _ := pnode.AttrValue(kStrAgents_AgentPropertyMember)
							staticStr, _ := pnode.AttrValue(kStrAgents_AgentPropertyStatic)
							var valueStr string
							var agentTypeMember string
							bIsMember, _ := strconv.ParseBool(memberStr)
							if !bIsMember {
								valueStr, _ = pnode.AttrValue(kStrAgents_AgentPropertyDefaultvalue)
							} else {
								agentTypeMember, _ = pnode.AttrValue(kStrAgents_AgentPropertyAgent)
							}
							bIsStatic, _ := strconv.ParseBool(staticStr)
							ap.AddProperty(typeStr, bIsStatic, nameStr, valueStr, agentTypeMember)
						}
					case kStrAgents_AgentMethods:
						methodsNodes := subNode.Nodes(kStrAgents_AgentMethod)
						for _, mnode := range methodsNodes {
							var paramTypes []string
							paramNodes := mnode.Nodes(kStrAgents_AgentParameter)
							for _, paramNode := range paramNodes {
								strType, _ := paramNode.AttrValue(kStrAgents_AgentParameterType)
								paramTypes = append(paramTypes, strType)
							}
							//todo:registe agent method
						}
					}
				}
			}
		}
	}
	return nil
}

func init() {
	RegisteLoader(EFF_xml, &XmlBehaviorTreeLoader{})
}
