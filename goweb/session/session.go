package session

import (
	"sync"
	"time"
	"encoding/base64"
	"math/rand"
	"net/http"
	"net/url"
)
/*
session 管理设计
·全局session 管理器
·保证sessionid 的全局唯一性
·为每个客户关联一个session
·session 的存储（可以存储到内存、文件、数据库等）
·session 过期处理
*/

//session管理器
type Manager struct {
	cookieName string //cookie 名称
	lock sync.Mutex //cookie 锁
	provider Provider //cookie 提供者
	maxLifeTime int64 //cookie 最大生存时间
}

//抽象出Provider 接口，用来表示session 管理器底层存储结构
type Provider interface {
	//SessionInit 函数实现Session 的初始化，返回新的Session 变量
	SessionInit(sid string)(Session,error)
	//SessionRead 函数返回sid 代表的Session 变量，如果不存在，那么将以sid 为参数调用SessionInit 函数创建并返回一个新的Session 变量
	SessionRead(sid string)(Session,error)
	//SessionDestroy 函数用来销毁sid 对应的Session 变量
	SessionDestroy(sid string) error
	//SessionGC 根据maxLifeTime 来删除过期的数据
	SessionGC(maxLifeTime int64)
}

//Session 的处理基本就 设置值、读取值、删除值和获取sessionid
type Session interface {
	Set(key,value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionID() string
}

var provides = make(map[string]Provider)

//创建manager
func NewManager(provideName,cookieName string,maxLifeTime int64) (*Manager,error){
	provider,ok := provides[provideName]
	if !ok {
		panic("session: unknown provide "+provideName+"(forgotten import?)")
	}
	return &Manager{provider:provider,cookieName:cookieName,maxLifeTime:maxLifeTime},nil
}

//提供名称,注册一个session 提供者
//如果注册两次都是重复的，或者驱动是nil ， 返回错误
func Register(name string,provider Provider){
	if provider == nil{
		panic("session:Register provider is nil")
	}
	if _,dup := provides[name];dup{
		panic("session:Register called twice for provider "+name)
	}
	provides[name] = provider

}

//全局唯一的 SessionID
func (manager *Manager) sessionId() string{
	b := make([]byte,32)
	if _,err := rand.Read(b);err != nil{
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (manager *Manager) SessionStart(w http.ResponseWriter,r *http.Request)(session Session){
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie,err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == ""{
		sid := manager.sessionId()
		session,_ = manager.provider.SessionInit(sid)
		cookie := http.Cookie{
			Name:manager.cookieName,
			Value:url.QueryEscape(sid),
			Path:"/",
			HttpOnly:true,
			//设置httpOnly属性（说明：Cookie的HttpOnly属性，指示浏览器不要在除HTTP（和 HTTPS)请求之外暴露Cookie。
			// 一个有HttpOnly属性的Cookie，不能通过非HTTP方式来访问，例如通过调用JavaScript(例如，引用 document.cookie），
			// 因此，不可能通过跨域脚本（一种非常普通的攻击技术）来偷走这种Cookie。尤其是Facebook 和 Google 正在广泛地使用HttpOnly属性。）
			MaxAge: int(manager.maxLifeTime),
		}
		http.SetCookie(w,&cookie)
	}else{
		sid,_ := url.QueryUnescape(cookie.Value)
		session,_ = manager.provider.SessionRead(sid)
	}
	return
}

func (manager *Manager) SessionDestroy(w http.ResponseWriter,r *http.Request){
	cookie,err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == ""{
		return
	}else{
		manager.lock.Lock()
		defer manager.lock.Unlock()
		manager.provider.SessionDestroy(cookie.Value)
		expiration := time.Now()
		cookie := http.Cookie{
			Name:manager.cookieName,
			Path:"/",
			HttpOnly:true,
			Expires:expiration,
			MaxAge:-1,
		}
		http.SetCookie(w,&cookie)
	}
}

func (manager *Manager) GC(){
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionGC(manager.maxLifeTime)
	time.AfterFunc(time.Duration(manager.maxLifeTime),func(){
		manager.GC()
	})
}