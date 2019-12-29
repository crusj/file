package file

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	EMPTY_TAG            = "tag miss"
	START_TAG_NOT_EXISTS = "start tag not exists"
	END_TAG_NOT_EXISTS   = "end tag not exists"
	TAG_MISS_MATCH       = "no match tag pair"
)

type file struct {
	Path string
	//所有内容行
	Lines []*Line
	//总行数
	Total int
	//大小
	Size int64
	//已打上的标签
	Tags map[Tag]Empty
}

//创建文件对象
func NewFile(path string) (*file, error) {
	file := new(file)
	file.Tags = make(map[Tag]Empty)
	if absPath, err := filepath.Abs(path); err != nil {
		return nil, err
	} else {
		file.Path = absPath
	}
	if fileInfo, err := os.Stat(file.Path); os.IsNotExist(err) {
		return nil, err
	} else {
		file.Size = fileInfo.Size()
	}
	return file, nil
}

type Tag = string
type Empty struct{}
type Line struct {
	Content string
	Tags    map[Tag]Empty
}

//为文件某一行打上标签，需要用户实现
type AddTag func(line int, content string) Tag

//遍历文件并打上标签
func (file *file) Scan(addTags ...AddTag) error {
	f, err := os.Open(file.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := Line{
			Content: scanner.Text(),
			Tags:    make(map[string]Empty),
		}
		file.Total++
		//打上标签
		for _, addTag := range addTags {
			if tag := addTag(file.Total, scanner.Text()); tag != "" {
				if _, ok := line.Tags[scanner.Text()]; !ok {
					//该行打上标签
					line.Tags[tag] = Empty{}
					//记录所有打上的标签
					if _, ok := file.Tags[tag]; !ok {
						file.Tags[tag] = Empty{}
					}
				}
			}
		}
		file.Lines = append(file.Lines, &line)
	}
	return nil
}

//在每个特定标签处中插入一行或多行,标签为空则替换整个文件内容
func (file *file) Insert(tags []Tag, contents []string) error {
	if tmpFile, err := createTmpFile(); err != nil {
		return err
	} else {
		defer removeTmpFile(tmpFile)

		if f, err := os.OpenFile(file.Path, os.O_WRONLY|os.O_TRUNC, 0644); err != nil {
			return err
		} else {
			//替换文件
			if len(tags) == 0 {
				if _, err := f.WriteString(strings.Join(contents, "\n")); err != nil {
					f.Close()
					return err
				}
			} else {
				//遍历文件行
				for _, line := range file.Lines {
					tmpFile.WriteString(line.Content + "\n")
					for _, tag := range tags {
						if _, ok := line.Tags[tag]; ok {
							//写入临时文件
							if _, err := tmpFile.WriteString(strings.Join(contents, "\n") + "\n"); err != nil {
								f.Close()
								return err
							}
						}
					}
				}
				f.Close()
				return os.Rename(tmpFile.Name(), file.Path)
			}

		}
	}
	return nil
}

//在开始和结束标签中插入一行或多行,如果存在多对标签只会取最内部那一对
func (file *file) InsertBetween(start, end Tag, contents []string) error {

	if start == "" || end == "" {
		return errors.New(EMPTY_TAG)
	}
	if _, exists := file.Tags[start]; !exists {
		return errors.New(START_TAG_NOT_EXISTS)
	}
	if _, exists := file.Tags[end]; !exists {
		return errors.New(END_TAG_NOT_EXISTS)
	}

	if tmpFile, err := createTmpFile(); err != nil {
		return err
	} else {
		defer removeTmpFile(tmpFile)

		matchStartLine, matchEndLine := 0, 0
		var readLine = 0
		for i, line := range file.Lines {
			readLine = i
			if _, exists := line.Tags[start]; exists {
				matchStartLine = i + 1
			} else {
				//已经匹配到开始
				if _, exists := line.Tags[end]; exists && matchStartLine > 0 {
					matchEndLine = i + 1
					if _, err := tmpFile.WriteString(strings.Join(contents, "\n") + "\n"); err != nil {
						return err
					}
					break
				}
			}
			tmpFile.WriteString(line.Content + "\n")
		}
		if matchEndLine == 0 {
			removeTmpFile(tmpFile)
			return errors.New(TAG_MISS_MATCH)
		}
		//未结束
		if readLine <= len(file.Lines)-1 {
			for _, line := range file.Lines[readLine:] {
				if _, err := tmpFile.WriteString(line.Content + "\n"); err != nil {
					return err
				}
			}
		}
		return os.Rename(tmpFile.Name(), file.Path)
	}
}

//在开始和结束标签插入一行或多行并且每一行内容不重复，
func (file *file) InsertBetweenNoRepeat(start, end Tag, contents []string) error {
	return nil
}

/*
	删除指定标签行,返回删除行数
	标签为空删除整个文件
*/
func (file *file) Delete([]Tag) (int, error) {

	return 0, nil
}

/*
	删除开始和结束标签之间的内容,返回删除行数
	开始标签和结束标签都为空删除整个文件
	结束标签为空,从开始标签开始删除到文件结束，包括开始标签所在行
	开始标签为空,从文件开始删除到结束标签,包括结束标签所在空行
*/
func (file *file) DeleteBetween(start, end Tag) (int, error) {
	return 0, nil
}

//临时文件
func createTmpFile(path ...string) (*os.File, error) {
	rand.Seed(time.Now().UnixNano())
	tmpFileName := fmt.Sprintf("/tmp/tmp_%s_%d", time.Now().String(), rand.Intn(10))
	tmpFileName = append(path, tmpFileName)[0]
	return os.OpenFile(tmpFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
}

//删除临时文件
func removeTmpFile(f *os.File) {
	f.Close()
	os.Remove(f.Name())
}
