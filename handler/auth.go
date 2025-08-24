package handler

import (
	"context"
	"fmt"
	"golang-whatsapp-clone/auth"
	db "golang-whatsapp-clone/database/gen"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

var oauthStateCookieName = "whatsappgio_oauth_state"
var oauthClientTypeCookieName = "whatsappgio_oauth_client_type"

func (h *Handler) GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {

	// Get client type (web or mobile)
	queryParams := r.URL.Query()
	// clientType := c.Query("client_type", "web")
	clientType := queryParams.Get("client_type")
	if clientType == "" {
		clientType = "web"
	}

	h.logger.Info().Msgf("Initializing oauth google flow from url: %s", r.URL.String())
	h.logger.Info().Msgf(">> Client type is: %s", clientType)

	state := ""
	// generate auth url and state
	authURL, state, err := h.oauthService.GenerateAuthURL()
	if err != nil {
		log.Printf("Failed to generate auth url: %v\n", err)
		// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 	"error": "Failed to generate auth url",
		// })

		h.logger.Error().Msg("failed to generate auth url")
		h.ServerErrorResponse(w, r, err)
		return
	}
	// mobile will send us a state as a query param
	// stateForMobile := c.Query("state")

	// if clientType == "web" {
	// 	state = stateForWeb
	// } else if clientType == "mobile" {
	// 	state = stateForMobile
	// }

	// if state == "" {
	// 	log.Print("Failed to get the state value from url\n")
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": "Failed to get the state value",
	// 	})

	// }

	// log.Printf("state value: %s", state)

	// log.Printf("AuthURL: %s\n", authURL)

	h.logger.Info().Msgf(">> State is: %s", state)

	// store state and client type in cookies for validation
	h.setStateCookies(w, r, state, clientType)

	// redirect to google page
	// return c.Redirect(authURL, fiber.StatusTemporaryRedirect)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (h *Handler) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info().Msgf("Callback run with url: %s", r.URL.String())

	queryParams := r.URL.Query()

	// verify state parameter
	state := queryParams.Get("state")
	// storedState := c.Cookies(oauthStateCookieName)
	storedState := ""
	storedStateCookie, err := r.Cookie(oauthStateCookieName)
	if err == nil {
		storedState = storedStateCookie.Value
	}
	// clientType := c.Cookies(oauthClientTypeCookieName, "web")
	clientType := "web"
	clientTypeCookie, err := r.Cookie(oauthClientTypeCookieName)
	if err == nil {
		clientType = clientTypeCookie.Value
	}

	h.logger.Info().Msgf(">>> Callback - client type is %s and state is %s", clientType, state)

	if state == "" || state != storedState {
		log.Printf("state: %s and storedState: %s\n", state, storedState)
		h.clearStateCookies(w, r)
		// return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		// 	"error": "invalid state parameter",
		// })
		h.errorResponse(w, r, http.StatusBadRequest, "invalid state parameter")
		return
	}

	// clear state cookies
	h.clearStateCookies(w, r)

	// exchange code for token
	code := queryParams.Get("code")
	if code == "" {
		// return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		// 	"error": "missing authorization code",
		// })
		h.errorResponse(w, r, http.StatusBadRequest, "missing authorization code")
		return
	}

	ctx := context.Background()
	token, err := h.oauthService.ExchangeCode(ctx, code)
	if err != nil {
		log.Printf("Failed to exchange token: %v\n", err)
		// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 	"error": "failed to exchange token",
		// })
		h.logger.Error().Msg("failed to exchange token")
		h.ServerErrorResponse(w, r, err)
		return
	}

	// get user info from google
	googleUser, err := h.oauthService.GetUserInfo(ctx, token)
	if err != nil {
		// log.Printf("Failed to get user info: %v\n", err)
		// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 	"error": "failed to get user information",
		// })
		h.logger.Error().Msgf("failed to get user information: %s", err)
		h.ServerErrorResponse(w, r, err)
		return

	}

	user, err := h.dbQueries.UpsertUserByGoogleAuthSafe(ctx, db.UpsertUserByGoogleAuthSafeParams{
		Name:      pgtype.Text{String: googleUser.Name, Valid: true},
		GoogleID:  pgtype.Text{String: googleUser.ID, Valid: true},
		Email:     googleUser.Email,
		AvatarUrl: pgtype.Text{String: googleUser.Picture, Valid: true},
	})
	if err != nil {
		// log.Printf("failed to upsert user: %v", err)
		// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 	"error": "failed to create/update user",
		// })
		h.logger.Error().Msg("failed to create/update the user")
		h.ServerErrorResponse(w, r, err)
		return
	}

	// generate jwt
	jwtToken, err := h.jwtService.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		log.Printf("failed to generate jwt: %v", err)
		// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 	"error": "failed to generate token",
		// })
		h.ServerErrorResponse(w, r, err)
		return
	}

	h.handleAuthSuccess(w, r, jwtToken, clientType)
}

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	expiredAuthCookie := h.createAuthCookie(w, r, "", -1, time.Now().Add(-1*time.Hour))
	http.SetCookie(w, expiredAuthCookie)

	// return c.JSON(fiber.Map{
	// 	"message": "logged out successfully",
	// })
	err := h.writeJson(w, http.StatusOK, envelop{"message": "logged out successfully"}, nil)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

