package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

type Node struct {
	Name  string
	Route string
	Cost  int
}

func (n Node) String() string {
	return fmt.Sprintf("Name: %s, Route: %s, Cost: %d", n.Name, n.Route, n.Cost)
}

type Update struct {
	RoutingTable map[string]Node
	From         string
}

func (u Update) String() string {
	var s string
	s = fmt.Sprintf("From: %s\nRouting Table:\n", u.From)
	for _, node := range u.RoutingTable {
		s += fmt.Sprint("  ")
		s += fmt.Sprintln(node)
	}
	return s
}

type RoutingTable struct {
	Table map[string]Node
	Self  string
}

func (r RoutingTable) String() string {
	var s string
	s = fmt.Sprintf("Routing table for: %s\n", r.Self)

	for _, val := range r.Table {
		s += fmt.Sprintf("  %s\n", val)
	}
	return s
}

func testClient() {
	conn, err := net.Dial("udp", "127.0.0.1:1337")
	if err != nil {
		fmt.Println("[testClient] Error dialing connection.", err.Error())
	}

	for {
		time.Sleep(time.Second * 2)

		// build a test update struct
		testRoutingTable := make(map[string]Node)
		testRoutingTable["t1"] = Node{"t1", "bob", 3}
		testRoutingTable["t2"] = Node{"t2", "joe", 5}
		testRoutingTable["t3"] = Node{"t3", "dan", 6}

		update := Update{testRoutingTable, "yoda"}

		u, err := json.Marshal(update) // u is []byte
		if err != nil {
			fmt.Println("[testClient] error marshaling update to JSON", err.Error())
		}

		conn.Write(u)
	}
}
