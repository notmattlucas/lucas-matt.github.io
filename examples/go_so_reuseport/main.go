package main;

import (
	"os"
	"fmt";
	"net/http";
	"log";
);

type FuncHandler func(w http.ResponseWriter, r *http.Request);

func main() {
	name := os.Args[1];
	handler := createHello(name);
	http.HandleFunc("/hello", handler);
	log.Fatal(http.ListenAndServe(":8080", nil));
};

func createHello(name string)FuncHandler {
	var handleHello FuncHandler = func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello ", name, "!");
	};
	return handleHello;
};
