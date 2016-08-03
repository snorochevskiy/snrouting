package snweb

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"time"
)

const SID_COOKIE_NAME = "SID"

var sessionMap map[string]*Session

func init() {
	sessionMap = make(map[string]*Session)
}

type Session struct {
	Sid     string
	User    interface{} //*UserEntity
	Created time.Time
}

type UserRenderInfo struct {
	LoggedIn bool
	Details  interface{}
}

func (session *Session) IsLoggedIn() bool {
	return session.User != nil
}

func (session *Session) GetUserRenderInfo() *UserRenderInfo {
	u := new(UserRenderInfo)
	if session.User == nil {
		u.LoggedIn = false
	} else {
		u.LoggedIn = true
		u.Details = session.User
	}
	return u
}

type SessionManagerService struct {
	sessionMap map[string]*Session
}

func InitSession(w http.ResponseWriter, userInfo interface{}) { //*UserEntity) {
	generatedSid := generateSid()
	session := &Session{Sid: generatedSid}
	session.User = userInfo

	sessionMap[generatedSid] = session

	cookie := &http.Cookie{Name: SID_COOKIE_NAME, Value: session.Sid, MaxAge: 0}
	http.SetCookie(w, cookie)
}

func ClearSession(r *http.Request, w http.ResponseWriter) {
	cookie, err := r.Cookie(SID_COOKIE_NAME)
	if err != nil {
		return
	}
	delete(sessionMap, cookie.Value)

	http.SetCookie(w, &http.Cookie{Name: SID_COOKIE_NAME, Value: "", Path: "/", MaxAge: -1})
}

func generateSid() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func GetSessionForRequest(r *http.Request) *Session {

	cookie, err := r.Cookie(SID_COOKIE_NAME)
	if err != nil {
		return new(Session)
	}
	sid := cookie.Value

	session, contains := sessionMap[sid]
	if contains {
		return session
	} else {
		// Outdated cookie
		return new(Session)
	}
}

func (session *Session) SetCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{Name: SID_COOKIE_NAME, Value: session.Sid, MaxAge: 0}
	http.SetCookie(w, cookie)
}

