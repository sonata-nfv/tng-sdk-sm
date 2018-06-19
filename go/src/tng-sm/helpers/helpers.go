package helpers

import (
	"os"
	"io/ioutil"
	"fmt"
	"io"
	"path/filepath"
	"errors"
	"strings"
    "tng-sm/structs"
    "gopkg.in/yaml.v2"
	"github.com/nu7hatch/gouuid"
	)

func CreateDirectory(name, sm_type, path string) (directory string, err error) {

	// TODO: Do we need to overwrite if the directory already exists?
	directory_name := name + "-" + sm_type
	directory = filepath.Join(path, directory_name)
	os.Mkdir(directory, os.FileMode(0777))

	return
}	

func RemoveDirectory(name string) (err error) {

	//TODO: remove the directory at the path
	return
}

func CopyTemplate(sm_dir string) (err error) {

	// locate template location
	tmp_path := os.Getenv("TNG_SM_PWD")
	tmp_folder := "templates"
	tmp_dir := filepath.Join(tmp_path, tmp_folder)

	// Check if path exists
	_, err = os.Stat(tmp_dir)

	if err != nil {
		// Delete the directory
		os.Remove(sm_dir)
		return errors.New("son-sm-template unreachable. TNG_SM_PWD set?") 
	}

	// list content of directory
	files, err := ioutil.ReadDir(tmp_dir)

	if err != nil {
		return 
	}

	// Start copying
	for _, file := range files {

		sourcePath := filepath.Join(tmp_dir, file.Name())
		destinationPath := filepath.Join(sm_dir, file.Name())

		if file.IsDir(){
			CopyDir(sourcePath, destinationPath)
		} else {
			err = CopyFile(sourcePath, destinationPath)
			if err != nil {
				return
			}
		}
	}

	return
}

func ReadFile(name string) (data []byte, err error) {
	data, err = ioutil.ReadFile(name)
	if err != nil {
		err = errors.New("File unreadable")
	}

	return
}

func WriteFile(data []byte, name string) (err error) {

	err = ioutil.WriteFile(name, data, 0644)
	return
}

func CopyFile(source, destination string) (err error) {

	// Open the source file
	in, err := os.Open(source)
	defer in.Close()
	if err != nil {
		return
	}

	// Create the destination file
	out, err := os.Create(destination)
	defer out.Close()
	if err != nil {
		return
	}

	// Copy the file
	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	// Save to disk
	err = out.Sync()
	if err != nil {
		return
	}

	// Get file info of source
	info, err := os.Stat(source)
	if err != nil {
		return
	}

	// Enforce file info on destination
	err = os.Chmod(destination, info.Mode())
	if err != nil {
		return
	}

	return
}

func CopyDir(sourceDir, destinationDir string) (err error) {

	// create the destination directory
	os.Mkdir(destinationDir, os.FileMode(0777))

	// Loop over files and directories in source directory
	files, err := ioutil.ReadDir(sourceDir)

	for _,file := range files {

		sourcePath := filepath.Join(sourceDir, file.Name())
		destinationPath := filepath.Join(destinationDir, file.Name())

		if file.IsDir(){
			CopyDir(sourcePath, destinationPath)
		} else {
			err = CopyFile(sourcePath, destinationPath)
			if err != nil {
				return
			}
		}
	}

	return
}

func CustomiseTemplate(dir, name, sm_type string) (err error) {

	name = strings.ToLower(name)

	// rename the _template directory
	generic_dir := filepath.Join(dir, "template")
	specific_dir := filepath.Join(dir, name)
	err = os.Rename(generic_dir, specific_dir)
	if err != nil {
		return
	}

	// rename _template file
	generic_file := filepath.Join(specific_dir, sm_type + "_template.py")
	specific_file := filepath.Join(specific_dir, name + ".py")
	err = os.Rename(generic_file, specific_file)
	if err != nil {
		return
	}

	// clean-up obsolete files
	if strings.Compare(sm_type, "ssm") == 0 {
		os.Remove(filepath.Join(specific_dir, "fsm_template.py"))
		os.Remove(filepath.Join(specific_dir, "ssh.py"))
	}
	if strings.Compare(sm_type, "fsm") == 0 {
		os.Remove(filepath.Join(specific_dir, "ssm_template.py"))
	}

	// customise _template file
	err = ReplaceTagInFile("<name>", name, specific_file)
	if err != nil {
		return
	}

	// customise main file
	err = ReplaceTagInFile("<name>", name, filepath.Join(specific_dir, "__main__.py"))
	if err != nil {
		return
	}

	// customise Dockerfile file
	err = ReplaceTagInFile("<name>", name, filepath.Join(dir, "Dockerfile"))
	if err != nil {
		return
	}

	// customise Readme file
	err = ReplaceTagInFile("<name>", name, filepath.Join(dir, "README.md"))
	if err != nil {
		return
	}
	return

}

