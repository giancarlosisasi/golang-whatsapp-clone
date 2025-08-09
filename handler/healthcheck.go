package handler

import (
	"net/http"
)

func (h *Handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.methodNotAllowedResponse(w, r)
		return
	}

	data := envelop{
		"status": "available",
		"system_info": map[string]string{
			"environment": h.appConfig.AppEnv,
		},
	}

	err := h.writeJson(w, http.StatusOK, data, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}
