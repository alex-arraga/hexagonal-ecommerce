package repository

import (
	"context"
	"go-ecommerce/internal/adapters/storage/database/postgres/database_dtos"
	"go-ecommerce/internal/adapters/storage/database/postgres/models"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) ports.ProductRepository {
	return &ProductRepo{db: db}
}

// SaveProduct implements ports.ProductRepository.
func (pr *ProductRepo) SaveProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	productDb := database_dtos.ConvertProductDomainToModel(product)

	// if exist product.ID update, else create new product
	if product.ID != uuid.Nil {
		if result := pr.db.WithContext(ctx).Where("id = ?", product.ID).Updates(productDb); result.Error != nil {
			if result.RowsAffected == 0 {
				return nil, domain.ErrProductNotFound
			}
			return nil, result.Error
		}
	} else {
		if result := pr.db.WithContext(ctx).Create(productDb); result.Error != nil {
			return nil, result.Error
		}
	}

	productDomain := database_dtos.ConvertProductModelToDomain(productDb)
	return productDomain, nil
}

// GetProductById implements ports.ProductRepository.
func (pr *ProductRepo) GetProductById(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	var productDb *models.ProductModel

	if result := pr.db.WithContext(ctx).First(productDb, "id = ?", id); result.Error != nil {
		if result.RowsAffected == 0 {
			return nil, domain.ErrProductNotFound
		}
		return nil, result.Error
	}

	productDomain := database_dtos.ConvertProductModelToDomain(productDb)
	return productDomain, nil
}

// ListProducts implements ports.ProductRepository.
func (pr *ProductRepo) ListProducts(ctx context.Context) ([]*domain.Product, error) {
	var productsDb []*models.ProductModel

	if result := pr.db.WithContext(ctx).Find(productsDb); result.Error != nil {
		return nil, result.Error
	}

	productsDomain := database_dtos.ConvertProductsModelsToDomain(productsDb)
	return productsDomain, nil
}

// DeleteProduct implements ports.ProductRepository.
func (pr *ProductRepo) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	var productDb *models.ProductModel

	if result := pr.db.WithContext(ctx).Delete(productDb, "id = ?", id); result.Error != nil {
		return result.Error
	}

	return nil
}
