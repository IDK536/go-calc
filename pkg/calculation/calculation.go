package calculation

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type addition struct {
	numbers []float64
	znak    []string
}

type staples struct {
	staples []string
	znak    []string
}

func parseExpression(expression string) (addition, error) {
	a := addition{}
	znakinx := 0

	for i := range expression {
		if strings.Contains("+-*/", string(expression[i])) {
			intnumber, err := strconv.ParseFloat(strings.TrimSpace(expression[znakinx:i]), 64)
			if err != nil {
				return a, err
			}
			a.numbers = append(a.numbers, intnumber)
			a.znak = append(a.znak, string(expression[i]))
			znakinx = i + 1
		}
	}

	if znakinx < len(expression) {
		intnumber, err := strconv.ParseFloat(strings.TrimSpace(expression[znakinx:]), 64)
		if err != nil {
			return a, err
		}
		a.numbers = append(a.numbers, intnumber)
	}
	return a, nil
}

func calculateMulDiv(a *addition) {
	for i := 0; i < len(a.znak); i++ {
		if a.znak[i] == "*" {
			a.numbers[i] = a.numbers[i] * a.numbers[i+1]
			a.numbers = append(a.numbers[:i+1], a.numbers[i+2:]...)
			a.znak = append(a.znak[:i], a.znak[i+1:]...)
			i--
		} else if a.znak[i] == "/" {
			a.numbers[i] = a.numbers[i] / a.numbers[i+1]
			a.numbers = append(a.numbers[:i+1], a.numbers[i+2:]...)
			a.znak = append(a.znak[:i], a.znak[i+1:]...)
			i--
		}
	}
}

func calculateAddSub(a *addition) float64 {
	calc := a.numbers[0]
	for i := 0; i < len(a.znak); i++ {
		if a.znak[i] == "+" {
			calc += a.numbers[i+1]
		} else if a.znak[i] == "-" {
			calc -= a.numbers[i+1]
		}
	}
	return calc
}

func Calc_without_brackets(expression string) (float64, error) {
	a, err := parseExpression(expression)
	// log.Println("-1 in ", a)
	if err != nil {
		return 0, err
	}

	calculateMulDiv(&a)
	result := calculateAddSub(&a)
	return result, nil
}

func splitByParentheses(s string) ([]string, error) {
	var result []string
	var current string
	var stack []rune

	for _, char := range s {
		if char == '(' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
			stack = append(stack, char)
		} else if char == ')' {
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
				if len(stack) == 0 {
					if current != "" {
						result = append(result, current)
						current = ""
					}
				}
			} else {
				return []string{}, ErrInvalidExpression
			}
		} else {
			current += string(char)
		}
	}

	if current != "" {
		result = append(result, current)
	}
	return result, nil
}

func hasConsecutiveOperators(s string) bool {
	operators := "+-*/"

	for i := 0; i < len(s)-1; i++ {
		if strings.ContainsRune(operators, rune(s[i])) && strings.ContainsRune(operators, rune(s[i+1])) {
			return true
		}
	}
	return false
}

