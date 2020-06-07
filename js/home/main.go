package main

import (
    "net/http"
    "text/template"
    "time"
    "github.com/computermouth/cookies/pkg/dynamic"
    "honnef.co/go/js/dom"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "bytes"
)

func main(){
	
	for {
	
		d := dom.GetWindow().Document()
		elements := d.GetElementsByClassName("dynamic")
		if len(elements) > 1 {
			fmt.Println("TODO: handleme")
		}
		
		res, _ := http.Get("/projects")
		respdata, _ := ioutil.ReadAll(res.Body)
		
		projects := []dynamic.Project{}
		json.Unmarshal(respdata, &projects)
		
		t := template.Must(template.New("homedynamic").Parse(dynamic.HomeBody))
		
		var tmplbytes bytes.Buffer
		err := t.Execute(&tmplbytes, projects)
		if err != nil {
			fmt.Println(err)
		}
		
		elements[0].SetInnerHTML(tmplbytes.String())
		
		time.Sleep(2 * time.Second)
	}
	
}
