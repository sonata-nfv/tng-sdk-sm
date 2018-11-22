// Copyright (c) 2015 SONATA-NFV, 2017 5GTANGO
// ALL RIGHTS RESERVED.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Neither the name of the SONATA-NFV, 5GTANGO
// nor the names of its contributors may be used to endorse or promote
// products derived from this software without specific prior written
// permission.

// This work has been performed in the framework of the SONATA project,
// funded by the European Commission under Grant number 671517 through
// the Horizon 2020 and 5G-PPP programmes. The authors would like to
// acknowledge the contributions of their colleagues of the SONATA
// partner consortium (www.sonata-nfv.eu).

// This work has been performed in the framework of the 5GTANGO project,
// funded by the European Commission under Grant number 761493 through
// the Horizon 2020 and 5G-PPP programmes. The authors would like to
// acknowledge the contributions of their colleagues of the 5GTANGO
// partner consortium (www.5gtango.eu).

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

var Version = "1.1"

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
		fmt.Printf("    generate       Generates a set of artefacts to be used as inputs when executing\n")
		fmt.Printf("                   specific managers\n")
		fmt.Printf("\n")
		fmt.Printf("usage: tng-sm new [OPTIONS] <specific manager name>\n\n")
		fmt.Printf("    --path         Path where new specific manager should be stored\n")
		fmt.Printf("    --type         Type of specific manager to be created: \"ssm\" or \"fsm\"\n")
		fmt.Printf("\n")
		fmt.Printf("usage: tng-sm delete [OPTIONS] <specific manager name>\n\n")
		fmt.Printf("    --path         Path where specific manager can be found\n")
		fmt.Printf("\n")
		fmt.Printf("usage: tng-sm execute [OPTIONS] <specific manager name>\n\n")
		fmt.Printf("    --path         Path where specific manager can be found\n")
		fmt.Printf("    --event        Event that needs to be executed: \"start\", \"stop\" or \"configure\"\n")
		fmt.Printf("    --payload      Payload for the execution\n")
		fmt.Printf("\n")
		fmt.Printf("usage: tng-sm generate [OPTIONS] <input file>\n\n")
		fmt.Printf("    --type         Type of specific manager the payload is for: \"fsm\" or \"ssm\"\n")
		fmt.Printf("    --input        Type of input file: \"descriptor\" or \"package\"\n")
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
	directory, err := helpers.CreateDirectory(name + *newTypePtr, *newPathPtr)

	// Copy the template to the new directory
	err = helpers.CopyTemplate(directory)

	if err != nil {
		fmt.Println(err)
		helpers.RemoveDir(directory)
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
	err = helpers.RemoveDir(dir)
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
	cmd_string := "import " + sm_name + "; import yaml; fsm=" + sm_name + "." + sm_name + "FSM(connect_to_broker=False);" + pyld_cmd + " event=fsm." + event_str + "(payload); print(\n\"Event output: \" + str(event))"
	cmd:= exec.Command("python3", "-c", cmd_string)
	cmd.Dir = code_dir
	out, err := cmd.CombinedOutput()

	if err != nil {
		globals.RedBold.Printf("\nError occured when executing command. Error:\n")
		fmt.Printf("%s\n", err)
		globals.RedBold.Printf("Python output:\n")
		fmt.Printf("%s\n", out)
		globals.RedBold.Printf("Python code:\n")
		fmt.Printf("%s\n", cmd_string)
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
	inputTypePtr := arg.String("input", "", "Type of input for the generation process: 'descriptor' or 'package'")
	managerTypePtr := arg.String("type", "", "Type of manager that generated content is for: 'fsm' or 'ssm'")

	arg.Parse(arg_list)	

	// Check if name was provided for output file
	if len(arg.Args()) != 1 {
		fmt.Printf("Generation failed. Provide input file.\n")
		os.Exit(1)
	} 

	// Check if input type is provided
	if *inputTypePtr == "" {
		fmt.Printf("Generation failed. Provide '--input' out of 'descriptor', 'package'.\n")
		os.Exit(1)
	}

	// Check if type is provided
	if *managerTypePtr == "" {
		fmt.Printf("Generation failed. Provide '--type' out of 'fsm', 'ssm'.\n")
		os.Exit(1)
	}

	// declaring
	inputFile := arg.Args()[0]
	extension := filepath.Ext(inputFile)
	outputName:= inputFile[0:len(inputFile)-len(extension)]
	outputDirectory := outputName + "_payloads"
	outputFiles := []string{}

	if exists,_ := helpers.Exists(inputFile); !exists {
		fmt.Printf("Provided input file does not exist.\n")
		os.Exit(1)
	}

	// create directory for output files
	helpers.CreateDirectory(outputDirectory, "")

	switch *inputTypePtr {
	case "package":

		// Unzip the package at temporary location
		tempLoc := "tng_sm_temp"
		_, err := helpers.Unzip(inputFile, tempLoc)
		if err != nil {
			fmt.Printf("Unzipping of package failed. Is input file really a 5GTANGO package?\n")
			os.Exit(1)
		}

		// process package descriptor
		pd_path := "TOSCA-Metadata/NAPD.yaml"
		total_path := filepath.Join(tempLoc, pd_path)
		pd, err := helpers.GetPd(total_path)
		if err != nil {
			fmt.Printf("No valid package descriptor in package.\n")
			os.Exit(1)
		}

		// // use package descriptor to get nsd and vnfd paths
		// nsd_path, err := helpers.GetNsdFromPackage(tempLoc, pd)
		// if err != nil {
		// 	fmt.Printf("Couldn't retrieve nsd path from package.\n")
		// 	os.Exit(1)			
		// }
		vnfd_paths, err := helpers.GetVnfdsFromPackage(tempLoc, pd)
		if err != nil {
			fmt.Printf("Couldn't retrieve vnfd paths from package.\n")
			os.Exit(1)			
		}

		// write the vnfrs and the start/stop payloads
		for _, vnfd_path := range vnfd_paths {
			vnfd_byte, err := helpers.ReadFile(vnfd_path)
			_, file := filepath.Split(vnfd_path)
			if err != nil {
				fmt.Printf("vnfd is not readable.\n")
				os.Exit(1)
			}
			vnfr_byte,_ := helpers.GenerateVnfrFromVnfd(vnfd_byte)		
			err = helpers.WriteFile(vnfr_byte, filepath.Join(outputDirectory, file + "_vnfr"))
			outputFiles = append(outputFiles, filepath.Join(outputDirectory, file + "_vnfr"))

			if err != nil {
				fmt.Printf("error when writing vnfr to file.\n")
				os.Exit(1)
			}

			output_byte,_ := helpers.GenerateStartStopOutput(vnfd_byte, vnfr_byte)
			err = helpers.WriteFile(output_byte, filepath.Join(outputDirectory, file + "_start_stop"))
			outputFiles = append(outputFiles, filepath.Join(outputDirectory, file + "_start_stop"))

			if err != nil {
				fmt.Printf("error when writing start/stop payload to file.\n")
				os.Exit(1)
			}
		}

		// write the configure payload
		config_byte,_ := helpers.GenerateConfigureOutput(vnfd_paths)
		err = helpers.WriteFile(config_byte, filepath.Join(outputDirectory, "ns_configure"))
		outputFiles = append(outputFiles, filepath.Join(outputDirectory, "ns_configure"))

		// cleanup temp directory
		err = helpers.RemoveDir(tempLoc)

	case "descriptor":
		desc_byte, err := helpers.ReadFile(inputFile)

		// Check if provided descriptor is readable
		if err != nil {
			fmt.Printf("Input file not readable.\n")
			os.Exit(1)
		}

		switch *managerTypePtr {
		case "fsm":
			vnfr_byte,err := helpers.GenerateVnfrFromVnfd(desc_byte)
			if err != nil {
				fmt.Printf("Descriptor processing failed. Is input file really a vnfd?\n")
				os.Exit(1)
			}

			// Generating vnfr
			err = helpers.WriteFile(vnfr_byte, filepath.Join(outputDirectory, outputName + "_vnfr"))
			outputFiles = append(outputFiles, filepath.Join(outputDirectory, outputName + "_vnfr"))

			if err != nil {
				fmt.Printf("error while writing vnfr to file.\n")
				os.Exit(1)
			}

			// Generating start/stop payload
			output_byte,_ := helpers.GenerateStartStopOutput(desc_byte, vnfr_byte)
			err = helpers.WriteFile(output_byte, filepath.Join(outputDirectory, outputName + "_start_stop"))
			outputFiles = append(outputFiles, filepath.Join(outputDirectory, outputName + "_start_stop"))

			if err != nil {
				fmt.Printf("error while writing start/stop payload to file.\n")
				os.Exit(1)
			}

		case "ssm":
			fmt.Printf("ssm payload generation not supported.\n")
		default:
			fmt.Printf(*managerTypePtr + " not supported as '--type'\n.")
			os.Exit(1)
		}

	default:
		fmt.Printf(*inputTypePtr + " not supported as '--input'.\n")
	}

	fmt.Printf("The following files were generated:\n\n")
	for _,filepath := range outputFiles {

		fmt.Printf("    %s\n", filepath)
	}
	fmt.Printf("\n")

	os.Exit(0)
}
