import React, { useState } from "react";
import Header from "../layout/Header";

const Join: React.FC = () => {
    const [nameValue, setNameValue] = useState<string>("");

    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setNameValue(event.target.value);
    };

    return (
        <div className="join-container">
            <Header />
            <main>
                <input placeholder="Enter name" value={nameValue} onChange={handleInputChange} autoFocus />
                <button type="button" disabled={!nameValue}>Join Room</button> 
            </main>
        </div>
    );
}

export default Join;

