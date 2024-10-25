package transporthttp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	commonmodel "booking-event/internal/common/model"
	"booking-event/internal/common/util"
	"booking-event/internal/modules/booking/model"
)

func TestEventHttpHandler_RetrieveEventDetail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		eventID          string
		mockEventService func(ctrl *gomock.Controller) *MockEventHandler
		expectedStatus   int
		expectedBody     commonmodel.Response
	}{
		{
			name:    "Successful event retrieval",
			eventID: "1",
			mockEventService: func(ctrl *gomock.Controller) *MockEventHandler {
				mock := NewMockEventHandler(ctrl)
				mock.EXPECT().RetrieveEventDetail(gomock.Any(), 1).Return(&model.Event{ID: 1, Name: "Test Event"}, nil)
				return mock
			},
			expectedStatus: http.StatusOK,
			expectedBody: commonmodel.Response{
				Success: true,
				Data:    &model.Event{ID: 1, Name: "Test Event"},
				Message: "event retrieved",
			},
		},
		{
			name:    "Event not found",
			eventID: "1",
			mockEventService: func(ctrl *gomock.Controller) *MockEventHandler {
				mock := NewMockEventHandler(ctrl)
				mock.EXPECT().RetrieveEventDetail(gomock.Any(), 1).Return(nil, assert.AnError)
				return mock
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: commonmodel.Response{
				Success: false,
				Data:    nil,
				Message: assert.AnError.Error(),
			},
		},
		{
			name:    "Invalid event ID",
			eventID: "invalid",
			mockEventService: func(ctrl *gomock.Controller) *MockEventHandler {
				return NewMockEventHandler(ctrl)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockEventService := tt.mockEventService(ctrl)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = gin.Params{{Key: "event_id", Value: tt.eventID}}
			c.Request, _ = http.NewRequest(http.MethodGet, "/events/"+tt.eventID, nil)

			handler := NewEventHandler(mockEventService)
			handler.(*EventHttpHandler).RetrieveEventDetail(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus != http.StatusBadRequest {
				var response commonmodel.Response
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody.Success, response.Success)
				assert.Equal(t, tt.expectedBody.Message, response.Message)
				if tt.expectedBody.Data != nil {
					bExpected, _ := json.Marshal(tt.expectedBody.Data)
					data := map[string]any{}
					err = json.Unmarshal(bExpected, &data)
					assert.NoError(t, err)
					assert.Equal(t, data, response.Data)
				}
			}
		})
	}
}

