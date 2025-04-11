package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

var connUrl = "postgres://postgres:azerty@localhost:5400/mc"

func root(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "server works !")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error")
	}
	fmt.Println(string(body))
}

func add_user(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	var jsonObj map[string]interface{}
	json.Unmarshal(body, &jsonObj)

	username := jsonObj["username"].(string)
	Utype := jsonObj["type"].(string)
	email := jsonObj["email"].(string)
	password := jsonObj["password"].(string)

	conn, err := pgx.Connect(context.Background(), connUrl)
	if err != nil {
		log.Fatal("failed to connect to mcdb, err:", err)
	}

	_, err1 := conn.Exec(context.Background(), "INSERT INTO public.\"Personne\"(username, type, email, password) VALUES('"+username+"', B'"+Utype+"', '"+email+"', '"+password+"')")
	if err1 != nil {
		fmt.Println("failed to insert, err:", err1)
	}

}
func add_guest(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	var jsonObj map[string]interface{}
	json.Unmarshal(body, &jsonObj)

	username := jsonObj["username"].(string)
	Utype := jsonObj["type"].(string)

	conn, err := pgx.Connect(context.Background(), connUrl)
	if err != nil {
		log.Fatal("failed to connect to mcdb, err:", err)
	}

	_, err1 := conn.Exec(context.Background(), "INSERT INTO public.\"Personne\"(username, type) VALUES('"+username+"', B'"+Utype+"')")
	if err1 != nil {
		fmt.Println("failed to insert, err:", err1)
	}
}

func add_product(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	var jsonObj map[string]interface{}
	json.Unmarshal(body, &jsonObj)

	title := jsonObj["title"].(string)
	description := jsonObj["description"].(string)
	category := jsonObj["category"].(string)
	source := jsonObj["source"].(string)
	email := jsonObj["email"].(string)
	price := jsonObj["price"].(string)
	currency := jsonObj["currency"].(string)
	phoneNumber := jsonObj["phoneNumber"].(string)
	personUsername := jsonObj["personUsername"].(string)

	conn, err := pgx.Connect(context.Background(), connUrl)
	if err != nil {
		log.Fatal("failed to connect to mcdb, err:", err)
	}

	var selectPersonId string = "(Select \"personId\" from public.\"Personne\" where username ='" + personUsername + "')"

	_, err1 := conn.Exec(context.Background(), "INSERT INTO public.\"Produit\"(title, description, category, source, email, price, currency, \"phoneNumber\", \"personId\", date) VALUES('"+title+"', '"+description+"', '"+category+"', B'"+source+"', '"+email+"', "+price+", '"+currency+"', '"+phoneNumber+"', "+selectPersonId+", NOW())")
	if err1 != nil {
		fmt.Println("failed to insert, err:", err1)
	}
}

func add_comment(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	var jsonObj map[string]interface{}
	json.Unmarshal(body, &jsonObj)

	ProductId := jsonObj["ProductId"].(string)
	text := jsonObj["text"].(string)
	personUsername := jsonObj["personUsername"].(string)

	conn, err := pgx.Connect(context.Background(), connUrl)
	if err != nil {
		log.Fatal("failed to connect to mcdb, err:", err)
	}

	var selectPersonId string = "(SELECT \"personId\" FROM public.\"Personne\" WHERE username ='" + personUsername + "')"
	fmt.Println("Sous-requête pour personId:", selectPersonId)

	query := "INSERT INTO public.\"Commenter\"(\"personId\", \"username\", \"productId\", date, text) VALUES(" + selectPersonId + ", '" + personUsername + "', '" + ProductId + "', NOW(), '" + text + "')"
	fmt.Println("Requête complète:", query)

	_, err1 := conn.Exec(context.Background(), query)
	if err1 != nil {
		fmt.Println("failed to insert, err:", err1)
	} else {
		fmt.Println("Insertion réussie.")
	}
}

