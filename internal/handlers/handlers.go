package handlers

import "net/http"

// функция webhook — обработчик HTTP-запроса
func Webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, _ = w.Write([]byte(`
      {
        "response": {
          "text": "Извините, я пока ничего не умею"
        },
        "version": "1.0"
      }
    `))
}
