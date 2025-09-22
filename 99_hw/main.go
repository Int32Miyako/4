package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
)

func main() {
	data, _ := getDataFromXML("dataset.xml")
	pretty, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(pretty))
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
// order_field - по какому полю сортировать. Работает по полям ID
// order_by - направление сортировки
// limit - сколько записей вернуть
// offset - сколько записей пропустить от начала
func SearchServer(query string, order_field string, order_by int, limit int, offset int) {
	xmlData, err := getDataFromXML("dataset.xml")
	if err != nil {
		panic(err)
	}
	fmt.Println(&xmlData)

	if query != "" {

	}

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
