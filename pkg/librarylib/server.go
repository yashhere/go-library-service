package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	library "personal-learning/go-library/api"

	ptypes "github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	libraryData *library.Library // store the library data
)

type server struct{}

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
		CurrentBorrowers: make([]*library.Borrower, 0),
	}
}

// get port number from environment variable PORT. If not set, use 50051.
func getPort() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "50051"
	}
	return ":" + port
}

// search Book by ISBN in the library. If found, return the index of the
// book alongwith a boolean indicating that the book exists.
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

// search for a user by ID in borrowers' list
func searchBorrowerData(id string) (int, bool) {
	i := 0
	exists := false
	var borrower *library.Borrower
	for i, borrower = range libraryData.CurrentBorrowers {
		if borrower.IdNo == id {
			exists = true
			break
		}
	}
	if i <= (len(libraryData.CurrentBorrowers)-1) && exists == true {
		return i, exists
	}
	return -1, false
}

// Delete a book from Library
func deleteBook(isbn string) bool {
	index, exists := searchBook(isbn)

	// we do not care about the order of the books in the library.
	if index != -1 && exists == true {
		libraryData.Books.Books[index] = libraryData.Books.Books[len(libraryData.Books.Books)-1]
		libraryData.Books.Books = libraryData.Books.Books[:len(libraryData.Books.Books)-1]
		return true
	}
	return false
}

// Issue a book to some borrower and add it to her borrower's list
func addToBorrowerList(query *library.QueryFormat, book *library.BookIssued) int {
	name := query.Name
	id := query.IdNo
	borrowerType := library.Borrower_BorrowerType(query.BorrowerType)

	index, exists := searchBorrowerData(id)
	if exists == false {
		var borrowerData *library.Borrower
		borrowerData = &library.Borrower{}
		borrowerData.Name = name
		borrowerData.IdNo = id
		borrowerData.BorrowerType = borrowerType
		borrowerData.BooksIssued = append(borrowerData.BooksIssued, book)
		libraryData.CurrentBorrowers = append(libraryData.CurrentBorrowers, borrowerData)
		return len(libraryData.CurrentBorrowers) - 1
	} else {
		libraryData.CurrentBorrowers[index].BooksIssued = append(libraryData.CurrentBorrowers[index].BooksIssued, book)
		return index
	}
}

// GRPC method to return a list of all library book
func (s *server) ListAllBooks(ctx context.Context, empty *library.Empty) (*library.Books, error) {
	b := new(library.Books)
	b.Books = libraryData.Books.Books
	return b, nil
}

// GRPC method to add a book to the library.
func (s *server) AddBook(ctx context.Context, book *library.Book) (*library.Response, error) {
	res := new(library.Response)
	_, exists := searchBook(book.Isbn)
	if exists == false {
		libraryData.Books.Books = append(libraryData.Books.Books, book)
		res.Action = "Add"
		res.Status = "200"
		res.Message = "Book successfully added"
	} else {
		res.Action = "Add"
		res.Status = "200"
		res.Message = "Book already exists in the library. To update, use /updateBook."
	}

	return res, nil
}

// GRPC method to delete a book from the library
func (s *server) DeleteBook(ctx context.Context, query *library.QueryFormat) (*library.Response, error) {
	res := new(library.Response)
	if query.Isbn == "" {
		res.Action = "Delete"
		res.Status = "403"
		res.Message = "ISBN not supplied."
	} else {
		ok := deleteBook(query.Isbn)
		if ok {
			res.Action = "Delete"
			res.Status = "403"
			res.Message = "Book successfully deleted."
		} else {
			res.Action = "Delete"
			res.Status = "403"
			res.Message = "ISBN not found."
		}
	}
	return res, nil
}

