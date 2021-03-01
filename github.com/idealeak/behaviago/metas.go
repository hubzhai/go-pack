package behaviago

import (
	"encoding/xml"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// ============================================================================
//parentType
const (
	PT_INVALID int = iota
	PT_AGENT
	PT_INSTANCE
	PT_PAR
)

var AgentMetasMap = make(map[string]*AgentMeta)

type AgentMeta struct {
	ClassFullName string
	Base          string
	DisplayName   string
	Desc          string
	IsRefType     bool
	v             reflect.Value
	members       map[string]*AgentMemberMetaStruct
	methods       map[string]*AgentMethodMetaStruct
}

// ============================================================================
type AgentMemberMetaStruct struct {
	*AgentMeta
	Type        string
	Name        string
	DisplayName string
	Desc        string
	Range       float64
	IsStatic    bool
	IsConst     bool
	parentType  int
}

type AgentMethodParamMetaStruct struct {
	Name        string
	DisplayName string
	Desc        string
}

type AgentMethodMetaStruct struct {
	*AgentMeta
	Name        string
	DisplayName string
	Desc        string
	InParams    []AgentMethodParamMetaStruct
	OutParams   []AgentMethodParamMetaStruct
}

func REGISTER_META_METHODS(metas ...*AgentMethodMetaStruct) []*AgentMethodMetaStruct {
	return metas
}

func REGISTER_META_METHOD(name string) *AgentMethodMetaStruct {
	return &AgentMethodMetaStruct{Name: name}
}

func (amms *AgentMethodMetaStruct) DESC(desc string) *AgentMethodMetaStruct {
	amms.Desc = desc
	return amms
}

func (amms *AgentMethodMetaStruct) DISPLAYNAME(displayName string) *AgentMethodMetaStruct {
	amms.DisplayName = displayName
	return amms
}

func (amms *AgentMethodMetaStruct) IPARAM_DISPLAY_INFO(name, paramDisplayName, paramDesc string) *AgentMethodMetaStruct {
	ampms := AgentMethodParamMetaStruct{Name: name, DisplayName: paramDisplayName, Desc: paramDesc}
	amms.InParams = append(amms.InParams, ampms)
	return amms
}

func (amms *AgentMethodMetaStruct) OPARAM_DISPLAY_INFO(name, paramDisplayName, paramDesc string) *AgentMethodMetaStruct {
	ampms := AgentMethodParamMetaStruct{Name: name, DisplayName: paramDisplayName, Desc: paramDesc}
	amms.OutParams = append(amms.OutParams, ampms)
	return amms
}

// ============================================================================
func REGISTER_AGENT_METAS(a interface{}, methodsMeta []*AgentMethodMetaStruct) {
	v := reflect.ValueOf(a)
	t := reflect.Indirect(v).Type()
	if t.Kind() != reflect.Struct {
		return
	}
	if _, is := t.FieldByName("Agent"); !is {
		return
	}
	k := fmt.Sprintf("%v/%v", t.PkgPath(), t.Name())
	if _, exist := AgentMetasMap[k]; !exist {
		superCls := agentSuperName(t, reflect.Indirect(v))
		am := &AgentMeta{
			ClassFullName: k,
			Base:          superCls,
			DisplayName:   t.Name(),
			Desc:          t.Name(),
			v:             v,
		}
		parseAgentMemberMeta(reflect.Indirect(v), t, am, "")
		parseAgentMethodMeta(am, methodsMeta)
		AgentMetasMap[k] = am
	}
}

func GetAgentMeta(agentName string) *AgentMeta {
	if am, exist := AgentMetasMap[agentName]; exist {
		return am
	}
	return nil
}

func GetAgentMemberMeta(agentName, memberName string) *AgentMemberMetaStruct {
	am := GetAgentMeta(agentName)
	if am != nil {
		return am.GetMemberMeta(memberName)
	}
	return nil
}

func (am *AgentMeta) GetMemberMeta(name string) *AgentMemberMetaStruct {
	if am.members != nil {
		if amms, ok := am.members[name]; ok {
			return amms
		}
	}
	return nil
}

func parseAgentMemberMeta(v reflect.Value, t reflect.Type, am *AgentMeta, prefix string) {
	am.members = make(map[string]*AgentMemberMetaStruct)
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		kind := reflect.Indirect(v.FieldByName(ft.Name)).Kind()
		tag := ft.Tag.Get("behaviago")
		if tag != "" && kind != reflect.Struct {
			mn := &AgentMemberMetaStruct{AgentMeta: am}
			if mn != nil {
				tags := agentParseTag(tag)
				mn.Name = ft.Name
				ftt := reflect.Indirect(v.FieldByName(ft.Name)).Type()
				mn.DisplayName = agentGetTag(tags, "DisplayName")
				mn.Desc = agentGetTag(tags, "Desc")
				mn.Type = editorRefType(ftt.Name())
				am.members[mn.Name] = mn
			}
		}
		//export nest fields
		if ft.Anonymous && kind == reflect.Struct {
			ftv := reflect.Indirect(v.FieldByName(ft.Name))
			ftt := ftv.Type()
			var prefixNext string
			if prefix != "" {
				prefixNext = prefix + "." + ft.Name
			} else {
				prefixNext = ft.Name
			}
			parseAgentMemberMeta(ftv, ftt, am, prefixNext)
		}
	}
}
func parseAgentMethodMeta(am *AgentMeta, methodsMeta []*AgentMethodMetaStruct) {
	am.methods = make(map[string]*AgentMethodMetaStruct)
	for _, method := range methodsMeta {
		method.AgentMeta = am
		am.methods[method.Name] = method
	}
}

