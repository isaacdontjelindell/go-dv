package main

import (
    "fmt"
    "os"
    "net"
    "time"
    "strconv"
    "encoding/json"
)

func main() {
    // get this station name and neighbor/cost
    configLines, err := readLines("conf.ini")
    if err != nil {
        fmt.Println("Error reading config file", err.Error())
        os.Exit(1)
    }

    station := configLines[0]
    fmt.Println("Station:", station)

    neighborCountStr := configLines[1]
    neighborCount, err := strconv.Atoi(neighborCountStr)
    if err != nil {
        fmt.Println("Error converting neighborCost to int.", err.Error())
        os.Exit(1)
    }
    fmt.Println("neighborCount:", neighborCount)


    neighbors := make([]Node, 0)
    for i := 2; i<neighborCount*2 + 1; i+=2 {
        name := configLines[i]
        costStr := configLines[i+1]

        cost, err := strconv.Atoi(costStr)
        if err != nil {
            fmt.Println("Error converting cost to int.", err.Error())
            os.Exit(1)
        }

        fmt.Printf("Neighbor: %s, Cost: %d\n", name, cost)

        n := Node{name, "", cost}

        neighbors = append(neighbors, n)
    }
    println("")

    // set up the threads
    quit := make(chan int)

    updateChan := make(chan []Update)

    go maintainRoutingTable(quit, updateChan, neighbors)
    go acceptUpdates(quit, updateChan)

    go testClient() // TODO remove

    <-quit // blocks
}

func maintainRoutingTable(quit chan int, updateChan chan []Update, initialTable []Node) {
    routingTable := make([]Node, 0)
    routingTable = append(routingTable, initialTable...)

    for {
        update := <-updateChan
        fmt.Println("[maintainRoutingTable] processing an update...", update)
    }

    quit <- -1
}

func acceptUpdates(quit chan int, updateChan chan []Update) {
    fmt.Println("[acceptUpdates] starting listener")

    // accept updates and pass them to maintainRoutingTable
    LISTEN_IP := net.ParseIP("127.0.0.1")
    LISTEN_PORT := 1337

    listenAddr := net.UDPAddr{LISTEN_IP, LISTEN_PORT, ""}
    listener, err := net.ListenUDP("udp", &listenAddr)
    if err != nil {
        fmt.Println("[acceptUpdates] error starting update listener!", err.Error())
        os.Exit(1)
    }

    for {
        data := make([]byte, 2048)
        n, from, err := listener.ReadFromUDP(data) //first param is number of bytes recieved
        if err != nil {
            fmt.Println("[acceptUpdates] Error accepting connection!", err.Error())
            return
        }

        fmt.Println("[acceptUpdates] recieved a packet from", from)

        // unmarshal update struct from data
        var update Update
        err = json.Unmarshal(data[:n], &update)
        if err != nil {
            fmt.Println("[acceptUpdates] error unmarshaling struct from JSON!", err.Error())
            os.Exit(1)
        }
        // TODO process update
    }

    quit <- -1
}


func testClient() {
    conn, err := net.Dial("udp",  "127.0.0.1:1337")
    if err != nil {
        fmt.Println("[testClient] Error dialing connection.", err.Error())
    }

    for {
        time.Sleep(time.Second * 2)

        // build a test update struct
        testRoutingTable := []Node{Node{"t1", "yoda", 3}, Node{"t2", "yoda", 5}}
        update := Update{testRoutingTable, "yoda"}

        u, err := json.Marshal(update)  // u is []byte
        if err != nil {
            fmt.Println("[testClient] error marshaling update to JSON", err.Error())
        }

        conn.Write(u)

    }
}

