package queue

type Queue struct {
	list []uint32
	set  map[uint32]struct{}
}

func New() *Queue {
	self := new(Queue)
	self.list = make([]uint32, 0, 10)
	self.set = make(map[uint32]struct{})
	return self
}

func (self *Queue) Empty() bool { return len(self.list) <= 0 }

func (self *Queue) Push(pc uint32) {
	if _, ok := self.set[pc]; ok {
		return
	}
	self.set[pc] = struct{}{}
	self.list = append(self.list, pc)
}

func (self *Queue) Pop() uint32 {
	pc := self.list[len(self.list)-1]
	self.list = self.list[:len(self.list)-1]
	delete(self.set, pc)
	return pc
}
