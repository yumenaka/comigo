import { Switch } from '@headlessui/react'


type PropsType = {
    name: string
    label: string
    boolValue: boolean
    setBoolValue: (checked: boolean, valueName: string) => void
}

const BoolSwitch = (props: PropsType) => {
    const { name, label: nameText, boolValue, setBoolValue } = props

    const onChange = (checked: boolean) => {
        setBoolValue(checked, name)
    }

    return (
        <div className="flex flex-row w-96 m-1 p-2 pl-8 font-semibold rounded-md shadow-md bg-yellow-300 justify-left items-center">
            <div className="w-32 border border-black rounded-md">{nameText}</div>
            <Switch
                className={`${boolValue ? 'bg-blue-500' : 'bg-gray-400'}
                    px-1.5 py-0 m-0 relative inline-flex h-[30px] w-[76px] shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus-visible:ring-2 focus-visible:ring-white focus-visible:ring-opacity-75`}
                onChange={onChange}
                name={name}
            >
                <span className="sr-only">Use setting</span>
                <span
                    aria-hidden="true"
                    className={`${boolValue ? 'translate-x-9' : 'translate-x-0'}
                        pointer-events-none inline-block h-[25px] w-[25px] transform rounded-full bg-white shadow-lg ring-0 transition duration-200 ease-in-out`}
                />
            </Switch>
        </div>
    )
}

export default BoolSwitch
