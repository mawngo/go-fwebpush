package main

import (
	"errors"
	"flag"
	"github.com/mawngo/go-fwebpush"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	addr := flag.String("addr", ":8080", "serving address")
	vapid := flag.String("vapid", "", "vapid keypair")
	flag.Parse()

	var priv, pub string
	if *vapid == "" {
		*vapid = mustReadFile(".vapid.txt")
	}

	if *vapid == "" {
		var err error
		priv, pub, err = fwebpush.GenerateVAPIDKeys()
		if err != nil {
			println("Error generating VAPID key pair:", err)
		}
	} else {
		keypair := strings.Split(*vapid, ":")
		priv = keypair[0]
		pub = keypair[1]
	}
	mustCreateFile(".vapid.txt", priv+":"+pub)

	http.HandleFunc("GET /", func(w http.ResponseWriter, _ *http.Request) {
		p := map[string]string{
			"Priv": priv,
			"Pub":  pub,
		}
		t, _ := template.ParseFiles("index.gohtml")
		if err := t.Execute(w, p); err != nil {
			panic(err)
		}
	})

	http.HandleFunc("POST /sub", func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		println(string(b))
		mustCreateFile(".subscription.json", string(b))
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("GET /service-worker.js", func(w http.ResponseWriter, _ *http.Request) {
		f, err := os.Open("service-worker.js")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		defer func() {
			err := f.Close()
			if err != nil {
				panic(err)
			}
		}()
		w.Header().Add("Content-Type", "text/javascript")
		w.WriteHeader(http.StatusOK)
		_, err = io.Copy(w, f)
		if err != nil {
			panic(err)
		}
	})

	println("Listening on " + *addr)
	err := http.ListenAndServe(*addr, nil)
	if !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func mustCreateFile(filename string, data string) {
	err := (func() (err error) {
		fi, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer func() {
			cerr := fi.Close()
			if err == nil {
				err = cerr
			}
		}()
		_, err = fi.WriteString(data)
		return err
	})()
	if err != nil {
		panic(err)
	}
}

func mustReadFile(filename string) string {
	b, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ""
		}
		panic(err)
	}
	return string(b)
}
