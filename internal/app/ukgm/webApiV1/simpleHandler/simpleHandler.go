package simple

import "net/http"

func (SimpeApi *SimpleApi) simpleHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}
