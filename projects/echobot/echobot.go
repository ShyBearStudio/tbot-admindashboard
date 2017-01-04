package echobot

import (
	"fmt"
	"sync"

	"github.com/ShyBearStudio/tbot-admindashboard/logger"
	"github.com/ShyBearStudio/tbot-admindashboard/projects/echobot/data"
	"github.com/mrd0ll4r/tbotapi"
)

type TBot interface {
	Start() error
	Stop() error
}

type EchoBot struct {
	operMutex *sync.Mutex
	isRunning bool
	stopped   chan struct{}
}

func New(configFileName string) (bot *EchoBot, err error) {
	bot = new(EchoBot)
	if err = initBot(configFileName); err != nil {
		return
	}

	bot.isRunning = false
	bot.operMutex = new(sync.Mutex)
	return
}

func (bot *EchoBot) Start() error {
	bot.operMutex.Lock()
	defer bot.operMutex.Unlock()
	if !bot.isRunning {
		bot.stopped = make(chan struct{})
		bot.isRunning = true
		go bot.RunEngine()
	}
	return nil
}

func (bot *EchoBot) Stop() error {
	_ = "breakpoint"
	fmt.Println("Stop called")
	bot.operMutex.Lock()
	defer bot.operMutex.Unlock()
	close(bot.stopped)
	return nil
}

func (bot *EchoBot) RunEngine() {
	_ = "breakpoint"
	api, err := tbotapi.New(config.Token)
	if err != nil {
		logger.Errorf("Cannot create tbot API: %v", err)
		return
	}
	defer api.Close()
	defer logger.Intoln("Echo bot has been stoped")
	for {
		select {
		case <-bot.stopped:
			_ = "breakpoint"
			return
		case botUpdate := <-api.Updates:
			bot.RunBody(api, botUpdate)
			fmt.Println("update processed")
		}
	}
}

func (bot *EchoBot) RunBody(api *tbotapi.TelegramBotAPI, botUpdate tbotapi.BotUpdate) {
	if err := botUpdate.Error(); err != nil {
		logger.Errorf("There was an error in update: %v", err)
		return
	}
	update := botUpdate.Update()
	if update.Type() == tbotapi.MessageUpdate {
		msg := update.Message
		if msg.Type() != tbotapi.TextMessage {
			return
		}

		is, err := data.IsRegisteredChat(msg.Chat.ID)
		if err != nil {
			logger.Errorf("Cannot check whether chat '%v' registered", msg.Chat)
		}
		if !is {
			data.CreateChat(msg.Chat)
		}

		// Now simply echo that back.
		_, err = api.NewOutgoingMessage(tbotapi.NewRecipientFromChat(msg.Chat), *msg.Text).Send()
		if err != nil {
			logger.Errorf("Error sending: %s\n", err)
			return
		}
	}
}

func initBot(congFileName string) error {
	_ = "breakpoint"
	if err := loadConfig(congFileName); err != nil {
		logger.Errorf("Cannot load configuration file with name '%s': %v", congFileName, err)
		return err
	}
	if err := logger.InitLog(config.Log.Dir); err != nil {
		logger.Errorf("Cannot initialize logger with log files in dir '%s': ", config.Log.Dir, err)
		return err
	}
	if err := data.InitDb(config.Db.Driver, config.Db.ConnectionString); err != nil {
		logger.Errorf("Cannot initialize dababase (Driver: '%s' | ConnString: '%s'): ", config.Db.Driver, config.Db.ConnectionString, err)
		return err
	}
	return nil
}
