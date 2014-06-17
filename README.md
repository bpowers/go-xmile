go-xmile - Idiomatic Go for reading and writing XMILE files
===========================================================

This is a project to experiment with and track the work the OASIS
XMILE Technical Committee is doing to standardize the XMILE file
format.

example
-------

The following simple Go program:

```go
package main

import (
	"encoding/xml"
	"github.com/bpowers/go-xmile/xmile"
	"log"
	"os"
)

func xname(name string) xml.Name {
	return xml.Name{"", name}
}

func main() {
	m := &xmile.Model{
		Variables: []*xmile.Variable{
			&xmile.Variable{
				XMLName:  xname("aux"),
				Name:     "net_inflows",
				Equation: "births + migrations",
				Units:    "people/yard",
			},
			&xmile.Variable{
				XMLName:  xname("flow"),
				Name:     "births",
				Equation: "population*.08",
				Units:    "people/year",
			},
			&xmile.Variable{
				XMLName:  xname("flow"),
				Name:     "deaths",
				Equation: "population*.07",
				Units:    "people/year",
			},
			&xmile.Variable{
				XMLName:  xname("flow"),
				Name:     "migrations",
				Equation: "10",
				Units:    "people/year",
			},
			&xmile.Variable{
				XMLName:  xname("stock"),
				Name:     "population",
				Equation: "100",
				Inflows:  []string{"births", "migrations"},
				Outflows: []string{"deaths"},
				Units:    "people",
			},
		},
		Views: []*xmile.View{},
	}

	f := xmile.NewFile(3, "hello xworld")
	f.Models = append(f.Models, m)
	f.SimSpec.TimeUnits = "year"

	output, err := xml.MarshalIndent(f, "", "    ")
	if err != nil {
		log.Fatalf("xml.MarshalIndent: %s", err)
	}

	os.Stdout.Write([]byte(xmile.XMLDeclaration + "\n"))
	os.Stdout.Write(output)
	os.Stdout.Write([]byte("\n"))
}
```

Produces the output:

```xml
<?xml version="1.0" encoding="utf-8" ?>
<xmile version="1.0" level="3">
    <header>
        <name>hello, xworld</name>
        <uuid>5ec6ca0d-74e5-a62c-e84a-527fb9753db1</uuid>
        <vendor>SDLabs</vendor>
        <product version="0.1" lang="en">go-xmile</product>
    </header>
    <sim_specs time_units="year">
        <start>0</start>
        <stop>0</stop>
        <dt>0</dt>
    </sim_specs>
    <model>
        <variables>
            <aux name="net_inflows">
                <eqn>births+migrations</eqn>
                <units>people/yard</units>
            </aux>
            <flow name="births">
                <eqn>population*.08</eqn>
                <units>people/year</units>
            </flow>
            <flow name="deaths">
                <eqn>population*.07</eqn>
                <units>people/year</units>
            </flow>
            <flow name="migrations">
                <eqn>10</eqn>
                <units>people/year</units>
            </flow>
            <stock name="population">
                <eqn>100</eqn>
                <inflow>births</inflow>
                <inflow>migrations</inflow>
                <outflow>deaths</outflow>
                <units>people</units>
            </stock>
        </variables>
    </model>
</xmile>
```
