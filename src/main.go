package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
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
	fmt.Println("TOKEN:", token)
	fmt.Println("USER:", user)
	fmt.Println("PASSWORD:", password)

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

	bot, err := tb.NewBot(tb.Settings{
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
		Token:  token,
	})
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/pagina1", func(m *tb.Message) {
		go func() {
			if err := chromedp.Run(ctx1,
				chromedp.Navigate("https://www.chess.com/login"),
			); err != nil {
				log.Println(err)
				return
			}

			bot.Send(m.Sender, "Se navegó a la página 1.")
		}()
	})

	bot.Handle("/Dafiti", func(m *tb.Message) {
		go Dafiti(ctx2, bot, m, "https://pkg.go.dev/gopkg.in/telebot.v3@v3.2.1")
	})

	bot.Handle("/SBS", func(m *tb.Message) {
		go Dafiti(ctx3, bot, m, "https://www.youtube.com/")
	})

	log.Println("Bot is running. Press CTRL+C to exit.")
	bot.Start()
}

func Dafiti(ctx context.Context, bot *tb.Bot, m *tb.Message, url string) error {

	// Crea un navegador Chrome
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

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
	// Crea un navegador Chrome
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

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