func Calc(expression string) (float64, error) {
	expression = strings.Replace(expression, " ", "", -1)
	// log.Println(expression)
	re := regexp.MustCompile(`[^1234567890+\-/*()]`)

	if re.MatchString(expression) {
		return 0, ErrInvalidExpression
	}
	if len(expression) == 0 {
		return 0, ErrEmptyExpression
	}
	if hasConsecutiveOperators(expression) {
		return 0, ErrInvalidExpression
	}
	// if hasConsecutiveOperators(expression) {
	// 	return 0, errors.New("nenene")
	// }
	s := staples{}
	buf, err1 := splitByParentheses(expression)
	// log.Println(buf)
	if err1 != nil {
		return 0, err1
	}
	if string(expression[0]) == "+" || string(expression[0]) == "-" || string(expression[0]) == "/" || string(expression[0]) == "*" || string(expression[len(expression)-1]) == "+" || string(expression[len(expression)-1]) == "-" || string(expression[len(expression)-1]) == "*" || string(expression[len(expression)-1]) == "/" {
		return 0, ErrInvalidExpression
	}
	for part := range buf {
		// log.Println("parts ", buf[part])
		if string(buf[part][0]) == "+" || string(buf[part][0]) == "-" || string(buf[part][0]) == "*" || string(buf[part][0]) == "/" {
			// s.staples = append(s.staples, buf[part][1:])
			s.znak = append(s.znak, string(buf[part][0]))
			bufa, err := parseExpression(buf[part][1:])
			if err != nil {
				return 0, err
			}
			// log.Println("bufa ", bufa)
			s.znak = append(s.znak, bufa.znak...)
			for i := range bufa.numbers {
				s.staples = append(s.staples, fmt.Sprint(bufa.numbers[i]))
			}
			if string(buf[part][0]) == "-" {
				s.staples[len(s.staples)-1] = strings.Replace(s.staples[len(s.staples)-1], "+", "?", -1)
				s.staples[len(s.staples)-1] = strings.Replace(s.staples[len(s.staples)-1], "-", "+", -1)
				s.staples[len(s.staples)-1] = strings.Replace(s.staples[len(s.staples)-1], "?", "-", -1)
			}
		} else if string(buf[part][len(string(buf[part]))-1]) == "+" || string(buf[part][len(string(buf[part]))-1]) == "-" || string(buf[part][len(string(buf[part]))-1]) == "*" || string(buf[part][len(string(buf[part]))-1]) == "/" {
			// s.staples = append(s.staples, buf[part][:len(string(buf[part]))-1])
			bufa, err := parseExpression(buf[part][:len(string(buf[part]))-1])
			if err != nil {
				return 0, err
			}
			// log.Println("bufaa ", bufa)
			s.znak = append(s.znak, bufa.znak...)
			for i := range bufa.numbers {
				s.staples = append(s.staples, fmt.Sprint(bufa.numbers[i]))
			}
			s.znak = append(s.znak, string(buf[part][len(string(buf[part]))-1]))
			if string(buf[part][len(string(buf[part]))-1]) == "-" {
				s.staples[len(s.staples)-1] = strings.Replace(s.staples[len(s.staples)-1], "+", "?", -1)
				s.staples[len(s.staples)-1] = strings.Replace(s.staples[len(s.staples)-1], "-", "+", -1)
				s.staples[len(s.staples)-1] = strings.Replace(s.staples[len(s.staples)-1], "?", "-", -1)
			}
		} else {
			res1, err := Calc_without_brackets(buf[part])
			if err != nil {
				return 0, err
			}
			s.staples = append(s.staples, fmt.Sprint(res1))
			// log.Println("pop1 " + fmt.Sprint(res1))
		}
		// log.Println("pop " + buf[part])
	}
	var StringWthDegeneration string
	// log.Println("Err ", s.staples, s.znak)
	for i := range s.znak {
		StringWthDegeneration += s.staples[i] + s.znak[i]
	}
	StringWthDegeneration += s.staples[len(s.staples)-1]
	// log.Println("tre ", StringWthDegeneration)

	a := addition{}
	a.znak = s.znak
	for i := range s.staples {
		bufper, err := strconv.ParseFloat(s.staples[i], 64)
		if err != nil {
			return 0, err
		}
		a.numbers = append(a.numbers, bufper)
	}
	// for i := range s.staples {
	// 	bufs, _ := Calc_without_brackets(s.staples[i])
	// 	a.numbers = append(a.numbers, bufs)
	// }
	// log.Println(a)
	// calculateMulDiv(&a)
	// result1 := calculateAddSub(&a)
	// var result1 float64
	// if len(s.znak) != 0 {
	// 	result1, _ = Calc_without_brackets(StringWthDegeneration)

	// } else {
	// 	result1, _ = strconv.ParseFloat(StringWthDegeneration, 64)
	// }1/0

	calculateMulDiv(&a)
	result1 := calculateAddSub(&a)

	if result1 == math.Inf(1) || result1 == math.Inf(-1) {
		return 0, ErrDivisionByZero
	}

	return result1, nil
}

// func main() {
// 	result, err := Calc("1/1/0")
// 	fmt.Println(result)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}
// 	fmt.Println("Result:", result)
// }
