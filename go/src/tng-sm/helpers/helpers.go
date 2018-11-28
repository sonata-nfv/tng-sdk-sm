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

package helpers

import (
	"os"
	"io/ioutil"
	// "fmt"
	"io"
	"path/filepath"
	"errors"
	"strings"
    "tng-sm/structs"
    "gopkg.in/yaml.v2"
	"github.com/nu7hatch/gouuid"
	"archive/zip"
	)

func CreateDirectory(name, path string) (directory string, err error) {

	directory = filepath.Join(path, name)
	exists, err := Exists(directory)
	if !exists {
		os.Mkdir(directory, os.FileMode(0777))
	}

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
	err = ReplaceTagInFile("<type>", sm_type, filepath.Join(dir, "Dockerfile"))
	if err != nil {
		return
	}

	// customise Readme file
	err = ReplaceTagInFile("<name>", name, filepath.Join(dir, "README.md"))
	if err != nil {
		return
	}

	// customise setup file
	err = ReplaceTagInFile("<name>", name, filepath.Join(dir, "setup.py"))
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

func RemoveDir(dir string) (err error) {
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

    return
}

func GenerateStartStopOutput(vnfd_byte []byte, vnfr_byte []byte) (output_byte []byte, err error){

	output := structs.StartStop{}
	err = yaml.Unmarshal(vnfd_byte, &output.Vnfd)
	err = yaml.Unmarshal(vnfr_byte, &output.Vnfr)

    output_byte, err = yaml.Marshal(&output)

	return
}

func GenerateConfigureOutput(vnfds []string) (output_byte []byte, err error){

	output:= structs.Configure{}
	for _, vnfd_file := range vnfds {

		vnfd_byte,_ := ReadFile(vnfd_file)
		vnfd := structs.Vnfd{}
		err = yaml.Unmarshal(vnfd_byte, &vnfd)
		output.Vnfds = append(output.Vnfds, &vnfd)

	    vnfr := structs.Vnfr{}
	    vnfr_byte,_ := GenerateVnfrFromVnfd(vnfd_byte)
		err = yaml.Unmarshal(vnfr_byte, &vnfr)
		output.Vnfrs = append(output.Vnfrs, &vnfr)
	}

	output_byte, err = yaml.Marshal(&output)
	return
}

func Unzip(src string, dest string) ([]string, error) {

    var filenames []string

    r, err := zip.OpenReader(src)
    if err != nil {
        return filenames, err
    }
    defer r.Close()

    for _, f := range r.File {

        rc, err := f.Open()
        if err != nil {
            return filenames, err
        }
        defer rc.Close()

        // Store filename/path for returning and using later on
        fpath := filepath.Join(dest, f.Name)

        if f.FileInfo().IsDir() {

            // Make Folder
            os.MkdirAll(fpath, os.ModePerm)

        } else {

	        filenames = append(filenames, fpath)

            // Make File
            if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
                return filenames, err
            }

            outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
            if err != nil {
                return filenames, err
            }

            _, err = io.Copy(outFile, rc)

            // Close the file without defer to close before next iteration of loop
            outFile.Close()

            if err != nil {
                return filenames, err
            }

        }
    }
    return filenames, nil
}


func GetPd(filepath string) (pd structs.PackageDescriptor, err error) {

		data, err := ReadFile(filepath)
		if err != nil {
			return
		}

		err = yaml.Unmarshal(data, &pd)
		return
}


func GetVnfdsFromPackage(dir string, pd structs.PackageDescriptor) (filepath_out []string, err error){

	for _, artefact := range pd.PackageContent {

		content := *artefact
		type_seg := strings.Split(content.ContentType, ".")

		if type_seg[len(type_seg) -1] == "vnfd" && type_seg[len(type_seg) -2] == "5gtango" {
			filepath_out = append(filepath_out, filepath.Join(dir, content.Source))
		}
	}

	if len(filepath_out) < 1 {
		err = errors.New("No vnfd in package")
	}

	return
} 

func GetNsdFromPackage(dir string, pd structs.PackageDescriptor) (filepath_out string, err error){

	for _, artefact := range pd.PackageContent {

		content := *artefact
		type_seg := strings.Split(content.ContentType, ".")

		if type_seg[len(type_seg) -1] == "nsd" && type_seg[len(type_seg) -2] == "5gtango" {
			filepath_out = filepath.Join(dir, content.Source)
			return
		}
	}
	err = errors.New("No nsd in package")
	return
} 
