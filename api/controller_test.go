package api_test

import (
	"bytes"
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
