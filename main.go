// Command upfile serves a directory over HTTP so it can be browsed and
// downloaded from any device on the local network.
//
// Usage:
//
//	upfile [--port N] [dir]
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"slices"
)

func main() {
	fs := flag.NewFlagSet("upfile", flag.ExitOnError)
	port := fs.Int("port", 8000, "port to listen on")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: upfile [--port N] [dir]")
		fs.PrintDefaults()
	}

	// Allow flags both before and after the directory, e.g. `upfile . --port 9000`.
	fs.Parse(os.Args[1:])
	dir := "."
	if fs.NArg() > 0 {
		dir = fs.Arg(0)
		fs.Parse(fs.Args()[1:])
	}

	root, err := filepath.Abs(dir)
	if err != nil {
		log.Fatal(err)
	}
	if info, err := os.Stat(root); err != nil || !info.IsDir() {
		log.Fatalf("%s is not a directory", dir)
	}

	fmt.Printf("Serving %s\n\n", root)
	for _, host := range hosts() {
		fmt.Printf("  http://%s:%d\n", host, *port)
	}
	fmt.Println()

	handler := logRequests(http.FileServer(http.Dir(root)))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), handler))
}

// hosts returns the names and addresses this machine can be reached at.
func hosts() []string {
	list := []string{"localhost"}
	if name, err := os.Hostname(); err == nil {
		list = append(list, name)
	}
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			if ip := ipnet.IP.String(); !slices.Contains(list, ip) {
				list = append(list, ip)
			}
		}
	}
	return list
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
