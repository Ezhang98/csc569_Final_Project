package main

import (
	"fmt"
	"log"
	"io"
	"os"
	"encoding/csv"
	"strconv"
	"time"
	// "github.com/dathoangnd/gonet"
	"github.com/navossoc/bayesian"
	// "github.com/cdipaolo/goml/cluster"
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
		if index != 0{
			floatarr := make([]float64, len(record))
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
	fmt.Println(data)
	return data
}



func main() {
	timer := time.Now()
	// parseCSV("../datasets/exams.csv")
	data := parseCSV("../datasets/exams.csv")
	
	dic := make(map[float64]bool)
	for i := 0; i < len(data);i++{
		dic[data[i][1][0]] = true
	}
	fmt.Println(dic)
	// model := NewKMeans(4, 15, data[0])

	// if model.Learn() != nil {
	// 	panic("Oh NO!!! There was an error learning!!")
	// }

	// // now you can predict like normal!
	// guess, err := model.Predict([]float64{-3, 6})
	// if err != nil {
	// 	panic("prediction error")
	// }


	// // or if you want to get the clustering
	// // results from the data
	// results := model.Guesses()

	// // you can also concat that with the
	// // training set and save it to a file
	// // (if you wanted to plot it or something)
	// err = model.SaveClusteredData("/tmp/.goml/KMeansResults.csv")
	// if err != nil {
	// 	panic("file save error")
	// }

	// // you can also persist the model to a
	// // file
	// err = model.PersistToFile("/tmp/.goml/KMeans.json")
	// if err != nil {
	// 	panic("file save error")
	// }

	// // and also restore from file (at a
	// // later time if you want)
	// err = model.RestoreFromFile("/tmp/.goml/KMeans.json")
	// if err != nil {
	// 	panic("file save error")
	// }
	fmt.Printf("\nRuntime: %.5f seconds\n", time.Since(timer).Seconds())
	// // XOR traning data
	// trainingData := [][][]float64{
	// 	{{0, 0}, {0}},
	// 	{{0, 1}, {1}},
	// 	{{1, 0}, {1}},
	// 	{{1, 1}, {0}},
	// }

	// // Create a neural network
	// // 2 nodes in the input layer
	// // 2 hidden layers with 4 nodes each
	// // 1 node in the output layer
	// // The problem is classification, not regression
	// nn := gonet.New(2, []int{4, 4}, 1, false)

	// // Train the network
	// // Run for 3000 epochs
	// // The learning rate is 0.4 and the momentum factor is 0.2
	// // Enable debug mode to log learning error every 1000 iterations
	// nn.Train(trainingData, 3000, 0.4, 0.2, true)

	// // Predict
	// testInput := []float64{1, 0}
	// fmt.Printf("%f XOR %f => %f\n", testInput[0], testInput[1], nn.Predict(testInput)[0])
	// // 1.000000 XOR 0.000000 => 0.943074

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