# chaikin

Go package implementing Chaikin's algorithm for line-smooting.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/aaronland/go-chaikin.svg)](https://pkg.go.dev/github.com/aaronland/go-chaikin)

## Example

```
package main

import (
	"github.com/aaronland/go-chaikin"
)

func TestSmooth(t *testing.T) {

	input := [][2]float64{
		[2]float64{185.6, 115.34581311988063},
		[2]float64{371.2, 426.3801223498234},
		[2]float64{556.8, 396.20750147765983},
		[2]float64{742.4, 447.7259263695299},
		[2]float64{185.6, 115.34581311988063},		
	}

	iterations := 6
	close := true
	
	output := chaikin.Smooth(input, iterations, close)
}
```

And then if you drew an image plotting `input` in red and `output` in black you'd see this:

![](docs/images/smooth.png)

## See also

* http://graphics.cs.ucdavis.edu/education/CAGDNotes/Chaikins-Algorithm/Chaikins-Algorithm.html
* https://observablehq.com/@pamacha/chaikins-algorithm