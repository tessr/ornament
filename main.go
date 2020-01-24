package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/tendermint/iavl"
	tmdb "github.com/tendermint/tm-db"
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
		// TODO: check for the existence of keys
		log.Println("unimplemented")

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

	log.Printf("added <%s, %s> to create tree #%d with hash %s / %x", key, value, versionNumber, hashToEmoji(hash), hash)
	return tree, nil
}

func printShape(tree *iavl.MutableTree) {
	shape := tree.RenderShape("  ", nodeEncoder)
	fmt.Println(strings.Join(shape, "\n"))
}

func nodeEncoder(id []byte, depth int, isLeaf bool) string {
	prefix := fmt.Sprintf("-%d ", depth)
	if isLeaf {
		prefix = fmt.Sprintf("\t*%d ", depth)
	}
	if len(id) == 0 {
		return fmt.Sprintf("%s<nil>", prefix)
	}
	return fmt.Sprintf("%s%s", prefix, encodeID(id))
}

// casts to a string if it is printable ascii, emoji-hex-encodes otherwise
func encodeID(id []byte) string {
	for _, b := range id {
		if b < 0x20 || b >= 0x80 {
			return hashToEmoji(id)
		}
	}
	return string(id)
}

func hashToEmoji(id []byte) string {
	shorty := hex.EncodeToString(id[:7])
	dec, _ := strconv.ParseInt(shorty, 16, 64)
	n := int(dec) % len(emoji)
	return fmt.Sprintf("%s %s", emoji[n], shorty)
}
