package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const intPrintWidth = 8

var (
	c = flag.Bool(

		"c",
		false,
		`The number of bytes in each input file is written to the standard output.  This will cancel out any prior usage
	of the -m option.`,
	)
	l = flag.Bool("l", false, "The number of lines in each input file is written to the standard output.")
	w = flag.Bool("w", false, "The number of words in each input file is written to the standard output.")
	m = flag.Bool(

		"m",
		false,
		`The number of characters in each input file is written to the standard output.  If the current locale does not
	support multibyte characters, this is equivalent to the -c option.  This will cancel out any prior usage of the
	-c option.`,
	)
)

type stats struct {
	bytes int
	chars int
	lines int
	words int
}

func printAndFail(format string, a ...any) {
	fmt.Printf(format, a...)
	os.Exit(1)
}

func getStatsFromReader(r io.Reader) stats {
	br := bufio.NewReader(r)
	st := stats{}
	var prev rune
	for {
		r, s, err := br.ReadRune()
		if err != nil && err != io.EOF {
			printAndFail("reading input: %s", err.Error())
		}
		if err == io.EOF || s == 0 {
			break
		}

		st.bytes += s
		st.chars++

		if r == '\n' {
			st.lines++
		}
		if unicode.IsSpace(r) && !unicode.IsSpace(prev) {
			st.words++
		}
		prev = r
	}
	return st
}

func getStatsFromFile(fileName string) stats {
	f, err := os.Open(fileName)
	if err != nil {
		printAndFail("open file: %s: %s", fileName, err.Error())
	}
	return getStatsFromReader(bufio.NewReader(f))
}

func getStats(fileName string) stats {
	if fileName != "" {
		return getStatsFromFile(fileName)
	}
	return getStatsFromReader(bufio.NewReader(os.Stdin))
}

func statsToString(c, l, w, m bool, fileName string, s stats) string {
	noOptionSet := !m && !c && !l && !w

	sb := strings.Builder{}
	if l || noOptionSet {
		v := strconv.Itoa(s.lines)
		sb.WriteString(strings.Repeat(" ", intPrintWidth-len(v)))
		sb.WriteString(v)
	}
	if w || noOptionSet {
		v := strconv.Itoa(s.words)
		sb.WriteString(strings.Repeat(" ", intPrintWidth-len(v)))
		sb.WriteString(v)
	}
	if c || noOptionSet {
		v := strconv.Itoa(s.bytes)
		sb.WriteString(strings.Repeat(" ", intPrintWidth-len(v)))
		sb.WriteString(v)
	} else if m {
		v := strconv.Itoa(s.chars)
		sb.WriteString(strings.Repeat(" ", intPrintWidth-len(v)))
		sb.WriteString(v)
	}
	if fileName != "" {
		sb.WriteByte(' ')
		sb.WriteString(fileName)
	}

	return sb.String()
}

func wc(c, l, w, m bool, fileName string) string {
	return statsToString(c, l, w, m, fileName, getStats(fileName))
}

func main() {
	flag.Parse()
	fileName := flag.Arg(0)
	fmt.Println(wc(*c, *l, *w, *m, fileName))
}
