import React, { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import Header from "../layout/Header";
import { createRoom, joinRoom } from "../api/Room";

const Home: React.FC = () => {
    const { code } = useParams<string>();

    const [name, setName] = useState<string>("");
    const [nameError, setNameError] = useState<string>("");
    const navigate = useNavigate();

    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setName(event.target.value);
    };

    const handleCreate = async() => {
        try {
            const data = await createRoom(name);
            sessionStorage.setItem("name", data.name);
            sessionStorage.setItem("gameCode", data.gameCode);
            navigate("/game/" + sessionStorage.getItem("gameCode"));
        } catch (err) {
            console.log(err);
        }
    };

    const handleJoin = async () => {
        setNameError("");
        try {
            const data = await joinRoom(name, code);
            sessionStorage.setItem("name", data.name);
            sessionStorage.setItem("gameCode", data.gameCode);
            navigate("/game/" + sessionStorage.getItem("gameCode"));
        } catch (err) {
            console.log(err);
            setNameError("Name must be unique.")
        }
    };

    return (
        <div className="home-container">
            <Header />
            <main>
                <input placeholder="ENTER NAME" value={name} onChange={handleInputChange} autoFocus />
                <div>{nameError}</div>
                <button type="button" disabled={!name || !code} onClick={handleJoin}>PLAY GAME</button>
                <button type="button" disabled={!name} onClick={handleCreate}>CREATE ROOM</button>
            </main>
        </div>
    );
}

export default Home;
