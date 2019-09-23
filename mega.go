package main

import (
	"github.com/manifoldco/promptui"
	"fmt"
	"github.com/t3rm1n4l/go-mega"
	"github.com/cheggaaa/pb/v3"
	"time"
	"os"
)

func auth() {
	config, err := parseConf()
	if err != nil {
		config = setConf()
		initConf(config)
	}

	err = m.Login(config.Username, config.Password)
	if err != nil {
		fmt.Printf("login failed!\ncheck your username&password!\nUsername: '%s'\nPassword: '%s'\n", config.Username, config.Password)
		return
	}
	fmt.Printf("login success!\nlogin as {%s}!\n", config.Username)
	root = m.FS.GetRoot()
}

func choose(parent, trees []Item) *Item {
	prompt := promptui.Select{
		Label:        "select file",
		Items:        trees,
		Templates:    chooseTemplate,
		HideSelected: true,
		Size:         10,
	}
	index, _, err := prompt.Run()
	if err != nil {
		os.Exit(0)
		return nil
	}
	node := trees[index]
	if node.Type != 0 {
		childList := getChildren(node.Node)
		return choose(trees, childList)
	}

	if node.Name == prev {
		if parent == nil {
			return choose(nil, getChildren(root))
		}
		return choose(nil, parent)
	}

	return &node
}

func operate(selectedFile *Item) {
	prompt := promptui.Select{
		Label:        "Select Operation",
		HideSelected: true,
		Items:        []string{"Download", "Delete", "Cancel"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	if result == "Download" {
		fmt.Printf("download %s\n", selectedFile.Name)
		downloadFile(selectedFile.Node, selectedFile.Name)
	} else if result == "Delete" {
		confirmed := confirmDelete()
		if confirmed {
			err := deleteFile(selectedFile.Node)
			if err != nil {
				fmt.Printf("delete %s failed\n%s", selectedFile.Name, err.Error())
			}
		}
	} else {

	}
}

func deleteFile(node *mega.Node) error {
	return m.Delete(node, false)
}

func downloadFile(node *mega.Node, name string) {
	progress := make(chan int)
	go m.DownloadFile(node, name, &progress)
	var totalSize int
	totalSize = int(node.GetSize())
	showProgress(totalSize, progress)
}

func uploadFile(fileName, filePath string, fileSize int64) {
	progress := make(chan int)
	go m.UploadFile(filePath, root, fileName, &progress)
	var totalSize int
	totalSize = int(fileSize)
	showProgress(totalSize, progress)
}

func confirmDelete() bool {
	prompt := promptui.Prompt{
		Label:     "Delete Resource",
		IsConfirm: true,
	}

	result, err := prompt.Run()

	if err != nil {
		return false
	}
	if result == "y" || result == "Y" {
		return true
	}
	return false
}

func getChildren(node *mega.Node) []Item {
	var childList []Item
	if node.GetType() != 2 {
		childList = append(childList, Item{Name: prev})
	}
	l, _ := m.FS.GetChildren(node)
	for _, n := range l {
		childList = append(childList, Item{
			n.GetName(),
			n.GetType(),
			getSize(n.GetSize()),
			n,
		})
	}
	return childList
}

func getSize(size int64) string {
	if size == 0 {
		return "0"
	}
	BSize := float64(size)
	KSize := BSize / 1024
	if KSize <= 1024 {
		return fmt.Sprintf("%.2fK", KSize)
	}
	MSize := KSize / 1024
	if MSize <= 1024 {
		return fmt.Sprintf("%.2fM", MSize)
	}
	GSize := MSize / 1024
	return fmt.Sprintf("%.2fG", GSize)
}

func showProgress(totalSize int, progress chan int) {
	progressTemplate := `{{string . "prefix" }} {{counters . }} {{ "bit" }} {{ bar . }} {{speed . | rndcolor }} {{percent .}} {{string . "my_green_string" | green}} {{rtime . "ETA %s"}} {{string . "my_blue_string" | blue}}`
	bar := pb.ProgressBarTemplate(progressTemplate).Start(totalSize)

	for {
		select {
		case <-progress:
			return
		default:
			bar.Add(<-progress)
			time.Sleep(time.Millisecond)
		}
		time.Sleep(time.Millisecond)
	}
	bar.Finish()
}
