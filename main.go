package main

import (
    "fmt"
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
    // accept updates and pass them to maintainRoutingTable
    for {
        time.Sleep(time.Second)
        fmt.Println("[acceptUpdates] recieved an update ")
    }

    quit <- -1
}


