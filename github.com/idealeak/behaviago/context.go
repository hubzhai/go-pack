package behaviago

import (
	"sort"
)

var BTGContexts = make(map[int]*Context)

// ============================================================================
type PriorityAgentsItem struct {
	priority int
	agents   map[int]*Agent
}

// ============================================================================
type SortablePriorityAgents []*PriorityAgentsItem

func (self SortablePriorityAgents) Len() int {
	return len(self)
}
func (self SortablePriorityAgents) Less(i, j int) bool {
	return self[i].priority < self[j].priority
}
func (self SortablePriorityAgents) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

/// The Context class
// ============================================================================
type Context struct {
	Id                 int
	bCreatedByMe       bool
	bExecuting         bool
	namedAgents        map[string]*Agent
	agents             []*PriorityAgentsItem
	agentsMap          map[int]*PriorityAgentsItem
	delayAddedAgents   []*Agent
	delayRemovedAgents []*Agent
}

func newContext(id int) *Context {
	c := &Context{
		Id:          id,
		namedAgents: make(map[string]*Agent),
		agentsMap:   make(map[int]*PriorityAgentsItem),
	}
	return c
}

func GetContext(contextId int) *Context {
	if c, exist := BTGContexts[contextId]; exist {
		return c
	} else {
		c := newContext(contextId)
		if c != nil {
			BTGContexts[contextId] = c
		}
		return c
	}
	return nil
}

func CleanupContext(contextId int) {
	if contextId == -1 {
		BTGContexts = make(map[int]*Context)
	} else {
		delete(BTGContexts, contextId)
	}
}

func ExecAgents(contextId int) {
	if contextId != 0 {
		c := GetContext(contextId)
		if c != nil {
			c.ExecAgents()
		}
	} else {
		for _, c := range BTGContexts {
			c.ExecAgents()
		}
	}
}

func (c *Context) AddAgent(a *Agent) {
	if a == nil {
		return
	}
	if c.bExecuting {
		c.delayAddedAgents = append(c.delayAddedAgents, a)
	} else {
		c.addAgent_(a)
	}
}

func (c *Context) RemoveAgent(a *Agent) {
	if a == nil {
		return
	}
	if c.bExecuting {
		c.delayRemovedAgents = append(c.delayRemovedAgents, a)
	} else {
		c.removeAgent_(a)
	}
}

func (c *Context) GetInstance(agentInstanceName string) *Agent {
	if a, exist := c.namedAgents[agentInstanceName]; exist {
		return a
	}
	return nil
}

func (c *Context) BindInstance(agentInstanceName string, a *Agent) bool {
	if _, exist := c.namedAgents[agentInstanceName]; exist {
		BTGLog.Warnf("%v the name has been bound to an instance already!", agentInstanceName)
		return false
	}
	c.namedAgents[agentInstanceName] = a
	return true
}

func (c *Context) UnbindInstance(agentInstanceName string) bool {
	if _, exist := c.namedAgents[agentInstanceName]; exist {
		delete(c.namedAgents, agentInstanceName)
		return true
	}
	return false
}

func (c *Context) DelayProcessingAgents() {
	if len(c.delayAddedAgents) != 0 {
		for _, a := range c.delayAddedAgents {
			c.addAgent_(a)
		}
		c.delayAddedAgents = nil
	}
	if len(c.delayRemovedAgents) != 0 {
		for _, a := range c.delayRemovedAgents {
			c.removeAgent_(a)
		}
		c.delayRemovedAgents = nil
	}
}

func (c *Context) ExecAgents() {
	if !BTGWorkspace.isExecAgents {
		return
	}
	c.bExecuting = true
	for _, item := range c.agents {
		for _, a := range item.agents {
			if a.IsActive {
				a.BTGExec()
			}
			if !BTGWorkspace.isExecAgents {
				break
			}
		}
	}
	if GetIdMask() != 0 {
		c.LogStaticVariables()
	}
	c.bExecuting = false
	c.DelayProcessingAgents()
}

func (c *Context) CleanupInstances() {
	c.namedAgents = make(map[string]*Agent)
}

func (c *Context) LogStaticVariables() {

}

func (c *Context) addAgent_(a *Agent) {
	id := a.Id
	priority := a.Priority
	item := c.findPriorityAgents(priority)
	if item == nil {
		item = &PriorityAgentsItem{
			priority: priority,
			agents:   make(map[int]*Agent),
		}
		c.agents = append(c.agents, item)
		c.agentsMap[priority] = item
		item.agents[id] = a
		sort.Sort(SortablePriorityAgents(c.agents))
	} else {
		item.agents[id] = a
	}
}

func (c *Context) removeAgent_(a *Agent) {
	id := a.Id
	priority := a.Priority
	item := c.findPriorityAgents(priority)
	if item != nil {
		delete(item.agents, id)
	}
}

func (c *Context) findPriorityAgents(priority int) *PriorityAgentsItem {
	if item, exist := c.agentsMap[priority]; exist {
		return item
	}
	return nil
}
