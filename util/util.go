package util

import (
	"bufio"
	"html"
	"strings"

	h "golang.org/x/net/html"
)

func PadRight(str string, length int) string {
	for {
		if len(str) >= length {
			return str
		}
		str = str + string(' ')
	}
}

// colorizeQuoteLines colorizes lines that start with ">" in a string.
func colorizeQuoteLines(s string) string {
	scanner := bufio.NewScanner(strings.NewReader(s))
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, ">") {
			line = QuoteStyle.Render(line)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// htmlToText converts a hackernews text message which may contain HTML (like <p> tags) to plain text.
func HtmlToText(s string) string {
	s = html.UnescapeString(s)
	// Replace <p> tags with newlines
	s = strings.ReplaceAll(s, "<p>", "\n")
	s = strings.ReplaceAll(s, "</p>", "\n")

	doc, _ := h.Parse(strings.NewReader(s))
	var f func(*h.Node)
	f = func(n *h.Node) {
		if n.Type == h.ElementNode && n.Data == "i" {
			// Italicize text within <i> tags
			if n.FirstChild != nil {
				italicText := n.FirstChild.Data
				styledText := ItalicStyle.Render(italicText)
				s = strings.Replace(s, "<i>"+italicText+"</i>", styledText, -1)
			}
		} else if n.Data == "a" {
			// Colorize and simplify links in <a> tags
			var href string
			var rel string

			for _, a := range n.Attr {
				if a.Key == "href" {
					href = a.Val
				}
				if a.Key == "rel" {
					rel = a.Val
				}
			}
			if n.FirstChild != nil {
				linkText := n.FirstChild.Data
				styledText := LinkStyle.Render(linkText)
				if rel != "" {
					s = strings.Replace(s, `<a href="`+href+`" rel="`+rel+`">`+linkText+`</a>`, styledText, -1)
				} else {
					s = strings.Replace(s, `<a href="`+href+`">`+linkText+`</a>`, styledText, -1)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	s = colorizeQuoteLines(s)
	return s
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
