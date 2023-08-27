package server

import "github.com/olahol/melody"


func newWsServer() *melody.Melody{
	m:= melody.New()
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.BroadcastFilter(msg, func(q *melody.Session) bool {
			return q.Request.URL.Path == s.Request.URL.Path
		})
	})
	return m
}
