package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/gocolly/colly"
)

var urls = []string{
	"https://timberartdesignuk.com/bedroom.html",
	"https://timberartdesignuk.com/dining-room.html",
	"https://timberartdesignuk.com/living-room.html",
	"https://timberartdesignuk.com/hallway.html",
	"https://timberartdesignuk.com/home-office.html",
}

type Product struct {
	Name    string
	InStock bool
	SKU     string
	Url     string
}

var c = colly.NewCollector()

func main() {
	productUrls := []string{}
	var products []*Product
	c.SetRequestTimeout(time.Minute)
	c.OnHTML("div.product-item", func(h *colly.HTMLElement) {
		productUrls = append(productUrls, h.ChildAttr("a", "href"))
	})
	c.OnHTML("div.product-info-main", func(h *colly.HTMLElement) {
		prod := new(Product)
		prod.Name = h.ChildText("h1")
		prod.InStock = h.ChildText(".instockcheck") == "In stock"
		prod.SKU = h.ChildText("div[itemprop='sku']")
		products = append(products, prod)
	})

	for i := range urls {
		fmt.Printf("visiting %s\n", urls[i])
		if err := c.Visit(urls[i]); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("removing duplicate product links")
	slices.Sort(productUrls)
	productUrls = slices.Compact(productUrls)

	for i := range productUrls {
		fmt.Printf("visiting %s\n", productUrls[i])
		if err := c.Visit(productUrls[i]); err != nil {
			log.Fatal(err)
		}
	}

	file, err := os.Create("products.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := gocsv.MarshalFile(&products, file); err != nil {
		log.Fatal(err)
	}
}
