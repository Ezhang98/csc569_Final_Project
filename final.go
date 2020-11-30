// Group Members: Evan Zhang, Alexander Garcia and Christina Monahan
// Distributed System for Tuning Hyperparameters of Neural Networks
// PAXOS election, consensus, and recovery algorithm

package main

import (
	// "encoding/json"
	"fmt"
	// "hash/fnv"
	// "io/ioutil"
	"log"
	"os"
	// "path/filepath"
	// "plugin"
	// "sort"
	"reflect"
	"strconv"
	"strings"
	"time"
	"bufio"
	"github.com/dathoangnd/gonet"
	"encoding/csv"
)

// number of workers
var numWorkers int

//struct to organize data into the master function
type MasterData struct {
	id 					int
	numWorkers			int					
	request				[]chan ModelConfig
	models				[] ModelConfig
	replies				[]chan string
	corpses				chan []bool
	working				[]string
	finished			[]string
	toShadowMasters     []chan string
	log 				[]string
	hb1					[]chan [][]int64
	hb2 				[]chan [][]int64
	test 				[][][]float64
	training			[][][]float64
}

type UIWindow struct {
	TrainData  string
	TestData   string
	ModelCount int
	Models     []ModelConfig
}

type ModelConfig struct {
	ModelID int
	Name    string
	Model1Params
	NeuralNet
	Model3Params
}

type Model1Params struct {
	Activation int
	Nodes      int
}

type NeuralNet struct {
	inputNodes   	int
	numHiddenLayers 	int
	outputLayer		int
	numEpochs		int
	learningRate	float64
	momentum		float64
}

type Model3Params struct {
	Trees    int
	MaxDepth int
}

func main() {
	//timer := time.Now()
	
	// check command line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run final.go <number of nodes>")
		return
	}

	// launches master and shadow master nodes
	launchServers(os.Args[1])

	// customize number of nodes to run a system on? on UI
	//fmt.Printf("\nRuntime: %.5f seconds\n", time.Since(timer).Seconds())
}

// Launches nodes and creates MasterData structure
func launchServers(userInput string) {
	numWorkers, _ = strconv.Atoi(userInput)

	var mrData MasterData
	mrData.numWorkers = numWorkers
	mrData.request = make([]chan string, mrData.numWorkers)		// worker <- master : for task assignment
	mrData.replies = make([]chan string, mrData.numWorkers)		// master <- worker : worker reply for task completion
	mrData.corpses = make(chan []bool, mrData.numWorkers)		// master <- heartbeat : workers that have died
	mrData.working = make([]string, mrData.numWorkers)			// which tasks assigned to which workers
	mrData.finished = make([]string, len(mrData.models))		// which models have completed
	
	var hb1 = make([]chan [][]int64, mrData.numWorkers+3)		// heartbeat channels to neighbors for read
	var hb2 = make([]chan [][]int64, mrData.numWorkers+3)		// heartbeat channels to neighbors for write
	killMaster := make(chan string, 10)							// channel to kill Master to verify replication and recovery from log

	// initialize heartbeat tables
	for i := 0; i < mrData.numWorkers+3; i++ {
		hb1[i] = make(chan [][]int64, 1024)
		hb2[i] = make(chan [][]int64, 1024)
	}
	mrData.hb1 = hb1
	mrData.hb2 = hb2

	// initialize worker channels
	for k := 0; k < mrData.numWorkers; k++ {
		mrData.request[k] = make(chan string)
		mrData.replies[k] = make(chan string)
		mrData.working[k] = ""	
	}
	
	// initialize shadowMasters
	numShadowMasters := 2
	// shadowMaster <- master : replication
	mrData.toShadowMasters = make([]chan string, numShadowMasters)	
	for j := 0; j < numShadowMasters; j++ {
		mrData.toShadowMasters[j] = make(chan string, 10)
	}

	// start nodes
	masterlog := make([]string, 0)
	go master(mrData, hb1, hb2, masterlog, killMaster)
	go shadowMaster(mrData.toShadowMasters[0], hb1, hb2, mrData.numWorkers, mrData.numWorkers + 1, mrData, killMaster)
	go shadowMaster(mrData.toShadowMasters[1], hb1, hb2, mrData.numWorkers, mrData.numWorkers + 2, mrData, killMaster)

	// wait here until "q" is entered from the command line
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "q"{
			break
		}
	}
}

