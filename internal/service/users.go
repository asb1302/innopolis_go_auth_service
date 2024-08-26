package service

import (
	"authservice/internal/config"
	"authservice/internal/domain"
	"authservice/internal/repository/tokendb"
	"authservice/internal/repository/userdb"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var users userdb.DB
var tokens tokendb.DB
var bot *tgbotapi.BotAPI

func Init(userDB userdb.DB, tokenDB tokendb.DB, telegramBot *tgbotapi.BotAPI) {
	users = userDB
	tokens = tokenDB
	bot = telegramBot
}

func SignUp(lp *domain.LoginPassword) (*domain.UserToken, error) {
	if _, ok := users.CheckExistLogin(lp.Login); ok {
		return nil, errors.New("login " + lp.Login + " already exists")
	}

	newUser := domain.User{
		ID:             primitive.NewObjectID(),
		Login:          lp.Login,
		Password:       hash(lp.Password),
		Role:           domain.UserRoleDefault,
		TelegramChatID: 0,
	}

	if err := users.SetUser(&newUser); err != nil {
		return nil, err
	}

	token := createToken(lp.Login)

	if err := tokens.SetUserToken(token, newUser.ID); err != nil {
		return nil, err
	}

	return &domain.UserToken{
		UserId: newUser.ID,
		Token:  token,
	}, nil
}

func SignIn(lp *domain.LoginPassword) (*domain.UserToken, error) {
	userId, ok := users.CheckExistLogin(lp.Login)
	if !ok {
		return nil, errors.New("user not found")
	}

	user, err := users.GetUser(*userId)
	if err != nil {
		return nil, err
	}

	if user.Password != hash(lp.Password) {
		return nil, errors.New("wrong password")
	}

	token := createToken(lp.Login)

	if err := tokens.SetUserToken(token, *userId); err != nil {
		return nil, err
	}

	return &domain.UserToken{
		UserId: *userId,
		Token:  token,
	}, nil
}

func SetUserInfo(ui *domain.UserInfo) error {

	user, err := users.GetUser(ui.ID)
	if err != nil {
		return err
	}

	user.Name = ui.Name

	return users.SetUser(user)
}

func ChangePsw(up *domain.UserPassword) error {
	user, err := users.GetUser(up.ID)
	if err != nil {
		return err
	}

	user.Password = hash(up.Password)

	return users.SetUser(user)
}

func GetUserShortInfo(id primitive.ObjectID) (*domain.UserInfo, error) {
	user, err := users.GetUser(id)
	if err != nil {
		return nil, err
	}

	ui := domain.UserInfo{
		ID:                 user.ID,
		Name:               user.Name,
		TelegramUsername:   user.TelegramUsername,
		TelegramUserChatID: user.TelegramChatID,
	}

	return &ui, nil
}

func GetUserFullInfo(id primitive.ObjectID) (*domain.User, error) {
	user, err := users.GetUser(id)
	return user, err
}

func GetUserIDByToken(token string) (*primitive.ObjectID, error) {
	return tokens.GetUserByToken(token)
}

func hash(str string) string {
	hp := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hp[:])
}

func createToken(login string) string {
	timeChs := md5.Sum([]byte(time.Now().String()))
	loginChs := md5.Sum([]byte(login))

	return hex.EncodeToString(timeChs[:]) + hex.EncodeToString(loginChs[:])
}

func BindTelegram(userID primitive.ObjectID, telegramUsername string) error {
	user, err := users.GetUser(userID)
	if err != nil {
		return err
	}

	user.TelegramUsername = telegramUsername

	return users.SetUser(user)
}

func BindTelegramChatID(userID primitive.ObjectID, chatID int64) error {
	user, err := users.GetUser(userID)
	if err != nil {
		return err
	}

	user.TelegramChatID = chatID

	return users.SetUser(user)
}

func LoginWithTelegram(telegramUsername string) error {
	user, err := users.GetUserByTelegramUsername(telegramUsername)
	if err != nil {
		return errors.New("пользователь с таким Telegram username не найден")
	}

	code := generateCode()

	err = SendCodeToTelegram(user.TelegramChatID, code)
	if err != nil {
		log.Printf("Ошибка отправки кода в Telegram для пользователя %s: %v", user.ID.Hex(), err)
		return err
	}

	// Сохраняем код и время его создания
	user.AuthCode = code
	user.AuthCodeTime = time.Now().Unix()
	err = users.SetUser(user)
	if err != nil {
		return err
	}

	log.Printf("Код для входа отправлен и сохранен для пользователя %s", user.ID.Hex())
	return nil
}

func ConfirmTelegramCode(telegramUsername, code string) (*domain.UserToken, error) {
	user, err := users.GetUserByTelegramUsername(telegramUsername)
	if err != nil {
		return nil, errors.New("пользователь с таким Telegram username не найден")
	}

	// Проверяем код и время его действия
	if user.AuthCode != code || time.Now().Unix()-user.AuthCodeTime > config.GetConfig().CodeExpiryDuration {
		return nil, errors.New("неверный или истекший код")
	}

	// Сбрасываем код после успешной проверки
	user.AuthCode = ""
	user.AuthCodeTime = 0
	if err := users.SetUser(user); err != nil {
		return nil, err
	}

	// Создаем токен для пользователя, как при обычном входе
	token := createToken(user.Login)
	if err := tokens.SetUserToken(token, user.ID); err != nil {
		return nil, err
	}

	return &domain.UserToken{
		UserId: user.ID,
		Token:  token,
	}, nil
}
