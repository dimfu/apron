# apron

General purpose language to describe recipes of any kind. Heavily inspired from [https://github.com/reciperium/recipe-lang](recipe-lang).

## Installation

```bash
go install github.com/dimfu/apron@latest
```

## Features

apron's supported syntaxes:

- Ingredients with the tag `{ingredient_name}` or with amount: `{ingredient_name}(200gr)`
- Materials: `&{pot}`
- Timers: `t{15 minutes}`
- Metadata: with `>> tags: abc, easy, high-fiber`
- Comments: with `/* my comment */` or `// comment`

## Usage

To run apron, you just need to provide the file path location.
For recipe content example you can see [exmaple/jean-claude.apr](example/jean-claude.apr)

```bash
apron $PATH_TO_RECIPE_FILE
```