func ReplaceTagInFile(old_tag, new_tag, filepath string) (err error) {

    input, err := ioutil.ReadFile(filepath)
    if err != nil {
    	return
    }

    lines := strings.Split(string(input), "\n")

    for i, line := range lines {
            if strings.Contains(line, old_tag) {

            		new_line := strings.Replace(line, old_tag, new_tag, -1)
                    lines[i] = new_line 
            }
    }
    output := strings.Join(lines, "\n")
    err = ioutil.WriteFile(filepath, []byte(output), 0644)
    if err != nil {
            return
    }

    return
}

func RemoveContents(dir string) (err error) {
    d, err := os.Open(dir)
    if err != nil {
        return
    }
    defer d.Close()
    names, err := d.Readdirnames(-1)
    if err != nil {
        return
    }
    for _, name := range names {
        err = os.RemoveAll(filepath.Join(dir, name))
        if err != nil {
            return
        }
    }

    err = os.Remove(dir)
    if err != nil {
        return
    }

    return
}

func Exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return true, err
}


func GenerateNsrFromNsd(nsd_byte []byte) (nsr_byte []byte, err error) {

	return
}

func GenerateVnfrFromVnfd(vnfd_byte []byte) (vnfr_byte []byte, err error){

	//declaring
	vnfd := structs.Vnfd{}
    vnfr := structs.Vnfr{}
    u, err := uuid.NewV4()

	err = yaml.Unmarshal(vnfd_byte, &vnfd)
    if err != nil {
            fmt.Printf("Error parsing VNFD: %v\n", err)
            return
    }
    // fmt.Printf("--- t:\n%v\n\n", vnfd)

    // VNFD commons
    vnfr.DescriptorVersion = vnfd.DescriptorVersion
    vnfr.VirtualLinks = vnfd.VirtualLinks

    numberVdu := len(vnfd.VirtualDeploymentUnits)

    for i:= 0; i<numberVdu; i++ {
    	vnfdVdu := *vnfd.VirtualDeploymentUnits[i]
    	vnfrVdu := structs.VnfrVirtualDeploymentUnits{}

    	vnfrVdu.ResourceRequirements = vnfdVdu.ResourceRequirements
    	vnfrVdu.VmImage = vnfdVdu.VmImage
    	vnfrVdu.Id = vnfdVdu.Id
    	vnfrVdu.NumberOfInstances = 1
    	vnfrVdu.VduReference = vnfd.Name + ":" + vnfdVdu.Id

    	vnfc := structs.VnfcInstance{}
    	numberCp := len(vnfdVdu.ConnectionPoints)

    	for j:= 0; j<numberCp; j++ {
    		vnfdCp := *vnfdVdu.ConnectionPoints[j]
    		vnfrCp := structs.VnfrConnectionPoints{}

    		vnfrCp.Id = vnfdCp.Id
    		vnfrCp.Type = vnfdCp.Type

    		vnfrInterface := structs.VnfrInterface{}
    		vnfrInterface.Address = "<ENTER IP ADDRESS HERE>"
    		vnfrInterface.HardwareAddress = "<ENTER MAC ADDRESS HERE>"
    		vnfrInterface.Netmask = "<ENTER NETMASK HERE>"

    		vnfrCp.Interface = &vnfrInterface

    		vnfc.ConnectionPoints = append(vnfc.ConnectionPoints, &vnfrCp)
    	}

    	u, err = uuid.NewV4()
    	vnfc.VimId = u.String()
    	vnfc.Id = "0"

    	vnfrVdu.VnfcInstance = append(vnfrVdu.VnfcInstance, &vnfc)
    	vnfr.VirtualDeploymentUnits = append(vnfr.VirtualDeploymentUnits, &vnfrVdu)
    }

    // VNFR uniques
    vnfr.Version = "1"
    vnfr.Status = "normal operation"

    u, err = uuid.NewV4()
    vnfr.Id = u.String()

    u, err = uuid.NewV4()
    vnfr.DescriptorReference = u.String()

    vnfr_byte, err = yaml.Marshal(&vnfr)
	fmt.Printf("\n%s\n", vnfr)

    return
}

func GenerateStartStopOutput(vnfd_byte []byte, vnfr_byte []byte) (output_byte []byte, err error){

	output := structs.StartStop{}
	err = yaml.Unmarshal(vnfd_byte, &output.Vnfd)
	err = yaml.Unmarshal(vnfr_byte, &output.Vnfr)

    output_byte, err = yaml.Marshal(&output)

	return
}