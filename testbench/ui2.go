package main

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
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
			inputNodes := ui.NewSpinbox(0, 100)
			inputNodes.OnChanged(func(*ui.Spinbox) {
				windowData.Models[m.ModelID].InputNodes = inputNodes.Value()
			})

			form1 := ui.NewForm()
			form1.SetPadded(true)
			hbox1.Append(form1, false)
			form1.Append("# Input Nodes", inputNodes, false)

			// # of hidden layers
			layers := ui.NewSpinbox(0, 100)
			layers.OnChanged(func(*ui.Spinbox) {
				windowData.Models[m.ModelID].NumHiddenLayers = layers.Value()
			})

			form2 := ui.NewForm()
			form2.SetPadded(true)
			hbox1.Append(form2, false)
			form2.Append("# Hidden Layers", layers, false)

			// # of output nodes
			outputNodes := ui.NewSpinbox(0, 100)
			outputNodes.OnChanged(func(*ui.Spinbox) {
				windowData.Models[m.ModelID].OutputNodes = outputNodes.Value()
			})

			form3 := ui.NewForm()
			form3.SetPadded(true)
			hbox1.Append(form3, false)
			form3.Append("# Output Nodes", outputNodes, false)

			// # of epochs
			epochs := ui.NewSpinbox(0, 100)
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

			layers := ui.NewSpinbox(0, 100)
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
			numTrees := ui.NewSpinbox(0, 100)
			numTrees.OnChanged(func(*ui.Spinbox) {
				windowData.Models[m.ModelID].Trees = numTrees.Value()
			})

			form1 := ui.NewForm()
			form1.SetPadded(true)
			hbox1.Append(form1, false)
			form1.Append("numTrees", numTrees, false)

			maxDepth := ui.NewSpinbox(0, 100)
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
		inputNodes := ui.NewSpinbox(0, 100)
		inputNodes.SetValue(windowData.Models[m.ModelID].InputNodes)
		inputNodes.OnChanged(func(*ui.Spinbox) {
			windowData.Models[m.ModelID].InputNodes = inputNodes.Value()
		})

		form1 := ui.NewForm()
		form1.SetPadded(true)
		hbox1.Append(form1, false)
		form1.Append("# Input Nodes", inputNodes, false)

		// # of hidden layers
		layers := ui.NewSpinbox(0, 100)
		layers.SetValue(windowData.Models[m.ModelID].NumHiddenLayers)
		layers.OnChanged(func(*ui.Spinbox) {
			windowData.Models[m.ModelID].NumHiddenLayers = layers.Value()
		})

		form2 := ui.NewForm()
		form2.SetPadded(true)
		hbox1.Append(form2, false)
		form2.Append("# Hidden Layers", layers, false)

		// # of output nodes
		outputNodes := ui.NewSpinbox(0, 100)
		outputNodes.SetValue(windowData.Models[m.ModelID].OutputNodes)
		outputNodes.OnChanged(func(*ui.Spinbox) {
			windowData.Models[m.ModelID].OutputNodes = outputNodes.Value()
		})

		form3 := ui.NewForm()
		form3.SetPadded(true)
		hbox1.Append(form3, false)
		form3.Append("# Output Nodes", outputNodes, false)

		// # of epochs
		epochs := ui.NewSpinbox(0, 100)
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
		layers := ui.NewSpinbox(0, 100)
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
		numTrees := ui.NewSpinbox(0, 100)
		numTrees.SetValue(windowData.Models[m.ModelID].Trees)
		numTrees.OnChanged(func(*ui.Spinbox) {
			windowData.Models[m.ModelID].Trees = numTrees.Value()
		})

		form1 := ui.NewForm()
		form1.SetPadded(true)
		hbox1.Append(form1, false)
		form1.Append("numTrees", numTrees, false)

		maxDepth := ui.NewSpinbox(0, 100)
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

	button := ui.NewButton("Import Config")
	button.OnClicked(func(*ui.Button) {
		filename := ui.OpenFile(mainwin)
		if filename != "" {
			file, _ := ioutil.ReadFile(filename)
			temp := UIWindow{}
			_ = json.Unmarshal([]byte(file), &temp)
			windowData = temp
			vbox.Delete(1)
			grid = generateFromState()
			vbox.Append(grid, false)
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
		vbox.Delete(1)
		grid = generateFromState()
		vbox.Append(grid, false)
	})
	msggrid.Append(button,
		3, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	button = ui.NewButton("Run Models")
	button.OnClicked(func(*ui.Button) {
		ui.MsgBoxError(mainwin,
			"This message box describes an error.",
			"More detailed information can be shown here.")
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

func main() {
	ui.Main(setupUI)
}
