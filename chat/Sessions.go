package main

import (
	"crypto/rand"
	"errors"
	"math/big"
	"time"
)

const LARGE_NUMBER_1 uint64 = 100000000000000
const LARGE_NUMBER_2 = 899999999999999
const SESSION_LIFETIME_SECONDS time.Duration = 30 * 60000000 // 30 minutes in nanoseconds

type Session struct {
	UserID  uint64
	Timeout time.Time
}

type SessionManager struct {
	sessions map[uint64]Session
}

func NewSessionManager() *SessionManager {
	var sm SessionManager
	sm.sessions = make(map[uint64]Session)

	return &sm
}

func (sm *SessionManager) NewSession(userId uint64) (uint64, error) {
	var sessionId uint64

	for exists := true; exists; {
		addition, err := rand.Int(rand.Reader, big.NewInt(LARGE_NUMBER_2))
		if err != nil {
			return 0, err
		}

		sessionId = LARGE_NUMBER_1 + addition.Uint64()
		_, exists = sm.sessions[sessionId]
	}

	sm.sessions[sessionId] = Session{
		UserID:  userId,
		Timeout: time.Now().Add(SESSION_LIFETIME_SECONDS),
	}

	return sessionId, nil
}

func (sm *SessionManager) GetUserIDFromSession(sessionId uint64) (uint64, error) {
	session, exists := sm.sessions[sessionId]
	if !exists {
		return 0, errors.New("Not logged in")
	}

	if session.Timeout.After(time.Now()) {
		delete(sm.sessions, sessionId)
		return 0, errors.New("Session expired")
	}

	return session.UserID, nil
}
