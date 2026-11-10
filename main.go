// Неоновый Разлом: настольная обёртка для Windows.
// Игра встроена в exe и показывается в собственном окне через WebView2
// (есть в Windows 10/11). Если WebView2 нет, игра откроется в браузере.
package main

import (
	"embed"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	webview2 "github.com/jchv/go-webview2"
)

//go:embed index.html fps.html three.min.js
var files embed.FS

// Задаются при сборке через -ldflags "-X main.startPage=... -X main.title=...".
var (
	startPage = "index.html"
	title     = "Неоновый Разлом"
)

// Постоянный порт нужен, чтобы рекорд (localStorage) сохранялся между запусками.
var fixedAddr = "127.0.0.1:47813"

func serve() (string, error) {
	ln, err := net.Listen("tcp", fixedAddr)
	if err != nil {
		if ln, err = net.Listen("tcp", "127.0.0.1:0"); err != nil {
			return "", err
		}
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/")
		if name == "" {
			name = startPage
		}
		data, err := files.ReadFile(name)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		ctype := "text/html; charset=utf-8"
		if strings.HasSuffix(name, ".js") {
			ctype = "text/javascript; charset=utf-8"
		}
		w.Header().Set("Content-Type", ctype)
		w.Header().Set("Cache-Control", "no-store")
		w.Write(data)
	})
	go http.Serve(ln, mux)
	return "http://" + ln.Addr().String() + "/", nil
}

func main() {
	url, err := serve()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	dataDir := ""
	if base, err := os.UserConfigDir(); err == nil {
		dataDir = filepath.Join(base, "NeonRift")
	}

	w := webview2.NewWithOptions(webview2.WebViewOptions{
		AutoFocus: true,
		DataPath:  dataDir,
		WindowOptions: webview2.WindowOptions{
			Title:  title,
			Width:  1280,
			Height: 800,
			Center: true,
		},
	})
	if w == nil {
		// WebView2 не установлен: открываем игру в браузере и продолжаем её раздавать.
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
		select {}
	}
	defer w.Destroy()
	w.Navigate(url)
	w.Run()
}
