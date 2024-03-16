pixelsorter (+ cli) written in golang  
used my (node)js implementation as reference, which uses [satyarth/pixelsort](https://github.com/satyarth/pixelsort/) as reference  
not at parity; no sorter/comparator validation... i must still be forgettin somethin  
would like to mix and match thresholds (sort by red with saturation threshold, or sort by sat with darkness threshold)
first time usin go, this is prolly super cursed  
it go nyoom tho :3

# example
![webb picture-of-the-month 2023/10](https://cdn.esawebb.org/archives/images/screen/potm2310a.jpg) - [src: ESA/Webb, NASA & CSA, A. Adamo (Stockholm University) and the FEAST JWST team](https://esawebb.org/images/potm2310a/)  
![mask](./webb-mask.jpg)  
![sorted](./webb-sort.jpg)  
did you know webb and hubble pics are cc4?  

"elapsed" is time spent actually sorting  
```
=> time pixelsort_go -i ~/Downloads/potm2310a.jpg -m ./webb-mask.jpg -o ./benchmark
Sorting 1 images with a config of {Sorter:row Comparator:lightness Quality:0 SectionLength:60 Randomness:1 Reverse:false Thresholds:{Lower:0 Upper:1} Angle:0}.
Loading image 1 ("/home/______/Downloads/potm2310a.jpg" -> "./benchmark")...
./benchmark Elapsed: 1.03s
0:02.28 real 2.20 usermode 99% 214980 Kb mem
```
