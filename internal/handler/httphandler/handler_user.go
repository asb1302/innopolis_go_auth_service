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

func LoginWithTelegram(resp http.ResponseWriter, req *http.Request) {
	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	var input struct {
		TelegramUsername string `json:"telegram_username"`
	}
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

func ConfirmTelegramCode(resp http.ResponseWriter, req *http.Request) {
	respBody := &HTTPResponse{}
	defer func() {
		resp.Write(respBody.Marshall())
	}()

	var input struct {
		TelegramUsername string `json:"telegram_username"`
		Code             string `json:"code"`
	}
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
