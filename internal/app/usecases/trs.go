package usecases

import "github.com/BaldiSlayer/rofl-lab1/pkg/trsparser"

func ExtractFormalTrs(request string) (trsparser.Trs, string, error) {

	// TODO: сходить в клиент formalize

	// TODO: сходить в парсер

	return trsparser.Trs{}, "", nil
}

func FixFormalTrs(request, formalTrs string, parseError trsparser.ParseError) (trsparser.Trs, string, error) {

	// TODO: сходить в клиент formalize

	// TODO: сходить в парсер

	return trsparser.Trs{}, "", nil
}

func InterpretFormalTrs(trsparser.Trs) (string, error) {

	// TODO: сходить в клиент interpret

	return "", nil
}
