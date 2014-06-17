// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"encoding/xml"
	"fmt"
	"github.com/bpowers/go-xmile/compat"
	"github.com/bpowers/go-xmile/xmile"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

const formTmpl = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN"
          "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html>
    <head>
	<meta charset="utf-8"></meta>
        <title>convert to TC XMILE</title>

        <link href="https://fonts.googleapis.com/css?family=Droid+Sans|Droid+Sans+Mono" rel="stylesheet" type="text/css" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
    </head>

    <body>
        <p>choose a file to convert from isee v10 XMILE-draft to current XMILE TC format</p>
        <form action="/api/v1/convert/" enctype="multipart/form-data" method="post">
            <input type="file" name="data">
            <input type="submit" value="Convert">
        </form>
    </body>
</html>
`

type rootHandler struct{}

func (*rootHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")

	form := template.Must(template.New("").Parse(formTmpl))
	if err := form.Execute(rw, nil); err != nil {
		log.Printf("login tmpl.Execute: %v\n", err)
	}

}

type convertHandler struct{}

func (*convertHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	contents, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("err: %s", err)
		fmt.Fprintf(rw, "an unknown error occured. please try a different file.")
		return
	}

	var iseeFile *compat.File
	if iseeFile, err = compat.ReadFile(contents); err != nil {
		log.Printf("compat.ReadFile: %s", err)
		fmt.Fprintf(rw, "an unknown error occured. please try a different file.")
		return
	}
	var f xmile.Node
	if f, err = compat.ConvertFromIsee(iseeFile, false); err != nil {
		log.Printf("compat.ConvertFromIsee: %s", err)
		fmt.Fprintf(rw, "an unknown error occured. please try a different file.")
		return
	}
	var output []byte
	if output, err = xml.MarshalIndent(f, "", "    "); err != nil {
		log.Printf("xml.MarshalIndent: %s", err)
		fmt.Fprintf(rw, "an unknown error occured. please try a different file.")
		return
	}
	rw.Header().Set("Content-Type", "application/xmile; charset=utf-8")
	rw.Header().Set("Content-Description", "File Transfer")
	rw.Header().Set("Content-Disposition", `attachment; filename="TC_Converted.xmile"`)
	rw.Header().Set("Content-Transfer-Encoding", "binary")
	rw.Write([]byte(xmile.XMLDeclaration + "\n"))
	rw.Write(output)
	rw.Write([]byte("\n"))
}

func main() {
	var err error

	http.Handle("/", &decacheHandler{&rootHandler{}})
	http.Handle("/api/v1/convert/", &decacheHandler{&convertHandler{}})

	err = http.ListenAndServe(
		":8010",
		nil)

	if err != nil {
		log.Printf("ListenAndServe:", err)
	}
}
