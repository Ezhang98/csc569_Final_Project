// Group Members: Evan Zhang, Alexander Garcia and Christina Monahan
// Distributed System for Tuning Hyperparameters of Neural Networks
// PAXOS election, consensus,and recovery algorithms

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// content_types: prepare 0, promise 1, propose 2, accept 3, election 4
type msg struct {
	content_type int
	from         int
	to           int
	cmd          Command
}

// channels
type buffer struct {
	id  int
	buf map[int]chan msg
}

type Command struct {
	id   int
	name string
	arg1 string
}

type Node struct {
	id        int
	position  string
	stored    int
	commitLog []Command
	phone     *buffer
	hbtable   *[][]int64
	hb1       []chan [][]int64
	hb2       []chan [][]int64
	backup    *os.File
}

// System of channels for communication
func makeNetwork(length int) *buffer {
	phone := buffer{buf: make(map[int]chan msg, 0)}
	for i := 0; i < length; i++ {
		phone.buf[i] = make(chan msg, 1024)
	}
	return &phone
}

// sending channel
func (phone *buffer) send(m msg) {
	phone.buf[m.to] <- m
}

// receive channel
func (phone *buffer) receive(id int) msg {
	select {
	case reply := <-phone.buf[id]:
		return reply
	case <-time.After(100 * time.Millisecond):
		return msg{content_type: -1}
	}
}

// run proposer function to run consensus
func (n *Node) runProposer(numNodes int, input chan string) {
	select {
	case text := <-input:
		n = consensus(n, text, numNodes)
	default:
	}
	n.position = "acceptor"
}

// run Acceptor function - accepts proposals, sends promise, starts election
func (n *Node) runAcceptor(count int, numNodes int) int {
	lowest := 0
	count2 := 0
	if count == 20 {
		// run elections
		for i := 0; i < numNodes; i++ {
			if i == n.id {
				continue
			}
			var c Command
			m := msg{from: n.id, to: i, cmd: c, content_type: 4}
			n.phone.send(m)
			// accept request next value
		}
		count = 0
		lowest = n.id
		for {
			resp := n.phone.receive(n.id)
			count++
			// if message type 4, check from value
			if resp.content_type == 4 {
				count2++
				if resp.from < lowest {
					lowest = resp.from
				}
			}
			if count == 20 || count2 == numNodes-1 {
				break
			}
		}
		if lowest == n.id {
			// become proposer
			n.position = "proposer"
			return 0
		}
		return 0
	}
	// acceptor receives message
	m := n.phone.receive(n.id)
	if m.content_type == -1 {
		count++
		return count
	}

	c := m.cmd
	if m.content_type == 0 {
		val := n.stored
		var resp Command
		resp.name = "halt"
		if val < c.id {
			resp.id = c.id
			n.stored = c.id
		} else {
			resp.id = val + 1
		}
		m2 := msg{from: n.id, to: m.from, cmd: resp, content_type: 1}

		n.phone.send(m2)
	} else if m.content_type == 4 {
		count2++
		for i := 0; i < numNodes; i++ {
			if i == n.id {
				continue
			}
			var c Command
			m := msg{from: n.id, to: i, cmd: c, content_type: 4}
			n.phone.send(m)
			// accept request next value
		}
		count = 0
		lowest = n.id
		for {
			resp := n.phone.receive(n.id)
			count++
			// if message type 4, check from value
			if resp.content_type == 4 {
				count2++
				if resp.from < lowest {
					lowest = resp.from
				}
			}
			if count == 20 || count2 == numNodes-1 {
				break
			}
		}
		if lowest == n.id {
			// become proposer
			n.position = "proposer"
			return 0
		}
	}
	if m.content_type == 3 {
		n.commitToLog(c)
	}
	return 0
}

