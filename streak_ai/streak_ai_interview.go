package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type requestBody struct {
	Numbers []int `json:"numbers"`
	Target int `json:"target"`
}

func calculate(nums []int, target int) [][]int {
	results := [][]int{}
	for i := 0; i < len(nums) - 1; i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i] + nums[j] == target {
				results = append(results, []int{i, j})
			}
		} 
	}
	return results
}

func server() {
    http.HandleFunc("/find-pairs", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost  {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		p := requestBody{}
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		result := calculate(p.Numbers, p.Target)
		responseBody := struct{
			Solutions [][]int `json:"solutions"`
		}{
			Solutions: result,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&responseBody)
	})
	if err := http.ListenAndServe(":8080", nil); err != http.ErrServerClosed {
        panic(err)
    }
}

func main() {
    // start the server in a goroutine
	fmt.Println("Start")
    go server()

    // wait 1 second to give the server time to start
    time.Sleep(time.Second)

	fmt.Println("Test")
    if err := client(); err != nil {
        fmt.Println(err)
    }
}

func client() error {
    req := &requestBody{
        Numbers: []int{1, 2, 3, 4, 5},
        Target:  6,
    }

    b := new(bytes.Buffer)
    err := json.NewEncoder(b).Encode(req)
    if err != nil {
        return err
    }

    resp, err := http.Post("http://localhost:8080/find-pairs", "application/json", b)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    fmt.Println(resp.Status)

	// read response body
	body, error := io.ReadAll(resp.Body)
	if error != nil {
	   fmt.Println(error)
	}
	// close response body
	resp.Body.Close()
	
	// print response body
	fmt.Println(string(body))
	fmt.Println(resp.Body)
    return nil
}