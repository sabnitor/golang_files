package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type txtObj struct {
	text     string
	strategy strategy
}

func (t *txtObj) Action() {
	t.strategy.Action(t.text)
}

func (t *txtObj) ChangeStrategy(newStrategy strategy) {
	t.strategy = newStrategy
}

func newTxtObj(text string, strategy strategy) *txtObj {
	return &txtObj{
		text:     text,
		strategy: strategy,
	}
}

type strategy interface {
	Action(text string)
}

type wordCount struct{}

func (w *wordCount) Action(text string) {
	// Опрацьовуємо текст і виконуємо підрахунок слів або інші дії
	words := strings.Fields(text)
	wordCount := len(words)
	fmt.Printf("Кількість слів: %d\n", wordCount)
}

type MostRepeatedWordsStrategy struct{}

func (w *MostRepeatedWordsStrategy) Action(text string) {
	// Опрацьовуємо текст і виконуємо знаходження найбільше повторюваних слів або інші дії
	words := strings.Fields(text)
	wordCount := make(map[string]int)
	for _, word := range words {
		wordCount[word]++
	}

	var maxWord string
	var maxCount int
	for word, count := range wordCount {
		if count > maxCount {
			maxWord = word
			maxCount = count
		}
	}

	fmt.Printf("Most repeated word: %s (repeated %d times)\n", maxWord, maxCount)
}

func removeExtraSpaces(text string) string {
	// Видаляє всі послідовності пробілів, окрім символів нового рядка
	re := regexp.MustCompile(`[^\S\n]+`)
	cleanedText := re.ReplaceAllString(text, " ")

	// Видаляє пробіли з початку і кінця рядка
	cleanedText = regexp.MustCompile(`(^\s+|\s+$)`).ReplaceAllString(cleanedText, "")

	return cleanedText
}

func extraSpacesDecorator(action strategy) strategy {
	return &extraSpacesDecoratorStrategy{strategy: action}
}

type extraSpacesDecoratorStrategy struct {
	strategy strategy
}

func (e *extraSpacesDecoratorStrategy) Action(text string) {
	text = removeExtraSpaces(text)
	fmt.Println(text)
	e.strategy.Action(text)
}

func main() {
	filename := flag.String("input", "input.txt", "Ім'я вхідного файлу")
	strategyType := flag.String("strategy", "wordCount", "Виберіть стратегію (wordCount або MostRepeatedWords)")
	decorator := flag.String("decorator", "none", "Видалити зайві пробіли removeExtraSpaces(за замовчуванням none)")

	flag.Parse()

	// Читаємо файл
	data, err := os.ReadFile(*filename)
	if err != nil {
		fmt.Println("Помилка читання файлу:", err)
		return
	}

	text := string(data)

	// Обираємо стратегію
	var selectedStrategy strategy

	switch *strategyType {
	case "wordCount":
		selectedStrategy = &wordCount{}
	case "MostRepeatedWords":
		selectedStrategy = &MostRepeatedWordsStrategy{}
	default:
		fmt.Println("Невідома стратегія")
		return
	}

	// Створюємо новий текстовий об'єкт
	newTO := newTxtObj(text, selectedStrategy)

	// Обираємо декоратор
	switch *decorator {
	case "none":

	case "removeExtraSpaces":
		newTO.strategy = extraSpacesDecorator(newTO.strategy)
	default:
		fmt.Println("Невідомий декоратор")
		return
	}

	newTO.Action()
}
