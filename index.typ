#import "./tufted/tufted.typ" as tufted

#let template = tufted.tufted-web.with(
  header-links: (
    "/": "Home",
    "/blog/": "Blog",
    "/friends/" : "Friends",
    "/cv/": "CV",
  ),
  title: "Tufted",
)

#import "@preview/cmarker:0.1.7"

#show: template.with(title: "Blog")

錦瑟無端五十絃，一絃一柱思華年。

2023 年北京集训好题分享讲了这题。

== 题意
<题意>
给定一个长为 $n$ 的序列
$a_n \, forall i in [1 \, n] \, 1 lt.eq a_i lt.eq n$。

定义 $f (l \, r) = (min_(i = l)^r a_i \, max_(i = l)^r a_i)$。

$q$ 次询问，每次给定 $l \, r$，询问最小的 $k$ 使得
$f^k (l \, r) = (1 \, n)$，无解输出 $- 1$。

== 题解
<题解>
首先两个十分显然的性质：

- 如果某次操作把区间变成了 $(1 \, n)$，那么无论再操作多少次这个区间都是
  $(1 \, n)$；

- 状态数是 $O (n^2)$ 的。

这引导我们想到，如果能求得 $f^k (l \, r)$ 在 $k gt.eq n^2$
时的结果，就能判定是否有解，同时也可以利用二分之类的方法求得答案。

Key Observation 1:
$f (l \, r) = union.big_(i = l)^(r - 1) f (i \, i + 1)$。

证明：考虑归纳，则只需证明：$[l_1 \, r_1] union [l_2 \, r_2] = [l \, r] \, [l_1 \, r_1] sect [l_2 \, r_2] eq.not diameter$，则
$f (l \, r) = f (l_1 \, r_1) union f (l_2 \, r_2)$。而这是显然的。

Key Observation 2:
$f^k (l \, r) = union.big_(i = l)^(r - 1) f^k (i \, i + 1)$。

证明：考虑上一页中结论，每次增加 $k$ 相邻两个区间仍然总是有交。

$  & f^k (l \, r)\
= & f^k ([l_1 \, r_1] union [l_2 \, r_2])\
= & f (f^(k - 1) (l_1 \, r_1) union f^(k - 1) (l_2 \, r_2))\
= & f (f^(k - 1) (l_1 \, r_1)) union f (f^(k - 1) (l_2 \, r_2))\
= & f^k (l_1 \, r_1) union f^k (l_2 \, r_2)\
 $

然后我们发现到最后相邻两项区间还是有交，因此我们最终区间到左、右端点就是这些区间左、右节点的极值。这允许我们通过维护
$[i \, i + 1]$ 的信息，并利用 st 表求得任意区间的结果。

我们令 $F \/ G_(k \, j \, i)$ 为 $f^k (j \, j + 2^i - 1)$
的左、右端点。那么对于 $i$ 这维的转移，我们有：

$ F_(k \, j \, i) & = min (F_(k \, j \, i - 1) \, F_(k \, j + 2^(i - 1) \, i - 1))\
G_(k \, j \, i) & = max (G_(k \, j \, i - 1) \, G_(k \, j + 2^(i - 1) \, i - 1)) $

对于 $k$ 这一维，我们有：

$  & f^k (l \, r) = f (f^(k - 1) (l_1 \, r_1) union f^(k - 1) (l_2 \, r_2))\
 $

那么我们知道 $f^k (l \, r)$ 的左右端点分别是：

$  & min (F_(k \, l \, l g) \, F_(k \, r - 2^(l g) + 1 \, l g))\
 & max (G_(k \, l \, l g) \, G_(k \, r - 2^(l g) + 1 \, l g)) $

至此，预处理后我们能够在 $O (1)$ 时间内求解
$f^k (l \, r)$。二分或倍增即可求得答案。

== 代码
<代码>
代码实现把 $F$ 和 $G$ 放在了同一个数组里来卡常。

```
#pragma GCC optimize("Ofast")

#include <bits/stdc++.h>

const int MAX = 1e5 + 5;
const int LG = 35;
const int MAXX = 37;

int n, q, l, r;
int a[MAX], lg2[MAX], f[MAXX][MAX][20][3];

inline int read()
{
    char c=getchar();int x=0;bool f=0;
    for(;!isdigit(c);c=getchar())f^=!(c^45);
    for(;isdigit(c);c=getchar())x=(x<<1)+(x<<3)+(c^48);
    if(f)x=-x;return x;
}

inline int min(int a, int b) { return a < b ? a : b; }
inline int max(int a, int b) { return a > b ? a : b; }
void init(int k)
{
    for (int i = 1; (1 << i) < n; ++i)
    {
        for (int j = 1; j + (1 << i) <= n; ++j)
        {
            f[k][j][i][0] = min(f[k][j][i - 1][0], f[k][j + (1 << (i - 1))][i - 1][0]);
            f[k][j][i][1] = max(f[k][j][i - 1][1], f[k][j + (1 << (i - 1))][i - 1][1]);
        }
    }
}

int getl(int l, int r, int k)
{
    int lg = lg2[r - l + 1];
    return min(f[k][l][lg][0], f[k][r - (1 << lg) + 1][lg][0]);
}

int getr(int l, int r, int k)
{
    int lg = lg2[r - l + 1];
    return max(f[k][l][lg][1], f[k][r - (1 << lg) + 1][lg][1]);
}

void solve()
{
    l = read(); r = read(); long long res = 0;
    if (l == 1 and r == n) { printf("0\n"); return ; }
    else if (l == r) { printf("-1\n"); return ; }
    int _l = getl(l, r - 1, LG), _r = getr(l, r - 1, LG);
    if (_l != 1 or _r != n) { printf("-1\n"); return ; }
    for (int i = LG - 1; i >= 0; --i)
    {
        _l = getl(l, r - 1, i), _r = getr(l, r - 1, i);
        if (_l != 1 or _r != n) { res += (1ll << i); l = _l; r = _r; }
    }
    _l = getl(l, r - 1, 0), _r = getr(l, r - 1, 0);
    if (_l == 1 and _r == n) { printf("%lld\n", res + 1); } else { printf("-1\n"); }
}

int main()
{
    n = read(); q = read();
    for (int i = 1; i <= n; ++i) { a[i] = read(); }
    for (int i = 2; i <= n; ++i) { lg2[i] = lg2[i >> 1] + 1; }
    for (int i = 1; i <  n; ++i) { f[0][i][0][0] = min(a[i], a[i + 1]); f[0][i][0][1] = max(a[i], a[i + 1]); }
    init(0);
    for (int i = 1; i <= LG; ++i)
    {
        for (int j = 1; j <  n; ++j)
        {
            f[i][j][0][0] = getl(f[i - 1][j][0][0], f[i - 1][j][0][1] - 1, i - 1);
            f[i][j][0][1] = getr(f[i - 1][j][0][0], f[i - 1][j][0][1] - 1, i - 1);
        }
        init(i);
    }
    while (q--) { solve(); }
    return 0;
}
```
