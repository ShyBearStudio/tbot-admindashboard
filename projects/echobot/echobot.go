package echobot

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/ShyBearStudio/tbot-admindashboard/configutils"
	"github.com/ShyBearStudio/tbot-admindashboard/logger"
	"github.com/ShyBearStudio/tbot-admindashboard/projects/data"
	"github.com/mrd0ll4r/tbotapi"
)

type EchoBot struct {
	operMutex *sync.Mutex
	isRunning bool
	stopped   chan struct{}
	db        *sql.DB
}

func New(configFileName string) (bot *EchoBot, err error) {
	bot = new(EchoBot)
	if err = configutils.LoadConfig(configFileName, configEnvVarName, &config); err != nil {
		logger.Errorf("Cannot load configuration file with name '%s': %v", configFileName, err)
		return
	}
	if err = logger.InitLogger(config.Log.Dir); err != nil {
		logger.Errorf("Cannot initialize logger with log files in dir '%s': ", config.Log.Dir, err)
		return
	}
	if bot.db, err = data.InitDb(config.Db.Driver, config.Db.ConnectionString); err != nil {
		logger.Errorf("Cannot initialize dababase (Driver: '%s' | ConnString: '%s'): ", config.Db.Driver, config.Db.ConnectionString, err)
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

		is, err := isRegisteredChat(bot.db, msg.Chat.ID)
		if err != nil {
			logger.Errorf("Cannot check whether chat '%v' registered", msg.Chat)
		}
		if !is {
			createChat(bot.db, msg.Chat)
		}

		// Now simply echo that back.
		_, err = api.NewOutgoingMessage(tbotapi.NewRecipientFromChat(msg.Chat), *msg.Text).Send()
		if err != nil {
			logger.Errorf("Error sending: %s\n", err)
			return
		}
	}
}

func (bot *EchoBot) Chats() (chat []Chat, err error) {
	return chats(bot.db)
}
