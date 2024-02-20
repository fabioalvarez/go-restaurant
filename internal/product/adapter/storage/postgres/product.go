package postgres

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"go-restaurant/internal/common/adapter/storage/postgres"
	cmdomain "go-restaurant/internal/common/domain"
	cmutil "go-restaurant/internal/common/util"
	"go-restaurant/internal/product/domain"
	"time"
)

/*ProductRepository implements port.ProductRepository interface
 * and provides access to the postgres database
 */
type ProductRepository struct {
	db *postgres.DB
}

// NewProductRepository creates a new product repository instance
func NewProductRepository(db *postgres.DB) *ProductRepository {
	return &ProductRepository{
		db,
	}
}

// CreateProduct creates a new product record in the database
func (pr *ProductRepository) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	query := pr.db.QueryBuilder.Insert("products").
		Columns("category_id", "name", "image", "price", "stock").
		Values(product.CategoryID, product.Name, product.Image, product.Price, product.Stock).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = pr.db.QueryRow(ctx, sql, args...).Scan(
		&product.ID,
		&product.CategoryID,
		&product.SKU,
		&product.Name,
		&product.Stock,
		&product.Price,
		&product.Image,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// GetProductByID retrieves a product record from the database by id
func (pr *ProductRepository) GetProductByID(ctx context.Context, id uint64) (*domain.Product, error) {
	var product domain.Product

	query := pr.db.QueryBuilder.Select("*").
		From("products").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = pr.db.QueryRow(ctx, sql, args...).Scan(
		&product.ID,
		&product.CategoryID,
		&product.SKU,
		&product.Name,
		&product.Stock,
		&product.Price,
		&product.Image,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, cmdomain.ErrDataNotFound
		}
		return nil, err
	}

	return &product, nil
}

// ListProducts retrieves a list of products from the database
func (pr *ProductRepository) ListProducts(ctx context.Context, search string, categoryId, skip, limit uint64) ([]domain.Product, error) {
	var product domain.Product
	var products []domain.Product

	query := pr.db.QueryBuilder.Select("*").
		From("products").
		OrderBy("id").
		Limit(limit).
		Offset((skip - 1) * limit)

	if categoryId != 0 {
		query = query.Where(sq.Eq{"category_id": categoryId})
	}

	if search != "" {
		query = query.Where(sq.ILike{"name": "%" + search + "%"})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := pr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(
			&product.ID,
			&product.CategoryID,
			&product.SKU,
			&product.Name,
			&product.Stock,
			&product.Price,
			&product.Image,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

// UpdateProduct updates a product record in the database
func (pr *ProductRepository) UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	categoryId := cmutil.NullUint64(product.CategoryID)
	name := cmutil.NullString(product.Name)
	image := cmutil.NullString(product.Image)
	price := cmutil.NullFloat64(product.Price)
	stock := cmutil.NullInt64(product.Stock)

	query := pr.db.QueryBuilder.Update("products").
		Set("name", sq.Expr("COALESCE(?, name)", name)).
		Set("category_id", sq.Expr("COALESCE(?, category_id)", categoryId)).
		Set("image", sq.Expr("COALESCE(?, image)", image)).
		Set("price", sq.Expr("COALESCE(?, price)", price)).
		Set("stock", sq.Expr("COALESCE(?, stock)", stock)).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": product.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = pr.db.QueryRow(ctx, sql, args...).Scan(
		&product.ID,
		&product.CategoryID,
		&product.SKU,
		&product.Name,
		&product.Stock,
		&product.Price,
		&product.Image,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// DeleteProduct deletes a product record from the database by id
func (pr *ProductRepository) DeleteProduct(ctx context.Context, id uint64) error {
	query := pr.db.QueryBuilder.Delete("products").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = pr.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
