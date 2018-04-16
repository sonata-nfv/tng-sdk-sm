package helpers

import (
	"os"
	"io/ioutil"
	// "fmt"
	"io"
	"path/filepath"
	"errors"
	"strings"
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
		return errors.New("son-sm-template unreachable. Is TNG_SM_PWD set?") 
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