func TestEventHttpHandler_QueryEvents(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		body             model.EventQuery
		mockEventService func(ctrl *gomock.Controller) *MockEventHandler
		expectedStatus   int
		expectedBody     commonmodel.Response
	}{
		{
			name: "Successful events query",
			body: model.EventQuery{Pagination: commonmodel.Pagination{Page: 1, Limit: 10}},
			mockEventService: func(ctrl *gomock.Controller) *MockEventHandler {
				mock := NewMockEventHandler(ctrl)
				mock.EXPECT().QueryEvents(gomock.Any(), gomock.Any()).Return([]model.Event{{ID: 1, Name: "Event 1"}, {ID: 2, Name: "Event 2"}}, nil)
				return mock
			},
			expectedStatus: http.StatusOK,
			expectedBody: commonmodel.Response{
				Success: true,
				Data:    []model.Event{{ID: 1, Name: "Event 1"}, {ID: 2, Name: "Event 2"}},
				Message: "events retrieved",
			},
		},
		{
			name: "Failed events query",
			body: model.EventQuery{Pagination: commonmodel.Pagination{Page: 1, Limit: 10}},
			mockEventService: func(ctrl *gomock.Controller) *MockEventHandler {
				mock := NewMockEventHandler(ctrl)
				mock.EXPECT().QueryEvents(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
				return mock
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: commonmodel.Response{
				Success: false,
				Data:    nil,
				Message: assert.AnError.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockEventService := tt.mockEventService(ctrl)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			bodyBytes, _ := json.Marshal(tt.body)
			c.Request, _ = http.NewRequest(http.MethodPost, "/search/events", bytes.NewBuffer(bodyBytes))
			q := c.Request.URL.Query()
			q.Add("limit", "10")
			q.Add("offset", "0")
			c.Request.URL.RawQuery = q.Encode()

			handler := NewEventHandler(mockEventService)
			handler.(*EventHttpHandler).QueryEvents(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response commonmodel.Response
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody.Success, response.Success)
			assert.Equal(t, tt.expectedBody.Message, response.Message)
			if tt.expectedBody.Data != nil {
				bExpected, _ := json.Marshal(tt.expectedBody.Data)
				data := []any{}
				err = json.Unmarshal(bExpected, &data)
				assert.NoError(t, err)
				assert.Equal(t, data, response.Data)
			}
		})
	}
}

func TestEventHttpHandler_CreateEvent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		body             model.CreateEventRequest
		mockEventService func(ctrl *gomock.Controller) *MockEventHandler
		expectedStatus   int
		expectedBody     commonmodel.Response
	}{
		{
			name: "Successful event creation",
			body: model.CreateEventRequest{Name: "New Event", AvailableSeats: 100, StartAt: time.Now().Add(time.Hour * 24), Location: "HCM", Category: model.EventCategoryMusic, Price: 100, ExecutorID: 1},
			mockEventService: func(ctrl *gomock.Controller) *MockEventHandler {
				mock := NewMockEventHandler(ctrl)
				mock.EXPECT().CreateEvent(gomock.Any(), gomock.Any()).Return(nil)
				return mock
			},
			expectedStatus: http.StatusOK,
			expectedBody: commonmodel.Response{
				Success: true,
				Message: "Event created successfully",
			},
		},
		{
			name: "Failed event creation",
			body: model.CreateEventRequest{Name: "New Event", AvailableSeats: 100, StartAt: time.Now().Add(time.Hour * 24), Location: "HCM", Category: model.EventCategoryMusic, Price: 100, ExecutorID: 1},
			mockEventService: func(ctrl *gomock.Controller) *MockEventHandler {
				mock := NewMockEventHandler(ctrl)
				mock.EXPECT().CreateEvent(gomock.Any(), gomock.Any()).Return(assert.AnError)
				return mock
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: commonmodel.Response{
				Success: false,
				Data:    nil,
				Message: assert.AnError.Error(),
			},
		},
		{
			name: "Invalid request",
			mockEventService: func(ctrl *gomock.Controller) *MockEventHandler {
				return NewMockEventHandler(ctrl)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockEventService := tt.mockEventService(ctrl)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			bodyBytes, _ := json.Marshal(tt.body)
			c.Request, _ = http.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Request = c.Request.WithContext(util.SetUserIDContext(c.Request.Context(), 1))

			handler := NewEventHandler(mockEventService)
			handler.(*EventHttpHandler).CreateEvent(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response commonmodel.Response
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

func TestEventHttpHandler_UpdateEvent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		eventID          string
		body             model.UpdateEventRequest
		mockEventService func(ctrl *gomock.Controller) *MockEventHandler
		expectedStatus   int
		expectedBody     commonmodel.Response
	}{
		{
			name:    "Successful event update",
			eventID: "1",
			body:    model.UpdateEventRequest{Status: model.EventStatusActive},
			mockEventService: func(ctrl *gomock.Controller) *MockEventHandler {
				mock := NewMockEventHandler(ctrl)
				mock.EXPECT().UpdateEvent(gomock.Any(), gomock.Any()).Return(nil)
				return mock
			},
			expectedStatus: http.StatusOK,
			expectedBody: commonmodel.Response{
				Success: true,
				Message: "Event updated successfully",
			},
		},
		{
			name:    "Failed event update",
			eventID: "1",
			body:    model.UpdateEventRequest{Status: model.EventStatusActive},
			mockEventService: func(ctrl *gomock.Controller) *MockEventHandler {
				mock := NewMockEventHandler(ctrl)
				mock.EXPECT().UpdateEvent(gomock.Any(), gomock.Any()).Return(assert.AnError)
				return mock
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: commonmodel.Response{
				Success: false,
				Data:    nil,
				Message: assert.AnError.Error(),
			},
		},
		{
			name:    "Invalid event ID",
			eventID: "invalid",
			body:    model.UpdateEventRequest{Status: model.EventStatusActive},
			mockEventService: func(ctrl *gomock.Controller) *MockEventHandler {
				return NewMockEventHandler(ctrl)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockEventService := tt.mockEventService(ctrl)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			bodyBytes, _ := json.Marshal(tt.body)
			c.Request, _ = http.NewRequest(http.MethodPut, "/events/"+tt.eventID, bytes.NewBuffer(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Params = gin.Params{{Key: "event_id", Value: tt.eventID}}
			c.Request = c.Request.WithContext(util.SetUserIDContext(c.Request.Context(), 1))

			handler := NewEventHandler(mockEventService)
			handler.(*EventHttpHandler).UpdateEvent(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus != http.StatusBadRequest {
				var response commonmodel.Response
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}
