package queue

import "container/list"

type Queue struct {
	list *list.List
	set  map[uint32]bool
}

func New() *Queue {
	self := new(Queue)
	self.list = list.New()
	self.set = make(map[uint32]bool)
	return self
}

func (self *Queue) Empty() bool { return self.list.Len() <= 0 }

func (self *Queue) Push(pc uint32) {
	if _, ok := self.set[pc]; ok {
		return
	}
	self.set[pc] = true
	self.list.PushBack(pc)
}

func (self *Queue) Pop() uint32 {
	e := self.list.Front()
	pc, _ := e.Value.(uint32)
	self.list.Remove(e)
	delete(self.set, pc)
	return pc
}
