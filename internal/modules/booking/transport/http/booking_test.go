package transporthttp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	commonmodel "booking-event/internal/common/model"
	"booking-event/internal/common/util"
	"booking-event/internal/modules/booking/model"
)

func TestBookingHttpHandler_CreateBooking(t *testing.T) {

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		body               model.CreateBookingRequest
		mockBookingService func(ctrl *gomock.Controller) *MockBookingHandler
		expectedStatus     int
		expectedBody       commonmodel.Response
	}{
		{
			name: "Successful booking creation",
			body: model.CreateBookingRequest{EventID: 1, UserID: 1, Quantity: 2},
			mockBookingService: func(ctrl *gomock.Controller) *MockBookingHandler {
				mock := NewMockBookingHandler(ctrl)
				mock.EXPECT().CreateBooking(gomock.Any(), gomock.Any()).Return(&model.Booking{ID: 1, Status: model.BookingStatusPending, UserID: 1, EventID: 1, InitialQuantity: 2, Quantity: 2}, nil)
				return mock
			},
			expectedStatus: http.StatusOK,
			expectedBody: commonmodel.Response{
				Success: true,
				Data:    &model.Booking{ID: 1, Status: model.BookingStatusPending, UserID: 1, EventID: 1, InitialQuantity: 2, Quantity: 2},
				Message: "booking created successfully",
			},
		},
		{
			name: "Failed booking creation",
			body: model.CreateBookingRequest{EventID: 1, UserID: 1, Quantity: 2},
			mockBookingService: func(ctrl *gomock.Controller) *MockBookingHandler {
				mock := NewMockBookingHandler(ctrl)
				mock.EXPECT().CreateBooking(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
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
			name:           "Invalid request",
			body:           model.CreateBookingRequest{},
			expectedStatus: http.StatusBadRequest,
			mockBookingService: func(ctrl *gomock.Controller) *MockBookingHandler {
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockBookingService := tt.mockBookingService(ctrl)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			bodyBytes, _ := json.Marshal(tt.body)
			c.Request, _ = http.NewRequest(http.MethodPost, "/bookings", bytes.NewBuffer(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			handler := NewBookingHandler(mockBookingService)
			handler.(*BookingHttpHandler).CreateBooking(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response commonmodel.Response
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			bExpected, _ := json.Marshal(tt.expectedBody.Data)
			data := map[string]any{}
			err = json.Unmarshal(bExpected, &data)
			assert.NoError(t, err)

			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, data, response.Data)
			}
		})
	}
}

func TestBookingHttpHandler_ConfirmBooking(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		bookingID          string
		userID             int
		mockBookingService func(ctrl *gomock.Controller) *MockBookingHandler
		expectedStatus     int
		expectedBody       commonmodel.Response
	}{
		{
			name:      "Successful booking confirmation",
			bookingID: "1",
			userID:    1,
			mockBookingService: func(ctrl *gomock.Controller) *MockBookingHandler {
				mock := NewMockBookingHandler(ctrl)
				mock.EXPECT().ConfirmBooking(gomock.Any(), 1, 1).Return(nil)
				return mock
			},
			expectedStatus: http.StatusOK,
			expectedBody: commonmodel.Response{
				Success: true,
				Message: "booking confirmed",
			},
		},
		{
			name:      "Failed booking confirmation",
			bookingID: "1",
			userID:    1,
			mockBookingService: func(ctrl *gomock.Controller) *MockBookingHandler {
				mock := NewMockBookingHandler(ctrl)
				mock.EXPECT().ConfirmBooking(gomock.Any(), 1, 1).Return(assert.AnError)
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
			name:      "Invalid booking ID",
			bookingID: "invalid",
			userID:    1,
			mockBookingService: func(ctrl *gomock.Controller) *MockBookingHandler {
				return NewMockBookingHandler(ctrl)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockBookingService := tt.mockBookingService(ctrl)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = gin.Params{{Key: "booking_id", Value: tt.bookingID}}
			c.Request, _ = http.NewRequest(http.MethodPut, "/bookings/"+tt.bookingID+"/confirm", nil)
			c.Request = c.Request.WithContext(util.SetUserIDContext(c.Request.Context(), tt.userID))

			handler := NewBookingHandler(mockBookingService)
			handler.(*BookingHttpHandler).ConfirmBooking(c)

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

func TestBookingHttpHandler_CancelBooking(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		bookingID          string
		userID             int
		mockBookingService func(ctrl *gomock.Controller) *MockBookingHandler
		expectedStatus     int
		expectedBody       commonmodel.Response
	}{
		{
			name:      "Successful booking cancellation",
			bookingID: "1",
			userID:    1,
			mockBookingService: func(ctrl *gomock.Controller) *MockBookingHandler {
				mock := NewMockBookingHandler(ctrl)
				mock.EXPECT().CancelBooking(gomock.Any(), 1, 1).Return(nil)
				return mock
			},
			expectedStatus: http.StatusOK,
			expectedBody: commonmodel.Response{
				Success: true,
				Message: "booking canceled",
			},
		},
		{
			name:      "Failed booking cancellation",
			bookingID: "1",
			userID:    1,
			mockBookingService: func(ctrl *gomock.Controller) *MockBookingHandler {
				mock := NewMockBookingHandler(ctrl)
				mock.EXPECT().CancelBooking(gomock.Any(), 1, 1).Return(assert.AnError)
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
			name:      "Invalid booking ID",
			bookingID: "invalid",
			userID:    1,
			mockBookingService: func(ctrl *gomock.Controller) *MockBookingHandler {
				return NewMockBookingHandler(ctrl)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockBookingService := tt.mockBookingService(ctrl)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = gin.Params{{Key: "booking_id", Value: tt.bookingID}}
			c.Request, _ = http.NewRequest(http.MethodPut, "/bookings/"+tt.bookingID+"/cancel", nil)
			c.Request = c.Request.WithContext(util.SetUserIDContext(c.Request.Context(), tt.userID))

			handler := NewBookingHandler(mockBookingService)
			handler.(*BookingHttpHandler).CancelBooking(c)

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

func TestBookingHttpHandler_GetBookingByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		bookingID          string
		mockBookingService func(ctrl *gomock.Controller) *MockBookingHandler
		expectedStatus     int
		expectedBody       commonmodel.Response
	}{
		{
			name:      "Successful get booking",
			bookingID: "1",
			mockBookingService: func(ctrl *gomock.Controller) *MockBookingHandler {
				mock := NewMockBookingHandler(ctrl)
				mock.EXPECT().GetBookingByID(gomock.Any(), 1).Return(&model.Booking{ID: 1, Status: model.BookingStatusConfirmed}, nil)
				return mock
			},
			expectedStatus: http.StatusOK,
			expectedBody: commonmodel.Response{
				Success: true,
				Data:    &model.Booking{ID: 1, Status: model.BookingStatusConfirmed},
			},
		},
		{
			name:      "Booking not found",
			bookingID: "1",
			mockBookingService: func(ctrl *gomock.Controller) *MockBookingHandler {
				mock := NewMockBookingHandler(ctrl)
				mock.EXPECT().GetBookingByID(gomock.Any(), 1).Return(nil, assert.AnError)
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
			name:      "Invalid booking ID",
			bookingID: "invalid",
			mockBookingService: func(ctrl *gomock.Controller) *MockBookingHandler {
				return NewMockBookingHandler(ctrl)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockBookingService := tt.mockBookingService(ctrl)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = gin.Params{{Key: "booking_id", Value: tt.bookingID}}
			c.Request, _ = http.NewRequest(http.MethodGet, "/bookings/"+tt.bookingID, nil)

			handler := NewBookingHandler(mockBookingService)
			handler.(*BookingHttpHandler).GetBookingByID(c)

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
