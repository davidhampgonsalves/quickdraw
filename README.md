# Quick Draw Grid
Generate a grid of randomly selected [Google Quick Draw Drawings](https://github.com/googlecreativelab/quickdraw-dataset) probably for pen plotting.

![Airplane grid](./airplane-grid.svg)

Simplified quick draw data(which this tool will use to populate the grid) can be found at https://console.cloud.google.com/storage/browser/quickdraw_dataset/full/simplified.

## Install
```
go get github.com/davidhampgonsalves/quickdraw
```

## Usage
Generate a 1000x2000 px grid of drawings with 20 in each row and a 15px margin based on the provided simplified quick draw ndjson data
```
quickdraw -w 1000 -h 2000 -c 20 -m 15 "./full_simplified_airplane.ndjson"
```

Help `quickdraw --help`
