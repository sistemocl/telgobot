package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"tel_gobot/scrapper"
	"tel_gobot/sshclient"
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
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error al cargar el archivo .env")
		return
	}
	godotenv.Load()
	token := os.Getenv("BOT_TOKEN")
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	admin_user := os.Getenv("ADMIN_USER")
	admin_pass := os.Getenv("ADMIN_PASS")
	admin2 := os.Getenv("ADMIN2")
	pass2 := os.Getenv("PASS2")
	gor := os.Getenv("GOR")
	ip := os.Getenv("IP")
	pass3 := os.Getenv("SBSPASS")

	if token == "" || user == "" || password == "" {
		fmt.Println("No se han definido las tres variables de entorno: TOKEN, USER o PASSWORD")
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
		"/tiendas":    "Consolidado de tiendas",
		"/mle":        "Grafana MLE FLR",
		"/sbs":        " Reporte de cierre SBS",
		"/sbs2":       "SBS General",
		"/soporte":    "Grafana Soporte FLR",
		"/kpi":        "Kpi SBS",
		"/servidores": "Grafana FLR servidores",
		"/s11":        "panel S11",
		"/operacion":  "FLR operacion",
		"/zone":       "zone_orders SBS",
		"/colas":      "colas pps",
		"/erlang":     "erlang",
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

	bot.Handle("/erlang", func(m *tb.Message) {
		// // Dividir el mensaje para extraer el comando a ejecutar
		// parts := strings.Split(m.Text, " ")
		// if len(parts) < 2 {
		// 	bot.Send(m.Chat, "Debes especificar un comando para Erlang. Ejemplo: /ejemplo status")
		// 	return
		// }
		// command := parts[1]

		// Conectar a la consola de Erlang y ejecutar el comando
		output, err := sshclient.ConexionErlang(gor, pass3, ip, "application:get_env(butler_server,streaming_pps_list).")
		if err != nil {
			bot.Send(m.Chat, fmt.Sprintf("Error conectando a Erlang: %v", err))
			return
		}

		// Enviar la salida de Erlang como respuesta de texto
		bot.Send(m.Chat, fmt.Sprintf("Salida de Erlang: %s", output))
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
			bot.Send(m.Chat, " Tomando captura")
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

	bot.Handle("/mle", func(m *tb.Message) {

		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		command := strings.ToLower(strings.ReplaceAll(m.Text, " ", ""))

		mu.Lock()
		if userLastCommand[int(m.Sender.ID)] == nil {
			userLastCommand[int(m.Sender.ID)] = make(map[string]time.Time)
		}
		lastExecTime := userLastCommand[int(m.Sender.ID)]["/mle"]
		mu.Unlock()

		if time.Since(lastExecTime).Seconds() < 15 {
			bot.Reply(m, "Debes esperar al menos 15 segundos entre ejecuciones de este comando.")
			return
		}
		if command == "/mle" {
			bot.Send(m.Chat, " Tomando captura")
			screenshot, err := MLE(ctx, "http://10.115.43.118:3008/il/grafana/login", user, password)
			if err != nil {
				log.Printf("Error al tomar captura de pantalla: %v", err)
				return
			}
			Photos_response(screenshot, m, bot)
		} else {
			bot.Send(m.Chat, "Comando no reconocido. Por favor, intenta nuevamente.")
		}
		mu.Lock()
		userLastCommand[int(m.Sender.ID)]["/mle"] = time.Now()
		mu.Unlock()
		cancel()

	})

	bot.Handle("/soporte", func(m *tb.Message) {

		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		command := strings.ToLower(strings.ReplaceAll(m.Text, " ", ""))

		mu.Lock()
		if userLastCommand[int(m.Sender.ID)] == nil {
			userLastCommand[int(m.Sender.ID)] = make(map[string]time.Time)
		}
		lastExecTime := userLastCommand[int(m.Sender.ID)]["/soporte"]
		mu.Unlock()

		if time.Since(lastExecTime).Seconds() < 15 {
			bot.Reply(m, "Debes esperar al menos 15 segundos entre ejecuciones de este comando.")
			return
		}
		if command == "/soporte" {
			bot.Send(m.Chat, " Tomando captura")
			screenshot, err := Soporte_FLR(ctx, "http://10.115.43.118:3008/il/grafana/login", user, password)
			if err != nil {
				log.Printf("Error al tomar captura de pantalla: %v", err)
				return
			}
			Photos_response(screenshot, m, bot)
		} else {
			bot.Send(m.Chat, "Comando no reconocido. Por favor, intenta nuevamente.")
		}
		mu.Lock()
		userLastCommand[int(m.Sender.ID)]["/soporte"] = time.Now()
		mu.Unlock()
		cancel()

	})

	bot.Handle("/sbs", func(m *tb.Message) {

		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		command := strings.ToLower(strings.ReplaceAll(m.Text, " ", ""))

		mu.Lock()
		if userLastCommand[int(m.Sender.ID)] == nil {
			userLastCommand[int(m.Sender.ID)] = make(map[string]time.Time)
		}
		lastExecTime := userLastCommand[int(m.Sender.ID)]["/sbs"]
		mu.Unlock()

		if time.Since(lastExecTime).Seconds() < 15 {
			bot.Reply(m, "Debes esperar al menos 15 segundos entre ejecuciones de este comando.")
			return
		}
		if command == "/sbs" {
			bot.Send(m.Chat, " Tomando captura")
			screenshot, err := SBS(ctx, "http://10.115.43.24:3000/login", admin2, pass2)
			if err != nil {
				log.Printf("Error al tomar captura de pantalla: %v", err)
				return
			}
			Photos_response(screenshot, m, bot)
		} else {
			bot.Send(m.Chat, "Comando no reconocido. Por favor, intenta nuevamente.")
		}
		mu.Lock()
		userLastCommand[int(m.Sender.ID)]["/sbs"] = time.Now()
		mu.Unlock()
		cancel()

	})

	bot.Handle("/sbs2", func(m *tb.Message) {

		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		command := strings.ToLower(strings.ReplaceAll(m.Text, " ", ""))

		mu.Lock()
		if userLastCommand[int(m.Sender.ID)] == nil {
			userLastCommand[int(m.Sender.ID)] = make(map[string]time.Time)
		}
		lastExecTime := userLastCommand[int(m.Sender.ID)]["/sbs2"]
		mu.Unlock()

		if time.Since(lastExecTime).Seconds() < 15 {
			bot.Reply(m, "Debes esperar al menos 15 segundos entre ejecuciones de este comando.")
			return
		}
		if command == "/sbs2" {
			bot.Send(m.Chat, " Tomando captura")
			screenshot, err := SBS_General(ctx, "http://10.115.43.24:3000/login", admin2, pass2)
			if err != nil {
				log.Printf("Error al tomar captura de pantalla: %v", err)
				return
			}
			Photos_response(screenshot, m, bot)
		} else {
			bot.Send(m.Chat, "Comando no reconocido. Por favor, intenta nuevamente.")
		}
		mu.Lock()
		userLastCommand[int(m.Sender.ID)]["/sbs2"] = time.Now()
		mu.Unlock()
		cancel()
	})

	bot.Handle("/kpi", func(m *tb.Message) {

		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		command := strings.ToLower(strings.ReplaceAll(m.Text, " ", ""))

		mu.Lock()
		if userLastCommand[int(m.Sender.ID)] == nil {
			userLastCommand[int(m.Sender.ID)] = make(map[string]time.Time)
		}
		lastExecTime := userLastCommand[int(m.Sender.ID)]["/kpi"]
		mu.Unlock()

		if time.Since(lastExecTime).Seconds() < 15 {
			bot.Reply(m, "Debes esperar al menos 15 segundos entre ejecuciones de este comando.")
			return
		}
		if command == "/kpi" {
			bot.Send(m.Chat, " Tomando captura")
			//bot.Send(m.Chat, "👉👈")
			screenshot, err := kpi(ctx, "http://10.115.43.24:3000/login", admin2, pass2)
			if err != nil {
				log.Printf("Error al tomar captura de pantalla: %v", err)
				return
			}
			Photos_response(screenshot, m, bot)
		} else {
			bot.Send(m.Chat, "Comando no reconocido. Por favor, intenta nuevamente.")
		}
		mu.Lock()
		userLastCommand[int(m.Sender.ID)]["/kpi"] = time.Now()
		mu.Unlock()
		cancel()

	})

	bot.Handle("/zone", func(m *tb.Message) {

		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		command := strings.ToLower(strings.ReplaceAll(m.Text, " ", ""))

		mu.Lock()
		if userLastCommand[int(m.Sender.ID)] == nil {
			userLastCommand[int(m.Sender.ID)] = make(map[string]time.Time)
		}
		lastExecTime := userLastCommand[int(m.Sender.ID)]["/zone"]
		mu.Unlock()

		if time.Since(lastExecTime).Seconds() < 15 {
			bot.Reply(m, "Debes esperar al menos 15 segundos entre ejecuciones de este comando.")
			return
		}
		if command == "/zone" {
			bot.Send(m.Chat, " Tomando captura")
			screenshot, err := ZoneOrders(ctx, "http://10.115.43.24:3000/login", admin2, pass2)
			if err != nil {
				log.Printf("Error al tomar captura de pantalla: %v", err)
				return
			}
			Photos_response(screenshot, m, bot)
		} else {
			bot.Send(m.Chat, "Comando no reconocido. Por favor, intenta nuevamente.")
		}
		mu.Lock()
		userLastCommand[int(m.Sender.ID)]["/zone"] = time.Now()
		mu.Unlock()
		cancel()

	})

	bot.Handle("/servidores", func(m *tb.Message) {

		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		command := strings.ToLower(strings.ReplaceAll(m.Text, " ", ""))

		mu.Lock()
		if userLastCommand[int(m.Sender.ID)] == nil {
			userLastCommand[int(m.Sender.ID)] = make(map[string]time.Time)
		}
		lastExecTime := userLastCommand[int(m.Sender.ID)]["/servidores"]
		mu.Unlock()

		if time.Since(lastExecTime).Seconds() < 15 {
			bot.Reply(m, "Debes esperar al menos 15 segundos entre ejecuciones de este comando.")
			return
		}
		if command == "/servidores" {
			bot.Send(m.Chat, " Tomando captura")
			screenshot, err := Servidores(ctx, "http://10.115.43.118:3008/il/grafana/login", user, password)
			if err != nil {
				log.Printf("Error al tomar captura de pantalla: %v", err)
				return
			}
			Photos_response(screenshot, m, bot)
		} else {
			bot.Send(m.Chat, "Comando no reconocido. Por favor, intenta nuevamente.")
		}
		mu.Lock()
		userLastCommand[int(m.Sender.ID)]["/servidores"] = time.Now()
		mu.Unlock()
		cancel()

	})

	bot.Handle("/s11", func(m *tb.Message) {

		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		command := strings.ToLower(strings.ReplaceAll(m.Text, " ", ""))

		mu.Lock()
		if userLastCommand[int(m.Sender.ID)] == nil {
			userLastCommand[int(m.Sender.ID)] = make(map[string]time.Time)
		}
		lastExecTime := userLastCommand[int(m.Sender.ID)]["/s11"]
		mu.Unlock()

		if time.Since(lastExecTime).Seconds() < 15 {
			bot.Reply(m, "Debes esperar al menos 15 segundos entre ejecuciones de este comando.")
			return
		}
		if command == "/s11" {
			bot.Send(m.Chat, " Tomando captura")
			screenshot, err := S11(ctx, "http://10.115.43.118:3008/il/grafana/login", user, password)
			if err != nil {
				log.Printf("Error al tomar captura de pantalla: %v", err)
				return
			}
			Photos_response(screenshot, m, bot)
		} else {
			bot.Send(m.Chat, "Comando no reconocido. Por favor, intenta nuevamente.")
		}
		mu.Lock()
		userLastCommand[int(m.Sender.ID)]["/s11"] = time.Now()
		mu.Unlock()
		cancel()

	})

	bot.Handle("/operacion", func(m *tb.Message) {

		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		command := strings.ToLower(strings.ReplaceAll(m.Text, " ", ""))

		mu.Lock()
		if userLastCommand[int(m.Sender.ID)] == nil {
			userLastCommand[int(m.Sender.ID)] = make(map[string]time.Time)
		}
		lastExecTime := userLastCommand[int(m.Sender.ID)]["/operacion"]
		mu.Unlock()

		if time.Since(lastExecTime).Seconds() < 15 {
			bot.Reply(m, "Debes esperar al menos 15 segundos entre ejecuciones de este comando.")
			return
		}
		if command == "/operacion" {
			bot.Send(m.Chat, " Tomando captura")
			screenshot, err := Operacion(ctx, "http://10.115.43.118:3008/il/grafana/login", user, password)
			if err != nil {
				log.Printf("Error al tomar captura de pantalla: %v", err)
				return
			}
			Photos_response(screenshot, m, bot)
		} else {
			bot.Send(m.Chat, "Comando no reconocido. Por favor, intenta nuevamente.")
		}
		mu.Lock()
		userLastCommand[int(m.Sender.ID)]["/operacion"] = time.Now()
		mu.Unlock()
		cancel()

	})

	bot.Handle("/colas", func(m *tb.Message) {
		command := strings.ToLower(strings.ReplaceAll(m.Text, " ", ""))

		// Obtener los datos desde el scrapper
		data, err := scrapper.FetchDataFromTable("http://localhost:8181/mnesia/tables/ppsinfo")
		if err != nil {
			log.Printf("Error calling FetchQueueSizeData: %v\n", err)
			bot.Send(m.Chat, "Ocurrió un error al obtener los datos.")
			return
		}

		// Convertir los datos a formato JSON
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			log.Printf("Error converting data to JSON: %v\n", err)
			bot.Send(m.Chat, "Ocurrió un error al convertir los datos a JSON.")
			return
		}

		// Enviar el JSON al chat como mensaje
		if command == "/colas" {
			bot.Send(m.Chat, string(jsonData))
		} else {
			bot.Send(m.Chat, "Comando no reconocido. Por favor, intenta nuevamente.")
		}
	})

	log.Println("ChatBot is running. Press CTRL+C to exit.")
	go bot.Start()

	select {}
}

func Photos_response(screenshot []byte, m *tb.Message, bot *tb.Bot) {
	photo := &tb.Photo{
		File: tb.FromReader(bytes.NewReader(screenshot)),
	}
	mention := "@" + m.Sender.FirstName

	bot.Send(m.Chat, mention+" Aqui tienes:")
	bot.SendAlbum(m.Chat, tb.Album{photo})
}

func MLE(ctx context.Context, url, user, password string) ([]byte, error) {
	var buf []byte

	task := chromedp.Tasks{
		emulation.SetDeviceMetricsOverride(2560, 1440, 1, false),
		chromedp.Navigate(url),
		chromedp.WaitVisible("input[name=password]", chromedp.BySearch),
		chromedp.SendKeys("input[name=user]", user, chromedp.BySearch),
		chromedp.SendKeys("input[name=password]", password, chromedp.BySearch),
		chromedp.Click(`button[aria-label="Login button"]`, chromedp.BySearch),
		chromedp.WaitVisible("body", chromedp.BySearch),
		chromedp.Sleep(1 * time.Second),
		chromedp.Navigate("http://10.115.43.118:3008/il/grafana/?orgId=1"),
		chromedp.WaitVisible("body", chromedp.BySearch),
		chromedp.Navigate("http://10.115.43.118:3008/il/grafana/d/sDmADcSIk/mle-flr?orgId=1&refresh=30s"),
		chromedp.WaitVisible("body", chromedp.BySearch),
		chromedp.Sleep(7 * time.Second),
		chromedp.FullScreenshot(&buf, 100),
	}

	err := chromedp.Run(ctx, task)
	if err != nil {
		log.Fatal(err)
	}

	return buf, nil
}

func Soporte_FLR(ctx context.Context, url, user, password string) ([]byte, error) {
	var buf []byte

	task := chromedp.Tasks{
		emulation.SetDeviceMetricsOverride(2560, 1440, 1, false),
		chromedp.Navigate(url),
		chromedp.WaitVisible("input[name=password]", chromedp.BySearch),
		chromedp.SendKeys("input[name=user]", user, chromedp.BySearch),
		chromedp.SendKeys("input[name=password]", password, chromedp.BySearch),
		chromedp.Click(`button[aria-label="Login button"]`, chromedp.BySearch),
		chromedp.WaitVisible("body", chromedp.BySearch),
		chromedp.Sleep(1 * time.Second),
		chromedp.Navigate("http://10.115.43.118:3008/il/grafana/?orgId=1"),
		chromedp.WaitVisible("body", chromedp.BySearch),
		chromedp.Navigate("http://10.115.43.118:3008/il/grafana/d/F2yEI13Vk/flr-operacion?orgId=1&refresh=1m"),
		chromedp.WaitVisible("body", chromedp.BySearch),
		chromedp.Sleep(5 * time.Second),
		chromedp.FullScreenshot(&buf, 100),
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
		emulation.SetDeviceMetricsOverride(2560, 1440, 1, false),
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
		chromedp.Sleep(4 * time.Second),
		chromedp.FullScreenshot(&buf, 100),
	}

	err := chromedp.Run(ctx, task)
	if err != nil {
		log.Fatal(err)
	}

	return buf, nil
}

func SBS(ctx context.Context, url, user, password string) ([]byte, error) {
	var buf []byte

	task := chromedp.Tasks{
		emulation.SetDeviceMetricsOverride(2560, 1440, 1, false),
		chromedp.Navigate(url),
		chromedp.WaitVisible("input[name=password]", chromedp.BySearch),
		chromedp.SendKeys("input[name=user]", user, chromedp.BySearch),
		chromedp.SendKeys("input[name=password]", password, chromedp.BySearch),
		chromedp.Click(`button[aria-label="Login button"]`, chromedp.BySearch),
		chromedp.WaitVisible("body", chromedp.BySearch),
		chromedp.Sleep(1 * time.Second),
		chromedp.Navigate("http://10.115.43.24:3000/"),
		chromedp.WaitVisible("body", chromedp.BySearch),
		chromedp.Navigate("http://10.115.43.24:3000/d/LRJXk-NSk/reporte-de-cierre?orgId=4&from=now-7h&to=now&var-PpsId=All"),
		chromedp.WaitVisible("body", chromedp.BySearch),
		chromedp.Sleep(7 * time.Second),
		chromedp.FullScreenshot(&buf, 100),
	}

	err := chromedp.Run(ctx, task)
	if err != nil {
		log.Fatal(err)
	}

	return buf, nil
}

func SBS_General(ctx context.Context, url, user, password string) ([]byte, error) {
	var buf []byte

	task := chromedp.Tasks{
		emulation.SetDeviceMetricsOverride(2560, 1440, 1, false),
		chromedp.Navigate(url),
		chromedp.WaitVisible("input[name=password]", chromedp.BySearch),
		chromedp.SendKeys("input[name=user]", user, chromedp.BySearch),
		chromedp.SendKeys("input[name=password]", password, chromedp.BySearch),
		chromedp.Click(`button[aria-label="Login button"]`, chromedp.BySearch),
		chromedp.WaitVisible("body", chromedp.BySearch),
		chromedp.Sleep(1 * time.Second),
		chromedp.Navigate("http://10.115.43.24:3000/d/1-Uft5w4k/greymatter-6-1-streaming-store-orders-dashboard?orgId=4&refresh=1m"),
		chromedp.WaitVisible("body", chromedp.BySearch),
		chromedp.Sleep(7 * time.Second),
		chromedp.FullScreenshot(&buf, 100),
	}

	err := chromedp.Run(ctx, task)
	if err != nil {
		log.Fatal(err)
	}

	return buf, nil
}

func kpi(ctx context.Context, url, user, password string) ([]byte, error) {
	var buf []byte

	task := chromedp.Tasks{
		emulation.SetDeviceMetricsOverride(2560, 1440, 1, false),
		chromedp.Navigate(url),
		chromedp.WaitVisible("input[name=password]", chromedp.BySearch),
		chromedp.SendKeys("input[name=user]", user, chromedp.BySearch),
		chromedp.SendKeys("input[name=password]", password, chromedp.BySearch),
		chromedp.Click(`button[aria-label="Login button"]`, chromedp.BySearch),
		chromedp.WaitVisible("body", chromedp.BySearch),
		chromedp.Sleep(1 * time.Second),
		chromedp.Navigate("http://10.115.43.24:3000/"),
		chromedp.Sleep(1 * time.Second),
		chromedp.Navigate("http://10.115.43.24:3000/d/F_8FShESk/kpi-de-seguimiento?orgId=4&refresh=1m"),
		chromedp.WaitVisible("body", chromedp.BySearch),
		chromedp.Sleep(9 * time.Second),
		chromedp.FullScreenshot(&buf, 100),
	}

	err := chromedp.Run(ctx, task)
	if err != nil {
		log.Fatal(err)
	}

	return buf, nil
}

func Servidores(ctx context.Context, url, user, password string) ([]byte, error) {
	var buf []byte

	task := chromedp.Tasks{
		emulation.SetDeviceMetricsOverride(2560, 1440, 1, false),
		chromedp.Navigate(url),
		chromedp.WaitVisible("input[name=password]", chromedp.BySearch),
		chromedp.SendKeys("input[name=user]", user, chromedp.BySearch),
		chromedp.SendKeys("input[name=password]", password, chromedp.BySearch),
		chromedp.Click(`button[aria-label="Login button"]`, chromedp.BySearch),
		chromedp.WaitVisible("body", chromedp.BySearch),
		chromedp.Sleep(1 * time.Second),
		chromedp.Navigate("http://10.115.43.118:3008/il/grafana/?orgId=1"),
		chromedp.WaitVisible("body", chromedp.BySearch),
		chromedp.Sleep(1 * time.Second),
		chromedp.Navigate("http://10.115.43.118:3008/il/grafana/d/jc_66BSIz/servidores?orgId=1&refresh=1m"),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.Sleep(4 * time.Second),
		chromedp.FullScreenshot(&buf, 100),
	}

	err := chromedp.Run(ctx, task)
	if err != nil {
		log.Fatal(err)
	}

	return buf, nil
}

func S11(ctx context.Context, url, user, password string) ([]byte, error) {
	var buf []byte

	task := chromedp.Tasks{
		emulation.SetDeviceMetricsOverride(2560, 1440, 1, false),
		chromedp.Navigate(url),
		chromedp.WaitVisible("input[name=password]", chromedp.BySearch),
		chromedp.SendKeys("input[name=user]", user, chromedp.BySearch),
		chromedp.SendKeys("input[name=password]", password, chromedp.BySearch),
		chromedp.Click(`button[aria-label="Login button"]`, chromedp.BySearch),
		chromedp.WaitVisible("body", chromedp.BySearch),
		chromedp.Sleep(1 * time.Second),
		chromedp.Navigate("http://10.115.43.118:3008/il/grafana/?orgId=1"),
		chromedp.WaitVisible("body", chromedp.BySearch),
		chromedp.Navigate("http://10.115.43.118:3008/il/grafana/d/G-6ygS3Vk/flr-s11-mcu?orgId=1&refresh=1m"),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.Sleep(5 * time.Second),
		chromedp.FullScreenshot(&buf, 100),
	}

	err := chromedp.Run(ctx, task)
	if err != nil {
		log.Fatal(err)
	}

	return buf, nil
}

func Operacion(ctx context.Context, url, user, password string) ([]byte, error) {
	var buf []byte

	task := chromedp.Tasks{
		emulation.SetDeviceMetricsOverride(2560, 1440, 1, false),
		chromedp.Navigate(url),
		chromedp.WaitVisible("input[name=password]", chromedp.BySearch),
		chromedp.SendKeys("input[name=user]", user, chromedp.BySearch),
		chromedp.SendKeys("input[name=password]", password, chromedp.BySearch),
		chromedp.Click(`button[aria-label="Login button"]`, chromedp.BySearch),
		chromedp.WaitVisible("body", chromedp.BySearch),
		chromedp.Sleep(1 * time.Second),
		chromedp.Navigate("http://10.115.43.118:3008/il/grafana/?orgId=1"),
		chromedp.WaitVisible("body", chromedp.BySearch),
		chromedp.Navigate("http://10.115.43.118:3008/il/grafana/d/F2yEI13Vk/flr-operacion?orgId=1&from=now-7h&to=now&refresh=1m"),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.Sleep(5 * time.Second),
		chromedp.FullScreenshot(&buf, 100),
	}

	err := chromedp.Run(ctx, task)
	if err != nil {
		log.Fatal(err)
	}

	return buf, nil
}

func ZoneOrders(ctx context.Context, url, user, password string) ([]byte, error) {
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
		chromedp.Navigate("http://10.115.43.24:3000/d/HLk0BAHSk/zone-orders-by-store-granel-and-store?orgId=4&refresh=1m"),
		chromedp.WaitVisible("body", chromedp.BySearch),
		chromedp.Sleep(11 * time.Second),
		chromedp.FullScreenshot(&buf, 100),
	}

	err := chromedp.Run(ctx, task)
	if err != nil {
		log.Fatal(err)
	}

	return buf, nil
}