// commitLog saves to stable storage, and provides persistent state
func (n *Node) commitToLog(c Command) {
	if n.commitLog[c.id].name == "" {
		n.commitLog[c.id] = c
		logMsg := fmt.Sprintf("%d %s %s\n", c.id, c.name, c.arg1)
		_, err := n.backup.WriteString(logMsg)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// launches a new node with a heartbeat as a goRoutine
func (n *Node) launch(num int, chan1 chan string, input chan string, killhb chan string) {
	go heartbeat(n.hb1, n.hb2, n.id, n.hbtable, num, killhb)
	count := 0
	// fmt.Println(n.id, "started")
	for {
		if n.position == "proposer" {
			n.runProposer(num, input)
		} else if n.position == "acceptor" {
			count = n.runAcceptor(count, num)
		}
		select {
		case reply := <-chan1:
			if reply == "exit" {
				break
			}
		default:
		}
	}
	n.backup.Close()
	fmt.Println(n.id, "died")
}

// consensus function consistently stores data
func consensus(proposer *Node, command string, numNodes int) *Node {
	for {
		count := 0
		accepted := make([]bool, numNodes)
		largest := 0
		for j := 0; j < numNodes; j++ {
			if proposer.id == j {
				continue
			}
			// prepare request next value
			var c Command
			c.id = proposer.stored + 1
			if largest > c.id {
				c.id = largest
			} else {
				largest = c.id
			}
			c.name = "hash"
			c.arg1 = command
			m := msg{from: proposer.id, to: j, cmd: c, content_type: 0}
			proposer.phone.send(m)
			// accept request next value
			resp := proposer.phone.receive(proposer.id)
			cmd := resp.cmd
			if cmd.id == c.id {
				count++
				accepted[j] = true
			} else {
				if cmd.id > largest {
					largest = cmd.id
				}
				accepted[j] = false
			}
		}
		if count >= 4 {
			proposer.stored = largest
			break
		}
	}
	var c Command
	for j := 0; j < numNodes; j++ {
		// accept request next value
		if proposer.id == j {
			continue
		}
		c.id = proposer.stored
		c.name = "hash"
		c.arg1 = command
		m2 := msg{from: proposer.id, to: j, cmd: c, content_type: 3}
		proposer.phone.send(m2)
	}
	proposer.commitToLog(c)
	return proposer
}

//creates a new node
func createNode(id int, position string, stored int, cLog []Command, hb1chan []chan [][]int64, hb2chan []chan [][]int64) *Node {
	server := Node{id: id, position: position, stored: stored, commitLog: cLog}
	now := time.Now().Unix() // current local time
	hbtable := make([][]int64, 8)
	// initialize hbtable
	for i := 0; i < 8; i++ {
		hbtable[i] = make([]int64, 2)
		hbtable[i][0] = 0
		hbtable[i][1] = now
	}
	server.hbtable = &hbtable
	server.hb1 = hb1chan
	server.hb2 = hb2chan
	mode := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	filename := fmt.Sprintf("backup%d", id)
	server.backup, _ = os.OpenFile(filename, mode, 0644)
	return &server
}

func max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func main() {
	// check command line arguments
	numNodes := 0
	end := false
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run paxos.go numNodes")
		return
	}
	numNodes, _ = strconv.Atoi(os.Args[1])
	if numNodes < 3 {
		fmt.Println("numNodes needs to be at least 3")
		return
	}

	// initialize nodes
	nodes := make([]*Node, numNodes)

	var hb1chan = make([]chan [][]int64, 8) // heartbeat channels to neighbors for read
	var hb2chan = make([]chan [][]int64, 8)

	// initialize heartbeat tables
	for i := 0; i < 8; i++ {
		hb1chan[i] = make(chan [][]int64, 1024)
		hb2chan[i] = make(chan [][]int64, 1024)
	}

	for i := 0; i < numNodes; i++ {
		commitLog := make([]Command, 1024)
		nodes[i] = createNode(i, "acceptor", 0, commitLog, hb1chan, hb2chan)
	}

	phone := makeNetwork(len(nodes))

	//killNodesChan and killhb channels to be used with kill calls
	killNodesChan := make([]chan string, numNodes)
	killhb := make([]chan string, numNodes)
	input := make(chan string, 1024)
	// start all nodes
	for i := 0; i < numNodes; i++ {
		nodes[i].phone = phone
		killNodesChan[i] = make(chan string, 10)
		killhb[i] = make(chan string, 10)
		go nodes[i].launch(numNodes, killNodesChan[i], input, killhb[i])
	}

	// function to accept input from user and will kill goRoutines when exit is typed
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("> ")
		for scanner.Scan() {
			text := scanner.Text()
			if text == "exit" {
				for i := 0; i < numNodes; i++ {
					killNodesChan[i] <- "exit"
					killhb[i] <- "exit"
				}
				end = true
				break
			}
			input <- text
			time.Sleep(250 * time.Millisecond)
			fmt.Print("> ")
		}
	}()

	// killLists to simulate a hardware failure by killing nodes
	nodeKillList := make([][]int, 20)
	index := 0
	for i := 0; i < 20; i++ {
		index = i % numNodes
		nodeKillList[i] = make([]int, 2)
		nodeKillList[i][0] = index
		nodeKillList[i][1] = 10 * (i + 1)
	}
	counter := 0
	next := 0
	for {
		if end {
			break
		}
		if nodeKillList[next][1] == counter {
			// send kill to nkl[next][0]
			fmt.Println("\nKilling node", nodeKillList[next][0])
			fmt.Print("> ")
			killNodesChan[nodeKillList[next][0]] <- "exit"
			killhb[nodeKillList[next][0]] <- "exit"
		}
		// restart node at nkl[next][0] after 10 cycles
		if nodeKillList[next][1]+5 == counter {
			i := nodeKillList[next][0]
			commitLog, lastStored := catchUpCommands(i)
			nodes[i] = createNode(i, "acceptor", lastStored, commitLog, hb1chan, hb2chan)
			nodes[i].phone = phone
			go nodes[i].launch(numNodes, killNodesChan[i], input, killhb[i])
			fmt.Println("\nRestarting Node", i)
			fmt.Print("> ")
			next++
		}
		counter++
		if next == 20 {
			next = 0
		}
		time.Sleep(1 * time.Second)
	}
}

