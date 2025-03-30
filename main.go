package main

import (
	"context"
	"fmt"
	"github.com/mymmrac/telego"
	t_handler "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	gm "gopkg.in/gomail.v2"
	"log"
	"os"
	"strconv"
	"time"
)

type MyBotData struct {
	isWaitingForId         map[int64]bool
	bookLimitExtensionTime time.Duration
	botToken               string
	email                  string
	emailPassword          string
	emailSmptHost          string
	emailSmptPort          int
}

func (b *MyBotData) handleStart(ctx *t_handler.Context, msg telego.Message) error {
	b.isWaitingForId[msg.Chat.ChatID().ID] = false
	chatID := tu.ID(msg.Chat.ID)

	keyboard := tu.Keyboard(
		tu.KeyboardRow(tu.KeyboardButton("Записаться на экскурсию"), tu.KeyboardButton("Продлить книгу")),
		tu.KeyboardRow(tu.KeyboardButton("Афиша мероприятий"), tu.KeyboardButton("Литрес")),
		tu.KeyboardRow(tu.KeyboardButton("Пушкинская карта"), tu.KeyboardButton("Электронный каталог")),
		tu.KeyboardRow(tu.KeyboardButton("Комплектуемся вместе")),
	)
	greeting := fmt.Sprintf("Добро пожаловать %s %s, выберите действие", msg.From.FirstName, msg.From.LastName)
	message := tu.Message(
		chatID,
		greeting,
	).WithReplyMarkup(keyboard)

	_, _ = ctx.Bot().SendMessage(ctx, message)
	return nil
}

func (b *MyBotData) handleSignUpForExcursion(ctx *t_handler.Context, msg telego.Message) error {

	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(tu.InlineKeyboardButton("Оплатить экскурсию").WithURL("https://rnd.kassir.ru/frame/organizer/view/41139?key=b9c95356-77a1-9c01-9e67-e34eea7606d5&WIDGET_2754445811=eo6gr0n55vq0l4entpoeq47q66")),
	)

	message := tu.Message(
		msg.Chat.ChatID(),
		"Будем рады вас видеть!",
	).WithReplyMarkup(keyboard)
	_, _ = ctx.Bot().SendMessage(ctx, message)
	return nil
}

func (b *MyBotData) handleExtendBookRentLimit(ctx *t_handler.Context, msg telego.Message) error {
	b.isWaitingForId[msg.Chat.ChatID().ID] = true
	keyboard := tu.Keyboard(
		tu.KeyboardRow(tu.KeyboardButton("/start")),
	)

	message := tu.Message(
		msg.Chat.ChatID(),
		"Введите номер читательского или нажмите /start для отмены продления",
	).WithReplyMarkup(keyboard)
	_, _ = ctx.Bot().SendMessage(ctx, message)
	return nil
}

func (b *MyBotData) handleEventsDisplay(ctx *t_handler.Context, msg telego.Message) error {
	chatID := tu.ID(msg.Chat.ID)

	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(tu.InlineKeyboardButton("Афиша мероприятий").WithURL("https://vk.com/bibl_chehova?w=app6819359_-89514391")),
	)

	message := tu.Message(
		chatID,
		"Находитесь в курсе всего!",
	).WithReplyMarkup(keyboard)

	_, _ = ctx.Bot().SendMessage(ctx, message)
	return nil
}

func (b *MyBotData) handleLitres(ctx *t_handler.Context, msg telego.Message) error {
	chatID := tu.ID(msg.Chat.ID)

	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(tu.InlineKeyboardButton("Литрес").WithURL("http://taglib.ru/news/Biblioteki_Taganroga_predlagaut_novii_format_chteniya.html")),
	)

	message := tu.Message(
		chatID,
		"Получите доступ к электронным книгам на платформе ЛитРес!",
	).WithReplyMarkup(keyboard)

	_, _ = ctx.Bot().SendMessage(ctx, message)
	return nil
}

func (b *MyBotData) handlePushkinCard(ctx *t_handler.Context, msg telego.Message) error {
	chatID := tu.ID(msg.Chat.ID)

	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(tu.InlineKeyboardButton("Пушкинская карта").WithURL("https://rnd.kassir.ru/frame/organizer/view/41139?key=b9c95356-77a1-9c01-9e67-e34eea7606d5&WIDGET_2754445811=eo6gr0n55vq0l4entpoeq47q66")),
	)

	message := tu.Message(
		chatID,
		"Пушкинская карта",
	).WithReplyMarkup(keyboard)

	_, _ = ctx.Bot().SendMessage(ctx, message)
	return nil
}

func (b *MyBotData) handleSearchingBook(ctx *t_handler.Context, msg telego.Message) error {
	chatID := tu.ID(msg.Chat.ID)
	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Электронный каталог").WithURL("http://taglib.ru/bd.html"),
		),
	)

	message := tu.Message(
		chatID,
		"Просмотр каталогов по ссылке",
	).WithReplyMarkup(keyboard)

	_, _ = ctx.Bot().SendMessage(ctx, message)
	return nil
}

func (b *MyBotData) handleBundledTogether(ctx *t_handler.Context, msg telego.Message) error {
	chatID := tu.ID(msg.Chat.ID)

	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(tu.InlineKeyboardButton("Комплектуемся вместе").WithURL("http://cbs-tag.ru/index.php/component/content/article?id=1453")),
	)

	message := tu.Message(
		chatID,
		"Комплектуемся вместе",
	).WithReplyMarkup(keyboard)

	_, _ = ctx.Bot().SendMessage(ctx, message)
	return nil
}

