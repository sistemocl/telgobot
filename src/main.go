package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/joho/godotenv"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	lastCommandTime map[int]time.Time
	lastCommandLock sync.Mutex
)

func main() {
	lastCommandTime = make(map[int]time.Time)

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

	bot, err := tb.NewBot(tb.Settings{
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
		Token:  token,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Crear un contexto padre
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	commands := map[string]string{
		"/tiendas": "Abre el sitio web de Tiendas.",
		// "/SBS":      "Abre el sitio web de SBS.",
		//"/Dafiti": "Abre el sitio web de Dafiti.",
		"/flr": "Abre el sitio de FLR.",
	}

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

	bot.Handle("/tiendas", func(m *tb.Message) {
		command := strings.ToLower(strings.ReplaceAll(m.Text, " ", ""))
		if command == "/tiendas" {
			if canExecuteCommand(int(m.Sender.ID)) {
				screenshot, err := Tiendas(ctx, "http://10.115.43.82:3002/login", admin_user, admin_pass)
				if err != nil {
					log.Printf("Error al tomar captura de pantalla: %v", err)
					return
				}
				Photos_response(screenshot, m, bot)
			} else {
				bot.Send(m.Sender, "Por favor espera al menos 15 segundos antes de volver a ejecutar el comando.")
			}
		} else {
			bot.Send(m.Chat, "Comando no reconocido. Por favor, intenta nuevamente.")
		}

	})

	bot.Handle("/test", func(m *tb.Message) {
		command := strings.ToLower(strings.ReplaceAll(m.Text, " ", ""))
		if command == "/test" {
			if canExecuteCommand(int(m.Sender.ID)) {
				screenshot, err := Example(ctx, "https://chatgpt.com/")
				if err != nil {
					log.Printf("Error al tomar captura de pantalla: %v", err)
					return
				}
				Photos_response(screenshot, m, bot)
			} else {
				bot.Send(m.Sender, "Por favor espera al menos 15 segundos antes de volver a ejecutar el comando.")
			}
		} else {
			bot.Send(m.Chat, "Comando no reconocido. Por favor, intenta nuevamente.")
		}
	})

	bot.Handle("/test2", func(m *tb.Message) {
		command := strings.ToLower(strings.ReplaceAll(m.Text, " ", ""))
		if command == "/test2" {
			if canExecuteCommand(int(m.Sender.ID)) {
				screenshot, err := Example(ctx, "https://www.perplexity.ai/")
				if err != nil {
					log.Printf("Error al tomar captura de pantalla: %v", err)
					return
				}
				Photos_response(screenshot, m, bot)
			} else {
				bot.Send(m.Sender, "Por favor espera al menos 15 segundos antes de volver a ejecutar el comando.")
			}
		} else {
			bot.Send(m.Chat, "Comando no reconocido. Por favor, intenta nuevamente.")
		}
	})

	go bot.Start()

	select {}
}

func canExecuteCommand(userID int) bool {
	lastCommandLock.Lock()
	defer lastCommandLock.Unlock()

	lastTime, ok := lastCommandTime[userID]
	if !ok {
		lastCommandTime[userID] = time.Now()
		return true
	}

	elapsed := time.Since(lastTime)
	if elapsed >= 20*time.Second {
		lastCommandTime[userID] = time.Now()
		return true
	}
	return false
}

func Photos_response(screenshot []byte, m *tb.Message, bot *tb.Bot) {
	photo := &tb.Photo{
		File: tb.FromReader(bytes.NewReader(screenshot)),
	}
	mention := ""
	if m.Sender.Username != "" {
		mention = "@" + m.Sender.Username
	} else {
		mention = m.Sender.FirstName
	}
	bot.Send(m.Chat, mention+" Aqui tienes:")
	bot.SendAlbum(m.Chat, tb.Album{photo})
}

func FLR(ctx context.Context, url, user, password string) ([]byte, error) {
	var buf []byte

	// Navega a la página
	task := chromedp.Tasks{
		chromedp.Navigate("http://10.115.43.118:3008/il/grafana/login"),
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

	return buf, nil
}

func Tiendas(ctx context.Context, url, user, password string) ([]byte, error) {
	var buf []byte

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

	return buf, nil
}

func Example(ctx context.Context, url string) ([]byte, error) {
	var buf []byte

	task := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Sleep(2 * time.Second),
		chromedp.FullScreenshot(&buf, 90),
	}

	err := chromedp.Run(ctx, task)
	if err != nil {
		log.Fatal(err)
	}

	return buf, nil
}
