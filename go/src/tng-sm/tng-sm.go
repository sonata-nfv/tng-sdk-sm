package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "tng-sm/helpers"
    "strings"
    "os/exec"
)

var Version = "0.1"

func main() {

	// root flags
	versionPtr := flag.Bool("version", false, "Returns the application version")

	// Subcommands definition
	newCommand := flag.NewFlagSet("new", flag.ExitOnError)
	deleteCommand := flag.NewFlagSet("delete", flag.ExitOnError)
	executeCommand := flag.NewFlagSet("execute", flag.ExitOnError)

	// Setting flags for ssm subsubcommand
	newPathPtr := newCommand.String("path", "", "Path where new specific manager should be stored")
	newTypePtr := newCommand.String("type", "", "One of ssm|fsm")

	deletePathPtr := deleteCommand.String("path", "", "Path where specific manager to be deleted can be found")

	executePathPtr := executeCommand.String("path", "", "Path where specific manager to be executed can be found")
	executeEventPtr := executeCommand.String("event", "", "Event that needs to be executed")
	executePayloadPtr := executeCommand.String("payload", "", "Payload for the execution")

	// Check wether subcommand is provided
	if len(os.Args) < 2 {
		fmt.Println("Please provide a flag or a subcommand\n")
		os.Exit(1)
	}

	// Analyse first subcommand
	switch os.Args[1] {
	case "new":
		newCommand.Parse(os.Args[2:])
	case "delete":
		deleteCommand.Parse(os.Args[2:])
	case "execute":
		executeCommand.Parse(os.Args[2:])
	default:
		flag.Parse()
		if *versionPtr == true {
			fmt.Printf("The version is: %v\n", Version)
			os.Exit(0)
		} else {
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	if newCommand.Parsed() {

		// Check if the name of the ssm or fsm was passed as an argument
		if len(newCommand.Args()) != 1 {
			fmt.Printf("Provide one argument for the name\n")
			os.Exit(1)
		} 

		// Extract the name of the specific manager
		name := newCommand.Args()[0]

		// Check if a path was provided
		if *newPathPtr != "" {
			// Evaluate if the provided path exists
			_, err := os.Stat(*newPathPtr)
			if err != nil {
				fmt.Printf("Provided path does not exist\n")
				os.Exit(1)
			}
		} 

		// Check if the specific manager type was provided
		if *newTypePtr == "" {
			fmt.Printf("Please provide a type, fsm or ssm\n")
			os.Exit(1)
		}

		// All flags are correct, creating new ssm/fsm

		// Check if directory already exists
		dir := filepath.Join(*newPathPtr, name + "-" + *newTypePtr)
		dir_exists, _ := helpers.Exists(dir)
		if dir_exists {
			if *newPathPtr != "" {
				fmt.Printf("\n    Specific manager already exists on that location,\n    use \"tng-sm delete --path %s %s\" to delete it\n\n", *newPathPtr, name + "-" + *newTypePtr)
			} else {
			fmt.Printf("\n    Specific manager already exists on that location,\n    use \"tng-sm delete %s\" to delete it\n\n", name + "-" + *newTypePtr)
			}
			os.Exit(1)
		}

		// Create the directory
		directory, err := helpers.CreateDirectory(name, *newTypePtr, *newPathPtr)

		// Copy the template to the new directory
		err = helpers.CopyTemplate(directory)

		if err != nil {
			fmt.Println(err)
			helpers.RemoveDirectory(directory)
			os.Exit(1)
		}
		// Customise the template
		err = helpers.CustomiseTemplate(directory, name, *newTypePtr)	

		fmt.Printf("Specific manager created\n")
	}

	if deleteCommand.Parsed() {

		// Check if the name of the ssm or fsm was passed as an argument
		if len(deleteCommand.Args()) != 1 {
			fmt.Printf("Provide name of specific manager to delete\n")
			os.Exit(1)
		} 

		// Extract the name of the specific manager
		name := deleteCommand.Args()[0]

		// Check if a path was provided
		if *deletePathPtr != "" {
			// Evaluate if the provided path exists
			_, err := os.Stat(*deletePathPtr)
			if err != nil {
				fmt.Printf("Provided path does not exist\n")
				os.Exit(1)
			}
		}

		// Check if specific manager exists
		dir := filepath.Join(*deletePathPtr, name)

		_, err := os.Stat(dir)

		if err != nil {
			fmt.Printf("No specific manager with name \"%s\" at provided location\n", name)
			os.Exit(1)
		}

		// Check if the directory is actually a specific manager
		dockerfile_path := filepath.Join(dir, "Dockerfile")
		setupfile_path := filepath.Join(dir, "setup.py")
		docker_exists, _ := helpers.Exists(dockerfile_path)
		setup_exists, _ := helpers.Exists(setupfile_path)

		if !(docker_exists && setup_exists) {
			fmt.Printf("%s is not a specific manager\n", dir)
			os.Exit(1)
		}

		// Delete the specific manager
		err = helpers.RemoveContents(dir)
		if err != nil {
			fmt.Printf("Specific manager deletion failed\n")
			os.Exit(1)
		}

		fmt.Printf("Specific manager successfully removed\n")
		os.Exit(0)
	}

	if executeCommand.Parsed() {

		// Check if the name of the ssm or fsm was passed as an argument
		if len(executeCommand.Args()) != 1 {
			fmt.Printf("Provide name of specific manager to execute\n")
			os.Exit(1)
		} 

		// Extract the name of the specific manager
		name := executeCommand.Args()[0]

		// Check if a path was provided
		if *executePathPtr != "" {
			// Evaluate if the provided path exists
			_, err := os.Stat(*executePathPtr)
			if err != nil {
				fmt.Printf("Provided path does not exist\n")
				os.Exit(1)
			}
		}

		// Check if specific manager exists
		dir := filepath.Join(*executePathPtr, name)

		_, err := os.Stat(dir)

		if err != nil {
			fmt.Printf("No specific manager with name \"%s\" at provided location\n", name)
			os.Exit(1)
		}

		if *executeEventPtr == "" {
			fmt.Printf("Provide an event to be executed through the \"--event\" flag.\n")
			os.Exit(1)
		}

		// build path to code
		sm_name := strings.Split(name, "-")[0]
		code_dir := filepath.Join(dir, sm_name)

		// If payload was provided, move it to execution dir
		pyld_cmd := " payload = None;"
		if *executePayloadPtr != "" {
			err = helpers.CopyFile(*executePayloadPtr, filepath.Join(code_dir, *executePayloadPtr))
			// fmt.Printf("%s\n", err)
			pyld_cmd = " payload = yaml.load(open(\"" + *executePayloadPtr + "\", 'rb'));"
		} 

		fmt.Printf("%s\n", pyld_cmd)

		// Execute the command and capture the output
		cmd_string := "import " + sm_name + "; import yaml; fsm=" + sm_name + "." + sm_name + "FSM();" + pyld_cmd + " event=fsm.start_event(payload); print(\n\"Event output: \" + str(event))"
		cmd:= exec.Command("python3", "-c", cmd_string)
		cmd.Dir = code_dir 
		out, err := cmd.CombinedOutput()

		if err != nil {
			fmt.Printf("Error occured when executing command: \n")
			fmt.Printf("\n%s\n", out)
			os.Exit(1)
		} else {
			fmt.Printf("Command executed successfully: \n")
			fmt.Printf("\n%s\n", out)			
		}

		// analyse output

		// Clean-up: remove payload from exeuction dir
	}
}
