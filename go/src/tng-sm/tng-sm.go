package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "tng-sm/globals"
    "tng-sm/helpers"
    "strings"
    "os/exec"
)

var Version = "1.0"

func main() {

	UsageMessage(flag.CommandLine)

	// root flags
	versionPtr := flag.Bool("version", false, "Returns the application version")

	// Subcommands definition
	newCommand := flag.NewFlagSet("new", flag.ExitOnError)
	deleteCommand := flag.NewFlagSet("delete", flag.ExitOnError)
	executeCommand := flag.NewFlagSet("execute", flag.ExitOnError)
	generateCommand := flag.NewFlagSet("generate", flag.ExitOnError)

	UsageMessage(newCommand)
	UsageMessage(deleteCommand)
	UsageMessage(executeCommand)
	UsageMessage(generateCommand)

	// Check wether subcommand is provided
	if len(os.Args) < 2 {
		fmt.Println("Please provide a flag or a subcommand\n")
		os.Exit(1)
	}

	// Analyse first subcommand
	switch os.Args[1] {
	case "new":
		HandleNewArg(newCommand, os.Args[2:])
	case "delete":
		HandleDeleteArg(deleteCommand, os.Args[2:])
	case "execute":
		HandleExecuteArg(executeCommand, os.Args[2:])
	case "generate":
		HandleGenerateArg(generateCommand, os.Args[2:])
	default:
		flag.Parse()
		if *versionPtr == true {
			fmt.Printf("tng-sm version %v\n", Version)
			os.Exit(0)
		} else {
			flag.PrintDefaults()
			os.Exit(1)
		}
	}
}

// Usage message
func UsageMessage(f *flag.FlagSet){

	f.Usage = func () {
		fmt.Printf("usage: tng-sm [--version] [--help]\n\n")
		fmt.Printf("These are the subcommands for tng-sm:\n\n")
		fmt.Printf("    new            Create a new specific manager\n")
		fmt.Printf("    delete         Delete an existing specific manager\n")
		fmt.Printf("    execute        Execute an event of a specific manager\n")
		fmt.Printf("    generate       Generate artefacts to be used when executing specific managers\n")
		fmt.Printf("\n")
		fmt.Printf("usage: tng-sm new <specific manager name>\n\n")
		fmt.Printf("    --path         Path where new specific manager should be stored\n")
		fmt.Printf("    --type         Type of specific manager to be created: \"ssm\" or \"fsm\"\n")
		fmt.Printf("\n")
		fmt.Printf("usage: tng-sm delete <specific manager name>\n\n")
		fmt.Printf("    --path         Path where specific manager can be found\n")
		fmt.Printf("\n")
		fmt.Printf("usage: tng-sm execute <specific manager name>\n\n")
		fmt.Printf("    --path         Path where specific manager can be found\n")
		fmt.Printf("    --event        Event that needs to be executed: \"start\", \"stop\" or \"configure\"\n")
		fmt.Printf("    --payload      Payload for the execution\n")
		fmt.Printf("\n")
		fmt.Printf("usage: tng-sm generate <name output file>\n\n")
		fmt.Printf("    --type         Type of payload to be generated: \"vnfr\" or \"nsr\"\n")
		fmt.Printf("    --descriptor   File that serves as input for generation, should be a vnfd or nsd\n")
		fmt.Printf("\n")
	}

}

// This function handles the new argument
func HandleNewArg(arg *flag.FlagSet, arg_list []string) () {

	// Setting flags for ssm subsubcommand
	newPathPtr := arg.String("path", "", "Path where new specific manager should be stored")
	newTypePtr := arg.String("type", "", "One of ssm|fsm")

	// parse the argument list
	arg.Parse(arg_list)

	// Check if the name of the ssm or fsm was passed as an argument
	if len(arg.Args()) != 1 {
		fmt.Printf("Provide one argument for the name\n")
		os.Exit(1)
	} 

	// Extract the name of the specific manager
	name := arg.Args()[0]

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
	os.Exit(0)

}	

