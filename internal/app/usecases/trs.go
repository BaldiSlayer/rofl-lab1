package usecases

import "github.com/BaldiSlayer/rofl-lab1/pkg/trsparser"

func ExtractFormalTrs(request string) (trsparser.Trs, string, error) {
	return trsparser.Trs{}, "", nil
}

func FixFormalTrs(request, formalTrs string, parseError trsparser.ParseError) (trsparser.Trs, string, error) {
	return trsparser.Trs{}, "", nil
}

func InterpretFormalTrs(trsparser.Trs) (string, error) {
	return "", nil
}
