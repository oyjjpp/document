由于某些JS事件频繁被触发，因而频繁执行DOM操作、资源加载等重行为，导致UI停顿甚至浏览器崩溃。

例如：

1、window对象的resize、scroll事件

2、mousemove事件

throttle：

定义：

预先设定一个执行周期，当调用动作的时刻大于等于执行周期则执行该动作，然后进入下一个新周期。


/**
 * @param {function} func, 需要节流的函数
 * @param {number} wait, 节流间隔时间，单位ms
 * @param {boolean} leading，默认首次执行, false禁用首次执行
 * @demo throttle(func, opts.time)
 *       throttle(func, opts.time, {leading: false})
 */
var throttle = function(func, wait, options) {

    var context, args, result = null;
    var timeout = null;
    var previous = 0;
    if (!options) options = {};
    var later = function() {
        previous = options.leading === false ? 0 : Date.now();
        timeout = null;
        result = func.apply(context, args);
        if (!timeout) context = args = null;
    };

    return function() {
        var now = Date.now();
        if (!previous && options.leading === false) previous = now;
        var remaining = wait - (now - previous);
        context = this;
        args = arguments;
        if (remaining <= 0 || remaining > wait) {
            if (timeout) {
                clearTimeout(timeout);
                timeout = null;
            }
            previous = now;
            result = func.apply(context, args);
            if (!timeout) context = args = null;
        } else if (!timeout) {
            timeout = setTimeout(later, remaining);
        }
        return result;
    };
}

/** 
 *  options的默认值
 *  表示首次调用返回值方法时，会马上调用func；否则仅会记录当前时刻，当第二次调用的时间间隔超过wait时，才调用func。
 *  options.leading = true;
 *                         表示当调用方法时，未到达wait指定的时间间隔，则启动计时器延迟调用func函数，若后续在既未达到wait指定的时间间隔和func函数又未被调用的情况下调用返回值方法，则被调用请求将被丢弃。
 *  options.trailing = true; 
 *  注意：当options.trailing = false时，效果与上面的简单实现效果相同
 */
_.throttle = function(func, wait, options) {
    
    var context, args, result;
    var timeout = null;
    var previous = 0;
    if (!options) options = {};
    var later = function() {
        previous = options.leading === false ? 0 : _.now();
        timeout = null;
        result = func.apply(context, args);
        if (!timeout) context = args = null;
    };
    return function() {
        var now = _.now();
        if (!previous && options.leading === false) previous = now;
        // 计算剩余时间
        var remaining = wait - (now - previous);
        context = this;
        args = arguments;
        // 当到达wait指定的时间间隔，则调用func函数
        // 精彩之处：按理来说remaining <=0已经足够证明已经到达wait的时间间隔，
        // 但这里还考虑到假如客户端修改了系统时间则马上执行func函数。
        if (remaining <= 0 || remaining > wait) {
            // 由于setTimeout存在最小时间精度问题，因此会存在到达wait的时间间隔，但之前设置的setTimeout操作还没被执行，因此为保险起见，这里先清理setTimeout操作
            if (timeout) {
                clearTimeout(timeout);
                timeout = null;
            }
            previous = now;
            result = func.apply(context, args);
            if (!timeout) context = args = null;
        } else if (!timeout && options.trailing !== false) {
            // options.trailing=true时，延时执行func函数
            timeout = setTimeout(later, remaining);
        }
        return result;
    };
};
  
### 参考

http://www.cnblogs.com/fsjohnhuang/p/4147810.html