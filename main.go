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
    DbUrl            string
    DbUser           string
    DbPwd            string
    DbName           string
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
    dbConnect := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
                                config.DbUser,
                                config.DbPwd,
                                config.DbUrl,
                                config.DbName)

    db, err := sql.Open("postgres", dbConnect)
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
        return
    }

    fmt.Fprintf(w, `{id : -1, msg : "not found"}`)

}

func (config *Configuration) deleteOneItem(w http.ResponseWriter, r *http.Request) {
    
    // Connect to the DB, panic if failed
    // "postgres://user:password@localhost/dbName?sslmode=disable"
    dbConnect := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
                                config.DbUser,
                                config.DbPwd,
                                config.DbUrl,
                                config.DbName)

    db, err := sql.Open("postgres", dbConnect)
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
        fmt.Fprintf(w, "{msg : success}")
    }else{
        fmt.Fprintf(w, "{msg : fail}")
    }

}

func (config *Configuration) createOneItem(w http.ResponseWriter, r *http.Request) {
    
    // Connect to the DB, panic if failed
    // "postgres://user:password@localhost/dbName?sslmode=disable"
    dbConnect := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
                                config.DbUser,
                                config.DbPwd,
                                config.DbUrl,
                                config.DbName)

    db, err := sql.Open("postgres", dbConnect)
    if err != nil {
        fmt.Println(`Could not connect to db`)
        panic(err)
    }
    defer db.Close()

    // Start to Create
    returnId := -1
    name := mux.Vars(r)["name"]
    price := mux.Vars(r)["price"]
    query := "INSERT INTO items(name, price) VALUES($1, $2) RETURNING id"
    err = db.QueryRow(query, name, price).Scan(&returnId)

    if err != nil {
        fmt.Fprintf(w, `{id : -1, msg : "item already exists"}`)
        // panic(err)
    }else{
        result := fmt.Sprintf("{id : %d}", returnId)
        fmt.Fprintf(w, result)
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

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", home)
    router.HandleFunc("/helloworld", helloworld)

    router.HandleFunc("/items/{id}", config.getOneItem).Methods("GET")
    router.HandleFunc("/items/{name}/{price}", config.createOneItem).Methods("POST")
    router.HandleFunc("/items/{id}", config.deleteOneItem).Methods("DELETE")

    address := fmt.Sprintf("%s:%d", config.Url, config.Port) 
    log.Fatal(http.ListenAndServe(address, router))
}