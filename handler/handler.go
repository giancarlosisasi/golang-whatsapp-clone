package handler

import (
	"encoding/json"
	"fmt"
	"golang-whatsapp-clone/auth"
	"golang-whatsapp-clone/config"
	db "golang-whatsapp-clone/database/gen"
	"maps"
	"net/http"

	"github.com/rs/zerolog"
)

type Handler struct {
	logger       *zerolog.Logger
	appConfig    *config.AppConfig
	dbQueries    *db.Queries
	oauthService *auth.OAuthService
	jwtService   *auth.JWTService
}

type envelop map[string]any

func NewHandler(logger *zerolog.Logger, appConfig *config.AppConfig, dbQueries *db.Queries, oauthService *auth.OAuthService, jwtService *auth.JWTService) *Handler {
	return &Handler{
		logger:       logger,
		appConfig:    appConfig,
		dbQueries:    dbQueries,
		oauthService: oauthService,
		jwtService:   jwtService,
	}
}

func (h *Handler) writeJson(w http.ResponseWriter, status int, data envelop, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	maps.Copy(w.Header(), headers)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(js)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s is not supported for this resource", r.Method)
	h.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (h *Handler) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.logger.Error().Msg(err.Error())

	message := "the server encountered a problem and could not process your request"
	h.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (h *Handler) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	data := envelop{"error": message}

	err := h.writeJson(w, status, data, nil)
	if err != nil {
		// fallback to internal server error
		h.logger.Err(err)
		w.WriteHeader(500)
	}
}

const (
	RedirectTypePermanent = iota
	RedirectTypeFound
	RedirectTypeSeeOther
	RedirectTypeTemporary
	RedirectTypeManual
)
