package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	resty "gopkg.in/resty.v1"

	library "personal-learning/go-library/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	libraryData *library.Library // store the library data
	opaURL      string
)

type server struct{}

type errorBody struct {
	Err string `json:"error,omitempty"`
}

func getOpaURL() string {
	port := os.Getenv("OPA_PORT")
	if len(port) == 0 {
		port = "8080"
	}
	return "http://localhost:" + port
}

/*
InitializeLibrary is used to Initialize the libraryData variable.
We are using a global variable to store our library data. This way
we avoid the need of a DB. In real world, it is always a good idea
to use a DB to persist data to prevent data loss due to unavoidable
circumstances.
*/
func InitializeLibrary() {
	libraryData = &library.Library{
		Books: &library.Books{
			Books: make([]*library.Book, 0),
		},
	}
}

/*
search Book by ISBN in the library. If found, return the index of the
book alongwith a boolean indicating that the book exists.
*/
func searchBook(isbn string) (int, bool) {
	i := 0
	exists := false
	var book *library.Book
	for i, book = range libraryData.Books.Books {
		if book.Isbn == isbn {
			exists = true
			break
		}
	}
	if i <= (len(libraryData.Books.Books)-1) && exists == true {
		return i, exists
	}
	return -1, false
}

// GRPC method to return a list of all library book
func (s *server) ListAllBooks(ctx context.Context, query *library.QueryFormat) (*library.Books, error) {
	res := new(library.Books)
	url := opaURL + "/v1/data/library/list_all_books"

	fmt.Println(query)
	input := make(map[string]*library.QueryFormat, 0)
	input["input"] = query

	fmt.Println(input)
	resp, err := resty.R().
		SetBody(input).
		Post(url)

	if err != nil {
		log.Fatalf("%v", err)
		return res, err
	}

	var raw map[string][]*library.Book
	err = json.Unmarshal(resp.Body(), &raw)

	fmt.Println(resp)
	fmt.Println(raw)

	if err != nil {
		log.Fatalf("%v", err)
		return res, err
	}

	if len(raw["result"]) != 0 {
		res.Books = raw["result"]
	}

	return res, nil
}

// GRPC method to add a book to the library.
func (s *server) AddBook(ctx context.Context, query *library.QueryFormat) (*library.Response, error) {
	res := new(library.Response)
	if query.User.UserType != library.User_UserType(2) {
		res.Action = "Add"
		res.Status = 403
		res.Message = "You are not allowed to add books."
		return res, nil
	}

	url := opaURL + "/v1/data/books/" + query.Book.Isbn
	resp, err := resty.R().
		SetBody(query.Book).
		Put(url)

	res.Action = "Add"
	res.Status = int32(resp.StatusCode())
	res.Message = resp.Status()
	fmt.Println(resp, err)

	return res, nil
}

// GRPC method to search a book in the library.
func (s *server) SearchBook(ctx context.Context, query *library.QueryFormat) (*library.Response, error) {
	res := new(library.Response)
	input := make(map[string]*library.QueryFormat, 0)

	input["input"] = query
	url := opaURL + "/v1/data/library/search_books"
	resp, _ := resty.R().
		SetBody(input).
		Post(url)

	var raw map[string][]*library.Book
	err := json.Unmarshal(resp.Body(), &raw)

	if err != nil {
		log.Fatalf("%v", err)
		return res, err
	}

	if len(raw["result"]) == 0 {
		res.Action = "Search"
		res.Status = int32(resp.StatusCode())
		res.Message = "Book not found."
	} else {
		res.Action = "Search"
		res.Status = int32(resp.StatusCode())
		res.Message = resp.Status()
		res.Value = &library.Response_Book{
			Book: raw["result"][0],
		}
	}
	return res, nil
}

// StartHTTPServer - Start the HTTP Server
func StartHTTPServer(clientAddr string) {
	log.Printf("Starting HTTP Server...")
	runtime.HTTPError = CustomHTTPError

	addr := ":8181"
	opts := []grpc.DialOption{grpc.WithInsecure()}
	mux := runtime.NewServeMux()
	if err := library.RegisterLibraryServiceHandlerFromEndpoint(context.Background(), mux, clientAddr, opts); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
	log.Printf("HTTP Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

// StartLibraryServer - Start the GRPC Server
func StartLibraryServer(lis net.Listener) {
	opaURL = getOpaURL()
	log.Printf("Starting GRPC server...")
	s := grpc.NewServer()
	library.RegisterLibraryServiceServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to Serve: %v", err)
		os.Exit(1)
	}
}

// CustomHTTPError - Custom HTTP error on errors in StartHTTPServer
func CustomHTTPError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	const fallback = `{"error": "failed to marshal error message"}`

	w.Header().Set("Content-type", marshaler.ContentType())
	w.WriteHeader(runtime.HTTPStatusFromCode(grpc.Code(err)))
	jErr := json.NewEncoder(w).Encode(errorBody{
		Err: grpc.ErrorDesc(err),
	})

	if jErr != nil {
		w.Write([]byte(fallback))
	}
}
