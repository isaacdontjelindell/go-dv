package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
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

	neighbors := make(map[string]Node)
	for i := 2; i < neighborCount*2+1; i += 2 {
		name := configLines[i]
		costStr := configLines[i+1]

		cost, err := strconv.Atoi(costStr)
		if err != nil {
			fmt.Println("Error converting cost to int.", err.Error())
			os.Exit(1)
		}

		n := Node{name, "_self", cost}

		neighbors[name] = n
	}
	println("")

	// the initial routing table is just the neighbor list
	routingTable := RoutingTable{neighbors, station}

	// set up the channels
	quit := make(chan int)
	updateChan := make(chan Update)               // handle incoming updates
	outgoingUpdateChan := make(chan RoutingTable) // handle sending RoutingTable to neighbors

	// set up the threads
	go maintainRoutingTable(quit, updateChan, outgoingUpdateChan, routingTable, neighbors)
	go acceptUpdates(quit, updateChan)
	go sendUpdates(quit, outgoingUpdateChan, neighbors)
	outgoingUpdateChan <- routingTable

	//go testClient() // TODO remove

	<-quit // blocks to keep main thread alive
}

func maintainRoutingTable(quit chan int, updateChan chan Update, outgoingUpdateChan chan RoutingTable, routingTable RoutingTable, neighbors map[string]Node) {
	fmt.Println("[maintainRoutingTable] Initial routing table:")
	fmt.Println(routingTable)

	for {
		update := <-updateChan // wait for updates

		from := update.From

		fmt.Printf("[maintainRoutingTable] processing an update. From: %s\n", from)
		updated := false // keep track if this update caused changes in the routing table

		_, ok := neighbors[from]
		if ok != true {
			fmt.Printf("[maintainRoutingTable] got an update from a stranger!\n")
			continue
		}

		for name, newNode := range update.RoutingTable {
			// ignore entry for this host in recieved routing table
			if name == routingTable.Self {
				continue
			} else {
				routeCost := neighbors[from].Cost

				newName := newNode.Name
				newCost := newNode.Cost + routeCost

				// check if newNode is in our current routing table
				node, present := routingTable.Table[newName]
				if present {
					// if it is, check cost
					cost := node.Cost // what we have now

					if newCost < cost {
						node.Cost = newCost
						node.Route = from
						routingTable.Table[newName] = node // update table w/ updated node
						updated = true
					} else {
						// if the cost we got in the update is more, ignore it
					}
				} else {
					// if newNode isn't in the current routing table, add it
					newNode.Cost = newCost
					newNode.Route = from
					routingTable.Table[newName] = newNode // update table w/ newly discovered station
					updated = true
				}

			}
		}

		if updated {
			fmt.Printf("[maintainRoutingTable] Updated routing table:\n")
			fmt.Println(routingTable)

			// send an update to neighbors
			outgoingUpdateChan <- routingTable

		} else {
			fmt.Printf("[maintainRoutingTable] No changes due to update from %s\n\n", update.From)
		}
	}

	quit <- -1
}

func sendUpdates(quit chan int, updateChan chan RoutingTable, neighbors map[string]Node) {
	fmt.Println("[sendUpdates] starting sendUpdates")

	connections := make(map[string]net.Conn)
	for name, _ := range neighbors {
		conn, err := net.Dial("udp", name+":1337")
		if err != nil {
			fmt.Println("[sendUpdates] Error dialing connection.", err.Error())
		}
		connections[name] = conn
	}

	routingTable := <-updateChan // block the first time, waiting for inital routing table
	for {
		select {
		// grab an updated routing table if one exists
		case routingTable = <-updateChan: // non-blocking
		// send the update to all neighbors
		default:
			fmt.Println("[sendUpdates] sending (possibly) updated routing table to neighbors")
			fmt.Println(routingTable)
			update := Update{routingTable.Table, routingTable.Self}

			u, err := json.Marshal(update)
			if err != nil {
				fmt.Println("[sendUpdates] error marshaling update to JSON", err.Error())
			}

			for _, conn := range connections {
				conn.Write(u)
			}

			time.Sleep(time.Second * 3)
		}

	}

	quit <- -1
}

func acceptUpdates(quit chan int, updateChan chan Update) {
	fmt.Println("[acceptUpdates] starting listener")

	// accept updates and pass them to maintainRoutingTable
	LISTEN_IP := net.ParseIP("0.0.0.0")
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

		updateChan <- update // pass the update object to maintainRoutingTable()
	}

	quit <- -1
}
