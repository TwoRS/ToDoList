package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync/atomic" 
	"github.com/jinzhu/gorm"
    _ "github.com/lib/pq"
    "log"
    "github.com/jeffail/gabs"
 
)

func SendMessageToSock(sockId uint64, msg string) {
  	ActiveClients[sockId].websocket.WriteMessage(websocket.TextMessage, []byte(msg)) //отправляем
}

func SendMessageToAll(msg  string) {
  for _, Client := range ActiveClients {
    ActiveClients[Client.Socket_Id].websocket.WriteMessage(websocket.TextMessage, []byte(msg))
  }
}



func SockServer(w http.ResponseWriter, r *http.Request) {
    
    conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
    if _, ok := err.(websocket.HandshakeError); ok {
      http.Error(w, "Not a websocket handshake", 400)
      return
    }

    db, err := gorm.Open("postgres", "user=Roman password=Roman dbname=DB1  sslmode=disable")  
    if err != nil {
        log.Fatalf("error: %v\n", err)
    }
    //db.CreateTable(&MyList{})  //создаём табл
	var allToDo []MyList

    Soc_Id := atomic.AddUint64(Max_SocId, 1)
    sockCli := ClientConn{conn, Soc_Id}
    ActiveClients[Soc_Id] = sockCli

    db.DB()
    db.Order("ID").Find(&allToDo)
    jsonObj:= gabs.New()
	jsonObj.Set("All_List", "Command")
	jsonObj.Set(&allToDo, "Data")
 	SendMessageToSock(Soc_Id,jsonObj.String())
    
    for { 
	        _, msg, err := conn.ReadMessage()
	        if err != nil {  
	            delete(ActiveClients,Soc_Id) 
	            return
	        }
	        var Command string
	        var ok bool
			jsonParsed, err_ := gabs.ParseJSON(msg)
		    if err_ != nil {Davay_Do_Svidaniya(Soc_Id,conn);return}
			Command, ok= jsonParsed.Path("Command").Data().(string)
			if ok == false {Davay_Do_Svidaniya(Soc_Id,conn);return} 

		  	switch Command {

				case "Todo_Add":
				  	var Name string
					Name, ok= jsonParsed.Path("Data.Name").Data().(string)
					if ok == false {Davay_Do_Svidaniya(Soc_Id,conn);return} 

					var todo MyList
					todo.Name=Name
					todo.IsTrue=false
					db.Create(&todo)

					jsonObj := gabs.New()
					jsonObj.Set("New_Todo", "Command")
					jsonObj.Set(&todo.ID, "Data","ID")
					jsonObj.Set(Name, "Data","Name")

				 	SendMessageToAll(jsonObj.String())

				case "Todo_Delete":
				    var ID float64
					ID, ok= jsonParsed.Path("Data.ID").Data().(float64)
					if ok == false {Davay_Do_Svidaniya(Soc_Id,conn);return} 
					db.Where("ID = ?", ID).Delete(&MyList{})

					var todo MyList
					todo.IsTrue=false

					jsonObj := gabs.New()
					jsonObj.Set("Todo_Deleted", "Command")
					jsonObj.Set(ID, "Data","ID")

				 	SendMessageToAll(jsonObj.String())


				case "Todo_Check":
				    var ID float64
				   	var Value bool
				   	var value_ string
					ID, ok= jsonParsed.Path("Data.ID").Data().(float64)
					if ok == false {Davay_Do_Svidaniya(Soc_Id,conn);return} 
					value_, ok= jsonParsed.Path("Data.IsTrue").Data().(string)
					Value=ParseBool(value_) 
					if ok == false {Davay_Do_Svidaniya(Soc_Id,conn);return} 
 				
					var my_list MyList
					db.First(&my_list, uint64(ID))
					my_list.IsTrue = Value
					db.Save(&my_list) 

					jsonObj := gabs.New()
					jsonObj.Set("Todo_Checked", "Command")
					jsonObj.Set(ID, "Data","ID")
					jsonObj.Set(value_, "Data","IsTrue")

				 	 SendMessageToAll(jsonObj.String())


				case "RemoveAll":
					var todo MyList
					db.Delete(&todo)

    				db.Order("ID").Find(&allToDo)
    				jsonObj:= gabs.New()
					jsonObj.Set("All_List", "Command")
					jsonObj.Set(&allToDo, "Data") 
				 	SendMessageToAll(jsonObj.String())  
			}
    	}

    conn.Close();     
}

 
//Если прислал абрукадабру- выкидываем
func Davay_Do_Svidaniya(Soc_Id uint64,conn *websocket.Conn){
	delete(ActiveClients,Soc_Id) 
	conn.Close(); 
}

//"github.com/jeffail/gabs"  не работает с логическими значениями. Поэтому добавил перевод
func ParseBool(str string) (value bool) {
    switch str {
    	case "1","true", "TRUE", "True":
    			return true
    	case "0","false", "FALSE", "False":
    			return false
    }
    return false
}

