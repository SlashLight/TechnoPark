package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Client struct {
	Id        int    `xml:"id"`
	Age       int    `xml:"age"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Gender    string `xml:"gender"`
	About     string `xml:"about"`
}

type Clients struct {
	List []Client `xml:"row"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	clients := new(Clients)
	_, err := w.Write([]byte("Hello User"))
	if err != nil {
		log.Printf("%v", err)
		return
	}
	file, err := os.Open("dataset.xml")
	if err != nil {
		fmt.Println("Cant open dataset")
		return
	}

	var data []byte
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, scanner.Bytes()...)
	}
	err = xml.Unmarshal(data, &clients)
	if err != nil {
		w.Write([]byte("Error at converting XML to bytes"))
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		w.Write([]byte("Cant convert limit to int"))
		return
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		w.Write([]byte("Cant convert offset to int"))
		return
	}
	query := r.URL.Query().Get("query")
	orderField := r.URL.Query().Get("order_field")
	if orderField != "Id" && orderField != "Age" && orderField != "Name" {
		w.Write([]byte("You can sort clients by Id, Age, Name"))
		return
	}
	orderBy, err := strconv.Atoi(r.URL.Query().Get("order_by"))
	if err != nil {
		w.Write([]byte("Cant convert orderBy to int"))
		return
	}
	if orderBy < -1 || orderBy > 1 {
		w.WriteHeader(400)
		w.Write([]byte("OrderBy should be 1, 0 or -1"))
		return
	}

	queryClients := new(Clients)
	if query == "" {
		queryClients.List = clients.List[:limit]
	} else {
		for _, client := range clients.List {
			name := client.FirstName + client.LastName
			if strings.Contains(name, query) || strings.Contains(client.About, query) {
				queryClients.List = append(queryClients.List, client)
			}
			if len(queryClients.List) == limit {
				break
			}
		}
	}

	if orderBy != 0 {
		switch orderField {
		case "Id":
			sort.Slice(queryClients.List, func(i, j int) bool {
				return (orderBy * queryClients.List[i].Id) < (orderBy * queryClients.List[j].Id)
			})
		case "Age":
			sort.Slice(queryClients.List, func(i, j int) bool {
				return (orderBy * queryClients.List[i].Age) < (orderBy * queryClients.List[j].Age)
			})
		case "Name":
			sort.Slice(queryClients.List, func(i, j int) bool {
				if orderBy == 1 {
					return (queryClients.List[i].FirstName + queryClients.List[i].LastName) < (queryClients.List[j].FirstName + queryClients.List[j].LastName)
				} else {
					return (queryClients.List[i].FirstName + queryClients.List[i].LastName) > (queryClients.List[j].FirstName + queryClients.List[j].LastName)
				}
			})
		}
	}
	result, err := json.Marshal(queryClients.List[offset:])
	if err != nil {
		fmt.Println("Cant convert result to Json")
		return
	}
	_, err = w.Write([]byte(result))
	if err != nil {
		log.Printf("%v", err)
		return
	}

}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Starting server at 8080 port")
	http.ListenAndServe(":8080", nil)
}
