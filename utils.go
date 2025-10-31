package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"iter"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/davecgh/go-spew/spew"
	. "github.com/samber/lo"
)

func Println(args ...any) { fmt.Println(args...) }
func Dump(args ...any) { spew.Config.Indent = "    "; spew.Dump(args...) }

func Sleep(sec int) { time.Sleep(1 * time.Second) }

func F(str string, args ...any) string { return fmt.Sprintf(str, args...) }
func TrimEmptySpace(str string) string { return strings.Join(strings.Fields(str), " ") }
func StrToInt(str string) int {
	str = strings.ReplaceAll(str, ",", "")
	str = strings.ReplaceAll(str, ".", "")
	str = strings.ReplaceAll(str, " ", "")
	return Ternary(str == "", 0, Must(strconv.Atoi(str)))
}

func TemplExec(item any, templ string) string {
	var buf bytes.Buffer
	template.Must(template.New("").Parse(templ)).Execute(&buf, item)
	return buf.String()
}

func CollectSeq[V any](seq iter.Seq[V]) []V {
    var values []V
    seq(func(v V) bool { values = append(values, v); return true })
    return values
}
func CollectSeq2[K comparable, V any](seq2 iter.Seq2[K, V]) []V {
    var values []V
    seq2(func(k K, v V) bool { values = append(values, v); return true })
    return values
}

func HttpGET(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil { return "", err }
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil { return "", err }

	return string(body), nil
}

func HttpGETDocument(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil { return nil, err }
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil { return nil, err }

	return doc, nil
}

func JsonRead[T any](path string) (T, error) {
	var result T

	file, err := os.ReadFile(path)
	if err != nil { return result, err }

	err = json.Unmarshal(file, &result)
	if err != nil { return result, err }

	return result, nil
}

func JsonWrite(path string, data any) error {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil { return err }

	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil { return err }

	tmp := path + ".tmp"
	err = os.WriteFile(tmp, jsonData, 0644)
	if err == nil { return err }

	err = os.Rename(tmp, path)
	if err != nil { os.Remove(tmp); return err }

	return nil
}