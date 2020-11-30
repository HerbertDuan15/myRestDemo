package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"myRestDemo/server/dohandle"
	"net/http"
	"net/http/httputil"
)

/*
// Book struct (Model)
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}
// Author struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}
*/

var cliAction = flag.String("method", "GETALL",
	"***methods***:\n "+
		"-method GETALL\n"+
		"GET, for example: -method GET -id 2\n"+
		"POST, for example: -method POST -name 'who are you' -isbn 19961124 -author duanhongjian \n"+
		"DELETE, for example: -method DELETE -id 1\n"+
		"UPDATE, for example: -method UPDATE -id 2 -name 'who are you' -isbn 19961124 -author duanhongjian \n")
var cliBookId = flag.String("id", "1", "input book id")
var cliBookName = flag.String("name", "Book One", "input book name")
var cliBookIsbn = flag.String("isbn", "110110110", "input book Isbn")
var cliBookAuthor = flag.String("author", "Duanhong Jian", "input book author")

func bookToBody(tmpbook dohandle.Book) io.Reader {
	tb, err := json.Marshal(tmpbook)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	fmt.Println(tb)
	reader := bytes.NewReader(tb)
	return reader
}

func main() {
	flag.Parse()

	var tmpbook = dohandle.Book{
		ID:     string(rune(rand.Int() % 100)),
		Isbn:   string(rune(rand.Int() % 100000)),
		Title:  "Book Three",
		Author: &dohandle.Author{Firstname: "Hongjian" + string(rune(rand.Int()%10)), Lastname: "Duan"}}

	request, err := http.NewRequest(http.MethodGet, "http://localhost:8000/books", nil)
	url := "http://localhost:8000/books"
	switch {
	case *cliAction == "GETALL":
		request, err = http.NewRequest(http.MethodGet, url, nil)
	case *cliAction == "GET":
		request, err = http.NewRequest(http.MethodGet, url+"/"+(*cliBookId), nil)
	case *cliAction == "POST":
		tmpbook = dohandle.Book{
			Isbn:   *cliBookIsbn,
			Title:  *cliBookName,
			Author: &dohandle.Author{Firstname: *cliBookAuthor, Lastname: ""}}
		fmt.Println(tmpbook)
		request, err = http.NewRequest(http.MethodPost, url, bookToBody(tmpbook))
	case *cliAction == "UPDATE":
		tmpbook = dohandle.Book{
			ID:     *cliBookId,
			Isbn:   *cliBookIsbn,
			Title:  *cliBookName,
			Author: &dohandle.Author{Firstname: *cliBookAuthor, Lastname: ""}}
		fmt.Println(tmpbook)
		request, err = http.NewRequest(http.MethodPut, url+"/"+(*cliBookId), bookToBody(tmpbook))
	case *cliAction == "DELETE":
		request, err = http.NewRequest(http.MethodDelete, url+"/"+(*cliBookId), nil)
	default:
		fmt.Println(flag.ErrHelp)
	}

	request.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 10_3 like Mac OS X) AppleWebKit/602.1.50 (KHTML, like Gecko) CriOS/56.0.2924.75 Mobile/14E5239e Safari/602.1")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	s, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", s)
}
