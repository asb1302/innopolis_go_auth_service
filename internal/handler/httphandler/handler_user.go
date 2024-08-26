package httphandler

import (
	"authservice/internal/domain"
	appError "authservice/internal/error"
	"authservice/internal/service"
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"strings"
)

// SignUp
// @Summary Регистрация пользователя
// @Description Регистрирует нового пользователя и возвращает токен.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   input  body      domain.LoginPassword  true  "Логин и пароль"
// @Success 200     {object} HTTPResponse{data=string}   "Токен пользователя"
// @Failure 400     {object} HTTPResponse               "Неправильные входные данные"
// @Failure 409     {object} HTTPResponse               "Пользователь уже существует"
// @Router /sign_up [post]
func SignUp(resp http.ResponseWriter, req *http.Request) {

	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	var input domain.LoginPassword
	if err := readBody(req, &input); err != nil {
		resp.WriteHeader(http.StatusUnprocessableEntity)
		respBody.SetError(err)
		return
	}

	if !input.IsValid() {
		resp.WriteHeader(http.StatusBadRequest)
		respBody.SetError(errors.New("invalid input"))
		return
	}

	userToken, err := service.SignUp(&input)
	if err != nil {
		resp.WriteHeader(http.StatusConflict)
		respBody.SetError(err)
		return
	}

	respBody.SetData(userToken)
}

// SignIn
// @Summary Авторизация пользователя
// @Description Авторизует пользователя и возвращает токен.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   input  body      domain.LoginPassword  true  "Логин и пароль"
// @Success 200     {object} HTTPResponse{data=string}   "Токен пользователя"
// @Failure 400     {object} HTTPResponse               "Неправильные входные данные"
// @Failure 404     {object} HTTPResponse               "Пользователь не найден"
// @Router /sign_in [post]
func SignIn(resp http.ResponseWriter, req *http.Request) {

	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	var input domain.LoginPassword
	if err := readBody(req, &input); err != nil {
		resp.WriteHeader(http.StatusUnprocessableEntity)
		respBody.SetError(err)
		return
	}

	if !input.IsValid() {
		resp.WriteHeader(http.StatusBadRequest)
		respBody.SetError(errors.New("invalid input"))
		return
	}

	userToken, err := service.SignIn(&input)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		respBody.SetError(err)
		return
	}

	respBody.SetData(userToken)
}

// GetUserInfo
// @Summary Получить информацию о пользователе
// @Description Возвращает короткую информацию о пользователе на основе его ID.
// @Tags users
// @Produce  json
// @Success 200 {object} HTTPResponse{data=domain.UserInfo} "Информация о пользователе"
// @Param User-ID header string true "ID пользователя"
// @Param Authorization header string true "Auth token" "Токен авторизации"
// @Failure 401 {object} HTTPResponse "Пользователь не авторизован"
// @Failure 404 {object} HTTPResponse                       "Пользователь не найден"
// @Router /get_user_info [get]
// @Security ApiKeyAuth
func GetUserInfo(resp http.ResponseWriter, req *http.Request) {

	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	userID, _ := primitive.ObjectIDFromHex(req.Header.Get(HeaderUserID))

	info, err := service.GetUserShortInfo(userID)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		respBody.SetError(err)
	}

	respBody.SetData(info)
}

// SetUserInfo
// @Summary Обновить информацию о пользователе
// @Description Обновляет информацию о пользователе (имя).
// @Tags users
// @Accept  json
// @Produce  json
// @Param   input  body      SetUserInfoReq  true  "Информация о пользователе"
// @Param User-ID header string true "ID пользователя"
// @Param Authorization header string true "Auth token" "Токен авторизации"
// @Success 200     {object} HTTPResponse     "Информация успешно обновлена"
// @Failure 400     {object} HTTPResponse     "Неправильные входные данные"
// @Failure 401 {object} HTTPResponse "Пользователь не авторизован"
// @Failure 404     {object} HTTPResponse     "Пользователь не найден"
// @Router /set_user_info [put]
// @Security ApiKeyAuth
func SetUserInfo(resp http.ResponseWriter, req *http.Request) {

	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	var input SetUserInfoReq

	if err := readBody(req, &input); err != nil {
		resp.WriteHeader(http.StatusUnprocessableEntity)
		respBody.SetError(err)
		return
	}

	if !input.IsValid() {
		resp.WriteHeader(http.StatusBadRequest)
		respBody.SetError(errors.New("invalid input"))
		return
	}

	userID, _ := primitive.ObjectIDFromHex(req.Header.Get(HeaderUserID))

	if err := service.SetUserInfo(&domain.UserInfo{
		ID:   userID,
		Name: input.Name,
	}); err != nil {
		resp.WriteHeader(http.StatusNotFound)
		respBody.SetError(err)
		return
	}
}

