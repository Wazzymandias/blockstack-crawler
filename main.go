package main

import "fmt"

func main() {
	run()
}

func run() error {
	ns, err := GetAllNamespaces()

	if err != nil {
		return err
	}

	for _, namespace := range ns {
		nsNames, err := GetNamespaceNames(namespace)

		if err != nil {
			return err
		}

		fmt.Println(nsNames)
	}

	return nil
}
