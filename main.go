package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var (
	signs    = []string{"+", "-", "/", "*"}
	romsNum  = []string{"I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X"}
	romsNum2 = []string{"XI", "XV", "XX", "XXX", "XL", "L", "LX", "LXX", "LXXX", "XC", "C", "CC", "CCC", "CD", "D", "DC", "DCC", "DCCC", "CM", "M"}
	units    = map[string]string{"0": "", "1": "I", "2": "II", "3": "III", "4": "IV", "5": "V", "6": "VI", "7": "VII", "8": "VIII", "9": "IX", "10": "X"}
	dozens   = map[string]string{"1": "X", "2": "XX", "3": "XXX", "4": "XL", "5": "L", "6": "LX", "7": "LXX", "8": "LXXX", "9": "XC"}
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("Введите операцию: ")

	// отслеживаем ввод с консоли
	for scanner.Scan() {
		line := trim(scanner.Text())
		nums := strings.Fields(line)
		numsType := check(nums)

		if numsType == "" {
			fmt.Println(errors.New("строка пустая"))
		} else {
			arith(nums, numsType)
		}
		fmt.Println("Введите ещё операцию: ")
	}
}

// Функция удаления лишних символов и проверка формата
func trim(str string) string {
	re1 := regexp.MustCompile(`([*/+-])`)
	re2 := regexp.MustCompile(`^- `)
	re3 := regexp.MustCompile(`^\+ |[[:^ascii:]]`)

	str = strings.ReplaceAll(str, ",", ".")
	str = re1.ReplaceAllString(str, " $1 ")
	str = strings.TrimSpace(str)
	str = strings.ToUpper(str)

	for strings.Contains(str, "  ") {
		str = strings.Replace(str, "  ", " ", -1)
	}

	str = re2.ReplaceAllString(str, "-")
	str = re3.ReplaceAllString(str, "")

	if len(strings.Fields(str)) > 3 {
		panic(errors.New("формат не удовлетворяет заданию" +
			" — два операнда и один оператор (+, -, /, *)"))
	}
	if len(strings.Fields(str)) < 3 {
		panic(errors.New("строка не является математической операцией"))
	}
	return str
}

// Функция проверки корректности введенных систем счисления
func check(opr []string) string {

	// проверка наличия знака оператора
	cSign := slices.Contains(signs, opr[1])

	// проверка наличия римских цифр
	cRom1, cRom2 := slices.Contains(romsNum, opr[0]), slices.Contains(romsNum, opr[2])

	// проверка наличия целых или дробных чисел
	regInt := regexp.MustCompile(`^-[0-9]+$|^[0-9]+$`)
	regFloat := regexp.MustCompile(`^[0-9]*[.,][0-9]+$|^-[0-9]*[.,][0-9]+$`)
	cInt1, cInt2 := regInt.MatchString(opr[0]), regInt.MatchString(opr[2])
	cFloat1, cFloat2 := regFloat.MatchString(opr[0]), regFloat.MatchString(opr[2])
	romConv1, romConv2 := romToInt(opr[0]), romToInt(opr[2])

	// если строка не содержит знак оператора
	if !cSign {
		panic(errors.New("строка не является математической операцией"))
	}

	// проверка отрицательных римских цифр и условия - число не больше 10
	if romConv1 != 0 && romConv2 != 0 {
		if romConv1 > 10 || romConv2 > 10 {
			panic(errors.New("калькулятор принимает на вход числа от 1(I) до 10(X)"))
		}
		if strings.Contains(opr[0], "-") {
			panic(errors.New("в римской системе нет отрицательных чисел"))
		}
	}
	// проверка на сочетания любого отрицательного римского и любого арабского числа
	if romConv1 != 0 && (cInt2 || cFloat2) {
		if strings.Contains(opr[0], "-") {
			panic(errors.New("используются одновременно разные системы счисления"))
		}
	}
	// проверка на сочетания римского числа больше 10 и любого арабского числа
	if (romConv1 > 10 && cInt2) ||
		(cInt1 && romConv2 > 10) ||
		(romConv1 > 10 && cFloat2) || (cFloat1 && romConv2 > 10) {
		panic(errors.New("используются одновременно разные системы счисления"))
	}

	// проверка на сочетание подходящих римских, арабских целых и дробных чисел
	switch {
	case (cRom1 && cInt2) || (cInt1 && cRom2) ||
		(cRom1 && cFloat2) || (cFloat1 && cRom2):
		panic(errors.New("используются одновременно разные системы счисления"))
	case (cInt1 && cFloat2) || (cFloat1 && cInt2):
		panic(errors.New("калькулятор умеет работать только с целыми числами"))
	case cInt1 && cInt2:
		a, e1 := strconv.Atoi(opr[0])
		b, e2 := strconv.Atoi(opr[2])
		if e1 != nil && e2 != nil {
			fmt.Printf("Error 1: %d\n Error 2: %d", e1, e2)
		}
		if a > 10 || b > 10 || a == 0 || b == 0 {
			panic(errors.New("калькулятор принимает на вход числа от 1 до 10"))
		}
		return "integer"
	case cRom1 && cRom2:
		return "roman"
	default:
		panic(errors.New("строка не является математической операцией"))
	}
}

