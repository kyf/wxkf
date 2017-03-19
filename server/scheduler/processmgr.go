package main

import (
	"encoding/json"
	"sync"
)

type Process struct {
	Name    string
	IP      string
	ConnNum int32
}

type ProcessMgr struct {
	list []*Process
	sync.Mutex
}

func (this *ProcessMgr) Register(p *Process) {
	this.Lock()
	defer this.Unlock()

	this.list = append(this.list, p)
}

func (this *ProcessMgr) Unregister(p *Process) {
	this.Lock()
	defer this.Unlock()

	for index, _p := range this.list {
		if _p == p {
			copy(this.list[index:], this.list[index+1:])
			this.list[len(this.list)-1] = nil
			this.list = this.list[:len(this.list)-1]
			return
		}
	}
}

var defaultProMgr = &ProcessMgr{list: make([]*Process, 0)}

func getBestProcessor() *Process {
	var result *Process
	for _, p := range defaultProMgr.list {
		if p.IP == "" && p.Name == "" {
			continue
		}
		if result == nil {
			result = p
			continue
		}

		if result.ConnNum > p.ConnNum {
			result = p
		}
	}

	return result
}
