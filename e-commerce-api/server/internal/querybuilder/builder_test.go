package querybuilder

import (
	"testing"
)

// =============================================================================
// SQL BUILDER TESTS
// =============================================================================

func TestNewSQLBuilder(t *testing.T) {
	builder := NewSQLBuilder()

	if builder == nil {
		t.Fatal("NewSQLBuilder returned nil")
	}

	if len(builder.allowedColumns) != 0 {
		t.Errorf("Expected 0 allowed columns, got %d", len(builder.allowedColumns))
	}

	if builder.defaultColumn != "id" {
		t.Errorf("Expected default column 'id', got %q", builder.defaultColumn)
	}

	if builder.defaultOrder != "DESC" {
		t.Errorf("Expected default order 'DESC', got %q", builder.defaultOrder)
	}
}

func TestSQLBuilder_AllowColumn(t *testing.T) {
	builder := NewSQLBuilder().
		AllowColumn("name").
		AllowColumn("price")

	if !builder.IsColumnAllowed("name") {
		t.Error("Expected 'name' to be allowed")
	}

	if !builder.IsColumnAllowed("price") {
		t.Error("Expected 'price' to be allowed")
	}

	if builder.IsColumnAllowed("invalid") {
		t.Error("Expected 'invalid' to NOT be allowed")
	}
}

func TestSQLBuilder_AllowColumns(t *testing.T) {
	builder := NewSQLBuilder().
		AllowColumns("name", "price", "stock", "created_at")

	allowedCols := builder.GetAllowedColumns()

	if len(allowedCols) != 4 {
		t.Errorf("Expected 4 allowed columns, got %d", len(allowedCols))
	}
}

func TestSQLBuilder_SetDefault(t *testing.T) {
	builder := NewSQLBuilder().
		SetDefault("created_at", "ASC")

	if builder.defaultColumn != "created_at" {
		t.Errorf("Expected default column 'created_at', got %q", builder.defaultColumn)
	}

	if builder.defaultOrder != "ASC" {
		t.Errorf("Expected default order 'ASC', got %q", builder.defaultOrder)
	}
}