// ChangePsw
// @Summary Изменить пароль
// @Description Изменяет пароль пользователя.
// @Tags users
// @Accept  json
// @Produce  json
// @Param   input  body      ChangePswReq  true  "Новый пароль пользователя"
// @Param User-ID header string true "ID пользователя"
// @Param Authorization header string true "Auth token" "Токен авторизации"
// @Success 200     {object} HTTPResponse   "Пароль успешно изменен"
// @Failure 400     {object} HTTPResponse   "Неправильные входные данные"
// @Failure 401 {object} HTTPResponse "Пользователь не авторизован"
// @Failure 404     {object} HTTPResponse   "Пользователь не найден"
// @Router /change_psw [put]
// @Security ApiKeyAuth
func ChangePsw(resp http.ResponseWriter, req *http.Request) {

	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	var input ChangePswReq

	if err := readBody(req, &input); err != nil {
		resp.WriteHeader(http.StatusUnprocessableEntity)
		respBody.SetError(err)
		return
	}

	if !input.IsValid() {
		resp.WriteHeader(http.StatusBadRequest)
		respBody.SetError(errors.New("invalid input"))
		return
	}

	userID, _ := primitive.ObjectIDFromHex(req.Header.Get(HeaderUserID))
	err := service.ChangePsw(&domain.UserPassword{
		ID:       userID,
		Password: input.Password,
	})
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		respBody.SetError(err)
		return
	}
}

// BindTelegramBot
// @Summary Привязать Telegram-бота
// @Description Привязывает Telegram username к пользователю.
// @Tags telegram
// @Accept  json
// @Produce  json
// @Param   input  body      BindTelegramData  true  "Данные Telegram"
// @Param User-ID header string true "ID пользователя"
// @Param Authorization header string true "Auth token" "Токен авторизации"
// @Success 200     {object} HTTPResponse      "Telegram username успешно привязан"
// @Failure 400     {object} HTTPResponse      "Неправильные входные данные"
// @Failure 401     {object} HTTPResponse      "Неавторизованный запрос"
// @Failure 500     {object} HTTPResponse      "Ошибка сервера"
// @Router /bind_telegram [post]
// @Security ApiKeyAuth
func BindTelegramBot(resp http.ResponseWriter, req *http.Request) {
	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	var input BindTelegramData
	if err := readBody(req, &input); err != nil {
		resp.WriteHeader(http.StatusUnprocessableEntity)
		respBody.SetError(err)
		return
	}

	userIDHex := req.Header.Get(HeaderUserID)
	if userIDHex == "" {
		resp.WriteHeader(http.StatusUnauthorized)
		respBody.SetError(errors.New("user ID is missing"))
		return
	}
	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		respBody.SetError(errors.New("invalid user ID"))
		return
	}

	err = service.BindTelegram(userID, input.TelegramUsername)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		respBody.SetError(err)
		return
	}

	respBody.SetData("Telegram username успешно привязан к пользователю")
}

// LoginWithTelegram
// @Summary Вход с помощью Telegram
// @Description Отправляет код авторизации в Telegram чат пользователя.
// @Tags telegram
// @Accept  json
// @Produce  json
// @Param   input  body      LoginWithTelegramData  true "Telegram username"
// @Success 200     {object} HTTPResponse        "Код отправлен в Telegram чат"
// @Failure 400     {object} HTTPResponse        "Неправильные входные данные"
// @Failure 422     {object} HTTPResponse        "Ошибка валидации"
// @Failure 500     {object} HTTPResponse        "Ошибка сервера"
// @Router /login_with_telegram [post]
func LoginWithTelegram(resp http.ResponseWriter, req *http.Request) {
	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	var input LoginWithTelegramData

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		respBody.SetError(err)
		return
	}

	if strings.TrimSpace(input.TelegramUsername) == "" {
		resp.WriteHeader(http.StatusBadRequest)
		respBody.SetError(errors.New("telegram username is required"))
		return
	}

	err := service.LoginWithTelegram(input.TelegramUsername)
	if err != nil {
		if errors.Is(err, appError.ErrChatIDEmpty) {
			resp.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			resp.WriteHeader(http.StatusInternalServerError)
		}
		respBody.SetError(err)
		return
	}

	respBody.SetData("Код отправлен в Telegram чат.")
}

// ConfirmTelegramCode
// @Summary Подтверждение кода Telegram
// @Description Подтверждает код авторизации из Telegram и возвращает токен.
// @Tags telegram
// @Accept  json
// @Produce  json
// @Param   input  body      ConfirmTelegramCodeData  true  "Telegram username и код"
// @Success 200     {object} HTTPResponse{data=string}  "Токен пользователя"
// @Failure 400     {object} HTTPResponse             "Неправильные входные данные"
// @Failure 401     {object} HTTPResponse             "Неверный код авторизации"
// @Router /confirm_telegram_code [post]
func ConfirmTelegramCode(resp http.ResponseWriter, req *http.Request) {
	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	var input ConfirmTelegramCodeData

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		respBody.SetError(err)
		return
	}

	if strings.TrimSpace(input.TelegramUsername) == "" || strings.TrimSpace(input.Code) == "" {
		resp.WriteHeader(http.StatusBadRequest)
		respBody.SetError(errors.New("telegram username and code are required"))
		return
	}

	userToken, err := service.ConfirmTelegramCode(input.TelegramUsername, input.Code)
	if err != nil {
		resp.WriteHeader(http.StatusUnauthorized)
		respBody.SetError(err)
		return
	}

	respBody.SetData(userToken)
}

func readBody(req *http.Request, s any) error {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, s)
}
