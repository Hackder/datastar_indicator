package main

import (
	"fmt"
	"github.com/starfederation/datastar-go/datastar"
	"log"
	"net/http"
	"time"
)

const INDEX_HTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
	<script type="module" src="https://cdn.jsdelivr.net/gh/starfederation/datastar@main/bundles/datastar.js"></script>
    <title>Document</title>
</head>
<body>
    <input
		type="text"
		data-bind-query
		data-on-input__debounce.200ms="@get('/search')"
		data-indicator-searching
	/>
	<div style="height: 1em; width: 1em; background: red;" data-show="$searching"></div>
	<div id="output"></div>

	<script>
	  function triggerBug() {
		const input = document.querySelector('input[data-bind-query]');

		// First "type"
		input.value = 'hello';
		input.dispatchEvent(new InputEvent('input', { bubbles: true, composed: true }));

		// Wait 300ms, then "type" again
		setTimeout(() => {
		  input.value = 'hello world';
		  input.dispatchEvent(new InputEvent('input', { bubbles: true, composed: true }));
		}, 500);
	  }
	</script>

	<button onclick="triggerBug()">Trigger Bug</button>
</body>
</html>`

type Store struct {
	Query string `json:"query"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, INDEX_HTML)
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)

		// Read signals from the request
		store := &Store{}
		if err := datastar.ReadSignals(r, store); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Create a Server-Sent Event writer
		sse := datastar.NewSSE(w, r)

		// Patch elements in the DOM
		sse.PatchElements(fmt.Sprintf(`<div id="output">%s</div>`, store.Query))
	})

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
