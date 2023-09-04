package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"

	"math"
	"net/http"
	"os"

	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/mailgun/mailgun-go/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
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
var pub_api_key = "pubkey-a874664de1f2657f0f5e1bc673d3c897"
var mailgun_api_key = "key-c652f3e23e6f6e0b6c09cc90bc6214ae"
var mg_domain = "sandbox320af4e7f297428bb363cca4c5b1a624.mailgun.org"

func main() {
	books := getData()
	var wg sync.WaitGroup
	results := make(chan Book, len(books))

	for _, book := range books {
		wg.Add(1)
		bookCopy := book
		go scrapeBookData(&wg, bookCopy, results, books)
	}

	wg.Wait()
	close(results)

	for book := range results {
		fmt.Println(book.OnlinePrice)
		fmt.Println(book.Name)
	}
}

func SendSimpleMessage(domain, apiKey string, book string, price int, books []Book) (string, error) {
	mg := mailgun.NewMailgun(domain, apiKey)
	from := "michaelfeggans@gmail.com"
	to := "michaelfeggans@gmail.com"
	subject := fmt.Sprintf("%s price alert", book)
	text := fmt.Sprintf("Ander, the price of %s has dropped to %d!", book, price)
	message := mg.NewMessage(from, subject, text, to)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := mg.Send(ctx, message)
	fmt.Println(err)
	index := 0
	spreadsheetId := "16vBeSyQTR5IxyOmSi1GHyI-dYXWXShKxGbrg-W0CBLM"
	for ind, val := range books {
		if val.Name == book {
			index = ind+1
		}
	}
	writeRange := fmt.Sprintf("Sheet1!C%d", index)
	var vr sheets.ValueRange
	myval := []interface{}{price}
	vr.Values = append(vr.Values, myval)
	srv := startSpreadsheet()
	_, err = srv.Spreadsheets.Values.Update(spreadsheetId, writeRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet. %v", err)
	}
	return id, err
}

func scrapeBookData(wg *sync.WaitGroup, book Book, results chan<- Book, books []Book) {
	time.Sleep(time.Duration(rand.Intn(120)) * time.Second)
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
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/12.246 Mozilla/5.0 (Windows NT 10.0; Win64; x64; Trident/7.0; AS; rv:11.0) like Gecko",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
		"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36 Edge/16.16299",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.1.2 Safari/603.3.8",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:60.0) Gecko/20100101 Firefox/60.0",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.106 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.101 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 YaBrowser/17.6.0.1633 Yowser/2.5 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.112 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	}

	url := book.SearchURL
	c := colly.NewCollector()
	randomUserAgent := userAgents[rand.Intn(len(userAgents))]
	c.SetRequestTimeout(30 * time.Second)
	c.UserAgent = randomUserAgent
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		r.Headers.Set("User-Agent", randomUserAgent)
	})
	onlinePrice := book.OnlinePrice

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
				if newnum != 0 {
					// fmt.Println(newnum)
					if newnum <= onlinePrice {
						onlinePrice = newnum
					}
					if onlinePrice == 0 {
						onlinePrice = newnum
					}
				}
			}
		})
	})
	if err := c.Visit(url); err != nil {
		fmt.Println("Error:", err)
	}
	book.OnlinePrice = onlinePrice
	if book.PricePoint > onlinePrice && onlinePrice != 0 {
		SendSimpleMessage(mg_domain, mailgun_api_key, book.Name, onlinePrice, books)
	}
	results <- book
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

func startSpreadsheet() *sheets.Service {
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	return srv
}

func getData() []Book {
	var books []Book
	srv := startSpreadsheet()
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
