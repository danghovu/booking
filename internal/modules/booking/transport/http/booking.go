//go:generate mockgen -source=booking.go -destination=booking_mock.go -package=transporthttp
package transporthttp

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"booking-event/internal/common/handler"
	commonmodel "booking-event/internal/common/model"
	"booking-event/internal/common/util"
	"booking-event/internal/modules/booking/model"
)

type BookingHandler interface {
	CreateBooking(ctx context.Context, booking model.CreateBookingRequest) (*model.Booking, error)
	ConfirmBooking(ctx context.Context, userID int, bookingID int) error
	CancelBooking(ctx context.Context, id int, executorID int) error
	GetBookingByID(ctx context.Context, id int) (*model.Booking, error)
}

type BookingHttpHandler struct {
	bookingService BookingHandler
}

func NewBookingHandler(bookingService BookingHandler) handler.HttpHandler {
	return &BookingHttpHandler{bookingService: bookingService}
}

func (h *BookingHttpHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/bookings", h.CreateBooking)
	router.PUT("/bookings/:booking_id/confirm", h.ConfirmBooking)
	router.PUT("/bookings/:booking_id/cancel", h.CancelBooking)
	router.GET("/bookings/:booking_id", h.GetBookingByID)
}

func (h *BookingHttpHandler) CreateBooking(c *gin.Context) {
	var booking model.CreateBookingRequest
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, commonmodel.Response{
			Success: false,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	booking.UserID = util.GetUserIDContext(c.Request.Context())
	resp, err := h.bookingService.CreateBooking(c.Request.Context(), booking)
	if err != nil {
		c.JSON(http.StatusInternalServerError, commonmodel.Response{
			Success: false,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, commonmodel.Response{
		Success: true,
		Message: "booking created",
		Data:    resp,
	})
}

func (h *BookingHttpHandler) ConfirmBooking(c *gin.Context) {
	var request model.ConfirmBookingRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, commonmodel.Response{
			Success: false,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	userID := util.GetUserIDContext(c.Request.Context())
	err := h.bookingService.ConfirmBooking(c.Request.Context(), userID, request.BookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, commonmodel.Response{
			Success: false,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, commonmodel.Response{
		Success: true,
		Message: "booking confirmed",
	})
}

func (h *BookingHttpHandler) CancelBooking(c *gin.Context) {
	var request model.CancelBookingRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, commonmodel.Response{
			Success: false,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	userID := util.GetUserIDContext(c.Request.Context())
	err := h.bookingService.CancelBooking(c.Request.Context(), request.BookingID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, commonmodel.Response{
			Success: false,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, commonmodel.Response{
		Success: true,
		Message: "booking canceled",
	})
}

func (h *BookingHttpHandler) GetBookingByID(c *gin.Context) {
	var request model.GetBookingByIDRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, commonmodel.Response{
			Success: false,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	booking, err := h.bookingService.GetBookingByID(c.Request.Context(), request.BookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, commonmodel.Response{
			Success: false,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, commonmodel.Response{
		Success: true,
		Data:    booking,
	})
}
