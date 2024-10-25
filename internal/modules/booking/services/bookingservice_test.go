package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"booking-event/internal/modules/booking/model"
)

func TestBookingService_CreateBooking(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                  string
		request               model.CreateBookingRequest
		mockEventService      func(ctrl *gomock.Controller) *MockEventServiceForBooking
		mockBookingRepo       func(ctrl *gomock.Controller) *MockBookingRepository
		mockEventTokenService func(ctrl *gomock.Controller) *MockBookingEventTokenService
		expectedResponse      *model.Booking
		expectedError         error
	}{
		{
			name: "Successful booking creation",
			request: model.CreateBookingRequest{
				EventID:  1,
				UserID:   1,
				Quantity: 2,
			},
			mockEventService: func(ctrl *gomock.Controller) *MockEventServiceForBooking {
				mock := NewMockEventServiceForBooking(ctrl)
				mock.EXPECT().GetEventByID(gomock.Any(), 1).Return(&model.Event{Status: model.EventStatusActive}, nil)
				return mock
			},
			mockBookingRepo: func(ctrl *gomock.Controller) *MockBookingRepository {
				mock := NewMockBookingRepository(ctrl)
				mock.EXPECT().CountBookingByUserID(gomock.Any(), 1, 1).Return(0, nil)
				mock.EXPECT().CreateBooking(gomock.Any(), &model.Booking{
					Status:          model.BookingStatusPending,
					UserID:          1,
					EventID:         1,
					InitialQuantity: 2,
					Quantity:        2,
				}, []model.BookingItem{
					{Token: "token1"},
					{Token: "token2"},
				}).DoAndReturn(func(ctx context.Context, booking *model.Booking, bookingItems []model.BookingItem) error {
					booking.ID = 1
					return nil
				})
				return mock
			},
			mockEventTokenService: func(ctrl *gomock.Controller) *MockBookingEventTokenService {
				mock := NewMockBookingEventTokenService(ctrl)
				mock.EXPECT().SelectAvailableToken(gomock.Any(), int32(1), 1, 2).Return([]string{"token1", "token2"}, nil)
				return mock
			},

			expectedResponse: &model.Booking{ID: 1, Status: model.BookingStatusPending, UserID: 1, EventID: 1, InitialQuantity: 2, Quantity: 2},
			expectedError:    nil,
		},
		{
			name: "Event not active",
			request: model.CreateBookingRequest{
				EventID:  1,
				UserID:   1,
				Quantity: 2,
			},
			mockEventService: func(ctrl *gomock.Controller) *MockEventServiceForBooking {
				mock := NewMockEventServiceForBooking(ctrl)
				mock.EXPECT().GetEventByID(gomock.Any(), 1).Return(&model.Event{Status: model.EventStatusInactive}, nil)
				return mock
			},
			expectedError: errors.New("event is not active"),
		},
		{
			name: "Max booking per user reached",
			request: model.CreateBookingRequest{
				EventID:  1,
				UserID:   1,
				Quantity: 2,
			},
			mockEventService: func(ctrl *gomock.Controller) *MockEventServiceForBooking {
				mock := NewMockEventServiceForBooking(ctrl)
				mock.EXPECT().GetEventByID(gomock.Any(), 1).Return(&model.Event{Status: model.EventStatusActive}, nil)
				return mock
			},
			mockBookingRepo: func(ctrl *gomock.Controller) *MockBookingRepository {
				mock := NewMockBookingRepository(ctrl)
				mock.EXPECT().CountBookingByUserID(gomock.Any(), 1, 1).Return(2, nil)
				return mock
			},
			expectedError: errors.New("max booking per user reached"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var mockEventService *MockEventServiceForBooking
			if tt.mockEventService != nil {
				mockEventService = tt.mockEventService(ctrl)
			}
			var mockEventTokenService *MockBookingEventTokenService
			if tt.mockEventTokenService != nil {
				mockEventTokenService = tt.mockEventTokenService(ctrl)
			}
			var mockBookingRepo *MockBookingRepository
			if tt.mockBookingRepo != nil {
				mockBookingRepo = tt.mockBookingRepo(ctrl)
			}

			service := NewBookingService(
				mockEventService,
				mockEventTokenService,
				mockBookingRepo,
				nil,
				nil,
				BookingConfig{MaxBookingPerUser: 2},
			)
			resp, err := service.CreateBooking(context.Background(), tt.request)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, resp)
			}
		})
	}
}

