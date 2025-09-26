package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"
)

func toUser(u UserXml) User {
	return User{
		Id:     u.Id,
		Name:   u.FirstName + " " + u.SecondName,
		Age:    u.Age,
		About:  u.About,
		Gender: u.Gender,
	}
}

func toUsers(xmlUsers []UserXml) []User {
	users := make([]User, len(xmlUsers))
	for i, u := range xmlUsers {
		users[i] = toUser(u)
	}
	return users
}

func TestSearchHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(searchHandler))
	defer server.Close()

	url := server.URL + "/search?query=Иван&orderField=Age&orderBy=1&limit=10&offset=0"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// вызов ф-ии код ниже

	query := r.URL.Query().Get("query")
	orderField := r.URL.Query().Get("orderField")
	orderBy, _ := strconv.Atoi(r.URL.Query().Get("orderBy"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	fmt.Fprintln(w, "you enter: ", query, orderField, orderBy, orderBy, limit, offset)

	users, _ := SearchServer(query, orderField, orderBy, limit, offset)

	pretty, _ := json.MarshalIndent(users, "", "  ")
	fmt.Fprintln(w, string(pretty))

}

type UsersXml struct {
	Users []UserXml `xml:"row"`
}

type UserXml struct {
	Id         int    `xml:"id"`
	FirstName  string `xml:"first_name"`
	SecondName string `xml:"second_name"`
	Age        int    `xml:"age"`
	About      string `xml:"about"`
	Gender     string `xml:"gender"`
}

// SearchServer занимается поиском данных в файле dataset.xml
// query - что искать. Если пустое - возвращаем все записи
// order_field - по какому полю сортировать. Работает по полям ID, Age, Name
// order_by - направление сортировки
// limit - сколько записей вернуть
// offset - сколько записей пропустить от начала
func SearchServer(query string, orderField string, orderBy int, limit int, offset int) (*[]User, error) {
	xmlData, err := getDataFromXML("dataset.xml")
	if err != nil {
		panic(err)
	}
	users := toUsers(xmlData.Users) // перевод полученных из xml данных в наш тип []User

	var result []User
	if query != "" {
		for _, user := range users {

			if strings.Contains(user.Name, query) || strings.Contains(user.About, query) {
				result = append(result, user)

			}

		}
	} else {
		result = users
	}

	if orderBy == 1 || orderBy == -1 {
		var less func(i, j int) bool

		switch orderField {
		case "Name", "":
			less = func(i, j int) bool {
				if orderBy == -1 {
					return result[i].Name < result[j].Name
				}
				return result[i].Name > result[j].Name
			}
		case "Id":
			less = func(i, j int) bool {
				if orderBy == -1 {
					return result[i].Id < result[j].Id
				}
				return result[i].Id > result[j].Id
			}
		case "Age":
			less = func(i, j int) bool {
				if orderBy == -1 {
					return result[i].Age < result[j].Age
				}
				return result[i].Age > result[j].Age
			}
		default:
			panic(`ErrorBadOrderField`)
		}

		sort.Slice(result, less)

	} else if orderBy != 0 {
		panic(`ErrorBadOrderBy`)
	}

	lenResult := len(result)
	if offset > lenResult {
		fmt.Println("[]")
		panic("offset > lenResult")
	}

	end := offset + limit
	if end > lenResult {
		end = lenResult
	}

	result = result[offset:end]

	return &result, nil
}

func getDataFromXML(fileName string) (*UsersXml, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	users := new(UsersXml)
	err = xml.Unmarshal(data, users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
