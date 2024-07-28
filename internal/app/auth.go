package app

import (
	"authservice/internal/config"
	"authservice/internal/handler/httphandler"
	"authservice/internal/repository/cache"
	"authservice/internal/server"
	"authservice/internal/service"
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Run() {
	// Инициализация конфигурации
	config.InitConfig()
	cfg := config.GetConfig()

	// Инициализация Telegram-бота
	bot, err := initTelegramBot(cfg.TelegramBotToken, cfg.ServerURL)
	if err != nil {
		log.Fatalf("Ошибка инициализации Telegram бота: %v", err)
	}
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup

	// Инициализация баз данных
	userDB, err := cache.UserCacheInit(ctx, &wg)
	if err != nil {
		log.Fatalf("ERROR failed to initialize user database: %v", err)
	}
	tokenDB, err := cache.TokenCacheInit(ctx, &wg)
	if err != nil {
		log.Fatalf("ERROR failed to initialize tokens database: %v", err)
	}

	// initialize service
	service.Init(userDB, tokenDB, bot)

	go func() {
		// изменено на 0.0.0.0, чтобы приложение было доступно извне
		err := server.Run("0.0.0.0", "8000", httphandler.NewRouter())
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("ERROR server run ", err)
		}
	}()

	log.Println("INFO auth service is running")

	<-ctx.Done()

	err = server.Shutdown()
	if err != nil {
		log.Fatal("ERROR server was not gracefully shutdown", err)
	}
	wg.Wait()

	log.Println("INFO auth service was gracefully shutdown")
}

func initTelegramBot(token, serverURL string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	go func() {
		for update := range updates {
			if update.Message != nil {
				service.HandleTelegramMessage(bot, update.Message, serverURL)
			}
		}
	}()

	return bot, nil
}