func report_product(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	var jsonObj map[string]interface{}
	json.Unmarshal(body, &jsonObj)

	selectPersonId := jsonObj["selectPersonId"].(string)
	ProductId := jsonObj["ProductId"].(string)
	personUsername := jsonObj["personUsername"].(string)

	conn, err := pgx.Connect(context.Background(), connUrl)
	if err != nil {
		log.Fatal("failed to connect to mcdb, err:", err)
	}

	// var selectPersonId string = "(SELECT \"personId\" FROM public.\"Personne\" WHERE username ='" + personUsername + "')"
	// fmt.Println("Sous-requête pour personId:", selectPersonId)

	query := "INSERT INTO public.\"Signaler\"(\"personId\", \"username\", \"productId\") VALUES(" + selectPersonId + ", '" + personUsername + "', '" + ProductId + "')"
	fmt.Println("Requête complète:", query)

	_, err1 := conn.Exec(context.Background(), query)
	if err1 != nil {
		fmt.Println("failed to insert, err:", err1)
	} else {
		fmt.Println("Insertion réussie.")
	}
}

func buy_product(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	var jsonObj map[string]interface{}
	json.Unmarshal(body, &jsonObj)

	ProductId := jsonObj["ProductId"].(string)
	personUsername := jsonObj["personUsername"].(string)
	cardNumber := jsonObj["cardNumber"].(string)
	CVV := jsonObj["CVV"].(string)
	expDate := jsonObj["expDate"].(string)
	name := jsonObj["name"].(string)

	conn, err := pgx.Connect(context.Background(), connUrl)
	if err != nil {
		log.Fatal("failed to connect to mcdb, err:", err)
	}

	var selectPersonId string = "(SELECT \"personId\" FROM public.\"Personne\" WHERE username ='" + personUsername + "')"
	fmt.Println("Sous-requête pour personId:", selectPersonId)

	query := "INSERT INTO public.\"Acheter\"(\"personId\", \"username\", \"productId\", date, \"cardNumber\", \"CVV\", \"expDate\", name) VALUES(" +
		selectPersonId + ", '" + personUsername + "', '" + ProductId + "', NOW(), '" + cardNumber + "', '" + CVV + "', '" + expDate + "', '" + name + "')"

	fmt.Println("Requête complète:", query)

	_, err1 := conn.Exec(context.Background(), query)
	if err1 != nil {
		fmt.Println("failed to insert, err:", err1)
	} else {
		fmt.Println("Insertion réussie.")
	}
}
func fetch_users(w http.ResponseWriter, r *http.Request) {

	conn, err := pgx.Connect(context.Background(), connUrl)
	if err != nil {
		log.Fatal("failed to connect to mcdb, err:", err)
	}

	rows, _ := conn.Query(context.Background(), `SELECT "personId", "username", "type"::varchar, "email", "password" FROM public."Personne"`)
	var users []map[string]interface{}

	for rows.Next() {
		var id int
		var username, email, password string
		var userType string

		rows.Scan(&id, &username, &userType, &email, &password)

		users = append(users, map[string]interface{}{
			"personId": id,
			"username": username,
			"type":     userType,
			"email":    email,
			"password": password,
		})
	}

	json.NewEncoder(w).Encode(users)
}

