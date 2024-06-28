package product

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/davidado/go-api-reference/types"
)

// Store : Product store
type Store struct {
	db *sql.DB
}

// NewStore : Create a new product store
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// GetProducts : Get all products
func (s *Store) GetProducts() ([]types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}

	products := make([]types.Product, 0)
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}

	return products, nil
}

// GetProductsByID : Get products by ID
func (s *Store) GetProductsByID(productIDs []int) ([]types.Product, error) {
	placeholders := strings.Repeat(",?", len(productIDs)-1)
	query := fmt.Sprintf("SELECT * FROM products WHERE id IN (?%s)", placeholders)

	// Convert productIDs to []interface{}
	args := make([]interface{}, len(productIDs))
	for i, v := range productIDs {
		args[i] = v
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	products := []types.Product{}
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}

	return products, nil
}

// UpdateProduct : Update a product
func (s *Store) UpdateProduct(product types.Product) error {
	_, err := s.db.Exec("UPDATE products SET name = ?, description = ?, image = ?, price = ?, quantity = ? WHERE id = ?", product.Name, product.Description, product.Image, product.Price, product.Quantity, product.ID)
	return err

}

func scanRowsIntoProduct(rows *sql.Rows) (*types.Product, error) {
	p := &types.Product{}
	err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Image, &p.Price, &p.Quantity, &p.CreatedAt)
	if err != nil {
		return &types.Product{}, err
	}

	return p, nil
}