// Master node
func master(mrData MasterData, hb1 []chan [][]int64, hb2 []chan [][]int64, log []string, killMaster chan string) {
	// initialize variables
	currentStep := "step start"
	mrData.log = log
	killHB := make(chan string, 10)
	message := ""

	// Master Recovery: resume step after Master failure
	for i := 0; i < len(log); i++{
		if log[i][0:4] == "step"{
			currentStep = log[i]
		}
	}

	// Run master heartbeat
	go masterHeartbeat(hb1, hb2, mrData.numWorkers, mrData.corpses, killHB, killMaster, mrData)
	
	// fmt.Println("Master running: ", currentStep) // print in UI
	for {
		if currentStep == "step start" {
			// If master died partway through launching workers, find the last logged k and start from there
			k := 0
			if len(log) > 0 {
				last := strings.Split(log[len(log) - 1], " ")
				if last[0] == "launch" {
					k, _ = strconv.Atoi(last[2])
				}
			}
			for ; k < mrData.numWorkers; k++ {
				go worker(mrData.request[k], mrData.replies[k], hb1, hb2, k)
				// message = fmt.Sprintf("launch worker %s", k)
				mrData.log = append(mrData.log, currentStep)	//appends launch worker step to log
				mrData.toShadowMasters[0] <- message			//sends launch worker message to first Shadow Master channel
				mrData.toShadowMasters[1] <- message			//sends launch worker message to second Shadow Master channel
			}
			currentStep = "step working"
			mrData.log = append(mrData.log, currentStep)	//appends load step to log
			mrData.toShadowMasters[0] <- currentStep		//sends load message to first Shadow Master channel
			mrData.toShadowMasters[1] <- currentStep		//sends load message to second Shadow Master channel
			killHB <- "die"
			return
			
		}  else if currentStep == "step working" {
			// manage the distributeTasks step 
			trainpath := "../datasets/mnist_train.csv"
			trainpath := "../datasets/mnist_test.csv"
			mrData.training = parseCSV(trainpath)
			mrData.test = parseCSV(testpath)
			mrData = distributeTasks(mrData)
			currentStep = "step cleanup"
			mrData.log = append(mrData.log, currentStep)	//appends master distributeTasks step to log
			mrData.toShadowMasters[0] <- currentStep		//sends distributeTasks message to first Shadow Master channel
			mrData.toShadowMasters[1] <- currentStep		//sends distributeTasks message to second Shadow Master channel
			killHB <- "die"
			return

		} else if currentStep == "step cleanup" {
			// cleanup workers who should now be done with all tasks
			mrData = cleanup(mrData)
			currentStep = "step end"
			mrData.log = append(mrData.log, currentStep)	//appends end message to log
			mrData.toShadowMasters[0] <- currentStep		//sends end message to first Shadow Master channel
			mrData.toShadowMasters[1] <- currentStep		//sends end message to second Shadow Master channel
			killHB <- "die"
			return
		} else {
			break
		}
	}
	killHB <- "die"
	// fmt.Println("Running master has died.")
}


