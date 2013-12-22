package main

import (
	"fmt"
	"github.com/tpjg/goriakpbc"
  "github.com/bitly/go-simplejson"
	"net/http"
  "html/template"
  "encoding/json"
  "log"
	_ "os"
)

var (
	riakUrl        = "127.0.0.1"
	riakPort       = "8087"
	riakMainBucket = "riaktor"
  httpPort       = "8888"
)

type Soda struct {
	SugarFree   bool   `riak:"sugar_free"`
	Flavor      string `riak:"flavor"`
	Brand       string `riak:"brand"`
	riak.Model
}

type Obj struct {
  Type string `riak:"type"`
  Data map[string]string `riak:"data"`
  riak.Model
}

func (o *Obj) Id() string {
  return o.Data["id"]
}

func RiakTransaction(f func(*riak.Client)) {
	client := riak.New(fmt.Sprintf("%s:%s", riakUrl, riakPort))
	err := client.Connect()
	if err != nil {
		fmt.Println("Cannot connect, is Riak running?")
		return
	}
	f(client) // run transaction on client
	client.Close()
}

func set(o []interface{}, client riak.Client) {
	err := client.Connect()
	if err != nil {
		fmt.Println("Error finding bucket.")
	}
	var soda *Soda
  soda = &Soda{SugarFree: true, Flavor: "cola", Brand: "Coca-Cola"} 
  err = client.New("sodas", "test_soda", soda)

  // savedO = bucket.new()
}

func update(name string, client riak.Client) {
	err := client.Connect()
  if err != nil {
    fmt.Println("Error finding bucket.")
  }
  var soda Soda
  err = client.Load("sodas", "test_soda", &soda)
  soda.Flavor = "lemon_lime"
  soda.Save()
}

func clone(name string, newName string, client riak.Client) {
  err := client.Connect()
  if err != nil {
    fmt.Println("Error finding bucket.")
  }
  var soda Soda
  err = client.Load("sodas", "test_soda", &soda)
  // dev.Description = "something else"
  soda.SaveAs(newName)
}

func insertTestObject(bucket *riak.Bucket) {
	obj := bucket.New("testobj")
	obj.ContentType = "application/json"
	obj.Data = []byte("{'name': 'Bob'}")
	obj.Store()
}

func handler(w http.ResponseWriter, r *http.Request) {
    // fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
    t, _ := template.ParseFiles("templates/index.html")
    t.Execute(w, nil)
}

type Flash struct {
  Message string
}

func newHandler(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  log.Println(r.Form)
  t, _ := template.ParseFiles("templates/index.html")
  f := &Flash{Message: "Object saved."}
  decoder := json.NewDecoder(r.Body)
  // var obj Obj
  var obj map[string]interface{}
  err := decoder.Decode(&obj)
  if err != nil {
    fmt.Println("Error saving object.")
    log.Fatal(err)
  }
  t.Execute(w, f)
}

func main() {
	RiakTransaction(func(client *riak.Client) {
  	bucket, _ := client.Bucket("riaktor")

    // insertTestObject(bucket)
  	obj := bucket.New("testobj")
  	obj.ContentType = "application/json"
  	obj.Data = []byte("{'name': 'Bob'}")
  	obj.Store()

    var soda *Soda
    soda = &Soda{SugarFree: true, Flavor: "cola", Brand: "Coca-Cola"} 
    err := client.New("riaktor", "test_soda", soda)
    err = soda.SaveAs("test_soda") // dev.Save()
    if err != nil {
      fmt.Println("Error saving object.")
    }

    // var soda2 Soda
    // err = client.Load("sodas", "abcdefghijklm", &soda2)
    // soda2.Brand = "Pepsi"
    // err = soda2.saveAs("newsoda")
    // if err != nil {
    //   fmt.Println("Error saving object.")
    // }
    
		fmt.Printf("Stored objects in Riak, vclock = %v\n", obj.Vclock)
	})

  http.HandleFunc("/", handler)
  http.HandleFunc("/new", newHandler)
  http.ListenAndServe(fmt.Sprintf(":%s",httpPort), nil)
  
	fmt.Println("Done.")
}
