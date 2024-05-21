package domyApi

import (
	"fmt"
	"testing"
)

type TestApi struct {
	Phone      string `json:"phone"`
	Password   string `json:"password"`
	FirebaseId string `json:"firebaseid"`
	DeviceId   string `json:"deviceid"`
}

type Sister struct {
	Id_sdm string `url:"id_sdm" json:"id_sdm"`
}

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type Data struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func TestGet(t *testing.T) {
	myData, err := Get[Data]("https://dog.ceo/api/breeds/image/random")
	fmt.Println(myData.Message, err)
}

func TestGetStruct(t *testing.T) {
	dt := Sister{
		Id_sdm: "8fe6735c-6e28-43e7-9eb3-3ae092bbcd62",
	}
	url := "https://httpbin.org/get"
	res := GetStruct(dt, url)
	fmt.Println("GetStruct : ", res)
}

func TestPostStruct(t *testing.T) {
	dt := TestApi{
		Phone:      "+6285155476774",
		Password:   "#P@ssw0rd",
		FirebaseId: "123",
		DeviceId:   "6580fb6e714844ca",
	}
	url := "https://httpbin.org/post"
	res, err := PostStruct[Response](dt, url)
	if err != "" {
		t.Fatalf("PostStruct failed: %s", err)
	}
	fmt.Println("PostStruct : ", res)
}

func TestRequestStructWithToken(t *testing.T) {
	dt := Sister{
		Id_sdm: "8fe6735c-6e28-43e7-9eb3-3ae092bbcd62",
	}
	urlGet := "https://httpbin.org/get"
	urlPost := "https://httpbin.org/post"

	var result interface{}
	var err string

	// Test GetStructWithToken
	result, err = GetStructWithToken[interface{}]("token", "dsfdsfdsfdsfdsf", dt, urlGet)
	if err != "" {
		t.Fatalf("GetStructWithToken failed: %s", err)
	}
	fmt.Println("GetStructWithToken result:", result)

	// Test PostStructWithToken
	dta := TestApi{
		Phone:      "+6285155476774",
		Password:   "#P@ssw0rd",
		FirebaseId: "123",
		DeviceId:   "6580fb6e714844ca",
	}
	result, err = PostStructWithToken[interface{}]("Login", "dsfdsfdsfdsfdsf", dta, urlPost)
	if err != "" {
		t.Fatalf("PostStructWithToken failed: %s", err)
	}
	fmt.Println("PostStructWithToken result:", result)

	// Test PostStructWithBearer
	result, err = PostStructWithBearer[interface{}]("dsfdsfdsfdsfdsf", dta, urlPost)
	if err != "" {
		t.Fatalf("PostStructWithBearer failed: %s", err)
	}
	fmt.Println("PostStructWithBearer result:", result)

	// Test GetStructWithBearer
	result, err = GetStructWithBearer[interface{}]("dsfdsfdsfdsfdsf", dt, urlGet)
	if err != "" {
		t.Fatalf("GetStructWithBearer failed: %s", err)
	}
	fmt.Println("GetStructWithBearer result:", result)
}
