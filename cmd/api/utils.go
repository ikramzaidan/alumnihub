package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil

}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1024 * 1024
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body cannot accepted")
	}

	return nil
}

func (app *application) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload JSONResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)
}

func (app *application) getFirstImageFromHtml(body string) (string, error) {
	// Memparsing HTML untuk mengekstrak atribut src dari tag img pertama
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return "", err
	}

	// Variabel untuk menyimpan src dari img pertama
	var imgSrc string

	// Fungsi rekursif untuk traversing node HTML
	var f func(*html.Node) bool
	f = func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, attr := range n.Attr {
				if attr.Key == "src" {
					imgSrc = attr.Val
					return true // Ditemukan, hentikan pencarian
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if found := f(c); found {
				return true // Jika sudah ditemukan, hentikan pencarian
			}
		}
		return false
	}

	f(doc)

	// Jika imgSrc masih kosong, kembalikan string kosong
	if imgSrc == "" {
		imgSrc = ""
	}

	return imgSrc, nil
}

// Helper function to sanitize the title for file name
func (app *application) sanitizeFileName(title string) string {
	// Replace spaces with underscores and remove special characters
	return strings.ReplaceAll(strings.Map(func(r rune) rune {
		if r == ' ' {
			return '_'
		}
		// Only allow alphanumeric characters and underscores
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_' {
			return r
		}
		// Replace any other character with an empty string
		return -1
	}, title), "__", "_")
}
