import React from "react";
import { useTranslation } from "react-i18next";

interface Props {
  group: string;
  setGroup: React.Dispatch<React.SetStateAction<string>>;
  InterfaceColor: string;
}

const Contained: React.FC<Props> = ({
  group,
  setGroup,
  InterfaceColor,
}) => {
  const { t } = useTranslation();

  return (
    <div
      style={{
        backgroundColor: InterfaceColor, // 绑定样式
      }}
      className={`min-w-72 mx-2 p-0 justify-center inline-flex`}>
      <button
        className={`mx-1 ${group === "bookstore" ? " bg-white text-blue-500" : "bg-gray-300 text-gray-500 hover:text-gray-700"} inline-flex items-center gap-2 rounded-md px-4 py-2 text-sm font-semibold shadow-sm focus:relative `}
        onClick={() => setGroup("bookstore")}
      >
        {/* https://heroicons.com/ */}
        <svg className="m-0 w-6 h-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" >
          <path strokeLinecap="round" strokeLinejoin="round" d="M12 6.042A8.967 8.967 0 006 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 016 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 016-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0018 18a8.967 8.967 0 00-6 2.292m0-14.25v14.25" />
        </svg>
        <span className="m-0">{t("book_shelf")}</span>
      </button>

      <button
        className={`mx-1 ${group === "internet" ? " bg-white text-blue-500" : "bg-gray-300 text-gray-500 hover:text-gray-700"} inline-flex items-center gap-2 rounded-md px-4 py-2 text-sm shadow-sm font-semibold focus:relative `}
        onClick={() => setGroup("internet")}
      >
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-6 h-6">
          <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />
        </svg>
        <span className="m-0">{t("network")}</span>
      </button>

      <button
        className={`mx-1 ${group === "other" ? " bg-white text-blue-500" : "bg-gray-300 text-gray-500 hover:text-gray-700"} inline-flex items-center gap-2 rounded-md px-4 py-2 text-sm shadow-sm font-semibold focus:relative `}
        onClick={() => setGroup("other")}
      >
        <svg className="w-6 h-6" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><g fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M16 4l3 3L6.65 19.35a1.5 1.5 0 0 1-3-3L16 4"></path><path d="M10 10h6"></path><path d="M19 15l1.5 1.6a2 2 0 1 1-3 0L19 15"></path></g></svg>
        <span className="m-0">{t("abs")}</span>
      </button>
    </div>
  );
};

export default Contained;
