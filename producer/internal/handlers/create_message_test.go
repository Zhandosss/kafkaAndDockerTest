package handlers

import (
	"Messaggio/producer/internal/handlers/mocks"
	"Messaggio/producer/internal/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateMessage(t *testing.T) {
	type mockBehavior func(service *mocks.IServices, producer *mocks.IProducer)

	var tests = []struct {
		name         string
		body         string
		mockBehavior mockBehavior
		wantAnswer   string
		wantStatus   int
	}{
		{
			name: "valid request",
			body: `{"content": "test"}`,
			mockBehavior: func(service *mocks.IServices, producer *mocks.IProducer) {
				service.On("SaveMessage", mock.MatchedBy(func(msg *model.Message) bool {
					return msg.Content == "test"
				})).Return("1", nil)
				producer.On("SendMessage", mock.Anything).Return(int32(0), int64(0), nil)
			},
			wantAnswer: `{"id":"1"}` + "\n",
			wantStatus: http.StatusCreated,
		}, {
			name:         "bad request body",
			body:         `{"content": ""`,
			mockBehavior: func(service *mocks.IServices, producer *mocks.IProducer) {},
			wantAnswer:   `{"message":"invalid request body"}` + "\n",
			wantStatus:   http.StatusBadRequest,
		}, {
			name:         "empty content",
			body:         `{"content": ""}`,
			mockBehavior: func(service *mocks.IServices, producer *mocks.IProducer) {},
			wantAnswer:   `{"message":"content is required"}` + "\n",
			wantStatus:   http.StatusBadRequest,
		}, {
			name: "failed to save message",
			body: `{"content": "test"}`,
			mockBehavior: func(service *mocks.IServices, producer *mocks.IProducer) {
				service.On("SaveMessage", mock.Anything).Return("", assert.AnError)
			},
			wantAnswer: `{"message":"failed to save message"}` + "\n",
			wantStatus: http.StatusInternalServerError,
		}, {
			name: "failed to send message",
			body: `{"content": "test"}`,
			mockBehavior: func(service *mocks.IServices, producer *mocks.IProducer) {
				service.On("SaveMessage", mock.Anything).Return("1", nil)
				producer.On("SendMessage", mock.Anything).Return(int32(0), int64(0), assert.AnError)
			},
			wantAnswer: `{"message":"failed to send message"}` + "\n",
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			serviceMock := mocks.NewIServices(t)
			producerMock := mocks.NewIProducer(t)

			test.mockBehavior(serviceMock, producerMock)

			e := echo.New()

			New(e, serviceMock, producerMock)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/messages", strings.NewReader(test.body))
			r.Header.Set("Content-Type", "application/json")

			e.ServeHTTP(w, r)

			assert.Equal(t, test.wantStatus, w.Code)
			assert.Equal(t, test.wantAnswer, w.Body.String())
		})
	}
}
