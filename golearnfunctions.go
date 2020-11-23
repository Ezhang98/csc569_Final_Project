// the IDX file format is a simple format for vectors and multidimensional matrices of various numerical types.
// The basic format is

// magic number
// size in dimension 0
// size in dimension 1
// size in dimension 2
// .....
// size in dimension N
// data

// The magic number is an integer (MSB first). The first 2 bytes are always 0.

// The third byte codes the type of the data:
// 0x08: unsigned byte
// 0x09: signed byte
// 0x0B: short (2 bytes)
// 0x0C: int (4 bytes)
// 0x0D: float (4 bytes)
// 0x0E: double (8 bytes)

// The 4-th byte codes the number of dimensions of the vector/matrix: 1 for vectors, 2 for matrices....

// The sizes in each dimension are 4-byte integers (MSB first, high endian, like in most non-Intel processors).

// The data is stored like in a C array, i.e. the index in the last dimension changes the fastest.
 
// 60000 28x28 matrices of unsigned ints
// Happy hacking.

package main

import 
(
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/knn"
)

// func main() {
// 	f, err := os.Open("t10k-images.idx3-ubyte")
// 	if err != nil {
// 			fmt.Println(err)
// 			return
// 	}
// 	defer f.Close()
// 	data := make([]byte, 4096)
// }



func knn() {
	// Load in a dataset, with headers. Header attributes will be stored.
	// Think of instances as a Data Frame structure in R or Pandas.
	// You can also create instances from scratch.
	rawData, err := base.ParseCSVToInstances("datasets/iris.csv", false)
	if err != nil {
		panic(err)
	}

	// Print a pleasant summary of your data.
	fmt.Println(rawData)

	//Initialises a new KNN classifier
	cls := knn.NewKnnClassifier("euclidean", "linear", 2)

	//Do a training-test split
	trainData, testData := base.InstancesTrainTestSplit(rawData, 0.50)
	cls.Fit(trainData)

	//Calculates the Euclidean distance and returns the most popular label
	predictions, err := cls.Predict(testData)
	if err != nil {
		panic(err)
	}

	// Prints precision/recall metrics
	confusionMat, err := evaluation.GetConfusionMatrix(testData, predictions)
	if err != nil {
		panic(fmt.Sprintf("Unable to get confusion matrix: %s", err.Error()))
	}
	fmt.Println(evaluation.GetSummary(confusionMat))
}

func main() {
	knn()
}