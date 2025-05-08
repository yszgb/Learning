package models

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type PathConfig struct {
	Cert         string
	Key          string
	Avatars      string
	CourseImages string `yaml:"courseimages"`
	Lectures     string
	Html         string
	// Home 为传入变量，不从数据库获取
	Home string
}

type DatabaseConfig struct {
	Host         string
	Port         int
	Dbname       string
	User         string
	Passwd       string
	MaxOpenConns int    `yaml:"maxopenconns"`
	MaxIdleConns int    `yaml:"maxidleconns"`
	MaxIdleTime  string `yaml:"maxidletime"`
}

// 生成数据库连接字符串
func (c *DatabaseConfig) Dsn() string {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", c.User, c.Passwd, c.Host, c.Port, c.Dbname)
	return dsn
}

func (c *DatabaseConfig) MaxIdleDuration() time.Duration {
	duration, err := time.ParseDuration(c.MaxIdleTime)
	if err != nil {
		duration, _ = time.ParseDuration("15m")
	}
	return duration
}

type LimiterConfig struct {
	Rps     int
	Burst   int
	Enabled bool
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Passwd   string
	Sender   string
}

// 总体配置
type Config struct {
	Server    string
	Port      int
	VideoPort int `yaml:"videoport"`
	Env       string
	Path      PathConfig     `yaml:"path"`
	Database  DatabaseConfig `yaml:"database"`
	Limiter   LimiterConfig  `yaml:"limiter"`
	Smtp      SMTPConfig     `yaml:"smtp"`
}

// 初始化配置。读取、解析配置文件，返回配置对象。
func NewConfig(path, home string) (*Config, error) {
	stats, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	bytes := make([]byte, stats.Size()) // 缓存
	buf := bufio.NewReader(f)
	_, err = buf.Read(bytes)
	if err != nil {
		return nil, err
	}

	cfg := Config{}
	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	} else {
		return &cfg, nil
	}
}

// 辅助函数

// 获取证书路径
func (c *PathConfig) CertPath() string {
	return filepath.Join(c.Home, c.Cert)
}

// 获取密钥路径
func (c *PathConfig) KeyPath() string {
	return filepath.Join(c.Home, c.Key)
}

// 获取头像路径
func (c *PathConfig) AvatarsPath() string {
	return filepath.Join(c.Home, c.Avatars)
}

// 获取课程图片路径
func (c *PathConfig) CourseImagePath() string {
	return filepath.Join(c.Home, c.CourseImages)
}

// 获取视频路径
func (c *PathConfig) LecturePath() string {
	return filepath.Join(c.Home, c.Lectures)
}

// 获取HTML路径
func (c *PathConfig) HtmlPath() string {
	return filepath.Join(c.Home, c.Html)
}
