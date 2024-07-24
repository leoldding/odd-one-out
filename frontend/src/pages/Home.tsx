import React, { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import Header from "../layout/Header";

const Home: React.FC = () => {
    const { code } = useParams<string>();

    const [nameValue, setNameValue] = useState<string>("");
    const navigate = useNavigate();

    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setNameValue(event.target.value);
    };

    const handleCreate = () => {
        navigate("/game/createdID");
    };

    const handleJoin = () => {
        navigate("/game/" + code);
    };

    return (
        <div className="home-container">
            <Header />
            <main>
                <input placeholder="ENTER NAME" value={nameValue} onChange={handleInputChange} autoFocus />
                <button type="button" disabled={!nameValue || !code} onClick={handleJoin}>PLAY GAME</button>
                <button type="button" disabled={!nameValue} onClick={handleCreate}>CREATE ROOM</button>
            </main>
        </div>
    );
}

export default Home;
