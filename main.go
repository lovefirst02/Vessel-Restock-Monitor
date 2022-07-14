package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"vessel/Models"
	"vessel/Service"

	"github.com/spf13/viper"
)

var checkProduct []string

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigName("Setting")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("no such config file")
		} else {
			fmt.Println("read config error")
		}
		fmt.Println(err)
	}

	file, err := os.Open("url.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		handle := strings.Split(scanner.Text(), "products/")[1]
		checkProduct = append(checkProduct, handle)
	}
}

func checkTwoMap(m1, m2 map[string][]Models.Variants, k string) (
	v1, v2 []Models.Variants, ok1, ok2 bool) {

	v1, ok1 = m1[k]
	v2, ok2 = m2[k]
	return
}

func getVariant(products Models.Product) map[string][]Models.Variants {
	Detail := make(map[string][]Models.Variants)
	for _, v := range products.Products {
		Detail[v.Handle] = v.Variants
	}
	return Detail
}

func getProduct(products Models.Product) map[string]Models.Detail {
	Detail := make(map[string]Models.Detail)
	for _, v := range products.Products {
		Detail[v.Handle] = v
	}
	return Detail
}

func compareProduct(p1 string, p2 string) bool {
	return p1 == p2
}

func main() {
	before_products := Service.GetProduct()
	before_allVariants := getVariant(before_products)
	for {
		products := Service.GetProduct()
		allProducts := getProduct(products)
		allVariants := getVariant(products)
		for _, v := range checkProduct {
			if data, before_data, ok1, ok2 := checkTwoMap(allVariants, before_allVariants, v); ok1 && ok2 {
				data, _ := json.Marshal(data)
				before_data, _ := json.Marshal(before_data)
				if !compareProduct(string(data), string(before_data)) {
					text := ""
					count := 0
					for _, v := range allVariants[v] {
						if v.Available {
							if count < 3 {
								text += fmt.Sprintf(`[%s](https://xvessel.co/cart/%d:1)  `, v.Title, v.ID)
								count = count + 1
							} else {
								text += fmt.Sprintf(`\n\n[%s](https://xvessel.co/cart/%d:1)  `, v.Title, v.ID)
								count = 1
							}
						}
					}
					before_allVariants[v] = allVariants[v]
					if p, ok := allProducts[v]; ok {
						fmt.Println(time.Now().Format("[2006-01-02 15:04:05]"), v, "-- SEND")
						Service.SendWebHook(p.Title, fmt.Sprintf("https://xvessel.co/collections/shoes/products/%s", v), text, p.Images[0].Src)
					}
				} else {
					fmt.Println(time.Now().Format("[2006-01-02 15:04:05]"), v, "-- Not Change")
				}
			}
		}
		time.Sleep(time.Duration(viper.GetInt("DELAY")) * time.Second)
	}
}
