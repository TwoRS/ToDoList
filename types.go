package main

import (
    "github.com/jinzhu/gorm"
    "github.com/gorilla/websocket"
    //"encoding/json"
    //"github.com/jeffail/gabs"
   // "time"
)



type MyList struct {
    gorm.Model
    //Id uint64 //`sql:"AUTO_INCREMENT" gorm:"primary_key"`
    Name  string  
    IsTrue bool 
}

/*type Model struct {
    Id        uint `gorm:"primary_key"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time
}*/

type ClientConn struct {
    websocket *websocket.Conn
    Socket_Id uint64
}

type struct_json struct { 
    Command string
    Data interface{} 
}


/*
type struct_json_Response struct { 
    Command string
    Data string
}
*/