// Функция арифметических операций
func arith(nums []string, numsType string) {

	var a, b, res int

	// преобразование если числа римские
	if numsType == "roman" {
		a, b = romToInt(nums[0]), romToInt(nums[2])
		if (b > a && nums[1] == "/") || (b == a && nums[1] == "-") {
			panic(errors.New("в римской системе могут быть только положительные числа больше нуля"))
		}
		if b > a && nums[1] == "-" {
			panic(errors.New("в римской системе нет отрицательных чисел"))
		}
	}
	// преобразование если числа арабские
	if numsType == "integer" {
		aNum, err1 := strconv.Atoi(nums[0])
		bNum, err2 := strconv.Atoi(nums[2])
		if err1 != nil && err2 != nil {
			fmt.Printf("Error 1: %d\n Error 2: %d", err1, err2)
		} else {
			a, b = aNum, bNum
		}
	}

	// основные операции
	switch {
	case nums[1] == "+":
		res = a + b
	case nums[1] == "-":
		res = a - b
	case nums[1] == "*":
		res = a * b
	case nums[1] == "/":
		res = a / b
	}

	// выводим результат в консоли
	if numsType == "roman" {
		fmt.Printf("Результат: \n%v\n", intToRom(res))
	} else {
		fmt.Printf("Результат: \n%v\n", res)
	}
}

// Функция - перевод арабского числа в римское
func intToRom(num int) string {

	// перевод числа в строку и создание среза чисел
	toStr := strconv.Itoa(num)
	number := strings.Split(toStr, "")
	result := ""

	// обход по карте, поиск совпадений
	for i, value := range number {
		// число больше 100 не может быть по условию
		if num == 100 {
			result = "C"
		}
		if num > 100 {
			panic(errors.New("результат больше 100, " +
				"калькулятор принимает на вход только числа от 1 до 10"))
		}
		// обходим десятки, потом единицы, добавляя в строку result
		if len(number) == 2 {
			if i == 0 {
				val, ok := dozens[value]
				if ok {
					result = result + val
				}
			} else {
				val, ok := units[value]
				if ok {
					result = result + val
				}
			}
		} else {
			val, ok := units[value]
			if ok {
				result = result + val
			}
		}
	}
	return result
}

// Функция - перевод римского числа в арабское
func romToInt(num string) int {

	// чистка лишних знаков
	num = strings.ReplaceAll(num, "-", "")
	num = strings.ReplaceAll(num, " ", "")

	// проверка корректности римских цифр
	regNum := regexp.MustCompile(`^[IVXLCDM]+$`)
	regInt := regexp.MustCompile(`^\d*[a-zA-Z][a-zA-Z\d]*$`)
	isRom := regNum.MatchString(num)
	isRom2 := regInt.MatchString(num)
	overCnt := 0

	// если числа нет в римской системе
	if !isRom {
		if isRom2 {
			panic(fmt.Errorf("в римской системе нет числа: %q", num))
		}
		return 0
	}

	// обход по карте, поиск совпадений чисел от 1 до 10
	for ku, vu := range units {
		if vu == num {
			res, err := strconv.Atoi(ku)
			if err != nil {
				panic(err)
			}
			return res
		}
	}

	// обход по карте, поиск совпадений чисел больше 10
	for _, rn2 := range romsNum2 {
		regOver := regexp.MustCompile(`^` + rn2 + `.*`)
		isOver := regOver.MatchString(num)

		if isOver {
			overCnt++
		}
	}

	// если есть числа больше 10, возвращаем условное число 11
	if overCnt > 0 {
		return 11
	} else {
		panic(fmt.Errorf("в римской системе нет числа: %q", num))

	}
}