func distributeTasks(mrData MasterData) MasterData {
	count := 0
	modelNumber := 0
	loop := true

	fmt.Println("Distributing Tasks Started...")
	for loop {
		for i := 0; i < mrData.numWorkers; i++ {
			// checks for available workers
			if mrData.working[i] == "" {
				for j := 0; j < len(mrData.models); j++ {
					if mrData.finished[j] == "not started" {
						modelNumber = j
						mrData.finished[j] = "started"
						mrData.request[i] <- ("m_" + strconv.Itoa(modelNumber))
						mrData.contents[i] <- (mrData.models[modelNumber])
						mrData.working[i] = strconv.Itoa(modelNumber)
						break
					}
				}
			}

			// checks for replies and dead workers
			select {
			case message := <-mrData.replies[i]:
				replied := strings.Split(message, "_")
				workerID, _ := strconv.Atoi(replied[1])
				mrData.working[workerID] = ""
				modelID, _ := strconv.Atoi(replied[0])
				mrData.finished[modelID] = "finished"
				count++
			case coffins := <-mrData.corpses:
				for j := 0; j < mrData.numWorkers; j++ {
					if coffins[j] == true {
						mrData.hb1[j] = make(chan [][]int64, numWorkers+3)
						mrData.hb2[j] = make(chan [][]int64, numWorkers+3)
						mrData.request[j] = make(chan string)
						mrData.contents[j] = make(chan string)
						mrData.replies[j] = make(chan string)
						go worker(mrData.request[j], mrData.contents[j], mrData.replies[j], mrData.hb1, mrData.hb2, j)
						tempModelID, _ := strconv.Atoi(mrData.working[j])
						mrData.finished[tempModelID] = "not started"
						mrData.working[j] = ""
						coffins[j] = false
					}
				}
			default:
			}
		}
		
		// checks that all models have completed
		if count >= mrData.modelNumber {
			check := true
			
			for a := 0; a < mrData.numWorkers; a++ {
				if mrData.working[a] != "" {
					check = false
				}
			}
			if check == true {
				loop = false
			}
		}
	}
	fmt.Println("\nDistributing Finished")
	return mrData
}


func worker(master chan string, content chan string, reply chan string,  hb1 []chan [][]int64, hb2 []chan [][]int64, k int) {
	go heartbeat(hb1, hb2, k)
	task := ""
	for {
		// read task from channel
		task = <-master
		tasks := strings.Split(task, "_")
		if tasks[0] == "end" {
			reply <- tasks[2]
			return
		}
		if tasks[0] == "m" {
			x := <-content
			distributeTasks(x)
		} 
		reply <- tasks[1] + "_" + tasks[2]
	}
}

// Shuts down worker nodes
func cleanup(mrData MasterData) MasterData {
	// Check that all workers shutdown
	for i := 0; i < mrData.numWorkers; i++ {
		mrData.request[i] <- "end_0_" + strconv.Itoa(i)
		msg := <-mrData.replies[i]
		num, _ := strconv.Atoi(msg)
		mrData.working[num] = ""
	}
	return mrData
}

func runModelType(model ModelConfig){

	switch model.Name{
	case "neuralnet":
		runNeuralNet(x)
	default:
	}
}

// Helper function to find max
func max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func parseCSV(path string) [][][]float64{
	data := make([][][]float64, 0)

	// Open the file
	csvfile, err := os.Open(path)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	r := csv.NewReader(csvfile)
	index := 0

	// Parse Data
	for {
		// Read each record from csv
		record, err := r.Read()
		if index != 0{
			floatarr := make([]float64, len(record) + 1)
			expected := make([]float64, 1)
			for i := 0; i < len(record); i++ {
				if s, err := strconv.ParseFloat(record[i], 64); err == nil {
					if i == 0{
						expected[0] = s
					}else{
						floatarr[i] = s
					}
					// fmt.Println(s)
				}
			}
			if len(floatarr[1:]) == 0{
				break
			}
			one_entry := [][]float64{floatarr[1:], expected}
			
			data = append(data, one_entry)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			// fmt.Println(index)
		}
		index++
	}
	// fmt.Println(data)
	return data
}

