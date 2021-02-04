package main

import (
	"log"
	"os"
	"testing"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime)
}

func TestHTTPForkJoin(t *testing.T) {

}
