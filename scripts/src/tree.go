package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func printTree(root string, prefix string) {
	entries, err := os.ReadDir(root)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	for i, entry := range entries {
		isLast := i == len(entries)-1
		indent := "├── "
		subPrefix := "│   "
		if isLast {
			indent = "└── "
			subPrefix = "    "
		}

		fmt.Println(prefix + indent + entry.Name())

		if entry.IsDir() {
			printTree(filepath.Join(root, entry.Name()), prefix+subPrefix)
		}
	}
}

func main() {
	root := "." // Текущая директория
	fmt.Println(root)
	printTree(root, "")
}
