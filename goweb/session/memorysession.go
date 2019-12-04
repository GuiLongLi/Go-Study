package session

import (
	"container/list"
	"sync"
	"time"
)

type ProviderStruct struct {
	lock sync.Mutex  //用来锁
	sessions map[string]*list.Element //用来存储在内存
	list *list.List //用来做gc
}

var ps = &ProviderStruct{list:list.New()}

type SessionStore struct {
	sid string //session id 唯一标识
	timeAccessed time.Time //最后访问时间
	value map[interface{}]interface{} //session 里面存储的值
}

func (st *SessionStore) Set(key,value interface{}) error{
	st.value[key] = value
	ps.SessionUpdate(st.sid)
	return nil
}

func (st *SessionStore) Get(key interface{}) interface{}{
	ps.SessionUpdate(st.sid)
	if v,ok := st.value[key]; ok{
		return v
	}else{
		return nil
	}
	return nil
}

func (st *SessionStore) Delete(key interface{}) error{
	delete(st.value, key)
	ps.SessionUpdate(st.sid)
	return nil
}

func (st *SessionStore) SessionID() string{
	return st.sid
}

func (ps *ProviderStruct) SessionInit(sid string) (Session,error){
	ps.lock.Lock()
	defer ps.lock.Unlock()
	v := make(map[interface{}]interface{},0)
	newsess := &SessionStore{sid:sid,timeAccessed:time.Now(),value:v}
	element := ps.list.PushBack(newsess)
	ps.sessions[sid] = element
	return newsess,nil
}

func (ps *ProviderStruct) SessionRead(sid string) (Session,error){
	if element,ok := ps.sessions[sid];ok{
		return element.Value.(*SessionStore),nil
	}else{
		sess,err := ps.SessionInit(sid)
		return sess,err
	}
	return nil,nil
}

func (ps *ProviderStruct) SessionDestroy(sid string) error{
	if element,ok := ps.sessions[sid];ok{
		delete(ps.sessions,sid)
		ps.list.Remove(element)
		return nil
	}
	return nil
}

func (ps *ProviderStruct) SessionGC(maxlifetime int64){
	ps.lock.Lock()
	defer ps.lock.Unlock()
	for{
		element := ps.list.Back()
		if element == nil{
			break
		}
		if (element.Value.(*SessionStore).timeAccessed.Unix() + maxlifetime) < time.Now().Unix(){
			ps.list.Remove(element)
			delete(ps.sessions,element.Value.(*SessionStore).sid)
		}else{
			break
		}
	}
}

func (ps *ProviderStruct) SessionUpdate(sid string) error{
	ps.lock.Lock()
	defer ps.lock.Unlock()
	if element,ok := ps.sessions[sid];ok{
		element.Value.(*SessionStore).timeAccessed = time.Now()
		ps.list.MoveToFront(element)
		return nil
	}
	return nil
}

func init(){
	ps.sessions = make(map[string]*list.Element,0)
	Register("memory",ps)
}