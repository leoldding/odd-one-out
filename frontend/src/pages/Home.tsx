import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import Header from "../layout/Header";

const Home: React.FC = () => {
    const [nameValue, setNameValue] = useState<string>("");
    const navigate = useNavigate();

    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setNameValue(event.target.value);
    };

    const handleSubmit = () => {
        navigate("/game/tempID");
    };

    return (
        <div className="home-container">
            <Header />
            <main>
                <form onSubmit={handleSubmit}>
                    <input placeholder="Enter name" value={nameValue} onChange={handleInputChange} autoFocus />
                    <button type="submit" disabled={!nameValue}>Create Room</button> {/* call backend to create room with uuid and token */}
                </form>
            </main>
        </div>
    );
}

export default Home;
