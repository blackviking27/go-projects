package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocolly/colly"
)

type item struct {
	Name   string `json:"name"`
	Price  string `json:"price"`
	ImgUrl string `json:"imgurl"`
}

func main() {
	fmt.Println("Scrapping started...")
	// allow the domain that can be crawled
	c := colly.NewCollector(
		colly.AllowedDomains("j2store.net"),
	)

	//store all the results
	var items []item

	// provide a css selector from the page and callback function

	// for the div result
	c.OnHTML("div[itemprop=itemListElement]", func(h *colly.HTMLElement) {
		item := item{
			Name:   h.ChildText("h2.product-title"),
			Price:  h.ChildText("div.sale-price"),
			ImgUrl: h.ChildAttr("img", "src"),
		}
		items = append(items, item)
	})

	// go to the next page
	c.OnHTML("[title=Next]", func(h *colly.HTMLElement) {
		next_page := h.Request.AbsoluteURL(h.Attr("href"))
		c.Visit(next_page)
	})

	// on new request
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Currently scraping", r.URL.String())
	})

	c.Visit("http://j2store.net/demo/index.php/shop")
	for _, item := range items {
		fmt.Println("Name:" + item.Name + "\tPrice:" + item.Price + "\tImage:" + item.ImgUrl + "\n")
	}

	// saving to json file
	content, err := json.Marshal(items)

	if err != nil {
		panic(err)
	}

	os.WriteFile("products.json", content, 0644)
	fmt.Println("Saved data to file successfully!!")
}
