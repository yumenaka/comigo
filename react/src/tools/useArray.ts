import { useState } from "react"
// https://github.com/WebDevSimplified/useful-custom-react-hooks/tree/main/src/5-useArray
// https://www.youtube.com/watch?v=0eXPG4Lep1o
export default function useArray<T>(defaultValue: T[]) {
    const [array, setArray] = useState<T[]>(defaultValue)

    function push(element: T): void {
        setArray((a: T[]) => [...a, element])
    }

    function filter(callback: (element: T, index: number, array: T[]) => boolean): void {
        setArray((a: T[]) => a.filter(callback))
    }

    function update(index: number, newElement: T): void {
        setArray((a: T[]) => [
            ...a.slice(0, index),
            newElement,
            ...a.slice(index + 1),
        ])
    }

    function remove(index: number): void {
        setArray((a: T[]) => [...a.slice(0, index), ...a.slice(index + 1)])
    }

    function clear(): void {
        setArray([])
    }

    return { array, set: setArray, push, filter, update, remove, clear }
}