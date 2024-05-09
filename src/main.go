package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/joho/godotenv"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {

	// Cargar variables de entorno
	godotenv.Load()
	token := os.Getenv("BOT_TOKEN")
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")

	// Verificar si las variables de entorno están definidas
	if token == "" || user == "" || password == "" {
		fmt.Println("No se han definido las tres variables de entorno: TOKEN, USER y PASSWORD")
		return
	}

	// Utilizar las variables de entorno
	// fmt.Println("TOKEN:", token)
	// fmt.Println("USER:", user)
	// fmt.Println("PASSWORD:", password)

	// Crea un contexto padre para el navegador
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Crea dos contextos hijo para dos pestañas
	ctx1, cancel1 := chromedp.NewContext(ctx)
	defer cancel1()

	ctx2, cancel2 := chromedp.NewContext(ctx)
	defer cancel2()

	ctx3, cancel3 := chromedp.NewContext(ctx)
	defer cancel3()

	// ctx4, cancel4 := chromedp.NewContext(ctx)
	// defer cancel4()

	bot, err := tb.NewBot(tb.Settings{
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
		Token:  token,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Mapa para almacenar los comandos y sus descripciones
	commands := map[string]string{
		"/Dafiti":  "Abre el sitio web de Dafiti.",
		"/SBS":     "Abre el sitio web de SBS.",
		"/flr":     "Abre el sitio de FLR.",
		"flr_info": "Retorna información delicada del sitio",
	}

	// Manejador de comando para mostrar todos los comandos disponibles
	bot.Handle("/comandos", func(m *tb.Message) {
		var reply strings.Builder
		reply.WriteString("Comandos disponibles:\n")
		for cmd, desc := range commands {
			reply.WriteString(cmd)
			reply.WriteString(" - ")
			reply.WriteString(desc)
			reply.WriteString("\n")
		}
		bot.Send(m.Chat, reply.String())
	})

	bot.Handle("/flr", func(m *tb.Message) {
		go flr(ctx1, bot, m, "https://www.chess.com/login", user, password)
	})

	bot.Handle("/Dafiti", func(m *tb.Message) {
		go Dafiti(ctx2, bot, m, "https://pkg.go.dev/gopkg.in/telebot.v3@v3.2.1")
	})

	bot.Handle("/SBS", func(m *tb.Message) {
		go Dafiti(ctx3, bot, m, "https://www.youtube.com/")
	})

	// bot.Handle("/flr_info", func(m *tb.Message) {
	// 	go flr_info(ctx4, bot, m, "https://www.chess.com/login", user, password)
	// })

	log.Println("Chat Bot is running. Press CTRL+C to exit.")
	bot.Start()
}

func flr(ctx context.Context, bot *tb.Bot, m *tb.Message, url string, user string, password string) error {

	var buf []byte

	// Navega a la página
	task := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible("input[type=password]", chromedp.ByQuery),
		chromedp.SendKeys("input[name=_username]", user, chromedp.ByQuery),
		chromedp.SendKeys("input[name=_password]", password, chromedp.ByQuery),
		chromedp.Click("button[type=submit]", chromedp.ByQuery),
		chromedp.Sleep(2 * time.Second),
		chromedp.FullScreenshot(&buf, 90),
	}

	err := chromedp.Run(ctx, task)
	if err != nil {
		log.Fatal(err)
	}

	// Crear un lector para el buffer de bytes
	reader := bytes.NewReader(buf)

	//msg := "@" + m.Sender.Username

	// Enviar la imagen como un mensaje de foto
	photo := &tb.Photo{File: tb.FromReader(reader)}
	_, err = bot.Send(m.Chat, photo)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func Dafiti(ctx context.Context, bot *tb.Bot, m *tb.Message, url string) error {

	var buf []byte

	// Navega a la página
	task := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.CaptureScreenshot(&buf),
	}

	err := chromedp.Run(ctx, task)
	if err != nil {
		log.Fatal(err)
	}

	// Crear un lector para el buffer de bytes
	reader := bytes.NewReader(buf)

	// Enviar la imagen como un mensaje de foto
	photo := &tb.Photo{File: tb.FromReader(reader)}
	_, err = bot.Send(m.Chat, photo)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func SBS(ctx context.Context, bot *tb.Bot, m *tb.Message, url string) error {

	var buf []byte

	// Navega a la página
	task := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.CaptureScreenshot(&buf),
	}

	err := chromedp.Run(ctx, task)
	if err != nil {
		log.Fatal(err)
	}

	// Crear un lector para el buffer de bytes
	reader := bytes.NewReader(buf)

	// Enviar la imagen como un mensaje de foto
	photo := &tb.Photo{File: tb.FromReader(reader)}
	_, err = bot.Send(m.Chat, photo)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// func flr_info(ctx context.Context, bot *tb.Bot, m *tb.Message, url string, user string, password string) error {

// 	// Navega a la página
// 	task := chromedp.Tasks{
// 		chromedp.Navigate(url),
// 		chromedp.WaitVisible("input[type=password]", chromedp.ByQuery),
// 		chromedp.SendKeys("input[name=_username]", user, chromedp.ByQuery),
// 		chromedp.SendKeys("input[name=_password]", password, chromedp.ByQuery),
// 		chromedp.Click("button[type=submit]", chromedp.ByQuery),
// 		chromedp.Sleep(2 * time.Second),
// 	}

// 	err := chromedp.Run(ctx, task)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	msg := "@" + m.Sender.Username

// 	bot.Send(m.Chat, msg)

// 	return nil
// }
