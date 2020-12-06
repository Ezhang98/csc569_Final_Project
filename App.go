// Group Members: Evan Zhang, Alexander Garcia and Christina Monahan
// Distributed System for Tuning Hyperparameters of Neural Networks
// PAXOS election, consensus, and recovery algorithm

package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/dathoangnd/gonet"
)

type UIWindow struct {
	TrainData  string
	TestData   string
	ModelCount int
	Models     []ModelConfig
}

type ModelConfig struct {
	ModelID int
	Name    string
	NeuralNet
	Model2Params
	Model3Params
}

type NeuralNet struct {
	InputNodes      int
	NumHiddenLayers int
	OutputNodes     int
	NumEpochs       int
	LearningRate    float64
	Momentum        float64
}

type Model2Params struct {
	Layers   int
	Learning string
}

type Model3Params struct {
	Trees    int
	MaxDepth int
}

var mainwin *ui.Window
var modelCount int = 0
var models []ui.Control = make([]ui.Control, 5)
var windowData UIWindow
var results *ui.Grid
var resList []*ui.Label = make([]*ui.Label, 35)

func makeModelParam(m ModelConfig) ui.Control {
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	grid := ui.NewGrid()
	grid.SetPadded(true)

	hbox.Append(grid, false)

	form := ui.NewForm()
	form.SetPadded(true)
	grid.Append(form,
		0, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	initial := 0
	if m.Name == "Model 2" {
		initial = 1
	}
	if m.Name == "Model 3" {
		initial = 2
	}
	alignment := ui.NewCombobox()
	// note that the items match with the values of the uiDrawTextAlign values
	alignment.Append("Neural Network")
	alignment.Append("Model 2")
	alignment.Append("Model 3")
	alignment.SetSelected(initial)
	alignment.OnSelected(func(*ui.Combobox) {
		s := alignment.Selected()
		if s == 0 {
			windowData.Models[m.ModelID].Name = "Neural Network"
			hbox1 := ui.NewHorizontalBox()
			hbox1.SetPadded(true)
			grid.Append(hbox1,
				0, 1, 1, 1,
				false, ui.AlignFill, false, ui.AlignFill)

			// InputNodes   	int
			// NumHiddenLayers 	int
			// OutputNodes		int
			// NumEpochs		int
			// LearningRate	float64
			// Momentum		float64

			// # of input nodes
			inputNodes := ui.NewSpinbox(0, 10000)
			inputNodes.OnChanged(func(*ui.Spinbox) {
				windowData.Models[m.ModelID].InputNodes = inputNodes.Value()
			})

			form1 := ui.NewForm()
			form1.SetPadded(true)
			hbox1.Append(form1, false)
			form1.Append("# Input Nodes", inputNodes, false)

			// # of hidden layers
			layers := ui.NewSpinbox(0, 10000)
			layers.OnChanged(func(*ui.Spinbox) {
				windowData.Models[m.ModelID].NumHiddenLayers = layers.Value()
			})

			form2 := ui.NewForm()
			form2.SetPadded(true)
			hbox1.Append(form2, false)
			form2.Append("# Hidden Layers", layers, false)

			// # of output nodes
			outputNodes := ui.NewSpinbox(0, 10000)
			outputNodes.OnChanged(func(*ui.Spinbox) {
				windowData.Models[m.ModelID].OutputNodes = outputNodes.Value()
			})

			form3 := ui.NewForm()
			form3.SetPadded(true)
			hbox1.Append(form3, false)
			form3.Append("# Output Nodes", outputNodes, false)

			// # of epochs
			epochs := ui.NewSpinbox(0, 10000)
			epochs.OnChanged(func(*ui.Spinbox) {
				windowData.Models[m.ModelID].NumEpochs = epochs.Value()
			})

			form4 := ui.NewForm()
			form4.SetPadded(true)
			hbox1.Append(form4, false)
			form4.Append("# of Epochs", epochs, false)

			// learning rate
			lrate := ui.NewEntry()
			windowData.Models[m.ModelID].LearningRate = 0.0
			lrate.OnChanged(func(*ui.Entry) {
				f, err := strconv.ParseFloat(lrate.Text(), 64)
				if err == nil {
					windowData.Models[m.ModelID].LearningRate = f
				} else {
					ui.MsgBoxError(mainwin,
						"Not a Number.",
						"Please enter a numerical value.")
					lrate.SetText("")
				}
			})

			form5 := ui.NewForm()
			form5.SetPadded(true)
			hbox1.Append(form5, false)
			form5.Append("Learning Rate", lrate, false)

			// momentum
			momentum := ui.NewEntry()
			windowData.Models[m.ModelID].Momentum = 0.0
			momentum.OnChanged(func(*ui.Entry) {
				f, err := strconv.ParseFloat(momentum.Text(), 64)
				if err == nil {
					windowData.Models[m.ModelID].Momentum = f
				} else {
					ui.MsgBoxError(mainwin,
						"Not a Number.",
						"Please enter a numerical value.")
					momentum.SetText("")
				}
			})

			form6 := ui.NewForm()
			form6.SetPadded(true)
			hbox1.Append(form6, false)
			form6.Append("Momentum", momentum, false)

			// end

		} else if s == 1 {
			windowData.Models[m.ModelID].Name = "Model 2"
			hbox1 := ui.NewHorizontalBox()
			hbox1.SetPadded(true)
			grid.Append(hbox1,
				0, 1, 1, 1,
				false, ui.AlignFill, false, ui.AlignFill)
			// layers and learning rate

			layers := ui.NewSpinbox(0, 10000)
			layers.OnChanged(func(*ui.Spinbox) {
				windowData.Models[m.ModelID].Layers = layers.Value()
			})

			form2 := ui.NewForm()
			form2.SetPadded(true)
			hbox1.Append(form2, false)
			form2.Append("numLayers", layers, false)

			lrate := ui.NewEntry()
			lrate.SetText("0.1")
			windowData.Models[m.ModelID].Learning = "0.1"
			lrate.OnChanged(func(*ui.Entry) {
				windowData.Models[m.ModelID].Learning = lrate.Text()
			})

			form1 := ui.NewForm()
			form1.SetPadded(true)
			hbox1.Append(form1, false)
			form1.Append("learning rate", lrate, false)
		} else {
			windowData.Models[m.ModelID].Name = "Model 3"
			hbox1 := ui.NewHorizontalBox()
			hbox1.SetPadded(true)
			grid.Append(hbox1,
				0, 1, 1, 1,
				false, ui.AlignFill, false, ui.AlignFill)
			// numTrees and max depth
			numTrees := ui.NewSpinbox(0, 10000)
			numTrees.OnChanged(func(*ui.Spinbox) {
				windowData.Models[m.ModelID].Trees = numTrees.Value()
			})

			form1 := ui.NewForm()
			form1.SetPadded(true)
			hbox1.Append(form1, false)
			form1.Append("numTrees", numTrees, false)

			maxDepth := ui.NewSpinbox(0, 10000)
			maxDepth.OnChanged(func(*ui.Spinbox) {
				windowData.Models[m.ModelID].MaxDepth = maxDepth.Value()
			})

			form2 := ui.NewForm()
			form2.SetPadded(true)
			hbox1.Append(form2, false)
			form2.Append("maxDepth", maxDepth, false)
		}
	})

	if m.Name == "Neural Network" {
		hbox1 := ui.NewHorizontalBox()
		hbox1.SetPadded(true)
		grid.Append(hbox1,
			0, 1, 1, 1,
			false, ui.AlignFill, false, ui.AlignFill)

		// # of input nodes
		inputNodes := ui.NewSpinbox(0, 10000)
		inputNodes.SetValue(windowData.Models[m.ModelID].InputNodes)
		inputNodes.OnChanged(func(*ui.Spinbox) {
			windowData.Models[m.ModelID].InputNodes = inputNodes.Value()
		})

		form1 := ui.NewForm()
		form1.SetPadded(true)
		hbox1.Append(form1, false)
		form1.Append("# Input Nodes", inputNodes, false)

		// # of hidden layers
		layers := ui.NewSpinbox(0, 10000)
		layers.SetValue(windowData.Models[m.ModelID].NumHiddenLayers)
		layers.OnChanged(func(*ui.Spinbox) {
			windowData.Models[m.ModelID].NumHiddenLayers = layers.Value()
		})

		form2 := ui.NewForm()
		form2.SetPadded(true)
		hbox1.Append(form2, false)
		form2.Append("# Hidden Layers", layers, false)

		// # of output nodes
		outputNodes := ui.NewSpinbox(0, 10000)
		outputNodes.SetValue(windowData.Models[m.ModelID].OutputNodes)
		outputNodes.OnChanged(func(*ui.Spinbox) {
			windowData.Models[m.ModelID].OutputNodes = outputNodes.Value()
		})

		form3 := ui.NewForm()
		form3.SetPadded(true)
		hbox1.Append(form3, false)
		form3.Append("# Output Nodes", outputNodes, false)

		// # of epochs
		epochs := ui.NewSpinbox(0, 10000)
		epochs.SetValue(windowData.Models[m.ModelID].NumEpochs)
		epochs.OnChanged(func(*ui.Spinbox) {
			windowData.Models[m.ModelID].NumEpochs = epochs.Value()
		})

		form4 := ui.NewForm()
		form4.SetPadded(true)
		hbox1.Append(form4, false)
		form4.Append("# of Epochs", epochs, false)

		// learning rate
		lrate := ui.NewEntry()
		s := strconv.FormatFloat(windowData.Models[m.ModelID].LearningRate, 'g', -1, 64)
		lrate.SetText(s)
		lrate.OnChanged(func(*ui.Entry) {
			f, err := strconv.ParseFloat(lrate.Text(), 64)
			if err == nil {
				windowData.Models[m.ModelID].LearningRate = f
			} else {
				ui.MsgBoxError(mainwin,
					"Not a Number.",
					"Please enter a numerical value.")
				lrate.SetText("")
			}
		})

		form5 := ui.NewForm()
		form5.SetPadded(true)
		hbox1.Append(form5, false)
		form5.Append("Learning Rate", lrate, false)

		// momentum
		momentum := ui.NewEntry()
		s = strconv.FormatFloat(windowData.Models[m.ModelID].Momentum, 'g', -1, 64)
		momentum.SetText(s)
		momentum.OnChanged(func(*ui.Entry) {
			f, err := strconv.ParseFloat(momentum.Text(), 64)
			if err == nil {
				windowData.Models[m.ModelID].Momentum = f
			} else {
				ui.MsgBoxError(mainwin,
					"Not a Number.",
					"Please enter a numerical value.")
				momentum.SetText("")
			}
		})

		form6 := ui.NewForm()
		form6.SetPadded(true)
		hbox1.Append(form6, false)
		form6.Append("Momentum", momentum, false)
	} else if m.Name == "Model 2" {
		hbox1 := ui.NewHorizontalBox()
		hbox1.SetPadded(true)
		grid.Append(hbox1,
			0, 1, 1, 1,
			false, ui.AlignFill, false, ui.AlignFill)
		// layers and learning rate
		layers := ui.NewSpinbox(0, 10000)
		layers.SetValue(windowData.Models[m.ModelID].Layers)
		layers.OnChanged(func(*ui.Spinbox) {
			windowData.Models[m.ModelID].Layers = layers.Value()
		})

		form2 := ui.NewForm()
		form2.SetPadded(true)
		hbox1.Append(form2, false)
		form2.Append("numLayers", layers, false)

		lrate := ui.NewEntry()
		lrate.SetText(windowData.Models[m.ModelID].Learning)
		lrate.OnChanged(func(*ui.Entry) {
			windowData.Models[m.ModelID].Learning = lrate.Text()
		})

		form1 := ui.NewForm()
		form1.SetPadded(true)
		hbox1.Append(form1, false)
		form1.Append("learning rate", lrate, false)
	} else if m.Name == "Model 3" {
		hbox1 := ui.NewHorizontalBox()
		hbox1.SetPadded(true)
		grid.Append(hbox1,
			0, 1, 1, 1,
			false, ui.AlignFill, false, ui.AlignFill)
		// numTrees and max depth
		numTrees := ui.NewSpinbox(0, 10000)
		numTrees.SetValue(windowData.Models[m.ModelID].Trees)
		numTrees.OnChanged(func(*ui.Spinbox) {
			windowData.Models[m.ModelID].Trees = numTrees.Value()
		})

		form1 := ui.NewForm()
		form1.SetPadded(true)
		hbox1.Append(form1, false)
		form1.Append("numTrees", numTrees, false)

		maxDepth := ui.NewSpinbox(0, 10000)
		maxDepth.SetValue(windowData.Models[m.ModelID].MaxDepth)
		maxDepth.OnChanged(func(*ui.Spinbox) {
			windowData.Models[m.ModelID].MaxDepth = maxDepth.Value()
		})

		form2 := ui.NewForm()
		form2.SetPadded(true)
		hbox1.Append(form2, false)
		form2.Append("maxDepth", maxDepth, false)
	}

	form.Append("Model", alignment, false)

	return hbox
}

func generateFromState() *ui.Grid {
	grid := ui.NewGrid()
	grid.SetPadded(true)

	button := ui.NewButton("Training Data")
	entry := ui.NewEntry()
	entry.SetReadOnly(true)
	entry.SetText(windowData.TrainData)
	button.OnClicked(func(*ui.Button) {
		filename := ui.OpenFile(mainwin)
		if filename == "" {
			filename = "(cancelled)"
		}
		entry.SetText(filename)
		windowData.TrainData = filename
	})
	grid.Append(button,
		0, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)
	grid.Append(entry,
		1, 0, 1, 1,
		true, ui.AlignFill, false, ui.AlignFill)

	button1 := ui.NewButton("Test Data")
	entry1 := ui.NewEntry()
	entry1.SetReadOnly(true)
	entry1.SetText(windowData.TestData)
	button1.OnClicked(func(*ui.Button) {
		filename := ui.OpenFile(mainwin)
		if filename == "" {
			filename = "(cancelled)"
		}
		entry1.SetText(filename)
		windowData.TestData = filename
	})
	grid.Append(button1,
		0, 1, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)
	grid.Append(entry1,
		1, 1, 1, 1,
		true, ui.AlignFill, false, ui.AlignFill)
	modelCount = windowData.ModelCount
	for i := 0; i < windowData.ModelCount; i++ {
		m := makeModelParam(windowData.Models[i])
		grid.Append(m,
			0, i+2, 2, 1,
			true, ui.AlignFill, false, ui.AlignFill)
	}
	return grid
}

func makeToolbar2() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	msggrid := ui.NewGrid()
	msggrid.SetPadded(true)
	vbox.Append(msggrid, false)

	grid := ui.NewGrid()
	grid.SetPadded(true)
	vbox.Append(grid, false)

	results = ui.NewGrid()
	results.SetPadded(true)
	vbox.Append(results, false)

	button := ui.NewButton("Import Config")
	button.OnClicked(func(*ui.Button) {
		filename := ui.OpenFile(mainwin)
		if filename != "" {
			file, _ := ioutil.ReadFile(filename)
			temp := UIWindow{}
			_ = json.Unmarshal([]byte(file), &temp)
			windowData = temp
			vbox.Delete(2)
			vbox.Delete(1)
			grid = generateFromState()
			vbox.Append(grid, false)
			results = ui.NewGrid()
			results.SetPadded(true)
			vbox.Append(results, false)
			for i := 0; i < windowData.ModelCount; i++ {
				for j := 0; j < 7; j++ {
					s := fmt.Sprintf("%s Results - ", generateLabel((i*7)+j))
					resList[(i*7)+j] = ui.NewLabel(s)
					results.Append(resList[(i*7)+j],
						0, ((i+1)*7)+j, 1, 1,
						false, ui.AlignFill, false, ui.AlignFill)
				}
			}
		}
	})
	msggrid.Append(button,
		0, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)
	button = ui.NewButton("Add Model")
	button.OnClicked(func(*ui.Button) {
		if modelCount < 5 {
			var m ModelConfig
			m.Name = ""
			m.ModelID = modelCount
			windowData.Models[modelCount] = m
			model := makeModelParam(m)
			grid.Append(model,
				0, modelCount+2, 2, 1,
				true, ui.AlignFill, false, ui.AlignFill)
			for j := 0; j < 7; j++ {
				s := fmt.Sprintf("%s Results - ", generateLabel((windowData.ModelCount*7)+j))
				resList[(windowData.ModelCount*7)+j] = ui.NewLabel(s)
				results.Append(resList[(windowData.ModelCount*7)+j],
					0, ((windowData.ModelCount+1)*7)+j, 1, 1,
					false, ui.AlignFill, false, ui.AlignFill)
			}
			modelCount++
			windowData.ModelCount++
		}
	})
	msggrid.Append(button,
		1, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	button = ui.NewButton("Save Config")
	button.OnClicked(func(*ui.Button) {
		filename := ui.SaveFile(mainwin)
		file, _ := json.MarshalIndent(windowData, "", " ")
		_ = ioutil.WriteFile(filename, file, 0644)
	})
	msggrid.Append(button,
		2, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	button = ui.NewButton("Clear All")
	button.OnClicked(func(*ui.Button) {
		temp := UIWindow{}
		windowData = temp
		windowData.Models = make([]ModelConfig, 5)
		windowData.ModelCount = 0
		vbox.Delete(2)
		vbox.Delete(1)
		grid = generateFromState()
		vbox.Append(grid, false)
		results = ui.NewGrid()
		results.SetPadded(true)
		vbox.Append(results, false)
	})
	msggrid.Append(button,
		3, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	button = ui.NewButton("Run Models")
	button.OnClicked(func(*ui.Button) {
		go launchServers(7 * windowData.ModelCount)
		// runNN()
	})
	msggrid.Append(button,
		4, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	button = ui.NewButton("Training Data")
	entry := ui.NewEntry()
	entry.SetReadOnly(true)
	button.OnClicked(func(*ui.Button) {
		filename := ui.OpenFile(mainwin)
		if filename == "" {
			filename = "(cancelled)"
		}
		entry.SetText(filename)
		windowData.TrainData = filename
	})
	grid.Append(button,
		0, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)
	grid.Append(entry,
		1, 0, 1, 1,
		true, ui.AlignFill, false, ui.AlignFill)

	button1 := ui.NewButton("Test Data")
	entry1 := ui.NewEntry()
	entry1.SetReadOnly(true)
	button1.OnClicked(func(*ui.Button) {
		filename := ui.OpenFile(mainwin)
		if filename == "" {
			filename = "(cancelled)"
		}
		entry1.SetText(filename)
		windowData.TestData = filename
	})
	grid.Append(button1,
		0, 1, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)
	grid.Append(entry1,
		1, 1, 1, 1,
		true, ui.AlignFill, false, ui.AlignFill)

	return vbox
}

func setupUI() {
	mainwin = ui.NewWindow("libui Control Gallery", 1800, 900, true)
	windowData.Models = make([]ModelConfig, 5)
	windowData.ModelCount = 0
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})
	mainwin.SetMargined(true)
	mainwin.SetChild(makeToolbar2())

	mainwin.Show()
}

func parseCSV(path string) [][][]float64 {
	data := make([][][]float64, 0)

	// Open the file

	// s := strings.Split(path, "\\")
	// fixedPath := "../datasets/" + s[len(s)-1]

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
		if index != 0 && len(record) > 1 {
			floatarr := make([]float64, len(record)-1)
			expected := make([]float64, 10)
			for i := 0; i < len(record); i++ {
				if s, err := strconv.ParseFloat(record[i], 64); err == nil {
					// fmt.Println(s)
					if i == 0 {
						expected[int(s)] = 1
					} else {
						floatarr[i-1] = s
					}

				}
			}
			if len(floatarr) == 0 {
				break
			}
			oneEntry := [][]float64{floatarr, expected}

			data = append(data, oneEntry)
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

// index 0 [[28x28], expected_value]
// index 1

func generateLabel(id int) string {
	num1 := id / 7
	num2 := id % 7
	s := ""
	if num2 == 0 {
		s = fmt.Sprintf("Model %d Default", num1)
	} else if num2 == 1 {
		s = fmt.Sprintf("Model %d Double Epoch", num1)
	} else if num2 == 2 {
		s = fmt.Sprintf("Model %d Half Epoch", num1)
	} else if num2 == 3 {
		s = fmt.Sprintf("Model %d Double Learning Rate", num1)
	} else if num2 == 4 {
		s = fmt.Sprintf("Model %d Half Learning Rate", num1)
	} else if num2 == 5 {
		s = fmt.Sprintf("Model %d Double Momentum", num1)
	} else if num2 == 6 {
		s = fmt.Sprintf("Model %d Half Momentum", num1)
	}
	return s
}

func runNN(m ModelConfig, train [][][]float64, test [][][]float64) {
	hidden := make([]int, m.NumHiddenLayers)
	for j := 0; j < m.NumHiddenLayers; j++ {
		hidden[j] = 32
	}
	nn := gonet.New(len(train[0][0]), hidden, 10, true)

	// Train the network
	// Run for 3000 epochs
	// The learning rate is 0.4 and the momentum factor is 0.2
	// Enable debug mode to log learning error every 1000 iterations
	timer := time.Now()
	nn.Train(train, m.NumEpochs, m.LearningRate, m.Momentum, true)
	s := fmt.Sprintf("%s Results - ", generateLabel(m.ModelID))
	runtime := fmt.Sprintf("Runtime: %.5f seconds, ", time.Since(timer).Seconds())

	// Predict
	totalcorrect := 0.0
	for i := 0; i < len(test); i++ {
		if MinMax(test[i][1]) == MinMax(nn.Predict(test[i][0])) {
			totalcorrect += 1.0
		}
	}

	acc := fmt.Sprintf("Accuracy: %.2f percent", totalcorrect/float64(len(test))*100.0)
	s = s + runtime + acc

	fmt.Println(s)
	ui.QueueMain(func() {
		resList[m.ModelID].SetText(s)
	})

}

func MinMax(array []float64) int {
	index := 0
	max := 0.0
	for i, value := range array {
		if max < value {
			index = i
			max = value
		}
	}
	return index
}

func main() {
	ui.Main(setupUI)
}

// ------------------------------ PAXOS ------------------------------------
// number of workers
var numWorkers int

//struct to organize data into the master function
type MasterData struct {
	id              int
	numWorkers      int
	request         []chan ModelConfig
	commands        []chan string
	dataOut         []chan [][][]float64
	models          []ModelConfig
	replies         []chan string
	corpses         chan []bool
	working         []string
	finished        []string
	toShadowMasters []chan string
	log             []string
	hb1             []chan [][]int64
	hb2             []chan [][]int64
	test            []chan [][][]float64
	training        []chan [][][]float64
}

// Launches nodes and creates MasterData structure
func launchServers(numW int) {
	numWorkers = numW

	var mrData MasterData
	mrData.numWorkers = numWorkers
	mrData.request = make([]chan ModelConfig, mrData.numWorkers) // worker <- master : for task assignment
	mrData.commands = make([]chan string, mrData.numWorkers)     // master -> worker : master sends commands
	mrData.replies = make([]chan string, mrData.numWorkers)      // master <- worker : worker reply for task completion
	mrData.corpses = make(chan []bool, mrData.numWorkers)        // master <- heartbeat : workers that have died
	mrData.working = make([]string, mrData.numWorkers)           // which tasks assigned to which workers
	mrData.training = make([]chan [][][]float64, mrData.numWorkers)
	mrData.test = make([]chan [][][]float64, mrData.numWorkers)

	var hb1 = make([]chan [][]int64, mrData.numWorkers+3) // heartbeat channels to neighbors for read
	var hb2 = make([]chan [][]int64, mrData.numWorkers+3) // heartbeat channels to neighbors for write
	killMaster := make(chan string, 10)                   // channel to kill Master to verify replication and recovery from log

	// initialize heartbeat tables
	for i := 0; i < mrData.numWorkers+3; i++ {
		hb1[i] = make(chan [][]int64, 1024)
		hb2[i] = make(chan [][]int64, 1024)
	}
	mrData.hb1 = hb1
	mrData.hb2 = hb2

	// initialize worker channels
	for k := 0; k < mrData.numWorkers; k++ {
		mrData.request[k] = make(chan ModelConfig)
		mrData.replies[k] = make(chan string)
		mrData.working[k] = ""
		mrData.commands[k] = make(chan string)
		mrData.training[k] = make(chan [][][]float64)
		mrData.test[k] = make(chan [][][]float64)
	}

	// initialize shadowMasters
	numShadowMasters := 2
	// shadowMaster <- master : replication
	mrData.toShadowMasters = make([]chan string, numShadowMasters)
	for j := 0; j < numShadowMasters; j++ {
		mrData.toShadowMasters[j] = make(chan string, 10)
	}

	for i := 0; i < windowData.ModelCount; i++ {
		if windowData.Models[i].Name != "" {
			vanilla := windowData.Models[i]
			vanilla.ModelID = i * 7
			mrData.models = append(mrData.models, vanilla)

			doubleEpoch := windowData.Models[i]
			doubleEpoch.NumEpochs *= 2
			doubleEpoch.ModelID = (i * 7) + 1
			mrData.models = append(mrData.models, doubleEpoch)

			halfEpoch := windowData.Models[i]
			halfEpoch.NumEpochs /= 2
			halfEpoch.ModelID = (i * 7) + 2
			mrData.models = append(mrData.models, halfEpoch)

			doubleLearningRate := windowData.Models[i]
			doubleLearningRate.LearningRate *= 2
			doubleLearningRate.ModelID = (i * 7) + 3
			mrData.models = append(mrData.models, doubleLearningRate)

			halfLearningRate := windowData.Models[i]
			halfLearningRate.LearningRate /= 2
			halfLearningRate.ModelID = (i * 7) + 4
			mrData.models = append(mrData.models, halfLearningRate)

			doubleMomentum := windowData.Models[i]
			doubleMomentum.Momentum *= 2
			doubleMomentum.ModelID = (i * 7) + 5
			mrData.models = append(mrData.models, doubleMomentum)

			halfMomentum := windowData.Models[i]
			halfMomentum.Momentum /= 2
			halfMomentum.ModelID = (i * 7) + 6
			mrData.models = append(mrData.models, halfMomentum)
		}
	}

	// example
	// for i := 0; i < 4; i ++{
	// 	var newModel ModelConfig
	// 	mrData.models[i] = newModel
	// 	mrData.models[i].ModelID = i
	// 	mrData.models[i].numHiddenLayers = rand.Intn(5)
	// 	mrData.models[i].numEpochs = rand.Intn(20)
	// 	mrData.models[i].learningRate = rand.Float64()
	// 	mrData.models[i].momentum = rand.Float64()
	// }

	// start nodes
	masterlog := make([]string, 0)
	var endrun = make(chan string)
	go master(mrData, hb1, hb2, masterlog, killMaster, endrun)

	go shadowMaster(mrData.toShadowMasters[0], hb1, hb2, mrData.numWorkers, mrData.numWorkers+1, mrData, killMaster)

	go shadowMaster(mrData.toShadowMasters[1], hb1, hb2, mrData.numWorkers, mrData.numWorkers+2, mrData, killMaster)

	for {
		select {
		case reply := <-endrun:
			// fmt.Println("end main thread")
			if reply == "end" {
				break
			}
		}
	}
}

// Master node
func master(mrData MasterData, hb1 []chan [][]int64, hb2 []chan [][]int64, log []string, killMaster chan string, end chan string) {
	// initialize variables
	currentStep := "step start"
	mrData.log = log
	killHB := make(chan string, numWorkers)
	message := ""

	// Master Recovery: resume step after Master failure
	for i := 0; i < len(log); i++ {
		if len(log[i]) >= 4 && log[i][0:4] == "step" {
			currentStep = log[i]
		}
	}

	// Run master heartbeat
	go masterHeartbeat(hb1, hb2, mrData.numWorkers, mrData.corpses, killHB, killMaster, mrData)

	for {
		// fmt.Println("Master running: ", currentStep) // print in UI
		if currentStep == "step start" {
			// If master died partway through launching workers, find the last logged k and start from there
			k := 0
			if len(log) > 0 {
				last := strings.Split(log[len(log)-1], " ")
				if last[0] == "launch" {
					k, _ = strconv.Atoi(last[2])
				}
			}
			// fmt.Println(mrData.numWorkers)
			for ; k < mrData.numWorkers; k++ {
				go worker(mrData.training[k], mrData.test[k], mrData.request[k], mrData.commands[k], mrData.replies[k], hb1, hb2, k)
				// fmt.Printf("launch worker %d\n", k)
				message = fmt.Sprintf("launch worker %s", k)
				mrData.log = append(mrData.log, currentStep) //appends launch worker step to log
				mrData.toShadowMasters[0] <- message         //sends launch worker message to first Shadow Master channel
				mrData.toShadowMasters[1] <- message         //sends launch worker message to second Shadow Master channel
			}
			currentStep = "step working"
			mrData.log = append(mrData.log, currentStep) //appends load step to log
			mrData.toShadowMasters[0] <- currentStep     //sends load message to first Shadow Master channel
			mrData.toShadowMasters[1] <- currentStep     //sends load message to second Shadow Master channel

		} else if currentStep == "step working" {
			// manage the distributeTasks step
			// trainingdata := parseCSV(windowData.TrainData)
			// testdata := parseCSV(windowData.TestData)
			trainingdata := parseCSV("datasets/mnist_train_short.csv")
			testdata := parseCSV("datasets/mnist_train_short.csv")
			// fmt.Println("starting getting data")
			for k := 0; k < mrData.numWorkers; k++ {
				// fmt.Println("sending to worker", k)
				mrData.training[k] <- trainingdata
				mrData.test[k] <- testdata
			}
			// fmt.Println("finished loading data")
			mrData = distributeTasks(mrData)
			currentStep = "step cleanup"
			mrData.log = append(mrData.log, currentStep) //appends master distributeTasks step to log
			mrData.toShadowMasters[0] <- currentStep     //sends distributeTasks message to first Shadow Master channel
			mrData.toShadowMasters[1] <- currentStep     //sends distributeTasks message to second Shadow Master channel

		} else if currentStep == "step cleanup" {
			// cleanup workers who should now be done with all tasks
			// mrData = cleanup(mrData)

			currentStep = "step end"
			mrData.log = append(mrData.log, currentStep) //appends end message to log
			mrData.toShadowMasters[0] <- currentStep     //sends end message to first Shadow Master channel
			mrData.toShadowMasters[1] <- currentStep     //sends end message to second Shadow Master channel

			killHB <- "die"
			for i := 0; i < mrData.numWorkers; i++ {
				mrData.commands[i] <- "end"
			}
			end <- "end"
			fmt.Println("Running master is done.")
			break
		}
	}
}

func distributeTasks(mrData MasterData) MasterData {
	count := 0
	loop := true

	fmt.Println("Distributing Tasks Started...")
	mrData.finished = make([]string, len(mrData.models)) // which models have completed
	for j := 0; j < len(mrData.models); j++ {
		mrData.finished[j] = "not started"
	}
	for loop {
		for i := 0; i < mrData.numWorkers; i++ {
			// checks for available workers
			if mrData.working[i] == "" {
				for j := 0; j < len(mrData.models); j++ {
					if mrData.finished[j] == "not started" {
						mrData.finished[j] = "started"
						mrData.request[i] <- (mrData.models[j])
						mrData.commands[i] <- "m"
						mrData.working[i] = strconv.Itoa(j)
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
				mrData.finished[workerID] = "finished"
				count++
			case coffins := <-mrData.corpses:
				for j := 0; j < mrData.numWorkers; j++ {
					if coffins[j] == true {
						mrData.hb1[j] = make(chan [][]int64, numWorkers+3)
						mrData.hb2[j] = make(chan [][]int64, numWorkers+3)
						mrData.request[j] = make(chan ModelConfig)
						mrData.replies[j] = make(chan string)
						go worker(mrData.training[j], mrData.test[j], mrData.request[j], mrData.commands[j], mrData.replies[j], mrData.hb1, mrData.hb2, j)
						tempModelID, _ := strconv.Atoi(mrData.working[j])
						mrData.finished[tempModelID] = "not started"
						mrData.working[j] = ""
						coffins[j] = false
					}
				}
			default:
			}
		}
		// fmt.Println(mrData.working)
		// checks that all models have completed
		if count >= len(mrData.models) {
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

func worker(train chan [][][]float64, test chan [][][]float64, frommaster chan ModelConfig, commands chan string, reply chan string, hb1 []chan [][]int64, hb2 []chan [][]int64, k int) {
	killHB := make(chan string, 1)
	go heartbeat(hb1, hb2, k, killHB)
	// task := ""
	var trainingdata [][][]float64
	var testdata [][][]float64
	var indivModel ModelConfig
	trainingdata = <-train
	testdata = <-test
	indivModel = <-frommaster
	for {
		select {
		case trainingdata = <-train:
			continue
		case testdata = <-test:
			continue
		case task := <-commands:
			// fmt.Println(task)
			tasks := strings.Split(task, "_")
			if tasks[0] == "end" {
				// fmt.Println("ending worker", k)
				killHB <- "end"
				return
			}
			if tasks[0] == "m" {

				if indivModel.NumHiddenLayers > 1 {
					runNN(indivModel, trainingdata, testdata)
				}
			}
			reply <- strconv.Itoa(k) + "_" + strconv.Itoa(indivModel.ModelID)
		default:
		}
	}

}

// Shuts down worker nodes
func cleanup(mrData MasterData) MasterData {
	// Check that all workers shutdown
	for i := 0; i < mrData.numWorkers; i++ {
		mrData.commands[i] <- "end_0_" + strconv.Itoa(i)
		msg := <-mrData.replies[i]
		num, _ := strconv.Atoi(msg)
		mrData.working[num] = ""
	}
	return mrData
}

// Helper function to find max
func max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func printModelParam(model ModelConfig) {
	fmt.Println("ID", model.ModelID, ", epochs", model.NumEpochs, ", learning rate", model.LearningRate,
		", momentum", model.Momentum, ", hidden layers", model.NumHiddenLayers)
}

// Update heartbeat tables for Master, 2 ShadowMasters and 8 workers
func updateTable(index int, hbtable [][]int64, counter int, hb1 []chan [][]int64, hb2 []chan [][]int64) [][]int64 {
	next := index + 1
	prev := index - 1
	if prev < 0 {
		prev = numWorkers + 2
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
	if counter%1 == 0 {
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
	// fmt.Println(index, " temp ", temp, "hbtable ", hbtable)
	temp[index][0] = hbtable[index][0] + 1
	temp[index][1] = now
	// send table
	hb1[index] <- temp
	hb2[index] <- temp
	return temp
}

// Heartbeat function for all workers
func heartbeat(hb1 []chan [][]int64, hb2 []chan [][]int64, k int, endHB chan string) {
	now := time.Now().Unix() // current local time
	counter := 0
	hbtable := make([][]int64, numWorkers+3)
	// fmt.Println(numWorkers)
	// initialize hbtable
	for i := 0; i < numWorkers+3; i++ {
		hbtable[i] = make([]int64, 2)
		hbtable[i][0] = 0
		hbtable[i][1] = now
	}

	for {
		fmt.Print("whb", k)
		time.Sleep(100 * time.Millisecond)
		hbtable = updateTable(k, hbtable, counter, hb1, hb2)
		counter++
		select {
		case reply := <-endHB:
			if reply == "end" {
				return
			}
		default:
		}
	}
}

// Heartbeat function for Master
func masterHeartbeat(hb1 []chan [][]int64, hb2 []chan [][]int64, k int, corpses chan []bool, kill chan string, killMaster chan string, mrData MasterData) {
	now := time.Now().Unix() // current local time
	counter := 0
	currentTable := make([][]int64, numWorkers+3)
	previousTable := make([][]int64, numWorkers+3)

	fmt.Println("master heartbeat started")

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
	for i := 0; i < numWorkers+3; i++ {
		deadWorkers[i] = false
	}
	for {
		select {
		case reply := <-kill:
			if reply == "die" {
				fmt.Println("master heartbeat ended")
				return
			}
		default:
		}
		time.Sleep(100 * time.Millisecond)
		currentTable = updateTable(k, previousTable, counter, hb1, hb2)
		for i := 0; i < numWorkers+3; i++ {
			if currentTable[k][1]-previousTable[i][1] > 2 {
				fmt.Println(currentTable[k][1], previousTable[i][1])
				if i == numWorkers+1 || i == numWorkers+2 {
					fmt.Println("Shadow master died")
					go shadowMaster(mrData.toShadowMasters[i-numWorkers], mrData.hb1, mrData.hb2, mrData.numWorkers, i, mrData, killMaster)
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
func shadowMaster(copier chan string, hb1 []chan [][]int64, hb2 []chan [][]int64, masterID int, selfID int, mrData MasterData, kill chan string) {
	// replicate logs to two shadowmasters that monitor if the master dies
	logs := make([]string, 0)
	killHB := make(chan string, 3)
	var isMasterDead = make(chan bool, 1)
	go shadowHeartbeat(hb1, hb2, masterID, isMasterDead, selfID, killHB)
	masterNotDead := true

	for masterNotDead {
		select {
		case copy := <-copier:
			// fmt.Println("shadow", selfID, copy)
			logs = append(logs, copy)
			// check if to die
			currentStep := copy
			if currentStep == "step start" {
				// update logs
				currentStep = "step load"
				mrData.log = append(mrData.log, currentStep)

			} else if currentStep == "step working" {
				// master working
				currentStep = "step cleanup"
				mrData.log = append(mrData.log, currentStep)

			} else if currentStep == "step cleanup" {
				// cleanup
				currentStep = "step end"
				mrData.log = append(mrData.log, currentStep)

				killHB <- "kill"
				// fmt.Println("shadow", selfID, "ending")
				return
			}
		case isDead := <-isMasterDead:
			masterNotDead = isDead
		default:
		}
		if !masterNotDead {
			fmt.Println("Shadow master becomes running master.")
			var endrun = make(chan string)
			go master(mrData, hb1, hb2, logs, kill, endrun)
			masterNotDead = true
		}
	}
}

// Heartbeat for Shadow Master
func shadowHeartbeat(hb1 []chan [][]int64, hb2 []chan [][]int64, masterID int, isMasterAlive chan bool, selfID int, killHB chan string) string {
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

	for {
		select {
		case reply := <-killHB:
			// fmt.Println("shadow heartbeat", selfID, "ending")
			return reply
		default:
		}
		time.Sleep(100 * time.Millisecond)
		currentTable = updateTable(selfID, previousTable, counter, hb1, hb2)

		if currentTable[selfID][1]-previousTable[masterID][1] > 2 {
			if selfID == masterID+1 {
				fmt.Println("\n----- The Running Master has died -----\n")
				isMasterAlive <- false
				currentTable = updateTable(masterID, currentTable, counter, hb1, hb2)
			}
		}
		previousTable = currentTable
		counter++
	}
}
