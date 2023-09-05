import React, { ChangeEvent } from "react";

interface BoolRadioButtonGroupProps {
    options: string[];
    selectedOption: string;
    onChange: (event: ChangeEvent<HTMLInputElement>) => void;
}

const BoolRadioButtonGroup: React.FC<BoolRadioButtonGroupProps> = ({
    options,
    selectedOption,
    onChange,
}) => {
    return (
        <div className="flex flex-row w-2/3 font-semibold rounded-md shadow-md bg-yellow-100 justify-start items-center">
            {options.map((option) => (
                <label key={option}>
                    <input
                        type="radio"
                        value={option}
                        checked={selectedOption === option}
                        onChange={onChange}
                    />
                    {option}
                </label>
            ))}
        </div>
    );
};

export default BoolRadioButtonGroup;
