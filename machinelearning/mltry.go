package main

import (
	"fmt"
	"log"
	"io"
	"os"
	"encoding/csv"
	"strconv"
	// "time"
	"github.com/dathoangnd/gonet"
	// "github.com/cdipaolo/goml/cluster"
	// "github.com/cdipaolo/goml/linear"
)



//TODO: parse data
//TODO: parse hyperparameters

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
		if index != 0 && len(record) > 1{
			floatarr := make([]float64, len(record) - 1)
			expected := make([]float64, 1)
			for i := 0; i < len(record); i++ {
				if s, err := strconv.ParseFloat(record[i], 64); err == nil {
					if i == 0{
						expected[0] = s
					}else{
						floatarr[i - 1] = s
					}
					// fmt.Println(s)
				}
			}
			if len(floatarr) == 0{
				break
			}
			
			one_entry := [][]float64{floatarr, expected}
			
			data = append(data, one_entry)
			// fmt.Println(one_entry)

		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		index++
	}
	// fmt.Println(data)
	return data
}



func main() {

	train := parseCSV("../datasets/mnist_train.csv")
	test := parseCSV("../datasets/mnist_test.csv")
	// train := parseCSV("../datasets/exams.csv")
	// test := parseCSV("../datasets/exam.csv")
	// XOR traning data
	// trainingData := [][][]float64{
	// 	{{0, 0}, {0}},
	// 	{{0, 1}, {1}},
	// 	{{1, 0}, {1}},
	// 	{{1, 1}, {0}},
	// }
	// for i := 0 ; i < 10; i++{
	// 	fmt.Println(test[i][1], "\n\n")
		
	// } 


	// Create a neural network
	// 2 nodes in the input layer
	// 2 hidden layers with 4 nodes each
	// 1 node in the output layer
	// The problem is classification, not regression
	
	fmt.Println("size of input", len(train[0][0]))
	nn := gonet.New(len(train[0][0]), []int{100, 50, 25}, 1, false)
	// func New(nInputs int, nHiddens []int, nOutputs int, isRegression bool) NN
	// Train the network
	// Run for 3000 epochs
	// The learning rate is 0.4 and the momentum factor is 0.2
	// Enable debug mode to log learning error every 1000 iterations
	nn.Train(train, 20, 0.1, 0.2, true)

	// Predict
	totalcorrect := 0.0
	for i := 0; i < len(test); i++ {
		// fmt.Println(test[i][0])
		fmt.Printf("actual: %f, predicted: %f\n", test[i][1][0], nn.Predict(test[i][0]))
		if nn.Predict(test[i][0])[0] == test[i][1][0]{
			totalcorrect += 1.0
		}
	}
	fmt.Printf("Percent correct: %f\n", totalcorrect/float64(len(test)))



	// // Save the model
	// nn.Save("model.json")

	// // Load the model
	// nn2, err := gonet.Load("model.json")
	// if err != nil {
	// 	log.Fatal("Load model failed.")
	// }
	// fmt.Printf("%f XOR %f => %f\n", testInput[0], testInput[1], nn2.Predict(testInput)[0])
	// 1.000000 XOR 0.000000 => 0.943074



}