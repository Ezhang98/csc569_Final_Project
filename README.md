# csc569_Final_Project
Distributed System for Tuning Hyperparameters of Neural Networks

Group Members: Evan Zhang, Alexander Garcia and Christina Monahan
PAXOS election, consensus, and recovery algorithms

Implemented a distributed application that can run multiple machine learning algorithms(building models and tuning hyperparameters) on different servers. 

Our system consists of two main components: a simple user interface for users to choose the models to be used for training with some starting hyperparameters, as well as the fault-tolerant distributed backend system that distributes training jobs to different processes.

- Using Go programming language for the native support for concurrency we use channels for all communication. This will mitigate the parallelism delays in building models. 
- Failure detection will be monitored via a heartbeat function using a Gossip protocol. 
- Manage membership using Paxos to replicate processes across the servers and will select the fastest process. 
- Customizable configuration file to specify which models and hyperparameters a user prefers to use. The user should be able to either create a new configuration using the client program or provide a configuration file for the distributed system to run. 
- The output of the application would be a list of models ranked by validation and accuracy.
- Progress bar illustrates running the distribution and training  of the Neural Networks. Upon completion of tasks, bar stops. 

Notes:
- Currently only the model Neural Network has been implemented
- When the UI first loads, in order to load parameters, you must open the model type dropdown and reselect Neural Network
- Import config works with any of the samples or test files in the configs directory
- On Mac, the results section cuts off but on windows it shows the full text. 


Running instructions
go run App.go 
    - Launches the User Interface to select Data Sets and set hyperparameters to run Neural Network