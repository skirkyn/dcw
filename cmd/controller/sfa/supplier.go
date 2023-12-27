package sfa

import (
	"encoding/json"
	"github.com/unknownfeature/dcw/cmd/common/config"
	"github.com/unknownfeature/dcw/cmd/util"
	"math"
	"sort"
	"sync"
)

var (
	alphabetCharacters = map[config.Alphabet][]rune{
		config.Decimals: []rune("0123456789"),
		config.Hex:      []rune("0123456789abcdef"),
		config.Base36:   []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"),
		config.Base64:   []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"),
	}
)

var (
	formattersFunctions = map[config.Formatter]func([]rune) (string, error){
		config.Simple: ToStringFromRunes,
		config.Uuid4:  ToUuid4StringFromRunes,
	}
)

type Config struct {
	Alphabet     []rune           `json:"alphabet"`
	ResultLength int              `json:"resultLength"`
	Formatter    config.Formatter `json:"formatter"`
}

type State struct {
	Config           Config `json:"config"`
	CurrentPositions []int  `json:"currentPositions"`
	Done             bool   `json:"done"`
}

type Supplier struct {
	state     *State
	stateLock *sync.RWMutex
}

// todo actually add state persistance
const StateFile = "/home/sfa_gen.json"

func ForCustom(resultLength int, alphabet []rune, formatter config.Formatter) (*Supplier, error) {

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
	state := State{Config: Config{stateAlphabet, resultLength, formatter}, CurrentPositions: make([]int, resultLength)}
	return StringFromAlphabetGeneratorFromState(state)

}

func ForStandard(alphabet config.Alphabet, resultLength int, formatter config.Formatter) (*Supplier, error) {

	if alphabet == config.Custom {
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

	return &Supplier{&state, &sync.RWMutex{}}, nil
}

func (g *Supplier) Apply(batchSize int) ([]string, error) {

	// todo enable this(temporary disabled due to bugs)

	//currentPositions, err := g.recalculatePositions(batchSize)

	//if err != nil {
	//	return nil, err
	//}
	//
	template := make([]rune, g.state.Config.ResultLength)
	chunk := make([]string, 0)
	g.stateLock.Lock()
	_, _, err := g.generateBatch(&chunk, template, batchSize, 0, g.state.CurrentPositions)
	g.stateLock.Unlock()

	return chunk, err
}

func (g *Supplier) CurrentState() ([]byte, error) {
	g.stateLock.RLock()
	res, e := json.Marshal(g.state)
	g.stateLock.RUnlock()
	return res, e
}

func (g *Supplier) generateBatch(res *[]string, current []rune, left int, depth int, currentIndices []int) (bool, int, error) {

	if left == 0 {
		return false, left, nil
	}
	alphabetLength := len(g.state.Config.Alphabet)

	if depth == len(current) {
		strRes, err := formattersFunctions[g.state.Config.Formatter](current)
		if err != nil {
			return false, left, err
		}
		*res = append(*res, strRes)

		return true, left - 1, nil
	}

	counter := left
	times := 0
	carryover := false

	for times < alphabetLength && counter > 0 {
		current[depth] = g.state.Config.Alphabet[(times+currentIndices[depth])%alphabetLength]
		newCarryover, newLeft, err := g.generateBatch(res, current, counter, depth+1, currentIndices)
		counter = newLeft
		carryover = carryover || newCarryover
		if err != nil {
			break
		}
		times++
	}
	if carryover {
		oldVal := currentIndices[depth]
		newVal := oldVal + times
		adjustedVal := newVal % alphabetLength
		currentIndices[depth] = adjustedVal
		return newVal != adjustedVal, counter, nil
	}
	return false, counter, nil

}

func (g *Supplier) updatePositions(positions []int, log int, sum int, index int) int {

	vocabLength := len(g.state.Config.Alphabet)

	if index == len(positions) {
		return 0
	}

	newLog := log
	adjustIndex := len(positions)-index == log
	newSum := sum
	newCarryover := 0
	if adjustIndex {
		iteration := int(math.Pow(float64(vocabLength), float64(log)))
		newSum = sum % iteration
		newCarryover = sum / iteration
		newLog = newLog - 1
	}

	carryover := g.updatePositions(positions, newLog, newSum, index+1)
	newValue := positions[index] + carryover
	if index == len(positions)-1 {
		newValue += newSum
	}
	positions[index] = int(math.Min(float64(newValue), float64(vocabLength-1)))
	if positions[index] < newValue && newCarryover == 0 {
		newCarryover++
	}
	return newCarryover
}

func (g *Supplier) recalculatePositions(batchSize int) ([]int, error) {

	g.stateLock.Lock()

	if g.state.Done {
		return nil, PotentialResultsExhaustedError
	}
	alphabetLength := len(g.state.Config.Alphabet)
	log := int(math.Log10(float64(batchSize)) / math.Log10(float64(alphabetLength)))

	oldPositions := make([]int, len(g.state.CurrentPositions))
	newPositions := make([]int, len(g.state.CurrentPositions))
	copy(oldPositions, g.state.CurrentPositions)
	copy(newPositions, g.state.CurrentPositions)
	carryover := g.updatePositions(newPositions, int(math.Min(float64(log), float64(g.state.Config.ResultLength))), batchSize, 0)

	if carryover > 0 {
		g.state.Done = true
	}
	for i := range g.state.CurrentPositions {
		g.state.CurrentPositions[i] = newPositions[i]
	}
	g.stateLock.Unlock()

	return oldPositions, nil
}
