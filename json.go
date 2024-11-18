package main
import (
	"net/http"
	"log"
	"encoding/json"
)
const badWords := ["kerfuffle", "vsharbert", "fornax"]

func jsonHandler(w http.ResponseWriter, r *http.Request){
    type parameters struct {
        // these tags indicate how the keys in the JSON should be mapped to the struct fields
        // the struct fields must be exported (start with a capital letter) if you want them parsed
        Body string `json:"body"`
    }
    type returnVals struct {
        Valid bool `json:"valid,omitempty"`
        ErrReturn string `json:"error,omitempty"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil  {
        // an error will be thrown if the JSON is invalid or has the wrong types
        // any missing fields will simply have their values in the struct set to their zero value
        returnVal := returnVals{
            ErrReturn:"something went wrong",
        }

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
        returnVal := returnVals{
            ErrReturn:"Chirp is too long",
        }
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
        returnVal := returnVals{
            Valid:true,
        }
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

