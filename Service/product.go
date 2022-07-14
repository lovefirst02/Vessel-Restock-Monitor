package Service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"vessel/Models"
)

func GetProduct() Models.Product {
	var proucts Models.Product

	res, err := http.Get("https://xvessel.co/collections/shoes/products.json")
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&proucts)
	if err != nil {
		fmt.Println(err)
	}
	return proucts
}
