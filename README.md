# GO--SASS
A simple SASS compiler written in go.

## Usage
Download the GO-SASS.exe from the releases page then to use write out a scss file like the example below:
```scss
$primary-color: red;
$seccondary-color: green;
$font: monospace;

// Nested example
nav {
    color: $primary-color;
    font: $font;
    a {
        color: $seccondary-color;
    }
}
```
Save that file then execute the program with the `-i` flag meaning include:
```bash
.\GO-SASS.exe -i text.scss
```
Then open the file output.css and find that tyour .scss file has been compiled into a valid .css file Like this example:
```css
nav {
    color: red;
    font-family: monospace;
}

nav a {
    color: green;
}
```

## Suppirted syntax features

| Feature name | Example |
| --- | --- |
| Nesting | Multiple selectors nested for styling specific elements in other elements |
| Parent selecting | Any nested element with & gets added say `&:hover` in `nav { a {}}` becomes nav a:hover{} |
| Variable support | Take any scss variable and turn it into css |
| Comment support | Take scss comments // and turn them into css commands /* */ |