// This function handles the delete argument
func HandleDeleteArg(arg *flag.FlagSet, arg_list []string) () {

	// Setting flags for ssm subsubcommand
	deletePathPtr := arg.String("path", "", "Path where specific manager to be deleted can be found")

	// parse the argument list
	arg.Parse(arg_list)

	// Check if the name of the ssm or fsm was passed as an argument
	if len(arg.Args()) != 1 {
		fmt.Printf("Provide name of specific manager to delete\n")
		os.Exit(1)
	} 

	// Extract the name of the specific manager
	name := arg.Args()[0]

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

// This function handles the execute argument
func HandleExecuteArg(arg *flag.FlagSet, arg_list []string) () {

	// Setting flags for ssm subsubcommand
	executePathPtr := arg.String("path", "", "Path where specific manager to be executed can be found")
	executeEventPtr := arg.String("event", "", "Event that needs to be executed")
	executePayloadPtr := arg.String("payload", "", "Payload for the execution")

	// parse the argument list
	arg.Parse(arg_list)

	// Check if the name of the ssm or fsm was passed as an argument
	if len(arg.Args()) != 1 {
		fmt.Printf("Provide name of specific manager to execute\n")
		os.Exit(1)
	} 

	// Extract the name of the specific manager
	name := arg.Args()[0]

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
	event_str := *executeEventPtr + "_event"

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

	// Execute the command and capture the output
	cmd_string := "import " + sm_name + "; import yaml; fsm=" + sm_name + "." + sm_name + "FSM();" + pyld_cmd + " event=fsm." + event_str + "(payload); print(\n\"Event output: \" + str(event))"
	cmd:= exec.Command("python3", "-c", cmd_string)
	cmd.Dir = code_dir 
	out, err := cmd.CombinedOutput()

	if err != nil {
		globals.RedBold.Printf("\nError occured when executing command. Python output: \n")
		fmt.Printf("\n%s\n", out)
		os.Exit(1)
	} else {
		globals.GreenBold.Printf("\nCommand executed successfully. Python output: \n")
		fmt.Printf("\n%s\n", out)			
	}

	// analyse output

	// Clean-up: remove payload from exeuction dir
}

// This function handles the generate argument
func HandleGenerateArg(arg *flag.FlagSet, arg_list []string) () {

	// Define possible flagsets
	payloadTypePtr := arg.String("type", "", "Type of payload to be generated: 'vnfr' or 'nsr'")
	payloadDescriptorPtr := arg.String("descriptor", "", "Descriptor used to generate the payload")

	arg.Parse(arg_list)	

	// Check if name was provided for output file
	if len(arg.Args()) != 1 {
		fmt.Printf("Provide name for output file\n")
		os.Exit(1)
	} 

	name := arg.Args()[0]

	// Check if type of payload that is to be generated is provided
	if *payloadTypePtr == "" {
		arg.PrintDefaults()
		os.Exit(1)
	}

	// Check if descriptor is provided
	if *payloadDescriptorPtr == "" {
		arg.PrintDefaults()
		os.Exit(1)
	} 

	data, err := helpers.ReadFile(*payloadDescriptorPtr)

	// Check if provided descriptor is readable
	if err != nil {
		fmt.Printf("Provided descriptor is not readable.\n")
		os.Exit(1)
	}

	switch *payloadTypePtr {
	case "nsr":
		// nsr,_ := helpers.GenerateNsrFromNsd(data)
		fmt.Printf("Functionality not yet supported.\n")
	case "vnfr":
		vnfr,_ := helpers.GenerateVnfrFromVnfd(data)
		output,_ := helpers.GenerateStartStopOutput(data, vnfr)
		err = helpers.WriteFile(output, name)

		if err != nil {
			fmt.Printf("error when writing file: %s", err)
		}
	default:
		fmt.Printf(*payloadTypePtr + " not supported as '--type'.\n")
	}
	

	os.Exit(0)
}
