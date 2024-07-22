import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import Header from "../layout/Header";

const Home: React.FC = () => {
    const [nameValue, setNameValue] = useState<string>("");
    const navigate = useNavigate();

    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setNameValue(event.target.value);
    };

    const handleCreate = () => {
        navigate("/game/tempID");
    };

    const handleJoin = () => {
        navigate("/join");
    };

    return (
        <div className="home-container">
            <Header />
            <main>
                <input placeholder="ENTER NAME" value={nameValue} onChange={handleInputChange} autoFocus />
                <button type="button" disabled={!nameValue} onClick={handleCreate}>CREATE ROOM</button>
                <button type="button" disabled={!nameValue} onClick={handleJoin}>JOIN ROOM</button>
            </main>
        </div>
    );
}

export default Home;
