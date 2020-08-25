package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

var sudokuGrid [9][9]int                              //Array for Storing Values in cells
var replicatedGrid [9][9]int                          //Replica of Original Grid Before removing values
var gridSize int = 9                                  //Size of Sudoku Puzzle
var blockSize int = int(math.Sqrt(float64(gridSize))) //Size of a block

// Constant Values for Setting Levels
const (
	low = iota
	medium
	high
)

// Function For Generating Random Values
func randomValueGenerator(upperLimit int) (RandomInt int) {
	var RangeMin int = 1
	var RangeMax int = upperLimit + 1
	rand.Seed(time.Now().UnixNano())
	RandomInt = rand.Intn(RangeMax-RangeMin) + RangeMin
	return
}

// Function for Generating Puzzle
func createPuzzle() {

	fillDiagonalBoxes()                                //Fill Diagonal Blocks
	fillRemainingCells(0, blockSize-1)                 //Fill Reamining Cells
	replicatedGrid = replicateOriginalGrid(sudokuGrid) //Copy Unaltered Grid in New Grid
	removeKCells(high)                                 //Remove Cells form Grid Based on Difficulty of Game

}

// Function for Replicating Winning Puzzle i.e Original Grid Before Removing Cell Values
func replicateOriginalGrid(sudokuGrid [9][9]int) (temporaryGrid [9][9]int) {
	temporaryGrid = sudokuGrid
	return
}

// Function to fill Boxes in Diagonal (Based on fact that Diagonal Boxes are Independent)
func fillDiagonalBoxes() {
	for i := 0; i < gridSize; i = i + blockSize {
		fillIndividualBox(i, i)
	}
}

// Function for Filling Individual Box
func fillIndividualBox(row int, col int) {
	var num int
	for i := 0; i < blockSize; i++ {
		for j := 0; j < blockSize; j++ {
			for {
				num = randomValueGenerator(gridSize)
				if uniqueBoxValidation(row, col, num) {
					break
				}
			}
			sudokuGrid[row+i][col+j] = num
		}
	}
}

// Function for Filling Remaining Cells in the Puzzle i.e Other than Diagonal Boxes
func fillRemainingCells(i int, j int) bool {

	if j >= gridSize && i < (gridSize-1) {
		i = i + 1
		j = 0
	}
	if i >= gridSize && j >= gridSize {
		return true
	}
	if i < blockSize {
		if j < blockSize {
			j = blockSize
		}
	} else if i < (gridSize - blockSize) {
		if j == int(i/blockSize)*blockSize {
			j = j + blockSize
		}
	} else {
		if j == (gridSize - blockSize) {
			i = i + 1
			j = 0
			if i >= gridSize {
				return true
			}
		}
	}

	for num := 1; num <= gridSize; num++ {
		if uniqueValidation(i, j, num) {
			sudokuGrid[i][j] = num
			if fillRemainingCells(i, j+1) {
				return true
			}
			sudokuGrid[i][j] = 0
		}
	}
	return false
}

//Function for removing K-values from grid based on difficulty Level
func removeKCells(difficultyLevel int) {
	var count int
	if difficultyLevel == low {
		count = 15
	} else if difficultyLevel == medium {
		count = 25
	} else {
		count = 40
	}

	for {
		var cellID = randomValueGenerator(gridSize*gridSize - 1)

		i := (cellID / gridSize)
		j := (cellID % gridSize)

		if sudokuGrid[i][j] != 0 {
			count = count - 1
			sudokuGrid[i][j] = 0
		}
		if count == 0 {
			break
		}
	}
}

// Function for Validating uniqueness of element in block
func uniqueBoxValidation(rowStart int, colStart int, num int) bool {
	for i := 0; i < blockSize; i++ {
		for j := 0; j < blockSize; j++ {
			if sudokuGrid[rowStart+i][colStart+j] == num {
				return false
			}
		}
	}
	return true
}

// Functiom for Validating Uniqueness of Element in Row
func uniqueRowValidation(i int, num int) bool {
	for j := 0; j < gridSize; j++ {
		if sudokuGrid[i][j] == num {
			return false
		}
	}
	return true
}

// Functiom for Validating Uniqueness of Element in Column
func uniqueColValidation(j int, num int) bool {
	for i := 0; i < gridSize; i++ {
		if sudokuGrid[i][j] == num {
			return false
		}
	}
	return true
}

// Functiom for Validating Uniqueness of Element in Row, Column and Box
func uniqueValidation(i int, j int, num int) bool {
	status := (uniqueRowValidation(i, num) && uniqueColValidation(j, num) && uniqueBoxValidation(i-(i%blockSize), j-(j%blockSize), num))
	return status
}

// Function for Printing Puzzle Grid on Terminal
func printSudoku(sudokuGrid [9][9]int) {
	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			if sudokuGrid[i][j] > 9 {
				fmt.Print(sudokuGrid[i][j], "|")
			} else {
				fmt.Print(sudokuGrid[i][j], " |")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {

	fmt.Println("Welcome To Sudoku Puzzle")
	createPuzzle()
	printSudoku(sudokuGrid)
	fmt.Println("----------------------------------")
	printSudoku(replicatedGrid)
	userDisplay()
}
