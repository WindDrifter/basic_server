package main
import (
	"net/http"
	"log"
	"encoding/json"
    "strings"
    "slices"
    // "github.com/winddrifter/basic_server/internal/database"

)


func jsonHandler(w http.ResponseWriter, r *http.Request){
    type parameters struct {
        // these tags indicate how the keys in the JSON should be mapped to the struct fields
        // the struct fields must be exported (start with a capital letter) if you want them parsed
        Body string `json:"body"`
    }
    type returnVals struct {
        CleanedBody string `json:"cleaned_body,omitempty"`
        ErrReturn string `json:"error,omitempty"`
    }

    returnVal := returnVals{
        ErrReturn:"something went wrong",
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil  {

        dat, returnErr := json.Marshal(returnVal)
        if returnErr != nil {
            log.Printf("Error marshaling return val: %s", err)
            w.WriteHeader(500)
            w.Write(dat)
            return
        } else {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(500)
            w.Write(dat)
            return
        }

    } else if len(params.Body) > 140 {
        returnVal.ErrReturn = "Chirp is too long"
        dat, returnErr := json.Marshal(returnVal)
        if returnErr != nil {
            
            log.Printf("Error marshaling return val: %s", err)
            w.WriteHeader(500)
            return
        } else {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(400)
            w.Write(dat)
            return

        }
    } else {
        returnVal.ErrReturn = ""
        returnVal.CleanedBody = maybeReplaceBadWord(params.Body)

        dat, returnErr := json.Marshal(returnVal)
        
        if returnErr != nil {
            log.Printf("Error marshaling return val: %s", err)
            w.WriteHeader(500)
            return
        } else {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(200)
            w.Write(dat)
            return

        }
    }
}

func maybeReplaceBadWord(body string) string{
    
    bodySplit := strings.Split(body, " ")
    badWords := []string{"kerfuffle", "sharbert", "fornax"}
    var returnVal []string
    for _, word := range bodySplit{
        if slices.Contains(badWords, strings.ToLower(word)) {

            returnVal = append(returnVal, "****")
        } else {
            returnVal = append(returnVal, word)
        }
    }
    return strings.Join(returnVal, " ")


}

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}