package londondi

import "testing"

const CASE_0 = 0x028d00011b8300011b820000001dffec8203
const CASE_1 = 0x8d00011b8300011b820000001dffec82
const RESULT_1 = 0x8d0001

func TestChecksum(t *testing.T) {

}

func TestSubstitutions(t *testing.T) {

	//test a slice with an invalid length
	result, err := MakeSubstitutions(CASE_0, DECODE)
	if err == nil {
		t.Error("Accepted array of invalid length")
		t.Fail()
	}

	//test a slice that does not need substitutions
	result, err = MakeSubstitutions(CASE_1, DECODE)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

}
