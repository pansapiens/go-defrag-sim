package scandisk

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

var (
	Width     int
	Height    int
	BaseDelay int
	Loop      bool
)

const (
	Empty            = '░'
	Full             = '█' // '◘'
	Frag             = '▒'
	BadBlock         = 'B'
	WritingBlock     = 'W'
	ReadingBlock     = 'r'
	UnmovableBlock   = 'X'
	MinChunkSizeFrag = 3
	MaxChunkSizeFrag = 15
	MinChunkSizeFull = 100
	MaxChunkSizeFull = 400
)

type Block struct {
	X, Y int
	Type rune
}

func InitializeBlocks() []Block {
	blocks := make([]Block, Width*Height)
	for i := range blocks {
		x := i % Width
		y := i / Width
		r := rand.Intn(100)
		var blockType rune
		switch {
		case r < 70:
			blockType = Empty
		case r < 90:
			blockType = Full
		case r < 98:
			blockType = Frag
		case r < 99:
			blockType = BadBlock
		default:
			blockType = Empty
		}
		blocks[i] = Block{X: x, Y: y, Type: blockType}
	}

	// Add unmovable blocks
	for i := range blocks {
		if rand.Intn(500) < 1 {
			blocks[i].Type = UnmovableBlock
		}
	}

	return blocks
}

func DrawBlocks(screen tcell.Screen, blocks []Block, scanIndex int, writeIndex int) {
	for i, block := range blocks {
		style := tcell.StyleDefault
		if i <= scanIndex {
			if i == writeIndex {
				block.Type = WritingBlock
				style = style.Foreground(tcell.ColorYellow)
			} else if block.Type == BadBlock || block.Type == UnmovableBlock {
				style = style.Foreground(tcell.ColorYellow)
			} else {
				block.Type = Full
				style = style.Foreground(tcell.ColorYellow)
			}
		} else {
			switch block.Type {
			case Full:
				style = style.Foreground(tcell.ColorWhite)
			case Frag:
				style = style.Foreground(tcell.ColorRed)
			case BadBlock:
				style = style.Foreground(tcell.ColorRed)
			case UnmovableBlock:
				style = style.Foreground(tcell.ColorFuchsia)
			case WritingBlock:
				style = style.Foreground(tcell.ColorGreen)
			case ReadingBlock:
				style = style.Foreground(tcell.ColorTeal)
			case Empty:
				style = style.Foreground(tcell.ColorBlack)
			}
		}
		screen.SetContent(block.X, block.Y+1, block.Type, nil, style)
	}
}

func DrawStatusBar(screen tcell.Screen, text string) {
	for x := 0; x < Width; x++ {
		screen.SetContent(x, 0, ' ', nil, tcell.StyleDefault.Background(tcell.ColorBlue))
	}
	for i, r := range text {
		screen.SetContent(i, 0, r, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue))
	}
}

func DrawLegendBar(screen tcell.Screen) {
	legend := "Legend: ░=Empty, █=Full, ▒=Fragmented, B=Bad, W=Writing, r=Reading, X=Unmovable"
	for x := 0; x < Width; x++ {
		screen.SetContent(x, Height+1, ' ', nil, tcell.StyleDefault.Background(tcell.ColorBlue))
	}
	for i, r := range legend {
		screen.SetContent(i, Height+1, r, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue))
	}
}

func FindNextChunk(blocks []Block, blockType rune, minChunkSize int, maxChunkSize int) []int {
	chunkSize := rand.Intn(maxChunkSize-minChunkSize+1) + minChunkSize
	fragmentedIndices := []int{}

	for i := range blocks {
		if blocks[i].Type == blockType {
			fragmentedIndices = append(fragmentedIndices, i)
			if len(fragmentedIndices) >= chunkSize {
				break
			}
		}
	}

	return fragmentedIndices
}

func FindFirstEmptyBlock(blocks []Block) int {
	for i := range blocks {
		if blocks[i].Type == Empty {
			return i
		}
	}
	return -1
}

