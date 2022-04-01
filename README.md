在前面章节中我们看到，语法生产式本质上指导了我们如何展开代码，例如对于生产式A->X Y Z,那么我们在解析A的时候，对于的代码就是：
```
func A() {
    X()
    Y()
    Z()
}
```
从逻辑上看，我们解析A时，实际上是将A分解成三部分X,Y,Z，然后依次解析X,Y,Z。这种将一个大元素分解成多个小元素，然后把小元素处理后将结果集合起来形成大元素处理结果的方法叫自上而下的解析。但这种解析方法存在一个问题，我们看这样的字符串规律：字符串要不只含有一个字符b，要不在b的后面包含任意个数的字符a,满足条件的字符串有"b", "ba", "baa", 依次类推，对该字符串组合规律的生产式定义如下：
```
A -> A "a" | ”b"
```
这样我们在解析时对应的代码为:
```
func A() {
    A()
    match("b")
}
```
其中match("b")表示当前读取的字符是否与"b"相等。我们看到代码有问题，那就是函数A执行时直接调用了它自己，于是就会形成无限递归最终以栈被撑爆结束。我们在前面章节中重新定义算术表达式时也有这个问题也就是：
```
list -> list + digit 
list -> list - digit 
digit -> "0" | "1" .... |"9"
```
同样list生产式也产生了左递归，因此我们的代码套路无法使用。这种情况叫语法定义的左递归，我们需要使用一些办法处理它，好在有固定的套路，其处理方法如下，例如有如下的左递归生产式：
```
X -> X Y Z | "x"
```
那么我们把 Y Z 用另一个非终结符α表示，也就是 α -> Y Z, 然后引入一个新的非终结符R,然后将生产式改下如下
```
X -> "x" R 
R -> αR | ε
```
其中ε表示”空“，也就是什么都不做。按照上面的方法，我们看看如何处理A -> A "a"这种情况，根据上面修改方法，生产式重新设置如下：
```
A -> "b" R
R -> "a" R | ε
```
我们慢慢品味一下，修改后的语法生产式其效果跟原来一样，同时它再也没有左递归的情形，当然它也产生另外一个问题， 那就是R -> "a" R | ε 包含了右递归，这种情况会在语法解析上产生新的问题，不过在这里我们先忽略。

有了上面的基础后，我们再次修改算术表达式的语法生产式，处理其中的歧义，处理左递归，最后我们给出它的解析代码。首先我们看看消除歧义后的算术表达式语法：
```
list -> list "+" digit   {print('+')}
list -> list "-" digit   {print('-')}
list -> digit
digit -> "0" {print("0")} | "1" {print("1")} |... | "9" {print("9")}
```
在上语法中，右边{}中的内容表示完成解析后所需要的操作，例如第一行的{print("+")}表示完成解析list -> list "+" digit 后，打印出字符”+"。由于语法中存在左递归，因此我们需要先处理。上面的语法在形式上为：
```
A -> A α 
A -> A β 
A -> γ
```
其中 α 对应  "+"  digit  {print("+")}, β对应 “-" digit {print("-")} , γ 对应 digit ,按照前面说法修改后的语法为：
```
 A -> γ R 
 R -> α R | β  R | ε
```
我们将 A = list, α =  "+"  digit  {print("+")}, β= “-" digit {print("-")} , γ = digit进行替换就有：
```
list -> digit rest 
rest -> "+" digit {print("+")} rest | "-" digit {print("-")} rest | ξ
digit -> "0" {print("0")} | ... | "9" {print("9")}
```
这里我们使用 rest来替代前面说的R，如果有不明白的地方，我们一会使用代码实现时，疑惑就会解开。我们再次打开parser.go,修改代码如下：
```
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

```
上面代码有几处需要注意，第一是代码针对print("+")对应的操作，我们代码中通过在控制台输出对应语法中的print操作，另外还需要注意的是rest->ε在代码中的实现，它实际上对应一个空操作，在代码里我们利用一个空的else{}来对应该生产式，我们看看主函数入口，在main.go中添加代码如下：
```
package main

import (
	"fmt"
	"io"
	"lexer"
	"simple_parser"
)

func main() {
	source := "9-5+2"
	my_lexer := lexer.NewLexer(source)
	parser := simple_parser.NewSimpleParser(my_lexer)
	err := parser.Parse()
	if err == io.EOF {
		fmt.Println("parsing success")
	} else {
		fmt.Println("source is ilegal : ", err)
	}
}

```
代码运行后输出为95-2+，该结构对应算术表达式9-5+2的后向表达形式，同时我们解决了上一节给出的语法歧义性。视频推演请在B站搜索Coding迪斯尼，更多有趣内容请扫描二维码：
![请添加图片描述](https://img-blog.csdnimg.cn/591eebc7dca544e186c0f8f22e3b8387.png)
