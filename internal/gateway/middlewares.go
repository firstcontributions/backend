package gateway

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/firstcontributions/backend/internal/gateway/session"
	"github.com/firstcontributions/backend/internal/models/usersstore"
	"github.com/firstcontributions/backend/internal/storemanager"
	"github.com/google/uuid"
)

// HandleSession checks session data
func (s *Server) HandleSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, err := s.getConetxWithSession(r)
		if err != nil {
			ErrorResponse(err, w)
			return
		}
		r = r.WithContext(
			storemanager.ContextWithStore(
				ctx,
				s.Store,
			),
		)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) getConetxWithSession(r *http.Request) (context.Context, *Error) {
	cookie, err := r.Cookie("fc_session")
	if err != nil {
		log.Print("error on reading cookie ", err)
		return r.Context(), nil
	}
	cookieValue := make(map[string]string)

	err = s.CookieManager.Decode("fc_session", cookie.Value, &cookieValue)
	if err != nil {
		log.Print("error on decoding cookie ", err)
		return nil, ErrInternalServerError()
	}
	var sessionData session.MetaData
	if err := s.SessionManager.Get(r.Context(), cookieValue["id"], &sessionData); err != nil {
		log.Print("error on getting session ", err)
		return r.Context(), nil
	}
	log.Println("with cookie", sessionData.Handle())
	return session.WithContext(r.Context(), &sessionData), nil
}

func (s *Server) setSession(w http.ResponseWriter, r *http.Request, profile *usersstore.User) (string, error) {

	sessionData := session.NewMetaData(profile)
	sessionID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	if err := s.SessionManager.Set(r.Context(), sessionID.String(), sessionData); err != nil {
		return "", err
	}
	cookieData := map[string]string{
		"id": sessionID.String(),
	}
	encoded, err := s.CookieManager.Encode("fc_session", cookieData)
	if err != nil {
		return "", err
	}
	cookie := &http.Cookie{
		Name:    "fc_session",
		Value:   encoded,
		Expires: time.Now().Add(time.Duration(*s.SessionTTLDays) * time.Hour * 24),
		Path:    "/",
		Domain:  "opensource.forum",
	}
	http.SetCookie(w, cookie)
	return sessionID.String(), nil
}

func (s *Server) SessionHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := r.Context().Value(session.CxtKeySession)
		JSONResponse(w, http.StatusOK, session)
	})
}
