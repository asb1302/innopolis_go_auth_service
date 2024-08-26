package cache

import (
	"authservice/internal/domain"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"sync"
)

type UserCache struct {
	userPull     map[primitive.ObjectID]*domain.User
	loginPull    map[string]primitive.ObjectID
	telegramPull map[string]primitive.ObjectID
	mtx          sync.RWMutex
}

const userDumpFileName = "users.json"

func UserCacheInit(ctx context.Context, wg *sync.WaitGroup) (*UserCache, error) {
	var c UserCache
	c.userPull = make(map[primitive.ObjectID]*domain.User)
	c.loginPull = make(map[string]primitive.ObjectID)
	c.telegramPull = make(map[string]primitive.ObjectID)

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		makeDump(userDumpFileName, c.userPull)
	}()

	if err := loadFromDump(userDumpFileName, &c.userPull); err != nil {
		return nil, err
	}

	for _, user := range c.userPull {
		c.loginPull[user.Login] = user.ID
		if user.TelegramUsername != "" {
			c.telegramPull[user.TelegramUsername] = user.ID
		}
	}

	return &c, nil
}

func (c *UserCache) CheckExistLogin(login string) (*primitive.ObjectID, bool) {
	c.mtx.RLock()
	id, ok := c.loginPull[login]
	c.mtx.RUnlock()
	return &id, ok
}

func (c *UserCache) GetUser(id primitive.ObjectID) (*domain.User, error) {
	c.mtx.RLock()
	user, ok := c.userPull[id]
	c.mtx.RUnlock()

	if !ok {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (c *UserCache) SetUser(newUserInfo *domain.User) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	// Удаляем старое значение telegramPull, если оно существует
	for username, id := range c.telegramPull {
		if id == newUserInfo.ID {
			delete(c.telegramPull, username)
		}
	}

	c.userPull[newUserInfo.ID] = newUserInfo
	c.loginPull[newUserInfo.Login] = newUserInfo.ID

	if newUserInfo.TelegramUsername != "" {
		log.Printf("Сохраняем Telegram username: %s для userID: %s", newUserInfo.TelegramUsername, newUserInfo.ID.Hex())
		c.telegramPull[newUserInfo.TelegramUsername] = newUserInfo.ID
	}

	log.Printf("Текущие данные пользователя: %+v", newUserInfo)
	return nil
}

func (c *UserCache) GetUserByTelegramUsername(username string) (*domain.User, error) {
	c.mtx.RLock()
	id, ok := c.telegramPull[username]
	c.mtx.RUnlock()

	if !ok {
		return nil, errors.New("user not found")
	}

	return c.GetUser(id)
}
