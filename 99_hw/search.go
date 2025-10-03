package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	orderField := r.URL.Query().Get("orderField")
	orderBy, _ := strconv.Atoi(r.URL.Query().Get("orderBy"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	users, err := SearchServer(query, orderField, orderBy, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
			return nil, errors.New("ErrorBadOrderField")
		}

		sort.Slice(result, less)

	} else if orderBy != 0 {
		return nil, errors.New(ErrorBadOrderBy)
	}

	lenResult := len(result)
	if offset > lenResult {
		return &result, nil
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
