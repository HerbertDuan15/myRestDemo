package main

import (
	"github.com/gorilla/mux"
	"log"
	"myRestDemo/server/dohandle"
	"net/http"
	"os"
)

type appHandler func(writer http.ResponseWriter, request *http.Request) error

type usrError interface {
	error
	Message() string
}

func errWrapper(handler appHandler) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic: %v", r)
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		err := handler(writer, request)

		if err != nil {
			log.Printf("Error occurred handling request: %s", err.Error())

			if userErr, ok := err.(usrError); ok {
				http.Error(writer, userErr.Message(), http.StatusBadRequest)
				return
			}

			code := http.StatusOK
			switch {
			case os.IsNotExist(err):
				code = http.StatusNotFound
			case os.IsPermission(err):
				code = http.StatusForbidden
			default:
				code = http.StatusInternalServerError
			}
			http.Error(writer, http.StatusText(code), code)
		}
	}
}

// Main function
func main() {
	dohandle.WriteLogs("this is log of all transactions\n")
	//init router
	r := mux.NewRouter()

	//book data
	dohandle.Books = append(dohandle.Books, dohandle.Book{ID: "1", Isbn: "438227", Title: "Book One", Author: &dohandle.Author{Firstname: "John", Lastname: "Doe"}})
	dohandle.Books = append(dohandle.Books, dohandle.Book{ID: "2", Isbn: "454555", Title: "Book Two", Author: &dohandle.Author{Firstname: "Steve", Lastname: "Smith"}})

	// Route handles & endpoints
	r.HandleFunc("/log.txt", errWrapper(dohandle.HandleFileList))
	r.HandleFunc("/books", dohandle.GetBooks).Methods("GET")
	r.HandleFunc("/books/{id}", dohandle.GetBook).Methods("GET")
	r.HandleFunc("/books", dohandle.CreateBook).Methods("POST")
	r.HandleFunc("/books/{id}", dohandle.UpdateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", dohandle.DeleteBook).Methods("DELETE")

	//done := make(chan bool)
	//quit := make(chan os.Signal, 1)
	//signal.Notify(quit, os.Interrupt)
	//go func() {
	//	<-quit
	//	//TODO: close gracefully
	//	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	//	defer cancel()
	//	fmt.Println("TODO: close gracefully")
	//	close(done)
	//}()
	// Start server
	log.Fatal(http.ListenAndServe(":8000", r))
	//<-done
	//fmt.Println("Server stopped")

}
