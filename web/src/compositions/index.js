// 懒加载，参考《Vue3项目中如何实现数据懒加载》：https://blog.csdn.net/qq_42651173/article/details/118913980

import { ref } from 'vue'
//vueuse,基于Vue组合式API的函数工具集: https://bbs.huaweicloud.com/blogs/324277
// 它的初衷就是将一切原本并不支持响应式的JS API变得支持响应式，省去程序员自己写相关代码。
// https://vueuse.org/integrations/useqrcode/
import { useIntersectionObserver } from '@vueuse/core'

/**
 * 功能: 数据懒加载
 *
 * @param {*} fn  当目标可见时，要调用一次的函数
 * @returns target: 要观察的目标（vue3的引用）
 */
const useLazyData = (fn) => {
    const target = ref(null)
    const { stop } = useIntersectionObserver(
        target, // target 是vue的对象引用。是观察的目标
        // isIntersecting 是否进入可视区域，true是进入 false是移出
        // observerElement 被观察的dom
        ([{ isIntersecting }], observerElement) => {
            if (isIntersecting) { // 目标可见，
                // 1. ajax可以发了，后续不需要观察了
                stop()
                // 2. 执行函数
                fn()
            }
        }, {
        threshold: 0//threshold ，容器和可视区交叉的占比（进入的面积/容器完整面积） 取值，0-1 之间，默认是0.1，所以需要滚动较多才能触发进入可视区域事件。
    }
    )
    return target
}

export default useLazyData