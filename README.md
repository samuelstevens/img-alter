# ImgTag Labeler

Uses MS Azure to automatically caption all `<img\>` tags in `html` files and save to the `alt` attribute.

## Installation

```bash
make install # installs to ~/go/bin
```

## Usage

```bash
# shows options.
gocaption --help 

# edit and add "endpoint" and "key" variables.
vim ~/.labelrc.json 

# outputs a label for selfie.png
gocaption selfie.png

# add alt captions to all images in html files in website-dir.
gocaption --filetypes html --write --silent ~/projects/website-dir/
```

## Future Features
* Making requests concurrently
