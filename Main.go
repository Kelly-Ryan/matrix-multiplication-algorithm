//Name: Kelly Ryan 
//Student No. 0347345
//Level: A1

package main

import (
	"bufio"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

//Global variables
var inputs = make([]int, 4)
var useNumCPUs int
var testMatrixDims = [][] int {{500,500,500,500}, {1000,1000,1000,1000}, {1500,1500,1500,1500}, {2000,2000,2000,2000}}
var performanceTimes = make([] float64, 4 * len(testMatrixDims) * runtime.NumCPU())

func inputMatrixDims() []int {
	fmt.Println("Enter the dimensions of 2 matrices to be multiplied.\nPress return after each number.")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("\nMatrix A - no. of rows:")
	scanner.Scan()
	inputs[0], _ = strconv.Atoi(scanner.Text())
	fmt.Print("Matrix A - no. of columns:")
	scanner.Scan()
	inputs[1], _ = strconv.Atoi(scanner.Text())
	fmt.Print("Matrix B - no. of rows:")
	scanner.Scan()
	inputs[2], _ = strconv.Atoi(scanner.Text())
	fmt.Print("Matrix B - no. of columns:")
	scanner.Scan()
	inputs[3], _ = strconv.Atoi(scanner.Text())
	fmt.Println()

	if scanner.Err() != nil {
		fmt.Println("Error: ", scanner.Err())
	}
	return inputs
}

func validate(inputs []int) bool {
	//check no. of columns in matrix A is equal to no. of rows in matrix B
	return inputs[1] == inputs[2]
}

func setNumCPU() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter number of available CPUs to use and press return:")
	scanner.Scan()
	input, _ := strconv.Atoi(scanner.Text())

	if input > 0 && input <= runtime.NumCPU() {
		useNumCPUs = input
	} else {
		setNumCPU()
	}
}

func matrixGeneration(inputs [] int) (a, b [][] int) {
	//assign validated row and column values
	aRows := inputs[0]
	aCols := inputs[1]
	bRows := inputs[2]
	bCols := inputs[3]

	//set seed to generate random matrix values on each function call
	rand.Seed(time.Now().UnixNano())

	//make 2D slices to hold values of input matrices
	aVals := make([][]int, aRows)
	for i := 0; i < aRows; i++ {
		aVals[i] = make([]int, aCols)
	}

	bVals := make([][]int, bRows)
	for i := 0; i < bRows; i++ {
		bVals[i] = make([]int, bCols)
	}

	//populate 2D slices with random numbers
	for i := 0; i < aRows; i++ {
		for j := 0; j < aCols; j++ {
			aVals[i][j] = rand.Intn(100)
		}
	}

	for i := 0; i < bRows; i++ {
		for j := 0; j < bCols; j++ {
			bVals[i][j] = rand.Intn(100)
		}
	}

	//display input matrices
	//fmt.Println("\n\nMatrix Generation:")
	//fmt.Println("\nMatrix A:")
	//for _, a := range aVals{
	//	fmt.Println(a)
	//}
	//fmt.Println("\nMatrix B:")
	//for _, b := range bVals{
	//	fmt.Println(b)
	//}

	return aVals, bVals
}

