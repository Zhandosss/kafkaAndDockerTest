package handlers

import (
	"Messaggio/producer/internal/handlers/mocks"
	"Messaggio/producer/internal/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMessage(t *testing.T) {
	type MockBehavior func(service *mocks.IServices)

	var tests = []struct {
		name         string
		id           string
		mockBehavior MockBehavior
		wantAnswer   string
		wantStatus   int
	}{
		{
			name: "valid request",
			id:   "1",
			mockBehavior: func(service *mocks.IServices) {
				service.On("GetMessage", mock.Anything).Return(&model.Message{
					ID:      "1",
					Content: "test",
				}, nil)
			},
			wantAnswer: `{"id":"1","content":"test","is_processed":false,"create_time":"0001-01-01T00:00:00Z"}` + "\n",
			wantStatus: http.StatusOK,
		}, {
			name: "failed to get message",
			id:   "1",
			mockBehavior: func(service *mocks.IServices) {
				service.On("GetMessage", mock.Anything).Return(nil, assert.AnError)
			},
			wantAnswer: `{"message":"failed to get message"}` + "\n",
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			serviceMock := mocks.NewIServices(t)

			test.mockBehavior(serviceMock)

			e := echo.New()

			New(e, serviceMock, nil)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/api/messages/"+test.id, nil)
			r.Header.Set("Content-Type", "application/json")

			e.ServeHTTP(w, r)

			assert.Equal(t, test.wantStatus, w.Code)
			assert.Equal(t, test.wantAnswer, w.Body.String())
		})
	}
}

func TestMessages(t *testing.T) {
	type MockBehavior func(service *mocks.IServices)

	var tests = []struct {
		name         string
		mockBehavior MockBehavior
		wantAnswer   string
		wantStatus   int
	}{
		{
			name: "valid request",
			mockBehavior: func(service *mocks.IServices) {
				service.On("GetMessages").Return([]*model.Message{
					{
						ID:      "1",
						Content: "test",
					},
				}, nil)
			},
			wantAnswer: `[{"id":"1","content":"test","is_processed":false,"create_time":"0001-01-01T00:00:00Z"}]` + "\n",
			wantStatus: http.StatusOK,
		}, {
			name: "failed to get messages",
			mockBehavior: func(service *mocks.IServices) {
				service.On("GetMessages").Return(nil, assert.AnError)
			},
			wantAnswer: `{"message":"failed to get messages"}` + "\n",
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			serviceMock := mocks.NewIServices(t)

			test.mockBehavior(serviceMock)

			e := echo.New()

			New(e, serviceMock, nil)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/api/messages", nil)
			r.Header.Set("Content-Type", "application/json")

			e.ServeHTTP(w, r)

			assert.Equal(t, test.wantStatus, w.Code)
			assert.Equal(t, test.wantAnswer, w.Body.String())
		})
	}
}
