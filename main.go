package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"math"
	"math/rand"
	"net/http"
	"os"

	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	// "google.golang.org/grpc/credentials"
	// "github.com/joho/godotenv"
)

type Book struct {
	Name        string
	SearchURL   string
	PricePoint  int
	OnlinePrice int
}

type onlineBook struct {
	Name        string
	SearchURL   string
	PricePoint  int
	Description string
	OnlineDesc  string
}

type Result struct {
	name  string
	price int
}

var wg sync.WaitGroup

// func main() {
// 	books := getData()
// 	x := books[]
// 	for ind, book := range books {
// 		url := book.SearchURL
// 		c := colly.NewCollector()
// 		c.SetRequestTimeout(30 * time.Second)
// 		c.OnRequest(func(r *colly.Request) {
// 			fmt.Println("Visiting", r.URL)
// 		})

// 		c.OnHTML("tbody", func(e *colly.HTMLElement) {
// 			// fmt.Println(e.ChildText("td.item-note"))
// 			e.ForEach("tr", func(i int, h *colly.HTMLElement) {
// 				desc := strings.ToLower(h.ChildText("td.item-note"))
// 				if !strings.Contains(desc, "fair") {
// 					price := h.ChildText("span.results-price a")
// 					price = strings.ReplaceAll(price, "A$", "")
// 					price = strings.ReplaceAll(price, ",", "")
// 					num, _ := strconv.ParseFloat(price, 64)
// 					fmt.Println(num)
// 					newnum := int(math.Ceil(num))
// 					fmt.Println(newnum)
// 					if newnum != 0 && newnum <= books[ind].PricePoint || books[ind].OnlinePrice == 0 {
// 						books[ind].PricePoint = newnum
// 						books[ind].OnlinePrice = newnum

// 					}
// 				}
// 			})
// 		})
// 		if err := c.Visit(url); err != nil {
// 			fmt.Println("Error:", err)
// 		}
// 		c.Wait()
// 		fmt.Println(books[ind].OnlinePrice)
// 		fmt.Println(books[ind].Name)
// 	}
// }

func main() {
	books := getData()
	var wg sync.WaitGroup
	results := make(chan Book, len(books))

	for _, book := range books {
		wg.Add(1)
		bookCopy := book
		go scrapeBookData(&wg, bookCopy, results)
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
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:88.0) Gecko/20100101 Firefox/88.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:87.0) Gecko/20100101 Firefox/87.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Version/14.0.3 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Version/14.0.3 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Version/13.1.2 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.59",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36 Edg/90.0.818.66",
		// Add more user agents as needed
	}
	rand.Seed(time.Now().UnixNano())

	url := book.SearchURL
	c := colly.NewCollector()
	randomUserAgent := userAgents[rand.Intn(len(userAgents))]
	c.SetRequestTimeout(30 * time.Second)
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		r.Headers.Set("User-Agent", randomUserAgent)
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
				// fmt.Println(newnum)
				if newnum != 0 {
					if newnum <= book.PricePoint || book.OnlinePrice == 0 {
						book.OnlinePrice = newnum
					}
				}
			}
		})
	})
	fmt.Println(book.Name)
	fmt.Println(book.OnlinePrice)
	results <- book

	// ... (same as before) ...
	// Update book data (book.PricePoint and book.OnlinePrice)

	if err := c.Visit(url); err != nil {
		fmt.Println("Error:", err)
	}
	c.Wait()
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// func scrape(url string) {

// }

func getData() []Book {
	var books []Book
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	spreadsheetId := "16vBeSyQTR5IxyOmSi1GHyI-dYXWXShKxGbrg-W0CBLM"
	readRange := "Sheet1"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	values := resp.Values
	// var results []string
	for _, val := range values {
		var book Book
		book.Name = val[0].(string)
		book.SearchURL = val[1].(string)
		book.PricePoint, _ = strconv.Atoi(val[2].(string))
		books = append(books, book)
		// results[i] = s.(string)
	}
	return books
}

// fmt.Println(results)

// for i, val := range resp.Values {
// 	book.Name = val[i]
// }

// for i := 0; i < len(resp.Values); i++ {
// 	book.Name = resp.Values[i][0]
// 	book.SearchURL = resp.Values[i][1]
// 	book.PricePoint = resp.Values[i][2]
// 	books = append(books, book)
// }

// for _, row := range resp.Values {
// 	&book.Name, &book.SearchURL, &book.PricePoint
// 	books = append(books, book)

// 	fmt.Println(books)
