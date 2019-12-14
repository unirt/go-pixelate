# go-tile
**Simple image pixelater**  
You can easily convert an image to *tile*.
## Usage
1. Download `go-tile` binary from [releases](https://github.com/unirt/go-tile/releases)  
   Recommended version: **v0.2**
2. Convert image to *tile*
   ```bash
   ./go-tile -i path/to/image -n 4
   # converted image path will be `./outputs/16tile_{originalImagePath}.png`
   ```
#### if you want to `go run`
1. Clone this repository
   ```bash
   git clone git@github.com:unirt/go-tile.git
   ```
2. Convert image to *tile*
   ```bash
   cd go-tile
   go run main.go -i path/to/image -n 4
   # converted image path will be `./outputs/16tile_{originalImagePath}.png`
   ```
### Args
|arg  |must/optional  |description  |
|:-:|:-:|:-:|
|i  |must  |Path to original image  |
|n  |optional(default `4`)|The square root of the number of tiles in converted image. <br> ex. Converted image will have 16 tiles if you provide `-n 4`|
## Example
```bash
./go-tile -i ~/Pictures/sample.png
```
Original(`sample.png`)  
![sample.png](.github/sample.jpg)  
Result(`16tile_sample.png`)  
![16tile_sample.png](.github/16tile_sample.png)
