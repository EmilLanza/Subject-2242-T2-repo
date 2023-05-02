//Filename: main.go
//Basic Middleware example 1

/*package main

import(
	"log"
	"net/http"
)
//Write a middleware
func middlewareA(next http.Handler) http.Handler { //type has to be http.Handler and returns a function
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		//this is executed on the way down to the handler
		log.Println("Executing middleware A")
		next.ServeHTTP(w,r)// Call the next handler in the chain
		//this executed on the way up to the client
		log.Println("Executing middleware A again")
	})
}

func middlewareB(next http.Handler) http.Handler { //type has to be http.Handler and returns a function
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		//this is executed on the way down to the handler
		log.Println("Executing middleware B")
		if r.URL.Path == "/cherry"{
			return //exit
		}
		next.ServeHTTP(w,r)
		//this executed on the way up to the client
		log.Println("Executing middleware B again")
	})
}

//create a handler function
func ourHandler(w http.ResponseWriter, r *http.Request){ //deal with HTTP methods
	log.Println("Executing the handler...")
	w.Write([]byte("CARROTS"))  //send carrots to the client
}

func main(){
	mux := http.NewServeMux() // create multiplexer
	mux.Handle("/check", middlewareA(middlewareB(http.HandlerFunc(ourHandler)))) 
	mux.Handle("/cherry", middlewareA(middlewareB(http.HandlerFunc(ourHandler))))//pass key "/" and value "home"
	log.Print("starting sever on : 4000") //print line
	err := http.ListenAndServe(": 4000", mux) //create server
	log.Fatal(err)
}*/

//Basic Middleware example 2
/*package main

import (
	"log"
	"net/http"
)

func middlewareOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("Executing middlewareOne")//this is executed on the way down to the handler
		next.ServeHTTP(w, r) //call the next handler in the chain
		log.Print("Executing middlewareOne again")//this executed on the way up to the client
	})
}

func middlewareTwo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("Executing middlewareTwo")
		if r.URL.Path == "/foo" {
			return
		}

		next.ServeHTTP(w, r)
		log.Print("Executing middlewareTwo again")
	})
}

func final(w http.ResponseWriter, r *http.Request) { //handler function
	log.Print("Executing finalHandler")
	w.Write([]byte("OK"))
}

func main() {
	mux := http.NewServeMux()

	finalHandler := http.HandlerFunc(final)
	mux.Handle("/", middlewareOne(middlewareTwo(finalHandler)))

	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}*/


//Example 3: Proper Middleware
/*package main

import (
	"log"
	"mime"
	"net/http"
)

func enforceJSONHandler(next http.Handler) http.Handler {//type has to be http.Handler and returns a function
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //return new handler
		contentType := r.Header.Get("Content-Type")// Get the "Content-Type" header from the request

		if contentType != "" {
			mt, _, err := mime.ParseMediaType(contentType)
			if err != nil { // return an error response if the header is incorrect
				http.Error(w, "Malformed Content-Type header", http.StatusBadRequest)
				return
			}

			if mt != "application/json" { // Return an error response if the media type is not "application/json"
				http.Error(w, "Content-Type header must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		}

		next.ServeHTTP(w, r) // Call the next handler in the chain if the "Content-Type" header is valid
	})
}

func final(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK")) // Send message to the client
}

func main() {
	mux := http.NewServeMux() //create multiplexer

	finalHandler := http.HandlerFunc(final) //create handler that write response
	mux.Handle("/", enforceJSONHandler(finalHandler)) //enforceJSONHandler to handler

	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", mux) //start server
	log.Fatal(err)
}	*/

//curl -i localhost:3000
//curl -i -H "Content-Type: application/xml" localhost:3000
//curl -i -H "Content-Type: application/json; charset=UTF-8" localhost:3000

//ThirdParty Middleware

//Example 4: goji/httpauth Security
/*package main

import (
	"log"
	"net/http"

	"github.com/goji/httpauth"
)

func main() {
	authHandler := httpauth.SimpleBasicAuth("alice", "pa$$word") //creates handler with username and password, also middleware

	mux := http.NewServeMux() //create multiplexer

	finalHandler := http.HandlerFunc(final) //create handler to write a response
	mux.Handle("/", authHandler(finalHandler)) //middleware that handles the path and request

	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", mux) //start server
	log.Fatal(err)
}

func final(w http.ResponseWriter, r *http.Request) { //handler
	w.Write([]byte("OK"))
}	
*/
//Example 4: Loggin Handler
/*package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func main() {
	logFile, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)//open log file to write logs
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux() //create multiplexer

	finalHandler := http.HandlerFunc(final) //create handler
	mux.Handle("/", handlers.LoggingHandler(logFile, finalHandler)) //write to log file

	log.Print("Listening on :3000...")
	err = http.ListenAndServe(":3000", mux) //start server
	log.Fatal(err)
}

func final(w http.ResponseWriter, r *http.Request) { //create handler
	w.Write([]byte("OK"))
}*/

//Example 5: Neat loggin handler
package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func newLoggingHandler(dst io.Writer) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(dst, h)//logs all request and response
	}
}

func main() {
	// open a log file in write-only mode, create the file if it doesn't exist, and append new logs to it
	logFile, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		log.Fatal(err)
	}

	loggingHandler := newLoggingHandler(logFile)

	mux := http.NewServeMux() //create multiplexer

	finalHandler := http.HandlerFunc(final) //create handler
	mux.Handle("/", loggingHandler(finalHandler))//path and handler

	log.Print("Listening on :3000...")
	err = http.ListenAndServe(":3000", mux)//start server
	log.Fatal(err)
}

func final(w http.ResponseWriter, r *http.Request) { //handler
	w.Write([]byte("OK"))
}