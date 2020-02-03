package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tendermint/iavl"
	tmdb "github.com/tendermint/tm-db"
	"github.com/tessr/hmoj"
)

var (
	dataDir = "dir"
)

func main() {

	args := os.Args[1:]
	if len(args) < 1 || (args[0] != "wipe" && args[0] != "add" && args[0] != "check" && args[0] != "print") {
		fmt.Fprintln(os.Stderr, "Usage: ornament <wipe|add|check|print>")
		os.Exit(1)
	}

	switch args[0] {
	case "wipe":
		err := wipe()
		if err != nil {
			panic(err)
		}

	case "add":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: ornament add <key> <value>")
			os.Exit(1)
		}
		key, value := args[1], args[2]
		tree, err := add(key, value)
		if err != nil {
			panic(err)
		}
		printShape(tree)

	case "check":
		// TODO: support versions
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: ornament check <key>")
			os.Exit(1)
		}

		key := args[1]

		tree, err := getTree()
		if err != nil {
			panic(err)
		}

		v, err := checkKey(tree, key)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		fmt.Printf("got value %s for key %s\n", v, key)

	case "print":
		// TODO: support versions
		tree, err := getTree()
		if err != nil {
			panic(err)
		}
		printShape(tree)
	}
}

func wipe() error {
	if tmdb.FileExists(dataDir) {
		log.Println("removing old directory")
		os.RemoveAll(dataDir)
		return nil
	}
	return errors.New("db does not exist")
}

func getTree() (*iavl.MutableTree, error) {
	db, err := tmdb.NewGoLevelDB("ornament", dataDir)
	tree, err := iavl.NewMutableTree(db, 128)
	if err != nil {
		return nil, err
	}

	ver, err := tree.LoadVersion(0)
	if err != nil {
		return nil, err
	}

	fmt.Printf("got version: %d\n", ver)
	return tree, nil
}

func add(key, value string) (*iavl.MutableTree, error) {
	tree, err := getTree()
	if err != nil {
		return nil, err
	}
	tree.Set([]byte(key), []byte(value))
	hash, versionNumber, err := tree.SaveVersion()
	if err != nil {
		return nil, err
	}

	fmt.Printf("added <%s, %s> to create tree #%d with hash %s / %x", key, value, versionNumber, hmoj.Encode(hash), hash)
	return tree, nil
}

func checkKey(tree *iavl.MutableTree, key string) (string, error) {
	val, proof, err := tree.GetVersionedWithProof([]byte(key), tree.Version())
	if err != nil {
		return "", err
	}

	err = proof.Verify(tree.Hash())
	if err != nil {
		// TODO: expand error message to explain that proof is invalid?
		log.Println("invalid proof")
		return "", err
	}

	err = proof.VerifyItem([]byte(key), val)
	if err != nil {
		log.Println("unable to verify item")
		return "", err
	}

	return string(val), nil
}

func printShape(tree *iavl.MutableTree) {
	shape := tree.RenderShape("  ", nodeEncoder)
	fmt.Println(strings.Join(shape, "\n"))
}
