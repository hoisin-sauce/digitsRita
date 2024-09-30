package main

import "fmt"

// sum a set of 10 numbers associated digit pairs
// if the digit pairs are in ascending order
func evaluate(numbers [10]uint16)uint16{
	var (
		previous [5]uint16
		current uint16
		sum uint16
	)

	// go through all 5 numbers
	for i := 0; i < 5; i++{
		// create 2 digit number
		current = numbers[i*2] * 10 + numbers[i*2 + 1]

		// check if number was seen
		for j := 0; j < i; j++{
			if previous[j] > current{
				return 0 // if number was seen exit
			}
		}

		sum += current

		previous[i] = current
	}

	// return sum if no issues
	return sum
}


// controls worker goroutine
func manageWorker(inputChanel chan [10]uint16, outputChanel chan uint64){
	var (
		sum uint16
	)

	// goes through chanel input until chanel close
	for nums := range inputChanel{
		sum = evaluate(nums)
		outputChanel <- uint64(sum)
	}

	// close output to indicate no sums remaining
	close(outputChanel)
}


// adds up all input in a chanel and gives it to an output chanel
func summation(inputChanel chan uint64, sumChannel chan uint64){
	var sum uint64

	for i := range inputChanel{
		sum += i
	}

	sumChannel <- sum
}

// creates all possible permutations of an input array of length 10
// sends the output into a specified chanel
func createData(current [10]uint16, index int, outputChanel chan [10]uint16){
	var swap uint16

	// if attempting to alter beyond the bounds, data is ready for output
	if index == 10{
		outputChanel <- current
		return
	}

	// swaps values from the current specified index
	for i := index; i < len(current); i++{
		// copy for memory fuckery (all hail the garbage collector)
		copy := current
		swap = copy[index]
		copy[index] = copy[i]
		copy[i] = swap

		createData(copy, index + 1, outputChanel)
	}

	// when the original call has completed close the channel to indicate to other goroutines
	// that there are no more numbers
	if index == 0{
		close(outputChanel)
	}else if index == 1{
		fmt.Println(current[0]*10, "%")
	}
}

func main(){
	//setup chanels
	inputChanel := make(chan [10]uint16)
	outputChanel := make(chan uint64)
	sumChannel := make(chan uint64)

	// base data
	baseData := [10]uint16{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	// start goroutines
	go summation(outputChanel, sumChannel)
	go manageWorker(inputChanel, outputChanel)
	go createData(baseData, 0 , inputChanel)

	// recieve output
	outputData := <- sumChannel
	fmt.Println(outputData)
}