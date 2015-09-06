package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/seiflotfy/skizze/sketches/wrappers/topk/go-topk"
)

func main() {

	k := flag.Int("n", 500, "k")
	f := flag.String("f", "", "file to read")
	counts := flag.Bool("c", false, "each item has a count associated with it")

	flag.Parse()

	var r io.Reader

	if *f == "" {
		r = os.Stdin
	} else {
		var err error
		r, err = os.Open(*f)
		if err != nil {
			log.Fatal(err)
		}
	}

	tk := topk.New(*k)
	sc := bufio.NewScanner(r)

	for sc.Scan() {
		line := sc.Text()

		var count int
		var item string

		if *counts {
			fields := strings.Fields(line)
			cint, err := strconv.Atoi(fields[1])
			if err != nil {
				log.Println("failed to parse count: ", fields[1], ":", err)
				continue
			}
			item = fields[0]
			count = cint
		} else {
			item = line
			count = 1
		}

		tk.Insert(item, count)
	}

	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}

	for _, v := range tk.Keys() {
		fmt.Println(v.Key, v.Count, v.Error)
	}
}