func TestBookingService_ConfirmBooking(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		userID              int
		bookingID           int
		mockEventService    func(ctrl *gomock.Controller) *MockEventServiceForBooking
		mockBookingRepo     func(ctrl *gomock.Controller) *MockBookingRepository
		mockBookingItemRepo func(ctrl *gomock.Controller) *MockBookingItemRepository
		expectedError       error
	}{
		{
			name:      "Successful booking confirmation",
			userID:    1,
			bookingID: 1,
			mockEventService: func(ctrl *gomock.Controller) *MockEventServiceForBooking {
				mock := NewMockEventServiceForBooking(ctrl)
				mock.EXPECT().GetEventByID(gomock.Any(), 1).Return(&model.Event{}, nil)
				return mock
			},
			mockBookingRepo: func(ctrl *gomock.Controller) *MockBookingRepository {
				mock := NewMockBookingRepository(ctrl)
				mock.EXPECT().GetBookingByID(gomock.Any(), 1).Return(&model.Booking{ID: 1, UserID: 1, EventID: 1, Status: model.BookingStatusPending}, nil)
				mock.EXPECT().ConfirmBooking(gomock.Any(), &model.Booking{ID: 1, UserID: 1, EventID: 1, Status: model.BookingStatusConfirmed}, &model.Event{}, []model.BookingItem{}).Return(nil)
				return mock
			},
			mockBookingItemRepo: func(ctrl *gomock.Controller) *MockBookingItemRepository {
				mock := NewMockBookingItemRepository(ctrl)
				mock.EXPECT().GetBookingItemsByBookingID(gomock.Any(), 1).Return([]model.BookingItem{}, nil)
				return mock
			},
			expectedError: nil,
		},
		{
			name:      "Unauthorized user",
			userID:    2,
			bookingID: 1,
			mockBookingRepo: func(ctrl *gomock.Controller) *MockBookingRepository {
				mock := NewMockBookingRepository(ctrl)
				mock.EXPECT().GetBookingByID(gomock.Any(), 1).Return(&model.Booking{UserID: 1, Status: model.BookingStatusPending}, nil)
				return mock
			},
			expectedError: errors.New("unauthorized user is not allowed to confirm this booking"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var mockEventService *MockEventServiceForBooking
			if tt.mockEventService != nil {
				mockEventService = tt.mockEventService(ctrl)
			}
			var mockBookingRepo *MockBookingRepository
			if tt.mockBookingRepo != nil {
				mockBookingRepo = tt.mockBookingRepo(ctrl)
			}
			var mockBookingItemRepo *MockBookingItemRepository
			if tt.mockBookingItemRepo != nil {
				mockBookingItemRepo = tt.mockBookingItemRepo(ctrl)
			}

			service := NewBookingService(
				mockEventService,
				nil,
				mockBookingRepo,
				mockBookingItemRepo,
				nil,
				BookingConfig{},
			)

			err := service.ConfirmBooking(context.Background(), tt.userID, tt.bookingID)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBookingService_CancelBooking(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEventService := NewMockEventServiceForBooking(ctrl)
	mockBookingRepo := NewMockBookingRepository(ctrl)

	service := NewBookingService(
		mockEventService,
		nil,
		mockBookingRepo,
		nil,
		nil,
		BookingConfig{},
	)

	tests := []struct {
		name          string
		bookingID     int
		executorID    int
		setupMocks    func()
		expectedError error
	}{
		{
			name:       "Successful booking cancellation",
			bookingID:  1,
			executorID: 1,
			setupMocks: func() {
				mockBookingRepo.EXPECT().GetBookingByID(gomock.Any(), 1).Return(&model.Booking{ID: 1, UserID: 1, EventID: 1, Status: model.BookingStatusConfirmed}, nil)
				mockEventService.EXPECT().GetEventByID(gomock.Any(), 1).Return(&model.Event{StartAt: time.Now().Add(24 * time.Hour)}, nil)
				mockBookingRepo.EXPECT().CancelBooking(gomock.Any(), 1).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:       "Unauthorized user",
			bookingID:  1,
			executorID: 2,
			setupMocks: func() {
				mockBookingRepo.EXPECT().GetBookingByID(gomock.Any(), 1).Return(&model.Booking{ID: 1, UserID: 1}, nil)
			},
			expectedError: errors.New("unauthorized user is not allowed to cancel this booking"),
		},
		{
			name:       "Already canceled booking",
			bookingID:  1,
			executorID: 1,
			setupMocks: func() {
				mockBookingRepo.EXPECT().GetBookingByID(gomock.Any(), 1).Return(&model.Booking{ID: 1, UserID: 1, Status: model.BookingStatusCanceled}, nil)
			},
			expectedError: errors.New("booking is already canceled"),
		},
		{
			name:       "Event already started",
			bookingID:  1,
			executorID: 1,
			setupMocks: func() {
				mockBookingRepo.EXPECT().GetBookingByID(gomock.Any(), 1).Return(&model.Booking{ID: 1, UserID: 1, EventID: 1, Status: model.BookingStatusConfirmed}, nil)
				mockEventService.EXPECT().GetEventByID(gomock.Any(), 1).Return(&model.Event{StartAt: time.Now().Add(-1 * time.Hour)}, nil)
			},
			expectedError: errors.New("event is already started"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := service.CancelBooking(context.Background(), tt.bookingID, tt.executorID)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
