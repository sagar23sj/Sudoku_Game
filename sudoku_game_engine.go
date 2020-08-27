package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Sudoku struct {
	sudokuGrid      [9][9]int //Array for Storing Values in cells
	replicatedGrid  [9][9]int //Replica of Original Grid Before removing values
	answerGrid      [9][9]int
	difficultyLevel map[string]int
	//gridSize int = 9                                  //Size of Sudoku Puzzle
	//blockSize int = int(math.Sqrt(float64(gridSize)) //Size of a block
	gridSize  int
	blockSize int
	gameLevel string
}

func (s *Sudoku) initializeGame(puzzleSize int, subBoxSize int, level string) {
	s.gridSize = puzzleSize
	s.blockSize = subBoxSize
	s.gameLevel = level
}

// Constant Values for Setting Levels
func (s *Sudoku) setKValue() {
	s.difficultyLevel = make(map[string]int)
	var key string
	for i := 0; i < 3; i++ {
		key = strconv.Itoa(i)
		s.difficultyLevel[key] = int(s.gridSize * 2 * (i + 1))
	}
}

// Function For Generating Random Values
func randomValueGenerator(upperLimit int) (RandomInt int) {
	var RangeMin int = 1
	var RangeMax int = upperLimit + 1
	rand.Seed(time.Now().UnixNano())
	RandomInt = rand.Intn(RangeMax-RangeMin) + RangeMin
	return
}

// Function for Generating Puzzle
func (s *Sudoku) createPuzzle(gameLevel string) {

	s.fillDiagonalBoxes()                                  //Fill Diagonal Blocks
	s.fillRemainingCells(0, s.blockSize-1)                 //Fill Reamining Cells
	s.replicatedGrid = replicateOriginalGrid(s.sudokuGrid) //Copy Unaltered Grid in New Grid
	s.setKValue()                                          //Set K values in map
	s.removeKCells(gameLevel)                              //Remove Cells form Grid Based on Difficulty of Game
}

// Function for Replicating Winning Puzzle i.e Original Grid Before Removing Cell Values
func replicateOriginalGrid(sudokuGrid [9][9]int) (temporaryGrid [9][9]int) {
	temporaryGrid = sudokuGrid
	return
}

// Function to fill Boxes in Diagonal (Based on fact that Diagonal Boxes are Independent)
func (s *Sudoku) fillDiagonalBoxes() {
	for i := 0; i < s.gridSize; i = i + s.blockSize {
		s.fillIndividualBox(i, i)
	}
}

// Function for Filling Individual Box
func (s *Sudoku) fillIndividualBox(row int, col int) {
	var num int
	for i := 0; i < s.blockSize; i++ {
		for j := 0; j < s.blockSize; j++ {
			for {
				num = randomValueGenerator(s.gridSize)
				if s.uniqueBoxValidation(row, col, num) {
					break
				}
			}
			s.sudokuGrid[row+i][col+j] = num
		}
	}
}

// Function for Filling Remaining Cells in the Puzzle i.e Other than Diagonal Boxes
func (s *Sudoku) fillRemainingCells(i int, j int) bool {

	if j >= s.gridSize && i < (s.gridSize-1) {
		i = i + 1
		j = 0
	}
	if i >= s.gridSize && j >= s.gridSize {
		return true
	}
	if i < s.blockSize {
		if j < s.blockSize {
			j = s.blockSize
		}
	} else if i < (s.gridSize - s.blockSize) {
		if j == int(i/s.blockSize)*s.blockSize {
			j = j + s.blockSize
		}
	} else {
		if j == (s.gridSize - s.blockSize) {
			i = i + 1
			j = 0
			if i >= s.gridSize {
				return true
			}
		}
	}

	for num := 1; num <= s.gridSize; num++ {
		if s.uniqueValidation(i, j, num) {
			s.sudokuGrid[i][j] = num
			if s.fillRemainingCells(i, j+1) {
				return true
			}
			s.sudokuGrid[i][j] = 0
		}
	}
	return false
}

//Function for removing K-values from grid based on difficulty Level
func (s *Sudoku) removeKCells(gameLevel string) {
	var count int = s.difficultyLevel[gameLevel]
	fmt.Println(gameLevel)
	for {
		var cellID = randomValueGenerator(s.gridSize*s.gridSize - 1)

		i := (cellID / s.gridSize)
		j := (cellID % s.gridSize)

		if s.sudokuGrid[i][j] != 0 {
			count = count - 1
			s.sudokuGrid[i][j] = 0
		}
		if count == 0 {
			break
		}
	}
}

// Function for Validating uniqueness of element in block
func (s *Sudoku) uniqueBoxValidation(rowStart int, colStart int, num int) bool {
	for i := 0; i < s.blockSize; i++ {
		for j := 0; j < s.blockSize; j++ {
			if s.sudokuGrid[rowStart+i][colStart+j] == num {
				return false
			}
		}
	}
	return true
}

// Functiom for Validating Uniqueness of Element in Row
func (s *Sudoku) uniqueRowValidation(i int, num int) bool {
	for j := 0; j < s.gridSize; j++ {
		if s.sudokuGrid[i][j] == num {
			return false
		}
	}
	return true
}

// Functiom for Validating Uniqueness of Element in Column
func (s *Sudoku) uniqueColValidation(j int, num int) bool {
	for i := 0; i < s.gridSize; i++ {
		if s.sudokuGrid[i][j] == num {
			return false
		}
	}
	return true
}

// Functiom for Validating Uniqueness of Element in Row, Column and Box
func (s *Sudoku) uniqueValidation(i int, j int, num int) bool {
	status := (s.uniqueRowValidation(i, num) && s.uniqueColValidation(j, num) && s.uniqueBoxValidation(i-(i%s.blockSize), j-(j%s.blockSize), num))
	return status
}

// Function for Printing Puzzle Grid on Terminal
func (s *Sudoku) printSudoku(sudokuGrid [9][9]int) {
	for i := 0; i < s.gridSize; i++ {
		for j := 0; j < s.gridSize; j++ {
			if s.sudokuGrid[i][j] > 9 {
				fmt.Print(s.sudokuGrid[i][j], "|")
			} else {
				fmt.Print(s.sudokuGrid[i][j], " |")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