// This middleware will try to attach the user data to the context request in case an authentication header exists
func (h *Handler) AuthenticateUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var tokenString string

			// try to get token from Authorization header first (for mobile)
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenString = authHeader[7:]
			} else {
				// try to get token from cookie (web)
				tokenStringCookie, err := r.Cookie(h.appConfig.CookieName)
				if err == nil {
					tokenString = tokenStringCookie.Value
				}
			}

			if tokenString == "" {
				// return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				// 	"error": "Missing authentication token",
				// })
				next.ServeHTTP(w, r)
				return
			}

			// validate jwt
			claims, err := h.jwtService.ValidateToken(tokenString)
			if err != nil {
				// return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				// 	"error": "invalid authentication token",
				// })
				h.errorResponse(w, r, http.StatusUnauthorized, "invalid authentication token")
				return
			}

			// store user info in context for graphql resolvers
			// c.Locals("user_id", claims.UserID)
			// c.Locals("user_email", claims.Email)
			ctx := r.Context()
			if claims.UserID != "" && claims.Email != "" {
				ctx = auth.WithUserContext(ctx, claims.UserID, claims.Email)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		},
	)
}

func (h *Handler) setStateCookies(w http.ResponseWriter, r *http.Request, state string, clientType string) {
	isSecure := h.appConfig.AppEnv == "production"

	oauthStateCookie := &http.Cookie{
		Name:     oauthStateCookieName,
		Value:    state,
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   300, // 5 minutes
		Path:     "/",
	}

	http.SetCookie(w, oauthStateCookie)

	oauthClientTypeCookie := &http.Cookie{
		Name:     oauthClientTypeCookieName,
		Value:    clientType,
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   300,
		Path:     "/",
	}
	http.SetCookie(w, oauthClientTypeCookie)
}

func (h *Handler) clearStateCookies(w http.ResponseWriter, r *http.Request) {
	isSecure := h.appConfig.AppEnv == "production"

	oauthStateCookie := &http.Cookie{
		// Name:     h.appConfig.CookieName,
		Name:     oauthStateCookieName,
		Value:    "",
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1, // 5 minutes
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour), // set expiry in the past
	}
	http.SetCookie(w, oauthStateCookie)

	oauthClientTypeCookie := &http.Cookie{
		Name:     oauthClientTypeCookieName,
		Value:    "",
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
	}
	http.SetCookie(w, oauthClientTypeCookie)
}

func (h *Handler) handleAuthSuccess(
	w http.ResponseWriter, r *http.Request,
	jwtToken string, clientType string) {

	h.logger.Info().Msgf("Handle success with client type: %s", clientType)

	if clientType == "mobile" {
		h.logger.Info().Msg("Redirecting to mobile schema after successfully google oauth login...")
		// redirect to mobile app with token
		redirectURL := fmt.Sprintf("%sauth/%s/", h.appConfig.MobileAppSchema, jwtToken)
		// return c.Redirect(redirectURL, fiber.StatusTemporaryRedirect)
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
	}

	h.logger.Info().Msg("Redirecting to web url after successfully google oauth login...")

	isSecure := h.appConfig.AppEnv == "production"
	cookie := &http.Cookie{
		Name:     h.appConfig.CookieName,
		Value:    jwtToken,
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteDefaultMode,
		MaxAge:   7 * 24 * 60 * 60, // 7 days
		Path:     "/",
	}
	http.SetCookie(w, cookie)

	// return c.Redirect("/chats", fiber.StatusTemporaryRedirect)
	http.Redirect(w, r, "/chats", http.StatusTemporaryRedirect)
}

func (h *Handler) createAuthCookie(
	w http.ResponseWriter, r *http.Request,
	value string, maxAge int, expires time.Time,
) *http.Cookie {
	isSecure := h.appConfig.AppEnv == "production"
	cookie := &http.Cookie{
		Name:     h.appConfig.CookieName,
		Value:    value,
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   maxAge,
		Path:     "/",
		Expires:  expires,
	}

	// http.SetCookie(w, cookie)
	return cookie
}
