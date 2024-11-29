package main

import (
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestFile(t *testing.T) {
	for _, test := range []struct {
		c, l, w, m bool
		wcCMD      []string
	}{
		{true, false, false, false, []string{"wc", "-c", "test.txt"}},
		{false, true, false, false, []string{"wc", "-l", "test.txt"}},
		{false, false, true, false, []string{"wc", "-w", "test.txt"}},
		{false, false, false, true, []string{"wc", "-m", "test.txt"}},
		{false, false, false, false, []string{"wc", "test.txt"}},
	} {
		want, err := exec.Command(test.wcCMD[0], test.wcCMD[1:]...).Output()
		if err != nil {
			t.Fatal(err)
		}
		if got := wc(test.c, test.l, test.w, test.m, "test.txt"); string(got) != string(want)[:len(want)-1] {
			t.Errorf("wc(%v, %v, %v, %v, %q) = %q, want %q", test.c, test.l, test.w, test.m, "test.txt", got, want)
		}
	}
}

func TestStdin(t *testing.T) {
	for _, test := range []struct {
		c, l, w, m bool
		wcCMD      []string
	}{
		{true, false, false, false, []string{"wc", "-c"}},
		{false, true, false, false, []string{"wc", "-l"}},
		{false, false, true, false, []string{"wc", "-w"}},
		{false, false, false, true, []string{"wc", "-m"}},
		{false, false, false, false, []string{"wc"}},
	} {
		want, err := exec.Command(
			"sh",
			"-c",
			strings.Join(append([]string{"cat", "test.txt", "|"}, test.wcCMD...), " "),
		).Output()
		if err != nil {
			t.Fatal(err)
		}

		rd, wr, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		orgStdin := os.Stdin
		os.Stdin = rd

		go func() {
			if got := wc(test.c, test.l, test.w, test.m, ""); string(got) != string(want)[:len(want)-1] {
				t.Errorf("wc(%v, %v, %v, %v, %q) = %q, want %q", test.c, test.l, test.w, test.m, "", got, want)
			}
		}()

		f, err := os.Open("test.txt")
		if err != nil {
			t.Fatal(err)
		}
		io.Copy(wr, f)
		wr.Close()
		os.Stdin = orgStdin
	}
}
