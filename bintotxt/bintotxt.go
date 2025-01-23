package main

import (
	"encoding/binary" // Для чтения бинарных данных из файла.
	"flag"           // Для парсинга аргументов командной строки.
	"fmt"            // Для форматированного ввода/вывода.
	"io"             // Для работы с вводом/выводом.
	"os"             // Для работы с файловой системой.
)

// Edge представляет ребро в графе.
type Edge struct {
	from   int16  // Начальная вершина ребра.
	to     int16  // Конечная вершина ребра.
	weight int16  // Вес ребра.
}

// readGraph считывает граф из бинарного файла.
// Формат файла: сначала количество вершин (int16),
// затем тройки (from, to, weight) в формате int16 каждая.
func readGraph(filename string) (int16, []Edge, error) {
	// Открываем бинарный файл для чтения.
	file, err := os.Open(filename)
	if err != nil {
		return 0, nil, fmt.Errorf("не удалось открыть файл %s: %v", filename, err)
	}
	defer file.Close() // Закрываем файл при выходе из функции.

	var numVertices int16
	// Читаем количество вершин.
	err = binary.Read(file, binary.LittleEndian, &numVertices)
	if err != nil {
		return 0, nil, fmt.Errorf("не удалось прочитать количество вершин: %v", err)
	}

	edges := []Edge{} // Срез для хранения ребер.

	for {
		var from, to, weight int16
		// Читаем тройку (from, to, weight).
		err = binary.Read(file, binary.LittleEndian, &from)
		if err == io.EOF {
			break // Достигнут конец файла.
		}
		if err != nil {
			return 0, nil, fmt.Errorf("ошибка при чтении вершины 'from': %v", err)
		}

		err = binary.Read(file, binary.LittleEndian, &to)
		if err != nil {
			return 0, nil, fmt.Errorf("ошибка при чтении вершины 'to': %v", err)
		}

		err = binary.Read(file, binary.LittleEndian, &weight)
		if err != nil {
			return 0, nil, fmt.Errorf("ошибка при чтении веса ребра: %v", err)
		}

		// Проверка корректности индексов вершин.
		if from < 0 || from >= numVertices || to < 0 || to >= numVertices {
			return 0, nil, fmt.Errorf("недопустимые индексы вершин: from=%d, to=%d", from, to)
		}

		// Добавляем ребро в срез.
		edges = append(edges, Edge{from, to, weight})
	}

	return numVertices, edges, nil // Возвращаем количество вершин и список ребер.
}

// writeEdges записывает все ребра в текстовый файл.
// Каждое ребро записывается в формате (from, to, weight) на отдельной строке.
func writeEdges(filename string, edges []Edge) error {
	// Создаём (или перезаписываем) выходной текстовый файл.
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("не удалось создать файл %s: %v", filename, err)
	}
	defer file.Close() // Закрываем файл при выходе из функции.

	// Записываем заголовок.
	_, err = fmt.Fprintln(file, "Список ребер графа:")
	if err != nil {
		return fmt.Errorf("ошибка при записи заголовка: %v", err)
	}

	// Записываем каждое ребро.
	for _, edge := range edges {
		// Предполагается, что вершины нумеруются с 0, добавляем 1 для удобства чтения.
		line := fmt.Sprintf("(%d, %d, %d)", edge.from, edge.to, edge.weight)
		_, err = fmt.Fprintln(file, line)
		if err != nil {
			return fmt.Errorf("ошибка при записи ребра: %v", err)
		}
	}

	return nil // Возвращаем nil, если всё прошло успешно.
}

// main является точкой входа в программу.
// Он парсит аргументы командной строки, считывает граф из бинарного файла,
// и записывает все ребра в текстовый файл.
func main() {
	// Определяем флаги командной строки.
	inputFile := flag.String("i", "", "Имя входного бинарного файла")
	outputFile := flag.String("o", "output.txt", "Имя выходного текстового файла")
	flag.Parse() // Парсим флаги.

	// Проверяем, что имя входного файла было предоставлено.
	if *inputFile == "" {
		fmt.Println("Использование: go run read_graph.go -i inputfile [-o outputfile]")
		return
	}

	// Считываем граф из бинарного файла.
	numVertices, edges, err := readGraph(*inputFile)
	if err != nil {
		fmt.Printf("Ошибка при чтении графа: %v\n", err)
		return
	}

	fmt.Printf("Граф успешно считан: %d вершин, %d ребер.\n", numVertices, len(edges))

	// Записываем ребра в текстовый файл.
	err = writeEdges(*outputFile, edges)
	if err != nil {
		fmt.Printf("Ошибка при записи ребер: %v\n", err)
		return
	}

	fmt.Printf("Ребра успешно записаны в файл %s.\n", *outputFile)
}
