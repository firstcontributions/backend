package gateway

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/firstcontributions/backend/internal/gateway/session"
	"github.com/firstcontributions/backend/internal/profile/proto"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// HandleSession checks session data
func (s *Server) HandleSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("fc_session")
		if err != nil {
			ErrorResponse(ErrUnauthorized(), w)
			return
		}
		cookieValue := make(map[string]string)

		err = s.CookieManager.Decode("fc_session", cookie.Value, &cookieValue)
		if err != nil {
			log.Print("error on decoding cookie ", err)
			ErrorResponse(ErrInternalServerError(), w)
			return
		}
		var sessionData session.MetaData
		if err := s.SessionManager.Get(r.Context(), cookieValue["id"], &sessionData); err != nil {
			if errors.Is(err, redis.Nil) {
				ErrorResponse(ErrUnauthorized(), w)
				return
			}
			log.Print("error on getting session ", err)
			ErrorResponse(ErrInternalServerError(), w)
			return
		}

		if sessionData.Handle == "" {
			ErrorResponse(ErrUnauthorized(), w)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), session.CxtKeySession, sessionData))
		log.Println(r.Context().Value(session.CxtKeySession))
		// Our middleware logic goes here...
		next.ServeHTTP(w, r)
	})
}

func (s *Server) setSession(w http.ResponseWriter, r *http.Request, profile *proto.Profile) error {

	sessionData := session.MetaData{
		UserID: profile.GetUuid(),
		Handle: profile.GetHandle(),
	}
	sessionID, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	if err := s.SessionManager.Set(r.Context(), sessionID.String(), sessionData); err != nil {
		return err
	}
	cookieData := map[string]string{
		"id": sessionID.String(),
	}
	encoded, err := s.CookieManager.Encode("fc_session", cookieData)
	if err != nil {
		return err
	}
	cookie := &http.Cookie{
		Name:    "fc_session",
		Value:   encoded,
		Expires: time.Now().Add(time.Duration(*s.SessionTTLDays) * time.Hour * 24),
		Path:    "/",
	}
	http.SetCookie(w, cookie)
	return nil
}

func (s *Server) SessionHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := r.Context().Value(session.CxtKeySession)
		JSONResponse(w, http.StatusOK, session)
	})
}
