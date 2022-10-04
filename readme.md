## Mandelbrot Maker
----

### Basic Algorithm
The simplest algorithm for generating a representation of the Mandelbrot set is known as the `"escape time"` algorithm. A repeating calculation is performed for each x, y point in the plot area and based on the behavior of that calculation, a color is chosen for that pixel.

- For each (x, y), the locations are used as starting values in 
a repeating calculation.
- The result of each iteration is used as the starting value for the next iteration.
- The values are checked during each iteration to see whether they have reached a `critical escape condition` or `bailout`.
- If the value has not reached the critical escape condition, the iteration continues. Any pixel in the set will never escape the critical condition, hence there is a maximum number of iterations that can be performed.
---
### In terms of Complex Numbers
- $z_{n+1}=z_{n}^2+z_0$
    - where  $z_0 = x_0 + y_0i$
- $z_{n+1}=(x_n + y_ni)^2 + (x_0 + y_0i)$
- $z_{n+1}=(x_n^2 - y_n^2 + x_0) + (2x_ny_n + y_0)i$
- $x_{n+1}=x_n^2 - y_n^2 + x_0$
- $y_{n+1}=2x_ny_n + y_0$
---
### Psuedo Code
1. Unoptimized naive escape time algorithm
```bash
LIMIT = 4

for each pixel (Px, Py) do
    x0 := scaled(Px) # scaled x-coord of pixel (scaled to lie in the Mandelbrot X scale (mx, Mx))
    y0 := scaled(Py) # scaled y-coord of pixel (scaled to lie in the Mandelbrot Y scale (my, My))

    x := 0.0
    y := 0.0
    iteration := 0
    max_iteration := MAX_ITERATIONS
    while (x*x + y*y) <= LIMIT AND iteration < MAX_ITERATIONS do  # 2 mul
        x_temp := x*x - y*y + x0 # 2 mul
        y := (x+x)*y + y0  # = 2*x*y + y0 (to reduce one multiplication) # 1 mul
        x := x_temp
        iteration := iteration + 1
    

    if iteration < MAX_ITERATIONS then
        color := palette[iteration]
    else
        color := palette[0]  # generally black

    plot(Px, Py, color)

```
2. OPTIMIZED escape time algorithm
Above code uses an unoptimized while operation. Aboev one must perform 5 multiplications per iteration.

Following simplication can be done by simplifying the complex multiplication.

- $(iy+x)^2 = x^2 - y^2 + 2iyx$


```bash
for each pixel do
    ...
    x2 := 0
    y2 := 0

    while x2 + y2 <= LIMIT and iteration < MAX_ITERATIONS do
        x := x2 - y2 + x0
        y := (x2 +x2)*y2  + y0  # 1 mul
        x2 := x*x  # 1 mul
        y2 := y*y  # 1 mul
        ...
```


### Usage
```bash
Flags: 
  -width int
        Width of the image (default 700)
  -height int
        Height of the image (default 400)
  -scale int
        Scale of the image
  -iter int
        Max Iterations (default 1000)
  -mode string
        Mode of the image (default "seq") (options: seq, pixel, row, workers)
  -out string
        Name of the output file with extension (default "mandelbrot.png")
  -threshold float
        Threshold for the mandelbrot set (default 4)
  -workers int
        Number of workers to use (default 4)
```

### Examples
```bash
go build -o ./build/mandel ./main.go && ./build/mandel \
    --out img/mandelbrot.png \
    --mode pixel \
    --scale 1 \
    --threshold 128 \
    --workers 8 \
    --iter 1000
```
`1x`
- `png`
- `406 KiloBytes`
- Using mode: pixel (one goroutine per pixel)
- Time taken to create image: 144.6938ms
- Time taken to save image: 80.7434ms
![mandelbrot 1x](./img/mandelbrot.png)

`10x`
- `png`
- `20 MegaBytes`
- Using mode: pixel (one goroutine per pixel)
- Time taken to create image: 14.8648101s
- Time taken to save image: 9.5717677s
![mandelbrot 10x](./img/mandelbrotx10.png)

`10x`
- `jpg`
- `6.1 MegaBytes`
- Using mode: pixel (one goroutine per pixel)
- Time taken to create image: 14.6557209s
- Time taken to save image: 2.0148509s
![mandelbrot 10x](./img/mandelbrotx10.jpg)



### Gallery
![mandelbrotx10](./img/mandelbrotx10.jpg)
![web-colored#i=120_t=1000_z=1e+06_x=0.243_y=0.8115](./img/web-colored#i=120_t=1000_z=1e+06_x=0.243_y=0.8115.jpg)
![web-colored#i=120_t=1000_z=20_x=0.2_y=0.8](./img/web-colored#i=120_t=1000_z=20_x=0.2_y=0.8.jpg)
![web-colored#i=120_t=1000_z=20_x=0.25_y=0.8](./img/web-colored#i=120_t=1000_z=20_x=0.25_y=0.8.jpg)
![web-colored#i=120_t=1000_z=100_x=0.245_y=0.81](./img/web-colored#i=120_t=1000_z=100_x=0.245_y=0.81.jpg)
![web-colored#i=120_t=1000_z=1000_x=0.243_y=0.8111](./img/web-colored#i=120_t=1000_z=1000_x=0.243_y=0.8111.jpg)
![web-colored#i=120_t=1000_z=10000_x=0.243_y=0.8113](./img/web-colored#i=120_t=1000_z=10000_x=0.243_y=0.8113.jpg)
![web-colored#i=120_t=1000_z=10000_x=0.243_y=0.8114](./img/web-colored#i=120_t=1000_z=10000_x=0.243_y=0.8114.jpg)
![web-colored#i=120_t=1000_z=10000_x=0.243_y=0.8115](./img/web-colored#i=120_t=1000_z=10000_x=0.243_y=0.8115.jpg)
![web-colored#i=120_t=1000_z=100000_x=0.243_y=0.8115](./img/web-colored#i=120_t=1000_z=100000_x=0.243_y=0.8115.jpg)
![web-colored#i=400_t=1000_z=2e+07_x=0.7_y=0.291](./img/web-colored#i=400_t=1000_z=2e+07_x=0.7_y=0.291.jpg)
![web-colored#i=500_t=1000_z=2e+07_x=0.7_y=0.291](./img/web-colored#i=500_t=1000_z=2e+07_x=0.7_y=0.291.jpg)
![web-colored#i=800_t=1000_z=2e+07_x=0.7_y=0.291](./img/web-colored#i=800_t=1000_z=2e+07_x=0.7_y=0.291.jpg)
![colored#i=100_t=1000_z=1_x=0.5_y=-0.1313](./img/colored#i=100_t=1000_z=1_x=0.5_y=-0.1313.jpg)
![colored#i=100_t=1000_z=100_x=0.7_y=0.29](./img/colored#i=100_t=1000_z=100_x=0.7_y=0.29.jpg)
![colored#i=100_t=1000_z=100_x=0.7_y=-0.29](./img/colored#i=100_t=1000_z=100_x=0.7_y=-0.29.jpg)
![colored#i=800_t=1000_z=2e+07_x=0.7_y=0.291](./img/colored#i=800_t=1000_z=2e+07_x=0.7_y=0.291.jpg)
![colored#i=1000_t=1000_z=1000_h=200_x=0.7435_y=-0.1315](./img/colored#i=1000_t=1000_z=1000_h=200_x=0.7435_y=-0.1315.jpg)
![colored](./img/colored.jpg)
![mandelbrot](./img/mandelbrot.png)