func runNeuralNet(training [][][]float64, test [][][]float64, inputNodes int, hiddenLayers []int, outputLayer int, numEpochs int, learningRate float64, momentum float64){
	// Create a neural network
	// 2 nodes in the input layer
	// 2 hidden layers with 4 nodes each []int{4, 4}
	// 1 node in the output layer
	// The problem is classification, not regression

	// initialize hiddenlayers
	nn := gonet.New(inputNodes, hiddenLayers, outputLayer, false)

	// Train the network
	// Run for 3000 epochs
	// The learning rate is 0.4 and the momentum factor is 0.2
	// Enable debug mode to log learning error every 1000 iterations
	nn.Train(training, numEpochs, learningRate, momentum, true)

	// Predict
	// testInput := []float64{1, 0}
	// fmt.Printf("test input: %f  %f => %f\n", testInput[0], testInput[1], nn.Predict(testInput)[0])
	// // Save the model
	// nn.Save("model.json")
	// // Load the model
	// nn2, err := gonet.Load("model.json")
	// if err != nil {
	// 	log.Fatal("Load model failed.")
	// }
	// fmt.Printf("%f XOR %f => %f\n", testInput[0], testInput[1], nn2.Predict(testInput)[0])
	// // 1.000000 XOR 0.000000 => 0.943074
}

