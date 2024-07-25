import React, { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import Header from "../layout/Header";
import { createRoom, joinRoom } from "../api/Room";

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
            sessionStorage.setItem("name", data.name);
            sessionStorage.setItem("roomCode", data.roomCode);
            sessionStorage.setItem("leader", data.leader.toString());
            navigate("/game/" + sessionStorage.getItem("roomCode"));
        } catch (err) {
            console.log(err);
        }
    };

    const handleJoin = async () => {
        try {
            const data = await joinRoom(name, code);
            sessionStorage.setItem("name", data.name);
            sessionStorage.setItem("roomCode", data.roomCode);
            sessionStorage.setItem("leader", data.leader.toString());
            navigate("/game/" + sessionStorage.getItem("roomCode"));
        } catch (err) {
            console.log(err);
        }
        
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
