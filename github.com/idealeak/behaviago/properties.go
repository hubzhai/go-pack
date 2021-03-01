package behaviago

type Variable struct {
	id       uint32
	name     string
	property *Property
	changed  bool
}

func (v *Variable) IsChanged() bool {
	return v.changed
}

func (v *Variable) IsLocal() bool {
	if v.property != nil {
		return v.property.IsLocal
	}
	return false
}

func (v *Variable) IsMember() bool {
	return false
}

func (v *Variable) GetTypeId() int {
	return 0
}

func (v *Variable) GetId() uint32 {
	return v.id
}

func (v *Variable) Name() string {
	return v.name
}

func (v *Variable) GetProperty() *Property {
	return v.property
}

func (v *Variable) SetProperty(p *Property) {
	v.property = p
}

func (v *Variable) Reset() {
	v.changed = false
}

func (v *Variable) CheckIfChanged(a *Agent) bool {
	return false
}

type Variables struct {
	vars map[uint32]*Variable
}

func (v *Variables) Clear(bFull bool) {
	if bFull {
		v.vars = make(map[uint32]*Variable)
	} else {
		for k, p := range v.vars {
			if p.IsLocal() {
				delete(v.vars, k)
			}
		}
	}
}

func (v *Variables) IsExisting(id uint32) bool {
	_, exist := v.vars[id]
	return exist
}