func (b *MyBotData) handleInput(ctx *t_handler.Context, msg telego.Message) error {
	if b.isWaitingForId[msg.Chat.ChatID().ID] {
		libraryCardNumber, err := strconv.Atoi(msg.Text)
		if err != nil {
			message := tu.Message(
				msg.Chat.ChatID(),
				"Введите число",
			)
			_, _ = ctx.Bot().SendMessage(ctx, message)
			return err
		}
		b.isWaitingForId[msg.Chat.ChatID().ID] = false

		err = b.sendEmail(libraryCardNumber)
		if err == nil {
			currentTime := time.Now()
			currentTime = currentTime.Add(b.bookLimitExtensionTime)
			message := tu.Message(
				msg.Chat.ChatID(),
				fmt.Sprintf("Вы успешно продлили книгу до %s", currentTime.Format("02-01-2006")),
			)
			_, _ = ctx.Bot().SendMessage(ctx, message)
		} else {
			message := tu.Message(
				msg.Chat.ChatID(),
				"Ошибка в отправке заявки, уведомление в тех.поддержку уже отправлено.",
			)
			_, _ = ctx.Bot().SendMessage(ctx, message)
		}
	}
	return nil
}

func (b *MyBotData) sendEmail(libraryCardNumber int) error {
	m := gm.NewMessage()

	m.SetHeader("From", b.email)
	m.SetHeader("To", b.email)
	emailSubjectText := fmt.Sprintf("От бота о продлении книги\n %d хочет продлить свои книги", libraryCardNumber)
	m.SetHeader("Subject", emailSubjectText)

	d := gm.NewDialer(b.emailSmptHost, b.emailSmptPort, b.email, b.emailPassword)

	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		fmt.Println("Свяжитесь с поддержкой!")
		return err
	}
	fmt.Println("Email sent successfully.")
	return nil
}

var (
	TG_BOT_TOKEN           string
	TG_BOT_EMAIL           string
	TG_BOT_EMAIL_PASSWORD  string
	TG_BOT_EMAIL_SMTP_HOST string
	TG_BOT_EMAIL_SMTP_PORT string
)

func setBotData() *MyBotData {
	botToken := os.Getenv("TG_BOT_TOKEN")
	if len(botToken) == 0 {
		botToken = TG_BOT_TOKEN
	}

	email := os.Getenv("TG_BOT_EMAIL")
	if len(email) == 0 {
		email = TG_BOT_EMAIL
	}

	emailPassword := os.Getenv("TG_BOT_EMAIL_PASSWORD")
	if len(emailPassword) == 0 {
		emailPassword = TG_BOT_EMAIL_PASSWORD
	}
	emailSmptHost := os.Getenv("TG_BOT_EMAIL_SMTP_HOST")
	if len(emailSmptHost) == 0 {
		emailSmptHost = TG_BOT_EMAIL_SMTP_HOST
	}

	emailSmptPortRaw := os.Getenv("TG_BOT_EMAIL_SMTP_PORT")
	if len(emailSmptPortRaw) == 0 {
		emailSmptPortRaw = TG_BOT_EMAIL_SMTP_PORT
	}

	emailSmptPort, err := strconv.Atoi(emailSmptPortRaw)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	myBotData := MyBotData{isWaitingForId: make(map[int64]bool),
		bookLimitExtensionTime: time.Hour * 360, /* 15 дней */
		botToken:               botToken,
		email:                  email,
		emailPassword:          emailPassword,
		emailSmptHost:          emailSmptHost,
		emailSmptPort:          emailSmptPort,
	}

	return &myBotData
}

/// TODO: Перенести все ссылки в env?
/// TODO: Подогнать полностью под интерфейс мити
/// TODO: Защита от DDOS

func main() {
	ctx := context.Background()

	myBotData := setBotData()

	bot, err := telego.NewBot(myBotData.botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		fmt.Println("Свяжитесь с поддержкой")
		os.Exit(1)
	}

	updates, _ := bot.UpdatesViaLongPolling(ctx, nil)
	bot_handler, _ := t_handler.NewBotHandler(bot, updates)

	defer bot_handler.Stop()
	defer bot.StopPoll(ctx, nil)

	bot_handler.HandleMessage(myBotData.handleStart, t_handler.CommandEqual("start"))

	bot_handler.HandleMessage(myBotData.handleSignUpForExcursion, t_handler.TextEqual("Записаться на экскурсию"))

	bot_handler.HandleMessage(myBotData.handleExtendBookRentLimit, t_handler.TextEqual("Продлить книгу"))

	bot_handler.HandleMessage(myBotData.handleEventsDisplay, t_handler.TextEqual("Афиша мероприятий"))

	bot_handler.HandleMessage(myBotData.handleLitres, t_handler.TextEqual("Литрес"))

	bot_handler.HandleMessage(myBotData.handlePushkinCard, t_handler.TextEqual("Пушкинская карта"))

	bot_handler.HandleMessage(myBotData.handleSearchingBook, t_handler.TextEqual("Электронный каталог"))

	bot_handler.HandleMessage(myBotData.handleBundledTogether, t_handler.TextEqual("Комплектуемся вместе"))

	bot_handler.HandleMessage(myBotData.handleInput, t_handler.AnyMessage())

	bot_handler.Start()
}
