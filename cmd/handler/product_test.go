package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gituhb.com/juajosserand/goweb/internal/domain"
	producti "gituhb.com/juajosserand/goweb/internal/product"
	"gituhb.com/juajosserand/goweb/pkg/web"
)

var testProduct = &domain.Product{
	Name:        "Test Product",
	Quantity:    1,
	CodeValue:   "A1B2C3",
	IsPublished: true,
	Expiration:  "20/01/2023",
	Price:       100.0,
}

func init() {
	_, filename, _, _ := runtime.Caller(0)
	rootDir := path.Join(path.Dir(filename), "../..")
	_ = os.Chdir(rootDir)

	_ = godotenv.Load(rootDir + "/.env")
}

func arrange(method string, endpoint string, headers map[string]string, body []byte) (func() *http.Response, error) {
	repo, err := producti.NewRepository()
	if err != nil {
		return nil, err
	}

	svc := producti.NewService(repo)

	mux := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	NewProduct(mux, svc)

	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	for h, v := range headers {
		req.Header.Set(h, v)
	}

	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	res := httptest.NewRecorder()

	return func() *http.Response {
		mux.ServeHTTP(res, req)
		return res.Result()
	}, nil
}

func TestGetAll(t *testing.T) {
	act, err := arrange(http.MethodGet, "/products/", nil, []byte(""))
	if err != nil {
		t.Fatal(err)
	}

	res := act()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetById(t *testing.T) {
	act, err := arrange(http.MethodGet, "/products/1", nil, []byte(""))
	if err != nil {
		t.Fatal(err)
	}

	res := act()

	assert.Equal(t, http.StatusFound, res.StatusCode)
}

func TestCreate(t *testing.T) {
	bytes, err := json.Marshal(testProduct)
	if err != nil {
		t.Fatal(err)
	}

	act, err := arrange(http.MethodPost, "/products/", map[string]string{"token": os.Getenv("TOKEN")}, bytes)
	if err != nil {
		t.Fatal(err)
	}

	res := act()

	assert.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestDelete(t *testing.T) {
	act, err := arrange(http.MethodDelete, "/products/501", map[string]string{"token": os.Getenv("TOKEN")}, []byte(""))
	if err != nil {
		t.Fatal(err)
	}

	res := act()

	assert.Equal(t, http.StatusNoContent, res.StatusCode)
}

func TestBadRequestErr(t *testing.T) {
	tests := []struct {
		method   string
		endpoint string
		headers  map[string]string
		body     *domain.Product
		expected string
	}{
		{
			http.MethodGet,
			"/products/501a",
			nil,
			nil,
			producti.ErrInvalidId.Error(),
		},
		{
			http.MethodPut,
			"/products/501",
			map[string]string{"token": os.Getenv("TOKEN")},
			testProduct,
			producti.ErrInvalidData.Error(),
		},
		{
			http.MethodPatch,
			"/products/501a",
			map[string]string{"token": os.Getenv("TOKEN")},
			testProduct,
			producti.ErrInvalidId.Error(),
		},
		{
			http.MethodDelete,
			"/products/501a",
			map[string]string{"token": os.Getenv("TOKEN")},
			nil,
			producti.ErrInvalidId.Error(),
		},
	}

	for _, test := range tests {
		t.Run("TestBadRequestErr"+test.method, func(t *testing.T) {
			if test.body != nil {
				test.body.Name = ""
			}

			bytes, err := json.Marshal(test.body)
			if err != nil {
				t.Fatal(err)
			}

			act, err := arrange(test.method, test.endpoint, test.headers, bytes)
			if err != nil {
				t.Fatal(err)
			}

			res := act()
			r := web.ErrResponse(0, "", "")
			if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, http.StatusBadRequest, res.StatusCode)
			assert.Equal(t, test.expected, r.Message)
		})
	}
}

func TestNotFoundErr(t *testing.T) {
	tests := []struct {
		method   string
		endpoint string
		headers  map[string]string
		body     *domain.Product
	}{
		{
			http.MethodGet,
			"/products/502",
			nil,
			nil,
		},
		{
			http.MethodPut,
			"/products/502",
			map[string]string{"token": os.Getenv("TOKEN")},
			testProduct,
		},
		{
			http.MethodPatch,
			"/products/502",
			map[string]string{"token": os.Getenv("TOKEN")},
			testProduct,
		},
		{
			http.MethodDelete,
			"/products/502",
			map[string]string{"token": os.Getenv("TOKEN")},
			nil,
		},
	}

	for _, test := range tests {
		t.Run("TestNotFoundErr"+test.method, func(t *testing.T) {
			if test.body != nil {
				test.body.Name = "Test Product Updated"
			}

			bytes, err := json.Marshal(testProduct)
			if err != nil {
				t.Fatal(err)
			}

			act, err := arrange(test.method, test.endpoint, test.headers, bytes)
			if err != nil {
				t.Fatal(err)
			}

			res := act()
			r := web.ErrResponse(0, "", "")
			if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, http.StatusNotFound, res.StatusCode)
			assert.Equal(t, producti.ErrNotFound.Error(), r.Message)
		})
	}
}

func TestUnauthorizedErr(t *testing.T) {
	tests := []struct {
		method   string
		endpoint string
		body     *domain.Product
	}{
		{
			http.MethodPost,
			"/products/",
			testProduct,
		},
		{
			http.MethodPut,
			"/products/501",
			testProduct,
		},
		{
			http.MethodPatch,
			"/products/501",
			testProduct,
		},
		{
			http.MethodDelete,
			"/products/501",
			nil,
		},
	}

	for _, test := range tests {
		t.Run("TestUnauthorizedErr"+test.method, func(t *testing.T) {
			bytes, err := json.Marshal(testProduct)
			if err != nil {
				t.Fatal(err)
			}

			act, err := arrange(test.method, test.endpoint, nil, bytes)
			if err != nil {
				t.Fatal(err)
			}

			res := act()
			r := web.ErrResponse(0, "", "")
			if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
			assert.Equal(t, "invalid token", r.Message)
		})
	}
}
