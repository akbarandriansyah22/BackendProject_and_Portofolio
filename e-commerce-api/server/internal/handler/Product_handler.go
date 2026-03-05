package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/models"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/observability"
	"github.com/akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api/server/internal/ports"
)

type ProductHandler struct {
	productService ports.ProductService
	logger         observability.Logger
}

func NewProductHandler(
	productService ports.ProductService,
	logger observability.Logger,
) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		logger:         logger,
	}
}

// =======================
// PUBLIC
// =======================

// GET /api/products
func (h *ProductHandler) ListProducts(c *fiber.Ctx) error {
	filter := &models.ProductFilter{
		Page:  c.QueryInt("page", 1),
		Limit: c.QueryInt("limit", 10),
	}

	result, err := h.productService.ListProducts(c.Context(), filter)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// GET /api/products/:id
func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid product id")
	}

	product, err := h.productService.GetProductByID(c.Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

// GET /api/products/slug/:slug
func (h *ProductHandler) GetProductBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")

	product, err := h.productService.GetProductBySlug(c.Context(), slug)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

// GET /api/products/search
func (h *ProductHandler) SearchProducts(c *fiber.Ctx) error {
	keyword := c.Query("q")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	result, err := h.productService.SearchProducts(
		c.Context(),
		keyword,
		page,
		limit,
	)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// GET /api/products/category/:id
func (h *ProductHandler) GetProductsByCategory(c *fiber.Ctx) error {
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid category id")
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	result, err := h.productService.GetProductsByCategory(
		c.Context(),
		categoryID,
		page,
		limit,
	)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// =======================
// ADMIN
// =======================

// POST /api/products
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req models.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	product, err := h.productService.CreateProduct(c.Context(), &req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

// PUT /api/products/:id
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid product id")
	}

	var req models.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	if err := h.productService.UpdateProduct(c.Context(), id, &req); err != nil {
		return h.handleError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

// DELETE /api/products/:id
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid product id")
	}

	if err := h.productService.DeleteProduct(c.Context(), id); err != nil {
		return h.handleError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

// PUT /api/products/:id/activate
func (h *ProductHandler) ActivateProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid product id")
	}

	if err := h.productService.ActivateProduct(c.Context(), id); err != nil {
		return h.handleError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

// PUT /api/products/:id/deactivate
func (h *ProductHandler) DeactivateProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid product id")
	}

	if err := h.productService.DeactivateProduct(c.Context(), id); err != nil {
		return h.handleError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

// PATCH /api/products/:id/stock
func (h *ProductHandler) UpdateStock(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid product id")
	}

	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	if err := h.productService.UpdateStock(c.Context(), id, req.Quantity); err != nil {
		return h.handleError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

// =======================
// ERROR HANDLER (PRIVATE)
// =======================

func (h *ProductHandler) handleError(c *fiber.Ctx, err error) error {
	switch err.Error() {
	case "product not found":
		// Aman di-forward ke client — ini user-facing message, bukan internal detail
		return fiber.NewError(fiber.StatusNotFound, "product not found")
	default:
		// Log error internal, tapi jangan expose detail ke client
		h.logger.Error("product handler error", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}
}