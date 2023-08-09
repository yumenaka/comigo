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
        <div>
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
