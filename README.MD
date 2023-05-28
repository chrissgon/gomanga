
# gomanga

✨ The project aims to make manga reading easier, allowing you to view chapters directly from popular sources. Enjoy the convenience of accessing your favorite manga in one place! 📚💻

## 📦 Install

```bash
go get github.com/chrissgon/gomanga
```

## 🚀 Quick Start

Search by all providers:

```golang
import "github.com/chrissgon/gomanga"

gomanga.SearchByProviders("naruto", "1")
```

Search by specific provider:

```golang
import "github.com/chrissgon/gomanga"

gomanga.SearchByProvider("naruto", "1", MUITOMANGA)
```

### ✅ Providers

| OK | SITE                                  | NAME       | LANGUAGE
| -- | ------------------------------------- | ---------- | --------
| ✅ | [MangaLivre](https://mangalivre.net/) | MANGALIVRE | PT-BR
| ✅ | [MuitoManga](https://muitomanga.com/) | MUITOMANGA | PT-BR

## 🔒 Environment Variables

When you run this project, you can set the following environment variables.

`GOMANGA_TIMEOUT`: Search timeout in seconds.


## 💪🏻 Contribution

This project is open source and welcomes community contributions. Feel free to fork, implement improvements, and submit a pull request. Every contribution is valued and appreciated!

We hope that gomanga proves useful to you and enhances your manga reading experience. Feel free to explore the source code, provide feedback, and report any issues you encounter.

## 💚 Authors

- [@chrissgon](https://www.github.com/chrissgon)