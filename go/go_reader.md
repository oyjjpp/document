# golang

## 环境安装

- [下载包](https://www.golangtc.com/download)

## 问题

- 切片、数组差别
- 指针、引用
- 字符串 byte长度 定义
- 执行流程

## 流行框架

- [beego](https://beego.me/)
- [Buffalo](https://github.com/gobuffalo/buffalo)
- [Echo](https://github.com/labstack/echo)
- [Gin](https://github.com/gin-gonic/gin)
- [Iris](https://github.com/kataras/iris)
- [Revel](https://github.com/revel/revel)
- [go实现的定时器](https://github.com/robfig/cron)

## 相关书籍

- [Go语言实战: 编写可维护Go语言代码建议](https://github.com/llitfkitfk/go-best-practice)
- [雨痕 go学习笔记](https://github.com/qyuhen/book)
- [Go 入门指南-github](https://github.com/Unknwon/the-way-to-go_ZH_CN)
- [Go 入门指南-kancloud](https://www.kancloud.cn/kancloud/the-way-to-go/72432)

## 常见包

### io相关

#### io

### 网络相关

#### net

#### sync

### 数据库相关

#### database/sql

#### database/sql/driver

### 编码相关

#### json

### 字符串相关

#### fmt

##### 1、简介

包实现了类似C语言printf和scanf的格式化I/O

##### 2、标准输出

>// Print 将参数列表 a 中的各个参数转换为字符串并写入到标准输出中。  
>// 非字符串参数之间会添加空格，返回写入的字节数。  
>func Print(a ...interface{}) (n int, err error)  

>// Println 功能类似 Print，只不过最后会添加一个换行符。  
>// 所有参数之间会添加空格，返回写入的字节数。  
>func Println(a ...interface{}) (n int, err error)

>// Printf 将参数列表 a 填写到格式字符串 format 的占位符中。  
>// 填写后的结果写入到标准输出中，返回写入的字节数。  
>func Printf(format string, a ...interface{}) (n int, err error)

```golang
//常见输出
func fmtPrint()  {
    fmt.Print("aaa", 12, []byte("fmt"))
    fmt.Println("aaa", 12, []byte("fmt"))
    fmt.Printf("one type : %v,two type : %v,three type : %v\n", "aaa", 12, []byte("fmt"))
}
```

##### 3、通过Sprint将参数转换为字符串

>// 功能同上面三个函数，只不过将转换结果以字符串形式返回。  
func Sprint(a ...interface{}) string  
func Sprintln(a ...interface{}) string  
func Sprintf(format string, a ...interface{}) string

```golang
//通过Sprint将参数转换为字符串
func fmtSPrint()  {
    s1 := fmt.Sprint("aaa", 12, []byte("fmt"))
    s2 := fmt.Sprintln("aaa", 12, []byte("fmt"))
    s3 := fmt.Sprintf("one type : %v,two type : %v,three type : %v\n", "aaa", 12, []byte("fmt"))

    fmt.Println(s1, s2, s3)
}
```

##### 4、将转换结果写入到write中

>// 功能同上面三个函数，只不过将转换结果写入到 w 中。  
func Fprint(w io.Writer, a ...interface{}) (n int, err error)  
func Fprintln(w io.Writer, a ...interface{}) (n int, err error)  
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error)

```golang
func fmtFPrint() {
    // unbuffered
    fmt.Fprintf(os.Stdout, "%s\n", "hello world! - unbuffered")
    // buffered: os.Stdout implements io.Writer
    buf := bufio.NewWriter(os.Stdout)
    // and now so does buf.
    fmt.Fprintf(buf, "%s\n", "hello world! - buffered")
    buf.Flush()
}
```

##### 5、自定义输出格式

```golang
func fmtCustom()  {
    u := User("Hello gopher")
    fmt.Printf("%-+ 0#8.5m\n", u)
    fmt.Printf("%+ 0#8.5M\n", u)
    fmt.Println(u)
    fmt.Printf("%s\n", u)
    fmt.Printf("%#v\n", u)
    fmt.Printf("%d\n", u)
}

type User string

func (us User) String() string {
    return strings.ToUpper(string(us))
}

func (us User) GoString() string {
    return `"` + strings.ToUpper(string(us)) + `"`
}

func (us User) Format(f fmt.State, c rune)  {
    write := func(s string) {
        f.Write([]byte(s))
    }

    switch c {
    case 'm', 'M':
        write("旗标：[")
        for s := "+- 0#"; len(s) > 0; s = s[1:] {
            if f.Flag(int(s[0])) {
                write(s[1:])
            }
        }
        write("]")

        if v, ok := f.Width(); ok{
            write(" | 宽度：" + strconv.FormatInt(int64(v), 10))
        }

        if v, ok := f.Precision(); ok{
            write(" | 精度：" + strconv.FormatInt(int64(v),10))

        }
    case 's','v':
        if c =='v' && f.Flag('#') {
            write(us.GoString())
        }else {
            write(us.String())
        }
    default:
        write("无效格式：" + string(c))
    }
}
```

#### strings

##### 1、简介

strings包实现了用于操作字符的简单函数。

##### 2、比较

>//比较两个字符串  相等返回0 大于返回1 小于 返回-1  
func Compare(a, b string) int  
//判断两个utf-8编码字符串（将unicode大写、小写、标题三种格式字符视为相同）是否相同。  
>func EqualFold(s, t string) bool

```golang
func stringCompare()  {
    a := "This is a strings"
    b := "THis is a string"
    if strings.EqualFold(a, b) {
        fmt.Println("a equal b")
    }

    rs := strings.Compare(a, b)
    fmt.Printf("result type %T, value %v\n", rs, rs)
}
```

##### 3、转换

>//大写、小写、Title格式转换  
func ToUpper(s string) string  
func ToLower(s string) string  
func ToTitle(s string) string  

>//根据特殊格式进行大写、小写、Title格式转换  
func ToUpperSpecial(_case unicode.SpecialCase, s string) string  
func ToLowerSpecial(_case unicode.SpecialCase, s string) string  
func ToTitleSpecial(_case unicode.SpecialCase, s string) string  

>//返回s中每个单词的首字母都改为标题格式的字符串拷贝。  
**Title用于划分单词的规则不能很好的处理Unicode标点符号。**  
func Title(s string) string

```golang
//大写、小写、Title转换
func stringUpperLower(){
    //返回将所有字母都转为对应的小写版本的拷贝
    fmt.Println(strings.ToLower("Gopher"))
    //返回将所有字母都转为对应的大写版本的拷贝。
    fmt.Println(strings.ToUpper("Gopher"))
    //返回将所有字母都转为对应的标题版本的拷贝。
    fmt.Println(strings.ToTitle("loud noises"))
    fmt.Println(strings.ToTitle("хлеб"))
    //返回s中每个单词的首字母都改为标题格式的字符串拷贝。
    fmt.Println(strings.Title("her royal highness"))
}
```

##### 4、清理

>//返回将s前后端所有cutset包含的utf-8码值都去掉的字符串。  
func Trim(s string, cutset string) string  
func TrimLeft(s string, cutset string) string  
func TrimRight(s string, cutset string) string

>//返回将s前后端所有满足f的unicode码值都去掉的字符串。  
func TrimFunc(s string, f func(rune) bool) string  
func TrimLeftFunc(s string, f func(rune) bool) string  
func TrimRightFunc(s string, f func(rune) bool) string

>//返回将s前后端所有空白（unicode.IsSpace指定）都去掉的字符串。  
func TrimSpace(s string) string  

>//返回去除s可能的前缀prefix(后缀suffix)的字符串。  
func TrimPrefix(s, prefix string) string  
func TrimSuffix(s, suffix string) string  

```golang
func stringTrim()  {
    origin := " !!! Achtung! Achtung! !!! "
    cutset := "! "
    //返回将s前后端所有cutset包含的utf-8码值都去掉的字符串。
    fmt.Printf("[%q] \n", strings.Trim(origin, cutset))
    fmt.Printf("[%q] \n", strings.TrimLeft(origin, cutset))
    fmt.Printf("[%q] \n", strings.TrimRight(origin, cutset))

    fmt.Println(strings.TrimSpace(" \t\n a lone gopher \n\t\r\n"))

    var s = "Goodbye,, world!"
    s = strings.TrimPrefix(s, "Goodbye,")
    s = strings.TrimPrefix(s, "Howdy,")
    fmt.Print("Hello" + s)
}
```

##### 5、分割及拼接

>//根据指定字符串进行切割 返回[]string类型切片  
func Split(s, sep string) []string  
func SplitN(s, sep string, n int) []string

>//根据指定字符串进行切割，切割后含有切割字符，返回[]string类型切片  
func SplitAfter(s, sep string) []string  
func SplitAfterN(s, sep string, n int) []string

>//返回将字符串按照空白（unicode.IsSpace确定，可以是一到多个连续的空白字符）分割的多个字符串。如果字符串全部是空白或者是空字符串的话，会返回空切片。  
func Fields(s string) []string  
func FieldsFunc(s string, f func(rune) bool) []string

>//将一系列字符串连接为一个字符串，之间用sep来分隔。  
func Join(a []string, sep string) string

>//返回count个s串联的字符串。  
func Repeat(s string, count int) string

```golang
//字符串的切割与连接
func stringSplit()  {
    //根据指定字符串进行切割，返回[]string
    fmt.Printf("%q\n", strings.Split("a,b,c", ","))
    fmt.Printf("%q\n", strings.Split("a man a plan a canal panama", "a "))
    fmt.Printf("%q\n", strings.Split(" xyz ", ""))
    fmt.Printf("%q\n", strings.Split("", "Bernardo O'Higgins"))

    //根据指定字符串进行切割，返回[]string 指定了返回的切片个数
    fmt.Printf("%q\n", strings.SplitN("a,b,c", ",", 2))
    z := strings.SplitN("a,b,c", ",", 0)
    fmt.Printf("%q (nil = %v)\n", z, z == nil)

    //切割都的字符串中含有切割字符
    fmt.Printf("%q\n", strings.SplitAfter("a,b,c", ","))

    //默认按空格切割字符串
    fmt.Printf("Fields are: %q", strings.Fields("  foo bar  baz   "))

    //将一个[]string类型切片连接起来
    s := []string{"foo", "bar", "baz"}
    fmt.Println(strings.Join(s, ", "))

    fmt.Println("ba" + strings.Repeat("na", 2))
}
```

##### 6、子串

>判断s是否有前缀字符串prefix(后缀字符串suffix)。  
func HasPrefix(s, prefix string) bool  
func HasSuffix(s, suffix string) bool

>判断字符串s是否包含（子串substr|utf-8码值r|字符串chars中的任一字符）。  
func Contains(s, substr string) bool  
func ContainsRune(s string, r rune) bool  
func ContainsAny(s, chars string) bool

>（子串sep|字符c|unicode码值r|字符串chars中的任一utf-8码值）在字符串s中第一次出现的位置，不存在则返回-1。  
func Index(s, sep string) int  
func IndexByte(s string, c byte) int  
func IndexRune(s string, r rune) int  
func IndexAny(s, chars string) int  
func IndexFunc(s string, f func(rune) bool) int

>（子串sep|字符c|unicode码值r|字符串chars中的任一utf-8码值）在字符串s中最后一次出现的位置，不存在则返回-1。  
func LastIndex(s, sep string) int  
func LastIndexByte(s string, c byte) int  
func LastIndexAny(s, chars string) int  
func LastIndexFunc(s string, f func(rune) bool) int

>返回字符串s中有几个不重复的sep子串。  
func Count(s, sep string) int

```golang
func stringSubstring()  {
    a := "this is a string"
    //前缀
    if strings.HasPrefix(a, "this") {
    fmt.Println("a has prefix 'this'")
    }

    //后缀
    if strings.HasSuffix(a, "string"){
    fmt.Println("a has suffix 'string'")
    }

    //包含子串
    if strings.Contains(a, "is"){
    fmt.Println("a contains 'is'")
    }

    //包含utf-8码值r
    if strings.ContainsRune(a,'a'){
    fmt.Println("a contains 'a'")
    }

    //包含字符串chars中的任一字符
    if strings.ContainsAny(a, "pre"){
    fmt.Println("a contains any string")
    }

    fmt.Println(strings.Count("cheese", "e"))
    fmt.Println(strings.Count("five", "")) // before & after each rune
}
```

##### 7、替换

>返回将s中前n个不重叠old子串都替换为new的新字符串，如果n<0会替换所有old子串。  
func Replace(s, old, new string, n int) string

>将s的每一个unicode码值r都替换为mapping(r)，返回这些新码值组成的字符串拷贝。如果mapping返回一个负值，将会丢弃该码值而不会被替换  
func Map(mapping func(rune) rune, s string) string

```golang
func stringReplace()  {
    fmt.Println(strings.Replace("oink oink oink", "k", "ky", 2))
    fmt.Println(strings.Replace("oink oink oink", "oink", "moo", -1))

    rs := func(r rune) rune {
        switch {
        case r >= 'A' && r <= 'Z':
            return 'A' + (r-'A'+13)%26
        case r >= 'a' && r <= 'z':
            return 'a' + (r-'a'+13)%26
        }
        return r
    }
    fmt.Println(strings.Map(rs, "'Twas brillig and the slithy gopher..."))
}
```

#### strconv

##### 1、简介

strconv包实现了基本数据类型和其字符串表示的相互转换。

##### 2、与bool类型相互转换

>// 将布尔值转换为字符串 true 或 false  
func FormatBool(b bool) string

>// 将字符串转换为布尔值  
// 它接受真值：1, t, T, TRUE, true, True  
// 它接受假值：0, f, F, FALSE, false, False  
// 其它任何值都返回一个错误。  
func ParseBool(str string) (bool, error)

```golang
func strconvBool()  {
    var a bool = true
    rs := strconv.FormatBool(a)
    fmt.Printf("Type is %T, value is %v\n", rs, rs)

    bo, err := strconv.ParseBool(rs)
    if err != nil {
        fmt.Printf("Type is %T, value is %v\n", bo, bo)
    }
}
```

##### 3、与整型相互转换

>// 将整数转换为字符串形式。base 表示转换进制，取值在 2 到 36 之间。  
// 结果中大于 10 的数字用小写字母 a - z 表示。  
func FormatInt(i int64, base int) string  
func FormatUint(i uint64, base int) string  

>// 将字符串解析为整数，ParseInt 支持正负号，ParseUint 不支持正负号。  
// base 表示进位制（2 到 36），如果 base 为 0，则根据字符串前缀判断，  
// 前缀 0x 表示 16 进制，前缀 0 表示 8 进制，否则是 10 进制。  
// bitSize 表示结果的位宽（包括符号位），0 表示最大位宽。  
func ParseInt(s string, base int, bitSize int) (i int64, err error)  
func ParseUint(s string, base int, bitSize int) (uint64, error)  

>// 将整数转换为十进制字符串形式（即：FormatInt(i, 10) 的简写）  
func Itoa(i int) string  

>// 将字符串转换为十进制整数，即：ParseInt(s, 10, 0) 的简写）  
func Atoi(s string) (int, error)

```golang
func strconvInt() {
    var a int64 = 1992
    var b uint64 = uint64(a)
    //整型转换到字符串
    rs1 := strconv.FormatInt(a, 2)
    rs2 := strconv.FormatUint(b, 2)
    fmt.Printf("Type is %T, value is %v\n", rs1, rs1)
    fmt.Printf("Type is %T, value is %v\n", rs2, rs2)

    rs3 := strconv.Itoa(1992)
    fmt.Printf("Type is %T, value is %v\n", rs3, rs3)

    //字符串换换到整型
    result1, err := strconv.ParseInt(rs1, 2, 0)
    if err == nil {
        fmt.Printf("Type is %T, value is %v\n", result1, result1)
    }

    result2, err := strconv.ParseUint(rs2, 2, 0)
    if err == nil {
        fmt.Printf("Type is %T, value is %v\n", result2, result2)
    }

    result3, err := strconv.Atoi(rs3)
    if err == nil {
        fmt.Printf("Type is %T, value is %v\n", result3, result3)
    }
}
```

##### 4、与浮点数相互转换

>// FormatFloat 将浮点数 f 转换为字符串形式  
// f：要转换的浮点数  
// fmt：格式标记（b、e、E、f、g、G）  
// prec：精度（数字部分的长度，不包括指数部分）  
// bitSize：指定浮点类型（32:float32、64:float64），结果会据此进行舍入。  
//  
// 格式标记：  
// 'b' (-ddddp±ddd，二进制指数)  
// 'e' (-d.dddde±dd，十进制指数)  
// 'E' (-d.ddddE±dd，十进制指数)  
// 'f' (-ddd.dddd，没有指数)  
// 'g' ('e':大指数，'f':其它情况)  
// 'G' ('E':大指数，'f':其它情况)  
//  
// 如果格式标记为 'e'，'E'和'f'，则 prec 表示小数点后的数字位数  
// 如果格式标记为 'g'，'G'，则 prec 表示总的数字位数（整数部分+小数部分）  
// 参考格式化输入输出中的旗标和精度说明  
func FormatFloat(f float64, fmt byte, prec, bitSize int) string  

>// 将字符串解析为浮点数，使用 IEEE754 规范进行舍入。  
// bigSize 取值有 32 和 64 两种，表示转换结果的精度。  
// 如果有语法错误，则 err.Error = ErrSyntax  
// 如果结果超出范围，则返回 ±Inf，err.Error = ErrRange  
func ParseFloat(s string, bitSize int) (float64, error)  

```golang
func strconvFloat()  {
    //浮点数转换成字符串
    or := 0.123456789012345
    fmt.Printf("Type is %T, value is %v\n", or, or)

    rs := strconv.FormatFloat(or,'e',20, 64)
    fmt.Printf("Type is %T, value is %v\n", rs, rs)


    //字符串转换成浮点数
    s := "0.12345678901234567890"

    f, err := strconv.ParseFloat(s, 32)
    fmt.Println(f, err)                // 0.12345679104328156
    fmt.Println(float32(f), err)       // 0.12345679

    f, err = strconv.ParseFloat(s, 64)
    fmt.Println(f, err)                // 0.12345678901234568
}
```

##### 5、将各种类型转换为字符串后追加到 dst 尾部

>// 将各种类型转换为字符串后追加到 dst 尾部。  
func AppendInt(dst []byte, i int64, base int) []byte  
func AppendUint(dst []byte, i uint64, base int) []byte  
func AppendFloat(dst []byte, f float64, fmt byte, prec, bitSize int) []byte  
func AppendBool(dst []byte, b bool) []byte  

```golang
func strconvAppend()  {
    var origin []byte = []byte("this is a string")
    var a int64 = 1992
    rs := strconv.AppendInt(origin,a,10)
    fmt.Printf("Type is %T\n, value is %v\n, string is %s\n", rs, rs, rs)
}
```

##### 6、特殊字符转换

>// 判断字符串是否可以不被修改的表示为一个单行的反引号字符串。  
// 字符串中不能含有控制字符（除了 \t）和“反引号”字符，否则返回 false  
func CanBackquote(s string) bool  

>// 判断 r 是否为可打印字符  
// 可否打印并不是你想象的那样，比如空格可以打印，而\t则不能打印  
func IsPrint(r rune) bool

>// 判断 r 是否为 Unicode 定义的图形字符。  
func IsGraphic(r rune) bool

>// 将 s 转换为双引号字符串  
func Quote(s string) string

>// 功能同上，非 ASCII 字符和不可打印字符会被转义  
func QuoteToASCII(s string) string

>// 功能同上，非图形字符会被转义  
func QuoteToGraphic(s string) string

>// 将 r 转换为单引号字符  
func QuoteRune(r rune) string

>// 功能同上，非 ASCII 字符和不可打印字符会被转义  
func QuoteRuneToASCII(r rune) string

>// 功能同上，非图形字符会被转义  
func QuoteRuneToGraphic(r rune) string

>// 将 r 转换为单引号字符  
func QuoteRune(r rune) string

>// 功能同上，非 ASCII 字符和不可打印字符会被转义  
func QuoteRuneToASCII(r rune) string

>// 功能同上，非图形字符会被转义  
func QuoteRuneToGraphic(r rune) string

>// Unquote 将“带引号的字符串” s 转换为常规的字符串（不带引号和转义字符）  
// s 可以是“单引号”、“双引号”或“反引号”引起来的字符串（包括引号本身）  
// 如果 s 是单引号引起来的字符串，则返回该该字符串代表的字符  
func Unquote(s string) (string, error)

>// UnquoteChar 将带引号字符串（不包含首尾的引号）中的第一个字符“取消转义”并解码  
//  
// s    ：带引号字符串（不包含首尾的引号）  
// quote：字符串使用的“引号符”（用于对字符串中的引号符“取消转义”）  
//  
// value    ：解码后的字符  
// multibyte：value 是否为多字节字符  
// tail     ：字符串 s 解码后的剩余部分  
// error    ：返回 s 中是否存在语法错误  
//  
// 参数 quote 为“引号符”  
// 如果设置为单引号，则 s 中允许出现 \'、" 字符，不允许出现单独的 ' 字符  
// 如果设置为双引号，则 s 中允许出现 \"、' 字符，不允许出现单独的 " 字符  
// 如果设置为 0，则不允许出现 \' 或 \" 字符，但可以出现单独的 ' 或 " 字符  
func UnquoteChar(s string, quote byte) (value rune, multibyte bool, tail string, err error)  

### bytes

##### 1、简介

bytes包实现了操作[]byte的常用函数

>type byte byte  
8位无符号整型，是uint8的别名，二者视为同一类型。

##### 2、转换（大小写）

> 将 s 中的所有字符修改为大写（小写、标题）格式返回。  
func ToUpper(s []byte) []byte  
func ToLower(s []byte) []byte  
func ToTitle(s []byte) []byte  
// 使用指定的映射表将 s 中的所有字符修改为大写（小写、标题）格式返回。  
func ToUpperSpecial(_case unicode.SpecialCase, s []byte) []byte  
func ToLowerSpecial(_case unicode.SpecialCase, s []byte) []byte  
func ToTitleSpecial(_case unicode.SpecialCase, s []byte) []byte  
// 将 s 中的所有单词的首字符修改为 Title 格式返回。  
// BUG: **不能很好的处理以 Unicode 标点符号分隔的单词。**  
func Title(s []byte) []byte  

```golang
func bytesUpperLower()  {
    a := []byte("I,am,is,a,gopher")
    upper := bytes.ToUpper(a)

    fmt.Printf("全部转换为大写:%q\n",upper)
    lower := bytes.ToLower(a)
    fmt.Printf("全部转换成小写:%q\n",lower)
    title := bytes.ToTitle([]byte("learn go languager"))
    fmt.Printf("全部转换成Title:%q\n",title)

    fmt.Println(bytes.ToLowerSpecial(unicode.TurkishCase, []byte("dünyanın ilk borsa yapısı Aizonai kabul edilir")))
}
```

##### 3、分割及拼接

>// Split 以 sep 为分隔符将 s 切分成多个子串，结果不包含分隔符。  
// 如果 sep 为空，则将 s 切分成 Unicode 字符列表。  
// SplitN 可以指定切分次数 n，超出 n 的部分将不进行切分。  
func Split(s, sep []byte) [][]byte  
func SplitN(s, sep []byte, n int) [][]byte  

>// 功能同 Split，只不过结果包含分隔符（在各个子串尾部）。  
func SplitAfter(s, sep []byte) [][]byte  
func SplitAfterN(s, sep []byte, n int) [][]byte  

>// 以连续空白为分隔符将 s 切分成多个子串，结果不包含分隔符。  
func Fields(s []byte) [][]byte  
// 以符合 f 的字符为分隔符将 s 切分成多个子串，结果不包含分隔符。  
func FieldsFunc(s []byte, f func(rune) bool) [][]byte  

>// 以 sep 为连接符，将子串列表 s 连接成一个字节串。  
func Join(s [][]byte, sep []byte) []byte  
// 将子串 b 重复 count 次后返回。  
func Repeat(b []byte, count int) []byte

```golang
//截断
func bytesSplit()  {
    a := []byte("  Hello   World !  ")
    b := []byte("I,am,is,a,gopher")
    c := []byte("I am is a gopher")
    d :=[][]byte{
    a,
    b,
    c,
    }
    fmt.Printf("%q\n", bytes.Split(b, []byte{','}))
    fmt.Printf("%q\n", bytes.SplitN(b,[]byte{','}, 3))

    fmt.Printf("%q\n", bytes.SplitAfter(b, []byte{','}))
    fmt.Printf("%q\n", bytes.SplitAfterN(b, []byte{','}, 3))

    fmt.Printf("%q\n", bytes.Fields(c))

    f := func(r rune) bool {
    return bytes.ContainsRune([]byte(" o"), r)
    }
    fmt.Printf("%q\n", bytes.FieldsFunc(a, f))
    //连接
    fmt.Printf("%q\n", bytes.Join(d, []byte{'x'}))
    //重复
    fmt.Printf("%q\n", bytes.Repeat(a,3))
}
```

##### 4、清理

>// 去掉 s 两边（左边、右边）包含在 cutset 中的字符（返回 s 的切片）  
func Trim(s []byte, cutset string) []byte  
func TrimLeft(s []byte, cutset string) []byte  
func TrimRight(s []byte, cutset string) []byte  

>// 去掉 s 两边（左边、右边）符合 f 要求的字符（返回 s 的切片）  
func TrimFunc(s []byte, f func(r rune) bool) []byte  
func TrimLeftFunc(s []byte, f func(r rune) bool) []byte  
func TrimRightFunc(s []byte, f func(r rune) bool) []byte  

>// 去掉 s 两边的空白（unicode.IsSpace）（返回 s 的切片）  
func TrimSpace(s []byte) []byte  
// 去掉 s 的前缀 prefix（后缀 suffix）（返回 s 的切片）  
func TrimPrefix(s, prefix []byte) []byte  
func TrimSuffix(s, suffix []byte) []byte  

```golang
//清理
func bytesTrim()  {
    bs := [][]byte{
    []byte("Hello World ！"),
    []byte("Hello 世界！"),
    []byte("hello golang ."),
    }

    f := func(r rune) bool{
    return bytes.ContainsRune([]byte("!！.。"), r)
    }

    for _,b := range bs{
    fmt.Printf("%q\n", bytes.TrimFunc(b, f))
    }

    for _, b := range bs {
    fmt.Printf("%q\n", bytes.TrimPrefix(b, []byte("Hello")))
    }

    a := []byte(",I am a gopher,")
    fmt.Printf("Origin:%q\n", a)
    fmt.Printf("Trim:%q\n", bytes.Trim(a,",I"))
    fmt.Printf("TrimRight:%q\n", bytes.TrimRight(a,","))
    fmt.Printf("TrimLeft:%q\n", bytes.TrimLeft(a,","))

    b := []byte(" I am a gopher ")
    fmt.Printf("Origin:%q\n", b)
    fmt.Printf("TrimSpace:%q\n", bytes.TrimSpace(b))
}
```

##### 5、比较

>// 比较两个 []byte，nil 参数相当于空 []byte。  
// a <  b 返回 -1  
// a == b 返回 0  
// a >  b 返回 1  
func Compare(a, b []byte) int  

>// 判断 a、b 是否相等，nil 参数相当于空 []byte。  
func Equal(a, b []byte) bool  

>// 判断 s、t 是否相似，忽略大写、小写、标题三种格式的区别。  
// 参考 unicode.SimpleFold 函数。  
func EqualFold(s, t []byte) bool  

```golang
func bytesCompare()  {
      a := []byte("I am a gopher!")
      b := []byte("i am a gopher!")
   // Interpret Compare's result by comparing it to zero.

   if bytes.Compare(a, b) < 0 {
    fmt.Println("a less b")
   }
   if bytes.Compare(a, b) <= 0 {
    fmt.Println("a less or equal b")
   }
   if bytes.Compare(a, b) > 0 {
    fmt.Println("a greater b")
   }
   if bytes.Compare(a, b) >= 0 {
    fmt.Println("a greater or equal b")
   }
   //返回a、b是否相等
   if bytes.Equal(a, b) {
    fmt.Println("a equal b")
   }
   if !bytes.Equal(a, b) {
    fmt.Println("a not equal b")
   }

   //忽略大小写.
   if bytes.EqualFold(a, b) {
    fmt.Println("a equal b")
   }
   if !bytes.EqualFold(a, b) {
    fmt.Println("a not equal b")
   }
}
```

##### 6、查询子串

>// 判断 s 是否有前缀 prefix（后缀 suffix）  
func HasPrefix(s, prefix []byte) bool  
func HasSuffix(s, suffix []byte) bool  

>// 判断 b 中是否包含子串 subslice（字符 r）  
func Contains(b, subslice []byte) bool  
func ContainsRune(b []byte, r rune) bool  

>// 判断 b 中是否包含 chars 中的任何一个字符  
func ContainsAny(b []byte, chars string) bool  

>// 查找子串 sep（字节 c、字符 r）在 s 中第一次出现的位置，找不到则返回 -1。  
func Index(s, sep []byte) int  
func IndexByte(s []byte, c byte) int  
func IndexRune(s []byte, r rune) int  

>// 查找 chars 中的任何一个字符在 s 中第一次出现的位置，找不到则返回 -1。  
func IndexAny(s []byte, chars string) int  

>// 查找符合 f 的字符在 s 中第一次出现的位置，找不到则返回 -1。  
func IndexFunc(s []byte, f func(r rune) bool) int  

>// 功能同上，只不过查找最后一次出现的位置。  
func LastIndex(s, sep []byte) int  
func LastIndexByte(s []byte, c byte) int  
func LastIndexAny(s []byte, chars string) int  
func LastIndexFunc(s []byte, f func(r rune) bool) int  

>// 获取 sep 在 s 中出现的次数（sep 不能重叠）。  
func Count(s, sep []byte) int  

##### 7、替换

>// 将 s 中前 n 个 old 替换为 new，n < 0 则替换全部。  
func Replace(s, old, new []byte, n int) []byte  

>// 将 s 中的字符替换为 mapping(r) 的返回值，  
// 如果 mapping 返回负值，则丢弃该字符。  
func Map(mapping func(r rune) rune, s []byte) []byte

>// 将 s 转换为 []rune 类型返回  
func Runes(s []byte) []rune

```golang
//替换
func bytesReplace(){
    a := []byte("I am a gopher~")
    //将 s 中前 n 个 old 替换为 new，n < 0 则替换全部。
    fmt.Printf("%q\n", bytes.Replace(a, []byte("gopher"), []byte("phper"), -1))

    //将 s 转换为 []rune 类型返回
    fmt.Printf("%q\n", bytes.Runes(a))


    b := []byte("abcabcabcabc")

    // Map 将 s 中满足 mapping(rune) 的字符替换为 mapping(rune) 的返回值
    // 如果 mapping(rune) 返回负数，则相应的字符将被删除
    c := bytes.Map(func(r rune) rune {
        if r == 'c' {
        return 'a'
        }
        return r
    }, b)
    println(string(c))
}
```

##### 8、字节缓存

>bytes.Buffer是一个实现了读写方法的可变大小的字节缓冲。本类型的零值是一个空的可用于读写的缓冲。

```golang
//字节缓存案列
func bytesBuffer()  {
    rd := bytes.NewBufferString("Hello World!")
    buf := make([]byte, 6)
    // 获取数据切片
    b := rd.Bytes()
    // 读出一部分数据，看看切片有没有变化
    rd.Read(buf)
    fmt.Printf("%s\n", rd.String()) // World!
    fmt.Printf("%s\n\n", b)         // Hello World!

    // 写入一部分数据，看看切片有没有变化
    rd.Write([]byte("abcdefg"))
    fmt.Printf("%s\n", rd.String()) // World!abcdefg
    fmt.Printf("%s\n\n", b)         // Hello World!

    // 再读出一部分数据，看看切片有没有变化
    rd.Read(buf)
    fmt.Printf("%s\n", rd.String()) // abcdefg
    fmt.Printf("%s\n", b)           // Hello World!
}
```

### 单元测试

#### testing

##### 1、简介

测试文件包含test functions（自动化测试）, benchmark functions（基准测试）, and example functions

##### 2、自动化测试

```golang
//执行所有_test.go文件
go test

//输出详细的测试信息
go test -v

//指定文件
go test -v main.go

//执行指定函数
go test -v -test.run TestAdd

//正则匹配测试函数
go test -v -run A
```

##### 3、基准测试

- 执行测试命令  

>//执行所有测试用例
>go test -bench=.

>//仅执行基准测试用例
>//-run=none
>go test -bench=. -run=none

>//执行具体的基准测试函数(BenchmarkArray)
>go test -bench=Array -run=none

>//内容消耗情况
>//-benchmem可以提供每次操作分配内存的次数，以及每次操作分配的字节数
>go test -bench=Array -run=none -benchmem

>//正则匹配所有Plus开头的基准测试函数
>go test -bench Plus

- 测试结果分析

>➜  test go test -v -bench=Array -run=none
>goos: linux
>goarch: amd64
>pkg: test/test
>BenchmarkArray-2 1000000  1342 ns/op
>PASS
>ok  test/test 1.368s

>//意味着循环执行了 1000000 次，每次循环花费 1342 纳秒(ns)。

- 常用函数

>b.ResetTimer() 是重置计时器，这样可以避免for循环之前的初始化代码的干扰

##### 4、示例测试

```golang
 func ExampleHello() {
  fmt.Println("hello")
  // Output: hello
 }

 //不考虑输出顺序
 func ExamplePerm() {
  for _, value := range Perm(4) {
  fmt.Println(value)
  }
  // Unordered output: 4
  // 2
  // 1
  // 3
  // 0
 }
```

- 注意

**没有输出注释的示例函数被编译但不执行**。

###### 5、参考

- 1 [基准测试](https://zhangwenbing.com/blog/golang/HyNlPG-Fq8Q)

### 时间

#### time

##### 1、基础功能

（1）、获取当前时间及当前时间戳

```golang
//获取当前时间
func getCurrentTime(layout string)  string{
    if layout == "" {
    layout = "2006-01-02 15:04:05"
    }
    return  time.Now().Format(layout)
}

//result : 2018-11-14 15:00:24



//获取当前时间戳
func getCurrentTimeUnix()  int64{
    return time.Now().Unix()
}

//result : 1542178824
```

（2）、将给定时间字符串转换成时间类型及时间戳

```golang
//将时间字符串转换成time类型，当前时区
func formatStringToTime(layout, value string) (time.Time, error) {
    if layout == "" {
        layout = "2006-01-02 15:04:05"
    }
    return time.ParseInLocation(layout, value, time.Now().Location())
}

//param : "", "2018-11-14 15:07:12"
//result : 2018-11-14 15:07:12 +0800 CST
```

（3）、将给定的时间戳转换成时间类型及时间戳

```golnag
//将整型转换为time类型
func formatIntToTime(value int64) time.Time{
    return time.Unix(value, 0)
}

var curUnix int64= 1541997444
curTime := formatIntToTime(curUnix)
fmt.Printf("curUnix:%T, cutTime:%T", curUnix, curTime)
//result : curUnix:int64, cutTime:time.Time#
```

（4）、计算两个时间的区段

```golang
//获取两个时间的时间间隔
func getTimeSub(startTime, endTime time.Time) time.Duration{
    return endTime.Sub(startTime)
}

//获取当前到指定时间的时间间隔
func getTimeSince(startTime time.Time) time.Duration {
    return time.Since(startTime)
}
```

### 参考

- 1 [Golang包](https://www.cnblogs.com/golove/tag/Golang%E5%8C%85/)
