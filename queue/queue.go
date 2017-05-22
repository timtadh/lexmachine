package queue

type Queue struct {
	list []uint32
	set  []uint32
}

func New(programLength int) *Queue {
	self := new(Queue)
	self.list = make([]uint32, 0, 10)
	self.set = make([]uint32, programLength)
	return self
}

func (self *Queue) Empty() bool { return len(self.list) <= 0 }

func (self *Queue) Has(pc uint32) bool {
	idx := self.set[pc]
	return idx < uint32(len(self.list)) && self.list[idx] == pc
}

func (self *Queue) Clear() {
	self.list = self.list[:0]
}

func (self *Queue) Push(pc uint32) {
	if self.Has(pc) {
		return
	}
	self.set[pc] = uint32(len(self.list))
	self.list = append(self.list, pc)
}

func (self *Queue) Pop() uint32 {
	pc := self.list[len(self.list)-1]
	self.list = self.list[:len(self.list)-1]
	return pc
}
