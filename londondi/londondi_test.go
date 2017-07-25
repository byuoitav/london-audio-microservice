package londondi

import (
	"bytes"
	"fmt"
	"log"
	"testing"
)

//decoding cases
var CASE_0 = []byte{0x8d, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x8c}
var RESULT_0 = []byte{0x8d, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x8c}

var CASE_1 = []byte{0x8d, 0x00, 0x01, 0x1b, 0x83, 0x00, 0x01, 0x1b, 0x82, 0x00, 0x00, 0x00, 0x1d, 0xff, 0xec}
var RESULT_1 = []byte{0x8d, 0x00, 0x01, 0x03, 0x00, 0x01, 0x02, 0x00, 0x00, 0x00, 0x1d, 0xff, 0xec}

var CASE_2 = []byte{0x8d, 0x00, 0x01, 0x1b, 0x83, 0x00, 0x01, 0x00, 0x00, 0x00, 0x1d, 0xff, 0xec, 0x82}
var RESULT_2 = []byte{0x8d, 0x00, 0x01, 0x03, 0x00, 0x01, 0x00, 0x00, 0x00, 0x1d, 0xff, 0xec, 0x82}

var CASE_3 = []byte{0x02, 0x8d, 0x00, 0x01, 0x1b, 0x83, 0x00, 0x01, 0x1b, 0x82, 0x00, 0x00, 0x00, 0x1d, 0xff, 0xec, 0x82, 0x03}
var RESULT_3 = []byte{0x02, 0x8d, 0x00, 0x01, 0x03, 0x00, 0x01, 0x02, 0x00, 0x00, 0x00, 0x1d, 0xff, 0xec, 0x82, 0x03}

//encoding cases
var CASE_4 = []byte{0x02}
var RESULT_4 = []byte{0x1b, 0x82}

var CASE_5 = []byte{0x02, 0x1b, 0x82}
var RESULT_5 = []byte{0x1b, 0x82, 0x1b, 0x9b, 0x82}

var CASE_6 = []byte{0x00, 0x01, 0x04}
var RESULT_6 = []byte{0x00, 0x01, 0x04}

//validating cases
var CASE_7 = []byte{0x02, 0x8d, 0x23, 0xda, 0x1b, 0x83, 0x00, 0x01, 0x1b, 0x82, 0x00, 0x00, 0x00, 0x31, 0xff, 0xa1, 0x1b, 0x9b, 0x03}
var RESULT_7 = []byte{0x8d, 0x23, 0xda, 0x03, 0x00, 0x01, 0x02, 0x00, 0x00, 0x00, 0x31, 0xff, 0xa1, 0x1b}

var CASE_8 = []byte{0x02, 0x03} //begins with STX and ends with ETX
var CASE_9 = []byte{0x02, 0x8d, 0x23, 0xda, 0x1b, 0x83, 0x02, 0x01, 0x1b, 0x82, 0x00, 0x00, 0x00, 0x31, 0xff, 0xa1, 0x1b, 0x9b, 0x03}

//Wrap results
var WRAP_1 = []byte{0x02, 0x00, 0x01, 0x04, 0x03}

//Unwrap results
var UNWRAP_1 = []byte{}

func TestWrap(t *testing.T) {

	//test pre-conditions
	result, err := Wrap(CASE_5)
	if err == nil {
		t.Error("Allowed message to contain STX byte")
		t.Fail()
	}

	result, err = Wrap(CASE_6)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	//test post conditions
	if !bytes.Equal(result, WRAP_1) {
		msg := fmt.Sprintf("Expected %x, returned %x", WRAP_1, result)
		t.Error(msg)
	}

}

func TestUnwrap(t *testing.T) {

	//test pre-conditions
	result, err := Unwrap(CASE_0)
	if err == nil {
		t.Error("Allowed message that did not begin with STX and end with ETX")
		t.Fail()
	}

	result, err = Unwrap(CASE_8)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if len(result) != 0 {
		msg := fmt.Sprintf("Expected %x, returned %x", UNWRAP_1, result)
		t.Error(msg)
	}

	result, err = Unwrap(CASE_9)
	if err == nil {
		t.Error("Did not throw error when message contained erroneous STX or ETX bytes")
		t.Fail()
	}
}

func TestSubstitutions(t *testing.T) {

	fmt.Printf("\n\nTesting decoding...\n\n")

	//test decoding a slice where there are no subsitutions necessary when decoding, e.g. no instances of escape byte
	log.Printf("Case 0: %x", CASE_0)
	result, err := MakeSubstitutions(CASE_0, DECODE)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !bytes.Equal(result, RESULT_0) {
		msg := fmt.Sprintf("Expected %x, returned %x", RESULT_0, result)
		t.Error(msg)
	}

	//slice with substitutions 0x1b82 -> 0x02, 0x1b83 -> 0x03
	log.Printf("Case 1: %x", CASE_1)
	result, err = MakeSubstitutions(CASE_1, DECODE)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !bytes.Equal(result, RESULT_1) {
		msg := fmt.Sprintf("Expected %x, returned %x", RESULT_1, result)
		t.Error(msg)
	}

	//slice with unescaped subsitution (no substitutions)
	log.Printf("Case 2: %x", CASE_2)
	result, err = MakeSubstitutions(CASE_2, DECODE)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !bytes.Equal(result, RESULT_2) {
		msg := fmt.Sprintf("Expected %x, returned %x", RESULT_2, result)
		t.Error(msg)
	}

	//the kitchen sink
	log.Printf("Case 3: %x", CASE_3)
	result, err = MakeSubstitutions(CASE_3, DECODE)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !bytes.Equal(result, RESULT_3) {
		msg := fmt.Sprintf("Expected %x, returned %x", RESULT_3, result)
		t.Error(msg)
	}

	fmt.Printf("\n\nTesting encoding...\n\n")

	log.Printf("Case 4: %x", CASE_4)
	result, err = MakeSubstitutions(CASE_4, ENCODE)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !bytes.Equal(result, RESULT_4) {
		msg := fmt.Sprintf("Expected %x, returned %x", RESULT_4, result)
		t.Error(msg)
	}

	log.Printf("Case 5: %x", CASE_5)
	result, err = MakeSubstitutions(CASE_5, ENCODE)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !bytes.Equal(result, RESULT_5) {
		msg := fmt.Sprintf("Expected %x, returned %x", RESULT_5, result)
		t.Error(msg)
	}

	log.Printf("Case 6: %x", CASE_6)
	result, err = MakeSubstitutions(CASE_6, ENCODE)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !bytes.Equal(result, RESULT_6) {
		msg := fmt.Sprintf("Expected %x, returned %x", RESULT_6, result)
		t.Error(msg)
	}

	//FINAL ROUND
	fmt.Printf("\n\nReversibility Test: %x\n\n", CASE_1)
	result, err = MakeSubstitutions(CASE_1, DECODE)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	log.Printf("Result: %x", result)
	result, err = MakeSubstitutions(result, ENCODE)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !bytes.Equal(CASE_1, result) {
		msg := fmt.Sprintf("Expecting %x, returned %x", CASE_1, result)
		t.Error("Function not reversible! " + msg)
	}

}

func TestValidate(t *testing.T) {
}
