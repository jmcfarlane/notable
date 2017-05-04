package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	lorem "github.com/drhodes/golorem"
)

type note struct {
	Content  string `json:"content"`
	Password string `json:"password"`
	Subject  string `json:"subject"`
	Tags     string `json:"tags"`
}

var tagsRegex = regexp.MustCompile(`[^a-z]+`)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	client := http.Client{}
	for i := 0; i < 50; i++ {
		tags := tagsRegex.ReplaceAll([]byte(lorem.Sentence(4, 15)), []byte(" "))
		noteJSON, _ := json.Marshal(note{
			Content: func() string {
				var s []string
				for i := 0; i < 5; i++ {
					s = append(s, lorem.Paragraph(1, 50))
				}
				return strings.Join(s, "\n\n")
			}(),
			Password: func() string {
				if rand.Intn(20) < 5 {
					return lorem.Word(0, 5)
				}
				return ""
			}(),
			Subject: strings.TrimSuffix(lorem.Sentence(5, 10), "."),
			Tags:    strings.ToLower(string(tags)),
		})
		req, _ := http.NewRequest("POST", "http://localhost:8080/api/note/create", bytes.NewBuffer(noteJSON))
		resp, err := client.Do(req)
		fmt.Println(resp, err)
	}
}
