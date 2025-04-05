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

func main() {

	http.HandleFunc("/add_user", add_user)
	http.HandleFunc("/add_guest", add_guest)
	http.HandleFunc("/add_product", add_product)
	http.HandleFunc("/add_comment", add_comment)
	// http.HandleFunc("/report_product", report_product);
	// http.HandleFunc("/buy_product", buy_product);
	// http.HandleFunc("/login", login);
	// http.HandleFunc("/fetch_products", fetch_products);
	// http.HandleFunc("/fetch_user_products", fetch_products);
	// http.HandleFunc("/fetch_users", fetch_users);
	// http.HandleFunc("/fetch_comments", fetch_comments);
	// http.HandleFunc("/fetch_reports", fetch_reports);
	// http.HandleFunc("/fetch_user", root);
	// http.HandleFunc("/update_user", root);
	// http.HandleFunc("/", root);
	// http.HandleFunc("/", root);
	// http.HandleFunc("/", root);
	http.HandleFunc("/", root)
	fmt.Println("server listening on 8000")
	http.ListenAndServe(":8000", nil)
}
