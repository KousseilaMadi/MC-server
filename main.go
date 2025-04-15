package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"context"
	"log"
	"encoding/base64"
	// "errors"
	"github.com/jackc/pgx/v5"
	// "github.com/jackc/pgconn"
	"github.com/google/uuid"
)

type Status struct {
	Status string 
}


var connUrl = "postgres://postgres:ukiyokaru@localhost:5432/mcdb"

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
	fullName := jsonObj["fullName"].(string)
	email := jsonObj["email"].(string)
	password := jsonObj["password"].(string)

	conn, err := pgx.Connect(context.Background(), connUrl)
	if err != nil {
		log.Fatal("failed to connect to mcdb, err:", err)
	}

	rows, err2 := conn.Query(context.Background(), "SELECT * FROM public.\"Personne\" WHERE username = $1", username)

	if err2 != nil {
		fmt.Println(err2)
	}
	defer rows.Close()
	count :=0
	for rows.Next(){
		count++
	}

	var status Status
	if count == 0{

		_, err1 := conn.Exec(context.Background(), "INSERT INTO public.\"Personne\"(username, \"fullName\", email, password, type) VALUES($1, $2, $3, $4, B'0')", username, fullName, email, password)
		if err1 != nil {
			fmt.Println("failed to insert, err:", err1)
		}else{

			status = Status{Status: "OK"}

			json.NewEncoder(w).Encode(status)
		}
	}else{
		
		status = Status{Status: "USERNAME_EXISTS"}

		json.NewEncoder(w).Encode(status)
	}

}
func add_guest(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	var jsonObj map[string]interface{}
	json.Unmarshal(body, &jsonObj)

	id := uuid.New()
	b := id[:]
	short := base64.RawURLEncoding.EncodeToString(b)
	username := "guest_"+short

	conn, err := pgx.Connect(context.Background(), connUrl)
	if err != nil {
		log.Fatal("failed to connect to mcdb, err:", err)
	}
	var status Status

	_, err1 := conn.Exec(context.Background(), "INSERT INTO public.\"Personne\"(username, type) VALUES($1, B'1')", username)
	if err1 != nil {
		fmt.Println("failed to insert, err:", err1)
	}else{
		
		status = Status{Status: "OK"}

		json.NewEncoder(w).Encode(status)
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
	var productId string
	switch v := jsonObj["productId"].(type){
	case float64:
		productId = fmt.Sprintf("%.0f", v)
	case string: productId = v
	default :http.Error(w, "Invalid ProductId", http.StatusBadRequest)
    return
	}
	text := jsonObj["text"].(string)
	personUsername := jsonObj["personUsername"].(string)

	conn, err := pgx.Connect(context.Background(), connUrl)
	if err != nil {
		log.Fatal("failed to connect to mcdb, err:", err)
	}

	var selectPersonId string = "(SELECT \"personId\" FROM public.\"Personne\" WHERE username ='" + personUsername + "')"
	fmt.Println("Sous-requête pour personId:", selectPersonId)

	query := "INSERT INTO public.\"Commenter\"(\"personId\", \"productId\", date, text) VALUES(" + selectPersonId + ", " + productId + ", NOW(), '" + text + "')"
	fmt.Println("Requête complète:", query)

	_, err1 := conn.Exec(context.Background(), query)
	if err1 != nil {
		fmt.Println("failed to insert, err:", err1)
	} else {
		fmt.Println("Insertion réussie.")
	}
}

func add_to_cart(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	var jsonObj map[string]interface{}
	json.Unmarshal(body, &jsonObj)

	productId := jsonObj["productId"].(float64)
	username := jsonObj["username"].(string)

	conn, err := pgx.Connect(context.Background(), connUrl)
	if err != nil {
		log.Fatal("failed to connect to mcdb, err:", err)
	}

	query := "INSERT INTO public.\"Panier\"(\"personId\", \"productId\", date) VALUES((SELECT \"personId\" FROM public.\"Personne\" WHERE username = $1), $2, NOW())"
	fmt.Println("Requête complète:", query)

	_, err1 := conn.Exec(context.Background(), query, username, productId)
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

	rows, _ := conn.Query(context.Background(), `SELECT "personId", "username", "fullName", "type"::varchar, "email", "password" FROM public."Personne"`)
	var users []map[string]interface{}

	for rows.Next() {
		var id int
		var username, fullName, email, password  string
		var userType string

		rows.Scan(&id, &username, &fullName, &userType, &email, &password)

		users = append(users, map[string]interface{}{
			"personId": id,
			"username": username,
			"fullName": fullName,
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	var jsonObj map[string]interface{}
	json.Unmarshal(body, &jsonObj)

	username := jsonObj["username"].(string)

	rows, _ := conn.Query(
		context.Background(),
		`SELECT 
		p."productId",
		p."title",
		p."description",
    	p."category",
    	p."source"::varchar,
    	p."email",
    	p."price",
    	p."currency",
    	p."phoneNumber",
    	p."personId",
    	p."date"::varchar,
    	r.rating
		FROM public."Produit" p
		LEFT JOIN public."Rating" r 
    	ON p."productId" = r."productId" AND r."personId" = (SELECT "personId" from public."Personne" where username = $1)`,
		username,
	)
	var products []map[string]interface{}

	for rows.Next() {
		var personId, productId, rating int
		var title, description, category, source, email, price, currency, phoneNumber, date string

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
			&date,
			&rating,
		)

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
			"date":        date,
			"rating":      rating,
		})
	}

	json.NewEncoder(w).Encode(products)
}


func fetch_product(w http.ResponseWriter, r *http.Request) {

	conn, err := pgx.Connect(context.Background(), connUrl)
	if err != nil {
		log.Fatal("failed to connect to mcdb, err:", err)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	var jsonObj map[string]interface{}
	json.Unmarshal(body, &jsonObj)

	username := jsonObj["username"].(string)
	productIdin := jsonObj["productId"].(float64)

	rows, _ := conn.Query(
		context.Background(),
		`SELECT 
		p."productId",
		p."title",
		p."description",
    	p."category",
    	p."source"::varchar,
    	p."email",
    	p."price",
    	p."currency",
    	p."phoneNumber",
    	p."personId",
    	p."date"::varchar,
    	r.rating
		FROM public."Produit" p
		JOIN public."Rating" r 
    	ON p."productId" = r."productId" AND r."personId" = (SELECT "personId" from public."Personne" where username = $1 AND p."productId" = $2)`,
		username,productIdin,
	)
	var product map[string]interface{}

		var personId, productId, rating int
		var title, description, category, source, email, price, currency, phoneNumber, date string

	if rows.Next() {
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
			&date,
			&rating,
		)

		product = map[string]interface{}{
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
			"date":        date,
			"rating":      rating,
		}

		json.NewEncoder(w).Encode(product)
	}else{
		http.Error(w, "Product not found or not rated by user", http.StatusNotFound)
	}
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


func login(w http.ResponseWriter, r *http.Request) {

	conn, err := pgx.Connect(context.Background(), connUrl)
	if err != nil {
		log.Fatal("failed to connect to mcdb, err:", err)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}


	var jsonObj map[string]interface{}
	json.Unmarshal(body, &jsonObj)

	username := jsonObj["username"].(string)
	password := jsonObj["password"].(string)



	var exists bool
	err1 := conn.QueryRow(context.Background(), `
	SELECT EXISTS(
		SELECT 1 FROM public."Personne" WHERE username = $1 AND password = $2
		)
	`, username, password).Scan(&exists)

	if err1 != nil {
		fmt.Print("error:", err1)
	}
	var status Status
	if exists {
		status = Status{Status: "OK"}
	} else{
		status = Status{Status: "DOES_NOT_EXIST"}
	}

	json.NewEncoder(w).Encode(status)
}


func fetch_user_products(w http.ResponseWriter, r *http.Request) {

	conn, err := pgx.Connect(context.Background(), connUrl)
	if err != nil {
		log.Fatal("failed to connect to mcdb, err:", err)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}


	var jsonObj map[string]interface{}
	json.Unmarshal(body, &jsonObj)

	username := jsonObj["username"].(string)

	rows, _ := conn.Query(context.Background(), `SELECT "productId", "title", "description","category", "source"::varchar,"email", "price" ,"currency","phoneNumber", "date"::varchar  FROM public."Produit" WHERE personId = (SELECT personId FROM public."Personne" WHERE username = $1)`, username)
	var products []map[string]interface{}

	for rows.Next() {
		var productId int
		var title, description, category, source, email, price, currency, phoneNumber, date string

		rows.Scan(&productId,
			&title,
			&description,
			&category,
			&source,
			&email,
			&price,
			&currency,
			&phoneNumber,
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
			"date":        date,
		})
	}

	json.NewEncoder(w).Encode(products)
}

func fetch_user(w http.ResponseWriter, r *http.Request) {

	conn, err := pgx.Connect(context.Background(), connUrl)
	if err != nil {
		log.Fatal("failed to connect to mcdb, err:", err)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	var jsonObj map[string]interface{}
	json.Unmarshal(body, &jsonObj)

	usernameIn := jsonObj["username"].(string)

	rows, _ := conn.Query(context.Background(), `SELECT "personId", "username", "fullName", "type"::varchar, "email" FROM public."Personne" WHERE username = $1`, usernameIn)

	if rows.Next(){

		var id int
		var username, fullName, email, userType string
		rows.Scan(&id, &username, &fullName, &userType, &email)
	user := map[string]interface{}{
		"personId" : id,
		"username" : username,
		"fullName" : fullName,
		"type" : userType,
		"email" : email,
	}
	
	json.NewEncoder(w).Encode(user)
}
}



// func update_user(w http.ResponseWriter, r *http.Request) {

// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	var jsonObj map[string]interface{}
// 	json.Unmarshal(body, &jsonObj)

// 	username := jsonObj["username"].(string)
// 	fullName := jsonObj["fullName"].(string)
// 	Utype := jsonObj["type"].(string)
// 	email := jsonObj["email"].(string)
// 	password := jsonObj["password"].(string)

// 	conn, err := pgx.Connect(context.Background(), connUrl)
// 	if err != nil {
// 		log.Fatal("failed to connect to mcdb, err:", err)
// 	}

// 	_, err1 := conn.Exec(context.Background(), "UPDATE public.\"Personne\"(username, fullName, type, email, password) VALUES('"+username+"', '" + fullName + "' B'"+Utype+"', '"+email+"', '"+password+"')")
// 	if err1 != nil {
// 		fmt.Println("failed to insert, err:", err1)
// 	}

// }


func fetch_cart(w http.ResponseWriter, r *http.Request) {

	conn, err := pgx.Connect(context.Background(), connUrl)
	if err != nil {
		log.Fatal("failed to connect to mcdb, err:", err)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}


	var jsonObj map[string]interface{}
	json.Unmarshal(body, &jsonObj)

	username := jsonObj["username"].(string)

	rows, _ := conn.Query(context.Background(), `
	SELECT 
	pr."productId", 
	"title", 
	"description",
	"category", 
	"source"::varchar,
	pr."email", 
	"price" ,
	"currency",
	"phoneNumber", 
	pr."date"::varchar  
	FROM public."Produit" pr
	JOIN public."Panier" pa ON pa."productId" = pr."productId"
	JOIN public."Personne" pe ON pe."personId" = pa."personId"
	WHERE pe."personId" = 
	(SELECT "personId" 
	FROM public."Personne" 
	WHERE username = $1)`, username)
	var products []map[string]interface{}

	for rows.Next() {
		var productId int
		var title, description, category, source, email, price, currency, phoneNumber, date string

		rows.Scan(&productId,
			&title,
			&description,
			&category,
			&source,
			&email,
			&price,
			&currency,
			&phoneNumber,
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
			"date":        date,
		})
	}

	json.NewEncoder(w).Encode(products)
}



func delete_from_cart(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	var jsonObj map[string]interface{}
	json.Unmarshal(body, &jsonObj)

	productId := jsonObj["productId"].(float64)
	personId := jsonObj["personId"].(float64)


	conn, err := pgx.Connect(context.Background(), connUrl)
	if err != nil {
		log.Fatal("failed to connect to mcdb, err:", err)
	}
	var status Status

	_, err1 := conn.Exec(context.Background(), "DELETE FROM public.\"Panier\" p WHERE p.\"personId\" = $1 AND p.\"productId\" = $2", int(personId), int(productId))
	if err1 != nil {
		fmt.Println("failed to insert, err:", err1)
	}else{
		
		status = Status{Status: "OK"}

		json.NewEncoder(w).Encode(status)
	}
}


func delete_user(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	var jsonObj map[string]interface{}
	json.Unmarshal(body, &jsonObj)

	username := jsonObj["username"].(string)
	password := jsonObj["password"].(string)


	conn, err := pgx.Connect(context.Background(), connUrl)
	if err != nil {
		log.Fatal("failed to connect to mcdb, err:", err)
	}
	var status Status

	_, err1 := conn.Exec(context.Background(), "DELETE FROM public.\"Personne\" p WHERE p.username = $1 AND p.password = $2", username, password)
	if err1 != nil {
		fmt.Println("failed to insert, err:", err1)
	}else{
		
		status = Status{Status: "OK"}

		json.NewEncoder(w).Encode(status)
	}
}


func delete_product(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	var jsonObj map[string]interface{}
	json.Unmarshal(body, &jsonObj)

	username := jsonObj["username"].(string)
	productId := jsonObj["productId"].(float64)


	conn, err := pgx.Connect(context.Background(), connUrl)
	if err != nil {
		log.Fatal("failed to connect to mcdb, err:", err)
	}
	var status Status

	_, err1 := conn.Exec(context.Background(), "DELETE FROM public.\"Produit\" p WHERE p.\"productId\" = $1 AND p.\"personId\" = (SELECT \"personId\" FROM \"Personne\" WHERE username = $2)", productId, username)
	if err1 != nil {
		fmt.Println("failed to insert, err:", err1)
	}else{
		
		status = Status{Status: "OK"}

		json.NewEncoder(w).Encode(status)
	}
}

func main() {

	http.HandleFunc("/add_user", add_user)
	http.HandleFunc("/add_guest", add_guest)
	http.HandleFunc("/add_product", add_product)
	http.HandleFunc("/add_comment", add_comment)
	http.HandleFunc("/add_to_cart", add_to_cart)
	http.HandleFunc("/report_product", report_product)
	http.HandleFunc("/buy_product", buy_product)
	http.HandleFunc("/fetch_users", fetch_users)
	http.HandleFunc("/fetch_comments", fetch_comments)
	http.HandleFunc("/fetch_products", fetch_products)
	http.HandleFunc("/fetch_product", fetch_product)
	http.HandleFunc("/fetch_reports", fetch_reports)
	http.HandleFunc("/fetch_payements", fetch_payements)
	http.HandleFunc("/login", login);
	http.HandleFunc("/fetch_user_products", fetch_products);
	http.HandleFunc("/fetch_user", fetch_user);
	// http.HandleFunc("/update_user", update_user);
	http.HandleFunc("/fetch_cart", fetch_cart);

	http.HandleFunc("/delete_from_cart", delete_from_cart);
	http.HandleFunc("/delete_user", delete_user);
	http.HandleFunc("/delete_product", delete_product);

	
	//think what other roots are needed and add them
	// make a small doc for rayane

	// http.HandleFunc("/", root);
	// http.HandleFunc("/", root);
	http.HandleFunc("/", root)
	fmt.Println("server listening on 8000")
	http.ListenAndServe(":8000", nil)
}
