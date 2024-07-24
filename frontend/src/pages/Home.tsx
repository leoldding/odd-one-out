import React, { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import Header from "../layout/Header";
import { createRoom } from "../api/Room";

const Home: React.FC = () => {
    const { code } = useParams<string>();

    const [name, setName] = useState<string>("");
    const navigate = useNavigate();

    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setName(event.target.value);
    };

    const handleCreate = async() => {
        try {
            const data = await createRoom(name);
            localStorage.setItem("player", JSON.stringify(data));
        } catch (err) {
            console.log(err);
        }
    };

    const handleJoin = () => {
        navigate("/game/" + code);
    };

    return (
        <div className="home-container">
            <Header />
            <main>
                <input placeholder="ENTER NAME" value={name} onChange={handleInputChange} autoFocus />
                <button type="button" disabled={!name || !code} onClick={handleJoin}>PLAY GAME</button>
                <button type="button" disabled={!name} onClick={handleCreate}>CREATE ROOM</button>
            </main>
        </div>
    );
}

export default Home;
