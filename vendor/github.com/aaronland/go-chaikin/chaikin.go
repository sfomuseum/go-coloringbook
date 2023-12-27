package chaikin

// Smooth will apply Chaikin's line-smooting algorithm to 'input' for as many iterations as defined
// by 'iterations'. If the 'close' argument is true that the output of each iterator will be returned
// as a closed path.
func Smooth(input [][2]float64, iterations int, close bool) [][2]float64 {

	if iterations == 1 {
		return input
	}

	return Smooth(smooth(input, close), iterations-1, close)
}

func smooth(input [][2]float64, close bool) [][2]float64 {

	output := make([][2]float64, 0)

	count := len(input)

	for i := 0; i < count-1; i++ {

		p1 := input[i]
		p2 := input[i+1%count]

		x1 := 0.75*p1[0] + 0.25*p2[0]
		y1 := 0.75*p1[1] + 0.25*p2[1]

		x2 := 0.25*p1[0] + 0.75*p2[0]
		y2 := 0.25*p1[1] + 0.75*p2[1]

		output = append(output, [2]float64{x1, y1})
		output = append(output, [2]float64{x2, y2})
	}

	if close {
		output = append(output, output[0])
	}

	return output
}
