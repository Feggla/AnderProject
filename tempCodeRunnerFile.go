func main() {
	books := getData()
	var wg sync.WaitGroup
	results := make(chan Book, len(books))

	for _, book := range books {
		wg.Add(1)
		go scrapeBookData(&wg, book, results)
	}

	wg.Wait()
	close(results)

	for book := range results {
		fmt.Println(book.OnlinePrice)
		fmt.Println(book.Name)
	}
}

func scrapeBookData(wg *sync.WaitGroup, book Book, results chan<- Book) {
	defer wg.Done()

	url := book.SearchURL
	c := colly.NewCollector()
	c.SetRequestTimeout(30 * time.Second)
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		// fmt.Println(e.ChildText("td.item-note"))
		e.ForEach("tr", func(i int, h *colly.HTMLElement) {
			desc := strings.ToLower(h.ChildText("td.item-note"))
			if !strings.Contains(desc, "fair") {
				price := h.ChildText("span.results-price a")
				price = strings.ReplaceAll(price, "A$", "")
				price = strings.ReplaceAll(price, ",", "")
				num, _ := strconv.ParseFloat(price, 64)
				newnum := int(math.Ceil(num))
				if newnum != 0 && newnum <= book.PricePoint || book.OnlinePrice == 0 {
					book.OnlinePrice = newnum

				}
			}
		})
		results <- book
	})

	// ... (same as before) ...
	// Update book data (book.PricePoint and book.OnlinePrice)

	if err := c.Visit(url); err != nil {
		fmt.Println("Error:", err)
	}
	c.Wait()
}