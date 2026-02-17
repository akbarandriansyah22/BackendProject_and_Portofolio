package security

import (
	"context"
	"fmt"
)

// RBAC defines role-based access control
type RBAC struct {
	roles       map[string][]string // role -> permissions
	permissions map[string][]string // resource -> roles
}

// NewRBAC creates new RBAC instance
func NewRBAC() *RBAC {
	return &RBAC{
		roles:       make(map[string][]string),
		permissions: make(map[string][]string),
	}
}

// AddRole adds a new role
func (r *RBAC) AddRole(role string, permissions []string) {
	r.roles[role] = permissions
}

// HasPermission checks if role has permission
func (r *RBAC) HasPermission(role, permission string) bool {
	permissions, exists := r.roles[role]
	if !exists {
		return false
	}

	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// GetRolePermissions returns all permissions for a role
func (r *RBAC) GetRolePermissions(role string) []string {
	if permissions, exists := r.roles[role]; exists {
		return permissions
	}
	return []string{}
}

// ValidateAccess checks if user has required permission
func (r *RBAC) ValidateAccess(ctx context.Context, userRole, requiredPermission string) error {
	if !r.HasPermission(userRole, requiredPermission) {
		return fmt.Errorf("insufficient permissions: role %s does not have %s", userRole, requiredPermission)
	}
	return nil
}

// InitializeDefaultRoles sets up standard RBAC roles
func (r *RBAC) InitializeDefaultRoles() {
	// Admin permissions
	adminPerms := []string{
		"read:users", "write:users", "delete:users",
		"read:products", "write:products", "delete:products",
		"read:categories", "write:categories", "delete:categories",
		"read:orders", "write:orders", "delete:orders",
		"read:payments", "write:payments", "delete:payments",
		"read:roles", "write:roles", "delete:roles",
		"read:reports", "write:reports",
	}

	// Customer permissions
	customerPerms := []string{
		"read:profile", "write:profile",
		"read:products",
		"read:categories",
		"read:cart", "write:cart",
		"read:orders", "write:orders",
		"read:payments", "write:payments",
	}

	// Seller permissions
	sellerPerms := []string{
		"read:profile", "write:profile",
		"read:products", "write:products",
		"read:categories",
		"read:orders",
		"read:payments",
		"read:reports",
	}

	r.AddRole("admin", adminPerms)
	r.AddRole("customer", customerPerms)
	r.AddRole("seller", sellerPerms)
}
