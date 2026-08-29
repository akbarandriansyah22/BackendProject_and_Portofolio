package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
)

// Thin adapters so existing repositories satisfy ports without rewriting SQL.

type UserRepositoryPort struct{ inner *UserRepository }

func NewUserRepositoryPort(inner *UserRepository) *UserRepositoryPort {
	return &UserRepositoryPort{inner: inner}
}

func (r *UserRepositoryPort) Create(ctx context.Context, user *models.User) error {
	_ = ctx
	return r.inner.Create(user)
}
func (r *UserRepositoryPort) GetByID(ctx context.Context, id int) (*models.User, error) {
	_ = ctx
	return r.inner.GetByID(id)
}
func (r *UserRepositoryPort) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	_ = ctx
	return r.inner.GetByEmail(email)
}
func (r *UserRepositoryPort) Update(ctx context.Context, user *models.User) error {
	_ = ctx
	return r.inner.Update(user)
}
func (r *UserRepositoryPort) Delete(ctx context.Context, id int) error {
	_ = ctx
	return r.inner.Delete(id)
}
func (r *UserRepositoryPort) List(ctx context.Context, page, limit int) ([]*models.User, int, error) {
	_ = ctx
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	users, total, err := r.inner.GetAll(limit, (page-1)*limit)
	if err != nil {
		return nil, 0, err
	}
	out := make([]*models.User, len(users))
	for i := range users {
		u := users[i]
		out[i] = &u
	}
	return out, int(total), nil
}
func (r *UserRepositoryPort) UpdatePassword(ctx context.Context, userID int, passwordHash string) error {
	_ = ctx
	return r.inner.UpdatePassword(userID, passwordHash)
}
func (r *UserRepositoryPort) VerifyEmail(ctx context.Context, userID int) error {
	_ = ctx
	return r.inner.VerifyEmail(userID)
}
func (r *UserRepositoryPort) SetActive(ctx context.Context, userID int, isActive bool) error {
	_ = ctx
	return r.inner.SetActive(userID, isActive)
}

type RoleRepositoryPort struct{ inner *RoleRepository }

func NewRoleRepositoryPort(inner *RoleRepository) *RoleRepositoryPort {
	return &RoleRepositoryPort{inner: inner}
}

func (r *RoleRepositoryPort) Create(ctx context.Context, role *models.Role) error {
	_ = ctx
	return r.inner.Create(role)
}
func (r *RoleRepositoryPort) GetByID(ctx context.Context, id int) (*models.Role, error) {
	_ = ctx
	return r.inner.GetByID(id)
}
func (r *RoleRepositoryPort) GetByName(ctx context.Context, name string) (*models.Role, error) {
	_ = ctx
	return r.inner.GetByName(name)
}
func (r *RoleRepositoryPort) Update(ctx context.Context, role *models.Role) error {
	_ = ctx
	return r.inner.Update(role)
}
func (r *RoleRepositoryPort) Delete(ctx context.Context, id int) error {
	_ = ctx
	return r.inner.Delete(id)
}
func (r *RoleRepositoryPort) List(ctx context.Context) ([]*models.Role, error) {
	_ = ctx
	roles, err := r.inner.GetAll()
	if err != nil {
		return nil, err
	}
	out := make([]*models.Role, len(roles))
	for i := range roles {
		role := roles[i]
		out[i] = &role
	}
	return out, nil
}

type CategoryRepositoryPort struct{ inner *CategoryRepository }

func NewCategoryRepositoryPort(inner *CategoryRepository) *CategoryRepositoryPort {
	return &CategoryRepositoryPort{inner: inner}
}