// ============================================================================
type AgentMemberXMLHolder struct {
	XMLName     xml.Name `xml:"Member"`
	Name        string   `xml:"Name,attr"`
	Type        string   `xml:"Type,attr"`
	Class       string   `xml:"Class,attr"`
	DisplayName string   `xml:"DisplayName,attr"`
	Desc        string   `xml:"Desc,attr"`
}
type AgentMethodParamXMLHolder struct {
	XMLName     xml.Name `xml:"Param"`
	Name        string   `xml:"Name,attr"`
	Type        string   `xml:"Type,attr"`
	DisplayName string   `xml:"DisplayName,attr"`
	Desc        string   `xml:"Desc,attr"`
}
type AgentMethodXMLHolder struct {
	XMLName     xml.Name `xml:"Method"`
	Name        string   `xml:"Name,attr"`
	ReturnType  string   `xml:"ReturnType,attr"`
	Class       string   `xml:"Class,attr"`
	DisplayName string   `xml:"DisplayName,attr"`
	Desc        string   `xml:"Desc,attr"`
	Params      []*AgentMethodParamXMLHolder
}
type AgentXMLHolder struct {
	XMLName       xml.Name `xml:"agent"`
	ClassFullName string   `xml:"classfullname,attr"`
	Base          string   `xml:"base,attr"`
	DisplayName   string   `xml:"DisplayName,attr"`
	Desc          string   `xml:"Desc,attr"`
	IsRefType     bool     `xml:"IsRefType,attr"`
	Member        []*AgentMemberXMLHolder
	Method        []*AgentMethodXMLHolder
}
type AgentsXMLHolder struct {
	XMLName xml.Name `xml:"agents"`
	Agent   []*AgentXMLHolder
}

func ExportAgentMetas(fileName string) bool {
	if !strings.HasSuffix(fileName, ".xml") {
		fileName = fileName + ".xml"
	}
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		BTGLog.Warnf("%v open failed", fileName)
		return false
	}
	defer f.Close()
	enc := xml.NewEncoder(f)
	enc.Indent("\t", "\t")
	agentsXML := AgentsXMLHolder{}
	for _, am := range AgentMetasMap {
		v := am.v
		vv := reflect.Indirect(v)
		t := vv.Type()
		superCls := agentSuperName(t, vv)
		agent := &AgentXMLHolder{
			ClassFullName: t.Name(),
			Base:          superCls,
			DisplayName:   t.Name(),
			Desc:          t.Name(),
			IsRefType:     true,
		}
		agentsXML.Agent = append(agentsXML.Agent, agent)
		agentXMLExportFields(am, agent, t.Name())
		agentXMLExportMethods(v, t, agent, am.methods, superCls, t.Name(), "")
	}
	f.WriteString(`<metas version="5" language="golang">`)
	f.WriteString("\n")
	if err = enc.Encode(agentsXML); err != nil {
		BTGLog.Errorf("error: %v\n", err)
		return false
	}
	f.WriteString("\n")
	f.WriteString(`</metas>`)
	return true
}

