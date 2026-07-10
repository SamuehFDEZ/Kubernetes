package main

// imports
import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

// Estructura de datos para el json
type HandsOn struct {
	Time     time.Time `json:"time"`
	Hostname string    `json:"hostname"`
}

// Funcion para crear servidor simple con validaciones en go
func ServeHTTP(w http.ResponseWriter, r *http.Request) {

    if r.URL.Path != "/"{
        http.NotFound(w,r)
        return
    }

	resp := HandsOn{
		Time:     time.Now(),
		Hostname: os.Getenv("HOSTNAME"),
	}
	jsonResp, err := json.Marshal(&resp)
	if err != nil {
		w.Write([]byte("Error"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

// funcion main
func main() {
	http.HandleFunc("/", ServeHTTP)
	log.Fatal(http.ListenAndServe(":9090", nil))
}