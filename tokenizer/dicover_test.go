package tokenizer_test

import (
	"os"
	"testing"

	"github.com/skyaxl/golactivity/reader"
	"github.com/skyaxl/golactivity/tokenizer"
)

func TestReadTokens(t *testing.T) {
	r := reader.New("./../../example", ".*")
	pgks, fset, _ := r.Read()
	c, _ := os.Getwd()
	println(c)
	tokenizer.ReadTokens(pgks, fset)
}