func (r *CategoryRepositoryPort) Create(ctx context.Context, category *models.Category) error {
	_ = ctx
	return r.inner.Create(category)
}
func (r *CategoryRepositoryPort) GetByID(ctx context.Context, id int) (*models.Category, error) {
	_ = ctx
	return r.inner.GetByID(id)
}
func (r *CategoryRepositoryPort) GetBySlug(ctx context.Context, slug string) (*models.Category, error) {
	_ = ctx
	return r.inner.GetBySlug(slug)
}
func (r *CategoryRepositoryPort) Update(ctx context.Context, category *models.Category) error {
	_ = ctx
	return r.inner.Update(category)
}
func (r *CategoryRepositoryPort) Delete(ctx context.Context, id int) error {
	_ = ctx
	return r.inner.Delete(id)
}
func (r *CategoryRepositoryPort) List(ctx context.Context) ([]*models.Category, error) {
	_ = ctx
	cats, err := r.inner.GetAll()
	if err != nil {
		return nil, err
	}
	out := make([]*models.Category, len(cats))
	for i := range cats {
		c := cats[i]
		out[i] = &c
	}
	return out, nil
}
func (r *CategoryRepositoryPort) GetChildren(ctx context.Context, parentID int) ([]*models.Category, error) {
	_ = ctx
	cats, err := r.inner.GetSubCategories(parentID)
	if err != nil {
		return nil, err
	}
	out := make([]*models.Category, len(cats))
	for i := range cats {
		c := cats[i]
		out[i] = &c
	}
	return out, nil
}
func (r *CategoryRepositoryPort) GetRootCategories(ctx context.Context) ([]*models.Category, error) {
	_ = ctx
	cats, err := r.inner.GetParentCategories()
	if err != nil {
		return nil, err
	}
	out := make([]*models.Category, len(cats))
	for i := range cats {
		c := cats[i]
		out[i] = &c
	}
	return out, nil
}
func (r *CategoryRepositoryPort) HasProducts(ctx context.Context, categoryID int) (bool, error) {
	_ = ctx
	return r.inner.HasProducts(categoryID)
}
func (r *CategoryRepositoryPort) SetActive(ctx context.Context, categoryID int, isActive bool) error {
	_ = ctx
	return r.inner.UpdateStatus(categoryID, isActive)
}
func (r *CategoryRepositoryPort) AssignProductToCategory(ctx context.Context, productID, categoryID int) error {
	_, err := r.inner.db.ExecContext(ctx,
		`INSERT INTO product_categories (product_id, category_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		productID, categoryID,
	)
	return err
}
func (r *CategoryRepositoryPort) RemoveProductFromCategory(ctx context.Context, productID, categoryID int) error {
	_, err := r.inner.db.ExecContext(ctx,
		`DELETE FROM product_categories WHERE product_id = $1 AND category_id = $2`,
		productID, categoryID,
	)
	return err
}
func (r *CategoryRepositoryPort) GetProductCategories(ctx context.Context, productID int) ([]*models.Category, error) {
	query := `
		SELECT c.id, c.name, c.slug, c.description, c.parent_id, c.is_active, c.created_at, c.updated_at
		FROM categories c
		INNER JOIN product_categories pc ON c.id = pc.category_id
		WHERE pc.product_id = $1
	`
	rows, err := r.inner.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var out []*models.Category
	for rows.Next() {
		c := &models.Category{}
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.ParentID, &c.IsActive, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, nil
}

type CartRepositoryPort struct{ inner *CartRepository }

func NewCartRepositoryPort(inner *CartRepository) *CartRepositoryPort {
	return &CartRepositoryPort{inner: inner}
}

func (r *CartRepositoryPort) Create(ctx context.Context, userID int) (*models.Cart, error) {
	_ = ctx
	return r.inner.GetOrCreateCart(userID)
}
func (r *CartRepositoryPort) GetByUserID(ctx context.Context, userID int) (*models.Cart, error) {
	_ = ctx
	return r.inner.GetByUserID(userID)
}
func (r *CartRepositoryPort) AddItem(ctx context.Context, cartID, productID, quantity int) error {
	var price float64
	err := r.inner.db.QueryRowContext(ctx, `SELECT price FROM products WHERE id = $1`, productID).Scan(&price)
	if err != nil {
		return err
	}
	return r.inner.AddItem(cartID, productID, quantity, price)
}
func (r *CartRepositoryPort) UpdateItemQuantity(ctx context.Context, cartID, cartItemID, quantity int) error {
	_ = ctx
	_ = cartID
	return r.inner.UpdateItemQuantity(cartItemID, quantity)
}
func (r *CartRepositoryPort) RemoveItem(ctx context.Context, cartID, cartItemID int) error {
	_ = ctx
	_ = cartID
	return r.inner.RemoveItem(cartItemID)
}
func (r *CartRepositoryPort) GetCartItems(ctx context.Context, cartID int) ([]*models.CartItem, error) {
	_ = ctx
	items, err := r.inner.GetCartItems(cartID)
	if err != nil {
		return nil, err
	}
	out := make([]*models.CartItem, len(items))
	for i, it := range items {
		out[i] = &models.CartItem{
			ID:        it.ID,
			CartID:    it.CartID,
			ProductID: it.ProductID,
			Quantity:  it.Quantity,
			Price:     it.Price,
			CreatedAt: it.CreatedAt,
			UpdatedAt: it.UpdatedAt,
		}
	}
	return out, nil
}
func (r *CartRepositoryPort) ClearCart(ctx context.Context, cartID int) error {
	_ = ctx
	return r.inner.ClearCart(cartID)
}
func (r *CartRepositoryPort) GetItemCount(ctx context.Context, cartID int) (int, error) {
	_ = ctx
	return r.inner.GetCartItemCount(cartID)
}
func (r *CartRepositoryPort) GetCartTotal(ctx context.Context, cartID int) (float64, error) {
	_ = ctx
	return r.inner.GetCartTotal(cartID)
}
func (r *CartRepositoryPort) CheckItemExists(ctx context.Context, cartID, productID int) (bool, int, error) {
	_ = ctx
	item, err := r.inner.GetCartItemByCartAndProduct(cartID, productID)
	if err != nil {
		return false, 0, nil
	}
	if item == nil {
		return false, 0, nil
	}
	return true, item.Quantity, nil
}
func (r *CartRepositoryPort) Delete(ctx context.Context, cartID int) error {
	_ = ctx
	return r.inner.DeleteCart(cartID)
}

type OrderRepositoryPort struct{ inner *OrderRepository }

func NewOrderRepositoryPort(inner *OrderRepository) *OrderRepositoryPort {
	return &OrderRepositoryPort{inner: inner}
}

func (r *OrderRepositoryPort) Create(ctx context.Context, order *models.Order) error {
	_ = ctx
	return r.inner.Create(order)
}
func (r *OrderRepositoryPort) CreateOrderItems(ctx context.Context, orderID int, items []*models.OrderItem) error {
	_ = ctx
	for _, item := range items {
		item.OrderID = orderID
		if item.Subtotal == 0 {
			item.Subtotal = item.Price * float64(item.Quantity)
		}
		if err := r.inner.AddOrderItem(item); err != nil {
			return err
		}
	}
	return nil
}
func (r *OrderRepositoryPort) GetByID(ctx context.Context, id int) (*models.Order, error) {
	_ = ctx
	return r.inner.GetByID(id)
}
func (r *OrderRepositoryPort) GetByOrderNumber(ctx context.Context, orderNumber string) (*models.Order, error) {
	_ = ctx
	return r.inner.GetByOrderNumber(orderNumber)
}
func (r *OrderRepositoryPort) GetOrderItems(ctx context.Context, orderID int) ([]*models.OrderItem, error) {
	_ = ctx
	items, err := r.inner.GetOrderItems(orderID)
	if err != nil {
		return nil, err
	}
	out := make([]*models.OrderItem, len(items))
	for i := range items {
		it := items[i]
		out[i] = &it
	}
	return out, nil
}
func (r *OrderRepositoryPort) Update(ctx context.Context, order *models.Order) error {
	_ = ctx
	return r.inner.Update(order)
}
func (r *OrderRepositoryPort) UpdateStatus(ctx context.Context, orderID int, status string) error {
	_ = ctx
	return r.inner.UpdateStatus(orderID, status)
}
func (r *OrderRepositoryPort) GetUserOrders(ctx context.Context, userID int, page, limit int) ([]*models.Order, int, error) {
	_ = ctx
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	orders, total, err := r.inner.GetByUserID(userID, limit, (page-1)*limit)
	if err != nil {
		return nil, 0, err
	}
	out := make([]*models.Order, len(orders))
	for i := range orders {
		o := orders[i]
		out[i] = &o
	}
	return out, int(total), nil
}
func (r *OrderRepositoryPort) GetAllOrders(ctx context.Context, filter *models.OrderFilter) ([]*models.Order, int, error) {
	_ = ctx
	page, limit := 1, 10
	if filter != nil {
		if filter.Page > 0 {
			page = filter.Page
		}
		if filter.Limit > 0 {
			limit = filter.Limit
		}
	}
	offset := (page - 1) * limit
	var orders []models.Order
	var total int64
	var err error
	if filter != nil && filter.Status != "" && filter.UserID > 0 {
		orders, total, err = r.inner.GetByUserIDAndStatus(filter.UserID, filter.Status, limit, offset)
	} else if filter != nil && filter.Status != "" {
		orders, total, err = r.inner.GetByStatus(filter.Status, limit, offset)
	} else if filter != nil && filter.UserID > 0 {
		orders, total, err = r.inner.GetByUserID(filter.UserID, limit, offset)
	} else {
		orders, total, err = r.inner.GetAll(limit, offset)
	}
	if err != nil {
		return nil, 0, err
	}
	out := make([]*models.Order, len(orders))
	for i := range orders {
		o := orders[i]
		out[i] = &o
	}
	return out, int(total), nil
}
func (r *OrderRepositoryPort) GenerateOrderNumber(ctx context.Context) (string, error) {
	_ = ctx
	for i := 0; i < 5; i++ {
		num := fmt.Sprintf("ORD-%d", time.Now().UnixNano())
		exists, err := r.inner.OrderNumberExists(num)
		if err != nil {
			return "", err
		}
		if !exists {
			return num, nil
		}
	}
	return fmt.Sprintf("ORD-%d", time.Now().UnixNano()), nil
}
func (r *OrderRepositoryPort) CalculateTotal(ctx context.Context, items []*models.OrderItem) float64 {
	_ = ctx
	var total float64
	for _, item := range items {
		if item.Subtotal > 0 {
			total += item.Subtotal
		} else {
			total += item.Price * float64(item.Quantity)
		}
	}
	return total
}
func (r *OrderRepositoryPort) Delete(ctx context.Context, id int) error {
	_ = ctx
	return r.inner.Delete(id)
}

type PaymentRepositoryPort struct{ inner *PaymentRepository }

func NewPaymentRepositoryPort(inner *PaymentRepository) *PaymentRepositoryPort {
	return &PaymentRepositoryPort{inner: inner}
}

func (r *PaymentRepositoryPort) Create(ctx context.Context, payment *models.Payment) error {
	_ = ctx
	return r.inner.Create(payment)
}
func (r *PaymentRepositoryPort) GetByID(ctx context.Context, id int) (*models.Payment, error) {
	_ = ctx
	return r.inner.GetByID(id)
}
func (r *PaymentRepositoryPort) GetByOrderID(ctx context.Context, orderID int) (*models.Payment, error) {
	_ = ctx
	return r.inner.GetByOrderID(orderID)
}
func (r *PaymentRepositoryPort) GetByTransactionID(ctx context.Context, transactionID string) (*models.Payment, error) {
	_ = ctx
	return r.inner.GetByTransactionID(transactionID)
}
func (r *PaymentRepositoryPort) Update(ctx context.Context, payment *models.Payment) error {
	_ = ctx
	return r.inner.Update(payment)
}
func (r *PaymentRepositoryPort) UpdateStatus(ctx context.Context, paymentID int, status string) error {
	_ = ctx
	return r.inner.UpdateStatus(paymentID, status)
}
func (r *PaymentRepositoryPort) ConfirmPayment(ctx context.Context, paymentID int, transactionID string) error {
	_ = ctx
	return r.inner.MarkAsPaid(paymentID, transactionID)
}
func (r *PaymentRepositoryPort) List(ctx context.Context, filter *models.PaymentFilter) ([]*models.Payment, int, error) {
	_ = ctx
	page, limit := 1, 10
	if filter != nil {
		if filter.Page > 0 {
			page = filter.Page
		}
		if filter.Limit > 0 {
			limit = filter.Limit
		}
	}
	payments, total, err := r.inner.GetAll(limit, (page-1)*limit)
	if err != nil {
		return nil, 0, err
	}
	out := make([]*models.Payment, len(payments))
	for i := range payments {
		p := payments[i]
		out[i] = &p
	}
	return out, int(total), nil
}
func (r *PaymentRepositoryPort) Delete(ctx context.Context, id int) error {
	_ = ctx
	return r.inner.Delete(id)
}

// Ensure unused import stays valid if sql is needed by callers.
var _ = sql.ErrNoRows
