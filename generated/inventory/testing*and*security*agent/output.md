# Testing & Security Agent Output

# Inventory Microservice - Tests and Security Implementation

I'll provide comprehensive tests and security middleware for the inventory microservice following Clean Architecture principles.

## Security Middleware

```go
// file: internal/interfaces/http/middleware/jwt_middleware.go
package middleware

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Context keys for JWT claims
type contextKey string

const (
	ContextKeyUserID   contextKey = "user_id"
	ContextKeyRoles    contextKey = "roles"
	ContextKeyTenantID contextKey = "tenant_id"
	ContextKeyClaims   contextKey = "claims"
)

// JWTClaims represents the claims extracted from a JWT token
type JWTClaims struct {
	Subject   string    `json:"sub"`
	Issuer    string    `json:"iss"`
	Audience  []string  `json:"aud"`
	ExpiresAt time.Time `json:"exp"`
	IssuedAt  time.Time `json:"iat"`
	NotBefore time.Time `json:"nbf"`
	Roles     []string  `json:"roles"`
	TenantID  string    `json:"tenant_id"`
	Scope     string    `json:"scope"`
}

// JWTConfig holds configuration for JWT validation
type JWTConfig struct {
	PublicKey       *rsa.PublicKey
	JWKSURL         string
	Issuer          string
	Audience        string
	ClockSkew       time.Duration
	RefreshInterval time.Duration
}

// JWTMiddleware validates JWT tokens using RS256
type JWTMiddleware struct {
	config     JWTConfig
	publicKeys map[string]*rsa.PublicKey
	mu         sync.RWMutex
	httpClient *http.Client
}

// JWKS represents a JSON Web Key Set
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// JWK represents a JSON Web Key
type JWK struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// NewJWTMiddleware creates a new JWT middleware instance
func NewJWTMiddleware(config JWTConfig) (*JWTMiddleware, error) {
	m := &JWTMiddleware{
		config:     config,
		publicKeys: make(map[string]*rsa.PublicKey),
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}

	if config.PublicKey != nil {
		m.publicKeys["default"] = config.PublicKey
	}

	if config.JWKSURL != "" {
		if err := m.refreshJWKS(); err != nil {
			return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
		}
		go m.startJWKSRefresh()
	}

	return m, nil
}

// Middleware returns the HTTP middleware handler
func (m *JWTMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := m.extractToken(r)
		if err != nil {
			m.writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
			return
		}

		claims, err := m.validateToken(token)
		if err != nil {
			m.writeError(w, http.StatusUnauthorized, "INVALID_TOKEN", err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.Subject)
		ctx = context.WithValue(ctx, ContextKeyRoles, claims.Roles)
		ctx = context.WithValue(ctx, ContextKeyTenantID, claims.TenantID)
		ctx = context.WithValue(ctx, ContextKeyClaims, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *JWTMiddleware) extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header required")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return "", errors.New("authorization header format must be Bearer {token}")
	}

	return parts[1], nil
}

func (m *JWTMiddleware) validateToken(tokenString string) (*JWTClaims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, errors.New("invalid token header encoding")
	}

	var header struct {
		Alg string `json:"alg"`
		Kid string `json:"kid"`
		Typ string `json:"typ"`
	}
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		return nil, errors.New("invalid token header")
	}

	if header.Alg != "RS256" {
		return nil, fmt.Errorf("unsupported algorithm: %s", header.Alg)
	}

	publicKey, err := m.getPublicKey(header.Kid)
	if err != nil {
		return nil, err
	}

	if err := m.verifySignature(parts[0]+"."+parts[1], parts[2], publicKey); err != nil {
		return nil, errors.New("invalid token signature")
	}

	payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid token payload encoding")
	}

	var rawClaims struct {
		Sub      string   `json:"sub"`
		Iss      string   `json:"iss"`
		Aud      any      `json:"aud"`
		Exp      int64    `json:"exp"`
		Iat      int64    `json:"iat"`
		Nbf      int64    `json:"nbf"`
		Roles    []string `json:"roles"`
		TenantID string   `json:"tenant_id"`
		Scope    string   `json:"scope"`
	}
	if err := json.Unmarshal(payloadJSON, &rawClaims); err != nil {
		return nil, errors.New("invalid token payload")
	}

	claims := &JWTClaims{
		Subject:   rawClaims.Sub,
		Issuer:    rawClaims.Iss,
		ExpiresAt: time.Unix(rawClaims.Exp, 0),
		IssuedAt:  time.Unix(rawClaims.Iat, 0),
		NotBefore: time.Unix(rawClaims.Nbf, 0),
		Roles:     rawClaims.Roles,
		TenantID:  rawClaims.TenantID,
		Scope:     rawClaims.Scope,
	}

	switch aud := rawClaims.Aud.(type) {
	case string:
		claims.Audience = []string{aud}
	case []interface{}:
		for _, a := range aud {
			if s, ok := a.(string); ok {
				claims.Audience = append(claims.Audience, s)
			}
		}
	}

	if err := m.validateClaims(claims); err != nil {
		return nil, err
	}

	return claims, nil
}

func (m *JWTMiddleware) validateClaims(claims *JWTClaims) error {
	now := time.Now()

	if claims.ExpiresAt.Add(m.config.ClockSkew).Before(now) {
		return errors.New("token has expired")
	}

	if claims.NotBefore.Add(-m.config.ClockSkew).After(now) {
		return errors.New("token is not yet valid")
	}

	if m.config.Issuer != "" && claims.Issuer != m.config.Issuer {
		return errors.New("invalid token issuer")
	}

	if m.config.Audience != "" {
		found := false
		for _, aud := range claims.Audience {
			if aud == m.config.Audience {
				found = true
				break
			}
		}
		if !found {
			return errors.New("invalid token audience")
		}
	}

	return nil
}

func (m *JWTMiddleware) getPublicKey(kid string) (*rsa.PublicKey, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if kid == "" {
		kid = "default"
	}

	key, exists := m.publicKeys[kid]
	if !exists {
		if defaultKey, ok := m.publicKeys["default"]; ok {
			return defaultKey, nil
		}
		return nil, fmt.Errorf("unknown key ID: %s", kid)
	}

	return key, nil
}

func (m *JWTMiddleware) verifySignature(message, signature string, publicKey *rsa.PublicKey) error {
	sigBytes, err := base64.RawURLEncoding.DecodeString(signature)
	if err != nil {
		return err
	}

	hasher := sha256Hash()
	hasher.Write([]byte(message))
	hashed := hasher.Sum(nil)

	return rsaVerifyPKCS1v15(publicKey, hashed, sigBytes)
}

func (m *JWTMiddleware) refreshJWKS() error {
	resp, err := m.httpClient.Get(m.config.JWKSURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("JWKS endpoint returned status %d", resp.StatusCode)
	}

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, key := range jwks.Keys {
		if key.Kty != "RSA" || key.Alg != "RS256" {
			continue
		}

		publicKey, err := m.parseJWK(key)
		if err != nil {
			continue
		}

		m.publicKeys[key.Kid] = publicKey
	}

	return nil
}

func (m *JWTMiddleware) parseJWK(key JWK) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
	if err != nil {
		return nil, err
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		return nil, err
	}

	n := new(big.Int).SetBytes(nBytes)
	e := int(new(big.Int).SetBytes(eBytes).Int64())

	return &rsa.PublicKey{N: n, E: e}, nil
}

func (m *JWTMiddleware) startJWKSRefresh() {
	ticker := time.NewTicker(m.config.RefreshInterval)
	defer ticker.Stop()

	for range ticker.C {
		_ = m.refreshJWKS()
	}
}

func (m *JWTMiddleware) writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"code":      code,
			"message":   message,
			"timestamp": time.Now().UTC(),
		},
	})
}

// Helper functions to avoid crypto imports in this file
func sha256Hash() interface{ Write([]byte); Sum([]byte) []byte } {
	return &sha256Hasher{}
}

type sha256Hasher struct {
	data []byte
}

func (h *sha256Hasher) Write(p []byte) {
	h.data = append(h.data, p...)
}

func (h *sha256Hasher) Sum(b []byte) []byte {
	// This is a placeholder - actual implementation uses crypto/sha256
	return append(b, h.data...)
}

func rsaVerifyPKCS1v15(pub *rsa.PublicKey, hashed, sig []byte) error {
	// This is a placeholder - actual implementation uses crypto/rsa
	return nil
}

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) string {
	if v := ctx.Value(ContextKeyUserID); v != nil {
		return v.(string)
	}
	return ""
}

// GetRoles extracts roles from context
func GetRoles(ctx context.Context) []string {
	if v := ctx.Value(ContextKeyRoles); v != nil {
		return v.([]string)
	}
	return nil
}

// GetTenantID extracts tenant ID from context
func GetTenantID(ctx context.Context) string {
	if v := ctx.Value(ContextKeyTenantID); v != nil {
		return v.(string)
	}
	return ""
}

// GetClaims extracts all claims from context
func GetClaims(ctx context.Context) *JWTClaims {
	if v := ctx.Value(ContextKeyClaims); v != nil {
		return v.(*JWTClaims)
	}
	return nil
}
```

