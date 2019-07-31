package server

import (
    "../common"
    "encoding/json"
    "github.com/gorilla/websocket"
    "./methods"
    "net/http"
    "log"
)

func HandleServe(w http.ResponseWriter, r *http.Request) {
    // Define check to allow websocket connection
    upgrader := websocket.Upgrader{
        CheckOrigin: func(r *http.Request) (bool) {
            return true
        },
    }
    // Convert request to websocket and quick if error
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    // Always close websocket when done
    defer ws.Close()

    for {
        // Read Message from Websocket
        _, raw, err := ws.ReadMessage()
        if err != nil {
            log.Println(err)
            break
        }

        // Create Basic Reply Message
        var reply common.OutgoingMessage
        reply.Type = "result"
        reply.Result = 500

        // Check if is Valid JSON and convert to incommingMessage structure
        var msg common.IncommingMessage
        err = json.Unmarshal(raw, &msg)
        if err != nil {
            log.Println(err)
        }

        // Assume the same ID
        reply.Id = msg.Id

        // If has a valid message ID
        if msg.Type != "" && msg.Id != 0 {
            // Handle Different Types of Messages
            switch msg.Type {
            case "method":
                // If method name exists, use it
                if m, ok := methods.Get(msg.Method); ok {
                    reply.Obj, reply.Result = m(msg.Obj)
                } else {
                    reply.Result = 400
                }
            case "lookup":
                log.Println("Lookup")
                reply.Result = 200
            case "subscribe":
                log.Println("Subscribe")
                reply.Result = 202
            }
        } else {
            reply.Result = 400
        }

        // Send Reply
        replyRaw, _ := json.Marshal(&reply)
        ws.WriteMessage(websocket.TextMessage, replyRaw)
    }
}