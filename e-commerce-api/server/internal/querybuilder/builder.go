package querybuilder

import (
	"fmt"
	"strings"
)

// =============================================================================
// SQL BUILDER - Prevents SQL Injection in Dynamic Queries
// =============================================================================
// SQLBuilder provides a safe way to build dynamic SQL ORDER BY clauses
// It uses whitelist-based validation to prevent SQL injection attacks
type SQLBuilder struct {
	allowedColumns map[string]bool
	allowedOrders  map[string]bool
	defaultColumn  string
	defaultOrder   string
}

// NewSQLBuilder creates a new SQL query builder with secure defaults
func NewSQLBuilder() *SQLBuilder {
	return &SQLBuilder{
		allowedColumns: make(map[string]bool),
		allowedOrders: map[string]bool{
			"ASC":  true,
			"DESC": true,
		},
		defaultColumn: "id",
		defaultOrder:  "DESC",
	}
}

// AllowColumn adds a single column to the whitelist
// Only whitelisted columns can be used in ORDER BY clauses
func (b *SQLBuilder) AllowColumn(column string) *SQLBuilder {
	b.allowedColumns[column] = true
	return b
}

// AllowColumns adds multiple columns to the whitelist at once
// This is a convenience method for fluent API usage
func (b *SQLBuilder) AllowColumns(columns ...string) *SQLBuilder {
	for _, col := range columns {
		b.allowedColumns[col] = true
	}
	return b
}

// SetDefault sets the default column and order for fallback scenarios
// These defaults are used when user input is invalid or malicious
func (b *SQLBuilder) SetDefault(column, order string) *SQLBuilder {
	b.defaultColumn = column
	b.defaultOrder = strings.ToUpper(order)
	return b
}

// BuildOrderClause safely constructs an ORDER BY clause
// It validates both column name and sort order against whitelists
// Returns: "ORDER BY column_name ORDER_DIRECTION"
func (b *SQLBuilder) BuildOrderClause(column, order string) string {
	// FULL fallback if column invalid
	if column == "" || !b.allowedColumns[column] {
		return fmt.Sprintf(
			"ORDER BY %s %s",
			b.defaultColumn,
			b.defaultOrder,
		)
	}

	// Column valid → validate order
	upperOrder := strings.ToUpper(order)
	if !b.allowedOrders[upperOrder] {
		upperOrder = b.defaultOrder
	}

	return fmt.Sprintf("ORDER BY %s %s", column, upperOrder)
}

// GetAllowedColumns returns a copy of allowed columns (for debugging)
func (b *SQLBuilder) GetAllowedColumns() []string {
	columns := make([]string, 0, len(b.allowedColumns))
	for col := range b.allowedColumns {
		columns = append(columns, col)
	}
	return columns
}

// IsColumnAllowed checks if a column is in the whitelist
func (b *SQLBuilder) IsColumnAllowed(column string) bool {
	return b.allowedColumns[column]
}

// =============================================================================
// WHERE CLAUSE BUILDER - For Safe Parameterized Queries
// =============================================================================

// WhereClause helps build safe WHERE clauses with parameterized queries
type WhereClause struct {
	conditions []string
	args       []interface{}
	argCount   int
}

// NewWhereClause creates a new WHERE clause builder
func NewWhereClause() *WhereClause {
	return &WhereClause{
		conditions: []string{},
		args:       []interface{}{},
		argCount:   1,
	}
}

// AddCondition adds a condition with parameterized arguments
// Example: AddCondition("name ILIKE $1", "%search%")
func (w *WhereClause) AddCondition(condition string, args ...interface{}) *WhereClause {
	w.conditions = append(w.conditions, condition)
	w.args = append(w.args, args...)
	w.argCount += len(args)
	return w
}

// Build constructs the final WHERE clause and returns it with arguments
// Returns: ("WHERE condition1 AND condition2", []args)
func (w *WhereClause) Build() (string, []interface{}) {
	if len(w.conditions) == 0 {
		return "", w.args
	}
	return "WHERE " + strings.Join(w.conditions, " AND "), w.args
}

// Args returns the accumulated arguments array
func (w *WhereClause) Args() []interface{} {
	return w.args
}

// ArgCount returns the current argument count for building queries
func (w *WhereClause) ArgCount() int {
	return w.argCount
}

// =============================================================================
// PAGINATION BUILDER - For Safe LIMIT/OFFSET
// =============================================================================

// PaginationBuilder helps build safe pagination clauses
type PaginationBuilder struct {
	maxLimit int
}

// NewPaginationBuilder creates a new pagination builder
func NewPaginationBuilder(maxLimit int) *PaginationBuilder {
	if maxLimit <= 0 {
		maxLimit = 100 // Default max limit
	}
	return &PaginationBuilder{maxLimit: maxLimit}
}

// Validate validates and sanitizes pagination parameters
// Returns: (safeLimit, safeOffset)
func (p *PaginationBuilder) Validate(page, limit int) (int, int) {
	// Validate page
	if page < 1 {
		page = 1
	}

	// Validate limit
	if limit < 1 {
		limit = 10 // Default
	}
	if limit > p.maxLimit {
		limit = p.maxLimit
	}

	// Calculate offset
	offset := (page - 1) * limit

	return limit, offset
}

// BuildClause builds the LIMIT/OFFSET clause
func (p *PaginationBuilder) BuildClause(argStart int) string {
	return fmt.Sprintf("LIMIT $%d OFFSET $%d", argStart, argStart+1)
}
