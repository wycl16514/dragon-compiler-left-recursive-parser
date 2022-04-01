package simple_parser

import (
	"errors"
	"fmt"
	"lexer"
)

type SimpleParser struct {
	lexer lexer.Lexer
}

func NewSimpleParser(lexer lexer.Lexer) *SimpleParser {
	return &SimpleParser{
		lexer: lexer,
	}
}

func (s *SimpleParser) list() error {
	/*
		根据生产式 list -> digit rest ，所以我们调用函数digit 和 rest
	*/
	err := s.digit()
	err = s.rest()

	return err
}

func (s *SimpleParser) digit() error {
	tok, err := s.lexer.Scan()
	if err != nil {
		return err
	}
	//判断当前读到的是否为数字字符
	if tok.Tag != lexer.NUM || len(s.lexer.Lexeme) > 1 {
		s := fmt.Sprintf("parsing digit error, got %s", s.lexer.Lexeme)
		return errors.New(s)
	}

	//digit -> "0" print("0")|.. ，匹配后执行print操作
	fmt.Print(s.lexer.Lexeme)

	return nil
}

func (s *SimpleParser) rest() error {
	tok, err := s.lexer.Scan()
	if err != nil {
		return err
	}

	if tok.Tag == lexer.PLUS {
		//rest -> "+" digit print("+") rest
		err = s.digit()
		if err != nil {
			return err
		}
		//执行操作 print("+")
		fmt.Print("+")
		err = s.rest()
		if err != nil {
			return err
		}
	} else if tok.Tag == lexer.MINUS {
		//rest -> "-" digit print("-") rest
		err = s.digit()
		if err != nil {
			return err
		}
		//执行操作 print("-")
		fmt.Print("-")
		err = s.rest()
		if err != nil {
			return err
		}
	} else {
		//rest -> ε , 这里对应空操作
	}

	return err
}

func (s *SimpleParser) Parse() error {
	return s.list()
}