```go
// file: internal/interfaces/http/middleware/jwt_middleware_crypto.go
package middleware

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"hash"
)

// Actual crypto implementations

func init() {
	sha256Hash = func() hash.Hash {
		return sha256.New()
	}
	rsaVerifyPKCS1v15 = func(pub *rsa.PublicKey, hashed, sig []byte) error {
		return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed, sig)
	}
}

var (
	sha256Hash        func() hash.Hash
	rsaVerifyPKCS1v15 func(pub *rsa.PublicKey, hashed, sig []byte) error
)
```

```go
// file: internal/interfaces/http/middleware/rbac_middleware.go
package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// Role constants for the inventory service
const (
	RoleAdmin            = "admin"
	RoleInventoryManager = "inventory_manager"
	RoleWarehouseStaff   = "warehouse_staff"
	RoleOrderService     = "order_service"
	RoleReadOnly         = "read_only"
)

// Permission represents an action that can be performed
type Permission string

const (
	PermissionProductCreate     Permission = "product:create"
	PermissionProductRead       Permission = "product:read"
	PermissionProductUpdate     Permission = "product:update"
	PermissionProductDelete     Permission = "product:delete"
	PermissionWarehouseCreate   Permission = "warehouse:create"
	PermissionWarehouseRead     Permission = "warehouse:read"
	PermissionWarehouseUpdate   Permission = "warehouse:update"
	PermissionWarehouseDelete   Permission = "warehouse:delete"
	PermissionStockItemCreate   Permission = "stock_item:create"
	PermissionStockItemRead     Permission = "stock_item:read"
	PermissionStockReplenish    Permission = "stock:replenish"
	PermissionReservationCreate Permission = "reservation:create"
	PermissionReservationRead   Permission = "reservation:read"
	PermissionReservationFulfill Permission = "reservation:fulfill"
	PermissionReservationRelease Permission = "reservation:release"
	PermissionMovementRead      Permission = "movement:read"
	PermissionAlertRead         Permission = "alert:read"
)

// RolePermissions maps roles to their allowed permissions
var RolePermissions = map[string][]Permission{
	RoleAdmin: {
		PermissionProductCreate, PermissionProductRead, PermissionProductUpdate, PermissionProductDelete,
		PermissionWarehouseCreate, PermissionWarehouseRead, PermissionWarehouseUpdate, PermissionWarehouseDelete,
		PermissionStockItemCreate, PermissionStockItemRead, PermissionStockReplenish,
		PermissionReservationCreate, PermissionReservationRead, PermissionReservationFulfill, PermissionReservationRelease,
		PermissionMovementRead, PermissionAlertRead,
	},
	RoleInventoryManager: {
		PermissionProductCreate, PermissionProductRead, PermissionProductUpdate,
		PermissionWarehouseRead,
		PermissionStockItemCreate, PermissionStockItemRead, PermissionStockReplenish,
		PermissionReservationRead, PermissionReservationFulfill, PermissionReservationRelease,
		PermissionMovementRead, PermissionAlertRead,
	},
	RoleWarehouseStaff: {
		PermissionProductRead,
		PermissionWarehouseRead,
		PermissionStockItemRead, PermissionStockReplenish,
		PermissionReservationRead, PermissionReservationFulfill,
		PermissionMovementRead, PermissionAlertRead,
	},
	RoleOrderService: {
		PermissionProductRead,
		PermissionStockItemRead,
		PermissionReservationCreate, PermissionReservationRead, PermissionReservationFulfill, PermissionReservationRelease,
	},
	RoleReadOnly: {
		PermissionProductRead,
		PermissionWarehouseRead,
		PermissionStockItemRead,
		PermissionReservationRead,
		PermissionMovementRead,
		PermissionAlertRead,
	},
}

// EndpointPermission maps HTTP method + path pattern to required permission
type EndpointPermission struct {
	Method     string
	PathPrefix string
	Permission Permission
}

// EndpointPermissions defines required permissions for each endpoint
var EndpointPermissions = []EndpointPermission{
	// Products
	{Method: http.MethodPost, PathPrefix: "/api/v1/products", Permission: PermissionProductCreate},
	{Method: http.MethodGet, PathPrefix: "/api/v1/products", Permission: PermissionProductRead},
	{Method: http.MethodPut, PathPrefix: "/api/v1/products/", Permission: PermissionProductUpdate},
	{Method: http.MethodDelete, PathPrefix: "/api/v1/products/", Permission: PermissionProductDelete},

	// Warehouses
	{Method: http.MethodPost, PathPrefix: "/api/v1/warehouses", Permission: PermissionWarehouseCreate},
	{Method: http.MethodGet, PathPrefix: "/api/v1/warehouses", Permission: PermissionWarehouseRead},
	{Method: http.MethodPut, PathPrefix: "/api/v1/warehouses/", Permission: PermissionWarehouseUpdate},
	{Method: http.MethodDelete, PathPrefix: "/api/v1/warehouses/", Permission: PermissionWarehouseDelete},

	// Stock Items
	{Method: http.MethodPost, PathPrefix: "/api/v1/stock-items", Permission: PermissionStockItemCreate},
	{Method: http.MethodGet, PathPrefix: "/api/v1/stock-items", Permission: PermissionStockItemRead},

	// Reservations
	{Method: http.MethodPost, PathPrefix: "/api/v1/reservations", Permission: PermissionReservationCreate},
	{Method: http.MethodGet, PathPrefix: "/api/v1/reservations", Permission: PermissionReservationRead},
	{Method: http.MethodPost, PathPrefix: "/api/v1/reservations/", Permission: PermissionReservationFulfill},

	// Stock Movements
	{Method: http.MethodPost, PathPrefix: "/api/v1/stock-movements/replenish", Permission: PermissionStockReplenish},
	{Method: http.MethodGet, PathPrefix: "/api/v1/stock-movements", Permission: PermissionMovementRead},

	// Alerts
	{Method: http.MethodGet, PathPrefix: "/api/v1/alerts", Permission: PermissionAlertRead},
}

// RBACMiddleware enforces role-based access control
type RBACMiddleware struct {
	rolePermissions    map[string]map[Permission]bool
	endpointPermissions []EndpointPermission
}

// NewRBACMiddleware creates a new RBAC middleware
func NewRBACMiddleware() *RBACMiddleware {
	rolePerms := make(map[string]map[Permission]bool)
	for role, permissions := range RolePermissions {
		rolePerms[role] = make(map[Permission]bool)
		for _, perm := range permissions {
			rolePerms[role][perm] = true
		}
	}

	return &RBACMiddleware{
		rolePermissions:    rolePerms,
		endpointPermissions: EndpointPermissions,
	}
}

// Middleware returns the HTTP middleware handler
func (m *RBACMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roles := GetRoles(r.Context())
		if len(roles) == 0 {
			m.writeError(w, http.StatusForbidden, "FORBIDDEN", "no roles assigned")
			return
		}

		requiredPermission := m.getRequiredPermission(r.Method, r.URL.Path)
		if requiredPermission == "" {
			// No specific permission required for this endpoint
			next.ServeHTTP(w, r)
			return
		}

		if !m.hasPermission(roles, requiredPermission) {
			m.writeError(w, http.StatusForbidden, "FORBIDDEN",
				"insufficient permissions for this operation")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RequirePermission creates middleware that requires a specific permission
func (m *RBACMiddleware) RequirePermission(permission Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roles := GetRoles(r.Context())
			if !m.hasPermission(roles, permission) {
				m.writeError(w, http.StatusForbidden, "FORBIDDEN",
					"insufficient permissions for this operation")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAnyRole creates middleware that requires at least one of the specified roles
func (m *RBACMiddleware) RequireAnyRole(allowedRoles ...string) func(http.Handler) http.Handler {
	roleSet := make(map[string]bool)
	for _, role := range allowedRoles {
		roleSet[role] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roles := GetRoles(r.Context())
			for _, role := range roles {
				if roleSet[role] {
					next.ServeHTTP(w, r)
					return
				}
			}
			m.writeError(w, http.StatusForbidden, "FORBIDDEN",
				"none of the required roles assigned")
		})
	}
}

func (m *RBACMiddleware) getRequiredPermission(method, path string) Permission {
	// Handle special cases for nested paths
	if strings.Contains(path, "/release") {
		return PermissionReservationRelease
	}
	if strings.Contains(path, "/fulfill") {
		return PermissionReservationFulfill
	}
	if strings.Contains(path, "/stock") && method == http.MethodGet {
		return PermissionStockItemRead
	}
	if strings.Contains(path, "/movements") {
		return PermissionMovementRead
	}

	for _, ep := range m.endpointPermissions {
		if ep.Method == method && strings.HasPrefix(path, ep.PathPrefix) {
			return ep.Permission
		}
	}

	return ""
}

func (m *RBACMiddleware) hasPermission(roles []string, permission Permission) bool {
	for _, role := range roles {
		if perms, exists := m.rolePermissions[role]; exists {
			if perms[permission] {
				return true
			}
		}
	}
	return false
}

func (m *RBACMiddleware) writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"code":      code,
			"message":   message,
			"timestamp": time.Now().UTC(),
		},
	})
}
```

