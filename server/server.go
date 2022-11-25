package server

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
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

type Response struct {
	List []User
}

var token = map[string]interface{}{
	"good": nil,
}

func handler(w http.ResponseWriter, r *http.Request) {
	clients := new(Clients)
	file, err := os.Open("dataset.xml")
	if err != nil {
		io.WriteString(w, `"Error": "Cant open dataset"`)
		return
	}
	defer file.Close()

	var data []byte
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, scanner.Bytes()...)
	}
	err = xml.Unmarshal(data, &clients)
	if err != nil {
		io.WriteString(w, `"Error": "Error at converting XML to bytes"`)
		return
	}

	access := r.Header.Get("AccessToken")
	if _, ok := token[access]; !ok {
		w.WriteHeader(401)
		io.WriteString(w, `"Error": "Unauthorized user"`)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, `"Error": "Cant convert limit to int"`)
		return
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, `"Error": "Cant convert offset to int"`)
		return
	}
	query := r.URL.Query().Get("query")
	orderField := r.URL.Query().Get("order_field")
	if orderField != "Id" && orderField != "Age" && orderField != "Name" {
		w.WriteHeader(400)
		io.WriteString(w, `"Error": "OrderField invalid"`)
		return
	}
	orderBy, err := strconv.Atoi(r.URL.Query().Get("order_by"))
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, `"Error": "Cant convert orderBy to int"`)
		return
	}
	if orderBy < -1 || orderBy > 1 {
		w.WriteHeader(400)
		io.WriteString(w, `"Error": "OrderBy should be 1, 0 or -1"`)
		return
	}

	queryClients := new(Response)
	if query == "" {
		for _, client := range clients.List {
			user := User{
				ID:     client.Id,
				Name:   client.FirstName + " " + client.LastName,
				Age:    client.Age,
				About:  client.About,
				Gender: client.Gender,
			}
			queryClients.List = append(queryClients.List, user)
		}
	} else {
		for _, client := range clients.List {
			user := User{
				ID:     client.Id,
				Name:   client.FirstName + " " + client.LastName,
				Age:    client.Age,
				About:  client.About,
				Gender: client.Gender,
			}
			if strings.Contains(user.Name, query) || strings.Contains(user.About, query) {
				queryClients.List = append(queryClients.List, user)
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
				return (orderBy * queryClients.List[i].ID) < (orderBy * queryClients.List[j].ID)
			})
		case "Age":
			sort.Slice(queryClients.List, func(i, j int) bool {
				return (orderBy * queryClients.List[i].Age) < (orderBy * queryClients.List[j].Age)
			})
		case "Name":
			sort.Slice(queryClients.List, func(i, j int) bool {
				if orderBy == 1 {
					return (queryClients.List[i].Name) < (queryClients.List[j].Name)
				} else {
					return (queryClients.List[i].Name) > (queryClients.List[j].Name)
				}
			})
		}
	}
	result, err := json.Marshal(queryClients.List[offset:])
	if err != nil {
		io.WriteString(w, `"Error": "Cant convert result to Json"`)
		return
	}
	w.Write(result)

}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Starting server at 8080 port")
	http.ListenAndServe(":8080", nil)
}
