package httphandler

import (
	"authservice/internal/service"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// AdminGetUserInfo godoc
// @Summary Получить информацию о пользователе
// @Description Возвращает полную информацию о пользователе по его ID. Доступно только для администраторов.
// @Tags admin
// @Accept  json
// @Produce  json
// @Param User-ID header string true "ID пользователя"
// @Param Authorization header string true "Auth token" "Токен авторизации"
// @Param user_id query string true "ID пользователя"
// @Success 200 {object} domain.User "Полная информация о пользователе"
// @Failure 400 {object} HTTPResponse "Неверный запрос"
// @Failure 401 {object} HTTPResponse "Пользователь не авторизован"
// @Failure 404 {object} HTTPResponse "Пользователь не найден"
// @Router /admin/get_user_info [get]
// @Security ApiKeyAuth
func AdminGetUserInfo(resp http.ResponseWriter, req *http.Request) {

	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	id := req.URL.Query().Get("user_id")
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		respBody.SetError(errors.New("invalid input"))
		return
	}

	info, err := service.GetUserFullInfo(userID)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		respBody.SetError(err)
	}

	respBody.SetData(info)
}