func DefragBlocks(
	blocks []Block,
	targetBlockType rune,
	minChunkSize int,
	maxChunkSize int,
	screen tcell.Screen) bool {

	fragmentedIndices := FindNextChunk(blocks, targetBlockType, minChunkSize, maxChunkSize)

	if len(fragmentedIndices) == 0 {
		return false
	}

	emptyIndex := FindFirstEmptyBlock(blocks)

	if emptyIndex == -1 {
		return false
	}

	for _, idx := range fragmentedIndices {
		if blocks[idx].Type == targetBlockType {
			blocks[idx].Type = ReadingBlock
			screen.Clear()
			DrawStatusBar(screen, "Reading fragmented block...")
			DrawBlocks(screen, blocks, -1, -1)
			DrawLegendBar(screen)
			screen.Show()
			time.Sleep(time.Duration(BaseDelay) * time.Millisecond)

			blocks[idx].Type = Empty
			blocks[emptyIndex].Type = Full
			emptyIndex++
			screen.Clear()
			DrawStatusBar(screen, fmt.Sprintf("Moving %d fragmented blocks...", len(fragmentedIndices)))
			DrawBlocks(screen, blocks, -1, -1)
			DrawLegendBar(screen)
			screen.Show()
			time.Sleep(time.Duration(BaseDelay) * time.Millisecond)
		}
		if emptyIndex >= len(blocks) {
			break
		}
	}

	return true
}

func CompactFullBlocks(blocks []Block, screen tcell.Screen) {
	for i := 0; i < len(blocks); i++ {
		if blocks[i].Type == Full {
			emptyIndex := FindFirstEmptyBlock(blocks)
			if emptyIndex == -1 || emptyIndex >= i {
				break
			}
			blocks[emptyIndex].Type = WritingBlock
			screen.Clear()
			DrawStatusBar(screen, "Compacting full blocks...")
			DrawBlocks(screen, blocks, i, emptyIndex)
			DrawLegendBar(screen)
			screen.Show()
			time.Sleep(time.Duration(BaseDelay) * time.Millisecond)

			blocks[emptyIndex].Type = Full
			blocks[i].Type = Empty
		}
	}
}

func ScanBlocks(blocks []Block, screen tcell.Screen) {
	for i := 0; i < len(blocks); i++ {
		if blocks[i].Type == Frag && rand.Intn(10) < 2 {
			blocks[i].Type = ReadingBlock
			screen.Clear()
			DrawStatusBar(screen, fmt.Sprintf("Reading block %d/%d", i+1, len(blocks)))
			DrawBlocks(screen, blocks, i, i)
			DrawLegendBar(screen)
			screen.Show()
			time.Sleep(time.Duration(BaseDelay/2) * time.Millisecond)
			blocks[i].Type = Empty
		} else {
			DrawStatusBar(screen, fmt.Sprintf("Scanning block %d/%d", i+1, len(blocks)))
			DrawBlocks(screen, blocks, i, i)
			DrawLegendBar(screen)
			screen.Show()
			time.Sleep(time.Duration(BaseDelay/2) / 2 * time.Millisecond)
		}
	}
}

func RunDefrag(skipScan bool, screen tcell.Screen) {
	blocks := InitializeBlocks()
	DrawBlocks(screen, blocks, -1, -1)
	DrawLegendBar(screen)
	screen.Show()

	if !skipScan {
		ScanBlocks(blocks, screen)
	}

	for {
		if !DefragBlocks(blocks, Frag, MinChunkSizeFrag, MaxChunkSizeFrag, screen) {
			break
		}
		if !DefragBlocks(blocks, Full, MinChunkSizeFull, MaxChunkSizeFull, screen) {
			break
		}
		if !DefragBlocks(blocks, Full, MinChunkSizeFull, MaxChunkSizeFull, screen) {
			break
		}
		if !DefragBlocks(blocks, Full, MinChunkSizeFull, MaxChunkSizeFull, screen) {
			break
		}
		if !DefragBlocks(blocks, Full, MinChunkSizeFull, MaxChunkSizeFull, screen) {
			break
		}
	}

	CompactFullBlocks(blocks, screen)
}
