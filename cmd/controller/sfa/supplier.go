package sfa

import (
	"encoding/json"
	"github.com/unknownfeature/dcw/cmd/util"
	"math"
	"sort"
	"sync"
)

type Alphabet int

const (
	Decimals Alphabet = iota
	Hex
	Uuid
	Base36
	Base64
	Custom
)

var (
	alphabetCharacters = map[Alphabet][]rune{
		Decimals: []rune("0123456789"),
		Hex:      []rune("0123456789abcdef"),
		Base36:   []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"),
		Base64:   []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"),
	}
)

type Formatter int

const (
	Simple Formatter = iota
	Uuid4
)

var (
	formattersFunctions = map[Formatter]func([]rune) (string, error){
		Simple: ToStringFromRunes,
		Uuid4:  ToUuid4StringFromRunes,
	}
)

type Config struct {
	Alphabet     []rune    `json:"alphabet"`
	ResultLength int       `json:"resultLength"`
	Formatter    Formatter `json:"formatter"`
}

type State struct {
	Config           Config `json:"config"`
	CurrentPositions []int  `json:"currentPositions"`
	Done             bool   `json:"done"`
}

type Supplier struct {
	state     State
	stateLock *sync.RWMutex
}

// todo actually add state persistance
const StateFile = "/home/sfa_gen.json"

func ForCustom(resultLength int, alphabet []rune, formatter Formatter) (*Supplier, error) {

	if resultLength <= 0 {
		return nil, IncorrectResultLengthError
	}
	if alphabet == nil || len(alphabet) == 0 {
		return nil, IncorrectAlphabetLengthError
	}
	if int(formatter) >= len(formattersFunctions) {
		return nil, IncorrectFormatterError
	}
	stateAlphabet := append([]rune(nil), alphabet...)
	sort.Slice(stateAlphabet, func(i, j int) bool {
		return stateAlphabet[i] < stateAlphabet[j]
	})
	state := State{Config: Config{stateAlphabet, resultLength, formatter}}
	return StringFromAlphabetGeneratorFromState(state)

}

func ForStandard(alphabet Alphabet, resultLength int, formatter Formatter) (*Supplier, error) {

	if alphabet == Custom {
		return nil, CustomNotSupportedError
	}
	return ForCustom(resultLength, alphabetCharacters[alphabet], formatter)
}

func Resume(stateFileLocation string) (*Supplier, error) {

	res, err := util.ReadToStruct[State](stateFileLocation, func() State { return State{} })
	if err != nil {
		return nil, err
	}
	return StringFromAlphabetGeneratorFromState(res)
}

func StringFromAlphabetGeneratorFromState(state State) (*Supplier, error) {

	return &Supplier{state, &sync.RWMutex{}}, nil
}

func (g *Supplier) Apply(batchSize int) ([]string, error) {

	currentPositions, err := g.recalculatePositions(batchSize)

	if err != nil {
		return nil, err
	}

	template := make([]rune, g.state.Config.ResultLength)
	chunk := make([]string, batchSize)
	err = g.generateBatch(&chunk, template, batchSize, 0, currentPositions)

	return chunk, err
}

func (g *Supplier) CurrentState() ([]byte, error) {
	g.stateLock.RLock()
	res, e := json.Marshal(g.state)
	g.stateLock.RUnlock()
	return res, e
}

func (g *Supplier) generateBatch(res *[]string, current []rune, batchSize int, depth int, currentIndices []int) error {

	if len(*res) == batchSize {
		return nil
	}
	alphabetLength := len(g.state.Config.Alphabet)

	if depth == alphabetLength {
		strRes, err := formattersFunctions[g.state.Config.Formatter](current)
		if err != nil {
			return err
		}
		*res = append(*res, strRes)
	}

	for i := currentIndices[depth]; i < alphabetLength; i++ {
		current[depth] = g.state.Config.Alphabet[i]
		currentIndices[depth] = i
		err := g.generateBatch(res, current, batchSize, depth+1, currentIndices)
		if err != nil {
			return err
		}
	}
	currentIndices[depth] = 0
	return nil
}

func (g *Supplier) updatePositions(positions *[]int, log int, sum int, index int) int {
	positionsDeref := *positions
	vocabLength := len(g.state.Config.Alphabet)

	if index == len(positionsDeref) {
		return 0
	}

	newLog := log
	adjustIndex := len(positionsDeref)-index == log
	newSum := sum
	newCarryover := 0
	if adjustIndex {
		iteration := int(math.Pow(float64(vocabLength), float64(log)))
		newSum = sum % iteration
		newCarryover = sum / iteration
		newLog = newLog - 1
	}

	carryover := g.updatePositions(positions, newLog, newSum, index+1)
	newValue := positionsDeref[index] + carryover
	if index == len(positionsDeref)-1 {
		newValue += newSum
	}
	positionsDeref[index] = int(math.Min(float64(newValue), float64(vocabLength)))

	return newCarryover
}

func (g *Supplier) recalculatePositions(batchSize int) ([]int, error) {

	g.stateLock.Lock()

	if g.state.Done {
		return nil, PotentialResultsExhaustedError
	}
	alphabetLength := len(g.state.Config.Alphabet)
	log := int(math.Log10(float64(batchSize)) / math.Log10(float64(alphabetLength)))

	oldPositions := append(make([]int, g.state.Config.ResultLength), g.state.CurrentPositions...)
	newPositions := append(make([]int, g.state.Config.ResultLength), g.state.CurrentPositions...)

	carryover := g.updatePositions(&newPositions, int(math.Min(float64(log), float64(g.state.Config.ResultLength))), batchSize, 0)

	if carryover > 0 {
		g.state.Done = true
	}

	g.state.CurrentPositions = newPositions

	g.stateLock.Unlock()

	return oldPositions, nil
}
