package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"
)

// Organism - represent candidate solution
type Organism struct {
	DNA     []byte
	Fitness float64
}

// simple character-maching function
func (d *Organism) calcFitness(target []byte) {
	score := 0
	for i := 0; i < len(d.DNA); i++ {
		if d.DNA[i] == target[i] {
			score++
		}
	}
	d.Fitness = float64(score) / float64(len(d.DNA))
	return
}

func crossover(d1 Organism, d2 Organism) Organism {
	child := Organism{
		DNA:     make([]byte, len(d1.DNA)),
		Fitness: 0,
	}
	mid := rand.Intn(len(d1.DNA))
	for i := 0; i < len(d1.DNA); i++ {
		if i > mid {
			child.DNA[i] = d1.DNA[i]
		} else {
			child.DNA[i] = d2.DNA[i]
		}

	}
	return child
}

func (d *Organism) mutate() {
	for i := 0; i < len(d.DNA); i++ {
		if rand.Float64() < MutationRate {
			d.DNA[i] = byte(rand.Intn(95) + 32)
		}
	}
}

// create random candidate solution with fitness due to the
// target solution
func createOrganism(target []byte) (organism Organism) {
	randomArray := make([]byte, len(target))

	for i := 0; i < len(target); i++ {
		//random ascii code (just letters)
		randomArray[i] = byte(rand.Intn(95) + 32)
	}

	organism = Organism{
		DNA:     randomArray,
		Fitness: 0,
	}

	organism.calcFitness(target)
	return
}

func createPopulation(target []byte) (population []Organism) {
	population = make([]Organism, PopSize)
	for i := 0; i < PopSize; i++ {
		population[i] = createOrganism(target)
	}
	return
}

// create a pool with copies of organisms according to its fitness
func createPool(population []Organism, target []byte, maxFitness float64) (pool []Organism) {
	pool = make([]Organism, 0)
	for i := 0; i < len(population); i++ {
		population[i].calcFitness(target)
		numberOfCopies := int((population[i].Fitness / maxFitness) * 100)
		for n := 0; n < numberOfCopies; n++ {
			pool = append(pool, population[i])
		}
	}
	return
}

func naturalSelection(pool []Organism, population []Organism, target []byte) []Organism {
	next := make([]Organism, len(population))

	for i := 0; i < len(population); i++ {
		r1, r2 := rand.Intn(len(pool)), rand.Intn(len(pool))
		a := pool[r1]
		b := pool[r2]

		child := crossover(a, b)
		child.mutate()
		child.calcFitness(target)

		next[i] = child
	}
	return next
}

func getBest(population []Organism) Organism {
	best := 0
	for i := 0; i < len(population); i++ {
		if population[best].Fitness < population[i].Fitness {
			best = i
		}
	}
	return population[best]
}

var PopSize int = 900
var MutationRate float64 = 0

func main() {
	start := time.Now()
	rand.Seed(time.Now().UTC().UnixNano())

	target := []byte("Be or not to be")
	population := createPopulation(target)

	found := false
	generation := 0
	for !found {
		generation++
		bestOrganism := getBest(population)
		fmt.Printf("\r generation: %d | %s | fitness: %2f", generation, string(bestOrganism.DNA), bestOrganism.Fitness)

		if bytes.Compare(bestOrganism.DNA, target) == 0 {
			found = true
		} else {
			maxFitness := bestOrganism.Fitness
			pool := createPool(population, target, maxFitness)
			population = naturalSelection(pool, population, target)
		}

	}
	elapsed := time.Since(start)
	fmt.Printf("\nTime taken: %s\n", elapsed)
}
