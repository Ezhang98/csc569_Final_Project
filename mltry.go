package main

import (
	"fmt"
	"log"

	"github.com/dathoangnd/gonet"
)

//TODO: parse data
//TODO: parse hyperparameters

func main() {
	// XOR traning data
	trainingData := [][][]float64{
		{{0, 0}, {0}},
		{{0, 1}, {1}},
		{{1, 0}, {1}},
		{{1, 1}, {0}},
	}

	// Create a neural network
	// 2 nodes in the input layer
	// 2 hidden layers with 4 nodes each
	// 1 node in the output layer
	// The problem is classification, not regression
	nn := gonet.New(2, []int{4, 4}, 1, false)

	// Train the network
	// Run for 3000 epochs
	// The learning rate is 0.4 and the momentum factor is 0.2
	// Enable debug mode to log learning error every 1000 iterations
	nn.Train(trainingData, 3000, 0.4, 0.2, true)

	// Predict
	testInput := []float64{1, 0}
	fmt.Printf("%f XOR %f => %f\n", testInput[0], testInput[1], nn.Predict(testInput)[0])
	// 1.000000 XOR 0.000000 => 0.943074

	// Save the model
	nn.Save("model.json")

	// Load the model
	nn2, err := gonet.Load("model.json")
	if err != nil {
		log.Fatal("Load model failed.")
	}
	fmt.Printf("%f XOR %f => %f\n", testInput[0], testInput[1], nn2.Predict(testInput)[0])
	// 1.000000 XOR 0.000000 => 0.943074
}