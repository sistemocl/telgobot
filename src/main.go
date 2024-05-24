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

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	"github.com/joho/godotenv"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	userLastCommand map[int]map[string]time.Time
	mu              sync.Mutex
)

func main() {
	userLastCommand = make(map[int]map[string]time.Time)

	godotenv.Load()
	token := os.Getenv("BOT_TOKEN")
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	admin_user := os.Getenv("ADMIN_USER")
	admin_pass := os.Getenv("ADMIN_PASS")

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

	commands := map[string]string{
		"/tiendas": "Abre el sitio web de Tiendas.",
		"/flr":     "Abre el sitio de FLR.",
		"/test":    "Abre una pagina de ejemplo",
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

		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		command := strings.ToLower(strings.ReplaceAll(m.Text, " ", ""))

		mu.Lock()
		if userLastCommand[int(m.Sender.ID)] == nil {
			userLastCommand[int(m.Sender.ID)] = make(map[string]time.Time)
		}
		lastExecTime := userLastCommand[int(m.Sender.ID)]["/tiendas"]
		mu.Unlock()

		if time.Since(lastExecTime).Seconds() < 15 {
			bot.Reply(m, "Debes esperar al menos 15 segundos entre ejecuciones de este comando.")
			return
		}
		if command == "/tiendas" {
			screenshot, err := Tiendas(ctx, "http://10.115.43.82:3002/login", admin_user, admin_pass)
			if err != nil {
				log.Printf("Error al tomar captura de pantalla: %v", err)
				return
			}
			Photos_response(screenshot, m, bot)
		} else {
			bot.Send(m.Chat, "Comando no reconocido. Por favor, intenta nuevamente.")
		}
		mu.Lock()
		userLastCommand[int(m.Sender.ID)]["/tiendas"] = time.Now()
		mu.Unlock()

		cancel()

	})

	bot.Handle("/flr", func(m *tb.Message) {

		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		command := strings.ToLower(strings.ReplaceAll(m.Text, " ", ""))

		mu.Lock()
		if userLastCommand[int(m.Sender.ID)] == nil {
			userLastCommand[int(m.Sender.ID)] = make(map[string]time.Time)
		}
		lastExecTime := userLastCommand[int(m.Sender.ID)]["/flr"]
		mu.Unlock()

		if time.Since(lastExecTime).Seconds() < 15 {
			bot.Reply(m, "Debes esperar al menos 15 segundos entre ejecuciones de este comando.")
			return
		}
		if command == "/flr" {
			screenshot, err := FLR(ctx, "http://10.115.43.118:3008/il/grafana/login", user, password)
			if err != nil {
				log.Printf("Error al tomar captura de pantalla: %v", err)
				return
			}
			Photos_response(screenshot, m, bot)
		} else {
			bot.Send(m.Chat, "Comando no reconocido. Por favor, intenta nuevamente.")
		}
		mu.Lock()
		userLastCommand[int(m.Sender.ID)]["/flr"] = time.Now()
		mu.Unlock()
		cancel()

	})

	bot.Handle("/test", func(m *tb.Message) {

		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		command := strings.ToLower(strings.ReplaceAll(m.Text, " ", ""))

		mu.Lock()
		if userLastCommand[int(m.Sender.ID)] == nil {
			userLastCommand[int(m.Sender.ID)] = make(map[string]time.Time)
		}
		lastExecTime := userLastCommand[int(m.Sender.ID)]["/test"]
		mu.Unlock()

		if time.Since(lastExecTime).Seconds() < 15 {
			bot.Reply(m, "Debes esperar al menos 15 segundos entre ejecuciones de este comando.")
			return
		}
		if command == "/test" {
			screenshot, err := Example(ctx, "https://chatgpt.com/")
			if err != nil {
				log.Printf("Error al tomar captura de pantalla: %v", err)
				return
			}
			Photos_response(screenshot, m, bot)
		} else {
			bot.Send(m.Chat, "Comando no reconocido. Por favor, intenta nuevamente.")
		}

		mu.Lock()
		userLastCommand[int(m.Sender.ID)]["/test"] = time.Now()
		mu.Unlock()

	})

	log.Println("Bot is running. Press CTRL+C to exit.")
	go bot.Start()

	select {}
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

	task := chromedp.Tasks{
		emulation.SetDeviceMetricsOverride(1920, 1080, 1, false),
		chromedp.Navigate(url),
		chromedp.WaitVisible("input[name=password]", chromedp.BySearch),
		chromedp.SendKeys("input[name=user]", user, chromedp.BySearch),
		chromedp.SendKeys("input[name=password]", password, chromedp.BySearch),
		chromedp.Click(`button[aria-label="Login button"]`, chromedp.BySearch),
		chromedp.WaitVisible("body", chromedp.BySearch),
		chromedp.Sleep(1 * time.Second),
		chromedp.Navigate("http://10.115.43.118:3008/il/grafana/?orgId=1"),
		chromedp.Navigate("http://10.115.43.118:3008/il/grafana/d/sDmADcSIk/mle-flr?orgId=1&refresh=30s"),
		chromedp.WaitVisible("body", chromedp.BySearch),
		chromedp.Sleep(3 * time.Second),
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
		emulation.SetDeviceMetricsOverride(1920, 1080, 1, false),
		chromedp.Navigate(url),
		chromedp.WaitVisible("input[name=password]", chromedp.BySearch),
		chromedp.SendKeys("input[name=user]", user, chromedp.BySearch),
		chromedp.SendKeys("input[name=password]", password, chromedp.BySearch),
		chromedp.Click(`button[aria-label="Login button"]`, chromedp.BySearch),
		chromedp.WaitVisible("body", chromedp.BySearch),
		chromedp.Sleep(1 * time.Second),
		chromedp.Navigate("http://10.115.43.82:3002/dashboards"),
		chromedp.WaitVisible("body", chromedp.BySearch),
		chromedp.Navigate("http://10.115.43.82:3002/d/yPbn4f2Sk/consolidado-v3-remoto?orgId=4&refresh=1m"),
		chromedp.WaitVisible("body", chromedp.BySearch),
		chromedp.Sleep(1 * time.Second),
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

//emulation.SetDeviceMetricsOverride(1920, 1080, 1, false),