func TestSQLBuilder_BuildOrderClause(t *testing.T) {
	tests := []struct {
		name           string
		allowedColumns []string
		defaultColumn  string
		defaultOrder   string
		inputColumn    string
		inputOrder     string
		expectedClause string
		description    string
	}{
		{
			name:           "Valid column and order",
			allowedColumns: []string{"name", "price"},
			defaultColumn:  "id",
			defaultOrder:   "DESC",
			inputColumn:    "price",
			inputOrder:     "ASC",
			expectedClause: "ORDER BY price ASC",
			description:    "Should accept valid whitelisted column",
		},
		{
			name:           "SQL injection in column name",
			allowedColumns: []string{"name", "price"},
			defaultColumn:  "id",
			defaultOrder:   "DESC",
			inputColumn:    "name; DROP TABLE products--",
			inputOrder:     "ASC",
			expectedClause: "ORDER BY id DESC",
			description:    "Should fallback to default when column contains SQL injection",
		},
		{
			name:           "SQL injection in order",
			allowedColumns: []string{"name", "price"},
			defaultColumn:  "id",
			defaultOrder:   "DESC",
			inputColumn:    "price",
			inputOrder:     "ASC; DELETE FROM users--",
			expectedClause: "ORDER BY price DESC",
			description:    "Should fallback to default order when order contains SQL injection",
		},
		{
			name:           "Case insensitive order - lowercase",
			allowedColumns: []string{"name"},
			defaultColumn:  "id",
			defaultOrder:   "DESC",
			inputColumn:    "name",
			inputOrder:     "asc",
			expectedClause: "ORDER BY name ASC",
			description:    "Should accept lowercase 'asc'",
		},
		{
			name:           "Case insensitive order - mixed case",
			allowedColumns: []string{"name"},
			defaultColumn:  "id",
			defaultOrder:   "DESC",
			inputColumn:    "name",
			inputOrder:     "AsC",
			expectedClause: "ORDER BY name ASC",
			description:    "Should accept mixed case 'AsC'",
		},
		{
			name:           "Invalid column with valid order",
			allowedColumns: []string{"name", "price"},
			defaultColumn:  "id",
			defaultOrder:   "DESC",
			inputColumn:    "invalid_column",
			inputOrder:     "ASC",
			expectedClause: "ORDER BY id DESC",
			description:    "Should use all defaults when column is invalid",
		},
		{
			name:           "Valid column with invalid order",
			allowedColumns: []string{"name", "price"},
			defaultColumn:  "id",
			defaultOrder:   "DESC",
			inputColumn:    "price",
			inputOrder:     "INVALID",
			expectedClause: "ORDER BY price DESC",
			description:    "Should use default order when order is invalid",
		},
		{
			name:           "Empty column name",
			allowedColumns: []string{"name", "price"},
			defaultColumn:  "id",
			defaultOrder:   "DESC",
			inputColumn:    "",
			inputOrder:     "ASC",
			expectedClause: "ORDER BY id DESC",
			description:    "Should use defaults for empty column",
		},
		{
			name:           "Empty order",
			allowedColumns: []string{"name", "price"},
			defaultColumn:  "id",
			defaultOrder:   "DESC",
			inputColumn:    "price",
			inputOrder:     "",
			expectedClause: "ORDER BY price DESC",
			description:    "Should use default order for empty order",
		},
		{
			name:           "UNION injection attempt",
			allowedColumns: []string{"name", "price"},
			defaultColumn:  "id",
			defaultOrder:   "DESC",
			inputColumn:    "name UNION SELECT * FROM users",
			inputOrder:     "DESC",
			expectedClause: "ORDER BY id DESC",
			description:    "Should prevent UNION injection",
		},
		{
			name:           "Comment injection attempt",
			allowedColumns: []string{"name", "price"},
			defaultColumn:  "id",
			defaultOrder:   "DESC",
			inputColumn:    "name--",
			inputOrder:     "DESC",
			expectedClause: "ORDER BY id DESC",
			description:    "Should prevent comment injection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewSQLBuilder().
				AllowColumns(tt.allowedColumns...).
				SetDefault(tt.defaultColumn, tt.defaultOrder)

			result := builder.BuildOrderClause(tt.inputColumn, tt.inputOrder)

			if result != tt.expectedClause {
				t.Errorf("%s\nExpected: %q\nGot:      %q",
					tt.description, tt.expectedClause, result)
			}
		})
	}
}

func TestSQLBuilder_FluentAPI(t *testing.T) {
	// Test method chaining
	builder := NewSQLBuilder().
		AllowColumn("name").
		AllowColumn("price").
		AllowColumns("stock", "created_at").
		SetDefault("id", "DESC")

	if len(builder.allowedColumns) != 4 {
		t.Errorf("Expected 4 allowed columns, got %d", len(builder.allowedColumns))
	}

	clause := builder.BuildOrderClause("price", "ASC")
	expected := "ORDER BY price ASC"
	if clause != expected {
		t.Errorf("Expected %q, got %q", expected, clause)
	}
}

// =============================================================================
// WHERE CLAUSE BUILDER TESTS
// =============================================================================

func TestNewWhereClause(t *testing.T) {
	where := NewWhereClause()

	if where == nil {
		t.Fatal("NewWhereClause returned nil")
	}

	if len(where.conditions) != 0 {
		t.Errorf("Expected 0 conditions, got %d", len(where.conditions))
	}

	if where.argCount != 1 {
		t.Errorf("Expected argCount 1, got %d", where.argCount)
	}
}

func TestWhereClause_AddCondition(t *testing.T) {
	where := NewWhereClause()

	where.AddCondition("name ILIKE $1", "%test%")
	where.AddCondition("price > $2", 100.0)

	if len(where.conditions) != 2 {
		t.Errorf("Expected 2 conditions, got %d", len(where.conditions))
	}

	if len(where.args) != 2 {
		t.Errorf("Expected 2 args, got %d", len(where.args))
	}
}

