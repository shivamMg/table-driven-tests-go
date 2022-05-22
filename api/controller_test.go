package api_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/shivamMg/table-driven-tests-go/api"
	"github.com/shivamMg/table-driven-tests-go/api/mock"
)

const (
	testToken = "example-auth-token"
)

func assertEqual(t *testing.T, expected any, actual any) {
	if expected == actual {
		return
	}
	t.Fatalf("%v (expected) != %v (actual)", expected, actual)
}

func TestController_CreateTODO_MethodNotAllowed(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	auth := mock.NewMockAuthenticator(c)
	db := mock.NewMockDatabase(c)
	ctrl := api.NewController(auth, db)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://example.com/todos", nil)

	ctrl.CreateTODO(w, r)

	resp := w.Result()
	respBody, _ := io.ReadAll(resp.Body)

	assertEqual(t, 405, resp.StatusCode)
	assertEqual(t, "method is not POST\n", string(respBody))
}

func TestController_CreateTODO_UnauthenticatedError(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	auth := mock.NewMockAuthenticator(c)
	db := mock.NewMockDatabase(c)
	ctrl := api.NewController(auth, db)

	w := httptest.NewRecorder()
	rBody := bytes.NewBufferString(`{"name": "task1", "category": "cat1"}`)
	r := httptest.NewRequest(http.MethodPost, "http://example.com/todos", rBody)
	r.Header.Add("AuthToken", testToken)

	auth.EXPECT().IsAuthenticated(testToken).Return(false)

	ctrl.CreateTODO(w, r)

	resp := w.Result()
	respBody, _ := io.ReadAll(resp.Body)

	assertEqual(t, 401, resp.StatusCode)
	assertEqual(t, "unauthenticated\n", string(respBody))
}

func TestController_CreateTODO_BadRequestErrors(t *testing.T) {
	testCases := []struct {
		name             string
		requestBody      string
		expectedResponse string
	}{
		{"invalid json", `{"name"}`, "invalid json: invalid character '}' after object key\n"},
		{"empty name", `{"name": ""}`, "invalid todo: empty name\n"},
		{"empty category", `{"name": "task1", "category": ""}`, "invalid todo: empty category\n"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			auth := mock.NewMockAuthenticator(c)
			db := mock.NewMockDatabase(c)
			ctrl := api.NewController(auth, db)

			w := httptest.NewRecorder()
			rBody := bytes.NewBufferString(tc.requestBody)
			r := httptest.NewRequest(http.MethodPost, "http://example.com/todos", rBody)
			r.Header.Add("AuthToken", testToken)

			auth.EXPECT().IsAuthenticated(testToken).Return(true)

			ctrl.CreateTODO(w, r)

			resp := w.Result()
			respBody, _ := io.ReadAll(resp.Body)

			assertEqual(t, 400, resp.StatusCode)
			assertEqual(t, tc.expectedResponse, string(respBody))
		})
	}
}

func TestController_CreateTODO_DBError(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	auth := mock.NewMockAuthenticator(c)
	db := mock.NewMockDatabase(c)
	ctrl := api.NewController(auth, db)

	w := httptest.NewRecorder()
	rBody := bytes.NewBufferString(`{"name": "task1", "category": "cat1"}`)
	r := httptest.NewRequest(http.MethodPost, "http://example.com/todos", rBody)
	r.Header.Add("AuthToken", testToken)

	auth.EXPECT().IsAuthenticated(testToken).Return(true)
	db.EXPECT().CreateTODO(&api.TODO{"task1", "cat1"}).Return(errors.New("failed to commit txn"))

	ctrl.CreateTODO(w, r)

	resp := w.Result()
	respBody, _ := io.ReadAll(resp.Body)

	assertEqual(t, 500, resp.StatusCode)
	assertEqual(t, "db error: failed to commit txn\n", string(respBody))
}

func TestController_CreateTODO_BadTableDrivenTest(t *testing.T) {
	testCases := []struct {
		name string

		method string

		expectAuthCall  bool
		isAuthenticated bool

		expectDBCall bool
		dbError      error

		expectedStatusCode int
		expectedResponse   string
	}{
		{"method not allowed", http.MethodGet, false, false, false, nil, 405, "method is not POST\n"},
		{"unauthenticated", http.MethodPost, true, false, false, nil, 401, "unauthenticated\n"},
		{"db error", http.MethodPost, true, true, true, errors.New("failed to commit txn"), 500, "db error: failed to commit txn\n"},
		{"success", http.MethodPost, true, true, true, nil, 201, "todo created"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			auth := mock.NewMockAuthenticator(c)
			db := mock.NewMockDatabase(c)
			ctrl := api.NewController(auth, db)

			w := httptest.NewRecorder()
			rBody := bytes.NewBufferString(`{"name": "task1", "category": "cat1"}`)
			r := httptest.NewRequest(tc.method, "http://example.com/todos", rBody)
			r.Header.Add("AuthToken", testToken)

			if tc.expectAuthCall {
				auth.EXPECT().IsAuthenticated(testToken).Return(tc.isAuthenticated)
			}
			if tc.expectDBCall {
				db.EXPECT().CreateTODO(&api.TODO{"task1", "cat1"}).Return(tc.dbError)
			}

			ctrl.CreateTODO(w, r)

			resp := w.Result()
			respBody, _ := io.ReadAll(resp.Body)

			assertEqual(t, tc.expectedStatusCode, resp.StatusCode)
			assertEqual(t, tc.expectedResponse, string(respBody))
		})
	}
}

func TestController_CreateTODO_Success(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	auth := mock.NewMockAuthenticator(c)
	db := mock.NewMockDatabase(c)
	ctrl := api.NewController(auth, db)

	w := httptest.NewRecorder()
	rBody := bytes.NewBufferString(`{"name": "task1", "category": "cat1"}`)
	r := httptest.NewRequest(http.MethodPost, "http://example.com/todos", rBody)
	r.Header.Add("AuthToken", testToken)

	auth.EXPECT().IsAuthenticated(testToken).Return(true)
	db.EXPECT().CreateTODO(&api.TODO{"task1", "cat1"}).Return(nil)

	ctrl.CreateTODO(w, r)

	resp := w.Result()
	respBody, _ := io.ReadAll(resp.Body)

	assertEqual(t, 201, resp.StatusCode)
	assertEqual(t, "todo created", string(respBody))
}
