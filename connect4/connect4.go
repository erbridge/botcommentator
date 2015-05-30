package connect4

import (
	"errors"
)

type Board [][]int

func NewBoard(rows, columns int) Board {
	b := make(Board, columns)

	for i := range b {
		b[i] = make([]int, 0, rows)
	}

	return b
}

func (b Board) Duplicate() Board {
	newBoard := make(Board, len(b))
	for i := range b {
		newBoard[i] = make([]int, len(b[i]), cap(b[i]))
		copy(newBoard[i], b[i])
	}

	return newBoard
}

func (b Board) Play(piece, colIdx int) (err error) {
	height := len(b[colIdx])
	if height == cap(b[colIdx]) {
		return errors.New("That column is full. No more pieces allowed.")
	}

	b[colIdx] = b[colIdx][0 : height+1]
	b[colIdx][height] = piece

	return
}

// CheckWin returns the value of the piece that has won, if there is a win.
// Otherwise, it returns 0.
func (b Board) CheckWin() int {
	// Vertical
	for _, col := range b {
		if len(col) < 4 {
			continue
		}

		count := 0
		piece := 0
		for _, p := range col {
			if piece == p {
				count++
			} else {
				count = 1
				piece = p
			}

			if count == 4 {
				return piece
			}
		}
	}

	// Horizontal
	for rowIdx := 0; rowIdx < cap(b[0]); rowIdx++ {
		count := 0
		piece := 0
		for _, col := range b {
			if rowIdx >= len(col) {
				count = 0
				continue
			}

			if p := col[rowIdx]; piece == p {
				count++
			} else {
				count = 1
				piece = p
			}

			if count == 4 {
				return piece
			}
		}
	}

	// Up-Right Diagonal
	for rowIdx := 0; rowIdx < cap(b[0])-3; rowIdx++ {
		for colIdx := 0; colIdx < cap(b)-3; colIdx++ {
			count := 0
			piece := 0
			for offset := 0; offset < 4; offset++ {
				y := rowIdx + offset
				x := colIdx + offset

				col := b[x]

				if y >= len(col) {
					break
				}

				if p := col[y]; piece == p {
					count++
				} else {
					count = 1
					piece = p
				}

				if count == 4 {
					return piece
				}
			}
		}
	}

	// Up-Left Diagonal
	for rowIdx := 0; rowIdx < cap(b[0])-3; rowIdx++ {
		for colIdx := 3; colIdx < cap(b); colIdx++ {
			count := 0
			piece := 0
			for offset := 0; offset < 4; offset++ {
				y := rowIdx + offset
				x := colIdx - offset

				col := b[x]

				if y >= len(col) {
					break
				}

				if p := col[y]; piece == p {
					count++
				} else {
					count = 1
					piece = p
				}

				if count == 4 {
					return piece
				}
			}
		}
	}

	return 0
}

func (b Board) CountWins(piece, iterations int) (weightedCount, count int) {
	for colIdx := range b {
		newBoard := b.Duplicate()

		newBoard.Play(piece, colIdx)

		if win := newBoard.CheckWin(); win != 0 {
			weightedCount += win
			count++
		} else if iterations > 1 {
			w, t := newBoard.CountWins(-piece, iterations-1)
			weightedCount += w
			count += t
		}
	}

	return
}