func TestWhereClause_Build(t *testing.T) {
	tests := []struct {
		name       string
		conditions []struct {
			condition string
			args      []interface{}
		}
		expectedClause    string
		expectedArgsCount int
	}{
		{
			name: "Empty where clause",
			conditions: []struct {
				condition string
				args      []interface{}
			}{},
			expectedClause:    "",
			expectedArgsCount: 0,
		},
		{
			name: "Single condition",
			conditions: []struct {
				condition string
				args      []interface{}
			}{
				{"name = $1", []interface{}{"test"}},
			},
			expectedClause:    "WHERE name = $1",
			expectedArgsCount: 1,
		},
		{
			name: "Multiple conditions",
			conditions: []struct {
				condition string
				args      []interface{}
			}{
				{"name ILIKE $1", []interface{}{"%test%"}},
				{"price > $2", []interface{}{100.0}},
				{"stock >= $3", []interface{}{10}},
			},
			expectedClause:    "WHERE name ILIKE $1 AND price > $2 AND stock >= $3",
			expectedArgsCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			where := NewWhereClause()

			for _, cond := range tt.conditions {
				where.AddCondition(cond.condition, cond.args...)
			}

			clause, args := where.Build()

			if clause != tt.expectedClause {
				t.Errorf("Expected clause %q, got %q", tt.expectedClause, clause)
			}

			if len(args) != tt.expectedArgsCount {
				t.Errorf("Expected %d args, got %d", tt.expectedArgsCount, len(args))
			}
		})
	}
}

// =============================================================================
// PAGINATION BUILDER TESTS
// =============================================================================

func TestNewPaginationBuilder(t *testing.T) {
	pb := NewPaginationBuilder(100)

	if pb.maxLimit != 100 {
		t.Errorf("Expected maxLimit 100, got %d", pb.maxLimit)
	}
}

func TestPaginationBuilder_Validate(t *testing.T) {
	tests := []struct {
		name           string
		maxLimit       int
		inputPage      int
		inputLimit     int
		expectedLimit  int
		expectedOffset int
	}{
		{
			name:           "Valid pagination",
			maxLimit:       100,
			inputPage:      2,
			inputLimit:     10,
			expectedLimit:  10,
			expectedOffset: 10,
		},
		{
			name:           "Page less than 1",
			maxLimit:       100,
			inputPage:      0,
			inputLimit:     10,
			expectedLimit:  10,
			expectedOffset: 0,
		},
		{
			name:           "Limit exceeds max",
			maxLimit:       100,
			inputPage:      1,
			inputLimit:     200,
			expectedLimit:  100,
			expectedOffset: 0,
		},
		{
			name:           "Limit less than 1",
			maxLimit:       100,
			inputPage:      1,
			inputLimit:     0,
			expectedLimit:  10,
			expectedOffset: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pb := NewPaginationBuilder(tt.maxLimit)
			limit, offset := pb.Validate(tt.inputPage, tt.inputLimit)

			if limit != tt.expectedLimit {
				t.Errorf("Expected limit %d, got %d", tt.expectedLimit, limit)
			}

			if offset != tt.expectedOffset {
				t.Errorf("Expected offset %d, got %d", tt.expectedOffset, offset)
			}
		})
	}
}

// =============================================================================
// BENCHMARK TESTS
// =============================================================================

func BenchmarkBuildOrderClause(b *testing.B) {
	builder := NewSQLBuilder().
		AllowColumns("name", "price", "stock", "created_at").
		SetDefault("id", "DESC")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder.BuildOrderClause("price", "ASC")
	}
}

func BenchmarkWhereClauseBuild(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		where := NewWhereClause()
		where.AddCondition("name ILIKE $1", "%test%")
		where.AddCondition("price > $2", 100.0)
		where.Build()
	}
}
