

type PropsType = {
    name: string
    label: string
    fieldDescription: string;
    boolValue: boolean
    setBoolValue: (checked: boolean, valueName: string) => void
}

const BoolSwitch = (props: PropsType) => {
    const { name, label: nameText, fieldDescription, boolValue, setBoolValue } = props
    const onChange = () => {
        setBoolValue(!boolValue, name)
    }
    return (
        <div className="w-full m-1 p-2 flex flex-col font-semibold rounded-md shadow-md bg-blue-100 justify-left items-left">
            <div className="w-32">{nameText}</div>

            <label htmlFor="AcceptConditions" className="relative h-8 w-14 cursor-pointer">
                <input type="checkbox"  checked={boolValue}  id="AcceptConditions" className="peer sr-only" onChange={onChange} />

                <span
                    className="absolute inset-0 rounded-full bg-gray-300 transition peer-checked:bg-green-500"
                ></span>

                <span
                    className="absolute inset-y-0 start-0 m-1 h-6 w-6 rounded-full bg-white transition-all peer-checked:start-6"
                ></span>
            </label>

            <div className="py-1 w-3/4 text-xs text-gray-500">{fieldDescription}</div>
        </div>
    )
}

export default BoolSwitch
