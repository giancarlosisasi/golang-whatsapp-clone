package auth

import (
	"context"
	"fmt"
	"golang-whatsapp-clone/config"
	db "golang-whatsapp-clone/database/gen"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthHandlers struct {
	appConfig    *config.AppConfig
	oauthService *OAuthService
	jwtService   *JWTService
	queries      *db.Queries
}

func NewAuthHandlers(appConfig *config.AppConfig, oauthService *OAuthService, jwtService *JWTService, queries *db.Queries) *AuthHandlers {
	return &AuthHandlers{
		appConfig:    appConfig,
		oauthService: oauthService,
		jwtService:   jwtService,
		queries:      queries,
	}
}

var oauthStateCookieName = "whatsappgio_oauth_state"
var oauthClientTypeCookieName = "whatsappgio_oauth_client_type"

func (h *AuthHandlers) GoogleLogin(c *fiber.Ctx) error {
	// make sure to clear any previous invalid cookie
	h.clearStateCookies(c)

	log.Printf("Original URL: %s\n", c.OriginalURL())

	// Get client type (web or mobile)
	clientType := c.Query("client_type", "web")
	log.Printf("Client type: %s --\n", clientType)

	state := ""
	// generate auth url and state
	authURL, stateForWeb, err := h.oauthService.GenerateAuthURL()
	if err != nil {
		log.Printf("Failed to generate auth url: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate auth url",
		})
	}
	// mobile will send us a state as a query param
	stateForMobile := c.Query("state")

	if clientType == "web" {
		state = stateForWeb
	} else if clientType == "mobile" {
		state = stateForMobile
	}

	if state == "" {
		log.Print("Failed to get the state value from url\n")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get the state value",
		})

	}

	log.Printf("state value: %s", state)

	log.Printf("AuthURL: %s\n", authURL)

	// store state and client type in cookies for validation
	h.setStateCookies(c, state, clientType)

	// redirect to google page
	return c.Redirect(authURL, fiber.StatusTemporaryRedirect)
}

func (h *AuthHandlers) GoogleCallback(c *fiber.Ctx) error {
	// verify state parameter
	state := c.Query("state")
	storedState := c.Cookies(oauthStateCookieName)
	clientType := c.Cookies(oauthClientTypeCookieName, "web")

	if state == "" || state != storedState {
		fmt.Printf("state: %s and storedState: %s\n", state, storedState)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid state parameter",
		})
	}

	// clear state cookies
	h.clearStateCookies(c)

	// exchange code for token
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing authorization code",
		})
	}

	ctx := context.Background()
	token, err := h.oauthService.ExchangeCode(ctx, code)
	if err != nil {
		log.Printf("Failed to exchange token: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to exchange token",
		})
	}

	// get user info from google
	googleUser, err := h.oauthService.GetUserInfo(ctx, token)
	if err != nil {
		log.Printf("Failed to get user info: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get user information",
		})
	}

	user, err := h.queries.UpsertUserByGoogleAuthSafe(ctx, db.UpsertUserByGoogleAuthSafeParams{
		Name:      pgtype.Text{String: googleUser.Name, Valid: true},
		GoogleID:  pgtype.Text{String: googleUser.ID, Valid: true},
		Email:     googleUser.Email,
		AvatarUrl: pgtype.Text{String: googleUser.Picture, Valid: true},
	})
	if err != nil {
		log.Printf("failed to upsert user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create/update user",
		})
	}

	// generate jwt
	jwtToken, err := h.jwtService.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		log.Printf("failed to generate jwt: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate token",
		})
	}

	return h.handleAuthSuccess(c, jwtToken, clientType)
}

func (h *AuthHandlers) Logout(c *fiber.Ctx) error {
	c.ClearCookie(h.appConfig.CookieName)

	return c.JSON(fiber.Map{
		"message": "logged out successfully",
	})
}

func (h *AuthHandlers) setStateCookies(c *fiber.Ctx, state string, clientType string) {
	isSecure := h.appConfig.AppEnv == "production"

	c.Cookie(
		&fiber.Cookie{
			// Name:     h.appConfig.CookieName,
			Name:     oauthStateCookieName,
			Value:    state,
			HTTPOnly: true,
			Secure:   isSecure,
			SameSite: "Lax",
			MaxAge:   300, // 5 minutes
			Path:     "/",
		},
	)

	c.Cookie(&fiber.Cookie{
		Name:     oauthClientTypeCookieName,
		Value:    clientType,
		HTTPOnly: true,
		Secure:   isSecure,
		SameSite: "Lax",
		MaxAge:   300,
		Path:     "/",
	})
}

func (h *AuthHandlers) clearStateCookies(c *fiber.Ctx) {
	isSecure := h.appConfig.AppEnv == "production"

	c.Cookie(
		&fiber.Cookie{
			// Name:     h.appConfig.CookieName,
			Name:     oauthStateCookieName,
			Value:    "",
			HTTPOnly: true,
			Secure:   isSecure,
			SameSite: "Lax",
			MaxAge:   -1, // 5 minutes
			Path:     "/",
			Expires:  time.Now().Add(-1 * time.Hour), // set expiry in the past
		},
	)

	c.Cookie(&fiber.Cookie{
		Name:     oauthClientTypeCookieName,
		Value:    "",
		HTTPOnly: true,
		Secure:   isSecure,
		SameSite: "Lax",
		MaxAge:   -1,
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
	})
}

func (h *AuthHandlers) handleAuthSuccess(c *fiber.Ctx, jwtToken string, clientType string) error {
	if clientType == "mobile" {
		// redirect to mobile app with token
		redirectURL := fmt.Sprintf("%s://auth/%s/", h.appConfig.MobileAppSchema, jwtToken)
		return c.Redirect(redirectURL, fiber.StatusTemporaryRedirect)
	}

	isSecure := h.appConfig.AppEnv == "production"
	c.Cookie(&fiber.Cookie{
		Name:     h.appConfig.CookieName,
		Value:    jwtToken,
		HTTPOnly: true,
		Secure:   isSecure,
		SameSite: "Strict",
		MaxAge:   7 * 24 * 60 * 60, // 7 days
		Path:     "/",
	})

	return c.Redirect("/chats", fiber.StatusTemporaryRedirect)
}
