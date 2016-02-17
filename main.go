package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "io/ioutil"
    //"sync/atomic"
    //"encoding/json"
    "fmt"
    "github.com/jinzhu/gorm"
    _ "github.com/lib/pq"
)

var Max_SocId *uint64 = new(uint64)

var ActiveClients = make(map[uint64 ] ClientConn)
var CacheIndexHtml=""



func main() {
    
    db, err := gorm.Open("postgres", "user=Roman password=Roman dbname=DB1  sslmode=disable")  
    if err != nil {
        log.Fatalf("error: %v\n", err)
    }
    db.DB()
    db.DB().Ping()
    db.DB().SetMaxIdleConns(10)
    db.DB().SetMaxOpenConns(100)
    db.SingularTable(true)

    //db.CreateTable(&MyList{})
    //db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&MyList{})
    //db.DropTable(&MyList{})
    //db.Model(&MyList{}).ModifyColumn("description", "text")
    //db.Model(&MyList{}).DropColumn("description")
    db.AutoMigrate(&MyList{})
    //db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&MyList{})
    //db.AutoMigrate(&MyList{})




    //просто кешируем страничку. Можно и без этого)))
    CacheIndexHtml_, err := ioutil.ReadFile("Index.html")
    if err != nil {
        return
    }
    CacheIndexHtml= string(CacheIndexHtml_)  

    
    router := mux.NewRouter()
    static := http.StripPrefix("/static/", http.FileServer(http.Dir("./files/")))
    router.PathPrefix("/static/").Handler(static)
    router.HandleFunc("/", Mainhandler) 
    http.HandleFunc("/websocket", SockServer) 

    http.Handle("/", router) 
    http.ListenAndServe(":5000", nil) 
}




func Mainhandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-type", "text/html; charset=utf-8")
    fmt.Fprintln(w, CacheIndexHtml)
}