func agentXMLExportFields(am *AgentMeta, pnode *AgentXMLHolder, className string) {
	for _, member := range am.members {
		mn := &AgentMemberXMLHolder{}
		if mn != nil {
			mn.Name = member.Name
			mn.DisplayName = member.DisplayName
			mn.Desc = member.Desc
			mn.Type = member.Type
			mn.Class = className
			pnode.Member = append(pnode.Member, mn)
		}
	}
}

func agentXMLExportMethods(v reflect.Value, t reflect.Type, pnode *AgentXMLHolder, methods map[string]*AgentMethodMetaStruct, superName, className, prefix string) {
	if am, exist := AgentMetasMap[superName]; exist {
		vv := reflect.Indirect(am.v)
		t := vv.Type()
		superCls := agentSuperName(t, vv)
		var prefixNext string
		if prefix != "" {
			prefixNext = prefix + "." + t.Name()
		} else {
			prefixNext = t.Name()
		}
		agentXMLExportMethods(am.v, t, pnode, am.methods, superCls, className, prefixNext)
	} else {
		for _, mmeta := range methods {
			mv := v.MethodByName(mmeta.Name)
			if !mv.IsValid() {
				BTGLog.Warnf("%v no have %v method", className, mmeta.Name)
				continue
			}
			mf := mv.Type()
			if !mf.IsVariadic() && mf.NumIn() != len(mmeta.InParams) {
				BTGLog.Warnf("%v.%v method in params count not fit", className, mmeta.Name)
				continue
			}
			mn := &AgentMethodXMLHolder{}
			if mn != nil {
				mn.Name = mmeta.Name
				mn.DisplayName = mmeta.DisplayName
				mn.Desc = mmeta.Desc
				mn.Class = className
				if mf.NumOut() != 0 {
					mn.ReturnType = editorRefType(mf.Out(0).Name())
				} else {
					mn.ReturnType = "void"
				}
				if !mf.IsVariadic() {
					for j := 0; j < mf.NumIn(); j++ {
						pn := &AgentMethodParamXMLHolder{}
						p := mf.In(j)
						if pn != nil {
							pn.Name = mmeta.InParams[j].Name
							pn.DisplayName = mmeta.InParams[j].DisplayName
							pn.Desc = mmeta.InParams[j].Desc
							pn.Type = editorRefType(p.Name())
							mn.Params = append(mn.Params, pn)
						}
					}
				}
				pnode.Method = append(pnode.Method, mn)
			}
		}
	}
}

func agentSuperName(at reflect.Type, av reflect.Value) string {
	for i := 0; i < at.NumField(); i++ {
		f := at.Field(i)
		tag := f.Tag.Get("behaviago")
		if f.Anonymous && tag != "" && strings.ContainsAny(tag, "IsSuper") {
			return f.Name
			//			ft := reflect.Indirect(av.FieldByName(f.Name)).Type()
			//			return ft.PkgPath() + "/" + f.Name
		}
	}
	return ""
}

func agentParseTag(tag string) map[string]string {
	tags := strings.Split(tag, "|")
	ret := make(map[string]string)
	for _, t := range tags {
		strs := strings.Split(t, "=")
		if len(strs) == 1 {
			ret[strs[0]] = ""
		} else if len(strs) == 2 {
			ret[strs[0]] = strs[1]
		}
	}
	return ret
}

func agentGetTag(m map[string]string, key string) string {
	if v, exist := m[key]; exist {
		return v
	}
	return ""
}
