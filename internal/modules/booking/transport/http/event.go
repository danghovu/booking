//go:generate mockgen -source=event.go -destination=event_mock.go -package=transporthttp
package transporthttp

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	_errors "booking-event/internal/common/errors"
	"booking-event/internal/common/handler"
	commonmodel "booking-event/internal/common/model"
	"booking-event/internal/common/util"
	"booking-event/internal/modules/booking/model"
)

type EventHandler interface {
	RetrieveEventDetail(ctx context.Context, eventID int) (*model.Event, error)
	QueryEvents(ctx context.Context, query model.EventQuery) ([]model.Event, error)
	CreateEvent(ctx context.Context, params model.CreateEventRequest) error
	UpdateEvent(ctx context.Context, params model.UpdateEventRequest) error
}

type EventHttpHandler struct {
	eventService EventHandler
}

func NewEventHandler(eventService EventHandler) handler.HttpHandler {
	return &EventHttpHandler{eventService: eventService}
}

func (h *EventHttpHandler) RetrieveEventDetail(c *gin.Context) {
	var request model.RetrieveEventDetailRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, commonmodel.Response{
			Success: false,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	event, err := h.eventService.RetrieveEventDetail(c.Request.Context(), request.EventID)
	if errors.Is(err, _errors.ErrNotFound) {
		c.JSON(http.StatusNotFound, commonmodel.Response{
			Success: false,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
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
		Data:    event,
		Message: "event retrieved",
	})
}

func (h *EventHttpHandler) QueryEvents(c *gin.Context) {
	var query model.EventQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, commonmodel.Response{
			Success: false,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	if query.StartFrom.After(query.StartTo) && !query.StartTo.IsZero() {
		c.JSON(http.StatusBadRequest, commonmodel.Response{
			Success: false,
			Data:    nil,
			Message: "start_from must be before start_to",
		})
		return
	}
	events, err := h.eventService.QueryEvents(c.Request.Context(), query)
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
		Data:    events,
		Message: "events retrieved",
	})
}

func (h *EventHttpHandler) CreateEvent(c *gin.Context) {
	var request model.CreateEventRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commonmodel.Response{
			Success: false,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	request.ExecutorID = util.GetUserIDContext(c.Request.Context())

	err := h.eventService.CreateEvent(c.Request.Context(), request)
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
		Message: "Event created successfully",
	})
}

func (h *EventHttpHandler) UpdateEvent(c *gin.Context) {
	var request model.UpdateEventRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commonmodel.Response{
			Success: false,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	var err error
	request.EventID, err = strconv.Atoi(c.Param("event_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, commonmodel.Response{
			Success: false,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	request.ExecutorID = util.GetUserIDContext(c.Request.Context())

	err = h.eventService.UpdateEvent(c.Request.Context(), request)
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
		Message: "Event updated successfully",
	})
}

func (h *EventHttpHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/events/:event_id", h.RetrieveEventDetail)
	router.POST("/search/events", h.QueryEvents)
	router.POST("/events", h.CreateEvent)
	router.PUT("/events/:event_id", h.UpdateEvent)
}
