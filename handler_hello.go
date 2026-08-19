package main

import "net/http"

func (cfg *apiConfig) HandlerHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`
		<html>
		<head></head>
			<body>
				<h1>Hello world</h1>
			</body>
		</html>
	`))
}
