package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

type trackerType struct {
	tracker map[string]bool
}

func (tracker *trackerType) add(name string) {
	tracker.tracker[name] = true
}

func (tracker *trackerType) clone() trackerType {
	newTracker := make(map[string]bool)
	for k, v := range tracker.tracker {
		newTracker[k] = v
	}
	return trackerType{newTracker}
}

func (tracker *trackerType) reset() {
	tracker.tracker = map[string]bool{}
}

func (t trackerType) String() string {
	keys := make([]string, 0)
	for key, _ := range t.tracker {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return fmt.Sprintf("%v", strings.Join(keys, ", "))
}

type gate interface {
	calc() int
	resetConnections()
	set(value int)
	clone(tracker *trackerType, wires map[string]gate) gate
	getInput1Name() string
	getInput2Name() string
	getTracker() *trackerType
	op() string
}

type basicGate struct {
	outName    string
	input1Name string
	input2Name string
	wires      map[string]gate
	tracker    *trackerType
}

func (g basicGate) getInputs() (gate, gate) {
	return g.wires[g.input1Name], g.wires[g.input2Name]
}

func (g basicGate) getInput1Name() string {
	return g.input1Name
}

func (g basicGate) getInput2Name() string {
	return g.input2Name
}

func (g basicGate) getTracker() *trackerType {
	return g.tracker
}

type andGate struct {
	basicGate
	input1 gate
	input2 gate
}

func (g andGate) calc() int {
	if g.input1 == nil || g.input2 == nil {
		g.input1, g.input2 = g.getInputs()
	}
	g.tracker.add(g.outName)
	return g.input1.calc() & g.input2.calc()
}

func (g *andGate) set(value int) {
}

func (g andGate) clone(tracker *trackerType, wires map[string]gate) gate {
	out := &andGate{g.basicGate, nil, nil}
	out.wires = wires
	out.tracker = tracker
	return out
}

func (g andGate) resetConnections() {
	g.input1 = nil
	g.input2 = nil
}

func (g andGate) op() string {
	return "AND"
}

type orGate struct {
	basicGate
	input1 gate
	input2 gate
}

func (g orGate) calc() int {
	if g.input1 == nil || g.input2 == nil {
		g.input1, g.input2 = g.getInputs()
	}
	g.tracker.add(g.outName)
	return g.input1.calc() | g.input2.calc()
}

func (g *orGate) set(value int) {
}

func (g orGate) clone(tracker *trackerType, wires map[string]gate) gate {
	out := &orGate{g.basicGate, nil, nil}
	out.wires = wires
	out.tracker = tracker
	return out
}

func (g orGate) resetConnections() {
	g.input1 = nil
	g.input2 = nil
}

func (g orGate) op() string {
	return "OR"
}

type xorGate struct {
	basicGate
	input1 gate
	input2 gate
}

func (g xorGate) calc() int {
	if g.input1 == nil || g.input2 == nil {
		g.input1, g.input2 = g.getInputs()
	}
	g.tracker.add(g.outName)
	return g.input1.calc() ^ g.input2.calc()
}

func (g *xorGate) set(value int) {
}

func (g xorGate) clone(tracker *trackerType, wires map[string]gate) gate {
	out := &xorGate{g.basicGate, nil, nil}
	out.wires = wires
	out.tracker = tracker
	return out
}

func (g xorGate) resetConnections() {
	g.input1 = nil
	g.input2 = nil
}

func (g xorGate) op() string {
	return "XOR"
}

type constGate struct {
	basicGate
	value int
}

func (g constGate) calc() int {
	return g.value
}

func (g *constGate) set(value int) {
	g.value = value
}

func (g constGate) clone(tracker *trackerType, wires map[string]gate) gate {
	out := &constGate{g.basicGate, g.value}
	out.wires = wires
	out.tracker = tracker
	return out
}

func (g constGate) resetConnections() {
}

func (g constGate) op() string {
	return "CONST"
}

// var forwardConnections map[string][]string = make(map[string][]string)

func cloneWires(tracker *trackerType, originalWires map[string]gate) map[string]gate {
	clonedWires := make(map[string]gate)
	for key, value := range originalWires {
		clonedWires[key] = value.clone(tracker, clonedWires)
	}
	return clonedWires
}

func getRegisterName(char string, registerNumber int) string {
	return fmt.Sprintf("%s%02d", char, registerNumber)
}

func allOutRegisters(wires map[string]gate) trackerType {
	out := trackerType{map[string]bool{}}
	for registerName, _ := range wires {
		if !strings.HasPrefix(registerName, "x") && !strings.HasPrefix(registerName, "y") {
			out.add(registerName)
		}
	}
	return out
}

func getValuesUntilBit(wires map[string]gate, char string, lastBit int) int {
	bit := 0
	value := 0
	registerName := getRegisterName(char, bit)
	_, exists := wires[registerName]

	for exists && bit <= lastBit {
		// log.Println("Adding", registerName, wires[registerName].calc(), "to", value, "at bit", bit)
		value += wires[registerName].calc() << uint(bit)
		bit++
		registerName = getRegisterName(char, bit)
		_, exists = wires[registerName]
	}
	// log.Println("Returning", value)
	return value
}

func getValues(wires map[string]gate, char string) int {
	bit := 0
	value := 0
	registerName := getRegisterName(char, bit)
	_, exists := wires[registerName]

	for exists {
		value += wires[registerName].calc() << uint(bit)
		bit++
		registerName = getRegisterName(char, bit)
		_, exists = wires[registerName]
	}
	return value
}

func setValues(wires map[string]gate, char string, value int) {
	bit := 0
	registerName := getRegisterName(char, bit)
	_, exists := wires[registerName]

	for exists {
		wires[registerName].set(value & 1)
		value >>= 1
		bit++
		registerName = getRegisterName(char, bit)
		_, exists = wires[registerName]
	}
}

func setToSlice(m map[string]bool) []string {
	res := make([]string, 0, len(m))
	for v := range m {
		res = append(res, v)
	}
	return res
}

func union(set1 []string, set2 []string) []string {
	m := make(map[string]bool)
	for _, v := range set1 {
		m[v] = true
	}
	for _, v := range set2 {
		m[v] = true
	}
	return setToSlice(m)
}

func addSets(set1 []string, set2 []string) []string {
	out := make([]string, 0)
	for _, v := range set1 {
		out = append(out, v)
	}
	for _, v := range set2 {
		out = append(out, v)
	}
	return out
}

func difference(set1 trackerType, set2 trackerType) []string {
	out := make([]string, 0)
	for key, _ := range set1.tracker {
		if _, exists := set2.tracker[key]; !exists {
			out = append(out, key)
		}
	}
	return out
}

func cloneArray(original []string) []string {
	clone := make([]string, len(original))
	copy(clone, original)
	return clone
}

func testAddition(wires map[string]gate, atBit int) bool {
	// testRange := []int{0, 1}
	// for bit := 0; bit <= atBit; bit++ {
	// 	testRange = append(testRange, 1<<bit, 1<<bit-1)
	// }
	testRange := []int{0, 1, 1 << atBit, 1<<atBit - 1}
	for _, i := range testRange {
		for _, j := range testRange {
			setValues(wires, "x", i)
			setValues(wires, "y", j)
			if getValuesUntilBit(wires, "z", atBit+1) != i+j {
				//log.Println("Failed at", i, "+", j, "=", i+j, "result:", getValues(wires, "z"))
				return false
			}
		}
	}

	return true
}

func testAdditionOld(tracker *trackerType, wires map[string]gate, atBit int) (bool, []string) {
	prevTracker := tracker.clone()

	tracker.reset()
	for i := 0; i < 4; i++ {
		xvalue := (i % 2) << atBit
		yvalue := (i >> 1) << atBit
		setValues(wires, "x", xvalue)
		setValues(wires, "y", yvalue)
		// log.Println("Testing", xvalue, "+", yvalue, "=", xvalue+yvalue, "result:", getValuesUntilBit(wires, "z", atBit+1))
		if getValuesUntilBit(wires, "z", atBit+1) != xvalue+yvalue {
			// log.Println("Failed at bit", atBit, "x:", xvalue, "y:", yvalue, "z:", getValuesUntilBit(wires, "z", atBit+1))
			return false, difference(*tracker, prevTracker)
		}
	}
	return true, difference(*tracker, prevTracker)
}

func testFullAddition(wires map[string]gate, untilBit int) bool {
	for i := 0; i < 1<<untilBit; i++ {
		for j := 0; j < 1<<untilBit; j++ {
			setValues(wires, "x", i)
			setValues(wires, "y", j)
			if getValues(wires, "z") != i+j {
				log.Println("Failed at", i, "+", j, "=", i+j, "result:", getValues(wires, "z"))
				return false
			}
		}
	}
	return true
}

type stringPair struct {
	first  string
	second string
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func chooseToSwitch(involved []string, rest []string, alreadySwitched []string, goodGates map[string]bool, maxPairs int) [][]stringPair {
	// generate a list of pairs, so that:
	// - take one, two, ...etc or all from involved
	// - pair each up with one of the unused involved or rest
	// - at most maxPairs pairs to be generated in one set

	var result [][]stringPair
	if len(involved) < maxPairs {
		maxPairs = len(involved)
	}

	var generatePairs func(prefix []stringPair, involved []string, rest []string, alreadySwitched []string, maxPairs int)
	generatePairs = func(prefix []stringPair, involved []string, rest []string, alreadySwitched []string, maxPairs int) {
		// log.Println("generatePairs", prefix)
		if maxPairs == 0 || len(alreadySwitched) == 8 || len(involved) == 0 || len(rest) == 0 {
			if len(prefix) > 0 {
				result = append(result, append([]stringPair(nil), prefix...))
			}
			return
		}

		for i := 0; i < len(involved); i++ {
			if contains(alreadySwitched, involved[i]) || goodGates[involved[i]] {
				continue
			}
			for j := i + 1; j < len(involved); j++ {
				if contains(alreadySwitched, involved[j]) || goodGates[involved[j]] {
					continue
				}
				newPrefix := append(prefix, stringPair{involved[i], involved[j]})
				newInvolved := append(append([]string(nil), involved[:i]...), involved[i+1:]...)
				newInvolved = append(append([]string(nil), involved[:j]...), involved[j+1:]...)
				newAlreadySwitched := append(alreadySwitched, involved[i], involved[j])
				generatePairs(newPrefix, newInvolved, rest, newAlreadySwitched, maxPairs-1)
			}
			for j := 0; j < len(rest); j++ {
				if contains(alreadySwitched, rest[j]) || goodGates[rest[j]] {
					continue
				}

				newPrefix := append(prefix, stringPair{involved[i], rest[j]})
				newInvolved := append(append([]string(nil), involved[:i]...), involved[i+1:]...)
				newRest := append(append([]string(nil), rest[:j]...), rest[j+1:]...)
				newAlreadySwitched := append(alreadySwitched, involved[i], rest[j])
				generatePairs(newPrefix, newInvolved, newRest, newAlreadySwitched, maxPairs-1)
			}
		}
	}

	//maxPairs -= len(alreadySwitched) / 2
	for pairs := 1; pairs <= maxPairs && pairs <= 4-len(alreadySwitched)/2; pairs++ {
		generatePairs(nil, involved, rest, alreadySwitched, pairs)
	}

	return result
}

func printWireConnections(wires map[string]gate) {
	keys := make([]string, 0)
	for key, _ := range wires {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		log.Println(wires[key].getInput1Name(), wires[key].getInput2Name(), "->", key)
	}
}

func isCircular(wires map[string]gate, node string) (bool, map[string]bool) {
	visited := make(map[string]bool)
	var visit func(string) bool
	visit = func(name string) bool {
		if strings.HasPrefix(name, "x") || strings.HasPrefix(name, "y") {
			return false
		}
		if _, exists := visited[name]; exists {
			return true
		}
		visited[name] = true
		if wires[name].getInput1Name() != "" && visit(wires[name].getInput1Name()) {
			return true
		}
		if wires[name].getInput2Name() != "" && visit(wires[name].getInput2Name()) {
			return true
		}
		return false
	}
	return visit(node), visited
}

func switchGates(wires map[string]gate, first string, second string) (bool, map[string]bool) {
	for _, gate := range wires {
		gate.resetConnections()
	}

	wires[first], wires[second] = wires[second], wires[first]

	circ, visited := isCircular(wires, first)
	if circ {
		return circ, visited
	}
	circ, visited = isCircular(wires, second)
	if circ {
		return circ, visited
	}
	return false, visited
}

func main() {
	file, err := os.Open("input.txt")
	defer file.Close()
	if err != nil {
		log.Fatalln("can't open file", err)
	}

	wires := make(map[string]gate)
	tracker := trackerType{map[string]bool{}}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		re := regexp.MustCompile(`^([a-z][0-9][0-9]): ([01])$`)
		matches := re.FindStringSubmatch(line)
		if matches != nil {
			name := matches[1]
			value := int(matches[2][0] - '0')
			wires[name] = &constGate{basicGate{name, "", "", wires, &tracker}, value}
			continue
		}

		re = regexp.MustCompile(`^([a-z0-9]+) ([A-Z]+) ([a-z0-9]+) -> ([a-z0-9]+)$`)
		matches = re.FindStringSubmatch(line)
		if matches != nil && len(matches) == 5 {
			input1Name := matches[1]
			op := matches[2]
			input2Name := matches[3]
			outName := matches[4]
			// forwardConnections[input1Name] = append(forwardConnections[input1Name], outName)
			// forwardConnections[input2Name] = append(forwardConnections[input2Name], outName)

			switch op {
			case "AND":
				wires[outName] = &andGate{basicGate{outName, input1Name, input2Name, wires, &tracker}, wires[input1Name], wires[input2Name]}
			case "OR":
				wires[outName] = &orGate{basicGate{outName, input1Name, input2Name, wires, &tracker}, wires[input1Name], wires[input2Name]}
			case "XOR":
				wires[outName] = &xorGate{basicGate{outName, input1Name, input2Name, wires, &tracker}, wires[input1Name], wires[input2Name]}
			}
		}
	}

	//num := getValues(wires, "z")
	//log.Println("Part 1:", num)
	// testing cloning:
	clonedWires := cloneWires(&tracker, wires)
	log.Println("Part 1 cloned:", getValues(clonedWires, "z"))
	wires = clonedWires

	// manual tests:
	// printWireConnections(wires)
	// setValues(wires, "x", 255)
	// setValues(wires, "y", 256)
	// tracker = trackerType{[]string{}}
	// wires = cloneWires(&tracker, wires)
	// log.Println("x:", getValues(wires, "x"), "y:", getValues(wires, "y"), "z:", getValuesUntilBit(wires, "z", 9))

	// log.Println("Tracker:", tracker)

	log.Println("Part 2...")
	// tracker.reset()
	// result, involved := testAddition(&tracker, cloneWires(&tracker, wires), 8)
	// log.Println("testing addition", result, involved)
	// log.Println("tracker:", tracker)
	// test addition bit-by-bit, find out what connections are involved, try to switch them. Once we found a working adder, we save it
	// and move on to the next bit
	// there can be multiple solutions --> we have to save all and continue testing with them
	// We can switch 4 pairs only -- ie we must keep track of the switches we made

	const MAX_BIT = 44 // x44, y44, z45 is the max

	// idea: store gates that are involved in good bits - so we don't mess with them
	goodGates := trackerType{map[string]bool{}}

	/*
	tracker.reset()
	wires = cloneWires(&tracker, wires)
	for bit := 0; bit <= MAX_BIT; bit++ {
		prevTracker := tracker.clone()
		result := testAddition(wires, bit)
		affected := difference(tracker, prevTracker)
		if result {
			log.Println("Bit", bit, "is good")
			for _, gate := range affected {
				goodGates.add(gate)
			}
		}
	}
	
	log.Println("Found", len(goodGates.tracker), "good gates")
	*/

	type wireTestType struct {
		wires           map[string]gate
		tracker         *trackerType
		alreadySwitched []string
	}
	//switchGates(wires, "z09", "pvm")
	//switchGates(wires, "z09", "kgr")
	// switchGates(wires, "z09", "nnf") // !!!!! ---> why are we sure: after bit 8, no other solution fixes it. Also, verified manually! (see notes.txt)
	// switchGates(wires, "z20", "nhs") // looks good too !--> z09 nnf z20 nhs fails after bit 30 :(
	// switchGates(wires, "ddn", "kqh") // found manually...
	// if !testFullAddition(wires, 10) {
	// 	log.Fatalln("Test failed")
	// }
	tracker.reset()
	// wireTests := []wireTestType{{wires, &tracker, []string{"z09", "nnf", "z20", "nhs", "ddn", "kqh"}}}
	wireTests := []wireTestType{{wires, &tracker, []string{}}}
	for bit := 0; bit <= MAX_BIT; bit++ {
		log.Println("Testing bit", bit, "numWireTests:", len(wireTests))
		newWireTests := make([]wireTestType, 0)
		prevPrevTracker := trackerType{map[string]bool{}}
		prevTracker := trackerType{map[string]bool{}}
		for _, toTest := range wireTests {
			log.Println("--Testing", toTest.alreadySwitched)
			prevPrevTracker = prevTracker.clone()
			prevTracker = toTest.tracker.clone()
			result := testAddition(toTest.wires, bit)
			affected := difference(*toTest.tracker, prevPrevTracker)

			if !result {
				unused := difference(allOutRegisters(toTest.wires), *toTest.tracker)
				unusedTracker := trackerType{map[string]bool{}}
				for _, gate := range unused {
					unusedTracker.add(gate)
				}
				log.Println("    Addition at", bit, "failed. Gates involved:", affected, "numUnused:", len(unused), "unused minus good:", len(difference(unusedTracker, goodGates)))
				log.Println("alll used so far:", toTest.tracker)
				if bit == 30 {
					log.Println("Unused:", unused)
					log.Println("Good gates:", goodGates)
					log.Println("all used so far:", toTest.tracker)
				}

				// note here: theoretically we should check for 4 switches -- tooks really, really long though.
				// we're lucky that at each bit at maximum 1 pair was switched.
				// checking for 2 pairs also runs in reasonable time
				// toSwitch := chooseToSwitch(affected, unused, toTest.alreadySwitched, map[string]bool{}, 4)
				toSwitch := chooseToSwitch(affected, unused, toTest.alreadySwitched, goodGates.tracker, 1) // excluding good gates
				// toSwitch := chooseToSwitch(affected, difference(allOutRegisters(wires), *toTest.tracker), toTest.pairsToTest)
				log.Println("    Testing", len(toSwitch), "switches, already switched:", toTest.alreadySwitched)
				numSwitches := 0
				for _, pairList := range toSwitch {
					if numSwitches%100000 == 0 {
						log.Println("    Switching gates", numSwitches, "/", len(toSwitch))
					}
					numSwitches++
					// log.Println("  Switching gates", pairList)
					newTracker := toTest.tracker.clone()
					newWires := cloneWires(&newTracker, toTest.wires)
					newAlreadySwitched := cloneArray(toTest.alreadySwitched)
					isCircular := false
					for _, pair := range pairList {
						isThisCircular, visited := switchGates(newWires, pair.first, pair.second)
						if isThisCircular && pair.first == "z09" && pair.second == "z10" {
							log.Println("  !! Circular dependency detected", pair.first, pair.second, visited)
						}
						isCircular = isCircular || isThisCircular
						newAlreadySwitched = append(newAlreadySwitched, pair.first, pair.second)
					}
					if isCircular {
						// log.Println("  !! Circular dependency detected")
						continue
					}
					if result = testAddition(newWires, bit); result {
						log.Println("    !! Found a solution by switching", pairList)

						newWireTest := wireTestType{newWires, &newTracker, newAlreadySwitched}
						newWireTests = append(newWireTests, newWireTest)
					}
				}
			} else {
				newWireTests = append(newWireTests, toTest)
			}
		}
		if len(newWireTests) == 0 {
			log.Fatalln("  No more solutions found")
		} else {
			log.Println("  PASS")
			wireTests = newWireTests
		}
	}
	sort.Strings(wireTests[0].alreadySwitched)
	log.Println("Part 2:", strings.Join(wireTests[0].alreadySwitched, ","))

	// log.Println(chooseToSwitch([]string{"a", "b", "c"}, []string{"d", "e", "f", "g", "h"}, 2))
	// printWireConnections(wires)
	// switchGates(wires, "z00", "z01")
	// printWireConnections(wires)

}
