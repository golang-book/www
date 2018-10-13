package main

import (
	"bytes"
	"regexp"
	"strings"
	"unicode"

	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"

	"github.com/alecthomas/chroma/formatters/html"
)

// StripIndentation removes indentation from a block of code
func StripIndentation(str string) string {
	lines := strings.Split(str, "\n")
	minIndent := -1
	for i := range lines {
		cnt := 0
		for _, c := range lines[i] {
			if unicode.IsSpace(c) {
				cnt++
			} else {
				break
			}
		}
		if cnt >= len(lines[i])-1 {
			// all space
			continue
		}
		if minIndent == -1 || cnt < minIndent {
			minIndent = cnt
		}
	}
	for i := range lines {
		if len(lines[i]) > minIndent {
			lines[i] = lines[i][minIndent:]
		}
	}
	return strings.Join(lines, "\n")
}

func IndentWithSpaces(src string) string {
	lines := strings.Split(src, "\n")
	for i := range lines {
		for j := 0; j < len(lines[i]); j++ {
			c := rune(lines[i][j])
			if !unicode.IsSpace(c) {
				break
			}
			if c == '\t' {
				lines[i] = lines[i][:j] + "    " + lines[i][j+1:]
			}
		}
	}
	return strings.Join(lines, "\n")
}

// SyntaxHighlightHTML highlights source code in an html document
func SyntaxHighlightHTML(src string) string {
	re := regexp.MustCompile(`(?s)<blockcode(?: language="(.*?)")?>(.*?)</blockcode>`)
	return re.ReplaceAllStringFunc(src, func(m string) string {
		ms := re.FindStringSubmatch(m)
		code := ms[2]
		code = strings.Replace(code, "&lt;", "<", -1)
		code = StripIndentation(code)
		code = IndentWithSpaces(code)
		code = strings.TrimSpace(code)

		lang := ms[1]
		if lang == "" {
			lang = "text"
		}
		lexer := lexers.Get(lang)
		if lexer == nil {
			lexer = lexers.Fallback
		}
		it, err := lexer.Tokenise(nil, code)
		if err != nil {
			panic(err)
		}

		var buf bytes.Buffer
		formatter := html.New(html.WithClasses())
		style := styles.Fallback
		err = formatter.Format(&buf, style, it)
		if err != nil {
			panic(err)
		}

		return buf.String()
	})
}
