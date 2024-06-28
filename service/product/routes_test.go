package product

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davidado/go-api-reference/types"
	"github.com/gorilla/mux"
)

func TestProductServiceHandlers(t *testing.T) {
	productStore := &mockProductStore{}
	// userStore := &mockUserStore{}
	handler := NewHandler(productStore)

	t.Run("should handle get products", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/products", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/products", handler.handleGetProducts).Methods(http.MethodGet)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

type mockProductStore struct{}

func (m *mockProductStore) GetProducts() ([]types.Product, error) {
	return []types.Product{}, nil
}

func (m *mockProductStore) UpdateProduct(_ types.Product) error {
	return nil
}

func (m *mockProductStore) GetProductsByID(_ []int) ([]types.Product, error) {
	return []types.Product{}, nil
}