func fetch_comments(w http.ResponseWriter, r *http.Request) {

	conn, err := pgx.Connect(context.Background(), connUrl)
	if err != nil {
		log.Fatal("failed to connect to mcdb, err:", err)
	}

	rows, _ := conn.Query(context.Background(), `SELECT "personId", "username", "productId","date"::varchar, "text" FROM public."Commenter"`)
	var comments []map[string]interface{}

	for rows.Next() {
		var personId, productId int
		var username, text, date string
		// var date time.Time

		rows.Scan(&personId, &username, &productId, &date, &text)

		comments = append(comments, map[string]interface{}{
			"personId":  personId,
			"username":  username,
			"productId": productId,
			"date":      date,
			"text":      text,
		})
	}

	json.NewEncoder(w).Encode(comments)
}
func fetch_products(w http.ResponseWriter, r *http.Request) {

	conn, err := pgx.Connect(context.Background(), connUrl)
	if err != nil {
		log.Fatal("failed to connect to mcdb, err:", err)
	}

	rows, _ := conn.Query(context.Background(), `SELECT "productId", "title", "description","category", "source"::varchar,"email", "price" ,"currency","phoneNumber","personId", "username", "date"::varchar  FROM public."Produit"`)
	var products []map[string]interface{}

	for rows.Next() {
		var personId, productId int
		var title, description, category, source, email, price, currency, phoneNumber, username, date string

		rows.Scan(&productId,
			&title,
			&description,
			&category,
			&source,
			&email,
			&price,
			&currency,
			&phoneNumber,
			&personId,
			&username,
			&date)

		products = append(products, map[string]interface{}{
			"productId":   productId,
			"title":       title,
			"description": description,
			"category":    category,
			"source":      source,
			"email":       email,
			"price":       price,
			"currency":    currency,
			"phoneNumber": phoneNumber,
			"personId":    personId,
			"username":    username,
			"date":        date,
		})
	}

	json.NewEncoder(w).Encode(products)
}

func fetch_reports(w http.ResponseWriter, r *http.Request) {

	conn, err := pgx.Connect(context.Background(), connUrl)
	if err != nil {
		log.Fatal("failed to connect to mcdb, err:", err)
	}

	rows, _ := conn.Query(context.Background(), `SELECT "personId", "productId", "username"  FROM public."Signaler"`)
	var reports []map[string]interface{}

	for rows.Next() {
		var personId, productId int
		var username string

		rows.Scan(&personId, &productId, &username)

		reports = append(reports, map[string]interface{}{

			"personId":  personId,
			"productId": productId,
			"username":  username,
		})
	}

	json.NewEncoder(w).Encode(reports)
}
func fetch_payements(w http.ResponseWriter, r *http.Request) {

	conn, err := pgx.Connect(context.Background(), connUrl)
	if err != nil {
		log.Fatal("failed to connect to mcdb, err:", err)
	}

	rows, _ := conn.Query(context.Background(), `SELECT "personId", "username", "productId","date"::varchar,"cardNumber", "CVV", "expDate" ,"name"  FROM public."Acheter"`)
	var payements []map[string]interface{}

	for rows.Next() {
		var personId, productId int
		var username, date, cardNumber, CVV, expDate, name string

		rows.Scan(&personId, &username, &productId, &date, &cardNumber, &CVV, &expDate, &name)

		payements = append(payements, map[string]interface{}{

			"personId":   personId,
			"username":   username,
			"productId":  productId,
			"date":       date,
			"cardNumber": cardNumber,
			"CVV":        CVV,
			"expDate":    expDate,
			"name":       name,
		})
	}

	json.NewEncoder(w).Encode(payements)
}

func main() {

	http.HandleFunc("/add_user", add_user)
	http.HandleFunc("/add_guest", add_guest)
	http.HandleFunc("/add_product", add_product)
	http.HandleFunc("/add_comment", add_comment)
	http.HandleFunc("/report_product", report_product)
	http.HandleFunc("/buy_product", buy_product)
	http.HandleFunc("/fetch_users", fetch_users)
	http.HandleFunc("/fetch_comments", fetch_comments)
	http.HandleFunc("/fetch_products", fetch_products)
	http.HandleFunc("/fetch_reports", fetch_reports)
	http.HandleFunc("/fetch_payements", fetch_payements)

	// http.HandleFunc("/login", login);

	// http.HandleFunc("/fetch_user_products", fetch_products);

	// http.HandleFunc("/fetch_user", root);
	// http.HandleFunc("/update_user", root);
	// http.HandleFunc("/", root);
	// http.HandleFunc("/", root);
	// http.HandleFunc("/", root);
	http.HandleFunc("/", root)
	fmt.Println("server listening on 8000")
	http.ListenAndServe(":8000", nil)
}
