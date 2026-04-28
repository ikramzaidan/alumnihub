package utils

import (
	"strings"

	"golang.org/x/net/html"
)

func GetFirstImageFromHTML(body string) (string, error) {
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return "", err
	}

	var imgSrc string
	var f func(*html.Node) bool
	f = func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, attr := range n.Attr {
				if attr.Key == "src" {
					imgSrc = attr.Val
					return true
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if found := f(c); found {
				return true
			}
		}
		return false
	}

	f(doc)
	return imgSrc, nil
}

func SanitizeFileName(title string) string {
	return strings.ReplaceAll(strings.Map(func(r rune) rune {
		if r == ' ' {
			return '_'
		}
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_' {
			return r
		}
		return -1
	}, title), "__", "_")
}
