package gmx

/*
Classes = [{"name":"xxx", "type":"string"},{"name":"yyy", "type":"int"}]
*/
type TypeDesc struct {
	Name       string          `json:"name"`
	Type       string          `json:"type"`
	IsExported bool            `json:"isexported"`
	Fields     []FieldTypeDesc `json:"fields"`
	Methods    []MethodDesc    `json:"methods"`
}

/*
Class = {"name":"xxx", "type":"string"},{"name":"yyy", "type":"int"}
*/
type FieldTypeDesc struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Kind       uint   `json:"kind"`
	Tag        string `json:"tag"`
	Offset     int    `json:"offset"`
	Index      []int  `json:"index"`
	Anonymous  bool   `json:"anonymous"`
	IsExported bool   `json:"isexported"`
}

type MethodDesc struct {
	Name      string        `json:"name"`
	Index     int           `json:"index"`
	RcvrIsPtr bool          `json:"rcvrisptr"`
	In        []MethodParam `json:"in"`
	Out       []MethodParam `json:"out"`
}

type MethodParam struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Kind  uint   `json:"kind"`
	Index int    `json:"index"`
}

type FieldValueDesc struct {
	FieldTypeDesc
	Value interface{} `json:"value"`
}

type ValueDesc struct {
	Name    string           `json:"name"`
	Type    string           `json:"type"`
	Kind    uint             `json:"kind"`
	CanSet  bool             `json:"canset"`
	Value   interface{}      `json:"value"`
	Fields  []FieldValueDesc `json:"fields"`
	Methods []MethodDesc     `json:"methods"`
}

type MapDesc struct {
	Name         string
	Type         string
	PageNo       int
	NumOfPerPage int
	TotalCount   int
	Keys         []*ValueDesc
	Vals         []*ValueDesc
}

type SliceDesc struct {
	Name         string
	Type         string
	PageNo       int
	NumOfPerPage int
	TotalCount   int
	ElemType     string
	Elems        interface{}
}

type InvokeResult struct {
	Name     string
	FuncName string
	RetVals  []*ValueDesc
}
