# Gitignore CLI

A simple CLI for generating gitignore files on the fly powered by [gitignore.io](https://gitignore.io).

## Install

```sh
go get github.com/fvumbaca/gitignore
```

### Autocompletion

#### Bash

```sh
echo "source <(gitignore --bash-autocomplete)" >> ~/.bashrc
# You might also need to refresh your source file
source ~/.bashrc
```

#### Zsh

```sh
echo "source <(gitignore --zsh-autocomplete)" >> ~/.zshrc
# You might also need to refresh your source file
source ~/.zshrc
```

## Usage

```sh
gitignore go visualstudiocode
```