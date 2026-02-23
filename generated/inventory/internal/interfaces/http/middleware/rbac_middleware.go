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