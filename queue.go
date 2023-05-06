package vfacore

import "sync"

var qu queue

type queue struct {
	rw  sync.RWMutex
	Map map[int64]queueMsg
}
type queueMsg struct {
	Chan  chan int
	MsgId int64
}

func (q *queue) AddKey(key int64) {
	q.rw.Lock()
	defer q.rw.Unlock()
	msg := queueMsg{}
	msg.Chan = make(chan int)
	q.Map[key] = msg
}

func (q *queue) IssetKey(key int64) bool {
	q.rw.Lock()
	defer q.rw.Unlock()
	_, ok := q.Map[key]
	return ok
}

func (q *queue) RemoveKey(key int64) {
	q.rw.Lock()
	defer q.rw.Unlock()
	delete(q.Map, key)
	return
}
func (q *queue) GetMsg(key int64) queueMsg {
	q.rw.Lock()
	defer q.rw.Unlock()
	return q.Map[key]
}

func (q *queue) SetMsgId(key, msgid int64) {
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

// инициализация очередей
func InitQ() {
	var q queue
	q.Map = make(map[int64]queueMsg)
	qu = q

}
