This program converts a mbtiles file to a customized map format

1. remove all index

2. rename tiles table to map

3. add type column to map table

4. add index to map table

5. add task table to be compatible with legacy webui


how to run it?

./mbtilesConverter -i input.mbtiles -t type

type represents the type of tiles in the mbtiles file, it can be one of the following values:

0 - vector tiles
1 - raster tiles