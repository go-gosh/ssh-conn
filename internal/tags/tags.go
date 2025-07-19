package tags

import (
	"os"
	"path/filepath"
	"sync"

	"io"

	"gopkg.in/yaml.v3"
)

var (
	tagsFileOnce sync.Once
	tagsFilePath string
)

func getHomeDir() string {
	h := os.Getenv("HOME")
	if h != "" {
		return h
	}
	home, _ := os.UserHomeDir()
	return home
}

func getTagsFilePath() string {
	tagsFileOnce.Do(func() {
		home := getHomeDir()
		dir := filepath.Join(home, ".ssh-conn")
		os.MkdirAll(dir, 0700)
		tagsFilePath = filepath.Join(dir, "tags.yaml")
	})
	return tagsFilePath
}

type TagsMap map[string][]string // tag -> []host

// 读取 tags.yaml
func loadTags() (TagsMap, error) {
	path := getTagsFilePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return TagsMap{}, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var m TagsMap
	err = yaml.NewDecoder(f).Decode(&m)
	if err != nil {
		if err == io.EOF {
			return TagsMap{}, nil
		}
		return nil, err
	}
	return m, nil
}

// 保存 tags.yaml
func saveTags(m TagsMap) error {
	path := getTagsFilePath()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := yaml.NewEncoder(f)
	defer enc.Close()
	return enc.Encode(m)
}

// AddTags 给主机添加标签
func AddTags(host string, tags []string) error {
	m, err := loadTags()
	if err != nil {
		return err
	}
	for _, tag := range tags {
		found := false
		for _, h := range m[tag] {
			if h == host {
				found = true
				break
			}
		}
		if !found {
			m[tag] = append(m[tag], host)
		}
	}
	return saveTags(m)
}

// RemoveTags 从主机移除标签
func RemoveTags(host string, tags []string) error {
	m, err := loadTags()
	if err != nil {
		return err
	}
	for _, tag := range tags {
		hosts := m[tag]
		var newHosts []string
		for _, h := range hosts {
			if h != host {
				newHosts = append(newHosts, h)
			}
		}
		if len(newHosts) == 0 {
			delete(m, tag)
		} else {
			m[tag] = newHosts
		}
	}
	return saveTags(m)
}

// SetTags 覆盖主机标签
func SetTags(host string, tags []string) error {
	m, err := loadTags()
	if err != nil {
		return err
	}
	// 先移除所有包含 host 的标签
	for tag, hosts := range m {
		var newHosts []string
		for _, h := range hosts {
			if h != host {
				newHosts = append(newHosts, h)
			}
		}
		if len(newHosts) == 0 {
			delete(m, tag)
		} else {
			m[tag] = newHosts
		}
	}
	// 再添加新标签
	for _, tag := range tags {
		m[tag] = append(m[tag], host)
	}
	return saveTags(m)
}

// GetTags 查询主机所有标签
func GetTags(host string) ([]string, error) {
	m, err := loadTags()
	if err != nil {
		return nil, err
	}
	var tags []string
	for tag, hosts := range m {
		for _, h := range hosts {
			if h == host {
				tags = append(tags, tag)
				break
			}
		}
	}
	return tags, nil
}

// AllTags 返回所有标签
func AllTags() ([]string, error) {
	m, err := loadTags()
	if err != nil {
		return nil, err
	}
	var tags []string
	for tag := range m {
		tags = append(tags, tag)
	}
	return tags, nil
}

// HostsByTags 支持 AND/OR 查询
func HostsByTags(tags []string, orMode bool) ([]string, error) {
	m, err := loadTags()
	if err != nil {
		return nil, err
	}
	hostSet := map[string]int{} // host -> 命中标签数
	for _, tag := range tags {
		for _, h := range m[tag] {
			hostSet[h]++
		}
	}
	var hosts []string
	for h, cnt := range hostSet {
		if orMode {
			hosts = append(hosts, h)
		} else {
			if cnt == len(tags) {
				hosts = append(hosts, h)
			}
		}
	}
	return hosts, nil
}

// 批量添加标签
func AddTagsBatch(hosts, tagsList []string) error {
	for _, host := range hosts {
		err := AddTags(host, tagsList)
		if err != nil {
			return err
		}
	}
	return nil
}

// 批量移除标签
func RemoveTagsBatch(hosts, tagsList []string) error {
	for _, host := range hosts {
		err := RemoveTags(host, tagsList)
		if err != nil {
			return err
		}
	}
	return nil
}

// 获取全部标签映射
func LoadAll() (TagsMap, error) {
	return loadTags()
}

// 如果 tags.yaml 不存在，自动生成空模板
func EnsureTagsFile() error {
	path := getTagsFilePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return saveTags(TagsMap{})
	}
	return nil
}
