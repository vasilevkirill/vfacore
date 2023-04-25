package queue

import "sync"

var Q Queue

type Queue struct {
	rw  sync.RWMutex
	Map map[int64]Msg
}
type Msg struct {
	Chan  chan int
	MsgId int64
}

func (q *Queue) AddKey(key int64) {
	q.rw.Lock()
	defer q.rw.Unlock()
	msg := Msg{}
	msg.Chan = make(chan int)
	q.Map[key] = msg
}

func (q *Queue) IssetKey(key int64) bool {
	q.rw.Lock()
	defer q.rw.Unlock()
	_, ok := q.Map[key]
	return ok
}

func (q *Queue) RemoveKey(key int64) {
	q.rw.Lock()
	defer q.rw.Unlock()
	delete(q.Map, key)
	return
}
func (q *Queue) GetMsg(key int64) Msg {
	q.rw.Lock()
	defer q.rw.Unlock()
	return q.Map[key]
}

func (q *Queue) SetMsgId(key, msgid int64) {
	q.rw.Lock()
	defer q.rw.Unlock()
	_, ok := q.Map[key]
	if !ok {
		return
	}
	msg := q.Map[key]
	msg.MsgId = msgid
	q.Map[key] = msg

}

func InitQ() {
	var q Queue
	q.Map = make(map[int64]Msg)
	Q = q

}
