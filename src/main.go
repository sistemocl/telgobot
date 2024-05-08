package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/joho/godotenv"
	"github.com/tucnak/telebot"
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

	bot, err := telebot.NewBot(telebot.Settings{
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		Token:  token,
	})
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/pagina1", func(m *telebot.Message) {
		go handlePage(bot, m, "https://www.chess.com/login")
	})

	bot.Handle("/pagina2", func(m *telebot.Message) {
		go handlePage(bot, m, "https://pkg.go.dev/gopkg.in/telebot.v3@v3.2.1")
	})

	log.Println("Bot is running. Press CTRL+C to exit.")
	bot.Start()
}

func handlePage(bot *telebot.Bot, m *telebot.Message, url string) {
	// Crea un contexto padre para el navegador
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Crea dos contextos hijo para dos pestañas
	ctx1, cancel1 := chromedp.NewContext(ctx)
	defer cancel1()

	ctx2, cancel2 := chromedp.NewContext(ctx)
	defer cancel2()

	// Utiliza WaitGroup para sincronizar las goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	// Ejecuta la primera goroutine para la primera página
	go func() {
		defer wg.Done()
		if err := navigateAndSendMessage(ctx1, bot, m, url); err != nil {
			log.Println("Error en la goroutine 1:", err)
		}
	}()

	// Ejecuta la segunda goroutine para la segunda página
	go func() {
		defer wg.Done()
		if err := Dafiti(ctx2, bot, m, url); err != nil {
			log.Println("Error en la goroutine 2:", err)
		}
	}()

	// Espera a que ambas goroutines finalicen
	wg.Wait()
}

func navigateAndSendMessage(ctx context.Context, bot *telebot.Bot, m *telebot.Message, url string) error {
	// Crea un navegador Chrome
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	// Navega a la página
	if err := chromedp.Run(ctx, chromedp.Navigate(url)); err != nil {
		return err
	}

	// Envía un mensaje de respuesta
	msg := "Se ha navegado a la página: " + url
	if _, err := bot.Send(m.Chat, msg); err != nil {
		return err
	}

	return nil
}

func Dafiti(ctx context.Context, bot *telebot.Bot, m *telebot.Message, url string) error {
	// Crea un navegador Chrome
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	// Navega a la página
	if err := chromedp.Run(ctx, chromedp.Navigate(url)); err != nil {
		return err
	}

	// Envía un mensaje de respuesta
	msg := "@" + m.Sender.Username
	if _, err := bot.Send(m.Chat, msg); err != nil {
		return err
	}

	return nil
}
