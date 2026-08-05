package main

import (
  "regexp"
)

type Scanner struct {
  rules map[string]*regexp.Regexp
}

type Match struct {
  Rule  string
  File  string
  Value string
}

func NewScanner() *Scanner {
  return &Scanner {
    rules: make(map[string]*regexp.Regexp),
  }
}

func (s *Scanner) AddRegex(name, pattern string) error {
  r, err := regexp.Compile(pattern)
  if err != nil {
    return err
  }
  s.rules[name] = r
  return nil
}

func (s *Scanner) LoadDefaultRegexes() error {
	rules := map[string]string{
	  "email":               `\b[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}\b`,
	  "url_absolute":        `https?://[^\s"'<>]+`,
	  "url_root_relative":   `/(?:[A-Za-z0-9._~!$&'()*+,;=:@%-]+/?)+`,
	  "url_dot_relative":    `(?:\./|\.\./)[^\s"'<>]+`,
	  "websocket_url":       `wss?://[^\s"'<>]+`,
	  "graphql_endpoint":    `/graphql(?:/)?(?:\?[^\s"'<>]*)?`,
	  "fetch_call":          `fetch\s*\(\s*["']([^"']+)["']`,
	  "axios_call":          `axios\.(?:get|post|put|delete|patch|request)\s*\(\s*["']([^"']+)["']`,
	  "xhr_call":            `open\s*\(\s*["'](GET|POST|PUT|DELETE|PATCH)["']\s*,\s*["']([^"']+)["']`,
	  "eventsource_call":    `new\s+EventSource\s*\(\s*["']([^"']+)["']`,
	  "api_key_assignment":  `(?i)\b[a-z0-9_.-]*key[a-z0-9_.-]*\b\s*[:=]\s*["'][^"']+["']`,
	  "secret_assignment":   `(?i)\b[a-z0-9_.-]*secret[a-z0-9_.-]*\b\s*[:=]\s*["'][^"']+["']`,
	  "token_assignment":    `(?i)\b[a-z0-9_.-]*token[a-z0-9_.-]*\b\s*[:=]\s*["'][^"']+["']`,
	  "password_assignment": `(?i)\b[a-z0-9_.-]*password[a-z0-9_.-]*\b\s*[:=]\s*["'][^"']+["']`,
	  "aws_access_key":      `AKIA[0-9A-Z]{16}`,
	  "google_api_key":      `AIza[0-9A-Za-z\-_]{35}`,
	  "jwt":                 `eyJ[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+`,
	  "github_pat":          `ghp_[A-Za-z0-9]{36}`,
	  "slack_token":         `xox[baprs]-[A-Za-z0-9-]+`,
	  "stripe_live_key":     `(?:pk|sk)_live_[A-Za-z0-9]{24,}`,
	  "firebase_url":        `https://[A-Za-z0-9-]+\.firebaseio\.com`,
	  "sentry_dsn":          `https://[A-Za-z0-9]+@[A-Za-z0-9.-]+/\d+`,
  }
	for name, pattern := range rules {
		if err := s.AddRegex(name, pattern); err != nil {
			return err
		}
	}
	return nil
}

func (s *Scanner) Scan(filename string, data []byte) []Match {
  matchList := []Match{}
  for name, r := range s.rules {
    matches := r.FindAll(data, -1)
    for _, match := range matches {
      matchList = append(matchList, Match{
        Rule: name,
        File: filename,
        Value: string(match),
      })
    }
  }
  return matchList
}

