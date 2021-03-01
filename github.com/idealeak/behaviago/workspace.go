package behaviago

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type EFileFormat int

const (
	EFF_xml EFileFormat = iota
)

var BTGConfig = Config{}
var BTGWorkspace = NewWorkspace()

// ============================================================================
type Config struct {
	IsDesktopPlatform  bool
	IsLogging          bool
	IsLoggingFlush     bool
	IsSocketing        bool
	IsProfiling        bool
	IsSocketIsBlocking bool
	IsHotReload        bool
	SocketPort         uint16
}

// ============================================================================
type property_t struct {
	name  string
	value string
}

// ============================================================================
type BehaviorTreeCreator interface {
	CreateBehaviorTree() *BehaviorTree
}

type BehaviorTreeCreatorWrapper func() *BehaviorTree

func (btcw BehaviorTreeCreatorWrapper) CreateBehaviorTree() *BehaviorTree {
	return btcw()
}

// ============================================================================
type BTGItem struct {
	bts    []*BehaviorTreeTask
	agents []*Agent
}

// ============================================================================
/**
  this is called for every behavior node, in which users can do some custom stuff
*/
type BehaviorNodeLoader func(string, []property_t)

// ============================================================================
type Workspace struct {
	ExportPath           string
	isInited             bool
	isExecAgents         bool
	FileFormat           EFileFormat
	monitor              *FileSysMonitor
	behaviorNodeLoader   BehaviorNodeLoader
	behaviorTrees        map[string]*BehaviorTree
	behaviortreeCreators map[string]BehaviorTreeCreator
	allBehaviorTreeTasks map[string]*BTGItem
	frameSinceStartup    int64
	timeSinceStartup     int64
}

func NewWorkspace() *Workspace {
	ws := &Workspace{
		ExportPath:           "./behaviac/workspace/exported/",
		isInited:             false,
		isExecAgents:         true,
		FileFormat:           EFF_xml,
		behaviorTrees:        make(map[string]*BehaviorTree),
		behaviortreeCreators: make(map[string]BehaviorTreeCreator),
		allBehaviorTreeTasks: make(map[string]*BTGItem),
	}
	return ws
}

func (c *Config) LogInfo() {
	BTGLog.Infof("Config::IsDesktopPlatform %v", c.IsDesktopPlatform)
	BTGLog.Infof("Config::IsLogging %v", c.IsLogging)
	BTGLog.Infof("Config::IsLoggingFlush %v", c.IsLoggingFlush)
	BTGLog.Infof("Config::IsSocketing %v", c.IsSocketing)
	BTGLog.Infof("Config::IsProfiling %v", c.IsProfiling)
	BTGLog.Infof("Config::IsSocketIsBlocking %v", c.IsSocketIsBlocking)
	BTGLog.Infof("Config::IsHotReload %v", c.IsHotReload)
	BTGLog.Infof("Config::SocketPort %v", c.SocketPort)
}

func (w *Workspace) SetBehaviorNodeLoader(loaderCallback BehaviorNodeLoader) {
	w.behaviorNodeLoader = loaderCallback
}

func (w *Workspace) BehaviorNodeLoaded(nodeType string, properties []property_t) {
	if w.behaviorNodeLoader != nil {
		w.behaviorNodeLoader(nodeType, properties)
	}
}

func (w *Workspace) LoadWorkspaceSetting(file string) (string, error) {
	type XmlWorkspace struct {
		XMLName xml.Name `xml:"workspace"`
		Path    string   `xml:"path,attr"`
	}
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	xmlws := &XmlWorkspace{}
	err = xml.Unmarshal(buf, xmlws)
	if err != nil {
		return "", err
	}
	return xmlws.Path, nil
}

func (w *Workspace) SetFilePath(exportPath string) {
	exportPath = strings.TrimSpace(exportPath)
	if strings.HasSuffix(exportPath, "/") || strings.HasSuffix(exportPath, "\\") {
		w.ExportPath = exportPath
	} else {
		w.ExportPath = exportPath + "/"
	}
}

func (w *Workspace) TryInit() error {
	if w.isInited {
		return nil
	}
	w.isInited = true

	BTGConfig.LogInfo()

	err := TryStartup()
	if err != nil {
		return err
	}

	if w.ExportPath == "" {
		return errors.New("No 'WorkspaceExportPath' is specified!")
	}

	BTGLog.Infof("'WorkspaceExportPath' is '%v'", w.ExportPath)
	if BTGConfig.IsHotReload {
		w.monitor = NewFileSysMonitor()
		if w.monitor != nil {
			w.monitor.SetMonitorDir(w.ExportPath)
			err := w.monitor.Start()
			if err != nil {
				return err
			}

		}
	}

	//
	RegisterCustomizedTypes()
	//
	LoadAgentProperties()

	return nil
}

func (w *Workspace) Cleanup() {
	if BTGConfig.IsHotReload {
		if w.monitor != nil {
			w.monitor.Stop()
			w.monitor = nil
		}
	}

	w.UnRegisterBehaviorTreeCreators()

	BaseCleanup()

	w.isInited = false
}

