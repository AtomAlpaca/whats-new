# whats-new

一个用于获取和分析网页内容的 CLI 工具。

## 安装

```bash
cd src
go build -o ../whats-new .
```

## 命令

### links - 获取网页所有链接

获取网页中的所有链接（自动去重）。

```bash
whats-new links "https://example.com"
```

输出：
```json
{
  "count": 1,
  "links": [
    "https://iana.org/domains/example"
  ],
  "url": "https://example.com"
}
```

### extract - 提取正文内容

提取网页的主要内容，输出为纯文本。

```bash
whats-new extract "https://golang.org/doc"
```

输出：
```json
{
  "content": "The Go programming language is an open source project to make programmers more productive.\n\nGo is expressive, concise, clean, and efficient...",
  "url": "https://golang.org/doc"
}
```

特性：
- 自动移除脚本、样式表、导航栏等无关内容
- 保留段落换行
- 过滤短文本噪音（导航链接、按钮文字等）
- 标题单独处理

### full-text - 获取完整 HTML

获取网页的完整 HTML 内容。

```bash
whats-new full-text "https://example.com"
```

输出：
```json
{
  "content": "<!DOCTYPE html><html lang=\"en\">...",
  "url": "https://example.com"
}
```

注意：JSON 会自动转义 HTML 特殊字符。

### metadata - 提取网页元数据

提取网页的元数据信息。

```bash
whats-new metadata "https://golang.org"
```

输出：
```json
{
  "description": "Go is an open source programming language that makes it simple to build secure, scalable systems.",
  "favicon": "https://golang.org/images/favicon-gopher.png",
  "title": "The Go Programming Language",
  "url": "https://golang.org"
}
```

支持的元数据：
- 基本信息：`title`, `description`, `keywords`, `author`
- Open Graph：`og:title`, `og:description`, `og:image`, `og:url`, `og:site_name`
- 其他：`favicon`, `url`

### memory - 网站记忆

将网站信息保存到本地 SQLite 数据库，支持自动补全信息。

#### 保存网站

```bash
whats-new memory save '{"url": "https://example.com"}'
```

手动指定信息：
```bash
whats-new memory save '{"url": "https://golang.org", "title": "Go Language", "report": "The Go programming language"}'
```

自动补全：
- 未提供 `title` → 自动从网页元数据获取
- 未提供 `report` → 自动从网页内容提取前 3 行

输出：
```json
{
  "content_hash": "9863ff025eeb17fe6b32072e32c711e2350a0b7b600cef4ff8d209e6324dafb1",
  "last_recorded": "2026-03-14T19:46:44.568584732+08:00",
  "saved": true,
  "title": "Example Domain",
  "url": "https://example.com"
}
```

#### 获取网站

```bash
whats-new memory get "https://example.com"
```

输出：
```json
{
  "content_hash": "9863ff025eeb17fe6b32072e32c711e2350a0b7b600cef4ff8d209e6324dafb1",
  "content_report": "<!DOCTYPE html><html lang=\"en\">...",
  "created_at": "2026-03-14T11:39:29Z",
  "id": 1,
  "last_recorded": "2026-03-14T19:46:44.568584732+08:00",
  "title": "Example Domain",
  "updated_at": "2026-03-14T19:46:44.568584732+08:00",
  "url": "https://example.com"
}
```

## 数据库

数据存储在 SQLite 数据库中（`~/.whats-new.db`）。

表结构：
```sql
CREATE TABLE websites (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    url TEXT UNIQUE NOT NULL,
    title TEXT,
    content_hash TEXT NOT NULL,
    last_recorded DATETIME NOT NULL,
    content_report TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
)
```

字段说明：
- `url` - 网站链接
- `title` - 网站标题
- `content_hash` - HTML 内容全文的 SHA256 哈希值
- `last_recorded` - 最后记录时间
- `content_report` - 网站内容摘要
- `created_at` - 首次记录时间
- `updated_at` - 最后更新时间

## 技术栈

- **Go** - 编程语言
- **cobra** - CLI 框架
- **goquery** - HTML 解析库
- **go-sqlite3** - SQLite 数据库驱动
- **golang.org/x/crypto** - 哈希计算
