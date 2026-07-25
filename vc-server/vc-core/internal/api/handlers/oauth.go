package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"vc-server/vc-core/internal/config"
	"vc-server/vc-core/internal/db"
	"vc-server/vc-core/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const oauthStateCookie = "oauth_state"
const oauthProviderCookie = "oauth_provider"

type oauthProvidersResponse struct {
	GitHub bool `json:"github"`
	Google bool `json:"google"`
}

// HandleOAuthProviders returns which OAuth providers are configured.
func HandleOAuthProviders(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(oauthProvidersResponse{
			GitHub: cfg.GitHubClientID != "" && cfg.GitHubClientSecret != "",
			Google: cfg.GoogleClientID != "" && cfg.GoogleClientSecret != "",
		})
	}
}

// HandleOAuthStart redirects the browser to the provider consent screen.
func HandleOAuthStart(cfg *config.Config, provider string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state, err := randomState()
		if err != nil {
			http.Error(w, "Failed to start OAuth", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     oauthStateCookie,
			Value:    state,
			Path:     "/",
			HttpOnly: true,
			Secure:   r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   600,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     oauthProviderCookie,
			Value:    provider,
			Path:     "/",
			HttpOnly: true,
			Secure:   r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   600,
		})

		var authURL string
		switch provider {
		case "github":
			if cfg.GitHubClientID == "" || cfg.GitHubClientSecret == "" {
				http.Error(w, "GitHub OAuth is not configured", http.StatusServiceUnavailable)
				return
			}
			q := url.Values{}
			q.Set("client_id", cfg.GitHubClientID)
			q.Set("redirect_uri", oauthCallbackURL(cfg, "github"))
			q.Set("scope", "read:user user:email")
			q.Set("state", state)
			authURL = "https://github.com/login/oauth/authorize?" + q.Encode()
		case "google":
			if cfg.GoogleClientID == "" || cfg.GoogleClientSecret == "" {
				http.Error(w, "Google OAuth is not configured", http.StatusServiceUnavailable)
				return
			}
			q := url.Values{}
			q.Set("client_id", cfg.GoogleClientID)
			q.Set("redirect_uri", oauthCallbackURL(cfg, "google"))
			q.Set("response_type", "code")
			q.Set("scope", "openid email profile")
			q.Set("state", state)
			q.Set("access_type", "online")
			q.Set("prompt", "select_account")
			authURL = "https://accounts.google.com/o/oauth2/v2/auth?" + q.Encode()
		default:
			http.Error(w, "Unknown provider", http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, authURL, http.StatusFound)
	}
}

// HandleOAuthCallback exchanges the code, upserts the user, and redirects to the client with tokens.
func HandleOAuthCallback(cfg *config.Config, database *db.DB, provider string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientURL := cfg.ClientURL
		if clientURL == "" {
			clientURL = "http://localhost:5173"
		}

		fail := func(msg string) {
			http.Redirect(w, r, clientURL+"/login?error="+url.QueryEscape(msg), http.StatusFound)
		}

		stateCookie, err := r.Cookie(oauthStateCookie)
		if err != nil || stateCookie.Value == "" {
			fail("OAuth session expired. Please try again.")
			return
		}
		state := r.URL.Query().Get("state")
		code := r.URL.Query().Get("code")
		if state == "" || code == "" || state != stateCookie.Value {
			fail("Invalid OAuth state. Please try again.")
			return
		}
		if errMsg := r.URL.Query().Get("error"); errMsg != "" {
			fail("Sign-in was cancelled.")
			return
		}

		// Clear state cookies
		http.SetCookie(w, &http.Cookie{Name: oauthStateCookie, Value: "", Path: "/", MaxAge: -1})
		http.SetCookie(w, &http.Cookie{Name: oauthProviderCookie, Value: "", Path: "/", MaxAge: -1})

		var profile oauthProfile
		switch provider {
		case "github":
			profile, err = fetchGitHubProfile(cfg, code)
		case "google":
			profile, err = fetchGoogleProfile(cfg, code)
		default:
			fail("Unknown provider")
			return
		}
		if err != nil {
			fail("Could not verify your account with the provider.")
			return
		}

		userID, username, err := upsertOAuthUser(database, provider, profile)
		if err != nil {
			fail("Could not create your account.")
			return
		}

		accessToken, err := utils.GenerateToken(userID)
		if err != nil {
			fail("Could not create session.")
			return
		}
		refreshToken, err := utils.GenerateRefreshToken(userID)
		if err != nil {
			fail("Could not create session.")
			return
		}

		setAuthCookies(w, accessToken, refreshToken)

		// Hash fragment keeps tokens out of server access logs on the SPA host.
		frag := url.Values{}
		frag.Set("accessToken", accessToken)
		frag.Set("refreshToken", refreshToken)
		frag.Set("uid", userID)
		frag.Set("username", username)
		http.Redirect(w, r, clientURL+"/login/callback#"+frag.Encode(), http.StatusFound)
	}
}

type oauthProfile struct {
	ProviderID string
	Email      string
	Username   string
	AvatarURL  string
	Name       string
}

func oauthCallbackURL(cfg *config.Config, provider string) string {
	base := strings.TrimRight(cfg.PublicAPIURL, "/")
	if base == "" {
		host := cfg.ServerHost
		if host == "" {
			host = "localhost"
		}
		port := cfg.ServerPort
		if port == "" {
			port = "5000"
		}
		base = fmt.Sprintf("http://%s:%s", host, port)
	}
	return fmt.Sprintf("%s/auth/oauth/%s/callback", base, provider)
}

func randomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func fetchGitHubProfile(cfg *config.Config, code string) (oauthProfile, error) {
	form := url.Values{}
	form.Set("client_id", cfg.GitHubClientID)
	form.Set("client_secret", cfg.GitHubClientSecret)
	form.Set("code", code)
	form.Set("redirect_uri", oauthCallbackURL(cfg, "github"))

	req, err := http.NewRequest(http.MethodPost, "https://github.com/login/oauth/access_token", strings.NewReader(form.Encode()))
	if err != nil {
		return oauthProfile{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return oauthProfile{}, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var tokenResp struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil || tokenResp.AccessToken == "" {
		return oauthProfile{}, fmt.Errorf("github token exchange failed: %s", string(body))
	}

	userReq, _ := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	userReq.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	userReq.Header.Set("Accept", "application/vnd.github+json")
	userResp, err := http.DefaultClient.Do(userReq)
	if err != nil {
		return oauthProfile{}, err
	}
	defer userResp.Body.Close()
	var ghUser struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
		Name      string `json:"name"`
	}
	if err := json.NewDecoder(userResp.Body).Decode(&ghUser); err != nil {
		return oauthProfile{}, err
	}

	email := ghUser.Email
	if email == "" {
		email = fetchGitHubPrimaryEmail(tokenResp.AccessToken)
	}

	return oauthProfile{
		ProviderID: fmt.Sprintf("%d", ghUser.ID),
		Email:      email,
		Username:   ghUser.Login,
		AvatarURL:  ghUser.AvatarURL,
		Name:       ghUser.Name,
	}, nil
}

func fetchGitHubPrimaryEmail(token string) string {
	req, _ := http.NewRequest(http.MethodGet, "https://api.github.com/user/emails", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return ""
	}
	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email
		}
	}
	for _, e := range emails {
		if e.Verified {
			return e.Email
		}
	}
	if len(emails) > 0 {
		return emails[0].Email
	}
	return ""
}

func fetchGoogleProfile(cfg *config.Config, code string) (oauthProfile, error) {
	form := url.Values{}
	form.Set("client_id", cfg.GoogleClientID)
	form.Set("client_secret", cfg.GoogleClientSecret)
	form.Set("code", code)
	form.Set("grant_type", "authorization_code")
	form.Set("redirect_uri", oauthCallbackURL(cfg, "google"))

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", form)
	if err != nil {
		return oauthProfile{}, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var tokenResp struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil || tokenResp.AccessToken == "" {
		return oauthProfile{}, fmt.Errorf("google token exchange failed: %s", string(body))
	}

	req, _ := http.NewRequest(http.MethodGet, "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	userResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return oauthProfile{}, err
	}
	defer userResp.Body.Close()
	var gUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(userResp.Body).Decode(&gUser); err != nil {
		return oauthProfile{}, err
	}

	username := gUser.Email
	if at := strings.Index(username, "@"); at > 0 {
		username = username[:at]
	}
	username = sanitizeUsername(username)

	return oauthProfile{
		ProviderID: gUser.ID,
		Email:      gUser.Email,
		Username:   username,
		AvatarURL:  gUser.Picture,
		Name:       gUser.Name,
	}, nil
}

func sanitizeUsername(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			b.WriteRune(r)
		}
	}
	out := b.String()
	if out == "" {
		out = "player"
	}
	if len(out) > 24 {
		out = out[:24]
	}
	return out
}

func upsertOAuthUser(database *db.DB, provider string, profile oauthProfile) (userID, username string, err error) {
	collection := database.Collection("users")
	ctx := context.TODO()

	providerField := provider + "_id"
	filter := bson.M{providerField: profile.ProviderID}

	var existing bson.M
	err = collection.FindOne(ctx, filter).Decode(&existing)
	if err == nil {
		id := existing["_id"].(primitive.ObjectID).Hex()
		uname, _ := existing["username"].(string)
		return id, uname, nil
	}
	if err != mongo.ErrNoDocuments {
		return "", "", err
	}

	// Link by email if a password account already exists.
	if profile.Email != "" {
		err = collection.FindOne(ctx, bson.M{"email": profile.Email}).Decode(&existing)
		if err == nil {
			id := existing["_id"].(primitive.ObjectID)
			_, err = collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
				"$set": bson.M{
					providerField:  profile.ProviderID,
					"avatar_url":   profile.AvatarURL,
					"auth_provider": provider,
				},
			})
			if err != nil {
				return "", "", err
			}
			uname, _ := existing["username"].(string)
			return id.Hex(), uname, nil
		}
		if err != mongo.ErrNoDocuments {
			return "", "", err
		}
	}

	username = uniqueUsername(ctx, collection, profile.Username)
	doc := bson.M{
		"username":      username,
		"email":         profile.Email,
		providerField:   profile.ProviderID,
		"avatar_url":    profile.AvatarURL,
		"auth_provider": provider,
		"created_at":    time.Now(),
	}
	res, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return "", "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), username, nil
}

func uniqueUsername(ctx context.Context, collection *mongo.Collection, base string) string {
	base = sanitizeUsername(base)
	candidate := base
	for i := 0; i < 20; i++ {
		var existing bson.M
		err := collection.FindOne(ctx, bson.M{"username": candidate}, options.FindOne().SetProjection(bson.M{"_id": 1})).Decode(&existing)
		if err == mongo.ErrNoDocuments {
			return candidate
		}
		suffix, _ := randomState()
		if len(suffix) > 4 {
			suffix = suffix[:4]
		}
		candidate = base + "_" + suffix
		if len(candidate) > 28 {
			candidate = candidate[:28]
		}
	}
	return base + fmt.Sprintf("_%d", time.Now().Unix()%10000)
}

func setAuthCookies(w http.ResponseWriter, accessToken, refreshToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})
}
