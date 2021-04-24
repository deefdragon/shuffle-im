package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"gonum.org/v1/gonum/mat"
)

type GenreateChances interface {
	//getChances returns the chance that card r ends up in place c
	GetChances(cards, iterations int) mat.Dense
}

type OrderGenerator interface {
	GetOrder(cards int) []int
}

type BaseGenerator struct {
	Orderer OrderGenerator
}

type BaseOrderer struct {
	Rand *rand.Rand
	//the values that could be gotten from grouping.
	GroupingChance []int
}

type CryptoSource struct {
}

func (cs *CryptoSource) Int63() int64 {
	b := make([]byte, 8)
	_, err := crand.Read(b)
	if err != nil {
		panic(err)
	}
	dat := binary.BigEndian.Uint64(b)
	return int64(dat)

}
func (cs *CryptoSource) Seed(seed int64) {

}

var OneOrTwo = []int{
	1, 2,
}
var Distributed = []int{
	0, 0,
	1, 1, 1, 1, 1, 1, 1, 1,
	2, 2, 2, 2,
	3, 3,
	4,
}
var ZeroOne = []int{
	0, 1,
}
var OneFive = []int{1, 5}

func main() {

	cards := 52
	cardIterations := 10000
	shuffles := 9

	//create output writer
	writer, err := os.OpenFile("data.csv", os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	//create output system
	bo := &BaseOrderer{
		GroupingChance: OneFive,
		Rand:           rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	g := BaseGenerator{
		Orderer: bo,
	}

	//create shuffle matrix
	//m^1
	multiplier := g.GetChances(cards, cardIterations)

	//generate initial output
	writer.WriteString("0\n")
	current := GenerateStartMatrix(cards)
	PrintMatrix(writer, current)

	//generate all subsequent outputs.
	for i := 0; i < shuffles; i++ {
		writer.WriteString(strconv.Itoa(i+1) + "\n")
		newMat := mat.NewDense(cards, cards, nil)
		newMat.Mul(current, multiplier)
		PrintMatrix(writer, newMat)
		current = newMat

	}

	writer.Close()

}

func GenerateStartMatrix(cards int) *mat.Dense {
	b := mat.NewDense(cards, cards, nil)
	for i := 0; i < cards; i++ {
		b.Set(i, i, 1)
	}
	return b
}

func PrintMatrix(writer io.Writer, matrix *mat.Dense) {
	rows, cols := matrix.Dims()
	b := strings.Builder{}
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			b.WriteString(fmt.Sprintf("%.9f,", matrix.At(i, j)))
		}
		b.WriteString("\n")
	}
	b.WriteString("\n\n")
	writer.Write([]byte(b.String()))

}

//returns the % chance that a deck will have a card end up in that position.
func (g *BaseGenerator) GetChances(cards, iterations int) *mat.Dense {
	b := mat.NewDense(cards, cards, nil)

	for i := 0; i < iterations; i++ {
		res := g.Orderer.GetOrder(cards)
		for i, v := range res {

			b.Set(v, i, b.At(v, i)+1.0)
		}
	}

	b.Apply(func(i, j int, v float64) float64 {
		return v / float64(iterations)
	}, b)

	return b
}

//returns a "shuffeled" deck
func (o *BaseOrderer) GetOrder(cards int) []int {
	dat := GenerateIncArr(cards)
	divider := cards / 2

	d1 := dat[:divider]
	d2 := dat[divider:]
	if o.Rand.Intn(2) == 0 {
		return o.Combine(d1, d2)
	}

	return o.Combine(d2, d1)
	//return append(d2, d1...)
}

func (o *BaseOrderer) Combine(d1, d2 []int) []int {
	d1idx := 0
	d2idx := 0
	d1rem := len(d1)
	d2rem := len(d2)
	res := make([]int, 0)
	for d1rem > 0 || d2rem > 0 {

		if d1rem > 0 {
			takeIdx := o.Rand.Intn(len(o.GroupingChance))
			take := o.GroupingChance[takeIdx]
			if take > 0 {
				if take >= d1rem {
					res = append(res, d1[d1idx:]...)
				} else {
					res = append(res, d1[d1idx:(d1idx+take)]...)

				}
				d1idx += take
				d1rem -= take
			}
		}
		if d2rem > 0 {
			takeIdx := o.Rand.Intn(len(o.GroupingChance))
			take := o.GroupingChance[takeIdx]
			if take > 0 {
				if take >= d2rem {
					res = append(res, d2[d2idx:]...)
				} else {
					res = append(res, d2[d2idx:(d2idx+take)]...)

				}
				d2idx += take
				d2rem -= take
			}
		}

	}

	return res
}

//GenerateIncArr returns an array counting from 0 to count-1
func GenerateIncArr(count int) []int {
	c := make([]int, count)
	for i := 0; i < count; i++ {
		c[i] = i
	}
	return c
}
