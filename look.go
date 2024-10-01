package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func init() {
}

func main() {
	var fn string
	var key string
	var rlen int
	flag.StringVar(&fn, "file", "/usr/share/dict/words", "name of sorted file to search")
	flag.StringVar(&key, "key", "", "search key")
	flag.IntVar(&rlen, "rlen", 0, "fixed length record length")
	pfold := flag.Bool("fold", false, "fold case")
	flag.Parse()
	if key == "" {
		fmt.Println("key required")
		return
	}
	if rlen != 0 {
		flsearch(fn, key, rlen)
	} else {
		search(fn, key, *pfold)
	}
}

func search(fn string, key string, fold bool) {
	f, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}

	var line, cline string
	var lo, mid, hi int64
	var n int
	st, _ := os.Stat(fn)
	hi = st.Size()
	const bsz int64 = 1 << 12
	var ba [bsz]byte
	var buf []byte = ba[0:]

	if fold {
		key = strings.ToLower(key)
	}
	var found int64
	var match bool
	for {
		mid = lo + (hi-lo)/2
		//log.Println("Looping", lo, mid, hi, hi-lo)
		n, err = f.ReadAt(buf, mid)
		if err != nil && !errors.Is(err, io.EOF) && n == 0 {
			log.Printf("ReadAt %v %v %v", mid, n, err)
			log.Fatal(err)
		}
		br := bytes.NewBuffer(buf[:n])

		if (hi - lo) < bsz {
			//log.Println("linear")

			n, err = f.ReadAt(buf, lo)
			if err != nil && !errors.Is(err, io.EOF) && n == 0 {
				log.Printf("ReadAt %v %v %v", lo, n, err)
				log.Fatal(err)
			}
			br := bytes.NewBuffer(buf[:n])

			curo := lo // current offset
			for {
				//log.Println("lo close to hi")
				line, err = br.ReadString('\n')
				if err != nil {
					if errors.Is(err, io.EOF) {
						return
					}
					log.Fatal(err)
				}
				if fold {
					cline = strings.ToLower(line)
				} else {
					cline = line
				}
				//log.Println("ReadString false", cline)
				curo += int64(len(line))

				if strings.HasPrefix(cline, key) {
					fmt.Print(line)
					match = true
					found = curo
					break
				}
				continue
			}

			if match == true {
				//log.Println("true")

				n, err = f.ReadAt(buf, found)
				if err != nil && !errors.Is(err, io.EOF) && n == 0 {
					log.Printf("ReadAt %v %v %v", lo, n, err)
					log.Fatal(err)
				}
				br := bytes.NewBuffer(buf[:n])

				for {
					line, err = br.ReadString('\n')
					if err != nil {
						if errors.Is(err, io.EOF) {
							return
						}
						log.Fatal(err)
					}
					if fold {
						cline = strings.ToLower(line)
					} else {
						cline = line
					}
					//log.Println("ReadString true", cline)
					if strings.HasPrefix(cline, key) {
						fmt.Print(line)
					} else {
						return
					}
				}
			}
			//return
		}

		line, err = br.ReadString('\n') // partial line?
		mid += int64(len(line))
		line, err = br.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			log.Fatal(err)
		}
		if fold {
			cline = strings.ToLower(line)
		} else {
			cline = line
		}
		//log.Println("ReadString binary", cline)

		if strings.HasPrefix(cline, key) {
			//log.Println("HasPrefix", cline, key)
			found = mid
			match = true
			hi = mid
			continue
		}
		if key < cline {
			//log.Println(key, "<", cline)
			hi = mid
			continue
		}
		if key > cline {
			//log.Println(key, ">", line)
			lo = mid
			continue
		}
	}
}

func flsearch(fn string, key string, reclen int) {
	f, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}

	var lo, mid, hi int64
	var n int
	st, _ := os.Stat(fn)
	hi = st.Size()

	var rsz int64 = int64(reclen) * 1 << 12
	var buf = make([]byte, rsz)
	var rbuf = make([]byte, reclen)

	var found int64
	var match bool

	for {
		mid = lo + (hi-lo)/2
		mid = mid - mid%int64(reclen)

		n, err = f.ReadAt(buf, mid)
		if err != nil && !errors.Is(err, io.EOF) && n == 0 {
			log.Printf("ReadAt %v %v %v", mid, n, err)
			log.Fatal(err)
		}
		br := bytes.NewBuffer(buf[:n])
		nr := bufio.NewReader(br)

		if (hi - lo) < int64(n) {
			n, err = f.ReadAt(buf, lo)
			if err != nil && !errors.Is(err, io.EOF) && n == 0 {
				log.Printf("ReadAt %v %v %v", lo, n, err)
				log.Fatal(err)
			}
			br := bytes.NewBuffer(buf[:n])
			nr := bufio.NewReader(br)
			curo := lo
			for {
				_, err := nr.Read(rbuf)
				if err != nil {
					log.Fatal(err)
				}
				curo += int64(reclen)
				if strings.HasPrefix(string(rbuf), key) {
					if strings.HasSuffix(string(rbuf), "\n") {
						fmt.Print(string(rbuf))
					} else {
						fmt.Println(string(rbuf))
					}
					match = true
					found = curo
					break
				}
				continue
			}
		}

		if match == true {
			n, err = f.ReadAt(buf, found)
			if err != nil && !errors.Is(err, io.EOF) && n == 0 {
				fmt.Print("ReadAt", lo, n, err)
				log.Fatal(err)
			}
			br := bytes.NewBuffer(buf[:n])
			nr := bufio.NewReader(br)
			for {
				_, err := nr.Read(rbuf)
				if err != nil {
					log.Fatal(err)
				}
				if strings.HasPrefix(string(rbuf), key) {
					if strings.HasSuffix(string(rbuf), "\n") {
						fmt.Print(string(rbuf))
					} else {
						fmt.Println(string(rbuf))
					}
				} else {
					return
				}
			}

		}

		_, err := nr.Read(rbuf)
		if err != nil {
			log.Fatal(err)
		}
		mid += int64(reclen)
		if strings.HasPrefix(string(buf), key) {
			found = mid
			match = true
			hi = mid
			continue
		}
		if key < string(buf) {
			hi = mid
			continue
		}
		if key > string(buf) {
			lo = mid
			continue
		}

	}
}