func (w *Workspace) RegisterBehaviorTreeCreator(relativePath string, creator BehaviorTreeCreator) bool {
	if relativePath != "" && creator != nil {
		w.behaviortreeCreators[relativePath] = creator
		return true
	}
	return false
}

func (w *Workspace) UnRegisterBehaviorTreeCreators() {
	w.behaviortreeCreators = make(map[string]BehaviorTreeCreator)
}

func (w *Workspace) GetBehaviors() map[string]*BehaviorTree {
	return w.behaviorTrees
}

func (w *Workspace) Load(relativePath string, bForce bool) error {
	err := w.TryInit()
	if err != nil {
		return err
	}
	var bt *BehaviorTree
	var exist bool
	if bt, exist = w.behaviorTrees[relativePath]; !exist && !bForce {
		return nil
	}

	if bt != nil && !bForce {
		return nil
	}

	fullPath := filepath.Join(w.ExportPath, relativePath)
	switch w.FileFormat {
	case EFF_xml:
		fullPath = fullPath + ".xml"
	}
	_, err = os.Stat(fullPath)
	if err != nil {
		return err
	}
	var bNewly bool
	if bt == nil {
		bNewly = true
		bt = NewBehaviorTree()
		w.behaviorTrees[relativePath] = bt
	}
	if bt == nil {
		BTGLog.Errorf("Workspace.Load(%v,%v) failed", relativePath, bForce)
		return nil
	}
	var bCleared bool
	buf, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return err
	}
	if !bNewly {
		bCleared = true
		bt.Clear()
	}

	BTGLog.Tracef("Workspace.Load fullPath=%v", fullPath)

	btl := GetLoader(w.FileFormat)
	if btl != nil {
		err = btl.LoadBehaviorTree(buf, bt)
		if err != nil {
			if bNewly || bCleared {
				delete(w.behaviorTrees, relativePath)
			}
			BTGLog.Warnf("'%v' is not loaded!(error:%v)", fullPath, err)
			return err
		}
	} else {
		BTGLog.Warnf("'%v' is not registe fit Loader(%v)!", fullPath, w.FileFormat)
	}
	return nil
}

/**
  hot reload the modified behaviors.
*/
func (w *Workspace) RecordBTGAgentMapping(relativePath string, a *Agent) {
	if v, exist := w.allBehaviorTreeTasks[relativePath]; exist {
		bfound := false
		for _, aa := range v.agents {
			if aa == a {
				bfound = true
				break
			}
		}
		if !bfound {
			v.agents = append(v.agents, a)
		}
	} else {
		w.allBehaviorTreeTasks[relativePath] = &BTGItem{agents: []*Agent{a}}
	}
}

func (w *Workspace) SetTimeSinceStartup(timeSinceStartup int64) {
	w.timeSinceStartup = timeSinceStartup
}

func (w *Workspace) GetTimeSinceStartup() int64 {
	return w.timeSinceStartup
}

func (w *Workspace) SetFrameSinceStartup(frames int64) {
	w.frameSinceStartup = frames
}

func (w *Workspace) GetFrameSinceStartup() int64 {
	return w.frameSinceStartup
}

/**
  uses the behavior tree in the cache, if not loaded yet, it loads the behavior tree first
*/
func (w *Workspace) CreateBehaviorTreeTask(relativePath string) *BehaviorTreeTask {
	if _, exist := w.behaviorTrees[relativePath]; !exist {
		err := w.Load(relativePath, false)
		if err != nil {
			BTGLog.Warnf("Workspace.CreateBehaviorTreeTask(%v) err:%v", relativePath, err)
			return nil
		}
	}
	var task BehaviorTask
	if tree, exist := w.behaviorTrees[relativePath]; exist {
		task = CreateAndInitTask(tree)
	}
	if treeTask, ok := task.(*BehaviorTreeTask); ok {
		if v, exist := w.allBehaviorTreeTasks[relativePath]; exist {
			v.bts = append(v.bts, treeTask)
		} else {
			w.allBehaviorTreeTasks[relativePath] = &BTGItem{bts: []*BehaviorTreeTask{treeTask}}
		}
		return treeTask
	}
	return nil
}

func (w *Workspace) DestroyBehaviorTreeTask(task *BehaviorTreeTask, a *Agent) {
	if task == nil || a == nil {
		return
	}

	relativePath := task.GetName()
	if item, exist := w.allBehaviorTreeTasks[relativePath]; exist {
		for i, v := range item.bts {
			if v == task {
				if i == 0 {
					item.bts = item.bts[1:]
				} else if i == len(item.bts)-1 {
					item.bts = item.bts[:len(item.bts)-1]
				} else {
					temp := item.bts[:i]
					temp = append(temp, item.bts[i+1:]...)
					item.bts = temp
				}
				break
			}
		}
		for i, v := range item.agents {
			if v == a {
				if i == 0 {
					item.agents = item.agents[1:]
				} else if i == len(item.agents)-1 {
					item.agents = item.agents[:len(item.agents)-1]
				} else {
					temp := item.agents[:i]
					temp = append(temp, item.agents[i+1:]...)
					item.agents = temp
				}
				break
			}
		}
	}
}

func (w *Workspace) Update() {
	if w.isExecAgents {
		ExecAgents(0)
	}
}
