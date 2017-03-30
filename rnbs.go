package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"localhost/minty/dtree"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cheggaaa/pb"
)

var old, new string
var source, dest *string

type Copy struct {
	Replace                            bool
	delimiter                          string
	Old, New, Message                  string
	Source, Dest, DestRoot, SourceRoot string
	bar                                *pb.ProgressBar
	wg                                 *sync.WaitGroup
}

func NewCopy(source, dest string, bar *pb.ProgressBar, wg *sync.WaitGroup) *Copy {
	return &Copy{
		bar:    bar,
		Source: source,
		Dest:   dest,
		wg:     wg,
	}
}

func (c *Copy) Copy() error {
	//	fmt.Println(c.Source)
	fl, err := ioutil.ReadFile(c.Source)
	ioutil.WriteFile(c.Dest, fl, 755)
	return err
}

// CopyFile will copy the files
func (c *Copy) CopyFile() (err error) {
	if c.Replace {
		err = newReplace(c.Old, c.New, c.Source, c.Dest).copy()
		c.bar.Increment()
		return
	}
	c.Copy()
	c.bar.Increment()
	return

}

// CopyFolder will recursively copy files from each folder if they exist
func (c *Copy) CopyFolder() (err error) {
	sourceinfo, err := os.Stat(c.Source)
	if err != nil {
		os.Exit(1)
	}

	c.Dest = folderName(c.Old, c.New, c.Source, c.SourceRoot, c.DestRoot)

	err = os.MkdirAll(c.Dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	folder, _ := os.Open(c.Source)

	objects, err := folder.Readdir(-1)
	sourcefile := c.Source
	destinationfile := c.Dest
	for _, obj := range objects {

		sourcefile := filepath.Join(sourcefile, obj.Name())
		destinationfile := filepath.Join(destinationfile, obj.Name())

		if obj.IsDir() {
			c.Source = sourcefile
			c.Dest = destinationfile
			err = c.CopyFolder()
			if err != nil {
				fmt.Println(err)
			}
			c.bar.Increment()
		} else {
			c.wg.Add(1)
			go func(sourcefile, destinationfile string) {
				c.Source = sourcefile
				c.Dest = destinationfile
				err = c.CopyFile()
				defer c.wg.Done()
				if err != nil {
					fmt.Println(err)

				}
			}(sourcefile, destinationfile)
			c.wg.Wait()
		}

	}

	return
}

func (c *Copy) dirExist() {
	_, err := os.Open(c.Dest)
	if !os.IsNotExist(err) {
		c.Message = fmt.Sprintf("Directory ` %v ` ALREADY EXISTS and has been overwritten. Please DELETE it and do it again to avoid any data corruption", c.Dest)
	}
}

func (c *Copy) getnamefromdir(path string) (name, root string) {
	if strings.Index(path, `\`) > -1 {
		c.delimiter = `\`
	}
	if path[len(path)-1:] == c.delimiter {
		path = path[:len(path)-1]

	}
	i := strings.LastIndex(path, c.delimiter)
	root = path[:i]
	name = path[(i + 1):]
	return
}

func (c *Copy) findOldNew(s ...string) {
	c.Old, c.SourceRoot = c.getnamefromdir(s[0])
	c.New, c.DestRoot = c.getnamefromdir(s[1])
}

func usage() {
	fmt.Println(" \n\n USAGE ")
	fmt.Printf(" \n\n %v  <<sourcePath>>  <<destDest>> ", os.Args[0])
	fmt.Println("\n\n EXAMPLE ")
	fmt.Printf("\n rnbs D:\\RN\\Boileplate D:\\RN\\NewProject ")
	fmt.Println("\n", `Paths can be relative`)
	fmt.Println(` Destination path NEED NOT exist. It will be created for you.`, "\n")
}

func flagparse(flags []string) {
	if DEV {
		flagDev()
		return
	}

	num := len(flags)
	switch num {
	case 2:
		source, dest = &flags[0], &flags[1]
	default:
		usage()
		os.Exit(1)
	}
}

func removeDir(paths ...string) {

	path := filepath.Join(paths...)
	err := dtree.Tree(path)
	if err != nil {
		fmt.Printf("\n %v doesnot exist...skipping", path)
		//fmt.Printf("\n Error while processing %v is \"%v\"", path, err)
		return
	}
	fmt.Printf("\n Removing >> %v  from  \"%v\" \n", dtree.DirFiles(), path)

	os.RemoveAll(path)
}

func tailer(c *Copy) {
	c.bar.FinishPrint(" New React Native project created at")
	fmt.Printf(" %v", filepath.Join(c.DestRoot, c.New))
	fmt.Println(`REMOVING '.gradle' path and 'app/build' path IF ANY`)
	removeDir(c.DestRoot, c.New, "android", ".gradle")
	removeDir(c.DestRoot, c.New, "android", "app", "build")
	removeDir(c.DestRoot, c.New, "ios", "build") //TODO: ios needs to be tested
	fmt.Println("\n\n DONE !!")
}

func main() {
	flag.Parse()
	flagparse(flag.Args())

	dPath, err := filepath.Abs(*dest)
	if err != nil {
		fmt.Printf("Could not find ABSOLUTE PATH for Destination %v", *dest)
		os.Exit(1)
	}
	*dest = dPath
	sPath, err := filepath.Abs(*source)
	if err != nil {
		fmt.Printf("Could not find ABSOLUTE PATH for Source %v", *source)
		os.Exit(1)
	}
	*source = sPath

	src, err := os.Stat(*source)
	if err != nil {
		fmt.Printf("Source file/folder \"%v\"  you provided does not exist. Please check your path", *source)
		os.Exit(1)
	}

	if !src.IsDir() {
		fmt.Println("Source is not a folder")
		os.Exit(1)
	}

	fmt.Printf("\n Scanning Source Directory %v ......", *source)
	err = dtree.Tree(*source)
	if err != nil {
		fmt.Printf("\n\nProblem with the source >> %v. Ended with an error >> %v\n\n", *source, err)
	}

	fmt.Println(dtree.DirFiles(), "\n")
	bar := pb.StartNew(int(dtree.Count()))

	var wg sync.WaitGroup

	c := NewCopy(*source, *dest, bar, &wg)
	c.dirExist()
	if c.Message != "" {
		fmt.Println("\n", c.Message)
	}
	c.Replace = true
	if c.Replace {
		c.findOldNew(*source, *dest)
	}
	c.bar.Increment()
	err = c.CopyFolder()

	if err != nil {
		fmt.Printf("\n Exiting with ERROR > ", err)
		os.Exit(1)
	}

	tailer(c)
}
