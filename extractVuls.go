package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

// 结构体定义组件的信息
type ComponentInfo struct {
	Component   string `json:"component"`
	Version     string `json:"version"`
	Path        string `json:"path"`
	Severity    string `json:"severity"`
	Description string `json:"desc"`
}

type SiteComponent struct {
	Url        string          `json:"url"`
	Components []ComponentInfo `json:"components"`
}

// 解析输入字符串的函数
func parseVuls(inputString string) []ComponentInfo {
	// fmt.Println(inputString)
	lines := strings.Split(strings.TrimSpace(inputString), "\n")
	numberRegex := regexp.MustCompile(`\d`)
	var results []ComponentInfo
	seenComponents := make(map[string]bool)

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		parts := strings.Fields(line)

		if len(parts) > 2 && numberRegex.MatchString(parts[1]) && strings.Contains(parts[2], "http") {
			component := parts[0]
			version := parts[1]
			path := parts[2]

			// 解析下一行
			if i+1 < len(lines) {
				nextLineParts := strings.Fields(strings.TrimSpace(lines[i+1]))
				if len(nextLineParts) > 1 {
					// 判断是漏洞还是新的组件
					if numberRegex.MatchString(nextLineParts[1]) && strings.Contains(nextLineParts[2], "http") {
						// 新组件，那么将旧组件完结
						uniqueKey := component + version
						if seenComponents[uniqueKey] {
							continue
						}

						results = append(results, ComponentInfo{
							Component:   component,
							Version:     version,
							Path:        path,
							Severity:    "",
							Description: "",
						})
						seenComponents[uniqueKey] = true
						continue
					}

					severity := nextLineParts[0]
					description := strings.Join(nextLineParts[1:], " ")

					uniqueKey := component + version
					if seenComponents[uniqueKey] {
						continue
					}

					results = append(results, ComponentInfo{
						Component:   component,
						Version:     version,
						Path:        path,
						Severity:    severity,
						Description: description,
					})
					seenComponents[uniqueKey] = true
					// fmt.Println("################################")
					// fmt.Println("Component: ", component)
					// fmt.Println("Version: ", version)
					// fmt.Println("Path: ", path)
					// fmt.Println("Severity: ", severity)
					// fmt.Println("Description: ", description)

					// 跳到下下行
					i++
				}
			} else {
				uniqueKey := component + version
				if seenComponents[uniqueKey] {
					continue
				}

				results = append(results, ComponentInfo{
					Component:   component,
					Version:     version,
					Path:        path,
					Severity:    "",
					Description: "",
				})
				seenComponents[uniqueKey] = true
			}
		}
	}

	return results
}

func GetListFromFile(path string) []string {
	// 打开文本文件
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("无法打开文件: %s", err)
	}
	defer file.Close()

	// 创建一个 scanner 来读取文件的每一行
	scanner := bufio.NewScanner(file)
	var lines []string

	// 读取文件的每一行并添加到字符串切片中
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	// 检查 scanner 是否出现错误
	if err := scanner.Err(); err != nil {
		log.Fatalf("读取文件时出错: %s", err)
	}

	return lines
}
