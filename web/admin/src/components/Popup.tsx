import { useState, useEffect } from "react";
type PropsType = {
    showPopup: boolean
    popupMessage: string
    onClosePopup: () => void
}
// import Popup from "./components/Popup";
//       {/* Popup消息 */}
//       <Popup showPopup={showPopup} popupMessage={popupMessage} onClosePopup={() => setShowPopup(false)}></Popup>

//   //Popup Message
//   const [showPopup, setShowPopup] = useState(true)
//   const [popupMessage, setPopupMessage] = useState("")
//   const onOpenPopup = (message: string) => {
//     setPopupMessage(message)
//     setShowPopup(true)
//   }


const Popup = (props: PropsType) => {
    const { showPopup, popupMessage: message, onClosePopup } = props
    const [opacity, setOpacity] = useState(1)
    const onClose = () => {
        onClosePopup()
    }

    useEffect(() => {
        setTimeout(() => {
            onClosePopup()
        }, 2000)
        let timer:number = 100.0
        const timerId = setInterval(() => {
            timer -= 5.0
            let tempOpacity = timer / 100.0
            if (tempOpacity <= 0) {
                tempOpacity = 0
            }
            console.log(tempOpacity)
            setOpacity(tempOpacity)
            if (showPopup === false||tempOpacity<=0) {
                setOpacity(1)
                clearInterval(timerId)
                console.log("clearInterval")
                return
            }
        }, 100)
    }, [showPopup]); // Add onClosePopup to the dependency array


    return (
        <div role="alert" className="m-4 rounded-xl border border-gray-100 bg-white p-4 fixed" style={showPopup ? { opacity: opacity, display: "block" } : { display: "none" }}>

            <div className="flex items-start gap-4">
                <span className="text-green-600">
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                        strokeWidth="1.5"
                        stroke="currentColor"
                        className="h-6 w-6"
                    >
                        <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                        />
                    </svg>
                </span>

                <div className="flex-1">
                    <strong className="block font-medium text-gray-900"> Changes saved </strong>

                    <p className="mt-1 text-sm text-gray-700">
                        Your product changes have been saved.
                        {message}
                    </p>
                </div>

                <button onClick={onClose} className="text-gray-500 transition hover:text-gray-600">
                    <span className="sr-only">关闭弹出窗口 Dismiss popup</span>

                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                        strokeWidth="1.5"
                        stroke="currentColor"
                        className="h-6 w-6"
                    >
                        <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            d="M6 18L18 6M6 6l12 12"
                        />
                    </svg>
                </button>
            </div>
        </div>

    )
}

export default Popup