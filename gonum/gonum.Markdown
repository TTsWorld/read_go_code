

gonum 库是一个 go语言 BLAS 的实现库。

## BLAS 的介绍
BLAS(Basic Linear Algebra Subprograms)是一组线性代数计算中通用的基本运算操作函数集合[1] 。BLAS Technical (BLAST) Forum负责规范BLAS函数接口, 并在网站[1]公布一个由Fortran语言编写的BLAS库。
这个Fortran版BLAS库通常被称为BLAS参考库(the reference implementation)。 BLAS参考库使用的算法能高效地给出正确的结果，但仍有许多优化潜力。要想获得更高的计算效率，可以使用优化的BLAS库。


BLAS 函数

BLAS库中函数根据运算对象的不同，分为三个类：
Level 1 函数处理单一向量的线性运算以及两个向量的二元运算。Level 1 函数最初出现在1979年公布的BLAS库中。
Level 2 函数处理 矩阵与向量的运算，同时也包含线性方程求解计算。 Level 2 函数公布于1988年。
Level 3 函数包含矩阵与矩阵运算。Level 3 函数发表于1990年。


### BLAS 函数接口命名规范

BLAS Level 1 接口函数名称由前缀+操作简称组成。
例如 SROTG，其中
S        -- 标明矩阵或向量中元素数据类型的前缀；
ROTG -- 向量运算简称.

前缀： 矩阵或向量内元素的数据类型，有以下几种：
S - 单精度浮点数
D - 双精度浮点数
C - 复数
Z - 16位复数


Level 2 和 Level 3函数涉及矩阵运算，接口函数名称由前缀 + 矩阵类型 + 操作简称组成。
例如: SGEMV
S     -- 标明矩阵或向量中元素数据类型的前缀；
GE   -- 矩阵类型
MV  -- 向量或矩阵运算简称

BLAS库中使用的矩阵类型有以下几种：
GE - GEnearl  稠密矩阵
GB - General Band 带状矩阵

    SY - SYmmetric    对称矩阵
    SB - Symmetric Band 对称带状矩阵
    SP - Symmetric Packed  压缩存储对称矩阵

    HE - HEmmitian     Hemmitian矩阵，自共轭矩阵
    HB - Hemmitian Band   带状Hemmitian矩阵
    HP - Hemmitian Packed  压缩存储Hemmitian矩阵

    TR - TRiangular      三角矩阵
    TB - Triangular Band  三角带状矩阵 
    TP - Triangular Packed  压缩存储三角矩阵



参考
[BLAS 整体概况](https://www.cnblogs.com/chen1846847163/p/15592828.html)
[gonum doc](https://pkg.go.dev/gonum.org/v1/gonum@v0.11.0#section-readme)
[矩阵运算库blas, cblas, openblas, atlas, lapack, mkl之间有什么关系，在性能上区别大吗？](https://www.zhihu.com/question/27872849)