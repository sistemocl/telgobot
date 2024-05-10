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
	admin_user := os.Getenv("ADMIN_USER")
	admin_pass := os.Getenv("ADMIN_PASS")

	// Verificar si las variables de entorno están definidas
	if token == "" || user == "" || password == "" {
		fmt.Println("No se han definido las tres variables de entorno: TOKEN, USER y PASSWORD")
		return
	}

	//Utilizar las variables de entorno
	fmt.Println("TOKEN:", token)
	fmt.Println("USER:", user)
	fmt.Println("PASSWORD:", password)
	fmt.Println("ADMIN_USER:", admin_user)
	fmt.Println("Admin_pass:", admin_pass)

	// Crea un contexto padre para el navegador
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Crea dos contextos hijo para dos pestañas
	ctx1, cancel1 := chromedp.NewContext(ctx)
	defer cancel1()

	ctx2, cancel2 := chromedp.NewContext(ctx)
	defer cancel2()

	// ctx3, cancel3 := chromedp.NewContext(ctx)
	// defer cancel3()

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
		"/Tiendas":  "Abre el sitio web de Tiendas.",
		"/SBS":      "Abre el sitio web de SBS.",
		"/flr":      "Abre el sitio de FLR.",
		"/flr_info": "Retorna información delicada del sitio",
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
		go flr(ctx1, bot, m, "http://10.115.43.118:3008/il/grafana/login", user, password)
	})

	bot.Handle("/Tiendas", func(m *tb.Message) {
		go Tiendas(ctx2, bot, m, "http://10.115.43.82:3002/login", admin_user, admin_pass)
	})

	// bot.Handle("/SBS", func(m *tb.Message) {
	// 	go Tiendas(ctx3, bot, m, "https://www.youtube.com/")
	// })

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
		chromedp.WaitVisible("input[name=password]", chromedp.BySearch),
		chromedp.SendKeys("input[name=user]", user, chromedp.BySearch),
		chromedp.SendKeys("input[name=password]", password, chromedp.BySearch),
		chromedp.Click("button[aria-label=Login button]", chromedp.BySearch),
		chromedp.WaitVisible("body", chromedp.BySearch),
		chromedp.Navigate("http://10.115.43.118:3008/il/grafana/d/sDmADcSIk/mle-flr?orgId=1&refresh=30s"),
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

func Tiendas(ctx context.Context, bot *tb.Bot, m *tb.Message, url string, user string, password string) error {

	var buf []byte

	// Navega a la página
	task := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible("input[name=password]", chromedp.BySearch),
		chromedp.SendKeys("input[name=user]", user, chromedp.BySearch),
		chromedp.SendKeys("input[name=password]", password, chromedp.BySearch),
		chromedp.Click(`button[aria-label="Login button"]`, chromedp.BySearch),
		chromedp.WaitVisible("body", chromedp.BySearch),
		chromedp.Sleep(1 * time.Second),
		chromedp.Navigate("http://10.115.43.82:3002/dashboards"),
		chromedp.Navigate("http://10.115.43.82:3002/d/yPbn4f2Sk/consolidado-v3-remoto?orgId=4&refresh=1m"),
		chromedp.Sleep(2 * time.Second),
		chromedp.FullScreenshot(&buf, 90),
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

// func SBS(ctx context.Context, bot *tb.Bot, m *tb.Message, url string) error {

// 	var buf []byte

// 	// Navega a la página
// 	task := chromedp.Tasks{
// 		chromedp.Navigate(url),
// 		chromedp.CaptureScreenshot(&buf),
// 	}

// 	err := chromedp.Run(ctx, task)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Crear un lector para el buffer de bytes
// 	reader := bytes.NewReader(buf)

// 	// Enviar la imagen como un mensaje de foto
// 	photo := &tb.Photo{File: tb.FromReader(reader)}
// 	_, err = bot.Send(m.Chat, photo)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return nil
// }

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
