package file

import (
	"bufio"
	"os"
	"regexp"
	"strings"
	"testing"
)

var (
	oldPoem = []string{
		"白日依山尽",
		"黄河入海流",
		"更上一层楼",
	}
	newPoem = []string{
		"白日依山尽",
		"黄河入海流",
		"欲穷千里目",
		"更上一层楼",
	}
	newPoem2 = []string{
		"白日依山尽",
		"黄河入海流",
		"欲穷千里目",
		"更上一层楼",
		"欲穷千里目",
	}
	newPoem3 = []string{
		"我是替换内容",
	}
	newPoem4 = []string{
		"白日依山尽",
		"黄河入海流",
		"君不见黄河之水天上来",
		"奔流到海不复回",
		"更上一层楼",
	}
)

//测试插入单个标签
func TestFile_Insert(t *testing.T) {
	f, _ := createTmpFile("./poem.txt")
	if file, err := NewFile(f.Name()); err != nil {
		t.Fatal(err)
	} else {
		if err := writePoem(f); err != nil {
			t.Error(err)
		}
		file.Scan([]AddTag{poemTag, poemTag2}...)

		//某个标签插入
		if err := file.Insert([]Tag{"insert"}, []string{"欲穷千里目"}); err != nil {
			t.Error(err)
		} else {
			if newFile, err := os.Open(f.Name()); err != nil {
				t.Error(err)
			} else {
				scanner := bufio.NewScanner(newFile)
				index := 0
				for scanner.Scan() {
					if scanner.Text() != newPoem[index] {
						t.Error(scanner.Text())
					} else {
						index++
					}
				}
			}
		}
	}
	os.Remove(f.Name())
}

//测试多个标签
func TestFile_InsertManyTags(t *testing.T) {
	f, _ := createTmpFile("./poem.txt")
	if file, err := NewFile(f.Name()); err != nil {
		t.Fatal(err)
	} else {
		if err := writePoem(f); err != nil {
			t.Error(err)
		}
		file.Scan([]AddTag{poemTag, poemTag2}...)

		//多个标签插入
		if err := file.Insert([]Tag{"insert", "insert2"}, []string{"欲穷千里目"}); err != nil {
			t.Error(err)
		} else {
			if newFile, err := os.Open(f.Name()); err != nil {
				t.Error(err)
			} else {
				scanner := bufio.NewScanner(newFile)
				index := 0
				for scanner.Scan() {
					if scanner.Text() != newPoem2[index] {
						t.Error(scanner.Text())
					} else {
						index++
					}
				}
			}
		}
	}
	os.Remove(f.Name())
}

//无标签测试替换
func TestFile_InsertTruncate(t *testing.T) {
	f, _ := createTmpFile("./poem.txt")
	if file, err := NewFile(f.Name()); err != nil {
		t.Fatal(err)
	} else {
		if err := writePoem(f); err != nil {
			t.Error(err)
		} else {
			file.Scan([]AddTag{poemTag, poemTag2}...)
			if err := file.Insert(nil, newPoem3); err != nil {
				t.Error(err)
			} else {
				if newFile, err := os.Open(f.Name()); err != nil {
					t.Error(err)
				} else {
					scanner := bufio.NewScanner(newFile)
					index := 0
					for scanner.Scan() {
						if scanner.Text() != newPoem3[index] {
							t.Error(scanner.Text())
						} else {
							index++
						}
					}
				}
			}
		}
	}
	os.Remove(f.Name())
}

//insert between测试
func TestFile_InsertBetween(t *testing.T) {
	f, _ := createTmpFile("./poem.txt")
	if file, err := NewFile(f.Name()); err != nil {
		t.Fatal(err)
	} else {
		if err := writePoem(f); err != nil {
			t.Error(err)
		} else {
			file.Scan([]AddTag{poemTag, poemTag2}...)
			content := []string{
				"君不见黄河之水天上来",
				"奔流到海不复回",
			}

			err := file.InsertBetween("", "", content)
			if err.Error() != EMPTY_TAG {
				os.Remove(f.Name())
				t.Fatal()
			}
			err = file.InsertBetween("x", "insert2", content)
			if err.Error() != START_TAG_NOT_EXISTS {
				os.Remove(f.Name())
				t.Fatal()
			}
			err = file.InsertBetween("insert", "x", content)
			if err.Error() != END_TAG_NOT_EXISTS {
				os.Remove(f.Name())
				t.Fatal()
			}
			err = file.InsertBetween("insert2", "insert", content)
			if err.Error() != TAG_MISS_MATCH {
				os.Remove(f.Name())
				t.Fatal()
			}

			if err := file.InsertBetween("insert", "insert2", content);
				err != nil {
				t.Error(err)
			} else {
				if newFile, err := os.Open(f.Name()); err != nil {
					t.Error(err)
				} else {
					scanner := bufio.NewScanner(newFile)
					index := 0
					for scanner.Scan() {
						if scanner.Text() != newPoem4[index] {
							t.Error(scanner.Text())
						} else {
							index++
						}
					}
				}
			}
		}
	}
	os.Remove(f.Name())
}

func writePoem(f *os.File) error {
	if _, err := f.WriteString(strings.Join(oldPoem, "\n")); err != nil {
		return err
	} else {
		return nil
	}
}
func poemTag(line int, content string) Tag {
	reg, _ := regexp.Compile(`^\s*黄河入海流\s*$`)
	if reg.FindString(content) != "" {
		return "insert"
	} else {
		return ""
	}
}
func poemTag2(line int, content string) Tag {
	reg, _ := regexp.Compile(`^\s*更上一层楼\s*$`)
	if reg.FindString(content) != "" {
		return "insert2"
	} else {
		return ""
	}
}