// catchUpCommands to restore a failed node to its previous state
func catchUpCommands(id int) ([]Command, int) {
	commands := make([]Command, 1024)
	stored := 0
	filename := fmt.Sprintf("backup%d", id)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	lines := strings.Split(string(content), "\n")

	for i := 0; i < len(lines)-1; i++ {
		x := strings.Split(lines[i], " ")
		if len(x) == 0 {
			continue
		}
		var c Command
		c.id, _ = strconv.Atoi(x[0])
		c.name = x[1]
		c.arg1 = x[2]
		commands[c.id] = c
		stored = c.id
	}
	return commands, stored
}

// update heartbeat table
func updateTable(index int, hbtableOG *[][]int64, counter int, hb1 []chan [][]int64, hb2 []chan [][]int64, numNodes int) {
	hbtable := *hbtableOG
	next := index + 1
	prev := index - 1
	if prev < 0 {
		prev = numNodes - 1
	}
	if next > numNodes-1 {
		next = 0
	}

	neighbor1 := make([][]int64, numNodes)
	neighbor2 := make([][]int64, numNodes)

	for i := 0; i < numNodes; i++ {
		neighbor1[i] = make([]int64, 2)
		neighbor1[i][0] = hbtable[i][0]
		neighbor1[i][1] = hbtable[i][1]
		neighbor2[i] = make([]int64, 2)
		neighbor2[i][0] = hbtable[i][0]
		neighbor2[i][1] = hbtable[i][1]

	}
	loop2 := true
	if counter%5 == 0 {
		for loop2 {
			select {
			case neighbor1_1 := <-hb1[next]:
				neighbor1 = neighbor1_1
			default:
				loop2 = false
			}
		}
		for i := 0; i < numNodes; i++ {
			(*hbtableOG)[i][0] = max(neighbor1[i][0], hbtable[i][0])
			(*hbtableOG)[i][1] = max(neighbor1[i][1], hbtable[i][1])
		}
		loop2 = true

		for loop2 {
			select {
			case neighbor2_1 := <-hb2[prev]:
				neighbor2 = neighbor2_1
			default:
				loop2 = false
			}
		}
		for i := 0; i < numNodes; i++ {
			(*hbtableOG)[i][0] = max(neighbor2[i][0], hbtable[i][0])
			(*hbtableOG)[i][1] = max(neighbor2[i][1], hbtable[i][1])
		}
	}

	now := time.Now().Unix() // current local time
	(*hbtableOG)[index][0] = hbtable[index][0] + 1
	(*hbtableOG)[index][1] = now

	// send table
	hb1[index] <- *hbtableOG
	hb2[index] <- *hbtableOG

}

// heartbeat function
func heartbeat(hb1 []chan [][]int64, hb2 []chan [][]int64, k int, hbtable *[][]int64, numNodes int, noPulse chan string) {
	counter := 0
	for {
		time.Sleep(100 * time.Millisecond)
		updateTable(k, hbtable, counter, hb1, hb2, numNodes)
		counter++
		select {
		case reply := <-noPulse:
			if reply == "exit" {
				break
			}
		}
	}
}
