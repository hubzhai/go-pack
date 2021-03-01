package behaviago

var BTLoaders = make(map[EFileFormat]BTreeLoader)

// ============================================================================
type BTreeLoader interface {
	LoadBehaviorTree(buf []byte, bt *BehaviorTree) error
	LoadAgents(buf []byte) error
}

func RegisteLoader(eff EFileFormat, btl BTreeLoader) {
	BTLoaders[eff] = btl
}

func GetLoader(eff EFileFormat) BTreeLoader {
	return BTLoaders[eff]
}
