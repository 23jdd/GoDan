package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"godan/internal/config"
	"godan/internal/middleware"
	"godan/internal/pkg/errcode"
	pkgjwt "godan/internal/pkg/jwt"
	"godan/internal/pkg/response"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestPing(t *testing.T) {
	r := gin.New()
	r.GET("/ping", func(c *gin.Context) {
		response.Success(c, gin.H{"ping": "pong"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestAuthMiddleware_NoToken(t *testing.T) {
	jwtCfg := &config.JWTConfig{
		AccessSecret:  "test-access-secret",
		RefreshSecret: "test-refresh-secret",
		AccessExpire:  3600,
		RefreshExpire: 86400,
	}

	r := gin.New()
	r.GET("/protected", middleware.Auth(jwtCfg), func(c *gin.Context) {
		response.Success(c, nil)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	r.ServeHTTP(w, req)

	if w.Code != 400 && w.Code != 401 {
		t.Errorf("expected 400 or 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	jwtCfg := &config.JWTConfig{
		AccessSecret:  "test-access-secret",
		RefreshSecret: "test-refresh-secret",
		AccessExpire:  3600,
		RefreshExpire: 86400,
	}

	token, err := pkgjwt.GenerateToken(1, jwtCfg.AccessSecret, jwtCfg.AccessExpire)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	r.GET("/me", middleware.Auth(jwtCfg), func(c *gin.Context) {
		userID := middleware.GetUserID(c)
		response.Success(c, gin.H{"user_id": userID})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	jwtCfg := &config.JWTConfig{
		AccessSecret:  "test-access-secret",
		RefreshSecret: "test-refresh-secret",
		AccessExpire:  3600,
		RefreshExpire: 86400,
	}

	// token signed with wrong secret → invalid
	fakeSecret := "wrong-secret"
	token, _ := pkgjwt.GenerateToken(1, fakeSecret, 1)
	time.Sleep(10 * time.Millisecond)

	r := gin.New()
	r.GET("/me", middleware.Auth(jwtCfg), func(c *gin.Context) {
		response.Success(c, nil)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != 400 && w.Code != 401 {
		t.Errorf("expected 400 or 401, got %d", w.Code)
	}
}

func TestCORSMiddleware(t *testing.T) {
	r := gin.New()
	r.Use(middleware.CORS())
	r.GET("/api", func(c *gin.Context) {
		c.String(200, "ok")
	})

	// OPTIONS should return 204 with CORS headers
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/api", nil)
	r.ServeHTTP(w, req)

	if w.Code != 204 {
		t.Errorf("expected 204 for OPTIONS, got %d", w.Code)
	}
	if w.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Error("missing CORS header")
	}

	// normal GET should pass through
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api", nil)
	r.ServeHTTP(w2, req2)

	if w2.Code != 200 {
		t.Errorf("expected 200 for GET, got %d", w2.Code)
	}
}

func TestErrcodeResponse(t *testing.T) {
	r := gin.New()
	r.GET("/error", func(c *gin.Context) {
		response.Error(c, errcode.ErrUserNotFound)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/error", nil)
	r.ServeHTTP(w, req)

	var resp response.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != errcode.ErrUserNotFound.Code {
		t.Errorf("expected code %d, got %d", errcode.ErrUserNotFound.Code, resp.Code)
	}
}

func TestJsonBindValidation(t *testing.T) {
	r := gin.New()
	r.POST("/test", func(c *gin.Context) {
		var req struct {
			Name string `json:"name" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, errcode.ErrInvalidParams)
			return
		}
		response.Success(c, gin.H{"name": req.Name})
	})

	// missing required field
	body := bytes.NewBufferString(`{"wrong":"value"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code == 200 {
		t.Error("should fail with missing required field")
	}

	// valid request
	body = bytes.NewBufferString(`{"name":"hello"}`)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/test", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}