// GRPC method to search a book in the library.
func (s *server) SearchBook(ctx context.Context, query *library.QueryFormat) (*library.Response, error) {
	res := new(library.Response)
	index, exists := searchBook(query.Isbn)
	if exists == false {
		res.Action = "Search"
		res.Status = "403"
		res.Message = "Book not found."
	} else {
		res.Action = "Search"
		res.Status = "200"
		res.Message = "Book found."
		res.Value = &library.Response_Book{
			Book: libraryData.Books.Books[index],
		}
	}
	return res, nil
}

// GRPC method to issue a book to the requestor.
func (s *server) IssueBook(ctx context.Context, query *library.QueryFormat) (*library.Response, error) {
	res := new(library.Response)
	if query.Name == "" || query.IdNo == "" || query.Isbn == "" || query.BorrowerType == library.QueryFormat_BorrowerType(0) {
		res.Action = "Issue"
		res.Status = "403"
		res.Message = "Data insufficient"
		return res, nil
	}
	_, exists := searchBook(query.Isbn)

	if exists == false {
		res.Action = "Search"
		res.Status = "403"
		res.Message = "Book not found."
		return res, nil
	}

	i, exists := searchBorrowerData(query.IdNo)

	if exists {
		for _, book := range libraryData.CurrentBorrowers[i].BooksIssued {
			if book.Isbn == query.Isbn {
				res.Action = "IssueBook"
				res.Status = "403"
				res.Message = fmt.Sprintf("Book ISBN %s is already issued to user %s.", book.Isbn, query.Name)
				return res, nil
			}
		}
	}

	res.Action = "IssueBook"
	res.Status = "200"
	res.Message = fmt.Sprintf("Book is issued to user %s.", query.Name)

	issueDate, _ := ptypes.TimestampProto(time.Now())
	issuedBook := &library.BookIssued{
		Isbn:      query.Isbn,
		IssueDate: issueDate,
	}
	i = addToBorrowerList(query, issuedBook)

	res.Value = &library.Response_BorrowerData{
		BorrowerData: &library.Borrower{
			Name:         query.Name,
			IdNo:         query.IdNo,
			BooksIssued:  libraryData.CurrentBorrowers[i].BooksIssued,
			BorrowerType: library.Borrower_BorrowerType(query.BorrowerType),
		},
	}

	return res, nil
}

// GRPC method to return a book
func (s *server) ReturnBook(ctx context.Context, query *library.QueryFormat) (*library.Response, error) {
	res := new(library.Response)
	if query.Name == "" || query.IdNo == "" || query.Isbn == "" || query.BorrowerType == library.QueryFormat_BorrowerType(0) {
		res.Action = "Return"
		res.Status = "403"
		res.Message = "Data insufficient"
		return res, nil
	}

	i, exists := searchBorrowerData(query.IdNo)

	bookFound := false
	if exists {
		for index, book := range libraryData.CurrentBorrowers[i].BooksIssued {
			if book.Isbn == query.Isbn {
				bookFound = true
				libraryData.CurrentBorrowers[i].BooksIssued[index] = libraryData.CurrentBorrowers[i].BooksIssued[len(libraryData.CurrentBorrowers[i].BooksIssued)-1]
				libraryData.CurrentBorrowers[i].BooksIssued = libraryData.CurrentBorrowers[i].BooksIssued[:len(libraryData.CurrentBorrowers[i].BooksIssued)-1]
				res.Action = "ReturnBook"
				res.Status = "200"
				res.Message = "Book successfully returned."
				res.Value = &library.Response_BorrowerData{
					BorrowerData: &library.Borrower{
						Name:         query.Name,
						IdNo:         query.IdNo,
						BooksIssued:  libraryData.CurrentBorrowers[i].BooksIssued,
						BorrowerType: library.Borrower_BorrowerType(query.BorrowerType),
					},
				}
				return res, nil
			}
		}
	}

	if bookFound == false {
		res.Action = "ReturnBook"
		res.Status = "403"
		res.Message = fmt.Sprintf("Book ISBN %s is not issued to the user %s. Can't return.", query.Isbn, query.Name)
	}
	return res, nil
}

// StartLibraryServer - Start the GRPC Server
func StartLibraryServer(lis net.Listener) {
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
