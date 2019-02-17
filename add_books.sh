#!/bin/bash

curl -X POST \
http://localhost:8181/addBook \
-H 'Content-Type: application/json' \
-H 'Postman-Token: 6400aa6b-4ea1-4dd0-8b70-08c569d400c0' \
-H 'cache-control: no-cache' \
-d '{"book":{"title":"The Predictions of Tycho Dodonus","author":"Tycho Dodonus","isbn":"1241284709","no_of_copies":10, "access_level": 2},"user":{"userType": "Administration"}}'

curl -X POST \
http://localhost:8181/addBook \
-H 'Content-Type: application/json' \
-H 'Postman-Token: d6576647-be9f-4071-a5b7-3fa2cf0dd705' \
-H 'cache-control: no-cache' \
-d '{"book":{"title":"Book of Admittance","isbn":"3071568533","no_of_copies":1, "access_level": 3},"user":{"userType": "Administration"}}'

curl -X POST \
http://localhost:8181/addBook \
-H 'Content-Type: application/json' \
-H 'Postman-Token: 9d8f6801-1a91-49e0-ad30-5a6e384b2b00' \
-H 'cache-control: no-cache' \
-d '{"book":{"title":"Magical Theory","author":"Adalbert Waffling","isbn":"8524504765","no_of_copies":912, "access_level": 1},"user":{"userType": "Administration"}}'

curl -X POST \
http://localhost:8181/addBook \
-H 'Content-Type: application/json' \
-H 'Postman-Token: c8554568-7fb2-4b81-939f-87580689ff33' \
-H 'cache-control: no-cache' \
-d '{"book":{"title":"The Mill on the Floss","author":"George Eliot","isbn":"9488900377","no_of_copies":973, "access_level": 1},"user":{"userType": "Administration"}}'

curl -X POST \
http://localhost:8181/addBook \
-H 'Content-Type: application/json' \
-H 'Postman-Token: 3c86f9e8-b3ca-4f33-90d5-b332a2b0d1b0' \
-H 'cache-control: no-cache' \
-d '{"book":{"title":"Newts of Bognor","author":"Walter Aragon","isbn":"1128959038","no_of_copies":845, "access_level": 1},"user":{"userType": "Administration"}}'

curl -X POST \
http://localhost:8181/addBook \
-H 'Content-Type: application/json' \
-H 'Postman-Token: 3c86g9e8-b3ca-4f33-90d5-b332a2b0d1b0' \
-H 'cache-control: no-cache' \
-d '{"book":{"title":"The Pilgrims Progress","author":"John Bunyan","isbn":"1928959673","no_of_copies":845, "access_level": 3},"user":{"userType": "Administration"}}'