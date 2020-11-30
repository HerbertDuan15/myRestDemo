package dohandle

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var bookId = 3

const filepath = "log.txt"

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

// Init books var as a slice Book struct
var Books []Book

type userError string

func (e userError) Error() string {
	return e.Message()
}

func (e userError) Message() string {
	return string(e)
}

func WriteLogs(s string) {
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		panic("write log error, please check log file log.txt")
	}
	defer file.Close()

	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	write.WriteString("\n" + s + "\n")
	for _, item := range Books {
		tmp, _ := json.Marshal(item)
		write.Write(tmp)
	}
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
}

func HandleFileList(writer http.ResponseWriter, request *http.Request) error {
	//fmt.Println(request.URL.Path)
	//writer.Header().Set("Content-Type", "application/json")
	if strings.Index(request.URL.Path, filepath) == -1 {
		return userError(fmt.Sprintf("path %s must start with %s", request.URL.Path, filepath))
	}
	WriteLogs("web get logs")
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	all, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	writer.Write(all)
	return nil
}

// Get all books
func GetBooks(w http.ResponseWriter, r *http.Request) {
	WriteLogs("client get all books")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Books)
}

// Get single book
func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params

	// Loop through books and find one with the id from the params
	for _, item := range Books {
		if item.ID == params["id"] {
			WriteLogs("client get single book id: " + item.ID + ", name: " + item.Title)
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

// Add new book
func CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	bookId++
	book.ID = strconv.Itoa(bookId)
	WriteLogs("client add new book: " + fmt.Sprint(book))
	Books = append(Books, book)
	json.NewEncoder(w).Encode(book)
}

// Update book
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range Books {
		if item.ID == params["id"] {
			Books = append(Books[:index], Books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			Books = append(Books, book)
			WriteLogs("client update book name: " + fmt.Sprint(book))
			json.NewEncoder(w).Encode(book)
			return
		}
	}
}

// Delete book
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range Books {
		if item.ID == params["id"] {
			WriteLogs("client delete book: " + fmt.Sprint(item))
			Books = append(Books[:index], Books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(Books)
}
