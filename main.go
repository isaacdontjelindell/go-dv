package main

import (
    "fmt"
    "os"
    "net"
    "time"
)

type Node struct {
    name string
    route string
    totalCost int // int cost + (cost to route)
}

func updateNodeRoute(n Node, route string) Node {
    n.route = route
    return n
}

func updateNodeCost(n Node, cost int) Node {
    n.totalCost = cost
    return n
}


type Update struct {
    routingTable []Node
    from string
}


func main() {
    // get this station name and neighbor/cost
    print("Enter this station name: ")
    var station string
    fmt.Scanf("%s", &station)
    println("")

    neighbors := make([]Node, 0)
    for i := 0; i < 3; i++ {
        print("Enter the neighbor's name: ")
        var name string
        fmt.Scanf("%s", &name)

        print("Enter the cost: ")
        var cost int
        fmt.Scanf("%d", &cost)

        n := Node{name, "", cost}

        neighbors = append(neighbors, n)
    }
    println("")

    println("Station:", station)
    println("Neighbors:", neighbors)


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
        _, from, err := listener.ReadFromUDP(data) //first param is number of bytes recieved
        if err != nil {
            fmt.Println("[acceptUpdates] Error accepting connection!", err.Error())
            return
        }
        fmt.Println("[acceptUpdates] recieved an update...", "FROM:", from)
        fmt.Println(string(data))
        // TODO process data
    }

    quit <- -1
}


func testClient() {

    conn, err := net.Dial("udp",  "127.0.0.1:1337")
    if err != nil {
        fmt.Println("Test client broke", err.Error())
    }
    for {
        msg := []byte("test message")
        _, err = conn.Write(msg)
        if err != nil {
            fmt.Println(err.Error())
        }
        time.Sleep(time.Second * 2)
    }
}

