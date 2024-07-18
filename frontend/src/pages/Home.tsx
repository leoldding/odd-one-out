import React, { useState } from "react";
import Header from "../layout/Header";

const Home: React.FC = () => {
    const [nameValue, setNameValue] = useState<string>("");

    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setNameValue(event.target.value);
    };

    return (
        <div className="home-container">
            <Header />
            <main>
                <input placeholder="Enter name" value={nameValue} onChange={handleInputChange} autoFocus />
                <button type="button" disabled={!nameValue}>Create Room</button> {/* call backend to create room with uuid and token */}
            </main>
        </div>
    );
}

export default Home;
