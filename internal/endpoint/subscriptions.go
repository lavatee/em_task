package endpoint

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lavatee/subs/internal/model"
)

// @Summary Получение подписок
// @Description Получение подписок с возможной фильтрацией по ID пользователя и названию сервиса
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string false "Фильтрация по ID пользователя"
// @Param service_name query string false "Фильтрация по названию сервиса"
// @Success 200 {object} model.SubscriptionListResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /subscriptions [get]
func (e *Endpoint) GetUserSubscriptions(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	serviceName := ctx.Query("service_name")

	var userUUID uuid.UUID
	var err error

	if userID != "" {
		userUUID, err = uuid.Parse(userID)
		if err != nil {
			e.logger.Warnf("Invalid user ID format: %s", err.Error())
			ctx.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Invalid user ID format",
			})
			return
		}
	}

	subscriptions, err := e.services.Subscriptions.GetUserSubscriptions(ctx, userUUID, serviceName)
	if err != nil {
		e.logger.Errorf("Failed to get subscriptions: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error: "Failed to get subscriptions",
		})
		return
	}

	ctx.JSON(http.StatusOK, model.SubscriptionListResponse{
		Subscriptions: subscriptions,
	})
}

// @Summary Создание подписки
// @Description Создание новой подписки пользователя
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body model.CreateSubscriptionRequest true "Данные о подписке"
// @Success 201 {object} model.SubscriptionResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /subscriptions [post]
func (e *Endpoint) CreateSubscription(ctx *gin.Context) {
	var req model.CreateSubscriptionRequest
	if err := ctx.BindJSON(&req); err != nil {
		e.logger.Warnf("Invalid request body: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	subscription, err := e.services.Subscriptions.CreateSubscription(ctx, req)
	if err != nil {
		e.logger.Errorf("Failed to create subscription: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error: "Failed to create subscription",
		})
		return
	}

	ctx.JSON(http.StatusCreated, model.SubscriptionResponse{
		Subscription: subscription,
	})
}

// @Summary Получение одной подписки
// @Description Получение данных об одной подписке по ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {object} model.SubscriptionResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /subscriptions/{id} [get]
func (e *Endpoint) GetSubscription(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		e.logger.Warnf("Invalid subscription ID format: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: "Invalid subscription ID format",
		})
		return
	}

	subscription, err := e.services.Subscriptions.GetSubscription(ctx, id)
	if err != nil {
		e.logger.Errorf("Failed to get subscription: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error: fmt.Sprintf("Failed to get subscription: %s", err.Error()),
		})
		return
	}

	ctx.JSON(http.StatusOK, model.SubscriptionResponse{
		Subscription: subscription,
	})
}

// @Summary Изменение подписки
// @Description Обновление данных подписки
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Param subscription body model.UpdateSubscriptionRequest true "Данные о подписке"
// @Success 200 {object} model.SubscriptionResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /subscriptions/{id} [put]
func (e *Endpoint) UpdateSubscription(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		e.logger.Warnf("Invalid subscription ID format: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: "Invalid subscription ID format",
		})
		return
	}

	var req model.UpdateSubscriptionRequest
	if err := ctx.BindJSON(&req); err != nil {
		e.logger.Warnf("Invalid request body: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	subscription, err := e.services.Subscriptions.UpdateSubscription(ctx, id, req)
	if err != nil {
		e.logger.Errorf("Failed to update subscription: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error: fmt.Sprintf("Failed to update subscription: %s", err.Error()),
		})
		return
	}

	ctx.JSON(http.StatusOK, model.SubscriptionResponse{
		Subscription: subscription,
	})
}

// @Summary Удаление подписки
// @Description Удаление подписки по ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Success 204
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /subscriptions/{id} [delete]
func (e *Endpoint) DeleteSubscription(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		e.logger.Warnf("Invalid subscription ID format: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: "Invalid subscription ID format",
		})
		return
	}

	err = e.services.Subscriptions.DeleteSubscription(ctx, id)
	if err != nil {
		e.logger.Errorf("Failed to delete subscription: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error: fmt.Sprintf("Failed to delete subscription: %s", err.Error()),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary Получение стоимости всех подписок
// @Description Подсчет суммарной стоимости всех подписок за выбранный период с фильтрацией по id пользователя и названию подписки
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string false "Фильтрация по ID пользователя"
// @Param service_name query string false "Фильтрация по названию сервиса"
// @Param start_date query string false "Начальная дата (MM-YYYY)"
// @Param end_date query string false "Конечная дата (MM-YYYY)"
// @Success 200 {object} model.TotalCostResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /subscriptions/total [get]
func (e *Endpoint) GetTotalCost(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	serviceName := ctx.Query("service_name")
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	var userUUID uuid.UUID
	var err error

	if userID != "" {
		userUUID, err = uuid.Parse(userID)
		if err != nil {
			e.logger.Warnf("Invalid user ID format: %s", err.Error())
			ctx.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Invalid user ID format",
			})
			return
		}
	}

	var startDate, endDate time.Time
	if startDateStr != "" {
		startDate, err = time.Parse("01-2006", startDateStr)
		if err != nil {
			e.logger.Warnf("Invalid start date format: %s", err.Error())
			ctx.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Invalid start date format, expected MM-YYYY",
			})
			return
		}
	}

	if endDateStr != "" {
		endDate, err = time.Parse("01-2006", endDateStr)
		if err != nil {
			e.logger.Warnf("Invalid end date format: %s", err.Error())
			ctx.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Invalid end date format, expected MM-YYYY",
			})
			return
		}
	}

	total, err := e.services.Subscriptions.GetTotalCost(ctx, userUUID, serviceName, startDate, endDate)
	if err != nil {
		e.logger.Errorf("Failed to calculate total cost: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error: "Failed to calculate total cost",
		})
		return
	}

	ctx.JSON(http.StatusOK, model.TotalCostResponse{
		TotalCost: total,
	})
}
