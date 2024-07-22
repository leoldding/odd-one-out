import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import Header from "../layout/Header";

const Join: React.FC = () => {
    const navigate = useNavigate();

    const [roomCode, setRoomCode] = useState<string>("");

    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setRoomCode(event.target.value);
    };

    const handleSubmit = () => {
        navigate("/game/" + roomCode );
    };

    return (
        <div className="join-container">
            <Header />
            <main>
                <form onSubmit={handleSubmit}>
                    <input placeholder="ENTER CODE" value={roomCode} onChange={handleInputChange} autoFocus />
                    <button type="submit" disabled={!roomCode}>PLAY</button>
                </form>
            </main>
        </div >
    );
}

export default Join;