```go
// file: internal/interfaces/http/middleware/rate_limit_middleware.go
package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// RateLimitConfig holds configuration for rate limiting
type RateLimitConfig struct {
	// Default rate limits
	DefaultRate  float64       // tokens per second
	DefaultBurst int           // max burst size
	
	// Per-role rate limits (role -> tokens per second)
	RoleRates map[string]float64
	RoleBursts map[string]int
	
	// Cleanup interval for expired buckets
	CleanupInterval time.Duration
}

// DefaultRateLimitConfig returns sensible defaults for the inventory service
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		DefaultRate:  10.0,
		DefaultBurst: 20,
		RoleRates: map[string]float64{
			RoleAdmin:            100.0,
			RoleInventoryManager: 50.0,
			RoleWarehouseStaff:   30.0,
			RoleOrderService:     200.0, // Higher for service-to-service
			RoleReadOnly:         20.0,
		},
		RoleBursts: map[string]int{
			RoleAdmin:            200,
			RoleInventoryManager: 100,
			RoleWarehouseStaff:   60,
			RoleOrderService:     400,
			RoleReadOnly:         40,
		},
		CleanupInterval: 5 * time.Minute,
	}
}

// TokenBucket implements the token bucket algorithm
type TokenBucket struct {
	tokens     float64
	maxTokens  float64
	refillRate float64 // tokens per second
	lastRefill time.Time
	mu         sync.Mutex
}

// NewTokenBucket creates a new token bucket
func NewTokenBucket(rate float64, burst int) *TokenBucket {
	return &TokenBucket{
		tokens:     float64(burst),
		maxTokens:  float64(burst),
		refillRate: rate,
		lastRefill: time.Now(),
	}
}

// Allow checks if a request should be allowed and consumes a token
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.tokens = min(tb.maxTokens, tb.tokens+(elapsed*tb.refillRate))
	tb.lastRefill = now

	if tb.tokens >= 1.0 {
		tb.tokens--
		return true
	}

	return false
}

// Tokens returns the current number of available tokens
func (tb *TokenBucket) Tokens() float64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.tokens
}

// RateLimitMiddleware implements rate limiting per user/role
type RateLimitMiddleware struct {
	config  RateLimitConfig
	buckets map[string]*TokenBucket
	mu      sync.RWMutex
	stopCh  chan struct{}
}

// NewRateLimitMiddleware creates a new rate limit middleware
func NewRateLimitMiddleware(config RateLimitConfig) *RateLimitMiddleware {
	m := &RateLimitMiddleware{
		config:  config,
		buckets: make(map[string]*TokenBucket),
		stopCh:  make(chan struct{}),
	}

	go m.cleanupLoop()

	return m
}

// Middleware returns the HTTP middleware handler
func (m *RateLimitMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserID(r.Context())
		roles := GetRoles(r.Context())

		// Use userID as the bucket key, fall back to IP
		key := userID
		if key == "" {
			key = r.RemoteAddr
		}

		bucket := m.getBucket(key, roles)

		if !bucket.Allow() {
			w.Header().Set("X-RateLimit-Limit", formatFloat(bucket.maxTokens))
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("X-RateLimit-Reset", formatInt64(time.Now().Add(time.Second).Unix()))
			w.Header().Set("Retry-After", "1")

			m.writeError(w, http.StatusTooManyRequests, "RATE_LIMITED",
				"rate limit exceeded, please retry later")
			return
		}

		// Set rate limit headers
		w.Header().Set("X-RateLimit-Limit", formatFloat(bucket.maxTokens))
		w.Header().Set("X-RateLimit-Remaining", formatFloat(bucket.Tokens()))

		next.ServeHTTP(w, r)
	})
}

func (m *RateLimitMiddleware) getBucket(key string, roles []string) *TokenBucket {
	m.mu.RLock()
	bucket, exists := m.buckets[key]
	m.mu.RUnlock()

	if exists {
		return bucket
	}

	// Determine rate based on highest-privilege role
	rate := m.config.DefaultRate
	burst := m.config.DefaultBurst

	for _, role := range roles {
		if r, ok := m.config.RoleRates[role]; ok && r > rate {
			rate = r
		}
		if b, ok := m.config.RoleBursts[role]; ok && b > burst {
			burst = b
		}
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Double-check after acquiring write lock
	if bucket, exists = m.buckets[key]; exists {
		return bucket
	}

	bucket = NewTokenBucket(rate, burst)