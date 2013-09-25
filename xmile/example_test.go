// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xmile_test

import (
	"encoding/xml"
	"github.com/bpowers/go-xmile/xmile"
	"log"
	"os"
)

func ExampleNewFile() {
	m := &xmile.Model{
		Variables: []*xmile.Variable{
			&xmile.Variable{
				XMLName:  xml.Name{Local: "flow"},
				Name:     "migrations",
				Equation: "10",
				Units:    "people/year",
			},
			&xmile.Variable{
				XMLName:  xml.Name{Local: "stock"},
				Name:     "population",
				Equation: "100",
				Inflows:  []string{"births", "migrations"},
				Outflows: []string{"deaths"},
				Units:    "people",
			},
		},
		Views: []*xmile.View{},
	}

	f := xmile.NewFile(1, "hello xworld")
	f.Models = append(f.Models, m)
	f.SimSpec.TimeUnits = "year"

	output, err := xml.MarshalIndent(f, "", "    ")
	if err != nil {
		log.Fatalf("xml.MarshalIndent: %s", err)
	}

	os.Stdout.Write([]byte(xmile.XMLDeclaration + "\n"))
	os.Stdout.Write(output)
	os.Stdout.Write([]byte("\n"))

	// Output:
	//<?xml version="1.0" encoding="utf-8" ?>
	//<xmile version="1.0" level="1">
	//    <header>
	//        <name>hello, xworld</name>
	//        <uuid>7a435517-ce5d-c816-9ec5-b34e44ec4fee</uuid>
	//        <vendor>XMILE TC</vendor>
	//        <product version="0.1" lang="en">go-xmile</product>
	//    </header>
	//    <sim_specs time_units="year">
	//        <start>0</start>
	//        <stop>0</stop>
	//        <dt>0</dt>
	//    </sim_specs>
	//    <model>
	//        <variables>
	//            <flow name="deaths">
	//                <eqn>population*.07</eqn>
	//                <units>people/year</units>
	//            </flow>
	//            <flow name="migrations">
	//                <eqn>10</eqn>
	//                <units>people/year</units>
	//            </flow>
	//            <stock name="population">
	//                <eqn>100</eqn>
	//                <inflow>births</inflow>
	//                <inflow>migrations</inflow>
	//                <outflow>deaths</outflow>
	//                <units>people</units>
	//            </stock>
	//        </variables>
	//    </model>
	//</xmile>
}
