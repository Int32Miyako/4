package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

var formTmpl = []byte(`
<html>
  <body>
    <h3>Search form</h3>
    <form action="/search" method="post">
      Query: <input type="text" name="query"><br>
      OrderField: <input type="text" name="orderField" value="Name"><br>
      OrderBy: <input type="number" name="orderBy" value="0"><br>
      Limit: <input type="number" name="limit" value="10"><br>
      Offset: <input type="number" name="offset" value="0"><br>
      <input type="submit" value="Search">
    </form>
  </body>
</html>
`)

func SearchServerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Write(formTmpl)
		return
	}
	query := r.FormValue("query")
	orderField := r.FormValue("orderField")
	orderBy, _ := strconv.Atoi(r.FormValue("orderBy"))
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	offset, _ := strconv.Atoi(r.FormValue("offset"))

	fmt.Fprintln(w, "you enter: ", query, orderField, orderBy, orderBy, limit, offset)

	data, _ := SearchServer(query, orderField, orderBy, limit, offset)

	pretty, _ := json.MarshalIndent(data, "", "  ")
	fmt.Fprintln(w, string(pretty))

}

func main() {
	http.HandleFunc("/", SearchServerHandler)
	http.ListenAndServe(":8080", nil)
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
func SearchServer(query string, orderField string, orderBy int, limit int, offset int) (*[]UserXml, error) {
	xmlData, err := getDataFromXML("dataset.xml")
	if err != nil {
		panic(err)
	}

	var result []UserXml
	if query != "" {
		for _, user := range xmlData.Users {

			if strings.Contains(user.FirstName+" "+user.SecondName, query) || strings.Contains(user.About, query) {
				result = append(result, user)

			}

		}
	} else {
		result = xmlData.Users
	}

	if orderField == "Name" || orderField == "" {
		sort.Slice(result, func(i, j int) bool {
			name := result[i].FirstName + " " + result[i].SecondName
			name2 := result[j].FirstName + " " + result[j].SecondName
			if orderBy == -1 {
				return name < name2
			} else if orderBy == 1 {
				return name > name2
			}
			return name == name
		})

	} else if orderField == "Id" {
		sort.Slice(result, func(i, j int) bool {
			if orderBy == -1 {
				return result[i].Id < result[j].Id
			} else if orderBy == 1 {
				return result[i].Id > result[j].Id
			}
			return result[i].Id == result[j].Id
		})

	} else if orderField == "Age" {
		sort.Slice(result, func(i, j int) bool {
			if orderBy == -1 {
				return result[i].Age < result[j].Age
			} else if orderBy == 1 {
				return result[i].Age > result[j].Age
			}
			return result[i].Age == result[j].Age
		})

	} else {
		panic(`ErrorBadOrderField`)
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
	//file, err := os.Open(fileName)
	//if err != nil {
	//	return nil, err
	//}
	//defer file.Close()
	//
	//data, err := io.ReadAll(file)
	//if err != nil {
	//	return nil, err
	//}

	// все что выше можно заменить на одну строчку
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