// Update heartbeat tables for Master, 2 ShadowMasters and 8 workers
func updateTable(index int, hbtable [][]int64, counter int, hb1 []chan [][]int64, hb2 []chan [][]int64) [][]int64 {
	next := index + 1
	prev := index - 1
	if prev < 0 {
		prev = numWorkers+2
	}
	if next > numWorkers+2 {
		next = 0
	}

	temp := make([][]int64, numWorkers+3)
	neighbor1 := make([][]int64, numWorkers+3)
	neighbor2 := make([][]int64, numWorkers+3)

	for i := 0; i < numWorkers+3; i++ {
		neighbor1[i] = make([]int64, 2)
		neighbor1[i][0] = hbtable[i][0]
		neighbor1[i][1] = hbtable[i][1]
		neighbor2[i] = make([]int64, 2)
		neighbor2[i][0] = hbtable[i][0]
		neighbor2[i][1] = hbtable[i][1]
		temp[i] = make([]int64, 2)
		temp[i][0] = hbtable[i][0]
		temp[i][1] = hbtable[i][1]
	}
	loop2 := true
	if counter % 1 == 0 {
		for loop2 {
			select {
				case neighbor1_1 := <-hb1[next]:
					neighbor1 = neighbor1_1
				default:
					loop2 = false
			}
		}
		for i := 0; i < numWorkers+3; i++ {
			temp[i][0] = max(neighbor1[i][0], hbtable[i][0])
			temp[i][1] = max(neighbor1[i][1], hbtable[i][1])
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
		for i := 0; i < numWorkers+3; i++ {
			temp[i][0] = max(neighbor2[i][0], hbtable[i][0])
			temp[i][1] = max(neighbor2[i][1], hbtable[i][1])
		}
	}
	now := time.Now().Unix() // current local time
	temp[index][0] = hbtable[index][0] + 1
	temp[index][1] = now
	// send table
	hb1[index] <- temp
	hb2[index] <- temp
	return temp
}

// Heartbeat function for all workers
func heartbeat(hb1 []chan [][]int64, hb2 []chan [][]int64, k int) {
	now := time.Now().Unix() // current local time
	counter := 0
	hbtable := make([][]int64, numWorkers+3)
	// initialize hbtable
	for i := 0; i < numWorkers+3; i++ {
		hbtable[i] = make([]int64, 2)
		hbtable[i][0] = 0
		hbtable[i][1] = now
	}

	for {
		time.Sleep(100 * time.Millisecond)
		hbtable = updateTable(k, hbtable, counter, hb1, hb2)
		counter++
	}
}

// Heartbeat function for Master
func masterHeartbeat(hb1 []chan [][]int64, hb2 []chan [][]int64, k int, corpses chan []bool, kill chan string, killMaster chan string, mrData MasterData) {
	now := time.Now().Unix() // current local time
	counter := 0
	currentTable := make([][]int64, numWorkers+3)
	previousTable := make([][]int64, numWorkers+3)

	// initialize hbtable
	for i := 0; i < numWorkers+3; i++ {
		currentTable[i] = make([]int64, 2)
		previousTable[i] = make([]int64, 2)
		currentTable[i][0] = 0
		previousTable[i][0] = 0
		currentTable[i][1] = now
		previousTable[i][1] = now
	}
	deadWorkers := make([]bool, numWorkers+3)
	for i := 0; i < 8; i++ {
		deadWorkers[i] = false
	}
	for {
		select{
		case reply := <- kill:
			if reply == "die" {
				return
			}
		default:
		}
		time.Sleep(100 * time.Millisecond)
		currentTable = updateTable(k, previousTable, counter, hb1, hb2)
		for i := 0; i < numWorkers+3; i++ {
			if currentTable[k][1] - previousTable[i][1] > 2 {
				fmt.Println(currentTable[k][1], previousTable[i][1])
				if i == numWorkers+1 || i == numWorkers+2{
					fmt.Println("Shadow master died")
					go shadowMaster(mrData.toShadowMasters[i-numWorkers+1], mrData.hb1, mrData.hb2, mrData.numWorkers, i, mrData, killMaster)
				} else {
					fmt.Println("\n\n-------------------killed worker :", i, "\n")
					deadWorkers[i] = true
				}
			}
		}
		previousTable = currentTable
		corpses <- deadWorkers
		counter++
	}
}

// Shadow Master Node
func shadowMaster(copier chan string, hb1 []chan [][]int64, hb2 []chan [][]int64, masterID int, selfID int, mrData MasterData, kill chan string){
	// replicate logs to two shadowmasters that monitor if the master dies
	logs := make([]string, 0)
	killHB := make(chan string, 3)
	var isMasterDead = make(chan bool, 1)
	go shadowHeartbeat(hb1, hb2, masterID, isMasterDead, selfID, killHB)
	masterNotDead := true

	for masterNotDead{
		select{
			case copy := <-copier:
				logs = append(logs, copy)
				// check if to die
				currentStep := copy
			if currentStep == "step start" {
				// update logs
				currentStep = "step load"
				mrData.log = append(mrData.log, currentStep)
				
			} else if currentStep == "step load" {
				// master load
				mrData = distributeTasks(mrData)
				currentStep = "step working"
				mrData.log = append(mrData.log, currentStep)
				
			} else if currentStep == "step working" {
				// master working
				currentStep = "step cleanup"
				mrData.log = append(mrData.log, currentStep)

			} else if currentStep == "cleanup" {
				// cleanup
				currentStep = "step end"
				mrData.log = append(mrData.log, currentStep)
				
			}else if currentStep == "step end"{
				killHB <- "kill"
				return
			}
			case isDead := <-isMasterDead:
				masterNotDead = isDead
			default:	
		}
		if !masterNotDead{
			fmt.Println("Shadow master becomes running master.")
			go master(mrData, hb1, hb2, logs, kill)
			masterNotDead = true
		}
	}
}

// Heartbeat for Shadow Master
func shadowHeartbeat(hb1 []chan [][]int64, hb2 []chan [][]int64, masterID int, isMasterAlive chan bool, selfID int, killHB chan string) string{
	now := time.Now().Unix() // current local time
	counter := 0
	currentTable := make([][]int64, numWorkers+3)
	previousTable := make([][]int64, numWorkers+3)

	// initialize hbtable
	for i := 0; i < 11; i++ {
		currentTable[i] = make([]int64, 2)
		previousTable[i] = make([]int64, 2)
		currentTable[i][0] = 0
		previousTable[i][0] = 0
		currentTable[i][1] = now
		previousTable[i][1] = now
	}
	
	for {
		select {
			case reply := <- killHB:
				return reply
			default:
		}
		time.Sleep(100 * time.Millisecond)
		currentTable = updateTable(selfID, previousTable, counter, hb1, hb2)

		if currentTable[selfID][1] - previousTable[masterID][1] > 2 {
			if selfID == masterID + 1{
				fmt.Println("\n----- The Running Master has died -----\n")
				isMasterAlive <- false
				currentTable = updateTable(masterID, currentTable, counter, hb1, hb2)
			}
		}
		previousTable = currentTable
		counter++
	}
}