func singleThreadCalc(A, B[][] int) time.Duration{
	//sequential matrix multiplication algorithm for comparison purposes
	runtime.GOMAXPROCS(useNumCPUs)
	aRows, aCols, bCols := len(A), len(A[0]), len(B[0])
	sum := 0

	//create 2D slice C to hold results of A X B
	C := make([][] int, aRows)
	for i := 0; i < aRows; i++ {
		C[i] = make([]int, bCols)
	}

	//time multiplication operation
	start := time.Now()
	for i := 0 ; i < aRows; i++ {
		for k := 0; k < bCols; k++ {
			for j := 0;  j < aCols; j++ {
				sum = sum + (A[i][j] * B[j][k])
			}
			C[i][k] = sum
			sum = 0
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("1. 1 thread - sequential: %v\n", elapsed)
	return elapsed

	//display result matrix
	//fmt.Println("\nMethod 2\nResult Matrix C: ")
	//for _, c := range C {
	//	fmt.Println(c)
	//}
}

func twoThreadCalc(A, B[][] int) time.Duration {
	//this function creates two goroutines, each taking half of the calculation
	runtime.GOMAXPROCS(useNumCPUs)
	aRows, aCols, bCols := len(A), len(A[0]), len(B[0])
	sum := 0
	splitVal := len(A) / 2

	//create 2D slice C to hold results of A X B
	C := make([][]int, aRows)
	for i := 0; i < aRows; i++ {
		C[i] = make([]int, bCols)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	start := time.Now()
	go func() {
		defer wg.Done()
		for i := 0; i < splitVal; i++ {
			for k := 0; k < bCols; k++ {
				for j := 0; j < aCols; j++ {
					sum = sum + (A[i][j] * B[j][k])
				}
				C[i][k] = sum
				sum = 0
			}
		}
	}()

	go func() {
		defer wg.Done()
		for i := splitVal; i < aRows; i++ {
			for k := 0; k < bCols; k++ {
				for j := 0; j < aCols; j++ {
					sum = sum + (A[i][j] * B[j][k])
				}
				C[i][k] = sum
				sum = 0
			}
		}
	}()

	wg.Wait()
	elapsed := time.Since(start)

	fmt.Printf("2. 2 threads - 1 thread per n/2 rows: %v\n", elapsed)
	return elapsed

	//output result matrix
	//fmt.Println("\nLoop Split\nResult Matrix C: ")
	//for _, c := range C {
	//	fmt.Println(c)
	//}
}

func loopSplitCalc(A, B[][] int) time.Duration {
	runtime.GOMAXPROCS(useNumCPUs)
	aRows, aCols, bCols := len(A), len(A[0]), len(B[0])

	//create 2D slice C to hold results of A X B
	C := make([][] int, aRows)
	for i := 0; i < aRows; i++ {
		C[i] = make([]int, bCols)
	}

	noOfRoutines := (aRows / 100) + 1
	beg, increment := 0, aRows % 100
	var end int

	var wg sync.WaitGroup
	wg.Add(noOfRoutines)

	var calc = func(beg, end int){
		var sum int
		go func(){
			defer wg.Done()
			for i := beg ; i < end; i++ {
				for k := 0; k < bCols; k++ {
					for j := 0;  j < aCols; j++ {
						sum = sum + (A[i][j] * B[j][k])
					}
					C[i][k] = sum
					sum = 0
				}
			}
		}()
	}

	start := time.Now()
	for n := 0; n < noOfRoutines; n++ {
		end = (n * 100) + increment
		calc(beg, end)
		beg += end
		end += increment
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("3. %d thread(s) - 1 thread per 100 rows: %v\n", noOfRoutines, elapsed)

	//output results
	//fmt.Println("\nLoop Split Algorithm\nResult Matrix C: ")
	//for _, c := range C {
	//	fmt.Println(c)
	//}

	return elapsed
}

func threadPerRowCalc(A, B[][] int) time.Duration {
	//this function generates a new goroutine for each dot product (row x column) calculation
	runtime.GOMAXPROCS(useNumCPUs)
	aRows, aCols, bCols := len(A), len(A[0]), len(B[0])
	sum, threadCount := 0, 0

	//create 2D slice C to hold results of A X B
	C := make([][] int, aRows)
	for i := 0; i < aRows; i++ {
		C[i] = make([]int, bCols)
	}

	var wg sync.WaitGroup
	start := time.Now()
	for i := 0 ; i < aRows; i++ {
		for k := 0; k < bCols; k++ {
			wg.Add(1)
			go func(i, k int) {
				defer wg.Done()
				for j := 0;  j < aCols; j++ {
					sum = sum + (A[i][j] * B[j][k])
					threadCount++
				}
			}(i, k)
			C[i][k] = sum
			sum = 0
		}
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("4. %d thread(s) - 1 thread per row: %v\n", aRows, elapsed)
	return elapsed
}

func generateGraphs() {
	j := 0
	for i := 1; i <= runtime.NumCPU();  i++ {
		alg1 := []opts.LineData {
			{Value:performanceTimes[j]},
			{Value:performanceTimes[j + 4]},
			{Value:performanceTimes[j + 8]},
			{Value:performanceTimes[j + 12]},
		}

		alg2 := []opts.LineData{
			{Value: performanceTimes[j + 1]},
			{Value: performanceTimes[j + 5]},
			{Value: performanceTimes[j + 9]},
			{Value: performanceTimes[j + 13]},
		}

		alg3 := []opts.LineData{
			{Value: performanceTimes[j + 2]},
			{Value: performanceTimes[j + 6]},
			{Value: performanceTimes[j + 10]},
			{Value: performanceTimes[j + 14]},
		}

		alg4 := []opts.LineData {
			{Value:performanceTimes[j + 3]},
			{Value:performanceTimes[j + 7]},
			{Value:performanceTimes[j + 11]},
			{Value:performanceTimes[j + 15]},
		}

		line := charts.NewLine()
		line.SetGlobalOptions(
			charts.WithTitleOpts(opts.Title{Title: fmt.Sprintf("%d CPU(s)", i), Subtitle: "Y-axis Seconds, X-axis Matrix Dimensions n x n"}),
			charts.WithLegendOpts(opts.Legend{Show: true}))
		line.SetXAxis([]int{500, 1000, 1500, 2000}).
			AddSeries("1 thread", alg1).
			AddSeries("2 threads", alg2).
			AddSeries("n/100 threads", alg3).
			AddSeries("n threads", alg4)

		j += 4 * len(testMatrixDims)

		f, _ := os.Create(fmt.Sprintf("%dCPU.html", i))
		err := line.Render(f)
		if err != nil {
			return
		}
	}
}

func performanceTest() {
	fmt.Println("\n********Performance Test********")
	for i := 0; i < len(performanceTimes); {
		for j := 1; j <= runtime.NumCPU(); j++ {
			fmt.Printf("\n----------%d CPU(s)----------\n", j)
			for k := range testMatrixDims {
				fmt.Printf("\nInput Matrix Dimensions: %d x %d", testMatrixDims[k][k], testMatrixDims[k][k])
				fmt.Println("\n=====================================")
				performanceTimes[i] = singleThreadCalc(matrixGeneration(testMatrixDims[k])).Seconds()
				performanceTimes[i + 1] = twoThreadCalc(matrixGeneration(testMatrixDims[k])).Seconds()
				performanceTimes[i + 2] = loopSplitCalc(matrixGeneration(testMatrixDims[k])).Seconds()
				performanceTimes[i + 3] = threadPerRowCalc(matrixGeneration(testMatrixDims[k])).Seconds()

				i += 4
			}
		}
	}

	fmt.Print("\n\n\t====================================\n\t\t\tSpeedup Calculation\n\t====================================\n")
	for j := 0; j < 16; {
		for i := 500; i <= 2000; {
			switch j {
			case 0:
				fmt.Println("\n----------------------------------------\n\t\tSequential Execution\n----------------------------------------")
			case 4:
				fmt.Println("\n----------------------------------------\n\t\t\t2 threads\n----------------------------------------")
			case 8:
				fmt.Println("\n----------------------------------------\n\t\t\tn/100 threads\n----------------------------------------")
			case 12:
				fmt.Println("\n----------------------------------------\n\t\t\tn threads\n----------------------------------------")
			}

			fmt.Printf("\nMATRIX SIZE(N) = %d\n", i)
			fmt.Printf("1 CPU(s) = %vs\n", performanceTimes[j])
			fmt.Printf("Speedup with 2 CPU(s) = %v\n", performanceTimes[j] / performanceTimes[j + 16])
			fmt.Printf("Speedup with 3 CPU(s) = %v\n", performanceTimes[j] / performanceTimes[j + 32])
			fmt.Printf("Speedup with 4 CPU(s) = %v\n\n", performanceTimes[j] / performanceTimes[j + 48])
			i += 500
			j++
		}
	}
	generateGraphs()
}

func main() {
	fmt.Print("**********Matrix Multiplier**********\n\nNo. of CPUs available: ", runtime.NumCPU(), "\n")

	var getUserInput func()
	getUserInput = func() {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("\nEnter:\n1. To enter your own matrix dimensions\n2. To run a performance test\n3. To exit" +
			"\n\nPress return after you have made your selection:")
		scanner.Scan()
		switch selection := scanner.Text(); selection {
		case "1":
			validDimensions := false
			for !validDimensions {
				fmt.Println("\nNo. of columns in Matrix A must be equal to no. of rows in Matrix B.")
				if validate(inputMatrixDims()) {
					validDimensions = true
				}
			}
			setNumCPU()

			if validDimensions {
				A, B := matrixGeneration(inputs)
				singleThreadCalc(A, B)
				twoThreadCalc(A, B)
				loopSplitCalc(A, B)
				threadPerRowCalc(A, B)
			}
		case "2":
			performanceTest()
		case "0":
			os.Exit(0)
		default:
			fmt.Println("\nInvalid input. Please try again.")
			getUserInput()
		}
	}
	getUserInput()
}

