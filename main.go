package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "github.com/gorilla/mux"
    "github.com/tkanos/gonfig"
    "database/sql"
    _ "github.com/lib/pq"
)

type Configuration struct {
    // Server
    Url               string
    Port              int
    // DateBase
    DB_Url            string
    DB_User           string
    DB_Pwd            string
    DB_Name           string
}

type ItemRow struct {
    id        int64
    name      string
    price     int64
}

func home(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome home!")
}

func helloworld(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello World!")
}

func (config *Configuration) getOneItem(w http.ResponseWriter, r *http.Request) {
    
    // Connect to the DB, panic if failed
    // "postgres://user:password@localhost/dbName?sslmode=disable"
    db_connect := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
                                config.DB_User,
                                config.DB_Pwd,
                                config.Url,
                                config.DB_Name)

    db, err := sql.Open("postgres", db_connect)
    if err != nil {
        fmt.Println(`Could not connect to db`)
        panic(err)
    }
    defer db.Close()

    // Start to Query
    itemID := mux.Vars(r)["id"]
    rows, err := db.Query(`SELECT * FROM items WHERE id=$1`, itemID)
    if err != nil {
        panic(err)
    }

    row := ItemRow{}

    for rows.Next() {
        rows.Scan(&row.id, &row.name, &row.price)
        output := fmt.Sprintf("{id:%d, name:%s, price:%d}", row.id, row.name, row.price)
        fmt.Fprintf(w, output)
    }

}

func (config *Configuration) deleteOneItem(w http.ResponseWriter, r *http.Request) {
    
    // Connect to the DB, panic if failed
    // "postgres://user:password@localhost/dbName?sslmode=disable"
    db_connect := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
                                config.DB_User,
                                config.DB_Pwd,
                                config.Url,
                                config.DB_Name)

    db, err := sql.Open("postgres", db_connect)
    if err != nil {
        fmt.Println(`Could not connect to db`)
        panic(err)
    }
    defer db.Close()

    // Start to Delete
    itemID := mux.Vars(r)["id"]
    query := "DELETE FROM items WHERE id=$1"
    rows, err := db.Exec(query, itemID)
    // fmt.Println(rows)
    if err != nil {
        panic(err)
    }

    count, _ := rows.RowsAffected()  
    if count ==1 {
        fmt.Fprintf(w, "success")
    }else{
        fmt.Fprintf(w, "fail")
    }

}

func (config *Configuration) createOneItem(w http.ResponseWriter, r *http.Request) {
    
    // Connect to the DB, panic if failed
    // "postgres://user:password@localhost/dbName?sslmode=disable"
    db_connect := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
                                config.DB_User,
                                config.DB_Pwd,
                                config.Url,
                                config.DB_Name)

    db, err := sql.Open("postgres", db_connect)
    if err != nil {
        fmt.Println(`Could not connect to db`)
        panic(err)
    }
    defer db.Close()

    // Start to Create
    name := mux.Vars(r)["name"]
    price := mux.Vars(r)["price"]
    query := "INSERT INTO items(name, price) VALUES($1, $2)"
    rows, err := db.Exec(query, name, price)
    // fmt.Println(rows)
    if err != nil {
        panic(err)
    }

    count, _ := rows.RowsAffected()  
    if count ==1 {
        fmt.Fprintf(w, "success")
    }else{
        fmt.Fprintf(w, "fail")
    }

}

func main() {

    config := Configuration{}
    err := gonfig.GetConf("./config/config.json", &config)

    if err != nil {
        os.Exit(0);
    }

    fmt.Printf("Using Url: %s  \n",config.Url);
    fmt.Printf("Using Port: %d \n",config.Port);

    // fmt.Printf("Using DB_Url: %s  \n",config.DB_Url);
    // fmt.Printf("Using DB_Port: %s \n",config.DB_User);
    // fmt.Printf("Using DB_Url: %s  \n",config.DB_Pwd);
    // fmt.Printf("Using DB_Port: %s \n",config.DB_Name);

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", home)
    router.HandleFunc("/helloworld", helloworld)

    router.HandleFunc("/items/{id}", config.getOneItem).Methods("GET")
    router.HandleFunc("/items/{name}/{price}", config.createOneItem).Methods("POST")
    router.HandleFunc("/items/{id}", config.deleteOneItem).Methods("DELETE")

    address := fmt.Sprintf("%s:%d", config.Url, config.Port) 
    log.Fatal(http.ListenAndServe(address, router))
}