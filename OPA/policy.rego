package library

import data.books
import data.users
import input

search_books[book] {
  input.book.isbn == books[i].isbn
  input.user.user_type >= books[i].access_level
  book = books[i]
}

list_all_books[books[i]] {
  input.user.user_type >= books[i].access_level
}